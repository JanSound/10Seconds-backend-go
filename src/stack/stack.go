package stack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

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
	stack(stacks)
	c.JSON(200, gin.H{})
}

type StackDTO struct {
	Filename string
	Type     string
}

func stack(stacks []BeatStackDTO) {
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
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
}

func uploadStackBeat(input string) {

}
