package store

import (
	"context"
	"database/sql"
	"hack2023/internal/app/model"
	"time"
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

func (s *Store) GetSlotList(ctx context.Context) (map[string][]model.Slot, error) {
	sl := make(map[string][]model.Slot)
	data, err := s.db.QueryContext(
		ctx,
		`SELECT 
			ID,
  		UF_TIME,
			UF_DATE
		FROM 
			z_slots
		`)
	if err != nil && err != sql.ErrNoRows {
		return sl, err
	}
	// Обход результатов
	for data.Next() {
		p := model.Slot{}
		err = data.Scan(
			&p.ID,
			&p.Time,
			&p.Date,
		)
		if err != nil {
			return sl, err
		}
		sl[p.Date] = append(sl[p.Date], p)
	}
	return sl, nil
}

func (s *Store) GetConsultationList(ctx context.Context) (model.Consultations, error) {
	cl := model.Consultations{}
	var (
		question     sql.NullString
		isNeedLetter sql.NullBool
		isConfirmed  sql.NullBool
	)
	data, err := s.db.QueryContext(
		ctx,
		`SELECT 
			ID,
  		UF_TIME,
			UF_DATE,
  		UF_QUESTION,
  		UF_NADZOR_ORGAN_ID,
  		UF_CONTROL_TYPE_ID,
  		UF_CONSULT_TOPIC_ID,
  		UF_USER_ID,
  		UF_IS_NEED_LATTER,
			UF_IS_CONFIRMED
		FROM 
			z_consultations
		`)
	if err != nil && err != sql.ErrNoRows {
		return cl, err
	}
	// Обход результатов
	for data.Next() {
		p := model.Consultation{}
		err = data.Scan(
			&p.ID,
			&p.Time,
			&p.Date,
			&question,
			&p.NadzonOrganID,
			&p.ControlTypeID,
			&p.ConsultTopicID,
			&p.UserID,
			&isNeedLetter,
			&isConfirmed,
		)
		if err != nil {
			return cl, err
		}
		p.Question = question.String
		p.IsConfirmed = isConfirmed.Bool
		p.IsNeedLetter = isNeedLetter.Bool
		if p.Date.Unix() > time.Now().Unix() {
			cl.Active = append(cl.Active, p)
		} else {
			cl.Finished = append(cl.Finished, p)
		}
	}
	return cl, nil
}
