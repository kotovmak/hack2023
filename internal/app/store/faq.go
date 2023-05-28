package store

import (
	"context"
	"database/sql"
	"hack2023/internal/app/model"
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
