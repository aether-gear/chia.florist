package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	paymentgateway "service-core/internal/infra/payment-gateway"
	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
	query "service-core/internal/shared/query"
	transaction "service-core/internal/shared/transaction"
)

func SyncPaymentMethods(
	ctx context.Context,
	methodRepo repository.PaymentMethodRepository,
	executor transaction.Executor,
	gateway paymentgateway.Provider,
) error {
	allowedMethods := gateway.AllowedPaymentMethods()
	providerName := gateway.Name()

	existingMethods, err := methodRepo.ListAll(ctx, executor, query.Sorts{})
	if err != nil {
		return fmt.Errorf("failed to list existing payment methods: %w", err)
	}

	type key struct {
		code     string
		provider string
		mType    string
	}
	existingMap := make(map[key]*domain.PaymentMethod)
	for i := range existingMethods {
		m := &existingMethods[i]
		existingMap[key{
			code:     string(m.Code),
			provider: m.Provider,
			mType:    string(m.Type),
		}] = m
	}

	activeKeys := make(map[key]bool)

	for _, am := range allowedMethods {
		k := key{
			code:     am.Code,
			provider: providerName,
			mType:    am.Type,
		}
		activeKeys[k] = true

		existing, ok := existingMap[k]
		if !ok {
			newMethod := domain.PaymentMethod{
				ID:            uuid.New(),
				Name:          am.Name,
				Code:          domain.PaymentMethodCode(am.Code),
				Provider:      providerName,
				Type:          domain.PaymentMethodType(am.Type),
				IsActive:      true,
				Description:   am.Description,
				FeeType:       domain.PaymentFeeType(am.FeeType),
				FeeFixed:      am.FeeFixed,
				FeePercentage: am.FeePercentage,
				CreatedAt:     time.Now(),
			}
			if err := methodRepo.Save(ctx, executor, newMethod); err != nil {
				return fmt.Errorf("failed to save new payment method %s: %w", am.Code, err)
			}
		} else {
			updatedMethod := *existing
			updatedMethod.Name = am.Name
			updatedMethod.Description = am.Description
			updatedMethod.FeeType = domain.PaymentFeeType(am.FeeType)
			updatedMethod.FeeFixed = am.FeeFixed
			updatedMethod.FeePercentage = am.FeePercentage

			if err := methodRepo.Save(ctx, executor, updatedMethod); err != nil {
				return fmt.Errorf("failed to update payment method %s: %w", am.Code, err)
			}
		}
	}

	for k, m := range existingMap {
		if !activeKeys[k] && m.IsActive {
			m.IsActive = false
			if err := methodRepo.Save(ctx, executor, *m); err != nil {
				return fmt.Errorf("failed to deactivate obsolete payment method %s: %w", m.Code, err)
			}
		}
	}

	return nil
}
