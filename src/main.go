package main

import (
	"github.com/JanSound/10Seconds-backend-go/beat"
	"github.com/gin-gonic/gin"
	// "log"
)

func main() {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	r := gin.Default()
	r.GET("/generate-presigned-url", beat.GeneratePresignedURL)

	r.Run(":8001")
}
