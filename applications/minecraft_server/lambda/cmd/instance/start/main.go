package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

var (
	region      = os.Getenv("AWS_REGION")
	instance_id = os.Getenv("MINECRAFT_SERVER_INSTANCE_ID")
)

func handler(ctx context.Context) error {
	client := ec2.New(ec2.Options{Region: region})

	_, err := client.StartInstances(ctx, &ec2.StartInstancesInput{
		InstanceIds: []string{instance_id},
	})
	if err != nil {
		log.Println("could not start instance")
		log.Println(err)
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
