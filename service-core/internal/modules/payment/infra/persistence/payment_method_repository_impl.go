package persistence

import (
	"context"
	"errors"
	database "service-core/internal/infra/db"
	"service-core/internal/modules/payment/domain"
	"service-core/internal/modules/payment/repository"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type paymentMethodRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewPaymentMethodRepository(conn *database.Connection) repository.PaymentMethodRepository {
	return &paymentMethodRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *paymentMethodRepositoryImpl) Save(method domain.PaymentMethod) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO payment_methods (
			id, name, type, is_active, description,
			fee_type, fee_amount, fee_rate,
			created_at, updated_at
		) VALUES (
		 	$1,$2,$3,$4,$5,$6,$7,$8,$9,$10
		)
	`

	_, err := r.db.Exec(ctx, query,
		method.ID,
		method.Name,
		method.Type,
		method.IsActive,
		method.Description,
		method.FeeType,
		method.FeeFixed,
		method.FeePercentage,
		method.CreatedAt,
		method.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *paymentMethodRepositoryImpl) FindByName(name string) (*domain.PaymentMethod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, type, is_active, description,
			   created_at, updated_at
		FROM payment_methods
		WHERE name LIKE $1 || '%'
		LIMIT 1
	`

	var method domain.PaymentMethod

	err := r.db.QueryRow(ctx, query, name).Scan(
		&method.ID,
		&method.Name,
		&method.Type,
		&method.IsActive,
		&method.Description,
		&method.CreatedAt,
		&method.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &method, nil
}

func (r *paymentMethodRepositoryImpl) GetByID(paymentID uuid.UUID) (*domain.PaymentMethod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT 
			id,
			name,
			type,
			is_active,
			description,
			created_at,
			updated_at
		FROM payment_methods
		WHERE id = $1
		LIMIT 1
	`

	var method domain.PaymentMethod

	err := r.db.QueryRow(ctx, query, paymentID).Scan(
		&method.ID,
		&method.Name,
		&method.Type,
		&method.IsActive,
		&method.Description,
		&method.CreatedAt,
		&method.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &method, nil
}

func (r *paymentMethodRepositoryImpl) ListAll() ([]domain.PaymentMethod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT 
			id,
			name,
			type,
			is_active,
			description,
			created_at,
			updated_at
		FROM payment_methods
		WHERE deleted_at IS NULL
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.PaymentMethod

	for rows.Next() {
		var row domain.PaymentMethod

		err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.Type,
			&row.IsActive,
			&row.Description,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, row)

	}

	return result, nil
}
