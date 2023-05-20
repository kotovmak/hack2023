package store

import (
	"context"
	"hack2023/internal/app/model"
	"time"
)

func (s *Store) SaveRefreshToken(ctx context.Context, userID int, token string) error {
	return s.db.QueryRowContext(
		ctx,
		`INSERT INTO 
			z_refresh_tokens 
			(UF_USER_ID, UF_REFRESH_TOKEN, UF_UPDATED_AT)
		VALUES 
			(?, ?, ?) 
		ON DUPLICATE KEY UPDATE
			UF_REFRESH_TOKEN = VALUES(UF_REFRESH_TOKEN),
			UF_UPDATED_AT = VALUES(UF_UPDATED_AT);`,
		userID,
		token,
		time.Now(),
	).Err()
}

func (s *Store) DeleteRefreshToken(ctx context.Context, userID int) error {
	return s.db.QueryRowContext(
		ctx,
		`DELETE from 
			z_refresh_tokens
		WHERE 
			UF_USER_ID = ?`,
		userID,
	).Err()
}

func (s *Store) GetUserByRefreshToken(ctx context.Context, refresh_token string) (user model.Account, err error) {
	var userID int
	if err := s.db.QueryRowContext(
		ctx,
		`SELECT 
			UF_USER_ID
		FROM 
			z_refresh_tokens
		WHERE
			UF_REFRESH_TOKEN = ?
		`,
		refresh_token,
	).Scan(
		&userID,
	); err != nil {
		return user, err
	}
	if err := s.db.QueryRowContext(
		ctx,
		`SELECT 
			ID,
			UF_EMAIL,
			UF_NAME,
			UF_PASSWORD,
			UF_LOGIN
		FROM 
			z_api_users
		WHERE
			ID = ?
		`,
		userID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.Login,
	); err != nil {
		return user, err
	}
	return user, nil
}

func (s *Store) GetUserByLogin(ctx context.Context, login string) (user model.Account, err error) {
	if err := s.db.QueryRowContext(
		ctx,
		`SELECT 
			ID,
			UF_EMAIL,
			UF_NAME,
			UF_LOGIN,
			UF_PASSWORD
		FROM 
			z_api_users
		WHERE
			UF_LOGIN = ? AND UF_ACTIVE = 'Y'`,
		login,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Login,
		&user.Password,
	); err != nil {
		return user, err
	}
	return user, nil
}
