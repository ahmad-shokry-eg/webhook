package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

var App struct {
	Port, Secret, ScriptFile string
}

func main() {
	// grab terminal commands
	processTerminal()

	// create the logs directory
	os.Mkdir("../logs", 0755)

	// create gin web server
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	// add the webhook route
	server.POST("/webhook", webhook)

	// wait for shut down signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// run the server
	go server.Run(App.Port)

	// keep it running until shutdown
	<-quit
	log.Println("Shutting down webhook...")
}
