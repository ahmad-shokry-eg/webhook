package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func webhook(c *gin.Context) {
	// create logs file
	logFile, err := os.Create(fmt.Sprintf("../logs/build_%s", time.Now().Format("2006-01-02-Mon-03-04-05-PM")))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "Failed to create log file")
	}

	// get request header
	sig := c.GetHeader("X-Hub-Signature-256")

	// get request body
	body, _ := io.ReadAll(c.Request.Body)

	// validate signature
	if !verifySignature(body, sig, App.Secret) {
		logFile.WriteString("invalid github signature")
		logFile.Close()
		c.AbortWithStatusJSON(http.StatusUnauthorized, "Invalid signature")
		return
	}

	// Check event type
	eventType := c.GetHeader("X-GitHub-Event")
	if eventType != "push" && eventType != "ping" {
		logFile.WriteString("Event type not supported")
		logFile.Close()
		c.AbortWithStatusJSON(http.StatusUnauthorized, "Event type not supported")
		return
	}

	// execute build script
	go execScript(logFile)

	// return success
	c.JSON(http.StatusOK, "executing pipeline..")
}
