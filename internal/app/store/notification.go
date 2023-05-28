package store

import (
	"context"
	"database/sql"
	"hack2023/internal/app/model"
	"time"
)

func (s *Store) GetNotificationList(ctx context.Context, userID int) ([]model.Notification, error) {
	cl := []model.Notification{}
	var (
		text sql.NullString
		date sql.NullTime
	)
	data, err := s.db.QueryContext(
		ctx,
		`SELECT 
			ID,
			UF_DATE,
  		UF_TEXT,
			UF_USER_ID
		FROM 
			z_notifications
		WHERE
			UF_USER_ID = ?
		ORDER BY
			UF_DATE DESC
		`, userID)
	if err != nil && err != sql.ErrNoRows {
		return cl, err
	}
	// Обход результатов
	for data.Next() {
		p := model.Notification{}
		err = data.Scan(
			&p.ID,
			&date,
			&text,
			&p.UserID,
		)
		if err != nil {
			return cl, err
		}
		p.Date = date.Time
		p.DateExport = date.Time.Format("2006-01-02")
		p.Text = text.String
		cl = append(cl, p)
	}
	return cl, nil
}

func (s *Store) AddNotification(ctx context.Context, n model.Notification) error {
	return s.db.QueryRowContext(
		ctx,
		`INSERT INTO 
			z_notifications 
			(UF_USER_ID, UF_DATE, UF_TEXT)
		VALUES 
			(?, ?, ?)`,
		n.UserID,
		time.Now(),
		n.Text,
	).Err()
}
