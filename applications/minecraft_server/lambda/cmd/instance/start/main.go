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

func Start(ctx context.Context) (string, error) {
	instance, err := serverInstance(ctx)
	if err != nil {
		return "", err
	}

	switch *instance.State.Code {
	case 0: // pending
		fallthrough
	case 32: // shutting-down
		fallthrough
	case 64: // stopping
		fallthrough
	case 48: // terminated
		log.Printf("instance is already %s", instance.State.Name)
		return fmt.Sprintf("server is already %s", instance.State.Name), nil
	default:
		client := ec2.New(ec2.Options{Region: region})
		_, err = client.StartInstances(ctx, &ec2.StartInstancesInput{
			InstanceIds: []string{instanceId},
		})
		if err != nil {
			log.Println("could not start instance")
			log.Println(err)
			return "", err
		}
		return "ok", nil
	}
}

func main() {
	lambda.Start(Start)
}
