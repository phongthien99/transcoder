package adapter

import (
	"example/src/module/origin/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	s3fs "github.com/fclairamb/afero-s3"
	"github.com/spf13/afero"
)

func NewS3Fs(config *model.S3Config) (afero.Fs, error) {
	// Khởi tạo phiên AWS
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(config.Region),
		Endpoint:         aws.String(config.Endpoint),
		Credentials:      credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, ""),
		DisableSSL:       aws.Bool(!config.UseSSL),
		S3ForcePathStyle: aws.Bool(config.ForcePathStyle),
	})
	if err != nil {
		return nil, err
	}
	// Khởi tạo afero S3
	fs := s3fs.NewFs(config.Bucket, sess)

	return fs, nil
}
