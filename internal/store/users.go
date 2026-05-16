package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	PasswordHash   string `json:"-"`
	TelegramChatID string `json:"telegram_chat_id"`
	AlpacaKey      string `json:"-"`
	AlpacaSecret   string `json:"-"`
	WebhookToken   string `json:"webhook_token"`
}

func (s *Store) CreateUser(ctx context.Context, u *User) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, webhook_token)
		 VALUES ($1, $2, $3, $4)`,
		u.ID, u.Email, u.PasswordHash, u.WebhookToken,
	)
	return err
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, COALESCE(telegram_chat_id,''),
		        COALESCE(alpaca_key,''), COALESCE(alpaca_secret,''), webhook_token
		 FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.TelegramChatID, &u.AlpacaKey, &u.AlpacaSecret, &u.WebhookToken)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (s *Store) GetUserByToken(ctx context.Context, token string) (*User, error) {
	u := &User{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, COALESCE(telegram_chat_id,''),
		        COALESCE(alpaca_key,''), COALESCE(alpaca_secret,''), webhook_token
		 FROM users WHERE webhook_token = $1`,
		token,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.TelegramChatID, &u.AlpacaKey, &u.AlpacaSecret, &u.WebhookToken)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (s *Store) GetUserByID(ctx context.Context, id string) (*User, error) {
	u := &User{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, COALESCE(telegram_chat_id,''),
		        COALESCE(alpaca_key,''), COALESCE(alpaca_secret,''), webhook_token
		 FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.TelegramChatID, &u.AlpacaKey, &u.AlpacaSecret, &u.WebhookToken)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (s *Store) GetUserByTelegramChatID(ctx context.Context, chatID string) (*User, error) {
	u := &User{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, COALESCE(telegram_chat_id,''),
		        COALESCE(alpaca_key,''), COALESCE(alpaca_secret,''), webhook_token
		 FROM users WHERE telegram_chat_id = $1`,
		chatID,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.TelegramChatID, &u.AlpacaKey, &u.AlpacaSecret, &u.WebhookToken)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (s *Store) UpdateUserProfile(ctx context.Context, id, telegramChatID, alpacaKey, alpacaSecret string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET telegram_chat_id=NULLIF($2,''), alpaca_key=NULLIF($3,''), alpaca_secret=NULLIF($4,'')
		 WHERE id = $1`,
		id, telegramChatID, alpacaKey, alpacaSecret,
	)
	return err
}
