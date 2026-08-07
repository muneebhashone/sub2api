//go:build unit

package handler

import (
	"testing"

	coderws "github.com/coder/websocket"
)

func TestIsExpectedGrokRealtimeClose(t *testing.T) {
	for _, status := range []coderws.StatusCode{
		coderws.StatusNormalClosure,
		coderws.StatusGoingAway,
		coderws.StatusNoStatusRcvd,
		coderws.StatusAbnormalClosure,
placeholder {
		if !isExpectedGrokRealtimeClose(coderws.CloseError{Code: statusplaceholder) {
			t.Fatalf("status %v should be treated as an expected session close", status)
	placeholder
placeholder
	if isExpectedGrokRealtimeClose(coderws.CloseError{Code: coderws.StatusPolicyViolationplaceholder) {
		t.Fatal("policy violations must not be treated as billable normal closes")
placeholder
placeholder
