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
	"errors"
	"strings"
)

type EmbeddingAnswerProvider struct {
	kbm *KnowledeBaseManager
	oai OpenAIHandler
	pm  *PluginManager
}

func NewEmbeddingAnswerProvider(kbm *KnowledeBaseManager, oai OpenAIHandler) AnswerProvider {
	answerProvider := EmbeddingAnswerProvider{
		kbm,
		oai,
		NewPluginManger(kbm, oai),
	}
	return &answerProvider
}

func (sap *EmbeddingAnswerProvider) GetAnswers(session *UserSession, question *Question) ([]*Answer, error) {
	answers := make([]*Answer, 0)
	embedding, err := sap.oai.GptGetEmbedding(question)
	if err != nil {
		return nil, err
	}
	ranking, err := sap.kbm.GetCurrentEmbeddingsBase().RankEmbeddings(embedding)
	if err != nil {
		return nil, err
	}
	if len(ranking.Embeddings) == 0 {
		return nil, errors.New("no matching fact")
	}
	fact := sap.kbm.GetCurrentKnowledgeBase().GetFact(ranking.Embeddings[0].FactName)
	if fact == nil {
		return nil, errors.New("no matching fact")
	}
	if fact.Plugin != "" {
		answers, err = sap.pm.GetAnswers(session, question, fact)
		if err != nil {
			return nil, err
		}
	} else {
		plausabilityPrompt := ""
		plausabilityPrompt += "Please check if the following answer is a plausible answer to the give question. Answer simply with yes or no.\n"
		plausabilityPrompt += "Question:\n" + question.Text + "\n"
		plausabilityPrompt += "Answer:\n"
		for _, a := range fact.Answers {
			answers = append(answers, &Answer{a, "", "", 0, 0})
			plausabilityPrompt += a + "\n"
		}
		plausabilityAnswers, err := sap.oai.GptGetCompletions(&Question{plausabilityPrompt})
		if err != nil {
			return nil, err
		}
		// answer deemed not plausible
		if len(plausabilityAnswers) == 1 && strings.ToLower(plausabilityAnswers[0].Text) == "no" {
			answers = make([]*Answer, 0)
		}
	}
	session.LastQuestion = question
	session.LastAnswer = answers
	return answers, nil
}
