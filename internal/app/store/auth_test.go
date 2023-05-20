package store

import (
	"context"
	"hack2023/internal/app/config"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthStore_GetUserByLogin(t *testing.T) {
	login := "user"
	conf := config.Get()

	store, err := New(conf)
	if err != nil {
		log.Fatal(err)
	}

	user, err := store.GetUserByLogin(context.Background(), login)
	assert.NoError(t, err)
	assert.Equal(t, user.Login, login)
}

func TestAuthStore_RefreshToken(t *testing.T) {
	userId := 1
	token := "some_token"
	conf := config.Get()

	store, err := New(conf)
	if err != nil {
		log.Fatal(err)
	}

	err = store.SaveRefreshToken(context.Background(), userId, token)
	assert.NoError(t, err)

	user, err := store.GetUserByRefreshToken(context.Background(), token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, userId)

	err = store.DeleteRefreshToken(context.Background(), userId)
	assert.NoError(t, err)

	user, err = store.GetUserByRefreshToken(context.Background(), token)
	assert.Error(t, err)
}
