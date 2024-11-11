package main

import (
	"strings"
	"testing"
)

func TestGptCompletions(t *testing.T) {
	secretProvider, err := NewJSONSecretProvider("secrets.json")
	if err != nil {
		t.Errorf("failed to create secret provider: %v", err)
		return
	}
	oai := NewOpenAIHandler(secretProvider)
	question := &Question{Text: "Please say this is a simple test!"}
	t.Logf("question: %s\n", question.Text)
	answers := oai.GptGetCompletions(question)
	for _, a := range answers {
		t.Logf("answer: %s\n", a.Text)
		if !strings.Contains(a.Text, "test") {
			t.Errorf("unexpected answer, expected 'test': %s", a.Text)
		}
	}
}

func TestGptEmbeddings(t *testing.T) {
	secretProvider, err := NewJSONSecretProvider("secrets.json")
	if err != nil {
		t.Errorf("failed to create secret provider: %v", err)
		return
	}
	oai := NewOpenAIHandler(secretProvider)
	question := &Question{Text: "cats and dogs"}
	t.Logf("question: %s\n", question.Text)
	embedding, err := oai.GptGetEmbedding(question)
	if err != nil {
		t.Errorf("failed to get embedding: %v\n", err)
		return
	} else {
		t.Logf("embedding dimensions: %d\n", len(embedding.Embedding))
	}
	if len(embedding.Embedding) != 1536 {
		t.Errorf("wrong number of dimensions in embedding: %d, expected 1536\n", len(embedding.Embedding))
		return
	}
}

func TestGptImage(t *testing.T) {
	secretProvider, err := NewJSONSecretProvider("secrets.json")
	if err != nil {
		t.Errorf("failed to create secret provider: %v", err)
		return
	}
	oai := NewOpenAIHandler(secretProvider)
	question := &Question{Text: "cats and dogs"}
	t.Logf("question: %s\n", question.Text)
	answers := oai.GptGetImage(question)
	for _, a := range answers {
		t.Logf("answer: %s\n", a.Text)
	}
}
