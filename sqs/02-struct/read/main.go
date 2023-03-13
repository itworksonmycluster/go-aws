package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/itworksonmycluster/go-aws/common"
	sqs_service "github.com/itworksonmycluster/go-aws/sqs"
)

type Machine struct {
	Architecture string `json:"architecture"`
	OSType       string `json:"osType"`
	Kernel       string `json:"kernel"`
	Distro       string `json:"distro"`
}

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
		MaxNumberOfMessages: 10,
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

	fmt.Printf("Messages received: %d\n", len(outputReceive.Messages))
	var machines []Machine
	for _, message := range outputReceive.Messages {
		var m Machine
		err := json.Unmarshal([]byte(*message.Body), &m)
		if err != nil {
			fmt.Printf("could not unmarshal body, %v\n", err)
			continue
		}
		machines = append(machines, m)
		inputDelete := &sqs.DeleteMessageInput{
			QueueUrl:      &queue,
			ReceiptHandle: message.ReceiptHandle,
		}

		err = sqs_service.Delete(ctx, *sqsClient, inputDelete)
		if err != nil {
			fmt.Printf("could not delete de message %s, %v\n", *message.MessageId, err)
		} else {
			fmt.Printf("message ID %s deleted\n", *message.MessageId)
		}
	}

	fmt.Printf("Number of processed messages: %d\n", len(machines))

	fmt.Println(machines)

}
