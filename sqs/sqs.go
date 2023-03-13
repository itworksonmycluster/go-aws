package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Attribute struct {
	Key   string
	Value string
	Type  string
}

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

func Send(ctx context.Context, client sqs.Client, input *sqs.SendMessageInput) (*string, error) {
	result, err := client.SendMessage(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("could not get queue URL, %w", err)
	}
	return result.MessageId, nil
}

func Get(ctx context.Context, client sqs.Client, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	result, err := client.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("coult not receive a message, %w", err)
	}

	return result, nil
}

func Delete(ctx context.Context, client sqs.Client, input *sqs.DeleteMessageInput) error {
	_, err := client.DeleteMessage(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func SendBatch(ctx context.Context, client sqs.Client, input *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	result, err := client.SendMessageBatch(ctx, input)
	if err != nil {
		fmt.Println(result)
		return nil, err
	}

	return result, nil
}
