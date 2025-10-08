package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"magicchat/pkg/config"
)

type StorageClient struct {
	s3Client *s3.S3
	bucket   string
	provider string
}

var Client *StorageClient

func InitStorage(cfg *config.Config) (*StorageClient, error) {
	var awsConfig *aws.Config

	if cfg.Storage.Provider == "minio" {
		awsConfig = &aws.Config{
			Credentials:      credentials.NewStaticCredentials(cfg.Storage.MinioAccessKey, cfg.Storage.MinioSecretKey, ""),
			Endpoint:         aws.String(cfg.Storage.MinioEndpoint),
			Region:           aws.String(cfg.Storage.AWSRegion),
			DisableSSL:       aws.Bool(!cfg.Storage.MinioUseSSL),
			S3ForcePathStyle: aws.Bool(true),
		}
	} else {
		awsConfig = &aws.Config{
			Credentials: credentials.NewStaticCredentials(cfg.Storage.AWSAccessKey, cfg.Storage.AWSSecretKey, ""),
			Region:      aws.String(cfg.Storage.AWSRegion),
		}
		if cfg.Storage.S3Endpoint != "" {
			awsConfig.Endpoint = aws.String(cfg.Storage.S3Endpoint)
		}
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	Client = &StorageClient{
		s3Client: s3.New(sess),
		bucket:   cfg.Storage.S3Bucket,
		provider: cfg.Storage.Provider,
	}

	// Create bucket if it doesn't exist (for MinIO)
	if cfg.Storage.Provider == "minio" {
		err = Client.createBucketIfNotExists()
		if err != nil {
			log.Printf("Warning: Could not create bucket: %v", err)
		}
	}

	log.Printf("Storage client initialized with provider: %s", cfg.Storage.Provider)
	return Client, nil
}

func (c *StorageClient) createBucketIfNotExists() error {
	ctx := context.Background()
	_, err := c.s3Client.HeadBucketWithContext(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(c.bucket),
	})
	if err == nil {
		return nil // Bucket exists
	}

	_, err = c.s3Client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(c.bucket),
	})
	return err
}

func (c *StorageClient) UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(filename)
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Read file content
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, file)
	if err != nil {
		return "", err
	}

	// Upload to S3/MinIO
	_, err = c.s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(uniqueFilename),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(contentType),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}

	// Generate URL
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", c.bucket, uniqueFilename)
	if c.provider == "minio" {
		// For MinIO, construct URL differently
		url = fmt.Sprintf("http://%s/%s/%s", c.s3Client.Endpoint, c.bucket, uniqueFilename)
	}

	return url, nil
}

func (c *StorageClient) DeleteFile(ctx context.Context, fileURL string) error {
	// Extract key from URL
	key := filepath.Base(fileURL)

	_, err := c.s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (c *StorageClient) GetPresignedURL(ctx context.Context, key string, expirationMinutes int64) (string, error) {
	req, _ := c.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})

	url, err := req.Presign(time.Duration(expirationMinutes) * time.Minute)
	return url, err
}
