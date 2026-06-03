package application

import (
	"github.com/google/uuid"
	"github.com/promotionos/analytics-service/internal/domain/model"
	"github.com/promotionos/analytics-service/internal/domain/repository"
	"github.com/promotionos/analytics-service/internal/domain/service"
)

type AnalyticsService struct {
	repo        repository.AnalyticsRepository
	liftCalc    service.LiftCalculator
	burnTracker service.BudgetBurnTracker
	onExhausted func(tenantID, campaignID string, metrics *model.CampaignMetrics) error
}

func NewAnalyticsService(
	repo repository.AnalyticsRepository,
	liftCalc service.LiftCalculator,
	burnTracker service.BudgetBurnTracker,
	onExhausted func(tenantID, campaignID string, metrics *model.CampaignMetrics) error,
) *AnalyticsService {
	return &AnalyticsService{
		repo:        repo,
		liftCalc:    liftCalc,
		burnTracker: burnTracker,
		onExhausted: onExhausted,
	}
}

func (s *AnalyticsService) HandleOfferRedeemed(campaignID, tenantID string, discountApplied float64) error {
	metrics, err := s.repo.FindByCampaign(campaignID, tenantID)
	if err != nil || metrics == nil {
		return nil
	}

	exhausted, err := s.burnTracker.Update(metrics, discountApplied)
	if err != nil {
		return err
	}

	// BUG: LiftCalculatorImpl.Calculate() will panic on zero-cost campaigns
	if err := s.liftCalc.Calculate(metrics); err != nil {
		return err
	}

	if err := s.repo.Save(metrics); err != nil {
		return err
	}

	if exhausted && s.onExhausted != nil {
		return s.onExhausted(tenantID, campaignID, metrics)
	}
	return nil
}

func (s *AnalyticsService) HandleCampaignPublished(campaignID, tenantID string) error {
	existing, _ := s.repo.FindByCampaign(campaignID, tenantID)
	if existing != nil {
		return nil // idempotent
	}
	metrics := &model.CampaignMetrics{
		ID:         uuid.New().String(),
		CampaignID: campaignID,
		TenantID:   tenantID,
	}
	return s.repo.Save(metrics)
}

func (s *AnalyticsService) GetReport(campaignID, tenantID string) (*model.CampaignMetrics, error) {
	return s.repo.FindByCampaign(campaignID, tenantID)
}

func (s *AnalyticsService) GetBurn(campaignID, tenantID string) (*model.CampaignMetrics, error) {
	// TODO Team 5 Sprint 2: implement real-time burn endpoint
	return s.repo.FindByCampaign(campaignID, tenantID)
}
