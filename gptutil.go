package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	GPT_CURRENT_MODEL                = GPT_MODEL_GPT_35_TURBO
	GPT_MODEL_DALL_E_3               = "dall-e-3"
	GPT_MODEL_GPT_35_TURBO           = "gpt-3.5-turbo"
	GPT_MODEL_TEXT_EMBEDDING_ADA_002 = "text-embedding-ada-002"
	GPT_ROLE_USER                    = "user"
	GPT_ROLE_SYSTEM                  = "system"
	GPT_TEMPERATURE                  = 0.5
	GTP_ENCODING_FLOAT               = "float"
)

type GptError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

type GptMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GptCompletionsRequest struct {
	Model       string       `json:"model"`
	Messages    []GptMessage `json:"messages"`
	Temperature float64      `json:"temperature"`
}

type GptCompletionsResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens            int `json:"prompt_tokens"`
		CompletionTokens        int `json:"completion_tokens"`
		TotalTokens             int `json:"total_tokens"`
		CompletionTokensDetails struct {
			ReasoningTokens          int `json:"reasoning_tokens"`
			AcceptedPredictionTokens int `json:"accepted_prediction_tokens"`
			RejectedPredictionTokens int `json:"rejected_prediction_tokens"`
		} `json:"completion_tokens_details"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
		Index        int         `json:"index"`
	} `json:"choices"`
	Error *GptError `json:"error"`
}

type GptEmbeddingRequest struct {
	Input          string `json:"input"`
	Model          string `json:"model"`
	EncodingFormat string `json:"encoding_format"`
}

type GptEmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
	Error *GptError `json:"error"`
}

type GptImageRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type GptImageResponse struct {
	Created int `json:"created"`
	Data    []struct {
		URL string `json:"url"`
	} `json:"data"`
	Error *GptError `json:"error"`
}

func gptGetCompletions(secretProvider SecretProvider, question *Question) []*Answer {
	a := new(Answer)
	reqObj := GptCompletionsRequest{
		Model: GPT_CURRENT_MODEL,
		Messages: []GptMessage{
			{
				Content: question.Text,
				Role:    GPT_ROLE_SYSTEM,
			},
		},
	}
	buf, err := json.MarshalIndent(reqObj, "", "\t")
	if err != nil {
		a.Text = fmt.Sprintf("gpt completions error: %s", err.Error())
		return []*Answer{a}
	}
	url := "https://api.openai.com/v1/chat/completions"
	client := http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(buf))
	if err != nil {
		a.Text = fmt.Sprintf("gpt completions error: %s", err.Error())
		return []*Answer{a}
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+secretProvider.GetSecret("openai"))
	resp, err := client.Do(req)
	if err != nil {
		a.Text = fmt.Sprintf("gpt completions error: %s", err.Error())
		return []*Answer{a}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.Text = fmt.Sprintf("gpt completions error: %s", err.Error())
		return []*Answer{a}
	}
	//fmt.Printf("%s\n", string(body))
	var respObj GptCompletionsResponse
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		a.Text = fmt.Sprintf("gpt completions error: %s", err.Error())
		return []*Answer{a}
	}
	if respObj.Error != nil {
		a.Text = fmt.Sprintf("gpt completions error: %s", respObj.Error.Message)
		return []*Answer{a}
	}
	if len(respObj.Choices) == 0 {
		a.Text = "no completion for you"
		return []*Answer{a}

	}
	answers := []*Answer{}
	for _, r := range respObj.Choices {
		a.Text = r.Message.Content
		answers = append(answers, a)
		a = new(Answer)
	}
	return answers
}

func gptGetEmbedding(secretProvider SecretProvider, question *Question) (*Embedding, error) {
	reqObj := GptEmbeddingRequest{
		Input:          question.Text,
		Model:          GPT_MODEL_TEXT_EMBEDDING_ADA_002,
		EncodingFormat: GTP_ENCODING_FLOAT,
	}
	buf, err := json.MarshalIndent(reqObj, "", "\t")
	if err != nil {
		return nil, err
	}
	url := "https://api.openai.com/v1/embeddings"
	client := http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+secretProvider.GetSecret("openai"))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("%s\n", string(body))
	var respObj GptEmbeddingResponse
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return nil, err
	}
	if respObj.Error != nil {
		return nil, errors.New(respObj.Error.Message)
	}
	if len(respObj.Data) == 0 {
		return nil, errors.New("no embedding for you")
	}
	return NewEmbedding("", question.Text, "", GPT_CURRENT_MODEL).WithEmbedding(respObj.Data[0].Embedding), nil
}

func gptGetImage(secretProvider SecretProvider, question string) []*Answer {
	a := new(Answer)
	reqObj := GptImageRequest{
		GPT_MODEL_DALL_E_3,
		question,
		1,
		"1024x1024",
	}
	buf, err := json.MarshalIndent(reqObj, "", "\t")
	if err != nil {
		a.Text = fmt.Sprintf("gpt image error: %s", err.Error())
		return []*Answer{a}
	}
	url := "https://api.openai.com/v1/images/generations"
	client := http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(buf))
	if err != nil {
		a.Text = fmt.Sprintf("gpt image error: %s", err.Error())
		return []*Answer{a}
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+secretProvider.GetSecret("openai"))
	resp, err := client.Do(req)
	if err != nil {
		a.Text = fmt.Sprintf("gpt image error: %s", err.Error())
		return []*Answer{a}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.Text = fmt.Sprintf("gpt image error: %s", err.Error())
		return []*Answer{a}
	}
	fmt.Printf("%s\n", string(body))
	var respObj GptImageResponse
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		a.Text = fmt.Sprintf("gpt image error: %s", err.Error())
		return []*Answer{a}
	}
	if respObj.Error != nil && respObj.Error.Message != "" {
		a.Text = fmt.Sprintf("gpt image error: %s", respObj.Error.Message)
		return []*Answer{a}
	}
	if len(respObj.Data) == 0 {
		a.Text = "no image for you"
		return []*Answer{a}

	}
	a.Text = respObj.Data[0].URL
	return []*Answer{a}
}
