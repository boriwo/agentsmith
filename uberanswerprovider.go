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

type UberAnswerProvider struct {
	kbm         *KnowledeBaseManager
	oai         OpenAIHandler
	answerChain []AnswerProvider
}

func NewUberAnswerProvider(kbm *KnowledeBaseManager, oai OpenAIHandler) AnswerProvider {
	answerProvider := UberAnswerProvider{
		kbm,
		oai,
		[]AnswerProvider{},
	}
	answerProvider.answerChain = append(answerProvider.answerChain, NewCommandAnswerProvider(kbm))
	answerProvider.answerChain = append(answerProvider.answerChain, NewEmbeddingAnswerProvider(kbm, oai))
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
