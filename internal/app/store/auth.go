package store

import (
	"context"
	"database/sql"
	"hack2023/internal/app/model"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

func (s *Store) SaveAppToken(ctx context.Context, token model.Token) (cl model.Token, err error) {
	cl = token
	stmt, err := s.db.PrepareContext(
		ctx,
		`INSERT INTO 
			z_app_tokens 
			(UF_USER_ID, UF_APP_TOKEN)
		VALUES 
			(?, ?)`,
	)
	if err != nil {
		return cl, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		token.UserID,
		token.Token,
	)
	if err != nil {
		return cl, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return cl, err
	}
	cl.ID = int(id)
	return cl, nil
}

func (s *Store) GetAppTokenList(ctx context.Context, userID int) ([]model.Token, error) {
	cl := []model.Token{}
	data, err := s.db.QueryContext(
		ctx,
		`SELECT 
			ID,
			UF_USER_ID,
  		UF_APP_TOKEN
		FROM 
			z_app_tokens
		WHERE
			UF_USER_ID = ?
		`, userID)
	if err != nil && err != sql.ErrNoRows {
		return cl, err
	}
	// Обход результатов
	for data.Next() {
		p := model.Token{}
		err = data.Scan(
			&p.ID,
			&p.UserID,
			&p.Token,
		)
		if err != nil {
			return cl, err
		}
		cl = append(cl, p)
	}
	return cl, nil
}

func (s *Store) DeleteAppToken(ctx context.Context, token string) error {
	return s.db.QueryRowContext(
		ctx,
		`DELETE from 
			z_app_tokens
		WHERE 
			UF_APP_TOKEN = ?`,
		token,
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
	var (
		isKNO         sql.NullBool
		nadzonOrganID sql.NullInt16
	)
	if err := s.db.QueryRowContext(
		ctx,
		`SELECT 
			ID,
			UF_EMAIL,
			UF_NAME,
			UF_LOGIN,
			UF_PASSWORD,
			UF_IS_KNO,
			UF_NADZOR_ORGAN_ID
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
		&isKNO,
		&nadzonOrganID,
	); err != nil {
		return user, err
	}
	user.IsKNO = isKNO.Bool
	user.NadzonOrganID = int(nadzonOrganID.Int16)
	return user, nil
}

func (s *Store) GetUserByID(ctx context.Context, id int) (user model.Account, err error) {
	var (
		isKNO         sql.NullBool
		nadzonOrganID sql.NullInt16
	)
	if err := s.db.QueryRowContext(
		ctx,
		`SELECT 
			ID,
			UF_EMAIL,
			UF_NAME,
			UF_LOGIN,
			UF_PASSWORD,
			UF_IS_KNO,
			UF_NADZOR_ORGAN_ID
		FROM 
			z_api_users
		WHERE
			ID = ? AND UF_ACTIVE = 'Y'`,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Login,
		&user.Password,
		&isKNO,
		&nadzonOrganID,
	); err != nil {
		return user, err
	}
	user.IsKNO = isKNO.Bool
	user.NadzonOrganID = int(nadzonOrganID.Int16)
	return user, nil
}
