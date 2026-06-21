package persistence

import (
	"context"
	"errors"
	"fmt"

	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type paymentInstructionRepositoryImpl struct{}

func NewPaymentInstructionRepositoryImpl() repository.PaymentInstructionRepository {
	return &paymentInstructionRepositoryImpl{}
}

func (r *paymentInstructionRepositoryImpl) Save(
	ctx context.Context,
	exec transaction.Executor,
	instruction domain.PaymentInstruction,
) error {
	query := `
		INSERT INTO payment_instructions (
			id,
			payment_method_id,
			content,
			created_at
		) VALUES ($1, $2, $3, $4)
	`

	_, err := exec.Exec(ctx, query,
		instruction.ID,
		instruction.PaymentID,
		instruction.Content,
		instruction.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("save payment instruction failed: %w", err)
	}

	return nil
}

func (r *paymentInstructionRepositoryImpl) GetByPaymentID(
	ctx context.Context,
	exec transaction.Executor,
	paymentID uuid.UUID,
) (*domain.PaymentInstruction, error) {
	query := `
		SELECT
			id,
			payment_method_id,
			content,
			created_at
		FROM payment_instructions
		WHERE payment_method_id = $1
		LIMIT 1
	`

	var inst domain.PaymentInstruction
	err := exec.QueryRow(ctx, query, paymentID).Scan(
		&inst.ID,
		&inst.PaymentID,
		&inst.Content,
		&inst.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query payment instruction by payment_method_id failed: %w", err)
	}

	return &inst, nil
}
