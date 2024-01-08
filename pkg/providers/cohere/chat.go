package cohere

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"glide/pkg/providers/clients"

	"glide/pkg/api/schemas"
	"go.uber.org/zap"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatHistory struct {
	Role    string `json:"role"`
	Message string `json:"message"`
	User    string `json:"user,omitempty"`
}

// ChatRequest is a request to complete a chat completion..
type ChatRequest struct {
	Model             string        `json:"model"`
	Message           string        `json:"message"`
	Temperature       float64       `json:"temperature,omitempty"`
	PreambleOverride  string        `json:"preamble_override,omitempty"`
	ChatHistory       []ChatHistory `json:"chat_history,omitempty"`
	ConversationID    string        `json:"conversation_id,omitempty"`
	PromptTruncation  string        `json:"prompt_truncation,omitempty"`
	Connectors        []string      `json:"connectors,omitempty"`
	SearchQueriesOnly bool          `json:"search_queries_only,omitempty"`
	CitiationQuality  string        `json:"citiation_quality,omitempty"`

	// Stream            bool                `json:"stream,omitempty"`
}

type Connectors struct {
	ID              string            `json:"id"`
	UserAccessToken string            `json:"user_access_token"`
	ContOnFail      string            `json:"continue_on_failure"`
	Options         map[string]string `json:"options"`
}

// NewChatRequestFromConfig fills the struct from the config. Not using reflection because of performance penalty it gives
func NewChatRequestFromConfig(cfg *Config) *ChatRequest {
	return &ChatRequest{
		Model:             cfg.Model,
		Temperature:       cfg.DefaultParams.Temperature,
		PreambleOverride:  cfg.DefaultParams.PreambleOverride,
		ChatHistory:       cfg.DefaultParams.ChatHistory,
		ConversationID:    cfg.DefaultParams.ConversationID,
		PromptTruncation:  cfg.DefaultParams.PromptTruncation,
		Connectors:        cfg.DefaultParams.Connectors,
		SearchQueriesOnly: cfg.DefaultParams.SearchQueriesOnly,
		CitiationQuality:  cfg.DefaultParams.CitiationQuality,
	}
}

// Chat sends a chat request to the specified cohere model.
func (c *Client) Chat(ctx context.Context, request *schemas.UnifiedChatRequest) (*schemas.UnifiedChatResponse, error) {
	// Create a new chat request
	chatRequest := c.createChatRequestSchema(request)

	chatResponse, err := c.doChatRequest(ctx, chatRequest)
	if err != nil {
		return nil, err
	}

	if len(chatResponse.ModelResponse.Message.Content) == 0 {
		return nil, ErrEmptyResponse
	}

	return chatResponse, nil
}

func (c *Client) createChatRequestSchema(request *schemas.UnifiedChatRequest) *ChatRequest {
	// TODO: consider using objectpool to optimize memory allocation
	chatRequest := c.chatRequestTemplate // hoping to get a copy of the template
	chatRequest.Message = request.Message.Content

	// Build the Cohere specific ChatHistory
	if len(request.MessageHistory) > 0 {
		chatRequest.ChatHistory = make([]ChatHistory, len(request.MessageHistory))
		for i, message := range request.MessageHistory {
			chatRequest.ChatHistory[i] = ChatHistory{
				// Copy the necessary fields from message to ChatHistory
				// For example, if ChatHistory has a field called "Text", you can do:
				Role:    message.Role,
				Message: message.Content,
				User:    "",
			}
		}
	}

	return chatRequest
}

func (c *Client) doChatRequest(ctx context.Context, payload *ChatRequest) (*schemas.UnifiedChatResponse, error) {
	// Build request payload
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal cohere chat request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.chatURL, bytes.NewBuffer(rawPayload))
	if err != nil {
		return nil, fmt.Errorf("unable to create cohere chat request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+string(c.config.APIKey))
	req.Header.Set("Content-Type", "application/json")

	// TODO: this could leak information from messages which may not be a desired thing to have
	c.telemetry.Logger.Debug(
		"cohere chat request",
		zap.String("chat_url", c.chatURL),
		zap.Any("payload", payload),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send cohere chat request: %w", err)
	}

	defer resp.Body.Close() // TODO: handle this error

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			c.telemetry.Logger.Error("failed to read cohere chat response", zap.Error(err))
		}

		// TODO: Handle failure conditions
		// TODO: return errors
		c.telemetry.Logger.Error(
			"cohere chat request failed",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(bodyBytes)),
			zap.Any("headers", resp.Header),
		)

		return nil, clients.ErrProviderUnavailable
	}

	// Read the response body into a byte slice
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.telemetry.Logger.Error("failed to read cohere chat response", zap.Error(err))
		return nil, err
	}

	// Parse the response JSON
	var responseJSON map[string]interface{}

	err = json.Unmarshal(bodyBytes, &responseJSON)
	if err != nil {
		c.telemetry.Logger.Error("failed to parse cohere chat response", zap.Error(err))
		return nil, err
	}

	// Parse the response JSON
	var cohereCompletion schemas.CohereChatCompletion

	err = json.Unmarshal(bodyBytes, &cohereCompletion)
	if err != nil {
		c.telemetry.Logger.Error("failed to parse openai chat response", zap.Error(err))
		return nil, err
	}

	// Map response to UnifiedChatResponse schema
	response := schemas.UnifiedChatResponse{
		ID:       cohereCompletion.ResponseID,
		Created:  int(time.Now().UTC().Unix()), // Cohere doesn't provide this
		Provider: providerName,
		Router:   "chat",          // TODO: this will be the router used
		Model:    "command-light", // TODO: this needs to come from config or router as Cohere doesn't provide this
		Cached:   false,
		ModelResponse: schemas.ProviderResponse{
			ResponseID: map[string]string{
				"generationId": cohereCompletion.GenerationID,
				"responseId":   cohereCompletion.ResponseID,
			},
			Message: schemas.ChatMessage{
				Role:    "model", // TODO: Does this need to change?
				Content: cohereCompletion.Text,
				Name:    "",
			},
			TokenCount: schemas.TokenCount{
				PromptTokens:   cohereCompletion.TokenCount.PromptTokens,
				ResponseTokens: cohereCompletion.TokenCount.ResponseTokens,
				TotalTokens:    cohereCompletion.TokenCount.TotalTokens,
			},
		},
	}

	return &response, nil
}
