package repository

import "github.com/promotionos/analytics-service/internal/domain/model"

type AnalyticsRepository interface {
	FindByCampaign(campaignID string, tenantID string) (*model.CampaignMetrics, error)
	Save(metrics *model.CampaignMetrics) error
	FindAll(tenantID string) ([]*model.CampaignMetrics, error)
}
