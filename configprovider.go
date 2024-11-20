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

type ConfigProvider interface {
	GetConfig(name string) string
}

type JSONConfigProvider struct {
	filePath string
	configs  map[string]string
	mu       sync.RWMutex
}

func NewJSONConfigProvider(filePath string) (ConfigProvider, error) {
	provider := &JSONConfigProvider{
		filePath: filePath,
		configs:  make(map[string]string),
	}
	err := provider.loadSecrets()
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (cp *JSONConfigProvider) loadSecrets() error {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	data, err := os.ReadFile(cp.filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	err = json.Unmarshal(data, &cp.configs)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err)
	}
	return nil
}

func (cp *JSONConfigProvider) GetConfig(name string) string {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	return cp.configs[name]
}
