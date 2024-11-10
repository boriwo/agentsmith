package main

import "testing"

func TestGptCompletions(t *testing.T) {
	secretProvider, err := NewJSONSecretProvider("secrets.json")
	if err != nil {
		t.Fatalf("failed to create secret provider: %v", err)
	}
	question := &Question{Text: "Please say this is a simple test!"}
	t.Logf("question: %s\n", question.Text)
	answers := gptGetCompletions(secretProvider, question)
	for _, a := range answers {
		t.Logf("answer: %s\n", a.Text)
	}
}

func TestGptEmbeddings(t *testing.T) {
	secretProvider, err := NewJSONSecretProvider("secrets.json")
	if err != nil {
		t.Fatalf("failed to create secret provider: %v", err)
	}
	question := &Question{Text: "cats and dogs"}
	t.Logf("question: %s\n", question.Text)
	embedding, err := gptGetEmbedding(secretProvider, question)
	if err != nil {
		t.Fatalf("failed to get embedding: %v\n", err)
	} else {
		t.Logf("embedding: %v\n", embedding)
	}
}

func TestGptImage(t *testing.T) {
	secretProvider, err := NewJSONSecretProvider("secrets.json")
	if err != nil {
		t.Fatalf("failed to create secret provider: %v", err)
	}
	question := &Question{Text: "cats and dogs"}
	t.Logf("question: %s\n", question.Text)
	answers := gptGetImage(secretProvider, question.Text)
	for _, a := range answers {
		t.Logf("answer: %s\n", a.Text)
	}
}
