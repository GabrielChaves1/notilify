package awssns

import (
	notifications "GabrielChaves1/notilify/internal/clients/notification"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSClient struct {
	client *sns.Client
}

func NewSNSClient(client *sns.Client) notifications.Client {
	return &SNSClient{
		client: client,
	}
}

func (s *SNSClient) SendSMS(ctx context.Context, phoneNumber, message string) error {
	input := &sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(phoneNumber),
	}

	_, err := s.client.Publish(ctx, input)
	if err != nil {
		return err
	}

	fmt.Println("Mensagem enviada com sucesso!")
	return nil
}
