package main

type UberAnswerProvider struct {
	kb          KnowledeBaseProvider
	eb          EmbeddingsBaseProvider
	oai         OpenAIHandler
	answerChain []AnswerProvider
}

func NewUberAnswerProvider(kb KnowledeBaseProvider, eb EmbeddingsBaseProvider, oai OpenAIHandler) AnswerProvider {
	answerProvider := UberAnswerProvider{
		kb,
		eb,
		oai,
		[]AnswerProvider{},
	}
	answerProvider.answerChain = append(answerProvider.answerChain, NewCommandAnswerProvider(kb, eb))
	answerProvider.answerChain = append(answerProvider.answerChain, NewEmbeddingAnswerProvider(kb, eb, oai))
	answerProvider.answerChain = append(answerProvider.answerChain, NewSimpleAnswerProvider())
	return &answerProvider
}

func (sap *UberAnswerProvider) GetAnswers(session *UserSession, question *Question) ([]*Answer, error) {
	session.LastQuestion = question
	for _, ap := range sap.answerChain {
		answers, err := ap.GetAnswers(session, question)
		if err != nil {
			return nil, err
		}
		if len(answers) > 0 {
			return answers, nil
		}
	}
	return nil, nil
}
