package beat

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	// swagger embed files
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Presigner struct {
	PresignClient *s3.PresignClient
}

type Payload struct {
	Key string
}

func (presigner Presigner) GetObject(
	bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		fmt.Println(err)
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}

func (presigner Presigner) PutObject(
	bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.PresignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		fmt.Println(err)
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}

func generateUniqueFilename() string {
	now := time.Now()
	rand.Seed(now.UnixNano())
	randomNumber := rand.Intn(1000000000)

	return fmt.Sprintf("%s-%d", now.Format("20060102150405"), randomNumber)
}

func GenerateGetObjectPresignedUrl(c *gin.Context) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile("tenseconds"),
		config.WithRegion("ap-northeast-2"),
	)

	if err != nil {
		fmt.Println(err)
		return
	}
	var payload Payload
	c.ShouldBind(&payload)
	client := s3.NewFromConfig(cfg)
	bucket := os.Getenv("aws_s3_bucket")

	presignClient := s3.NewPresignClient(client)
	presigner := Presigner{PresignClient: presignClient}

	fileKey := payload.Key
	presignedGetRequest, err := presigner.GetObject(bucket, fileKey, 60)
	presignedURL := presignedGetRequest.URL
	c.JSON(200, gin.H{
		"presigned_url": presignedURL,
	})
}

// @Schemes
// @Description create presigned url to upload beats (m4a audio file)
// @Tags beats
// @Router /beats/presigned-url/put [post]
func GeneratePutObjectPresignedURL(c *gin.Context) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile("tenseconds"),
		config.WithRegion("ap-northeast-2"),
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	bucket := os.Getenv("aws_s3_bucket")
	file_root := os.Getenv("aws_s3_file_root")
	objectKey := file_root + generateUniqueFilename() + ".m4a"

	client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(client)
	presigner := Presigner{PresignClient: presignClient}
	presignedPutRequest, err := presigner.PutObject(bucket, objectKey, 60)

	if err != nil {
		fmt.Println(err)
	}

	presignedURL := presignedPutRequest.URL
	c.JSON(200, gin.H{
		"presigned_url": presignedURL,
		"key":           objectKey,
	})
}
