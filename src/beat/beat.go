package beat

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Beat struct {
	gorm.Model
	id    uint `gorm:"primaryKey"`
	title string
}

// @Schemes
// @Description create beat
// @Tags beats
// @Router /beats [post]
func PostBeat(c *gin.Context) {
	CreateBeat("test", "test2")
	c.JSON(200, gin.H{
		"message": "pong",
	})
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

// @Schemes
// @Description create beat
// @Tags beats
// @Router /beats [get]
func DeleteBeat(c *gin.Context) {
	fmt.Println(c.Param("beat_id"))
	c.JSON(200, gin.H{
		"message": "pong",
	})

}
