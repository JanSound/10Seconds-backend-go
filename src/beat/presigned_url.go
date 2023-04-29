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

func (presigner Presigner) GetObject(
	bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
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

// func generatePutObjectPresignedURL(sess *session.Session, content_type string, objectKey string) (string, error) {
// 	svc := s3.New(sess)
// 	expiration := 180 * time.Minute // 서비스화 할 경우 조정할 예정

// 	bucket := os.Getenv("aws_s3_bucket")

// 	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
// 		Bucket:      aws.String(bucket),
// 		Key:         aws.String(objectKey),
// 		ContentType: aws.String(content_type),
// 		// ContentType: aws.String("audio/x-m4a"), // or "audio/x-m4a"
// 	})
// 	presignedURL, err := req.Presign(expiration)
// 	return presignedURL, err
// }

// @Schemes
// @Description create presigned url to upload beats (m4a audio file)
// @Tags beats
// @Router /beats/generate-presigned-url [post]
func GeneratePutObjectPresignedURL(c *gin.Context) {

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile("tenseconds"),
		config.WithRegion("ap-northeast-2"),
	)

	if err != nil {
		return
	}
	client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(client)
	presigner := Presigner{PresignClient: presignClient}
	bucket := os.Getenv("aws_s3_bucket")
	file_root := os.Getenv("aws_s3_file_root")
	objectKey := file_root + generateUniqueFilename() + ".m4a"
	presignedPutRequest, err := presigner.PutObject(bucket, objectKey, 60)
	fmt.Printf(presignedPutRequest.URL)
	if err != nil {
		panic(err)
	}

	presignedURL := presignedPutRequest.URL
	c.JSON(200, gin.H{
		"presigned_url": presignedURL,
	})
}
