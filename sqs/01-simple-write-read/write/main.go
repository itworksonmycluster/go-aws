package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
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

	// create a context (timeout)
	duration := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(), duration)
	defer cancel()

	url, err := sqs_service.GetQueueUrl(ctx, sqsClient, queue)
	if err != nil {
		panic(err)
	}

	// print QUEUE URL
	fmt.Printf("Queue: %s\tURL: %s\n", queue, *url)

	// generate message
	inputSend := &sqs.SendMessageInput{
		QueueUrl:    url,
		MessageBody: aws.String("My first message"),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Architecture": {DataType: aws.String("String"), StringValue: aws.String("x86_64")},
			"OSType":       {DataType: aws.String("String"), StringValue: aws.String("linux")},
			"Kernel":       {DataType: aws.String("String"), StringValue: aws.String("6.2.2-arch1-1")},
			"Distro":       {DataType: aws.String("String"), StringValue: aws.String("EndeavourOS")},
		},
	}

	// create a context (timeout)
	duration = 5 * time.Second
	ctx, cancel = context.WithTimeout(context.TODO(), duration)
	defer cancel()

	//Send a message and return id
	id, err := sqs_service.Send(ctx, *sqsClient, inputSend)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Sent message ID: %s\n", *id)

}

func DecodeBase64(encoded string) string {
	decoded, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil {
		log.Printf("could not decode, %v", err)
		return encoded
	}

	return string(decoded)
}
