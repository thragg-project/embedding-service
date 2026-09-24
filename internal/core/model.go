package core

const (
	EntityAggregateType = "entity"

	EntityEventCreated = "created"
)

type (
	Entity struct {
		Id     *int64
		Field1 string
		Field2 int
	}
	Model struct {
		Name           string
		Provider       string
		UpstreamModel  string
		Dimensions     int
		MaxInputTokens int
		Multilingual   bool
		Enabled        bool
	}
	InputType      string
	EmbeddingInput struct {
		ID   string
		Text string
	}
	Embedding struct {
		ID     string
		Vector []float32
	}
)

const (
	InputTypeQuery    InputType = "query"
	InputTypeDocument InputType = "document"
)
