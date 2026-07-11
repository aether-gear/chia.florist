package usecase

import (
	"context"

	"service-core/internal/modules/threat_intel/domain"
	"service-core/internal/modules/threat_intel/repository"
)

type GetGeoIPUsecase struct {
	threatIntelProv repository.ThreatIntelProvider
}

func NewGetGeoIPUsecase(
	threatIntelProv repository.ThreatIntelProvider,
) *GetGeoIPUsecase {
	return &GetGeoIPUsecase{
		threatIntelProv: threatIntelProv,
	}
}

func (u *GetGeoIPUsecase) Execute(
	ctx context.Context,
	ip string,
) (*domain.GeoReport, error) {
	if ip == "" {
		return nil, domain.ErrInvalidIP
	}

	report, err := u.threatIntelProv.GetGeolocation(ctx, ip)
	if err != nil {
		return nil, err
	}

	return &domain.GeoReport{
		IP:        ip,
		RawReport: report,
	}, nil
}
