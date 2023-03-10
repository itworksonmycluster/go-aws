package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func GetQueueUrl(ctx context.Context, client *sqs.Client, name string) (*string, error) {
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &name,
	}

	result, err := client.GetQueueUrl(ctx, gQInput)
	if err != nil {
		return nil, fmt.Errorf("could not get queue URL, %w", err)
	}

	return result.QueueUrl, nil
}
