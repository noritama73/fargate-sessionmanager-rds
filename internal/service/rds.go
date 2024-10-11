package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type rdsClient struct {
	svc *rds.Client
}

type RDSService interface{
	GetClusterEndpoint(ctx context.Context, clusterName string) (string, error)
}

func NewRDSClient(ctx context.Context, cfg aws.Config) (RDSService, error) {
	svc := rds.NewFromConfig(cfg)
	return &rdsClient{svc: svc}, nil
}

func (c *rdsClient) GetClusterEndpoint(ctx context.Context, clusterName string) (string, error) {
	input := &rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(clusterName),
	}

	output, err := c.svc.DescribeDBClusters(ctx, input)
	if err != nil {
		return "", ServiceError{Code: ErrUnexpected, OriginalError: err}
	}

	if len(output.DBClusters) == 0 {
		return "", ServiceError{Code: ErrNotFound, OriginalError: err}
	}

	return *output.DBClusters[0].Endpoint, nil
}
