package event

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/promotionos/analytics-service/internal/application"
)

const (
	ChannelOfferRedeemed      = "promotionos.redemption.redeemed"
	ChannelCampaignPublished  = "promotionos.campaign.published"
	ChannelBudgetUpdated      = "promotionos.campaign.budget.updated"
)

type OfferRedeemedEvent struct {
	EventID    string `json:"eventId"`
	TenantID   string `json:"tenantId"`
	CampaignID string `json:"campaignId"`
	Payload    struct {
		RedemptionID    string  `json:"redemptionId"`
		CustomerID      string  `json:"customerId"`
		DiscountApplied float64 `json:"discountApplied"`
		CartTotal       float64 `json:"cartTotal"`
		RedeemedAt      string  `json:"redeemedAt"`
	} `json:"payload"`
}

type CampaignPublishedEvent struct {
	EventID    string `json:"eventId"`
	TenantID   string `json:"tenantId"`
	CampaignID string `json:"campaignId"`
}

func StartConsumers(client *redis.Client, svc *application.AnalyticsService) {
	go consumeOfferRedeemed(client, svc)
	go consumeCampaignPublished(client, svc)
	log.Println("Analytics: event consumers started")
}

func consumeOfferRedeemed(client *redis.Client, svc *application.AnalyticsService) {
	pubsub := client.Subscribe(context.Background(), ChannelOfferRedeemed)
	for msg := range pubsub.Channel() {
		var event OfferRedeemedEvent
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			log.Printf("Analytics: failed to parse OfferRedeemed: %v", err)
			continue
		}
		if err := svc.HandleOfferRedeemed(
			event.CampaignID,
			event.TenantID,
			event.Payload.DiscountApplied,
		); err != nil {
			log.Printf("Analytics: HandleOfferRedeemed error: %v", err)
		}
	}
}

func consumeCampaignPublished(client *redis.Client, svc *application.AnalyticsService) {
	pubsub := client.Subscribe(context.Background(), ChannelCampaignPublished)
	for msg := range pubsub.Channel() {
		var event CampaignPublishedEvent
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			log.Printf("Analytics: failed to parse CampaignPublished: %v", err)
			continue
		}
		if err := svc.HandleCampaignPublished(event.CampaignID, event.TenantID); err != nil {
			log.Printf("Analytics: HandleCampaignPublished error: %v", err)
		}
	}
}
