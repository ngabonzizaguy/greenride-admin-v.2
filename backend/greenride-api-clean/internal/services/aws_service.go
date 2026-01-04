package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	appconfig "greenride/internal/config"
	"greenride/internal/utils"
)

var (
	awsServiceInstance *AWSService
	awsServiceOnce     sync.Once

	// 允许的图片扩展名
	AllowedImageExts = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	// 允许的文件扩展名
	AllowedFileExts = map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".txt":  true,
		".csv":  true,
		".xlsx": true,
	}
)

// AWSService AWS 服务
type AWSService struct {
	s3Client  *s3.Client
	bucket    string
	region    string
	available bool
}

// GetAWSService 获取 AWS 服务单例
func GetAWSService() *AWSService {
	awsServiceOnce.Do(func() {
		SetupAWSService()
	})
	return awsServiceInstance
}

// SetupAWSService 初始化 AWS 服务
func SetupAWSService() {
	// 获取配置
	cfg := appconfig.Get()
	if cfg == nil || cfg.AWS.GetS3Bucket() == "" || cfg.AWS.GetS3Region() == "" {
		log.Printf("AWS service initialization failed: missing required configuration (bucket or region)")
		awsServiceInstance = &AWSService{
			available: false,
		}
		return
	}

	// 使用 IAM 角色自动获取凭证，适合生产环境
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AWS.GetS3Region()),
	)

	if err != nil {
		log.Printf("AWS service initialization failed: %v", err)
		awsServiceInstance = &AWSService{
			available: false,
		}
		return
	}

	// 验证 AWS 凭证
	svc := sts.NewFromConfig(awsConfig)
	_, err = svc.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Printf("AWS credentials validation failed: %v", err)
		awsServiceInstance = &AWSService{
			available: false,
		}
		return
	}

	// 创建 S3 客户端
	s3Client := s3.NewFromConfig(awsConfig)

	awsServiceInstance = &AWSService{
		s3Client:  s3Client,
		bucket:    cfg.AWS.GetS3Bucket(),
		region:    cfg.AWS.GetS3Region(),
		available: true,
	}

	log.Printf("AWS service initialized successfully with bucket: %s, region: %s", awsServiceInstance.bucket, awsServiceInstance.region)
}

// GetAWSServiceSafe 安全获取 AWS 服务实例，返回是否可用
func GetAWSServiceSafe() (*AWSService, bool) {
	service := GetAWSService()
	if service == nil {
		return nil, false
	}
	return service, service.IsAvailable()
}

// IsAvailable 检查 AWS 服务是否可用
func (s *AWSService) IsAvailable() bool {
	return s != nil && s.available && s.s3Client != nil
}

// UploadImage 上传图片到S3
func (s *AWSService) UploadImage(ctx context.Context, file *os.File, folder string) (string, error) {
	if !s.IsAvailable() {
		return "", fmt.Errorf("AWS S3 client not available")
	}

	// 检查文件扩展名
	ext := strings.ToLower(path.Ext(file.Name()))
	if !AllowedImageExts[ext] {
		return "", fmt.Errorf("unsupported image format: %s", ext)
	}

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %v", err)
	}

	// 读取文件内容
	buffer := make([]byte, fileInfo.Size())
	if _, err := file.Read(buffer); err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// 生成文件路径
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "image/jpeg"
	}

	// 生成唯一的对象键
	objectKey := fmt.Sprintf("%s/%s%s", folder, utils.GenerateID(), ext)

	// 上传到S3
	_, err = s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(buffer),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %v", err)
	}

	return s.GenerateObjectURL(objectKey), nil
}

// UploadImageFromURL 从URL下载图片并上传到S3
func (s *AWSService) UploadImageFromURL(ctx context.Context, imageURL string, folder, name string) (string, error) {
	if !s.IsAvailable() {
		return "", fmt.Errorf("AWS S3 client not available")
	}

	// 下载图片
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image from URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image, status code: %d", resp.StatusCode)
	}

	// 读取图片内容
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %v", err)
	}

	if name == "" {
		name = utils.GenerateID()
	}

	// 生成对象键
	objectKey := fmt.Sprintf("%s/%s", folder, name)

	// 上传到S3
	_, err = s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(imageData),
		ContentType: aws.String(resp.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %v", err)
	}

	return s.GenerateObjectURL(objectKey), nil
}

// UploadFile 上传任意文件到S3
func (s *AWSService) UploadFile(ctx context.Context, file *os.File, folder, fileName string) (string, error) {
	if !s.IsAvailable() {
		return "", fmt.Errorf("AWS S3 client not available")
	}

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %v", err)
	}

	// 读取文件内容
	buffer := make([]byte, fileInfo.Size())
	if _, err := file.Read(buffer); err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// 检测 Content-Type
	contentType := http.DetectContentType(buffer)
	if contentType == "" {
		ext := strings.ToLower(path.Ext(file.Name()))
		contentType = mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	// 生成对象键
	var objectKey string
	if fileName != "" {
		objectKey = fmt.Sprintf("%s/%s", folder, fileName)
	} else {
		objectKey = fmt.Sprintf("%s/%s", folder, file.Name())
	}

	// 上传到S3
	_, err = s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(buffer),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %v", err)
	}

	return s.GenerateObjectURL(objectKey), nil
}

// UploadFromReader 从 io.Reader 上传文件
func (s *AWSService) UploadFromReader(ctx context.Context, reader io.Reader, objectKey, contentType string) (string, error) {
	if !s.IsAvailable() {
		return "", fmt.Errorf("AWS S3 client not available")
	}

	// 上传到S3
	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(objectKey),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %v", err)
	}

	return s.GenerateObjectURL(objectKey), nil
}

// GenerateObjectURL 生成对象访问URL
func (s *AWSService) GenerateObjectURL(key string) string {
	// 使用正确的 S3 URL 格式，因为存储桶是公开读取的
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)
}

// GeneratePresignedURL 生成预签名URL
func (s *AWSService) GeneratePresignedURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	if !s.IsAvailable() {
		return "", fmt.Errorf("AWS S3 client not available")
	}

	presignClient := s3.NewPresignClient(s.s3Client)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = duration
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	return request.URL, nil
}

// UploadUserAvatar 上传用户头像的便捷方法
func (s *AWSService) UploadUserAvatar(ctx context.Context, userID string, reader io.Reader, fileExtension string) (string, error) {
	if !s.IsAvailable() {
		return "", fmt.Errorf("AWS S3 client not available")
	}

	objectKey := fmt.Sprintf("avatars/%s%s", userID, fileExtension)

	// 检测内容类型
	contentType := mime.TypeByExtension(fileExtension)
	if contentType == "" {
		contentType = "image/jpeg"
	}

	// 上传到S3，不设置ACL（依赖Bucket Policy进行公共访问）
	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(objectKey),
		Body:        reader,
		ContentType: aws.String(contentType),
		// 移除ACL设置，改为使用Bucket Policy控制访问权限
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar to S3: %v", err)
	}

	return s.GenerateObjectURL(objectKey), nil
}

// DeleteObject 删除S3对象
func (s *AWSService) DeleteObject(ctx context.Context, key string) error {
	if !s.IsAvailable() {
		return fmt.Errorf("AWS S3 client not available")
	}

	_, err := s.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	return err
}
