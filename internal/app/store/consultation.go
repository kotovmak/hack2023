package store

import (
	"context"
	"database/sql"
	"hack2023/internal/app/model"

	_ "github.com/go-sql-driver/mysql"
)

func (s *Store) GetConsultationList(ctx context.Context) (map[int]model.Consultation, error) {
	cl := make(map[int]model.Consultation)
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

		cl[p.SlotID] = p
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
				UF_SLOT_ID,
				UF_VKS_LINK
			 )
		VALUES 
			(?,?,?,?,?,?,?,?,?,?);
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
		cl.VKSLink,
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

func (s *Store) ApplyConsultation(ctx context.Context, consultationID string) error {
	if _, err := s.db.QueryContext(
		ctx,
		`UPDATE 
			z_consultations
		SET
			UF_IS_CONFIRMED = 1
		WHERE
			ID = ? AND (UF_IS_DELETED IS NULL OR UF_IS_DELETED = 0)
		`,
		consultationID,
	); err != nil {
		return err
	}
	return nil
}
