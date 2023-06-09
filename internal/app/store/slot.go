package store

import (
	"context"
	"database/sql"
	"hack2023/internal/app/model"
)

func (s *Store) GetSlotList(ctx context.Context, isKNO bool) (sl []model.Slot, err error) {
	query := `SELECT 
			ID,
  		UF_TIME,
			UF_DATE
		FROM 
			z_slots
		WHERE
			UF_DATE > CONCAT(CURDATE(), ' 00:00:00')
		`
	if !isKNO {
		query += `
			AND UF_IS_BUSY IS NULL OR UF_IS_BUSY = 0
		`
	}
	data, err := s.db.QueryContext(ctx, query)
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
		sl = append(sl, p)
	}
	return sl, nil
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
