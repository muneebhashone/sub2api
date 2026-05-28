//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAdminService_GetUserIncludeDeleted(t *testing.T) {
	ts := time.Date(2026, 5, 28, 0, 0, 0, 0, time.UTC)
	repo := &userRepoStub{user: &User{ID: 7, Email: "del@test.com", DeletedAt: &tsplaceholderplaceholder
	svc := &adminServiceImpl{userRepo: repoplaceholder

	got, err := svc.GetUserIncludeDeleted(context.Background(), 7)
placeholder
	require.Equal(t, int64(7), got.ID)
	require.NotNil(t, got.DeletedAt)
placeholder
