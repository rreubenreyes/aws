package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

var (
	region      = os.Getenv("AWS_REGION")
	instance_id = os.Getenv("MINECRAFT_SERVER_INSTANCE_ID")
)

func handler(ctx context.Context) (string, error) {
	client := ec2.New(ec2.Options{Region: region})

	result, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instance_id},
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

	launchtime := *result.Reservations[0].Instances[0].LaunchTime
	uptime := time.Since(launchtime)

	return fmt.Sprintf("%s (%f hours)", launchtime.Format(time.UnixDate), uptime.Hours()), nil
}

func main() {
	lambda.Start(handler)
}
