package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type ecsClient struct {
	svc *ecs.Client
}

type ECSService interface {
	GetSessionTarget(ctx context.Context, cluster string, taskID string) (string, error)
	ResolveTaskID(ctx context.Context, cluster string, service string) (string, error)
}

func NewECSClient(ctx context.Context, cfg aws.Config) (ECSService, error) {
	svc := ecs.NewFromConfig(cfg)
	return &ecsClient{svc: svc}, nil
}

func (c *ecsClient) GetSessionTarget(ctx context.Context, cluster string, taskID string) (string, error) {
	input := &ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks:   []string{taskID},
	}

	output, err := c.svc.DescribeTasks(ctx, input)
	if err != nil {
		return "", err
	}

	if len(output.Tasks) == 0 {
		return "", ServiceError{Code: ErrNotFound, OriginalError: errors.New("task not found")}
	}

	if len(output.Tasks[0].Containers) == 0 {
		return "", ServiceError{Code: ErrNotFound, OriginalError: errors.New("container not found")}
	}

	targetID := fmt.Sprintf("ecs:%s_%s_%s", cluster, taskID, *output.Tasks[0].Containers[0].RuntimeId)

	return targetID, nil
}

func (c *ecsClient) ResolveTaskID(ctx context.Context, cluster string, service string) (string, error) {
	input := &ecs.ListTasksInput{
		Cluster:     aws.String(cluster),
		ServiceName: aws.String(service),
	}

	output, err := c.svc.ListTasks(ctx, input)
	if err != nil {
		return "", ServiceError{Code: ErrUnexpected, OriginalError: err}
	}

	if len(output.TaskArns) == 0 {
		return "", ServiceError{Code: ErrNotFound, OriginalError: errors.New("task not found")}
	}

	describeInput := &ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks:   output.TaskArns,
	}

	describeOutput, err := c.svc.DescribeTasks(ctx, describeInput)
	if err != nil {
		return "", err
	}

	if len(describeOutput.Tasks) == 0 {
		return "", ServiceError{Code: ErrNotFound, OriginalError: errors.New("task not found")}
	}

	if len(describeOutput.Tasks) != 1 {
		return "", ServiceError{Code: ErrInvalid, OriginalError: errors.New("multiple tasks found")}
	}

	arn := *describeOutput.Tasks[0].TaskArn
	taskID := strings.Split(arn, "/")[2]

	return taskID, nil
}
