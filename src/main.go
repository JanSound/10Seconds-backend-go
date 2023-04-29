package main

import (
	"log"

	"github.com/JanSound/10Seconds-backend-go/beat"
	docs "github.com/JanSound/10Seconds-backend-go/docs"
	swaggerFiles "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// read .env file
	err := godotenv.Load()

	docs.SwaggerInfo.BasePath = "/api/v1"

	if err != nil {
		log.Fatal("Error loading .env5 file")
	}

	r := gin.Default()

	v1 := r.Group("api/v1")
	{
		eg := v1.Group("beats")
		{
			eg.POST("generate-presigned-url", beat.GeneratePutObjectPresignedURL)
			eg.POST("", beat.PostBeat)
			eg.GET("", beat.GetBeatList)
			eg.GET(":beat_id", beat.GetBeatDetail)
			eg.DELETE(":beat_id", beat.DeleteBeat)
		}
	}

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8001")
}
