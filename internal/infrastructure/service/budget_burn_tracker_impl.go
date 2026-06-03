package infraservice

import "github.com/promotionos/analytics-service/internal/domain/model"

type BudgetBurnTrackerImpl struct{}

func NewBudgetBurnTracker() *BudgetBurnTrackerImpl {
	return &BudgetBurnTrackerImpl{}
}

func (b *BudgetBurnTrackerImpl) Update(metrics *model.CampaignMetrics,
	discountApplied float64) (bool, error) {
	metrics.AddRedemption(discountApplied)

	exhausted := b.IsExhausted(metrics)
	if exhausted && !metrics.BudgetExhaustedEmitted {
		metrics.BudgetExhaustedEmitted = true
		return true, nil
	}
	return false, nil
}

func (b *BudgetBurnTrackerImpl) IsExhausted(metrics *model.CampaignMetrics) bool {
	return metrics.BudgetBurnPercent >= 95.0
}
