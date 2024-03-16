package main

import "fmt"

func Print(f string, args ...interface{}) {
	fmt.Printf(">> "+f+"\n", args...)
}
