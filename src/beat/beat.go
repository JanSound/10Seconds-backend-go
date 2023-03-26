package beat

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"

	"fmt"
	"os"
	"time"
)

func GeneratePresignedURL(c *gin.Context) {
	region := "ap-northeast-2"
	bucket := "tenseconds"
	objectKey := "tenseconds-demo"

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

	svc := s3.New(sess)
	expiration := 180 * time.Minute // 서비스화 할 경우 조정할 예정

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})
	presignedURL, err := req.Presign(expiration)

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
