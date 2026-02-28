package service

import (
	"bytes"
	"context"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"github.com/zhenglizhi/policy-fit/internal/config"
	"github.com/zhenglizhi/policy-fit/internal/domain"
)

func TestStorageServiceLocalSaveAndDelete(t *testing.T) {
	tmpDir := t.TempDir()
	svc, err := NewStorageService(config.StorageConfig{
		Type: "local",
		Path: tmpDir,
	})
	if err != nil {
		t.Fatalf("NewStorageService error: %v", err)
	}

	fileHeader := createPDFFileHeader(t, "report.pdf", []byte("%PDF-1.4 dummy"))
	key, err := svc.SaveUploadedFile(context.Background(), 10, domain.DocTypeReport, fileHeader)
	if err != nil {
		t.Fatalf("SaveUploadedFile error: %v", err)
	}

	fullPath := filepath.Join(tmpDir, filepath.FromSlash(key))
	if _, err := os.Stat(fullPath); err != nil {
		t.Fatalf("saved file not found: %v", err)
	}

	if err := svc.Delete(context.Background(), key); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		t.Fatalf("expected file removed, err=%v", err)
	}
}

func TestStorageServiceS3AdapterInit(t *testing.T) {
	svc, err := NewStorageService(config.StorageConfig{
		Type:      "s3",
		Endpoint:  "http://localhost:9000",
		Bucket:    "policy-fit-test",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin123",
	})
	if err != nil {
		t.Fatalf("NewStorageService(s3) error: %v", err)
	}
	if svc == nil {
		t.Fatalf("expected storage service")
	}
}

func createPDFFileHeader(t *testing.T, fileName string, content []byte) *multipart.FileHeader {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("write content: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	reader := multipart.NewReader(bytes.NewReader(body.Bytes()), writer.Boundary())
	form, err := reader.ReadForm(int64(len(body.Bytes())))
	if err != nil {
		t.Fatalf("ReadForm: %v", err)
	}
	files := form.File["file"]
	if len(files) == 0 {
		t.Fatalf("file header not found")
	}
	return files[0]
}
