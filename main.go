package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

type envConfig struct {
	SlackToken             string `envconfig:"SLACK_TOKEN" required:"true"`
	SlackVerificationToken string `envconfig:"SLACK_VERIFICATION_TOKEN" required:"true"`
	SlackChannelID         string `envconfig:"SLACK_CHANNEL_ID"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Println("[ERROR] failed to process env var: ", err)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	client := slack.New(env.SlackToken)
	slackClient := &SlackClient{
		client:            client,
		verificationToken: env.SlackVerificationToken,
		channelID:         env.SlackChannelID,
	}

	result, err := slackClient.handleEvent(request.Body)
	if err != nil {
		log.Println("[ERROR] Processing failed: ", err)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}
	if result == "" {
		result = `{"result": "ok"}`
	}

	return events.APIGatewayProxyResponse{Body: result, StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
