package persistence

import (
	"context"
	"fmt"
	"strings"
	"time"

	database "service-core/internal/infra/db"

	"service-core/internal/modules/courier/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type courierRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewCourierRepositoryImpl(conn *database.Connection) repository.CourierRepository {
	return &courierRepositoryImpl{
		db: conn.Pool,
	}
}

func (r *courierRepositoryImpl) GetActiveCodes(codes []string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT code
		FROM couriers
		WHERE code = ANY($1)
		AND is_active = true
	`

	rows, err := r.db.Query(ctx, query, codes)
	if err != nil {
		return nil, fmt.Errorf("query active couriers failed: %w", err)
	}
	defer rows.Close()

	var result []string

	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, fmt.Errorf("scan courier code failed: %w", err)
		}
		result = append(result, code)
	}

	return result, nil
}

func (r *courierRepositoryImpl) ValidateCouriers(codes []string) ([]string, error) {
	if len(codes) == 0 {
		return nil, fmt.Errorf("courier codes cannot be empty")
	}

	normalized := r.normalizeCodes(codes)
	validCodes, err := r.GetActiveCodes(normalized)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve couriers: %w", err)
	}

	if len(validCodes) != len(normalized) {
		return nil, fmt.Errorf("one or more courier codes are invalid")
	}

	return validCodes, nil
}

func (r *courierRepositoryImpl) normalizeCodes(codes []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(codes))

	for _, code := range codes {
		code = strings.ToLower(strings.TrimSpace(code))

		if code == "" {
			continue
		}

		if _, exists := seen[code]; !exists {
			result = append(result, code)
			seen[code] = struct{}{}
		}
	}

	return result
}
