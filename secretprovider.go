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
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type SecretProvider interface {
	GetSecret(name string) string
}

type JSONSecretProvider struct {
	filePath string
	secrets  map[string]string
	mu       sync.RWMutex
}

func NewJSONSecretProvider(filePath string) (*JSONSecretProvider, error) {
	provider := &JSONSecretProvider{
		filePath: filePath,
		secrets:  make(map[string]string),
	}
	err := provider.loadSecrets()
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (sp *JSONSecretProvider) loadSecrets() error {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	data, err := os.ReadFile(sp.filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	err = json.Unmarshal(data, &sp.secrets)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err)
	}
	return nil
}

func (sp *JSONSecretProvider) GetSecret(name string) string {
	sp.mu.RLock()
	defer sp.mu.RUnlock()
	return sp.secrets[name]
}
