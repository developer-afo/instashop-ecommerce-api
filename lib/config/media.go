package config

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/developer-afo/instashop-ecommerce-api/dto"
	"github.com/developer-afo/instashop-ecommerce-api/lib/constants"
	"github.com/developer-afo/instashop-ecommerce-api/lib/helper"
)

type MediaInterface interface {
	UploadFile(file *multipart.FileHeader) (string, error)
	GetObject(fileName string) (dto.GetMediaDTO, error)
}

type media struct {
	bucket, folder string
	service        *s3.S3
}

func NewMediaHelper(env constants.Env) MediaInterface {
	return &media{
		bucket:  env.AWS_BUCKET,
		folder:  env.AWS_BUCKET_FOLDER,
		service: s3.New(AWSConfig(env.AWS_REGION, env.AWS_ACCESS_KEY, env.AWS_SECRET_KEY)),
	}
}

func AWSConfig(region string, accessKey string, secretKey string) *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	}))
}

func (m *media) UploadFile(file *multipart.FileHeader) (string, error) {
	fileName := m.FileName(file.Filename)
	fileOpen, openErr := file.Open()

	if openErr != nil {
		return "", openErr
	}

	defer fileOpen.Close()

	fileContent := new(bytes.Buffer)

	_, copyErr := io.Copy(fileContent, fileOpen)

	if copyErr != nil {
		return "", copyErr
	}

	// Uploads the object to S3
	_, err := m.service.PutObject(&s3.PutObjectInput{
		Bucket: helper.StringToPointer(m.bucket),
		Key:    helper.StringToPointer(m.GetObjectKey(fileName)),
		Body:   bytes.NewReader(fileContent.Bytes()),
	})

	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (m *media) GetObject(key string) (dto.GetMediaDTO, error) {
	var media dto.GetMediaDTO

	// Downloads the object to a file
	obj, err := m.service.GetObject(&s3.GetObjectInput{
		Bucket: helper.StringToPointer(m.bucket),
		Key:    helper.StringToPointer(m.GetObjectKey(key)),
	})

	if err != nil {
		return dto.GetMediaDTO{}, err
	}

	media.Body = obj.Body
	media.ContentType = obj.ContentType
	media.ContentLength = obj.ContentLength

	return media, nil
}

func (m *media) FileName(name string) string {
	filename, _ := helper.GenerateSnowflakeID()
	fileExt := strings.Split(name, ".")[1]

	return fmt.Sprintf("%d.%s", filename, fileExt)
}

func (m *media) GetObjectKey(key string) string {
	return m.folder + "/" + key
}
