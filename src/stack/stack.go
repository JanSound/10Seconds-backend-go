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

func getStackResultFilename(stacks []BeatStackDTO) string {

	return ""
}

type StackDTO struct {
	Filename string
	Type     string
}

func stack(stacks []BeatStackDTO) {
	url := os.Getenv("core_host") + "/beats/stack"
	var data []StackDTO
	for _, value := range stacks {

		replacedKey := strings.Replace(value.Key, "beats/", "", -1)
		substrings := strings.Split(replacedKey, "/")
		data = append(data, StackDTO{
			Filename: substrings[0],
			Type:     substrings[1],
		})
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

func uploadStackBeat(input string) {

}
