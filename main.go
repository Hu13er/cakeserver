package main

import (
	"log"
	"os"
)

func main() {
	addr := os.Getenv("CAKE_ADDR")
	if addr == "" {
		addr = ":2128"
	}

	secret := os.Getenv("CAKE_SECRET")
	if secret == "" {
		log.Fatalln("CAKE_SECRET not presented.")
	}

	serve(addr, secret)
}
