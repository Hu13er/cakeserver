package main

import (
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

var grd *guard

func serve(addr, secret string) {
	log.Println("Start serving on", addr)

	grd = newGuard(secret)
	router := mux.NewRouter()

	router.HandleFunc("/generate/", grd.generatorHandler).Methods("GET")
	router.HandleFunc("/command/{command}/", grd.guardHandler(cmdHandler)).Methods("GET")

	log.Fatalln(
		http.ListenAndServe(addr, router),
	)
}

func cmdHandler(w http.ResponseWriter, r *http.Request) {
	name, ok := mux.Vars(r)["command"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	command := GetCommand(string(name))
	if command == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ip := getIP(r.RemoteAddr)
	if isSuper := context.Get(r, "super").(bool); !isSuper && !command.IsTrusted(ip) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := command.Run(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
