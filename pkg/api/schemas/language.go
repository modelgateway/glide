package schemas

// UnifiedChatRequest defines Glide's Chat Request Schema unified across all language models
type UnifiedChatRequest struct {
	Message        ChatMessage   `json:"message"`
	MessageHistory []ChatMessage `json:"messageHistory"`
}

func NewChatFromStr(message string) *UnifiedChatRequest {
	return &UnifiedChatRequest{
		Message: ChatMessage{
			"human",
			message,
			"roma",
		},
	}
}

// UnifiedChatResponse defines Glide's Chat Response Schema unified across all language models
type UnifiedChatResponse struct {
	ID            string           `json:"id,omitempty"`
	Created       int              `json:"created,omitempty"`
	Provider      string           `json:"provider,omitempty"`
	Router        string           `json:"router,omitempty"`
	Model         string           `json:"model,omitempty"`
	Cached        bool             `json:"cached,omitempty"`
	ModelResponse ProviderResponse `json:"modelResponse,omitempty"`
}

// ProviderResponse is the unified response from the provider.

type ProviderResponse struct {
	ResponseID map[string]string `json:"responseId,omitempty"`
	Message    ChatMessage       `json:"message"`
	TokenCount TokenCount        `json:"tokenCount"`
}

type TokenCount struct {
	PromptTokens   float64 `json:"promptTokens"`
	ResponseTokens float64 `json:"responseTokens"`
	TotalTokens    float64 `json:"totalTokens"`
}

// ChatMessage is a message in a chat request.
type ChatMessage struct {
	// The role of the author of this message. One of system, user, or assistant.
	Role string `json:"role"`
	// The content of the message.
	Content string `json:"content"`
	// The name of the author of this message. May contain a-z, A-Z, 0-9, and underscores,
	// with a maximum length of 64 characters.
	Name string `json:"name,omitempty"`
}

// OpenAI Chat Response
// TODO: Should this live here?
type OpenAIChatCompletion struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int      `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     float64 `json:"prompt_tokens"`
	CompletionTokens float64 `json:"completion_tokens"`
	TotalTokens      float64 `json:"total_tokens"`
}

// Cohere Chat Response
type CohereChatCompletion struct {
	Text          string                 `json:"text"`
	GenerationID  string                 `json:"generation_id"`
	ResponseID    string                 `json:"response_id"`
	TokenCount    CohereTokenCount       `json:"token_count"`
	Citations     []Citation             `json:"citations"`
	Documents     []Documents            `json:"documents"`
	SearchQueries []SearchQuery          `json:"search_queries"`
	SearchResults []SearchResults        `json:"search_results"`
	Meta          Meta                   `json:"meta"`
	ToolInputs    map[string]interface{} `json:"tool_inputs"`
}

type CohereTokenCount struct {
	PromptTokens   float64 `json:"prompt_tokens"`
	ResponseTokens float64 `json:"response_tokens"`
	TotalTokens    float64 `json:"total_tokens"`
	BilledTokens   float64 `json:"billed_tokens"`
}

type Meta struct {
	APIVersion struct {
		Version string `json:"version"`
	} `json:"api_version"`
	BilledUnits struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"billed_units"`
}

type Citation struct {
	Start      int      `json:"start"`
	End        int      `json:"end"`
	Text       string   `json:"text"`
	DocumentID []string `json:"documentId"`
}

type Documents struct {
	ID   string            `json:"id"`
	Data map[string]string `json:"data"` // TODO: This needs to be updated
}

type SearchQuery struct {
	Text         string `json:"text"`
	GenerationID string `json:"generationId"`
}

type SearchResults struct {
	SearchQuery []SearchQueryObject  `json:"searchQuery"`
	Connectors  []ConnectorsResponse `json:"connectors"`
	DocumentID  []string             `json:"documentId"`
}

type SearchQueryObject struct {
	Text         string `json:"text"`
	GenerationID string `json:"generationId"`
}

type ConnectorsResponse struct {
	ID              string            `json:"id"`
	UserAccessToken string            `json:"user_access_token"`
	ContOnFail      string            `json:"continue_on_failure"`
	Options         map[string]string `json:"options"`
}
