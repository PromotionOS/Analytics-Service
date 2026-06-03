package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const ChannelBudgetExhausted = "promotionos.analytics.budget.exhausted"

type Publisher struct {
	client *redis.Client
}

func NewPublisher(client *redis.Client) *Publisher {
	return &Publisher{client: client}
}

type BudgetExhaustedEvent struct {
	EventID       string  `json:"eventId"`
	TenantID      string  `json:"tenantId"`
	CampaignID    string  `json:"campaignId"`
	OccurredAt    string  `json:"occurredAt"`
	SchemaVersion int     `json:"schemaVersion"`
	Payload       BudgetExhaustedPayload `json:"payload"`
}

type BudgetExhaustedPayload struct {
	TotalAmount       float64 `json:"totalAmount"`
	BurnedAmount      float64 `json:"burnedAmount"`
	BudgetBurnPercent float64 `json:"budgetBurnPercent"`
	RedemptionCount   int     `json:"redemptionCount"`
	ExhaustedAt       string  `json:"exhaustedAt"`
}

func (p *Publisher) PublishBudgetExhausted(tenantID, campaignID string,
	totalAmount, burnedAmount, burnPercent float64, redemptionCount int) error {
	event := BudgetExhaustedEvent{
		EventID:       uuid.New().String(),
		TenantID:      tenantID,
		CampaignID:    campaignID,
		OccurredAt:    time.Now().UTC().Format(time.RFC3339),
		SchemaVersion: 1,
		Payload: BudgetExhaustedPayload{
			TotalAmount:       totalAmount,
			BurnedAmount:      burnedAmount,
			BudgetBurnPercent: burnPercent,
			RedemptionCount:   redemptionCount,
			ExhaustedAt:       time.Now().UTC().Format(time.RFC3339),
		},
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.client.Publish(context.Background(), ChannelBudgetExhausted, payload).Err()
}
