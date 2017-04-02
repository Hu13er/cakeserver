package main

import (
	"math/rand"
	"strings"
)

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
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[0:idx]
	}
	return ""
}
