package convert

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	beats, err := convertBeat(payload.Key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(200, beats)
}

func convertBeat(key string) ([]ConvertedBeatDTO, error) {
	// TODO(@shmoon): 음성 처리 서버로부터 3개의 응답 받아오기

	err := convert(key)
	if err != nil {
		return nil, err
	}
	presignedUrlList := getPresignedUrlList(key)
	replacedKey := strings.Replace(key, "voice/", "", -1)
	beatList := []ConvertedBeatDTO{
		{Key: "beat/bass/" + replacedKey, BeatType: "bass", PresignedUrl: presignedUrlList[0]},
		{Key: "beat/piano/" + replacedKey, BeatType: "piano", PresignedUrl: presignedUrlList[1]},
		{Key: "beat/drum/" + replacedKey, BeatType: "drum", PresignedUrl: presignedUrlList[2]},
	}
	return beatList, nil
}

func convert(key string) error {
	url := os.Getenv("core_host") + "/beats/convert"

	key = strings.Replace(key, "voice/", "", -1)
	key = strings.Replace(key, ".m4a", "", -1)

	var data map[string]string = map[string]string{
		"filename": key,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Create a request with the JSON payload
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body))
	}

	return nil
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
