# Analytics Service

**Bounded context:** Analytics  
**Language:** Go 1.21 / Gin  
**Status:** Pre-built with bug planted in `LiftCalculatorImpl.Calculate()`

## Bug Planted (Sprint 1 target)
- **Location:** `internal/infrastructure/service/lift_calculator_impl.go` line 32
- **Bug:** Divides by `TotalFundingCost` without zero guard — panics on 100% Kroger-funded campaigns
- **Symptom:** `GET /analytics/campaigns/camp-003/report` returns 500
- **Validation:** Scenarios 18, 19, 20 from test-data

## What Teams Build Per Sprint
- **Sprint 1:** RCA + fix the divide-by-zero bug
- **Sprint 2:** Real-time burn tracking via OfferRedeemed events, BudgetExhausted event
- **Sprint 3:** Full lift/ROI calculation, BudgetUpdated event handling
- **Sprint 4:** All scenarios passing, regression verified

## Local Development
```bash
export DB_URL=postgresql://localhost:5432/railway
export REDIS_URL=redis://localhost:6379
export TENANT_ID=tenant-kroger-001
go run cmd/main.go
```
