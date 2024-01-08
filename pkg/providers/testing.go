package providers

import (
	"context"

	"glide/pkg/api/schemas"
)

type ResponseMock struct {
	Msg string
	Err *error
}

func (m *ResponseMock) Resp() *schemas.UnifiedChatResponse {
	return &schemas.UnifiedChatResponse{
		ID: "rsp0001",
		ModelResponse: schemas.ProviderResponse{
			ResponseID: map[string]string{
				"ID": "0001",
			},
			Message: schemas.ChatMessage{
				Content: m.Msg,
			},
		},
	}
}

type ProviderMock struct {
	idx       int
	responses []ResponseMock
}

func NewProviderMock(responses []ResponseMock) *ProviderMock {
	return &ProviderMock{
		idx:       0,
		responses: responses,
	}
}

func (c *ProviderMock) Chat(_ context.Context, _ *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error) {
	response := c.responses[c.idx]
	c.idx++

	if response.Err != nil {
		return nil, *response.Err
	}

	return response.Resp(), nil
}

func (c *ProviderMock) Provider() string {
	return "provider_mock"
}
