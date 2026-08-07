package client

import (
	"context"
	"fmt"
	"net/http"

	httpclient "github.com/Enziofael/nutrigo/backend/pkg/HTTPclient"
	models "github.com/Enziofael/nutrigo/shared/models/v1"
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

// ============================================================
// User
// ============================================================

func (c *Client) CreateUser(ctx context.Context, tgID int64, tgTag string) (*models.User, error) {
	req := models.UserCreateRequest{
		TgID:  tgID,
		TgTag: tgTag,
	}

	path := "/v1/user/"

	var user models.User

	if err := c.inner.POST(ctx, path, req, &user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

func (c *Client) DeleteUser(ctx context.Context, tgID int64) error {
	path := fmt.Sprintf("/v1/user/%d", tgID)

	if err := c.inner.DELETE(ctx, path, nil); err != nil {
		if apiErr, ok := err.(*httpclient.APIError); ok && apiErr.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("delete user %d: %w", tgID, err)
	}
	return nil
}

func (c *Client) GetUser(ctx context.Context, tgID int64, tgTag string) (*models.User, error) {
	path := fmt.Sprintf("/v1/user/%d", tgID)

	var user models.User

	if err := c.inner.GET(ctx, path, &user); err != nil {
		if apiErr, ok := err.(*httpclient.APIError); ok && apiErr.StatusCode == http.StatusNotFound {
			return &models.User{
				TgID:  tgID,
				TgTag: tgTag,
			}, nil // пользователь не найден в базе
		}
		return nil, fmt.Errorf("get user %d: %w", tgID, err)
	}
	return &user, nil
}

func (c *Client) GetUserStrict(ctx context.Context, tgID int64) (*models.User, error) {
	path := fmt.Sprintf("/v1/user/%d", tgID)

	var user models.User

	if err := c.inner.GET(ctx, path, &user); err != nil {
		return nil, fmt.Errorf("get user %d: %w", tgID, err)
	}
	return &user, nil
}

func (c *Client) PatchUser(ctx context.Context, req models.UserPatchRequest) (*models.User, error) {
	path := fmt.Sprintf("/v1/user/%d", req.TgID)

	var user models.User

	if err := c.inner.PATCH(ctx, path, req, &user); err != nil {
		return nil, fmt.Errorf("patch user %d: %w", req.TgID, err)
	}
	return &user, nil
}

// PatchUserStatus – shortcut for PatchUser()
func (c *Client) PatchUserStatus(ctx context.Context, tgID int64, status string) (*models.User, error) {
	req := models.UserPatchRequest{
		TgID:   tgID,
		Status: &status,
	}
	return c.PatchUser(ctx, req)
}

// PatchUserContext – shortcut for PatchUser()
func (c *Client) PatchUserContext(ctx context.Context, tgID int64, context string, contextData models.ContextData) (*models.User, error) {
	req := models.UserPatchRequest{
		TgID:        tgID,
		Context:     &context,
		ContextData: &contextData,
	}
	return c.PatchUser(ctx, req)
}

// ============================================================
// УПРАЖНЕНИЯ (Exercise)
// ============================================================

func (c *Client) CreateExercise(ctx context.Context, tgID int64, name string) (*models.Exercise, error) {
	req := models.ExerciseCreateRequest{
		TgID: tgID,
		Name: name,
	}

	path := "/v1/exercise/"

	var exercise models.Exercise

	if err := c.inner.POST(ctx, path, req, &exercise); err != nil {
		return nil, fmt.Errorf("create exercise: %w", err)
	}
	return &exercise, nil
}

func (c *Client) DeleteExercise(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/v1/exercise/%d", id)

	if err := c.inner.DELETE(ctx, path, nil); err != nil {
		if apiErr, ok := err.(*httpclient.APIError); ok && apiErr.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("delete exercise %d: %w", id, err)
	}
	return nil
}

func (c *Client) GetExercise(ctx context.Context, id int64) (*models.Exercise, error) {
	path := fmt.Sprintf("/v1/exercise/%d", id)

	var exercise models.Exercise

	if err := c.inner.GET(ctx, path, &exercise); err != nil {
		return nil, fmt.Errorf("get exercise %d: %w", id, err)
	}
	return &exercise, nil
}

func (c *Client) ListExercises(ctx context.Context, tgID int64, offset, limit int, sortBy string, order string) (*[]models.Exercise, error) {
	path := fmt.Sprintf("/v1/exercise/u/%d?offset=%d&limit=%d&sort=%s&order=%s",
		tgID, offset, limit, sortBy, models.Order(order))

	var exercises []models.Exercise

	if err := c.inner.GET(ctx, path, &exercises); err != nil {
		return nil, fmt.Errorf("list exercises for user %d: %w", tgID, err)
	}
	return &exercises, nil
}

func (c *Client) CountExercises(ctx context.Context, tgID int64) (int, error) {
	path := fmt.Sprintf("/v1/exercise/u/%d/count", tgID)

	var Resp struct {
		Count int `json:"count"`
	}

	if err := c.inner.GET(ctx, path, &Resp); err != nil {
		return 0, fmt.Errorf("count exercises for user %d: %w", tgID, err)
	}
	return Resp.Count, nil
}

func (c *Client) PatchExercise(ctx context.Context, req models.ExercisePatchRequest) (*models.Exercise, error) {
	path := fmt.Sprintf("/v1/exercise/%d", req.ID)

	var exercise models.Exercise

	if err := c.inner.PATCH(ctx, path, req, &exercise); err != nil {
		return nil, fmt.Errorf("patch exercise %d: %w", req.ID, err)
	}
	return &exercise, nil
}

// PatchExerciseRating – shortcut for PatchExercise()
func (c *Client) PatchExerciseRating(ctx context.Context, id int64, rating int) (*models.Exercise, error) {
	req := models.ExercisePatchRequest{
		ID:     id,
		Rating: &rating,
	}
	return c.PatchExercise(ctx, req)
}

// PatchExerciseWeightUnitPreference – shortcut for PatchExercise()
func (c *Client) PatchExerciseWeightUnitPreference(ctx context.Context, id int64, unit models.WeightUnit) (*models.Exercise, error) {
	s := string(unit)
	req := models.ExercisePatchRequest{
		ID:         id,
		WeightUnit: &s,
	}
	return c.PatchExercise(ctx, req)
}

/*
func (c *Client) ListExerciseEntries(exerciseID int, ctx context.Context) (*[]models.ExerciseEntry, error) {
	path := fmt.Sprintf("/v1/exercise_entry/e/%d?offset=%d&limit=%d&sort=%s&order=%s", exerciseID)

}
*/
