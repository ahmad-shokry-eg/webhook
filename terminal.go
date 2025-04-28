package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"time"
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
	} else {
		App.Secret = *secret
	}

	// get the IP for localhost
	var url string
	cmd := exec.Command("curl", "ifconfig.me")

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
	}
	defer stdoutPipe.Close()

	stdoutChan := make(chan string, 1)
	done := make(chan struct{})
	defer close(done)

	// Goroutine to read stdout
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			select {
			case stdoutChan <- scanner.Text():
			case <-done:
				return
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading output: %v", err)
		}
		close(stdoutChan)
	}()

	if err := cmd.Start(); err != nil {
		log.Fatal("failed to start command:", err)
	}

	// Add timeout for the HTTP request
	select {
	case baseURL, ok := <-stdoutChan:
		if ok {
			url = fmt.Sprintf("http://%s/webhook", baseURL)
		}
	case <-time.After(10 * time.Second):
		log.Fatal("timeout waiting for IP address")
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal("command failed:", err)
	}

	if url == "" {
		log.Fatal("failed to get IP address")
	}

	fmt.Printf("Use github secret:\n%s\n", defaultSecret)
	fmt.Println()
	fmt.Printf("Use github webhook url:\n%s\n", url)
}
