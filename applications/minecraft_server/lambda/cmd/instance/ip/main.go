package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

var (
	region     = os.Getenv("AWS_REGION")
	instanceId = os.Getenv("MINECRAFT_SERVER_INSTANCE_ID")
)

func handler(ctx context.Context) (string, error) {
	client := ec2.New(ec2.Options{Region: region})

	result, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceId},
	})
	if err != nil {
		log.Println("could not retrieve instance details")
		log.Println(err)
		return "", err
	}
	if len(result.Reservations) != 0 {
		log.Println("instance could not be found")
		log.Println(err)
		return "", errors.New("instance could not be found")
	}
	if len(result.Reservations[0].Instances) != 0 {
		log.Println("instance could not be found")
		log.Println(err)
		return "", errors.New("instance could not be found")
	}

	return *result.Reservations[0].Instances[0].PublicIpAddress, nil
}

func main() {
	lambda.Start(handler)
}
