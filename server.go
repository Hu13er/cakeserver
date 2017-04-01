package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var grd *guard

func serve(addr, secret string) {
	log.Println("Start serving on", addr)

	grd = newGuard(secret)
	mux := http.NewServeMux()

	mux.HandleFunc("/generate/", grd.generatorHandler)
	mux.HandleFunc("/command/", grd.guardHandler(cmdHandler))

	log.Fatalln(
		http.ListenAndServe(addr, mux),
	)
}

func cmdHandler(w http.ResponseWriter, r *http.Request) {
	command, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := runCommand(
		strings.Split(string(command), " ")...,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(result))
}
