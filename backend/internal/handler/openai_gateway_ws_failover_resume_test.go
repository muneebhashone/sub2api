package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenAIWSNextAttemptMessageUsesCurrentTurnPayload(t *testing.T) {
	firstMessage := []byte(`{"type":"response.create","input":"first"placeholder`)
	currentTurn := []byte(`{"type":"response.create","input":"turn-281"placeholder`)

	next, ok := openAIWSNextAttemptMessage(firstMessage, currentTurn, true)

	require.True(t, ok)
	require.Equal(t, currentTurn, next)
	next[0] = 'x'
	require.Equal(t, byte('{'), currentTurn[0], "retry payload must be cloned")
placeholder

func TestOpenAIWSNextAttemptMessageRejectsMissingCurrentTurnPayload(t *testing.T) {
	next, ok := openAIWSNextAttemptMessage([]byte(`{"type":"response.create"placeholder`), nil, true)

	require.False(t, ok)
	require.Nil(t, next)
placeholder

func TestOpenAIWSNextAttemptMessageKeepsInitialMessageForFirstTurnFailover(t *testing.T) {
	firstMessage := []byte(`{"type":"response.create","input":"first"placeholder`)

	next, ok := openAIWSNextAttemptMessage(firstMessage, nil, false)

	require.True(t, ok)
	require.Equal(t, firstMessage, next)
placeholder
