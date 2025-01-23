package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type minioRepo struct {
	client *minio.Client
}

func New(endpoint, accessKey, secretAccessKey string, useSSL bool) (*minioRepo, error) {
	client, err := minio.New(endpoint, &minio.Options{

		Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		return nil, err
	}

	repo := &minioRepo{
		client: client,
	}

	return repo, nil
}

func (db *minioRepo) Get(ctx context.Context, bucketName string, objectName string, filePath string) ([]byte, error) {
	logrus.Debugf("trying to load file from Minio. Bucket: %s. Object: %s. Filepath: %s", bucketName, objectName, filePath)

	// Retrieve the file from MinIO
	object, err := db.client.GetObject(ctx, bucketName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	data := []byte{}

	n, err := object.Read(data)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("loaded %d byte(s)", n)

	return data, nil
}
