package stack

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/JanSound/10Seconds-backend-go/beat"
	"github.com/gin-gonic/gin"
)

type BeatStackDTO struct {
	Key string
}

// @Schemes
// @Description stack beat
// @Tags stack
// @Param body body []BeatStackDTO true "병합하려는 key 들의 리스트를 입력해주세요."
// @Router /stack-beat [post]
func StackBeat(c *gin.Context) {
	var stacks []BeatStackDTO

	// Bind the JSON request body to the items slice
	if err := c.ShouldBindJSON(&stacks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := stack(stacks)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{})
}

type StackDTO struct {
	Filename string
	Type     string
}

func stack(stacks []BeatStackDTO) error {
	url := os.Getenv("core_host") + "/beats/stack"
	var data []map[string]interface{}
	for _, value := range stacks {
		replacedKey := strings.Replace(value.Key, "beat/", "", -1)
		substrings := strings.Split(replacedKey, "/")
		// TODO(@shmoon): 정의한 인터페이스로 수정
		filename := strings.Replace(substrings[1], ".m4a", "", -1)
		data = append(data, map[string]interface{}{
			"filename": filename,
			"type":     substrings[0],
		})
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return errors.New("error marshaling json")
	}
	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return errors.New("error making request")
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(body))
	}

	stackFilename, _ := ioutil.ReadAll(resp.Body)

	unquotedFileName, err := strconv.Unquote(string(stackFilename))
	if err != nil {
		return err
	}
	filename := "beat/stack/" + unquotedFileName + ".m4a"
	beat.CreateBeat(string(filename), "stack")
	defer resp.Body.Close()
	return nil
}
