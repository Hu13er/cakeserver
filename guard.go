package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
)

type handlerFunc func(http.ResponseWriter, *http.Request)

type guard struct {
	secret string
	token  [sha256.Size]byte
	tried  map[string]int
}

func newGuard(secret string) *guard {
	return &guard{secret: secret, tried: make(map[string]int)}
}

func (g *guard) newToken() []byte {

	random := randomSeq()
	raw := append(random, '.')
	raw = append(raw)

	g.token = sha256.Sum256(raw)

	return random
}

func (g *guard) isForbidden(ip string) bool {
	return g.tried[ip] > 3
}

func (g *guard) guardHandler(handler handlerFunc) handlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		context.Set(r, "super", false)

		ip := getIP(r.RemoteAddr)
		if g.tried[ip] > 3 {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		token := r.Header.Get("token")
		if token != "" {
			mustBe := fmt.Sprintf("% x", g.token[:])
			if token != mustBe {
				g.tried[ip]++
				w.WriteHeader(http.StatusForbidden)
				return
			}

			g.newToken()
			g.tried[ip] = 0
			context.Set(r, "super", true)
		}

		handler(w, r)
	}
}

func (g *guard) generatorHandler(w http.ResponseWriter, r *http.Request) {

	ip := getIP(r.RemoteAddr)

	if g.isForbidden(ip) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	g.tried[ip] = 0
	random := g.newToken()

	w.Write(random)
}
