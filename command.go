package main

import (
	"os/exec"
	"strings"
)

var Commands []Command

func GetCommand(name string) *Command {
	for k, v := range Commands {
		if v.Name == name {
			return &Commands[k]
		}
	}
	return nil
}

type Command struct {
	Name    string   `yaml:"name"`
	Trusted []string `yaml:"trusted"`
	Script  string   `yaml:"script"`
}

func (c *Command) Run() error {
	script := c.Script
	script = strings.Replace(script, "\n", "; ", -1)
	cmd := exec.Command("/bin/sh", "-c", script)
	return cmd.Run()
}

func (c *Command) IsTrusted(ip string) bool {
	for _, v := range c.Trusted {
		if v == ip {
			return true
		}
	}
	return false
}
