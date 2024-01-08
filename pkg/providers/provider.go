package providers

import (
	"context"

	"glide/pkg/routers/health"

	"glide/pkg/api/schemas"
)

// LangModelProvider defines an interface a provider should fulfill to be able to serve language chat requests
type LangModelProvider interface {
	Provider() string
	Chat(ctx context.Context, request *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error)
}

// LangModel
type LangModel struct {
	modelID     string
	client      LangModelProvider
	rateLimit   *health.RateLimitTracker
	errorBudget *health.TokenBucket // TODO: centralize provider API health tracking in the registry
}

func NewLangModel(modelID string, client LangModelProvider, budget health.ErrorBudget) *LangModel {
	return &LangModel{
		modelID:     modelID,
		client:      client,
		rateLimit:   health.NewRateLimitTracker(),
		errorBudget: health.NewTokenBucket(budget.TimePerTokenMicro(), budget.Budget()),
	}
}

func (m *LangModel) ID() string {
	return m.modelID
}

func (m *LangModel) Provider() string {
	return m.client.Provider()
}

func (m *LangModel) Healthy() bool {
	return !m.rateLimit.Limited() && m.errorBudget.HasTokens()
}

func (m *LangModel) Chat(ctx context.Context, request *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error) {
	resp, err := m.client.Chat(ctx, request)

	if err == nil {
		// successful response
		resp.Model = m.modelID

		return resp, err
	}

	// TODO: track all availability issues

	return resp, err
}
