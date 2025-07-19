package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileUploadService struct {
	uploader *manager.Uploader
	bucket    string
}

type UploadResult struct {
	URL string
	Key string
	Location string
}

func NewFileUploadService(bucket, region string) (*FileUploadService, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	return &FileUploadService{
        uploader: uploader,
        bucket:   bucket,
    }, nil
}

func (s *FileUploadService) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, userID string, fileType string) (*UploadResult, error) {
    log.Printf("Starting file upload - UserID: %s, FileType: %s, Filename: %s", userID, fileType, header.Filename)
    
    ext := filepath.Ext(header.Filename)
    if ext == "" {
        ext = ".jpg"
        log.Printf("No extension found, defaulting to .jpg")
    }
    log.Printf("File extension: %s", ext)

    key := fmt.Sprintf("%s/%s/%d%s", fileType, userID, time.Now().Unix(), ext)
    log.Printf("Generated S3 key: %s", key)
    
    contentType := s.getContentType(ext)
    log.Printf("Content type: %s", contentType)
    
    log.Printf("Uploading to bucket: %s", s.bucket)

    result, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(key),
        Body:        file,
        ContentType: aws.String(contentType),
    })
    if err != nil {
        log.Printf("Upload failed: %v", err)
        return nil, fmt.Errorf("failed to upload file: %w", err)
    }

    log.Printf("Upload successful - Location: %s", result.Location)

    return &UploadResult{
        URL: result.Location,
        Key: key,
        Location: result.Location,
    }, nil
}

func (s *FileUploadService) UploadFromReader(ctx context.Context, reader io.Reader, key, contentType string) (*UploadResult, error) {
    result, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(key),
        Body:        reader,
        ContentType: aws.String(contentType),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to upload file: %w", err)
    }

    return &UploadResult{
        URL:      result.Location,
        Key:      key,
        Location: result.Location,
    }, nil
}

func (s *FileUploadService) getContentType(ext string) string {
    switch strings.ToLower(ext) {
    case ".jpg", ".jpeg":
        return "image/jpeg"
    case ".png":
        return "image/png"
    case ".gif":
        return "image/gif"
    case ".webp":
        return "image/webp"
    case ".mp4":
        return "video/mp4"
    case ".webm":
        return "video/webm"
    case ".pdf":
        return "application/pdf"
    default:
        return "application/octet-stream"
    }
}
