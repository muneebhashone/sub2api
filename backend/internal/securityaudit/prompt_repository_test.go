package securityaudit

import "testing"

func TestShouldStorePromptAuditEvent(t *testing.T) {
	tests := []struct {
		name            string
		storePassEvents bool
		decision        EventDecision
		want            bool
placeholder{
		{name: "pass disabled", storePassEvents: false, decision: EventPass, want: falseplaceholder,
		{name: "flag disabled", storePassEvents: false, decision: EventFlag, want: trueplaceholder,
		{name: "critical disabled", storePassEvents: false, decision: EventCritical, want: trueplaceholder,
		{name: "pass enabled", storePassEvents: true, decision: EventPass, want: trueplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldStorePromptAuditEvent(tt.decision, tt.storePassEvents); got != tt.want {
				t.Fatalf("shouldStorePromptAuditEvent(%q, %t) = %t, want %t", tt.decision, tt.storePassEvents, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder
