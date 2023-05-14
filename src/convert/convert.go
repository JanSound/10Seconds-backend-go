package convert

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/JanSound/10Seconds-backend-go/beat"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

type BeatConvertDTO struct {
	Key string
}

type ConvertedBeatDTO struct {
	Key          string
	BeatType     string
	PresignedUrl string
}

func ConvertBeat(c *gin.Context) {
	var payload BeatConvertDTO
	c.ShouldBind(&payload)
	beats := convertBeat(payload.Key)
	c.JSON(200, beats)
}

// @Schemes
// @Description convert beat
// @Tags convert
// @Param body body BeatConvertDTO true "변환하려는 파일의 key 를 입력해주세요."
// @Router /convert-beat [post]
func convertBeat(key string) []ConvertedBeatDTO {
	// TODO(@shmoon): 음성 처리 서버로부터 3개의 응답 받아오기

	tempPresignedUrl := getPresignedUrl(key)

	beatList := []ConvertedBeatDTO{
		{Key: key, BeatType: "guitar", PresignedUrl: tempPresignedUrl},
		{Key: key, BeatType: "base", PresignedUrl: tempPresignedUrl},
		{Key: key, BeatType: "drum", PresignedUrl: tempPresignedUrl},
	}

	return beatList
}

func getPresignedUrl(key string) string {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile("tenseconds"),
		config.WithRegion("ap-northeast-2"),
	)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	client := s3.NewFromConfig(cfg)
	bucket := os.Getenv("aws_s3_bucket")

	presignClient := s3.NewPresignClient(client)
	presigner := beat.Presigner{PresignClient: presignClient}

	presignedGetRequest, err := presigner.GetObject(bucket, key, 60)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	presignedURL := presignedGetRequest.URL
	return presignedURL
}
