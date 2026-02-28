package service

import "errors"

var (
	// ErrInvalidArgument 请求参数错误
	ErrInvalidArgument = errors.New("invalid argument")
	// ErrTaskNotFound 任务不存在
	ErrTaskNotFound = errors.New("task not found")
	// ErrTaskStatusConflict 任务状态冲突
	ErrTaskStatusConflict = errors.New("task status conflict")
	// ErrRequiredDocumentsMissing 缺少必要文档
	ErrRequiredDocumentsMissing = errors.New("required documents missing")
	// ErrUnsupportedFileType 不支持的文件类型
	ErrUnsupportedFileType = errors.New("unsupported file type")
	// ErrFileTooLarge 文件过大
	ErrFileTooLarge = errors.New("file too large")
	// ErrForbidden 无权限访问
	ErrForbidden = errors.New("forbidden")
)
