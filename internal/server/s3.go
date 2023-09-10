package server

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

// TODO: probably I should somehow reuse connections

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
		return nil, fmt.Errorf("unable to download %q from %q: %w", filename, bucket, err)
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

func DeleteVideo(bucket, filename string, logger *Logger) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return err
	}

	svc := s3.New(sess)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(filename)})
	if err != nil {
		return fmt.Errorf("unable to delete %q from %q: %w", filename, bucket, err)
	}

	// maybe without wait? :)
	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})

	if err != nil {
		return fmt.Errorf("unable to delete %q from %q: %w", filename, bucket, err)
	}

	logger.Infof("Object %s successfully deleted from %s", filename, bucket)
	return nil
}
