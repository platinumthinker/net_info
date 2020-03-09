package main

import (
	"context"
	"log"
	"sort"
	"time"

	routeros "gopkg.in/routeros.v2"
)

type Host struct {
	IP     string
	Label  string
	Subnet string
}

type state struct {
	hosts []Host
	req   chan chan []Host
}

func NewState(staticHosts map[string]string) *state {
	hosts := make([]Host, 0, len(staticHosts))
	for host, ip := range staticHosts {
		hosts = append(hosts, Host{
			IP:    ip,
			Label: host,
		})
	}

	sort.Slice(hosts, func(i, j int) bool {
		return hosts[i].IP > hosts[j].IP
	})

	return &state{
		hosts: hosts,
		req:   make(chan chan []Host),
	}
}

func (state *state) Run(ctx context.Context, cfg MikroticConfig) {
	c, err := routeros.Dial(cfg.Address, cfg.User, cfg.Password)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case reply := <-state.req:
			reply <- state.hosts
		case <-time.After(10 * time.Second):
			state.updateHosts(c)
		case <-ctx.Done():
			return
		}
	}
}

func (state *state) updateHosts(c *routeros.Client) {
	properties := "address,mac-address,server,last-seen,active-address,host-name,bloked"
	leases, err := c.Run("/ip/dhcp-server/lease/print", "=.proplist="+properties)
	if err != nil {
		log.Fatal(err)
	}

	for _, re := range leases.Re {
		state.updateHost(re.Map["address"], re.Map["host-name"])
	}
}

func (state *state) updateHost(ip, host string) {
	if len(host) == 0 {
		host = ip
	}

	for i, h := range state.hosts {
		if h.IP == ip {
			if h.Label == host {
				return
			}

			state.hosts[i].Label = host
			return
		}
	}

	state.hosts = append(state.hosts, Host{
		IP:    ip,
		Label: host,
	})
}
