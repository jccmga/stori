package filereader

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"stori/model"
)

const awsRegion = "us-east-1"

type S3 struct {
	fileURI string
}

func NewS3Reader(fileURI string) S3 {
	return S3{fileURI: fileURI}
}

func (reader S3) ReadTransactions() ([]model.Transaction, error) {
	fail := func(err error) ([]model.Transaction, error) {
		return nil, fmt.Errorf("filereader: S3: ReadTransactions: %w", err)
	}

	bucket, key, err := parseS3URI(reader.fileURI)
	if err != nil {
		return fail(fmt.Errorf("%w: %w", ErrInvalidURI, err))
	}

	destPath := filepath.Join(os.TempDir(), key)
	if errDownload := downloadFileFromS3(bucket, key, destPath); errDownload != nil {
		return fail(errDownload)
	}

	transactions, err := NewLocalReader(destPath).ReadTransactions()
	if err != nil {
		return fail(err)
	}

	return transactions, nil
}

func parseS3URI(s3URI string) (bucket, key string, err error) {
	fail := func(err error) (string, string, error) {
		return "", "", fmt.Errorf("filereader: parseS3URI: %w", err)
	}

	parsedURL, err := url.Parse(s3URI)
	if err != nil {
		return fail(fmt.Errorf("failed to parse S3 URI: %w", err))
	}

	if parsedURL.Scheme != "s3" {
		return fail(fmt.Errorf("invalid S3 URI scheme: %s", parsedURL.Scheme))
	}

	bucket = parsedURL.Host
	key = strings.TrimPrefix(parsedURL.Path, "/")

	return bucket, key, nil
}

func downloadFileFromS3(bucket, key, destPath string) error {
	fail := func(err error) error {
		return fmt.Errorf("filereader: downloadFileFromS3: %w", err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.AnonymousCredentials,
	})
	if err != nil {
		return fail(fmt.Errorf("%w: %w", ErrS3Connection, err))
	}

	file, err := os.Create(destPath)
	if err != nil {
		return fail(fmt.Errorf("%w (%s): %w", ErrFileCreation, destPath, err))
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(
		file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		return fail(fmt.Errorf("%w: %w", ErrDownloadFile, err))
	}

	return nil
}
