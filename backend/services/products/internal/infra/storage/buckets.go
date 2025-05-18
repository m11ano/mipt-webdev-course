package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/m11ano/mipt-webdev-course/backend/services/products/internal/infra/config"
)

type BucketName string

const (
	BucketProducts BucketName = "products"
)

func CreateBuckets(ctx context.Context, client Client) error {
	list := []struct {
		bucket BucketName
		policy string
	}{
		{
			bucket: BucketProducts,
			policy: fmt.Sprintf(`{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "AWS": [
                    "*"
                ]
            },
            "Action": [
                "s3:GetObject"
            ],
            "Resource": [
                "arn:aws:s3:::%s/*"
            ]
        }
    ]
}`, BucketProducts),
		},
	}

	for _, item := range list {
		err := client.CreateBucket(ctx, item.bucket, item.policy)
		if err != nil && !errors.Is(err, ErrBucketAlreadyExists) {
			return err
		}
	}
	return nil
}

func GetBucketURL(bucket BucketName, cfg *config.Config) string {
	var builder strings.Builder
	builder.WriteString(cfg.Storage.S3URL)
	builder.WriteString("/")
	builder.WriteString(string(bucket))

	return builder.String()
}
