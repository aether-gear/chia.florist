package persistence

import (
	"context"
	"fmt"

	"service-core/internal/modules/merchant/domain"
	"service-core/internal/modules/merchant/repository"
	transaction "service-core/internal/shared/transaction"
)

type merchantRepositoryImpl struct{}

func NewMerchantRepositoryImpl() repository.MerchantRepository {
	return &merchantRepositoryImpl{}
}

func (r *merchantRepositoryImpl) Create(
	ctx context.Context,
	exec transaction.Executor,
	merchant domain.Merchant,
) error {
	query := `
		INSERT INTO merchants (
			id,
			name,
			description,
			logo_url,
			banner_url,
			created_at
		) VALUES ($1,$2,$3,$4,$5,$6)
	`

	_, err := exec.Exec(ctx, query,
		merchant.ID,
		merchant.Name,
		merchant.Description,
		merchant.LogoUrl,
		merchant.BannerUrl,
		merchant.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("insert merchant failed: %w", err)
	}
	return nil
}
