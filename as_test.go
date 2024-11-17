/**
 * Copyright 2024 Boris Wolf
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	answers, err := oai.GptGetCompletions(question)
	if err != nil {
		t.Errorf("failed get completion: %v", err)
		return
	}
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

/*func TestGptImage(t *testing.T) {
	secretProvider, err := NewJSONSecretProvider("secrets.json")
	if err != nil {
		t.Errorf("failed to create secret provider: %v", err)
		return
	}
	oai := NewOpenAIHandler(secretProvider)
	question := &Question{Text: "cats and dogs"}
	t.Logf("question: %s\n", question.Text)
	answers, err := oai.GptGetImage(question)
	if err != nil {
		t.Errorf("failed get image: %v", err)
		return
	}
	for _, a := range answers {
		t.Logf("answer: %s\n", a.Text)
	}
}*/
