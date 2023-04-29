package beat

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"

	// swagger embed files
	"fmt"
	"math/rand"
	"os"
	"time"
)

func generateUniqueFilename() string {
	now := time.Now()
	rand.Seed(now.UnixNano())
	randomNumber := rand.Intn(1000000000)

	return fmt.Sprintf("%s-%d", now.Format("20060102150405"), randomNumber)
}

func generatePutObjectPresignedURL(sess *session.Session, content_type string, objectKey string) (string, error) {
	svc := s3.New(sess)
	expiration := 180 * time.Minute // 서비스화 할 경우 조정할 예정

	bucket := os.Getenv("aws_s3_bucket")

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(objectKey),
		ContentType: aws.String(content_type),
		// ContentType: aws.String("audio/x-m4a"), // or "audio/x-m4a"
	})
	presignedURL, err := req.Presign(expiration)
	return presignedURL, err
}

// @Schemes
// @Description create presigned url to upload beats (m4a audio file)
// @Tags beats
// @Router /beats/generate-presigned-url [post]
func GeneratePutObjectPresignedURL(c *gin.Context) {
	region := os.Getenv("aws_s3_region")
	file_root := os.Getenv("aws_s3_file_root")
	objectKey := file_root + generateUniqueFilename() + ".m4a"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("aws_access_key_id"), os.Getenv("aws_secret_access_key"), "",
		),
	})
	if err != nil {
		c.JSON(500, gin.H{
			"error": fmt.Sprintf("Failed to create session: %v", err),
		})
		return
	}
	presignedURL, err := generatePutObjectPresignedURL(sess, "audio/x-m4a", objectKey)

	if err != nil {
		c.JSON(500, gin.H{
			"error": fmt.Sprintf("Failed to generate presigned URL: %v", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"presigned_url": presignedURL,
	})
}
