package store

import (
	"context"
	"database/sql"
	"hack2023/internal/app/model"
	"log"
)

func (s *Store) GetMessagesList(ctx context.Context, userID int) ([]model.Message, error) {
	cl := []model.Message{}
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
			UF_SEND_BY_ID,
			UF_USER_ID
		FROM 
			z_messages
		WHERE
			UF_USER_ID = ?
		ORDER BY
			UF_DATE ASC
		`, userID)
	if err != nil {
		return cl, err
	}
	log.Println(data, err)
	// Обход результатов
	for data.Next() {
		p := model.Message{}
		err = data.Scan(
			&p.ID,
			&date,
			&text,
			&p.SendByID,
			&p.UserID,
		)
		if err != nil {
			return cl, err
		}
		p.Date = date.Time
		p.Text = text.String
		cl = append(cl, p)
	}
	return cl, nil
}

func (s *Store) AddMessage(ctx context.Context, cl model.Message) (model.Message, error) {
	stmt, err := s.db.PrepareContext(
		ctx,
		`INSERT INTO
			 z_messages (
				UF_DATE,
				UF_TEXT,
				UF_SEND_BY_ID,
				UF_USER_ID
			 )
		VALUES 
			(?,?,?,?);
		`,
	)
	if err != nil {
		return cl, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		cl.Date,
		cl.Text,
		cl.SendByID,
		cl.UserID,
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

func (s *Store) GetButtonList(ctx context.Context) ([]model.Button, error) {
	cl := []model.Button{}
	var (
		text sql.NullString
		link sql.NullString
	)
	data, err := s.db.QueryContext(
		ctx,
		`SELECT 
			ID,
  		UF_TEXT,
			UF_LINK
		FROM 
			z_chat_buttons
		`)
	if err != nil && err != sql.ErrNoRows {
		return cl, err
	}
	// Обход результатов
	for data.Next() {
		p := model.Button{}
		err = data.Scan(
			&p.ID,
			&text,
			&link,
		)
		if err != nil {
			return cl, err
		}
		p.Text = text.String
		p.Link = link.String
		cl = append(cl, p)
	}
	return cl, nil
}
