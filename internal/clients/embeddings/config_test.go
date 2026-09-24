package embeddings

import (
	"encoding/base64"
	"testing"
	"time"
)

func TestLoadConfigFromEnvParsesYAML(t *testing.T) {
	raw := []byte(`
providers:
  - name: openai
    type: openai_compatible
    base_url: https://api.openai.com/v1
    api_key: secret
    timeout: 15s
models:
  - name: text-embedding-3-small
    provider: openai
    upstream_model: text-embedding-3-small
    dimensions: 1536
    max_input_tokens: 8191
    multilingual: true
    enabled: true
`)
	t.Setenv(ConfigEnvName, base64.StdEncoding.EncodeToString(raw))

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if got := cfg.Providers[0].Timeout; got != 15*time.Second {
		t.Fatalf("timeout = %s, want 15s", got)
	}
	if got := cfg.Models[0].Name; got != "text-embedding-3-small" {
		t.Fatalf("model name = %q, want text-embedding-3-small", got)
	}
}
