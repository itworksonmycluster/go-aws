package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/itworksonmycluster/go-aws/common"
	sqs_service "github.com/itworksonmycluster/go-aws/sqs"
)

func main() {
	var queue string
	flag.StringVar(&queue, "queue", "", "The name of the queue")
	flag.Parse()

	if queue == "" {
		fmt.Println("you need to pass the queue name via flag (-queue QUEUE) or env")
		return
	}

	// load AWS config
	cfg, err := common.LoadConfig()
	if err != nil {
		panic(err)
	}

	// get sqs client
	sqsClient := common.NewSqsClient(*cfg)

	// Get a message
	inputReceive := &sqs.ReceiveMessageInput{
		QueueUrl: &queue,
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		MaxNumberOfMessages: 1,
		VisibilityTimeout:   10,
	}

	// create a context (timeout)
	duration := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(), duration)
	defer cancel()

	outputReceive, err := sqs_service.Get(ctx, *sqsClient, inputReceive)
	if err != nil {
		panic(err)
	}

	if outputReceive.Messages != nil {
		fmt.Printf("Received message ID: %s\n", *outputReceive.Messages[0].MessageId)
		fmt.Printf("Message body: %s\n", *outputReceive.Messages[0].Body)
		fmt.Println("Attributes:")
		for k, v := range outputReceive.Messages[0].MessageAttributes {
			fmt.Printf("\tkey: %s\tvalue: %s\n", k, *v.StringValue)
		}
		inputDelete := &sqs.DeleteMessageInput{
			QueueUrl:      &queue,
			ReceiptHandle: outputReceive.Messages[0].ReceiptHandle,
		}

		err := sqs_service.Delete(ctx, *sqsClient, inputDelete)
		if err != nil {
			fmt.Printf("could not delete de message %s, %v\n", *outputReceive.Messages[0].MessageId, err)
		} else {
			fmt.Printf("message ID %s deleted\n", *outputReceive.Messages[0].MessageId)
		}

	} else {
		fmt.Println("No messages found")
	}
}
