package common

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func NewSqsClient(cfg aws.Config) *sqs.Client {
	return sqs.NewFromConfig(cfg)
}
