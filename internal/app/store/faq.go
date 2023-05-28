package store

import (
	"context"
	"database/sql"
	"hack2023/internal/app/model"

	_ "github.com/go-sql-driver/mysql"
)

func (s *Store) GetFAQList(ctx context.Context) ([]model.FAQ, error) {
	cl := []model.FAQ{}
	var (
		question sql.NullString
		answer   sql.NullString
		date     sql.NullTime
	)
	data, err := s.db.QueryContext(
		ctx,
		`SELECT 
			ID,
			UF_DATE,
			UF_QUESTION,
  		UF_ANSWER,
			UF_NADZOR_ORGAN_ID,
			UF_CONTROL_TYPE_ID,
			UF_LIKE
		FROM 
			z_faq
		`)
	if err != nil && err != sql.ErrNoRows {
		return cl, err
	}
	// Обход результатов
	for data.Next() {
		p := model.FAQ{}
		err = data.Scan(
			&p.ID,
			&date,
			&question,
			&answer,
			&p.NadzorOrganID,
			&p.ControlTypeID,
			&p.Likes,
		)
		if err != nil {
			return cl, err
		}
		p.Date = date.Time
		p.Answer = answer.String
		p.Question = question.String
		cl = append(cl, p)
	}
	return cl, nil
}

func (s *Store) SearchFAQ(ctx context.Context, text string) (cl model.FAQ, err error) {
	var (
		question sql.NullString
		answer   sql.NullString
		date     sql.NullTime
	)
	if err := s.db.QueryRowContext(
		ctx,
		`SELECT 
			ID,
			UF_DATE,
			UF_QUESTION,
  		UF_ANSWER,
			UF_NADZOR_ORGAN_ID,
			UF_CONTROL_TYPE_ID,
			UF_LIKE
		FROM 
			z_faq
		WHERE
			MATCH (UF_QUESTION,UF_ANSWER) AGAINST (? IN BOOLEAN MODE)
		LIMIT 1
		`, text).Scan(
		&cl.ID,
		&date,
		&question,
		&answer,
		&cl.NadzorOrganID,
		&cl.ControlTypeID,
		&cl.Likes,
	); err != nil && err != sql.ErrNoRows {
		return cl, err
	}
	cl.Date = date.Time
	cl.Answer = answer.String
	cl.Question = question.String
	return cl, nil
}

func (s *Store) SearchFAQList(ctx context.Context, text string) (cl []model.FAQ, err error) {
	var (
		question sql.NullString
		answer   sql.NullString
		date     sql.NullTime
	)
	data, err := s.db.QueryContext(
		ctx,
		`SELECT 
			ID,
			UF_DATE,
			UF_QUESTION,
  		UF_ANSWER,
			UF_NADZOR_ORGAN_ID,
			UF_CONTROL_TYPE_ID,
			UF_LIKE
		FROM 
			z_faq
		WHERE
			MATCH (UF_QUESTION,UF_ANSWER) AGAINST (? IN BOOLEAN MODE)
		`, text)
	if err != nil && err != sql.ErrNoRows {
		return cl, err
	}
	// Обход результатов
	for data.Next() {
		p := model.FAQ{}
		err = data.Scan(
			&p.ID,
			&date,
			&question,
			&answer,
			&p.NadzorOrganID,
			&p.ControlTypeID,
			&p.Likes,
		)
		if err != nil {
			return cl, err
		}
		p.Date = date.Time
		p.Answer = answer.String
		p.Question = question.String
		cl = append(cl, p)
	}
	return cl, nil
}
