package openai

import (
	"fmt"
	"net/http"
	"io"
	"bytes"
)

type OpenAiProviderConfig struct {
	Model            string           `json:"model" validate:"required,lowercase"`
	Messages         string           `json:"messages" validate:"required"` // does this need to be updated to []string?
	MaxTokens        int              `json:"max_tokens" validate:"omitempty,gte=0"`
	Temperature      int              `json:"temperature" validate:"omitempty,gte=0,lte=2"`
	TopP             int              `json:"top_p" validate:"omitempty,gte=0,lte=1"`
	N                int              `json:"n" validate:"omitempty,gte=1"`
	Stream           bool             `json:"stream" validate:"omitempty, boolean"`
	Stop             interface{}      `json:"stop"`
	PresencePenalty  int              `json:"presence_penalty" validate:"omitempty,gte=-2,lte=2"`
	FrequencyPenalty int              `json:"frequency_penalty" validate:"omitempty,gte=-2,lte=2"`
	LogitBias        *map[int]float64 `json:"logit_bias" validate:"omitempty"`
	User             interface{}      `json:"user"`
	Seed             interface{}      `json:"seed" validate:"omitempty,gte=0"`
	Tools            []string         `json:"tools"`
	ToolChoice       interface{}      `json:"tool_choice"`
	ResponseFormat   interface{}      `json:"response_format"`
}

type OpenAiClient struct {
	apiKey   string
	baseURL  string
	params OpenAiProviderConfig
	http     *http.Client
}

func NewOpenAiClient(apiKey string) *OpenAiClient {
	return &OpenAiClient{
		apiKey:   apiKey,
		baseURL:  "https://api.openai.com/v1",
		params:   OpenAiChatDefaultConfig(),
		http:     http.DefaultClient,
	}
}

var defaultMessage = `[
	{
	  "role": "system",
	  "content": "You are a helpful assistant."
	},
	{
	  "role": "user",
	  "content": "Hello!"
	}
  ]`
  
func OpenAiChatDefaultConfig() OpenAiProviderConfig {
	return OpenAiProviderConfig{
		Model:            "gpt-3.5-turbo",
		Messages:         defaultMessage,
		MaxTokens:        100,
		Temperature:      1,
		TopP:             1,
		N:                1,
		Stream:           false,
		Stop:             nil,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		LogitBias:        nil,
		User:             nil,
		Seed:             nil,
		Tools:            nil,
		ToolChoice:       nil,
		ResponseFormat:   nil,
	}
}
func (c *OpenAiClient) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

func (c *OpenAiClient) SetHTTPOpenAiClient(httpOpenAiClient *http.Client) {
	c.http = httpOpenAiClient
}

func (c *OpenAiClient) GetAPIKey() string {
	return c.apiKey
}

func (c *OpenAiClient) Get(endpoint string) (string, error) {
	// Implement the logic to make a GET request to the OpenAI API

	return "", nil
}

func (c *OpenAiClient) Post(endpoint string, payload []byte) (string, error) {
	// Implement the logic to make a POST request to the OpenAI API

	// Create the full URL
	url := c.baseURL + endpoint

	// Create a new request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send the request using http Client
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseBody), nil
}

// Add more methods to interact with OpenAI API

func main() {
	// Example usage of the OpenAI OpenAiClient
	OpenAiClient := NewOpenAiClient("YOUR_API_KEY")
	
	// Call methods on the OpenAiClient to interact with the OpenAI API
	// For example:
	payrload := []byte(`{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "Hello!"}]}`)
	response, err := OpenAiClient.Post("/chat", payrload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	fmt.Println("Response:", response)
}