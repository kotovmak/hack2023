package store

import (
	"context"
	"database/sql"
	"hack2023/internal/app/model"
)

func (s *Store) GetServiceList(ctx context.Context) (tl []model.Service, err error) {
	data, err := s.db.QueryContext(
		ctx,
		`SELECT 
			ID,
  		UF_NAME,
			UF_DESCRIPTION
		FROM 
			z_services
		`)
	if err != nil && err != sql.ErrNoRows {
		return tl, err
	}
	// Обход результатов
	for data.Next() {
		var (
			description sql.NullString
		)
		p := model.Service{}
		err = data.Scan(
			&p.ID,
			&p.Name,
			&description,
		)
		if err != nil {
			return tl, err
		}
		p.Description = description.String
		tl = append(tl, p)
	}
	return tl, nil
}

func (s *Store) GetPravActList(ctx context.Context) (tl []model.PravAct, err error) {
	data, err := s.db.QueryContext(
		ctx,
		`SELECT 
			ID,
  		UF_NAME,
			UF_NADZOR_ORGAN_ID,
			UF_CONTROL_TYPE_ID
		FROM 
			z_prav_acts
		`)
	if err != nil && err != sql.ErrNoRows {
		return tl, err
	}
	// Обход результатов
	for data.Next() {
		p := model.PravAct{}
		var no_id, ct_id int
		err = data.Scan(
			&p.ID,
			&p.Name,
			&no_id,
			&ct_id,
		)
		if err != nil && err != sql.ErrNoRows {
			return tl, err
		}
		tl = append(tl, p)
	}
	return tl, nil
}

func (s *Store) SearchPravAct(ctx context.Context, text string) (cl model.PravAct, err error) {
	if err := s.db.QueryRowContext(
		ctx,
		`SELECT 
			ID,
			UF_NAME,
			UF_NADZOR_ORGAN_ID,
			UF_CONTROL_TYPE_ID
		FROM 
			z_prav_acts
		WHERE
			MATCH (UF_NAME) AGAINST (? IN BOOLEAN MODE)
		LIMIT 1
		`, text).Scan(
		&cl.ID,
		&cl.Name,
		&cl.NadzorOrganID,
		&cl.ControlTypeID,
	); err != nil && err != sql.ErrNoRows {
		return cl, err
	}
	return cl, nil
}

func (s *Store) GetNadzorOrganList(ctx context.Context) (tl map[int]model.NadzonOrgan, err error) {
	tl = make(map[int]model.NadzonOrgan)
	data, err := s.db.QueryContext(
		ctx,
		`SELECT 
			ID,
  		UF_NAME
		FROM 
			z_nadzor_organs
		`)
	if err != nil && err != sql.ErrNoRows {
		return tl, err
	}
	// Обход результатов
	for data.Next() {
		p := model.NadzonOrgan{}
		err = data.Scan(
			&p.ID,
			&p.Name,
		)
		if err != nil && err != sql.ErrNoRows {
			return tl, err
		}
		tl[p.ID] = p
	}
	return tl, nil
}

func (s *Store) GetNadzorOrganFilteredList(ctx context.Context, text string) (tl map[int]model.NadzonOrgan, err error) {
	tl = make(map[int]model.NadzonOrgan)
	data, err := s.db.QueryContext(
		ctx,
		`SELECT 
			ID,
  		UF_NAME
		FROM 
			z_nadzor_organs
		WHERE 
			MATCH (UF_NAME) AGAINST (? IN BOOLEAN MODE)
		`, text)
	if err != nil && err != sql.ErrNoRows {
		return tl, err
	}
	// Обход результатов
	for data.Next() {
		p := model.NadzonOrgan{}
		err = data.Scan(
			&p.ID,
			&p.Name,
		)
		if err != nil && err != sql.ErrNoRows {
			return tl, err
		}
		tl[p.ID] = p
	}
	return tl, nil
}

func (s *Store) GetConsultTopicList(ctx context.Context) (tl map[int]model.ConsultTopic, err error) {
	tl = make(map[int]model.ConsultTopic)
	data, err := s.db.QueryContext(
		ctx,
		`SELECT 
			ID,
  		UF_NAME,
			UF_NADZOR_ORGAN_ID,
			UF_CONTROL_TYPE_ID
		FROM 
			z_consult_topics
		`)
	if err != nil && err != sql.ErrNoRows {
		return tl, err
	}
	// Обход результатов
	for data.Next() {
		p := model.ConsultTopic{}
		var no_id int
		err = data.Scan(
			&p.ID,
			&p.Name,
			&no_id,
			&p.ControlTypeID,
		)
		if err != nil && err != sql.ErrNoRows {
			return tl, err
		}
		tl[p.ID] = p
	}
	return tl, nil
}

func (s *Store) GetControlTypeList(ctx context.Context) (tl map[int]model.ControlType, err error) {
	tl = make(map[int]model.ControlType)
	data, err := s.db.QueryContext(
		ctx,
		`SELECT 
			ID,
  		UF_NAME,
			UF_NADZOR_ORGAN_ID
		FROM 
			z_control_types
		`)
	if err != nil && err != sql.ErrNoRows {
		return tl, err
	}
	// Обход результатов
	for data.Next() {
		p := model.ControlType{}
		err = data.Scan(
			&p.ID,
			&p.Name,
			&p.NadzonOrganID,
		)
		if err != nil && err != sql.ErrNoRows {
			return tl, err
		}
		tl[p.ID] = p
	}
	return tl, nil
}
