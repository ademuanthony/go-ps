package main

import (
	"fmt"
)

import "github.com/ademuanthony/ps"

func main() {
	proc, err := ps.ProcessByName("dcrwallet")
	if err == nil {
		ports, err := ps.AssociatedPorts(proc.Pid())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(ports)
	}
	fmt.Println(err)
}
