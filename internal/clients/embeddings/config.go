package embeddings

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const ConfigEnvName = "EMBEDDINGS_CONFIG_BASE64"

type (
	Config struct {
		Providers []ProviderConfig `yaml:"providers"`
		Models    []ModelConfig    `yaml:"models"`
	}
	ProviderConfig struct {
		Name    string        `yaml:"name"`
		Type    ProviderType  `yaml:"type"`
		BaseURL string        `yaml:"base_url"`
		APIKey  string        `yaml:"api_key"`
		Timeout time.Duration `yaml:"-"`

		TimeoutRaw string `yaml:"timeout,omitempty"`
	}
	ProviderType string
	ModelConfig  struct {
		Name           string `yaml:"name"`
		Provider       string `yaml:"provider"`
		UpstreamModel  string `yaml:"upstream_model"`
		Dimensions     int    `yaml:"dimensions"`
		MaxInputTokens int    `yaml:"max_input_tokens"`
		Multilingual   bool   `yaml:"multilingual"`
		Enabled        bool   `yaml:"enabled"`
	}
)

const (
	ProviderTypeOpenAICompatible ProviderType = "openai_compatible"

	defaultProviderTimeout = 30 * time.Second
)

func LoadConfigFromEnv() (Config, error) {
	encoded := strings.TrimSpace(os.Getenv(ConfigEnvName))
	if encoded == "" {
		return Config{}, fmt.Errorf("%s is required", ConfigEnvName)
	}

	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return Config{}, fmt.Errorf("decode %s: %w", ConfigEnvName, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse %s yaml: %w", ConfigEnvName, err)
	}

	if err := cfg.validateAndNormalize(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) validateAndNormalize() error {
	if len(c.Providers) == 0 {
		return errors.New("embeddings config requires at least one provider")
	}
	if len(c.Models) == 0 {
		return errors.New("embeddings config requires at least one model")
	}

	providers := make(map[string]struct{}, len(c.Providers))
	for i := range c.Providers {
		provider := &c.Providers[i]
		if provider.Name == "" {
			return errors.New("embeddings provider name is required")
		}
		if _, exists := providers[provider.Name]; exists {
			return fmt.Errorf("duplicate embeddings provider %q", provider.Name)
		}
		providers[provider.Name] = struct{}{}

		switch provider.Type {
		case ProviderTypeOpenAICompatible:
		default:
			return fmt.Errorf("unsupported embeddings provider type %q", provider.Type)
		}
		if provider.BaseURL == "" {
			return fmt.Errorf("embeddings provider %q base_url is required", provider.Name)
		}
		if provider.TimeoutRaw == "" {
			provider.Timeout = defaultProviderTimeout
			continue
		}
		timeout, err := time.ParseDuration(provider.TimeoutRaw)
		if err != nil {
			return fmt.Errorf("parse embeddings provider %q timeout: %w", provider.Name, err)
		}
		if timeout <= 0 {
			return fmt.Errorf("embeddings provider %q timeout must be positive", provider.Name)
		}
		provider.Timeout = timeout
	}

	models := make(map[string]struct{}, len(c.Models))
	for i := range c.Models {
		model := &c.Models[i]
		if model.Name == "" {
			return errors.New("embeddings model name is required")
		}
		if _, exists := models[model.Name]; exists {
			return fmt.Errorf("duplicate embeddings model %q", model.Name)
		}
		models[model.Name] = struct{}{}

		if model.Provider == "" {
			return fmt.Errorf("embeddings model %q provider is required", model.Name)
		}
		if _, exists := providers[model.Provider]; !exists {
			return fmt.Errorf(
				"embeddings model %q references unknown provider %q",
				model.Name,
				model.Provider,
			)
		}
		if model.UpstreamModel == "" {
			model.UpstreamModel = model.Name
		}
		if model.Dimensions <= 0 {
			return fmt.Errorf("embeddings model %q dimensions must be positive", model.Name)
		}
		if model.MaxInputTokens <= 0 {
			return fmt.Errorf("embeddings model %q max_input_tokens must be positive", model.Name)
		}
	}

	return nil
}
