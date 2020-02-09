package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/routeros.v2"
)

var (
	address    = flag.String("address", "192.168.0.1:8728", "Address")
	username   = flag.String("username", "admin", "Username")
	password   = flag.String("password", "admin", "Password")
	properties = flag.String("properties", "address,mac-address,server,last-seen,active-address,host-name,bloked", "Properties")
	interval   = flag.Duration("interval", 1*time.Second, "Interval")
)

func main() {
	flag.Parse()

	c, err := routeros.Dial(*address, *username, *password)
	if err != nil {
		log.Fatal(err)
	}

	for {
		reply, err := c.Run("/ip/dhcp-server/lease/print", "=.proplist="+*properties)
		if err != nil {
			log.Fatal(err)
		}

		for _, re := range reply.Re {
			for _, p := range strings.Split(*properties, ",") {
				fmt.Print(re.Map[p], "\t")
			}
			fmt.Print("\n")
		}
		fmt.Print("\n")

		time.Sleep(*interval)
	}
}
