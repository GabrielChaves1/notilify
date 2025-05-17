package awsses

import (
	"GabrielChaves1/notilify/internal/clients/mailer"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SESClient struct {
	client          *ses.Client
	fromSourceEmail string
}

func NewSESClient(client *ses.Client, fromSourceEmail string) mailer.Client {
	return &SESClient{
		client:          client,
		fromSourceEmail: fromSourceEmail,
	}
}

func (s SESClient) SendMail(ctx context.Context, to, subject, text string) error {
	input := &ses.SendEmailInput{
		Source: aws.String(s.fromSourceEmail),
		Destination: &types.Destination{
			ToAddresses: []string{
				to,
			},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(subject),
			},
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(text),
				},
			},
		},
	}

	if _, err := s.client.SendEmail(ctx, input); err != nil {
		return err
	}

	return nil
}
