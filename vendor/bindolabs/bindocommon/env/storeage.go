package env

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Storage struct {
	Host              string
	S3AccessKeyID     string
	S3AccessKeySecret string
	S3Region          string
	Cdn               string
	Buckets           *struct {
		Asserts *struct {
			Bucket string
			Host   string
		}
	}
	DefaultBucket string
	AwsSession    *session.Session
}

func GetAwsSession() *session.Session {
	if Env.Storage != nil && Env.Storage.AwsSession == nil {
		Env.Storage.AwsSession = session.Must(session.NewSession(&aws.Config{
			Region: aws.String(Env.Storage.S3Region),
			Credentials: credentials.NewStaticCredentials(Env.Storage.S3AccessKeyID,
				Env.Storage.S3AccessKeySecret, ""),
		}))
	}
	return Env.Storage.AwsSession
}
