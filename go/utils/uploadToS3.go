package utils

import (
	"crypto/sha256"
	"time"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GenerateKeyFromFilename(filename string) string {
	// get file extension and base name
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	// generate unique hash from extension, base name, and current timestamp (nanoseconds so collision is extremely unlikely)
	hash := sha256.New()
    hash.Write([]byte(fmt.Sprintf("%s_%d", base, time.Now().UnixNano())))
	hashStr := fmt.Sprintf("%x", hash.Sum(nil))[:12]	// first 12 chars; reduce length for readability

	// clean up filename and combine with hash
	cleanName := strings.ReplaceAll(base, " ", "_")
	cleanName = strings.ReplaceAll(cleanName, "/", "_")

	// return final key with cleanName_hash.ext
    return fmt.Sprintf("%s_%s%s", cleanName, hashStr, ext)
}

func UploadToS3(ctx context.Context, s3Client *s3.Client, bucket, key, content string) error {
	log.Printf("[S3] Uploading content to bucket '%s' with key '%s'", bucket, key)

	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(content),
	})
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	return nil
}
