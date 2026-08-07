//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type missingThresholdSnapshotCleanerRepo struct {
	AccountRepository
placeholder

func TestClearAccountSchedulingThresholdSnapshots_RequiresRepositorySupport(t *testing.T) {
	err := clearAccountSchedulingThresholdSnapshots(context.Background(), missingThresholdSnapshotCleanerRepo{placeholder, 1)

placeholder
	require.Contains(t, err.Error(), "does not support account scheduling threshold snapshot cleanup")
placeholder
