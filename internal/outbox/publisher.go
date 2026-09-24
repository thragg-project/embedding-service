package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fr33dman/gobox"

	"github.com/fr33dman/go-template/internal/services"
	"github.com/fr33dman/go-template/pkg/database"
)

type Publisher struct{}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Publish(ctx context.Context, event services.Event) error {
	tx, ok := database.Tx(ctx)
	if !ok {
		return fmt.Errorf("publish outbox event: transaction is required")
	}

	rawPayload, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}

	outbox := gobox.NewOutbox(tx)
	_, err = outbox.Produce(ctx, gobox.Event{
		AggregateType: event.AggregateType,
		AggregateId:   strconv.FormatInt(event.AggregateID, 10),
		EventType:     event.EventType,
		Payload:       string(rawPayload),
		Headers:       "",
		DedupKey: fmt.Sprintf(
			"%s.%s.%d",
			event.AggregateType,
			event.EventType,
			event.AggregateID,
		),
	})

	return err
}
