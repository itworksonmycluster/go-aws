package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/itworksonmycluster/go-aws/common"
	sqs_service "github.com/itworksonmycluster/go-aws/sqs"
)

type Machine struct {
	Architecture string `json:"architecture"`
	OSType       string `json:"osType"`
	Kernel       string `json:"kernel"`
	Distro       string `json:"distro"`
}

var data = []Machine{
	{Architecture: "x86_64", OSType: "linux", Kernel: "6.2.2-arch2-1", Distro: "EndeavourOS"},
	{Architecture: "x86_64", OSType: "linux", Kernel: "5.15.74-1-pve", Distro: "Ubuntu 22.04.2 LTS"},
	{Architecture: "aarch64", OSType: "linux", Kernel: "5.10.63-v8+", Distro: "Debian GNU/Linux 11"},
	{Architecture: "aarch64", OSType: "linux", Kernel: "5.13.0-1007-raspi", Distro: "Pop!_OS 21.10"},
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

	// create a context (timeout)
	duration := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(), duration)
	defer cancel()

	url, err := sqs_service.GetQueueUrl(ctx, sqsClient, queue)
	if err != nil {
		panic(err)
	}

	for _, d := range data {
		body, err := json.Marshal(d)
		if err != nil {
			fmt.Printf("could not marshal %v, %v", d, err)
			continue
		}

		// create input message
		inputSend := &sqs.SendMessageInput{
			QueueUrl:    url,
			MessageBody: aws.String(string(body)),
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

		fmt.Printf("Message ID sent: %s, %v\n", *id, d)
	}

}
