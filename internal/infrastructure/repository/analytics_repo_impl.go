package repository

import (
	"errors"

	"gorm.io/gorm"
	"github.com/promotionos/analytics-service/internal/domain/model"
)

type AnalyticsRepositoryImpl struct {
	db *gorm.DB
}

func NewAnalyticsRepositoryImpl(db *gorm.DB) *AnalyticsRepositoryImpl {
	return &AnalyticsRepositoryImpl{db: db}
}

type CampaignMetricsDB struct {
	ID                     string  `gorm:"column:id;primaryKey"`
	CampaignID             string  `gorm:"column:campaign_id"`
	TenantID               string  `gorm:"column:tenant_id"`
	BaselineSalesPerDay    float64 `gorm:"column:baseline_sales_per_day"`
	ActualSalesPerDay      float64 `gorm:"column:actual_sales_per_day"`
	TotalFundingCost       float64 `gorm:"column:total_funding_cost"`
	BudgetBurnPercent      float64 `gorm:"column:budget_burn_percent"`
	BurnedAmount           float64 `gorm:"column:burned_amount"`
	TotalAmount            float64 `gorm:"column:total_amount"`
	RedemptionCount        int     `gorm:"column:redemption_count"`
	BudgetExhaustedEmitted bool    `gorm:"column:budget_exhausted_emitted"`
}

func (CampaignMetricsDB) TableName() string { return "campaign_metrics" }

func (r *AnalyticsRepositoryImpl) FindByCampaign(campaignID, tenantID string) (*model.CampaignMetrics, error) {
	var record CampaignMetricsDB
	result := r.db.Where("campaign_id = ? AND tenant_id = ?", campaignID, tenantID).First(&record)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &model.CampaignMetrics{
		ID:                     record.ID,
		CampaignID:             record.CampaignID,
		TenantID:               record.TenantID,
		BaselineSalesPerDay:    record.BaselineSalesPerDay,
		ActualSalesPerDay:      record.ActualSalesPerDay,
		TotalFundingCost:       record.TotalFundingCost,
		BudgetBurnPercent:      record.BudgetBurnPercent,
		BurnedAmount:           record.BurnedAmount,
		TotalAmount:            record.TotalAmount,
		RedemptionCount:        record.RedemptionCount,
		BudgetExhaustedEmitted: record.BudgetExhaustedEmitted,
	}, nil
}

func (r *AnalyticsRepositoryImpl) Save(metrics *model.CampaignMetrics) error {
	record := CampaignMetricsDB{
		ID:                     metrics.ID,
		CampaignID:             metrics.CampaignID,
		TenantID:               metrics.TenantID,
		BaselineSalesPerDay:    metrics.BaselineSalesPerDay,
		ActualSalesPerDay:      metrics.ActualSalesPerDay,
		TotalFundingCost:       metrics.TotalFundingCost,
		BudgetBurnPercent:      metrics.BudgetBurnPercent,
		BurnedAmount:           metrics.BurnedAmount,
		TotalAmount:            metrics.TotalAmount,
		RedemptionCount:        metrics.RedemptionCount,
		BudgetExhaustedEmitted: metrics.BudgetExhaustedEmitted,
	}
	return r.db.Save(&record).Error
}

func (r *AnalyticsRepositoryImpl) FindAll(tenantID string) ([]*model.CampaignMetrics, error) {
	var records []CampaignMetricsDB
	if err := r.db.Where("tenant_id = ?", tenantID).Find(&records).Error; err != nil {
		return nil, err
	}
	var metrics []*model.CampaignMetrics
	for _, rec := range records {
		metrics = append(metrics, &model.CampaignMetrics{
			ID:                     rec.ID,
			CampaignID:             rec.CampaignID,
			TenantID:               rec.TenantID,
			BudgetBurnPercent:      rec.BudgetBurnPercent,
			BurnedAmount:           rec.BurnedAmount,
			TotalAmount:            rec.TotalAmount,
			RedemptionCount:        rec.RedemptionCount,
			BudgetExhaustedEmitted: rec.BudgetExhaustedEmitted,
		})
	}
	return metrics, nil
}
