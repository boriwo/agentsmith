package main

import (
	"errors"
	"strings"
)

type EmbeddingAnswerProvider struct {
	kb  KnowledeBaseProvider
	eb  EmbeddingsBaseProvider
	oai OpenAIHandler
	pm  *PluginManager
}

func NewEmbeddingAnswerProvider(kb KnowledeBaseProvider, eb EmbeddingsBaseProvider, oai OpenAIHandler) AnswerProvider {
	answerProvider := EmbeddingAnswerProvider{
		kb,
		eb,
		oai,
		NewPluginManger(kb, eb, oai),
	}
	return &answerProvider
}

func (sap *EmbeddingAnswerProvider) GetAnswers(session *UserSession, question *Question) ([]*Answer, error) {
	answers := make([]*Answer, 0)
	embedding, err := sap.oai.GptGetEmbedding(question)
	if err != nil {
		return nil, err
	}
	ranking, err := sap.eb.RankEmbeddings(embedding)
	if err != nil {
		return nil, err
	}
	if len(ranking.Embeddings) == 0 {
		return nil, errors.New("no matching fact")
	}
	fact := sap.kb.GetFact(ranking.Embeddings[0].FactName)
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
		// answer deemed no plausible
		if len(plausabilityAnswers) == 1 && strings.ToLower(plausabilityAnswers[0].Text) == "no" {
			answers = make([]*Answer, 0)
		}
	}
	session.LastQuestion = question
	session.LastAnswer = answers
	return answers, nil
}
