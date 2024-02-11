package storage

import (
	"context"
	"dencoder/internal/logging"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Logger = logging.Logger

func DownloadVideo(ctx context.Context, bucket string, sess *session.Session, filename string, logger *Logger) ([]byte, error) {
	downloader := s3manager.NewDownloader(sess)

	buf := aws.NewWriteAtBuffer([]byte{})

	numBytes, err := downloader.DownloadWithContext(ctx,
		buf,
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

func UploadVideo(ctx context.Context, bucket string, sess *session.Session, filename string, video io.Reader, logger *Logger) error {
	uploader := s3manager.NewUploader(sess)
	_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
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

func DeleteVideo(ctx context.Context, bucket string, sess *session.Session, filename string, logger *Logger) error {
	svc := s3.New(sess)

	_, err := svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename)},
	)
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

func VideosCount(bucket string, sess *session.Session, logger *Logger) (int, error) {
	svc := s3.New(sess)

	objects, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: &bucket})
	if err != nil {
		return 0, err
	}

	return len(objects.Contents), nil
}
