package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

var (
	region     = os.Getenv("AWS_REGION")
	instanceId = os.Getenv("MINECRAFT_SERVER_INSTANCE_ID")
)

func serverInstance(ctx context.Context) (types.Instance, error) {
	client := ec2.New(ec2.Options{Region: region})
	result, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceId},
	})
	if err != nil {
		log.Println("error retriving instance details")
		log.Println(err)
		return types.Instance{}, err
	}
	if len(result.Reservations) != 0 {
		log.Println("reservation not found")
		return types.Instance{}, errors.New("server could not be found")
	}
	if len(result.Reservations[0].Instances) != 0 {
		log.Println("reservation present but instance not found")
		return types.Instance{}, errors.New("server could not be found")
	}

	return result.Reservations[0].Instances[0], nil
}

func ip(ctx context.Context) (string, error) {
	instance, err := serverInstance(ctx)
	if err != nil {
		return "", err
	}

	switch *instance.State.Code {
	case 16:
		if *instance.PublicIpAddress == "" {
			log.Println("instance is running but publicIpAddr is not present")
			msg := "server is running, but ip address is not available, please check again later"
			return msg, nil
		}
		return *instance.PublicIpAddress, nil
	default:
		log.Printf("instance is %s", instance.State.Name)
		msg := fmt.Sprintf("cannot retrieve ip (server is %s)", instance.State.Name)
		return msg, nil
	}
}

func main() {
	lambda.Start(ip)
}
