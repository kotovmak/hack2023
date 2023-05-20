package store

import (
	"context"
	"hack2023/internal/app/model"
)

func (s *Store) GetTypeList(ctx context.Context, inn string) (map[string]model.TypeList, error) {
	ps := make(map[string]model.TypeList)
	return ps, nil
}
