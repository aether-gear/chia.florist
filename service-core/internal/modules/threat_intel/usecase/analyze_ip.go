package usecase

import (
	"context"

	"service-core/internal/modules/threat_intel/domain"
	"service-core/internal/modules/threat_intel/repository"
)

type AnalyzeIPUsecase struct {
	threatIntelProv repository.ThreatIntelProvider
}

func NewAnalyzeIPUsecase(
	threatIntelProv repository.ThreatIntelProvider,
) *AnalyzeIPUsecase {
	return &AnalyzeIPUsecase{
		threatIntelProv: threatIntelProv,
	}
}

func (u *AnalyzeIPUsecase) Execute(
	ctx context.Context,
	ip string,
	apiKey string,
) (*domain.ReputationReport, error) {
	if ip == "" {
		return nil, domain.ErrInvalidIP
	}

	report, err := u.threatIntelProv.GetReputation(ctx, ip, apiKey)
	if err != nil {
		return nil, err
	}

	return &domain.ReputationReport{
		IP:        ip,
		RawReport: report,
	}, nil
}
