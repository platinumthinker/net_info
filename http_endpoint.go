package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type Page struct {
	Hosts []Host
}

type httpState struct {
	address  string
	req      chan chan []Host
	template *template.Template
}

func NewHttpState(address string, req chan chan []Host) *httpState {
	template := template.Must(template.ParseFiles("view/index.html"))
	return &httpState{
		req:      req,
		address:  address,
		template: template,
	}
}

func (h *httpState) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := make(chan []Host)
	h.req <- resp
	select {
	case hosts := <-resp:
		page := Page{
			Hosts: hosts,
		}

		err := h.template.Execute(w, &page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case <-time.After(1 * time.Second):
		http.Error(w, "Something bad", http.StatusInternalServerError)
	}
}

func httpRun(ctx context.Context, state *httpState) {
	http.Handle("/", state)
	fmt.Printf("Listen %s\n", state.address)
	http.ListenAndServe(state.address, nil)
}
