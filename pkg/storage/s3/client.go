package s3

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	internalerror "github.com/igilgyrg/betera-test/internal/error"
	"io/ioutil"
	"net/url"
	"strings"
)

const urlMedia = "https://%s.s3.amazonaws.com/%s"

type S3Storage interface {
	List(filename string) ([]Media, error)
	Download(filename string) ([]byte, error)
	Store(filename string, bytes []byte) error
}

type s3Client struct {
	client     *s3.S3
	downloader *s3manager.Downloader
	uplodaer   *s3manager.Uploader
	cfg        *s3Config
}

func NewS3Storage(cfg *s3Config) (S3Storage, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(cfg.accessKey, cfg.secretKey, ""),
		Region:      aws.String(cfg.region)},
	)
	if err != nil {
		return nil, errors.New("error of init session of s3")
	}

	client := s3.New(sess)
	downloader := s3manager.NewDownloader(sess)
	uploader := s3manager.NewUploader(sess)
	return &s3Client{client: client, downloader: downloader, uplodaer: uploader, cfg: cfg}, nil
}

func (s s3Client) List(filename string) ([]Media, error) {
	objects, err := s.client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s.cfg.bucketName),
		Prefix: aws.String(filename),
	})
	if err != nil {
		return nil, err
	}

	if len(objects.Contents) == 0 {
		return []Media{}, nil
	}

	result := make([]Media, 0, len(objects.Contents))

	for i := range objects.Contents {
		object := objects.Contents[i]
		if objects != nil && *object.Size != 0 {
			u, _ := url.Parse(fmt.Sprintf(urlMedia, s.cfg.bucketName, *object.Key))
			media := Media{
				Filename:     strings.Split(*object.Key, "/")[2],
				Date:         strings.Split(*object.Key, "/")[1],
				Url:          u.String(),
				LastModified: *object.LastModified,
			}
			result = append(result, media)
		}
	}

	return result, nil
}

func (s s3Client) Download(filename string) ([]byte, error) {
	resp, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.cfg.bucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		var s3Erro s3.RequestFailure
		if errors.As(err, &s3Erro) {
			return nil, internalerror.New(s3Erro.StatusCode(), s3Erro.Message(), s3Erro.Message())
		}
		return nil, internalerror.NewImageNotFound(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

func (s s3Client) Store(filename string, contentBytes []byte) error {

	result, err := s.uplodaer.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.cfg.bucketName),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(contentBytes),
		ACL:    aws.String("public-read"),
	})

	if err != nil {
		return err
	}

	if result.Location == "" {
		return errors.New("error of store file")
	}

	return nil
}
