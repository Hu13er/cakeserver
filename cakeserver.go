package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type handlerFunc func(http.ResponseWriter, *http.Request)

var (
	addr      string
	secret    = []byte{}
	generated [sha256.Size]byte
	tried     = map[string]int{}
)

func generate() []byte {
	rand := randomSeq()

	row := append(rand, '.')
	row = append(row, secret...)
	generated = sha256.Sum256(row)

	return rand
}

func guard(f handlerFunc) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("key")

		if tried[getIP(r.RemoteAddr)] > 3 {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		value := fmt.Sprintf("% x", generated[:])

		if key != value {
			tried[getIP(r.RemoteAddr)]++
			w.WriteHeader(http.StatusForbidden)
			return
		}

		generate()

		tried[getIP(r.RemoteAddr)] = 0

		f(w, r)

	}
}

func generateHandler(w http.ResponseWriter, r *http.Request) {

	if tried[getIP(r.RemoteAddr)] > 3 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	tried[getIP(r.RemoteAddr)]++

	rand := generate()
	w.Write(rand)
}

func makeHandler(w http.ResponseWriter, r *http.Request) {
	command, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := makefile(
		strings.Split(string(command), " ")...,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(result))
}

func serve() {
	log.Println("Start serving on", addr)
	mux := http.NewServeMux()

	mux.HandleFunc("/generate/", generateHandler)
	mux.HandleFunc("/make/", guard(makeHandler))

	log.Fatalln(
		http.ListenAndServe(addr, mux),
	)
}

func main() {
	addr = os.Getenv("CAKE_ADDR")
	if addr == "" {
		addr = ":2128"
	}

	sec := os.Getenv("CAKE_SECRET")
	if sec == "" {
		log.Fatalln("CAKE_SECRET not presented.")
	}
	secret = []byte(sec)

	serve()
}

func randomSeq() []byte {
	var (
		charset = "abcdefghijklmnopqrstuvwxyz1234567890!@#$%^&*"
		size    = 32
		random  = []byte{}
	)

	for i := 0; i < size; i++ {
		random = append(random, charset[rand.Intn(len(charset))])
	}

	return random
}

func isSliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for k := range a {
		if a[k] != b[k] {
			return false
		}
	}

	return true
}

func getIP(addr string) string {
	return strings.Split(addr, ":")[0]
}

func makefile(args ...string) (string, error) {
	command := exec.Command("make", args...)
	result, err := command.Output()
	return string(result), err
}
