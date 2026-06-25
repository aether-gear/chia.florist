package usecase

import (
	"context"
	"fmt"
	"strings"

	apperrors "service-core/internal/common/errors"
	shipping "service-core/internal/infra/shipping"
	"service-core/internal/modules/shipment/domain"
	transaction "service-core/internal/shared/transaction"
)

type EstimateShippingOptionsUsecase struct {
	shipping shipping.Provider
	executor transaction.Executor
}

func NewEstimateShippingOptionsUsecase(
	shipping shipping.Provider,
	executor transaction.Executor,
) *EstimateShippingOptionsUsecase {
	return &EstimateShippingOptionsUsecase{
		shipping: shipping,
		executor: executor,
	}
}

type EstimateShippingOptionsInput struct {
	Origin      int
	Destination int
	Weight      int
	Couriers    []string
	PriceFilter *string
}

func (u *EstimateShippingOptionsUsecase) Execute(
	ctx context.Context,
	input EstimateShippingOptionsInput,
) ([]shipping.RateOption, error) {
	if input.Origin == input.Destination {
		return nil, apperrors.NewInvalidInput(domain.ErrInvalidRoute.Error())
	}

	if len(input.Couriers) == 0 {
		return nil, apperrors.NewInvalidInput(domain.ErrNoCourierSelected.Error())
	}

	seen := make(map[string]struct{})
	result := make([]string, 0, len(input.Couriers))
	for _, courier := range input.Couriers {
		c := strings.TrimSpace(courier)
		if c == "" {
			return nil, apperrors.NewInvalidInput(domain.ErrInvalidCourier.Error())
		}

		if _, exist := seen[c]; exist {
			continue
		}

		seen[c] = struct{}{}
		result = append(result, c)
	}

	if input.Weight <= 0 {
		return nil, apperrors.NewInvalidInput(domain.ErrInvalidWeight.Error())
	}

	query := shipping.CalculateRatesInput{
		OriginID:      input.Origin,
		DestinationID: input.Destination,
		Weight:        input.Weight,
		Couriers:      result,
		PriceFilter:   input.PriceFilter,
	}

	costOptions, err := u.shipping.CalculateRates(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate shipping cost: %w", err)
	}

	return costOptions, nil
}
