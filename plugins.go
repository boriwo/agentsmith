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

import "errors"

const (
	COMMAND_PLUGIN      = "COMMAND_PLUGIN"
	IMAGE_PLUGIN        = "IMAGE_PLUGIN"
	PARAM_TYPE_CONSTANT = "constant"
)

type PluginManager struct {
	kb      KnowledeBaseProvider
	eb      EmbeddingsBaseProvider
	oai     OpenAIHandler
	plugins map[string]AnswerProvider
}

func NewPluginManger(kb KnowledeBaseProvider, eb EmbeddingsBaseProvider, oai OpenAIHandler) *PluginManager {
	mgr := &PluginManager{
		kb,
		eb,
		oai,
		make(map[string]AnswerProvider),
	}
	mgr.plugins[COMMAND_PLUGIN] = NewCommandAnswerProvider(kb, eb)
	mgr.plugins[IMAGE_PLUGIN] = NewImageAnswerProvider(oai)
	return mgr
}

func (pm *PluginManager) GetAnswers(session *UserSession, question *Question, fact *Fact) ([]*Answer, error) {
	if question == nil || question.Text == "" {
		return nil, errors.New("missing question")
	}
	if fact == nil || fact.Plugin == "" {
		return nil, errors.New("missing fact or blank plugin")
	}
	answerProvider, ok := pm.plugins[fact.Plugin]
	if !ok {
		return nil, errors.New("unknown plugin " + fact.Plugin)
	}
	// rewrite question
	q := question.Text
	if len(fact.Params) > 0 {
		q = ""
		for idx, param := range fact.Params {
			if param.Type == PARAM_TYPE_CONSTANT {
				q += param.Value
				if idx < len(fact.Params)-1 {
					q += " "
				}
			} else {
				return nil, errors.New("unsupported param type " + param.Type)
			}
		}
	}
	return answerProvider.GetAnswers(session, &Question{q})
}
