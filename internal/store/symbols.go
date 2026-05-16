package store

import (
	"context"
)

type UserSymbol struct {
	ID                   string
	UserID               string
	Symbol               string
	TraderType           string
	MinROI               float64
	MinRR                float64
	MinConfidence        float64
	TradeXSubscriptionID string
	Active               bool
}

func (s *Store) CreateUserSymbol(ctx context.Context, us *UserSymbol) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO user_symbols (id, user_id, symbol, trader_type, min_roi, min_rr, min_confidence, tradex_subscription_id, active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, TRUE)`,
		us.ID, us.UserID, us.Symbol, us.TraderType, us.MinROI, us.MinRR, us.MinConfidence, nullStr(us.TradeXSubscriptionID),
	)
	return err
}

func (s *Store) GetUserSymbols(ctx context.Context, userID string) ([]*UserSymbol, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, user_id, symbol, trader_type, min_roi, min_rr, min_confidence,
		        COALESCE(tradex_subscription_id,''), active
		 FROM user_symbols WHERE user_id = $1 AND active = TRUE
		 ORDER BY created_at`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*UserSymbol
	for rows.Next() {
		us := &UserSymbol{}
		if err := rows.Scan(&us.ID, &us.UserID, &us.Symbol, &us.TraderType,
			&us.MinROI, &us.MinRR, &us.MinConfidence, &us.TradeXSubscriptionID, &us.Active); err != nil {
			return nil, err
		}
		result = append(result, us)
	}
	return result, rows.Err()
}

func (s *Store) GetUserSymbolByID(ctx context.Context, id, userID string) (*UserSymbol, error) {
	us := &UserSymbol{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, symbol, trader_type, min_roi, min_rr, min_confidence,
		        COALESCE(tradex_subscription_id,''), active
		 FROM user_symbols WHERE id = $1 AND user_id = $2 AND active = TRUE`,
		id, userID,
	).Scan(&us.ID, &us.UserID, &us.Symbol, &us.TraderType,
		&us.MinROI, &us.MinRR, &us.MinConfidence, &us.TradeXSubscriptionID, &us.Active)
	if err != nil {
		return nil, err
	}
	return us, nil
}

func (s *Store) UpdateUserSymbol(ctx context.Context, id, userID string, minROI, minRR, minConfidence float64) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE user_symbols SET min_roi=$3, min_rr=$4, min_confidence=$5
		 WHERE id=$1 AND user_id=$2 AND active=TRUE`,
		id, userID, minROI, minRR, minConfidence,
	)
	return err
}

func (s *Store) DeactivateUserSymbol(ctx context.Context, id, userID string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE user_symbols SET active=FALSE WHERE id=$1 AND user_id=$2`,
		id, userID,
	)
	return err
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
