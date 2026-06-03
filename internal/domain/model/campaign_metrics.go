package model

import "time"

type CampaignMetrics struct {
	ID                     string
	CampaignID             string
	TenantID               string
	BaselineSalesPerDay    float64
	ActualSalesPerDay      float64
	Lift                   *float64
	LiftPercentage         *float64
	TotalFundingCost       float64
	IncrementalMargin      *float64
	ROI                    *float64
	BudgetBurnPercent      float64
	BurnedAmount           float64
	TotalAmount            float64
	RedemptionCount        int
	BudgetExhaustedEmitted bool
	UpdatedAt              time.Time
}

func (m *CampaignMetrics) AddRedemption(discountApplied float64) {
	m.BurnedAmount += discountApplied
	m.RedemptionCount++
	if m.TotalAmount > 0 {
		m.BudgetBurnPercent = (m.BurnedAmount / m.TotalAmount) * 100
	}
}
