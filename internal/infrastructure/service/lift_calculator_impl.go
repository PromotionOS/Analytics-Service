package infraservice

import "github.com/promotionos/analytics-service/internal/domain/model"

type LiftCalculatorImpl struct{}

func NewLiftCalculator() *LiftCalculatorImpl {
	return &LiftCalculatorImpl{}
}

func (l *LiftCalculatorImpl) Calculate(metrics *model.CampaignMetrics) error {
	// Calculate lift
	lift := metrics.ActualSalesPerDay - metrics.BaselineSalesPerDay
	metrics.Lift = &lift

	// Calculate lift percentage
	if metrics.BaselineSalesPerDay > 0 {
		liftPct := (lift / metrics.BaselineSalesPerDay) * 100
		metrics.LiftPercentage = &liftPct
	}
	// If baseline is 0 — LiftPercentage stays nil (correct)

	// Calculate incremental margin (30% average margin)
	if metrics.Lift != nil {
		margin := *metrics.Lift * 0.30
		metrics.IncrementalMargin = &margin
	}

	// BUG PLANTED: missing null guard on TotalFundingCost
	// Causes divide-by-zero panic on 100% Kroger-funded campaigns
	// where vendorShare = 0 and TotalFundingCost = 0
	// Teams find this in Sprint 1 via RCA — Scenario 20
	if metrics.IncrementalMargin != nil {
		roi := *metrics.IncrementalMargin / metrics.TotalFundingCost // PANIC when TotalFundingCost = 0
		metrics.ROI = &roi
	}

	return nil
}
