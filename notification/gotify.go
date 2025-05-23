package notification

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ahobsonsayers/twigots"
	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	"github.com/gotify/go-api-client/v2/models"
)

type GotifyClient struct {
	url   *url.URL
	token string

	client *client.GotifyREST
}

var _ Client = GotifyClient{}

func (g GotifyClient) SendTicketNotification(ticket twigots.TicketListing) error {
	notificationMessage, err := RenderMessage(ticket, WithFooter())
	if err != nil {
		return err
	}

	params := message.NewCreateMessageParams()
	params.Body = &models.MessageExternal{
		Title:   ticket.Event.Name,
		Message: notificationMessage,
		Extras: map[string]any{
			"client::display": map[string]any{
				"contentType": "text/markdown",
			},
			"client::notification": map[string]any{
				"click": map[string]any{
					"url": ticket.URL(),
				},
			},
		},
		Priority: 5,
	}

	_, err = g.client.Message.CreateMessage(
		params,
		auth.TokenAuth(g.token),
	)
	if err != nil {
		return err
	}

	return nil
}

type GotifyConfig struct {
	Url   string `json:"url"`
	Token string `json:"token"`
}

func NewGotifyClient(config GotifyConfig) (GotifyClient, error) {
	gotifyUrl, err := url.Parse(config.Url)
	if err != nil {
		return GotifyClient{}, fmt.Errorf("failed to parse gotify url: %v", err)
	}

	return GotifyClient{
		url:   gotifyUrl,
		token: config.Token, // TODO validate this token?

		client: gotify.NewClient(gotifyUrl, &http.Client{}),
	}, nil
}
