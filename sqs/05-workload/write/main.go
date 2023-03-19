package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/itworksonmycluster/go-aws/common"
	sqs_service "github.com/itworksonmycluster/go-aws/sqs"
)

// some fields must not be public either exported in real world.
type User struct {
	Uuid     string `json:"id"`
	Username string `json:"username"`
	Title    string `json:"title"`
	Company  string `json:"company"`
	Active   bool   `json:"active"`
	Password string `json:"password"`
}

func loadData() []User {
	var users []User
	data, err := os.ReadFile("./test/data/users219.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &users); err != nil {
		panic(err)
	}

	return users
}

var successful = make(chan []types.SendMessageBatchResultEntry)
var failed = make(chan []types.BatchResultErrorEntry)

func main() {

	var queue string
	flag.StringVar(&queue, "queue", "", "The name of the queue")
	flag.Parse()

	if queue == "" {
		fmt.Println("you need to pass the queue name via flag (-queue QUEUE) or env")
		return
	}

	// load data
	users := loadData()

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
	fmt.Println(*url)

	var wg sync.WaitGroup
	size := 10
	chunk := float64(len(users)) / float64(size)
	n := int(math.Ceil(chunk))
	var j int

	wg.Add(n)
	fmt.Printf("number of workers: %d\n", n)
	time.Sleep(5 * time.Second)
	for i := 0; i < len(users); i += size {
		j += size
		if j > len(users) {
			j = len(users)
		}
		go func(i int, url *string, data []User) {
			defer wg.Done()
			entries := make([]types.SendMessageBatchRequestEntry, 0, len(data))
			for i, d := range data {
				body, err := json.Marshal(d)
				if err != nil {
					fmt.Printf("could not marshal %v, %v", d, err)
					continue
				}

				id := strconv.Itoa(i)

				entry := types.SendMessageBatchRequestEntry{
					Id:          aws.String(id),
					MessageBody: aws.String(string(body)),
				}
				entries = append(entries, entry)
			}

			for _, e := range entries {
				fmt.Println(*e.Id, *e.MessageBody)
			}

			// create input message
			inputSendBatch := &sqs.SendMessageBatchInput{
				QueueUrl: url,
				Entries:  entries,
			}

			// create a context (timeout)
			duration = 5 * time.Second
			ctx, cancel = context.WithTimeout(context.TODO(), duration)
			defer cancel()

			//Send a message and return id
			output, err := sqs_service.SendBatch(ctx, *sqsClient, inputSendBatch)
			if err != nil {
				panic(err)
			}

			failed <- output.Failed
			successful <- output.Successful
			fmt.Printf("end go: %d\n", i)

		}(i, url, users[i:j])
	}
	wg.Wait()
}

func sendBatchUsers(wg *sync.WaitGroup, client sqs.Client) {

}
