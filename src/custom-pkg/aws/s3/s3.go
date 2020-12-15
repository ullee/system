package s3

import (
	"bytes"
	. "constants"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"net/http"
	"os"
)

type Context struct {
	fileDir   string
	uploadDir string
	bucket    string
}

func (c *Context) SetFileDir(fileDir string) {
	c.fileDir = fileDir
}

func (c *Context) GetFileDir() string {
	return c.fileDir
}

func (c *Context) SetUploadDir(uploadDir string) {
	c.uploadDir = uploadDir
}

func (c *Context) GetUploadDir() string {
	return c.uploadDir
}

func (c *Context) Upload() error {

	s, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(S3_REGION),
			Credentials: credentials.NewStaticCredentials(S3_ACCESS_KEY, S3_SECRET_KEY, ""),
		},
	})
	if err != nil {
		return err
	}

	// Upload
	err = c.addFileToS3(s)
	if err != nil {
		return err
	}

	return err
}

func (c *Context) addFileToS3(s *session.Session) error {

	// Open the file for use
	file, err := os.Open(c.fileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	c.bucket = S3_BUCKET_STAGING
	if os.Getenv("APP_ENV") == "production" {
		c.bucket = S3_BUCKET_PRODUCTION
	}

	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(c.bucket),
		Key:                  aws.String(c.uploadDir),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}
