package convert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

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

// @Schemes
// @Description convert beat
// @Tags convert
// @Param body body BeatConvertDTO true "변환하려는 파일의 key 를 입력해주세요."
// @Router /convert-beat [post]
func ConvertBeat(c *gin.Context) {
	var payload BeatConvertDTO
	c.ShouldBind(&payload)
	beats := convertBeat(payload.Key)
	c.JSON(200, beats)
}

func convertBeat(key string) []ConvertedBeatDTO {
	// TODO(@shmoon): 음성 처리 서버로부터 3개의 응답 받아오기

	convert(key)
	presignedUrlList := getPresignedUrlList(key)
	replacedKey := strings.Replace(key, "voice/", "", -1)
	beatList := []ConvertedBeatDTO{
		{Key: "beat/bass/" + replacedKey, BeatType: "bass", PresignedUrl: presignedUrlList[0]},
		{Key: "beat/piano/" + replacedKey, BeatType: "piano", PresignedUrl: presignedUrlList[1]},
		{Key: "beat/drum/" + replacedKey, BeatType: "drum", PresignedUrl: presignedUrlList[2]},
	}

	return beatList
}

type ConvertDTO struct {
	filename string
}

func convert(key string) {
	url := os.Getenv("core_host") + "/beats/convert"
	var data map[string]string
	data = map[string]string{
		"filename": key,
	}
	// data := ConvertDTO{
	// 	filename: key,
	// }

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

func getPresignedUrlList(key string) []string {
	// 임시 로직
	replacedKey := strings.Replace(key, "voice/", "", -1)
	replacedKey = strings.Replace(replacedKey, ".m4a", "", -1)

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
	myList := []string{"bass", "piano", "drum"}
	for _, value := range myList {
		url := "beat/" + value + "/" + replacedKey + ".m4a"
		presignedClient := s3.NewPresignClient(client)
		presigner := beat.Presigner{PresignClient: presignedClient}

		presignedGetRequest, err := presigner.GetObject(bucket, url, 60)
		if err != nil {
			fmt.Println(err)
		}
		ret = append(ret, presignedGetRequest.URL)

	}
	return ret
}
