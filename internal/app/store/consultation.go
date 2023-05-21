package store

import (
	"context"
	"database/sql"
	"hack2023/internal/app/model"
)

func (s *Store) GetTypeList(ctx context.Context) (tl model.TypeList, err error) {
	no := make(map[int]model.NadzonOrgan)
	ct := make(map[int]model.ControlType)
	topics := make(map[int][]model.ConsultTopic)

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
		tl.NadzonOrgans = append(tl.NadzonOrgans, p)
		no[p.ID] = p
	}

	data, err = s.db.QueryContext(
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
		var no_id int
		err = data.Scan(
			&p.ID,
			&p.Name,
			&no_id,
		)
		if err != nil && err != sql.ErrNoRows {
			return tl, err
		}
		p.NadzonOrgan = no[no_id]
		ct[p.ID] = p
		tl.ControlTypes = append(tl.ControlTypes, p)
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
		p.ControlType = ct[ct_id]
		topics[no_id] = append(topics[no_id], p)
		tl.ConsultTopics = append(tl.ConsultTopics, p)
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
		p.ControlType = ct[ct_id]
		tl.PravActs = append(tl.PravActs, p)
	}

	for i, v := range tl.NadzonOrgans {
		tl.NadzonOrgans[i].ConsultTopic = topics[v.ID]
	}

	return tl, nil
}
