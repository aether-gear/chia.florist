package persistence

import (
	"context"
	"errors"
	"fmt"

	appclock "service-core/internal/common/clock"
	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type paymentChannelDataRepositoryImpl struct{}

func NewPaymentChannelDataRepositoryImpl() repository.PaymentChannelDataRepository {
	return &paymentChannelDataRepositoryImpl{}
}

func (r *paymentChannelDataRepositoryImpl) Save(
	ctx context.Context,
	exec transaction.Executor,
	data domain.PaymentChannelData,
) error {
	query := `
		INSERT INTO payment_channel_data (
			id,
			payment_id,
			channel_type,
			display_name,
			action_url,
			expires_at,
			created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (payment_id)
		DO UPDATE SET
			channel_type = EXCLUDED.channel_type,
			display_name = EXCLUDED.display_name,
			action_url   = EXCLUDED.action_url,
			expires_at   = EXCLUDED.expires_at
	`

	_, err := exec.Exec(ctx, query,
		data.ID,
		data.PaymentID,
		data.ChannelType,
		data.DisplayName,
		data.ActionURL,
		data.ExpiresAt,
		data.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("save payment channel data failed: %w", err)
	}

	return nil
}

func (r *paymentChannelDataRepositoryImpl) GetByPaymentID(
	ctx context.Context,
	exec transaction.Executor,
	paymentID uuid.UUID,
) (*domain.PaymentChannelData, error) {
	query := `
		SELECT
			id,
			payment_id,
			channel_type,
			display_name,
			action_url,
			expires_at,
			created_at
		FROM payment_channel_data
		WHERE
			payment_id = $1
		LIMIT 1
	`

	var d domain.PaymentChannelData
	err := exec.QueryRow(ctx, query, paymentID).Scan(
		&d.ID,
		&d.PaymentID,
		&d.ChannelType,
		&d.DisplayName,
		&d.ActionURL,
		&d.ExpiresAt,
		&d.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query payment channel data by payment_id failed: %w", err)
	}

	return &d, nil
}

func (r *paymentChannelDataRepositoryImpl) ListByPaymentIDs(
	ctx context.Context,
	exec transaction.Executor,
	paymentIDs []uuid.UUID,
) (map[uuid.UUID]*domain.PaymentChannelData, error) {
	if len(paymentIDs) == 0 {
		return map[uuid.UUID]*domain.PaymentChannelData{}, nil
	}

	query := `
		SELECT
			id,
			payment_id,
			channel_type,
			display_name,
			action_url,
			expires_at,
			created_at
		FROM payment_channel_data
		WHERE
			payment_id = ANY($1::uuid[])
	`

	idStrings := make([]string, len(paymentIDs))
	for i, id := range paymentIDs {
		idStrings[i] = id.String()
	}

	rows, err := exec.Query(ctx, query, idStrings)
	if err != nil {
		return nil, fmt.Errorf("query payment channel data by payment ids failed: %w", err)
	}
	defer rows.Close()

	result := make(map[uuid.UUID]*domain.PaymentChannelData)
	for rows.Next() {
		var d domain.PaymentChannelData
		if err := rows.Scan(
			&d.ID,
			&d.PaymentID,
			&d.ChannelType,
			&d.DisplayName,
			&d.ActionURL,
			&d.ExpiresAt,
			&d.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan payment channel data failed: %w", err)
		}
		copied := d
		result[d.PaymentID] = &copied
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate payment channel data rows failed: %w", err)
	}

	// Initialise the created_at zero value
	// for rows that somehow lack it
	for k, v := range result {
		if v.CreatedAt.IsZero() {
			v.CreatedAt = appclock.Now()
			result[k] = v
		}
	}

	return result, nil
}
