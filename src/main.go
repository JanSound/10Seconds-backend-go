package main

import (
	"fmt"

	"github.com/JanSound/10Seconds-backend-go/beat"
	"github.com/JanSound/10Seconds-backend-go/convert"
	docs "github.com/JanSound/10Seconds-backend-go/docs"
	"github.com/JanSound/10Seconds-backend-go/stack"
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
		fmt.Println(err)
	}

	r := gin.Default()

	v1 := r.Group("api/v1")
	{
		beats := v1.Group("beats")
		{
			presigned_url := beats.Group("presigned-url")
			{
				presigned_url.POST("put/", beat.GeneratePutObjectPresignedURL)
				// presigned_url.POST("get/", beat.GenerateGetObjectPresignedUrl)
			}

			beats.POST("", beat.PostBeat)
			beats.GET("", beat.GetBeatList)
			beats.GET(":beat_id", beat.GetBeatDetail)
			beats.DELETE(":id", beat.DeleteBeat)
			beats.DELETE("all", beat.DeleteAllBeat)
		}

		convertBeat := v1.Group("convert-beat")
		{
			convertBeat.POST("", convert.ConvertBeat)
		}

		stackBeat := v1.Group("stack-beat")
		{
			stackBeat.POST("", stack.StackBeat)
		}
	}

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8001")
}
