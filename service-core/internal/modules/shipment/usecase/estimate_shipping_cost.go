package usecase

import (
	"fmt"
	"strings"

	appErr "service-core/internal/common/errors"

	"service-core/internal/modules/shipment/domain"
	"service-core/internal/modules/shipment/repository"
)

type EstimateShippingCostUsecase struct {
	shippingCostProvider repository.ShippingCostProvider
}

func NewEstimateShippingCostUsecase(
	shippingCostProvider repository.ShippingCostProvider,
) *EstimateShippingCostUsecase {
	return &EstimateShippingCostUsecase{
		shippingCostProvider: shippingCostProvider,
	}
}

type EstimateShippingCostInput struct {
	Origin      int
	Destination int
	Weight      int
	Couriers    []string
	PriceFilter *string
}

func (u *EstimateShippingCostUsecase) Execute(
	input EstimateShippingCostInput,
) ([]repository.CostOption, error) {
	if input.Origin == input.Destination {
		return nil, appErr.NewInvalidInput(domain.ErrInvalidRoute.Error())
	}

	if len(input.Couriers) == 0 {
		return nil, appErr.NewInvalidInput(domain.ErrNoCourierSelected.Error())
	}
	for _, courier := range input.Couriers {
		if strings.TrimSpace(courier) == "" {
			return nil, appErr.NewInvalidInput(domain.ErrInvalidCourier.Error())
		}
	}

	if input.Weight <= 0 {
		return nil, appErr.NewInvalidInput(domain.ErrInvalidWeight.Error())
	}

	query := repository.CalculateCostInput{
		OriginID:      input.Origin,
		DestinationID: input.Destination,
		Weight:        input.Weight,
		Couriers:      input.Couriers,
		PriceFilter:   input.PriceFilter,
	}

	costOptions, err := u.shippingCostProvider.CalculateCost(query)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate shipping cost: %w", err)
	}

	return costOptions, nil
}
