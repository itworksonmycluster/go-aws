package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/itworksonmycluster/go-aws/common"
	"github.com/itworksonmycluster/go-aws/sqs"
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

	url, err := sqs.GetQueueUrl(ctx, sqsClient, queue)
	if err != nil {
		panic(err)
	}

	fmt.Println(url)

}
