package convert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

	convert(key)
	presignedUrlList := getPresignedUrl(key)

	beatList := []ConvertedBeatDTO{
		{Key: key, BeatType: "base", PresignedUrl: presignedUrlList[0]},
		{Key: key, BeatType: "piano", PresignedUrl: presignedUrlList[1]},
		{Key: key, BeatType: "drum", PresignedUrl: presignedUrlList[2]},
	}

	return beatList
}

type ConvertDTO struct {
	filename string
}

func convert(key string) {
	url := os.Getenv("core_host") + "/beats/convert"
	data := ConvertDTO{
		filename: key,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Create a request with the JSON payload
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	fmt.Println("Response status code:", resp.StatusCode)

}

func getPresignedUrl(key string) []string {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile("tenseconds"),
		config.WithRegion("ap-northeast-2"),
	)
	if err != nil {
		fmt.Println(err)
	}
	client := s3.NewFromConfig(cfg)
	bucket := os.Getenv("aws_s3_bucket")

	ret := []string{}
	myList := []string{"base", "piano", "drum"}
	for _, value := range myList {
		url := value + "/" + key
		presignClient := s3.NewPresignClient(client)
		presigner := beat.Presigner{PresignClient: presignClient}

		presignedGetRequest, err := presigner.GetObject(bucket, url, 60)
		if err != nil {
			fmt.Println(err)
		}
		ret = append(ret, presignedGetRequest.URL)

	}
	return ret
}
