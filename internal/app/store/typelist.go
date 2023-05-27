package store

import (
	"context"
	"database/sql"
	"hack2023/internal/app/model"
)

func (s *Store) GetTypeList(ctx context.Context) (tl model.TypeList, err error) {
	topics := make(map[int][]model.ConsultTopic)
	types := make(map[int][]model.ControlType)

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
		tl.Services = append(tl.Services, p)
	}

	data, err = s.db.QueryContext(
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
		topics[ct_id] = append(topics[ct_id], p)
	}

	data, err = s.db.QueryContext(
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
		var no_id int
		p := model.ControlType{}
		err = data.Scan(
			&p.ID,
			&p.Name,
			&no_id,
		)
		if err != nil && err != sql.ErrNoRows {
			return tl, err
		}
		p.ConsultTopics = topics[p.ID]
		types[no_id] = append(types[no_id], p)
	}

	data, err = s.db.QueryContext(
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
		p.ControlTypes = types[p.ID]
		tl.NadzonOrgans = append(tl.NadzonOrgans, p)
	}

	data, err = s.db.QueryContext(
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
		tl.PravActs = append(tl.PravActs, p)
	}

	return tl, nil
}
