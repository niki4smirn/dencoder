package server

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

func DownloadVideo(bucket, filename string, logger *Logger) ([]byte, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}

	downloader := s3manager.NewDownloader(sess)

	buf := aws.NewWriteAtBuffer([]byte{})

	numBytes, err := downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filename),
		})
	if err != nil {
		return nil, err
	}

	logger.Infof("Downloaded %s (%v bytes)", filename, numBytes)

	return buf.Bytes(), nil
}

func UploadVideo(bucket, filename string, video io.Reader, logger *Logger) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   video,
	})
	if err != nil {
		return fmt.Errorf("unable to upload %q to %q: %w", filename, bucket, err)
	}

	logger.Infof("Successfully uploaded %q to %q\n", filename, bucket)

	return nil
}
