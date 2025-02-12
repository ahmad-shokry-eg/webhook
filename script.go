package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func execScript(logFile *os.File) {
	defer logFile.Close()

	logFile.WriteString("executing pipeline..\n\n")

	// Create a command to execute the script
	cmd := exec.Command("sh", "-C", fmt.Sprintf("./%s", App.ScriptFile))

	// Create pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		logFile.WriteString(fmt.Sprintf("Error creating stdout pipe: %v\n", err))
		return
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		logFile.WriteString(fmt.Sprintf("Error creating stderr pipe: %v\n", err))
		return
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		logFile.WriteString(fmt.Sprintf("Error starting command: %v\n", err))
		return
	}

	// Channels to receive lines from stdout and stderr
	stdoutChan := make(chan string)
	stderrChan := make(chan string)

	// Goroutines to read stdout and stderr line by line
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			stdoutChan <- scanner.Text()
		}
		close(stdoutChan)
	}()
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			stderrChan <- scanner.Text()
		}
		close(stderrChan)
	}()

	// Process lines from stdout and stderr
	for {
		select {
		case line, ok := <-stdoutChan:
			if !ok {
				stdoutChan = nil // Channel closed, stop reading
			} else {
				logFile.WriteString(fmt.Sprintf("%s: %s\n", time.Now().Format("2006-01-02-Mon-03-04-05-PM"), line))
			}
		case line, ok := <-stderrChan:
			if !ok {
				stderrChan = nil // Channel closed, stop reading
			} else {
				logFile.WriteString(fmt.Sprintf("%s: %s\n", time.Now().Format("2006-01-02-Mon-03-04-05-PM"), line))
			}
		}

		// Exit when both channels are closed
		if stdoutChan == nil && stderrChan == nil {
			break
		}
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		logFile.WriteString(fmt.Sprintf("Command finished with error: %v\n", err))
		return
	}

	logFile.WriteString("\nPipeline executed successfully.\n")
}
