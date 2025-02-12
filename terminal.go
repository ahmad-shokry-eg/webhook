package main

import (
	"flag"
	"fmt"
)

func processTerminal() {
	defaultSecret := generateSecret()

	// define flags
	flag.StringVar(&App.Port, "p", ":1099", "specify webhook port")
	flag.StringVar(&App.ScriptFile, "e", "script.sh", "specify pipeline file")
	secret := flag.String("s", "", "specify secret text ( empty to auto-generate new secret)")

	// parse flags
	flag.Parse()

	if *secret == "" {
		App.Secret = defaultSecret
		fmt.Printf("use github secret:\n%s\n", defaultSecret)
	} else {
		App.Secret = *secret
	}
}
