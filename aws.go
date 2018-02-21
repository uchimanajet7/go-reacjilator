package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	tl "github.com/aws/aws-sdk-go/service/translate"
)

// AwsClient is
type AwsClient struct {
	comprehendClient *comprehend.Comprehend
	translateClient  *tl.Translate
	languageCode     *map[string]string
}

// see also
// https://docs.aws.amazon.com/translate/latest/dg/how-it-works.html
var languageCodeMap = map[string]string{
	"Arabic":             "ar",
	"Chinese Simplified": "zh",
	"French":             "fr",
	"German":             "de",
	"Portuguese":         "pt",
	"Spanish":            "es",
	"English":            "en",
}

// GetLanguageCode is
func GetLanguageCode(language string) string {
	return languageCodeMap[language]
}

// NewAwsClient is
func NewAwsClient() *AwsClient {
	sess := session.Must(session.NewSession())
	client := &AwsClient{}
	client.comprehendClient = comprehend.New(sess, aws.NewConfig().WithRegion("us-west-2"))
	client.translateClient = tl.New(sess, aws.NewConfig().WithRegion("us-west-2"))

	return client
}

func (c *AwsClient) detectLanguageCode(text string) (string, error) {
	input := &comprehend.BatchDetectDominantLanguageInput{}
	input.SetTextList([]*string{&text})

	output, err := c.comprehendClient.BatchDetectDominantLanguage(input)
	if err != nil {
		log.Println("[Error] failed to aws comprehend detect language: ", err)
		return "", err
	}

	code := ""
	for _, i := range output.ResultList {
		for _, j := range i.Languages {
			if *j.LanguageCode != "" {
				code = *j.LanguageCode
			}
			break
		}
	}

	return code, nil
}

func (c *AwsClient) translate(text string, source string, target string) (string, error) {
	input := &tl.TextInput{}
	input.SetSourceLanguageCode(source)
	input.SetTargetLanguageCode(target)
	input.SetText(text)

	output, err := c.translateClient.Text(input)
	if err != nil {
		log.Println("[Error] failed to aws translate translation message: ", err)
		return "", err
	}

	return *output.TranslatedText, nil
}
