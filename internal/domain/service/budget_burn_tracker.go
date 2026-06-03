package service

import "github.com/promotionos/analytics-service/internal/domain/model"

type BudgetBurnTracker interface {
	Update(metrics *model.CampaignMetrics, discountApplied float64) (exhausted bool, err error)
	IsExhausted(metrics *model.CampaignMetrics) bool
}
