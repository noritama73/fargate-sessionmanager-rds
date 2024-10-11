package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type ssmClient struct {
	svc *ssm.Client
}

type SSMService interface {
	StartSession(ctx context.Context, target, host, port string) (string, error)
	TerminateSession(ctx context.Context, sessionID string) error
}

func NewSSMClient(ctx context.Context, cfg aws.Config) (SSMService, error) {
	svc := ssm.NewFromConfig(cfg)
	return &ssmClient{svc: svc}, nil
}

func (c *ssmClient) StartSession(ctx context.Context, target, host, port string) (string, error) {
	input := &ssm.StartSessionInput{
		Target: aws.String(target),
		Parameters: map[string][]string{
			"host":            {host},
			"portNumber":      {"3306"},
			"localPortNumber": {port},
		},
		DocumentName: aws.String("AWS-StartPortForwardingSessionToRemoteHost"),
	}

	output, err := c.svc.StartSession(ctx, input)
	if err != nil {
		return "", ServiceError{Code: ErrUnexpected, OriginalError: err}
	}

	return *output.SessionId, nil
}

func (c *ssmClient) TerminateSession(ctx context.Context, sessionID string) error {
	input := &ssm.TerminateSessionInput{
		SessionId: aws.String(sessionID),
	}

	_, err := c.svc.TerminateSession(ctx, input)
	if err != nil {
		return ServiceError{Code: ErrUnexpected, OriginalError: err}
	}

	return nil
}
