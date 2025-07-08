package inits

import (
    "fmt"
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)
func NewS3Client(s3Region string) (*s3.Client, error) {
	if s3Region == "" {
		return nil, fmt.Errorf("S3_REGION environemnt variable is required")
	}

    cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(s3Region),
	)
    if err != nil {
        return nil, fmt.Errorf("failed to load AWS config: %w", err)
    }

	return s3.NewFromConfig(cfg), nil
}