# go-reacjilator
Translate Slack message with flag emoji(e.g. :jp: :us: :uk:) reaction. :smile: Use only AWS Products(AWS Lambda functions in Go, Amazon Comprehend, Amazon Translate)

## Description
This is implemented [slack/reacjilator](https://github.com/slackapi/reacjilator) with AWS Products. ([AWS Lambda functions in Go](https://github.com/aws/aws-lambda-go), [Amazon Comprehend](https://aws.amazon.com/comprehend/), [Amazon Translate](https://aws.amazon.com/translate/))

You can easily translate Slack messages with flag emojiflag emoji(e.g. :jp: :us: :uk:) reaction.

The translated messages is sent back to the thread, so it does not interfere with other conversations.

## Features
- It works easily as Slack bot.
- Use only AWS Products.
- Of course, even when using other translation API you can easily respond.
- Since it is a serverless configuration, it is easy to build and deploy.

## Requirement
- golang 1.10+
- Packages in use
	- aws/aws-sdk-go: AWS SDK for the Go programming language.
		- https://github.com/aws/aws-sdk-go
	- kelseyhightower/envconfig: Golang library for managing configuration data from environment variables
		- https://github.com/kelseyhightower/envconfig
	- nlopes/slack: Slack API in Go
		- https://github.com/nlopes/slack
			- ***If the following pull requests are not merged it may not work with "invalid_limit" errors.***
				- Fixed "invalid_limit" error occurs. by uchimanajet7 · Pull Request #272 · nlopes/slack
					- https://github.com/nlopes/slack/pull/272

- Services in use
	- Amazon API Gateway 
		- https://aws.amazon.com/api-gateway/
	- AWS Lambda – Serverless Compute - Amazon Web Services
		- https://aws.amazon.com/lambda/
	- Amazon Translate – Neural Machine Translation - AWS
		- https://aws.amazon.com/translate/
	- Amazon Comprehend - Natural Language Processing (NLP) and Machine Learning (ML)
		- https://aws.amazon.com/comprehend/
	- Slack API | Slack 
		- https://api.slack.com/

- Tools in use
	- AWS Command Line Interface
		- https://aws.amazon.com/cli/
	- awslabs/serverless-application-model: AWS Serverless Application Model (AWS SAM) prescribes rules for expressing Serverless applications on AWS.
		- https://github.com/awslabs/serverless-application-model

## Usage
Install this bot on "Slack".

If you react to the message with the emoji of the flag, this bot translate the original message and post it under the message thread.

However, whether it can be translated depends on the product used for translation.

### Demo
![2017-12-17 18_11_16](https://user-images.githubusercontent.com/6448792/34078201-01eec178-e359-11e7-8494-17d044371c5f.gif)


## Installation

1. Refer to the following document and install this bot in "Slack".

	- slackapi/reacjilator: A translation bot that translates a message when a user reacted with an emoji 🇨🇳 🇮🇹 🇹🇭 🇫🇷 🇯🇵 🇮🇳 🇺🇸 🇧🇬 🇹🇼 🇦🇪 🇰🇷
		- https://github.com/slackapi/reacjilator
	- [Japanese] Developing a bot for your workspace 翻訳ボットを作る!
		- https://www.slideshare.net/tomomi/japanese-developing-a-bot-for-your-workspace-82133038

1. If you build from source yourself.

```sh
$ go get github.com/uchimanajet7/go-reacjilator
$ cd $GOPATH/src/github.com/uchimanajet7/go-reacjilator
$ go build
```

1. Build and Deploy bot using `AWS CLI` and `AWS SAM`.

*For deployment you need to be able to run the AWS CLI and prepare the AWS S3 resources necessary for AWS Lambda's operation in advance.*

- Deploying Lambda-based Applications - AWS Lambda 
	- https://docs.aws.amazon.com/lambda/latest/dg/deploying-lambda-apps.html

```sh
$ GOARCH=amd64 GOOS=linux go build -v -o build/go-reacjilator
$ aws cloudformation package \
    --template-file template.yml \
    --s3-bucket <YOUR_BUCKET_NAME> \
    --s3-prefix go-reacjilator \
    --output-template-file .template.yml
$ export SLACK_TOKEN=<YOUR_SLACK_TOKEN>
$ export SLACK_VERIFICATION_TOKEN=<YOUR_ SLACK_VERIFICATION_TOKEN>
$ export SLACK_CHANNEL_ID=<YOUR_SLACK_CHANNEL_ID>
$ aws cloudformation deploy \
    --template-file .template.yml \
    --stack-name go-reacjilator \
    --capabilities CAPABILITY_IAM \
    --parameter-overrides "SlackToken=$SLACK_TOKEN" "SlackVerificationToken=$SLACK_VERIFICATION_TOKEN" "SlackChannelID=$SLACK_CHANNEL_ID"
```

- About setting items
	- `s3-bucket <YOUR_BUCKET_NAME>`: **required**
		- Specify the bucket name of AWS S3 prepared in advance.
	- `SLACK_TOKEN=<YOUR_SLACK_TOKEN>`: **required**
		- Specify the token to use "Slack API".
	- `SLACK_VERIFICATION_TOKEN=<YOUR_ SLACK_VERIFICATION_TOKEN>`: **required**
		- Specify the verification token to use "Slack API".
	- `SLACK_CHANNEL_ID=<YOUR_SLACK_CHANNEL_ID>`: optional
		- Specify slack channel ID to allow.
			- *When you want to use channel filter*
				- If you want to limit the channels that respond to flag emoji reactions,  you need to specify one allowed channel.
				- Specify the `channel ID`.
				- Refer to the following for how to check the `channel ID`. 
					- Slack | APIに使う「チャンネルID」を取得する方法 - Qiita 
						- https://qiita.com/Yinaura/items/bd28c7b9ef614696fb7e


## Author
[uchimanajet7](https://github.com/uchimanajet7)

## Licence
[MIT License](https://github.com/uchimanajet7/go-reacjilator/blob/master/LICENSE)
