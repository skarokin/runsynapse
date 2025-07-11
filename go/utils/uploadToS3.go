package utils

import (
	"crypto/sha256"
	"time"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var allowedTypes = map[string]bool{
	// Images
    "image/jpeg":    true,
    "image/png":     true,
    "image/gif":     true,
    "image/webp":    true,
    "image/svg+xml": false, // SVG can contain scripts - avoid unless needed
    
    // Documents
    "application/pdf":     true,
    "text/plain":          true,
    "text/csv":            true,
    "text/markdown":       true,
    
    // Modern Office formats (safer than old ones)
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document":   true, // .docx
    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true, // .xlsx
    "application/vnd.openxmlformats-officedocument.presentationml.presentation": true, // .pptx
    
    // Archive formats
    "application/zip": true,
    "application/x-tar": true,
    "application/gzip": true,
    
    // Audio (for voice messages)
    "audio/mpeg":  true, // .mp3
    "audio/wav":   true,
    "audio/ogg":   true,
    "audio/webm":  true,
    
    // Video (for video messages)
    "video/mp4":   true,
    "video/webm":  true,
    "video/ogg":   true,
    
    // Code files
    "application/json": true,
    "text/javascript":  true,
    "text/css":         true,
    "text/html":        true,
    
    // Avoid these dangerous types:
    "application/msword":           false, // .doc - can contain macros
    "application/vnd.ms-excel":     false, // .xls - can contain macros
    "application/vnd.ms-powerpoint": false, // .ppt - can contain macros
    "application/x-msdownload":     false, // .exe
    "application/x-executable":     false, // executables
    "application/x-shockwave-flash": false, // .swf
}

const maxFileSize = 12 * 1024 * 1024 // 12 MB

func generateKeyFromFilename(filename string) string {
	// get file extension and base name
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	// generate unique hash from extension, base name, nanoseconds, and uuid
	hash := sha256.New()
    hash.Write([]byte(fmt.Sprintf("%s_%d_%s", base, time.Now().UnixNano(), uuid.New().String())))
	hashStr := fmt.Sprintf("%x", hash.Sum(nil))[:12]	// first 12 chars; reduce length for readability

	// clean up filename and combine with hash
	cleanName := strings.ReplaceAll(base, " ", "_")
	cleanName = strings.ReplaceAll(cleanName, "/", "_")

	// return final key with cleanName_hash.ext
    return fmt.Sprintf("%s_%s%s", cleanName, hashStr, ext)
}

func UploadFilesToS3(ctx context.Context, s3Client *s3.Client, bucket string, files []*multipart.FileHeader) ([]string, error) {
	log.Printf("[S3] Uploading %d files to bucket '%s'", len(files), bucket)

    var attachmentURLs []string
    
    if len(files) == 0 {
        return attachmentURLs, nil
    }
    
    log.Printf("Processing %d file(s)", len(files))
    
    for _, fileHeader := range files {
		contentType := fileHeader.Header.Get("Content-Type")
		if !validateFileType(contentType) {
			log.Printf("Invalid file type for %s: %s", fileHeader.Filename, contentType)
			return nil, fmt.Errorf("invalid file type for %s: %s", fileHeader.Filename, contentType)
		}

		if fileHeader.Size > maxFileSize {
			log.Printf("File %s exceeds maximum size of %d bytes: %d bytes", fileHeader.Filename, maxFileSize, fileHeader.Size)
			return nil, fmt.Errorf("file %s exceeds maximum size of %d bytes", fileHeader.Filename, maxFileSize)
		}

        file, err := fileHeader.Open()
        if err != nil {
            log.Printf("Error opening file %s: %v", fileHeader.Filename, err)
            return nil, fmt.Errorf("failed to open file %s: %w", fileHeader.Filename, err)
        }
        defer file.Close()

        key := generateKeyFromFilename(fileHeader.Filename)
        
        err = uploadFileToS3(
            ctx,
            s3Client,
            bucket, 
            key,
            file,
            fileHeader.Header.Get("Content-Type"),
        )
        if err != nil {
            log.Printf("Error uploading file %s: %v", fileHeader.Filename, err)
            return nil, fmt.Errorf("failed to upload file %s: %w", fileHeader.Filename, err)
        }
        
        attachmentURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key)
        attachmentURLs = append(attachmentURLs, attachmentURL)
    }
    
    return attachmentURLs, nil
}

func uploadFileToS3(ctx context.Context, s3Client *s3.Client, bucket, key string, file io.Reader, contentType string) error {
    log.Printf("[S3] Uploading file to bucket '%s' with key '%s'", bucket, key)

    _, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(bucket),
        Key:         aws.String(key),
        Body:        file,
        ContentType: aws.String(contentType),
		CacheControl: aws.String("max-age=31536000, public"), // 1 year cache
    })
    if err != nil {
        return fmt.Errorf("failed to upload file to S3: %w", err)
    }

    log.Printf("[S3] Successfully uploaded file with key: %s", key)
	
    return nil
}

func DeleteFromS3(ctx context.Context, s3Client *s3.Client, bucket, key string) error {
	log.Printf("[S3] Deleting file with key '%s' from bucket '%s'", key, bucket)

	_, err := s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	log.Printf("[S3] Successfully deleted file with key: %s", key)
	return nil
}

func validateFileType(fileType string) bool {
	if fileType == "" {
		return false
	}

	if allowed, exists := allowedTypes[fileType]; exists {
		return allowed
	}
	
	return false
}