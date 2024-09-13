package cloud

import (
	"context"
	"dhanushs3366/my-portfolio/utils"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsClient struct {
	client *s3.Client
}

func NewS3Client() *AwsClient {
	bucketRegion := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY")
	accessSecretKey := os.Getenv("AWS_ACCESS_SECRET_KEY")

	options := s3.Options{
		Region: bucketRegion,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			accessKey,
			accessSecretKey,
			"",
		)),
	}

	client := s3.New(options)

	return &AwsClient{client: client}
}

// obj key: <file-type>/<file-name>/uuid

func GenerateKeyForUpload(filename string) (string, error) {
	fileType, err := utils.GetFileType(filename)
	if err != nil {
		return "", err
	}
	uuidKey := utils.GetNewUUID()

	return fmt.Sprintf("%s/%s/%s", fileType, filename, uuidKey), nil
}

// pass in key as param call key function outside cuz we also need to add the key in db
func (c *AwsClient) S3Upload(file io.Reader, key string) error {
	bucketName := os.Getenv("AWS_BUCKET")
	inputOptions := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
	}
	_, err := c.client.PutObject(context.TODO(), inputOptions)
	if err != nil {
		log.Printf("Error uploading to S3: %v", err)
	}
	return err
}
