package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	data, err := ioutil.ReadFile(sp.filePath)
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
