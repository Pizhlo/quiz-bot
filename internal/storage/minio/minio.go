package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type minioRepo struct {
	client *minio.Client
	bucket string
}

func New(endpoint, accessKey, secretAccessKey string, useSSL bool, bucket string) (*minioRepo, error) {
	client, err := minio.New(endpoint, &minio.Options{

		Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		return nil, err
	}

	repo := &minioRepo{
		client: client,
		bucket: bucket,
	}

	return repo, nil
}

func (db *minioRepo) Get(ctx context.Context, filePath string) error {
	logrus.Debugf("trying to load file from Minio. Bucket: %s. Filepath: %s", db.bucket, filePath)

	err := db.client.FGetObject(ctx, db.bucket, filePath, filePath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
