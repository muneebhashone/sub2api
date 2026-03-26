package service

import "testing"

func TestParseDebugEnvBool(t *testing.T) {
	t.Run("empty is false", func(t *testing.T) {
		if parseDebugEnvBool("") {
			t.Fatalf("expected false for empty string")
	placeholder
placeholder)

	t.Run("true-like values", func(t *testing.T) {
		for _, value := range []string{"1", "true", "TRUE", "yes", "on"placeholder {
			t.Run(value, func(t *testing.T) {
				if !parseDebugEnvBool(value) {
					t.Fatalf("expected true for %q", value)
			placeholder
		placeholder)
	placeholder
placeholder)

	t.Run("false-like values", func(t *testing.T) {
		for _, value := range []string{"0", "false", "off", "debug"placeholder {
			t.Run(value, func(t *testing.T) {
				if parseDebugEnvBool(value) {
					t.Fatalf("expected false for %q", value)
			placeholder
		placeholder)
	placeholder
placeholder)
placeholder
