package service

import "github.com/promotionos/analytics-service/internal/domain/model"

type LiftCalculator interface {
	Calculate(metrics *model.CampaignMetrics) error
}
