//go:build unit

package service

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitBackupFile_ReassemblesExactBytes(t *testing.T) {
	src := writeBackupArchiveFixture(t, []byte("0123456789abcdefg"))
	parts, err := splitBackupFile(src, 5)
placeholder
	require.Len(t, parts, 4)

	var got bytes.Buffer
	for i, part := range parts {
		require.Equal(t, i+1, part.Index)
		require.LessOrEqual(t, part.SizeBytes, int64(5))
		data, readErr := os.ReadFile(part.Path)
		require.NoError(t, readErr)
		require.Equal(t, fmt.Sprintf("%x", sha256.Sum256(data)), part.SHA256)
		got.Write(data)
placeholder
	require.Equal(t, []byte("0123456789abcdefg"), got.Bytes())
placeholder

func TestSplitBackupFile_RejectsInvalidInput(t *testing.T) {
	src := writeBackupArchiveFixture(t, []byte("data"))

	_, err := splitBackupFile(src, 0)
placeholder

	empty := writeBackupArchiveFixture(t, nil)
	_, err = splitBackupFile(empty, 5)
placeholder

	_, err = splitBackupFile(filepathForMissingBackupArchive(t), 5)
placeholder
placeholder

func writeBackupArchiveFixture(t *testing.T, content []byte) string {
placeholder
	path := filepathForBackupArchive(t)
	require.NoError(t, os.WriteFile(path, content, 0o600))
	return path
placeholder

func filepathForBackupArchive(t *testing.T) string {
placeholder
	return t.TempDir() + "/archive.gz"
placeholder

func filepathForMissingBackupArchive(t *testing.T) string {
placeholder
	return t.TempDir() + "/missing.gz"
placeholder
