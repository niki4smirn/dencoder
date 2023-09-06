package server

import (
	"context"
	"net/http"
)

type S3Bucket struct {
	Name string
}

const contextS3BucketKey = "s3-bucket"

func setS3Bucket(ctx context.Context, b *S3Bucket) context.Context {
	return context.WithValue(ctx, contextS3BucketKey, b)
}

func GetS3Bucket(ctx context.Context) *S3Bucket {
	user, ok := ctx.Value(contextS3BucketKey).(*S3Bucket)

	if !ok {
		return nil
	}

	return user
}

func WithLogAndErr(handler func(http.ResponseWriter, *http.Request, *Logger) error, logger *Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r, logger); err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func SetS3BucketName(bucketName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := setS3Bucket(r.Context(), &S3Bucket{
				Name: bucketName,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
