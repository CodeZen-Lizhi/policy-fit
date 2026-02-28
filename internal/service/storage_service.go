package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zhenglizhi/policy-fit/internal/config"
	"github.com/zhenglizhi/policy-fit/internal/domain"
)

const (
	// DefaultMaxUploadSize 默认文件上传大小限制：30MB
	DefaultMaxUploadSize int64 = 30 << 20
)

// StorageService 文件存储服务
type StorageService struct {
	adapter       storageAdapter
	maxUploadSize int64
}

type storageAdapter interface {
	Save(ctx context.Context, key string, file multipart.File, size int64) error
	Delete(ctx context.Context, key string) error
}

// NewStorageService 创建存储服务
func NewStorageService(cfg config.StorageConfig) (*StorageService, error) {
	var adapter storageAdapter
	switch cfg.Type {
	case "local":
		adapter = &localStorageAdapter{
			basePath: cfg.Path,
		}
	case "s3":
		s3, err := newS3StorageAdapter(cfg)
		if err != nil {
			return nil, err
		}
		adapter = s3
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}

	return &StorageService{
		adapter:       adapter,
		maxUploadSize: DefaultMaxUploadSize,
	}, nil
}

// SaveUploadedFile 保存上传文件
func (s *StorageService) SaveUploadedFile(
	ctx context.Context,
	taskID int64,
	docType domain.DocumentType,
	fileHeader *multipart.FileHeader,
) (string, error) {
	if fileHeader == nil {
		return "", fmt.Errorf("%w: empty file", ErrInvalidArgument)
	}
	if fileHeader.Size <= 0 {
		return "", fmt.Errorf("%w: empty file size", ErrInvalidArgument)
	}
	if fileHeader.Size > s.maxUploadSize {
		return "", ErrFileTooLarge
	}
	if !isPDFFile(fileHeader.Filename) {
		return "", ErrUnsupportedFileType
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("open uploaded file: %w", err)
	}
	defer file.Close()

	storageKey := buildStorageKey(taskID, docType, fileHeader.Filename)
	if err := s.adapter.Save(ctx, storageKey, file, fileHeader.Size); err != nil {
		return "", err
	}

	return storageKey, nil
}

// Delete 删除存储对象
func (s *StorageService) Delete(ctx context.Context, storageKey string) error {
	return s.adapter.Delete(ctx, storageKey)
}

type localStorageAdapter struct {
	basePath string
}

func (a *localStorageAdapter) Save(ctx context.Context, key string, file multipart.File, _ int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	base := a.basePath
	if strings.TrimSpace(base) == "" {
		base = "./storage"
	}

	fullPath := filepath.Join(base, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	out, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("create target file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return fmt.Errorf("copy uploaded file: %w", err)
	}

	return nil
}

func (a *localStorageAdapter) Delete(_ context.Context, key string) error {
	base := a.basePath
	if strings.TrimSpace(base) == "" {
		base = "./storage"
	}
	fullPath := filepath.Join(base, filepath.FromSlash(key))
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete local object: %w", err)
	}
	return nil
}

type s3StorageAdapter struct {
	client *minio.Client
	bucket string
}

func newS3StorageAdapter(cfg config.StorageConfig) (*s3StorageAdapter, error) {
	endpoint, secure, err := normalizeEndpoint(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, fmt.Errorf("create s3 client: %w", err)
	}

	return &s3StorageAdapter{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

func (a *s3StorageAdapter) Save(ctx context.Context, key string, file multipart.File, size int64) error {
	if err := a.ensureBucket(ctx); err != nil {
		return err
	}

	_, err := a.client.PutObject(ctx, a.bucket, key, file, size, minio.PutObjectOptions{
		ContentType: "application/pdf",
	})
	if err != nil {
		return fmt.Errorf("upload object to s3: %w", err)
	}
	return nil
}

func (a *s3StorageAdapter) Delete(ctx context.Context, key string) error {
	if err := a.client.RemoveObject(ctx, a.bucket, key, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("delete object from s3: %w", err)
	}
	return nil
}

func (a *s3StorageAdapter) ensureBucket(ctx context.Context) error {
	exists, err := a.client.BucketExists(ctx, a.bucket)
	if err != nil {
		return fmt.Errorf("check bucket exists: %w", err)
	}
	if exists {
		return nil
	}
	if err := a.client.MakeBucket(ctx, a.bucket, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("create bucket: %w", err)
	}
	return nil
}

func buildStorageKey(taskID int64, docType domain.DocumentType, fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == "" {
		ext = ".pdf"
	}
	return path.Join(
		fmt.Sprintf("task/%d", taskID),
		string(docType),
		fmt.Sprintf("%s_%d%s", randomHex(8), time.Now().Unix(), ext),
	)
}

func randomHex(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "random"
	}
	return hex.EncodeToString(buf)
}

func isPDFFile(fileName string) bool {
	return strings.EqualFold(filepath.Ext(fileName), ".pdf")
}

func normalizeEndpoint(endpoint string) (string, bool, error) {
	raw := strings.TrimSpace(endpoint)
	if raw == "" {
		return "", false, fmt.Errorf("%w: empty s3 endpoint", ErrInvalidArgument)
	}

	if strings.Contains(raw, "://") {
		u, err := url.Parse(raw)
		if err != nil {
			return "", false, fmt.Errorf("parse s3 endpoint: %w", err)
		}
		return u.Host, u.Scheme == "https", nil
	}

	return raw, false, nil
}
