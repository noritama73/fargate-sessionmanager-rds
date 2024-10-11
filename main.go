package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/noritama73/fargate-sessionmanager-rds/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln(err)
	}

	profile := os.Getenv("PROFILE")
	cluster := os.Getenv("CLUSTER")
	serviceName := os.Getenv("SERVICE")
	rdsCluster := os.Getenv("RDS_CLUSTER")

	ctx := context.Background()

	cfg, err := service.NewSharedConfigProfile(ctx, profile)
	if err != nil {
		log.Fatalln(err)
	}

	ecsClient, err := service.NewECSClient(ctx, cfg)
	if err != nil {
		log.Fatalln(err)
	}

	taskID, err := ecsClient.ResolveTaskID(ctx, cluster, serviceName)
	if err != nil {
		log.Fatalln(err)
	}

	rdsClient, err := service.NewRDSClient(ctx, cfg)
	if err != nil {
		log.Fatalln(err)
	}

	rdsEndpoint, err := rdsClient.GetClusterEndpoint(ctx, rdsCluster)
	if err != nil {
		log.Fatalln(err)
	}

	target, err := ecsClient.GetSessionTarget(ctx, cluster, taskID)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("aws ssm start-session --target %s --profile %s --region ap-northeast-1 --document-name AWS-StartPortForwardingSessionToRemoteHost --parameters '{\"host\":[\"%s\"],\"portNumber\":[\"3306\"], \"localPortNumber\":[\"3306\"]}'\n",
		target, profile, rdsEndpoint)

	// ssmClient, err := service.NewSSMClient(ctx, cfg)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// sessionID, err := ssmClient.StartSession(ctx, target, rdsEndpoint, "3306")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// log.Printf("Session started: %s\n", sessionID)

	// tty, err := tty.Open()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer tty.Close()

	// log.Println("Press any key to terminate the session")

	// if _, err = tty.ReadRune(); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := ssmClient.TerminateSession(ctx, sessionID); err != nil {
	// 	log.Fatalln(err)
	// }

	// log.Println("Session terminated")
}
