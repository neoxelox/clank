package brevo

import (
	"context"

	"github.com/neoxelox/errors"
	"github.com/neoxelox/kit"
	"github.com/neoxelox/kit/util"

	brevo "github.com/getbrevo/brevo-go/lib"

	"backend/pkg/config"
)

var (
	ErrBrevoServiceGeneric  = errors.New("brevo service failed")
	ErrBrevoServiceTimedOut = errors.New("brevo service timed out")
)

type BrevoService struct {
	config   config.Config
	observer *kit.Observer
	client   *brevo.APIClient
}

func NewBrevoService(observer *kit.Observer, config config.Config) *BrevoService {
	clientConfig := brevo.NewConfiguration()
	clientConfig.AddDefaultHeader("Accept", "application/json")
	clientConfig.AddDefaultHeader("Api-Key", config.Brevo.APIKey)

	client := brevo.NewAPIClient(clientConfig)

	return &BrevoService{
		config:   config,
		observer: observer,
		client:   client,
	}
}

type BrevoServiceSendEmailParams struct {
	Receivers []string
	Subject   string
	Body      string
}

func (self *BrevoService) SendEmail(ctx context.Context, params BrevoServiceSendEmailParams) error {
	if self.config.Service.Environment != kit.EnvProduction {
		self.observer.Infof(ctx, "Sent email '%s' to %v", params.Subject, params.Receivers)
		return nil
	}

	receivers := make([]brevo.SendSmtpEmailTo, 0, len(params.Receivers))
	for _, receiver := range params.Receivers {
		receivers = append(receivers, brevo.SendSmtpEmailTo{Email: receiver})
	}

	_, response, err := self.client.TransactionalEmailsApi.SendTransacEmail(ctx, brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Email: self.config.Brevo.SenderEmail,
			Name:  self.config.Brevo.SenderName,
		},
		ReplyTo: &brevo.SendSmtpEmailReplyTo{
			Email: self.config.Brevo.ReplierEmail,
			Name:  self.config.Brevo.ReplierName,
		},
		To:          receivers,
		Subject:     params.Subject,
		HtmlContent: params.Body,
	})
	defer response.Body.Close()
	if err != nil {
		return ErrBrevoServiceGeneric.Raise().Cause(err)
	}

	return nil
}

func (self *BrevoService) Close(ctx context.Context) error {
	err := util.Deadline(ctx, func(exceeded <-chan struct{}) error {
		// Dummy log in order to mantain consistency although Brevo has no close() method
		self.observer.Info(ctx, "Closing Brevo service")
		self.observer.Info(ctx, "Closed Brevo service")

		return nil
	})
	if err != nil {
		if util.ErrDeadlineExceeded.Is(err) {
			return ErrBrevoServiceTimedOut.Raise().Cause(err)
		}

		return err
	}

	return nil
}
