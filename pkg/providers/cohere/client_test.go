package cohere

import (
	"context"
	"fmt"
	"testing"

	"glide/pkg/api/schemas"

	"glide/pkg/telemetry"

	"github.com/stretchr/testify/require"
)

func TestOpenAIClient_ChatRequest(t *testing.T) {
	ctx := context.Background()
	cfg := DefaultConfig()

	client, err := NewClient(cfg, telemetry.NewTelemetryMock())
	require.NoError(t, err)

	request := schemas.UnifiedChatRequest{Message: schemas.ChatMessage{
		Role:    "human",
		Content: "What's the biggest animal?",
	}}

	response, err := client.Chat(ctx, &request)
	require.NoError(t, err)

	fmt.Print(response)

	// require.Equal(t, "chatcmpl-123", response.ID)
}
