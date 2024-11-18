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
	"os"
	"path/filepath"
	"strings"
)

const (
	DEFAULT_KNOWLEDGE_BASE_NAME = "system"
	DEFAULT_EMBEDDING_BASE_NAME = "system"
	DEFAULT_KNOWLEDGE_BASE_PATH = "kb/facts"
	DEFAULT_EMBEDDING_BASE_PATH = "kb/embeddings"
)

type (
	KnowledeBaseManager struct {
		currentBaseName string
		secretProvider  SecretProvider
		oai             OpenAIHandler
		factsStores     map[string]KnowledeBaseProvider
		embeddingStores map[string]EmbeddingsBaseProvider
	}
)

func NewKnowledgeBaseManager(secretProvider SecretProvider, oai OpenAIHandler) (*KnowledeBaseManager, error) {
	kbm := KnowledeBaseManager{
		DEFAULT_EMBEDDING_BASE_NAME,
		secretProvider,
		oai,
		make(map[string]KnowledeBaseProvider, 0),
		make(map[string]EmbeddingsBaseProvider, 0),
	}
	err := kbm.loadAll()
	if err != nil {
		return nil, err
	}
	return &kbm, nil
}

func (kbm *KnowledeBaseManager) GetCurrentKnowledgeBase() KnowledeBaseProvider {
	return kbm.factsStores[kbm.currentBaseName]
}

func (kbm *KnowledeBaseManager) GetCurrentEmbeddingsBase() EmbeddingsBaseProvider {
	return kbm.embeddingStores[kbm.currentBaseName]
}

func (kbm *KnowledeBaseManager) GetCurrentBaseName() string {
	return kbm.currentBaseName
}

func (kbm *KnowledeBaseManager) SetCurrentBaseName(name string) error {
	_, ok := kbm.factsStores[name]
	if !ok {
		return errors.New("no knowledge base for " + name)
	}
	_, ok = kbm.embeddingStores[name]
	if !ok {
		return errors.New("no embeddings base for " + name)
	}
	kbm.currentBaseName = name
	return nil
}

func (kbm *KnowledeBaseManager) ListBaseNames() []string {
	names := make([]string, 0)
	for k, _ := range kbm.factsStores {
		names = append(names, k)
	}
	return names
}

func (kbm *KnowledeBaseManager) loadAll() error {
	files, err := os.ReadDir(DEFAULT_KNOWLEDGE_BASE_PATH)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			fkb := NewFileKnowledgeBase(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
			err = fkb.Load()
			if err != nil {
				return err
			}
			kbm.factsStores[fkb.GetName()] = fkb
		}
	}
	files, err = os.ReadDir(DEFAULT_EMBEDDING_BASE_PATH)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			feb := NewFileEmbeddingBase(kbm.secretProvider, kbm.oai, strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
			err = feb.Load()
			if err != nil {
				return err
			}
			kb, ok := kbm.factsStores[feb.GetName()]
			if !ok {
				return errors.New("no knowledge base for embedding base " + feb.GetName())
			}
			err = feb.SyncEmbeddings(kb)
			kbm.embeddingStores[feb.GetName()] = feb
		}
	}
	return nil
}
