package main

import (
	"log"

	"github.com/JanSound/10Seconds-backend-go/beat"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// .env 파일을 읽어서 환경변수를 설정
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env5 file")
	}

	r := gin.Default()
	r.GET("/generate-presigned-url", beat.GeneratePresignedURL)

	r.Run(":8001")
}
