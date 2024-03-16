package lib

import (
	"fmt"
	"log"
)

func Print(f string, args ...interface{}) {
	if Cfg.Quiet {
		return
	}
	log.Printf(">> "+Decode(f), args...)
}

func Sprint(f string, args ...interface{}) string {
	return fmt.Sprintf(Decode(f), args...)
}
