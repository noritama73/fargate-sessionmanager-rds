package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
)

func NewSharedConfigProfile(ctx context.Context, profile string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithSharedConfigProfile(profile),
		config.WithAssumeRoleCredentialOptions(func(options *stscreds.AssumeRoleOptions) {
			options.TokenProvider = func() (string, error) {
				return stscreds.StdinTokenProvider()
			}
		}),
	)
	if err != nil {
		return aws.Config{}, ServiceError{Code: ErrUnexpected, OriginalError: err}
	}

	return cfg, nil
}
