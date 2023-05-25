package store

import (
	"context"
	"database/sql"
	"hack2023/internal/app/model"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
		WHERE
			UF_IS_BUSY IS NULL OR UF_IS_BUSY = 0
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
		p.DateExport = p.Date.Format("2006-01-02")
		sl[p.DateExport] = append(sl[p.DateExport], p)
	}
	return sl, nil
}

func (s *Store) GetConsultationList(ctx context.Context) (model.Consultations, error) {
	cl := model.Consultations{}
	var (
		question     sql.NullString
		date         sql.NullTime
		isNeedLetter sql.NullBool
		isConfirmed  sql.NullBool
		VKSLink      sql.NullString
		videoLink    sql.NullString
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
			UF_IS_CONFIRMED,
			UF_VKS_LINK,
			UF_VIDEO_LINK,
			UF_SLOT_ID
		FROM 
			z_consultations
		WHERE
			UF_IS_DELETED IS NULL
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
			&date,
			&question,
			&p.NadzonOrganID,
			&p.ControlTypeID,
			&p.ConsultTopicID,
			&p.UserID,
			&isNeedLetter,
			&isConfirmed,
			&VKSLink,
			&videoLink,
			&p.SlotID,
		)
		if err != nil {
			return cl, err
		}
		p.Date = date.Time
		p.DateExport = date.Time.Format("2006-01-02")
		p.Question = question.String
		p.IsConfirmed = isConfirmed.Bool
		p.IsNeedLetter = isNeedLetter.Bool
		p.VKSLink = VKSLink.String
		p.VideoLink = videoLink.String
		if p.Date.Unix() > time.Now().Unix() {
			cl.Active = append(cl.Active, p)
		} else {
			cl.Finished = append(cl.Finished, p)
		}
	}
	return cl, nil
}

func (s *Store) AddConsultation(ctx context.Context, cl model.Consultation) (model.Consultation, error) {
	stmt, err := s.db.PrepareContext(
		ctx,
		`INSERT INTO
			 z_consultations (
				UF_TIME,
				UF_DATE,
				UF_QUESTION,
				UF_NADZOR_ORGAN_ID,
				UF_CONTROL_TYPE_ID,
				UF_CONSULT_TOPIC_ID,
				UF_USER_ID,
				UF_IS_NEED_LATTER,
				UF_SLOT_ID
			 )
		VALUES 
			(?,?,?,?,?,?,?,?,?);
		`,
	)
	if err != nil {
		return cl, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		cl.Time,
		cl.Date,
		cl.Question,
		cl.NadzonOrganID,
		cl.ControlTypeID,
		cl.ConsultTopicID,
		cl.UserID,
		cl.IsNeedLetter,
		cl.SlotID,
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

func (s *Store) CloseSlot(ctx context.Context, slotID int) error {
	if _, err := s.db.QueryContext(
		ctx,
		`UPDATE 
			 z_slots 
		SET
			UF_IS_BUSY = 1
		WHERE
			ID = ? AND (UF_IS_BUSY IS NULL OR UF_IS_BUSY = 0)
		`,
		slotID,
	); err != nil {
		return err
	}
	return nil
}

func (s *Store) OpenSlot(ctx context.Context, slotID int) error {
	if _, err := s.db.QueryContext(
		ctx,
		`UPDATE 
			 z_slots 
		SET
			UF_IS_BUSY = 0
		WHERE
			ID = ? AND UF_IS_BUSY = 1
		`,
		slotID,
	); err != nil {
		return err
	}
	return nil
}

func (s *Store) DeleteConsultation(ctx context.Context, consultationID string) error {
	if _, err := s.db.QueryContext(
		ctx,
		`UPDATE 
			z_consultations
		SET
			UF_IS_DELETED = 1
		WHERE
			ID = ? AND (UF_IS_DELETED IS NULL OR UF_IS_DELETED = 0)
		`,
		consultationID,
	); err != nil {
		return err
	}
	return nil
}

func (s *Store) GetConsultation(ctx context.Context, consultationID string) (p model.Consultation, err error) {
	var (
		question     sql.NullString
		date         sql.NullTime
		isNeedLetter sql.NullBool
		isConfirmed  sql.NullBool
		VKSLink      sql.NullString
		videoLink    sql.NullString
		isDeleted    sql.NullBool
	)
	if err := s.db.QueryRowContext(ctx,
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
			UF_IS_CONFIRMED,
			UF_VKS_LINK,
			UF_VIDEO_LINK,
			UF_SLOT_ID,
			UF_IS_DELETED
		FROM 
			z_consultations
		WHERE
			ID = ? AND (UF_IS_DELETED IS NULL OR UF_IS_DELETED = 0)`,
		consultationID,
	).Scan(
		&p.ID,
		&p.Time,
		&date,
		&question,
		&p.NadzonOrganID,
		&p.ControlTypeID,
		&p.ConsultTopicID,
		&p.UserID,
		&isNeedLetter,
		&isConfirmed,
		&VKSLink,
		&videoLink,
		&p.SlotID,
		&isDeleted,
	); err != nil {
		return p, err
	}
	p.Date = date.Time
	p.DateExport = date.Time.Format("2006-01-02")
	p.Question = question.String
	p.IsConfirmed = isConfirmed.Bool
	p.IsNeedLetter = isNeedLetter.Bool
	p.IsDeleted = isDeleted.Bool
	p.VKSLink = VKSLink.String
	p.VideoLink = videoLink.String
	return p, nil
}

func (s *Store) GetSlot(ctx context.Context, slotID int) (p model.Slot, err error) {
	if err := s.db.QueryRowContext(ctx,
		`SELECT 
			ID,
  		UF_TIME,
			UF_DATE
		FROM 
			z_slots
		WHERE
			ID = ? AND (UF_IS_BUSY IS NULL OR UF_IS_BUSY = 0)`,
		slotID,
	).Scan(
		&p.ID,
		&p.Time,
		&p.Date,
	); err != nil {
		return p, err
	}
	p.DateExport = p.Date.Format("2006-01-02")
	return p, nil
}
