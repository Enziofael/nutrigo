package v1

import (
	"context"
	"fmt"

	httpclient "github.com/Enziofael/nutrigo/backend/pkg/HTTPclient"
	models "github.com/Enziofael/nutrigo/shared/models/v1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Client struct {
	inner *httpclient.Client
}

func New(baseURL string, cfg httpclient.ClientConfig) *Client {
	return &Client{inner: httpclient.NewClient(baseURL, cfg)}
}

func (c *Client) Close() error {
	return c.inner.Close()
}

func (c *Client) GetUser(ctx context.Context, u tgbotapi.Update) (*models.User, error) {
	tgID := u.Message.From.ID
	path := fmt.Sprintf("/v1/user/%d", tgID)

	var user models.User
	if err := c.inner.GET(ctx, path, &user); err != nil {
		if apiErr, ok := err.(*httpclient.APIError); ok && apiErr.StatusCode == 404 {
			return nil, nil // пользователь не найден
		}
		return nil, fmt.Errorf("get user by telegram id %d: %w", tgID, err)
	}
	return &user, nil
}

func (c *Client) CreateUser(ctx context.Context, u tgbotapi.Update) (*models.User, error) {
	req := models.UserCreateRequest{
		TgID:  u.Message.From.ID,
		TgTag: u.Message.From.UserName,
	}
	path := "/v1/user/"

	var user models.User
	if err := c.inner.POST(ctx, path, req, &user); err != nil {
		return nil, &httpclient.APIError{Message: fmt.Sprintf("User creation failed:\n\t%s", err)}
	}
	return &user, nil
}
