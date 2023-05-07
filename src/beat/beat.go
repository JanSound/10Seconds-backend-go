package beat

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Beat struct {
	gorm.Model
	id    uint `gorm:"primaryKey"`
	title string
}

type BeatCreateDTO struct {
	Key      string
	BeatType string
}

// @Schemes
// @Description create beat
// @Tags beats
// @Param body body BeatCreateDTO true "생성하려는 파일의 key와 beatType 를 넣어주세요."
// @Router /beats [post]
func PostBeat(c *gin.Context) {
	var payload BeatCreateDTO
	c.ShouldBind(&payload)
	CreateBeat(payload.Key, payload.BeatType)
	c.JSON(201, gin.H{})
}

func GetBeatList(c *gin.Context) {
	beats, _ := GetAllBeats()
	c.JSON(200, gin.H{
		"message": beats,
	})
}

// @Schemes
// @Description create beat
// @Tags beats
// @Router /beats [get]
func GetBeatDetail(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

type User struct {
	ID int `uri:"id" binding:"required,uuid"`
}

// @Schemes
// @Description delete beat`
// @Tags beats
// @Param        id   path      int  true  "user id"
// @Router /beats/{id} [delete]
func DeleteBeat(c *gin.Context) {
	user_id := c.Param("id")
	DeleteBeatById(user_id)
	c.JSON(200, gin.H{})
}
