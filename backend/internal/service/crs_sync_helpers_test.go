package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildSelectedSet(t *testing.T) {
	tests := []struct {
		name     string
		ids      []string
		wantNil  bool
		wantSize int
placeholder{
		{
			name:    "nil input returns nil (backward compatible: create all)",
			ids:     nil,
			wantNil: true,
	placeholder,
		{
			name:     "empty slice returns empty map (create none)",
			ids:      []string{placeholder,
			wantNil:  false,
			wantSize: 0,
	placeholder,
		{
			name:     "single ID",
			ids:      []string{"abc-123"placeholder,
			wantNil:  false,
			wantSize: 1,
	placeholder,
		{
			name:     "multiple IDs",
			ids:      []string{"a", "b", "c"placeholder,
			wantNil:  false,
			wantSize: 3,
	placeholder,
		{
			name:     "duplicate IDs are deduplicated",
			ids:      []string{"a", "a", "b"placeholder,
			wantNil:  false,
			wantSize: 2,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildSelectedSet(tt.ids)
			if tt.wantNil {
				if got != nil {
					t.Errorf("buildSelectedSet(%v) = %v, want nil", tt.ids, got)
			placeholder
				return
		placeholder
			if got == nil {
				t.Fatalf("buildSelectedSet(%v) = nil, want non-nil map", tt.ids)
		placeholder
			if len(got) != tt.wantSize {
				t.Errorf("buildSelectedSet(%v) has %d entries, want %d", tt.ids, len(got), tt.wantSize)
		placeholder
			// Verify all unique IDs are present
			for _, id := range tt.ids {
				if _, ok := got[id]; !ok {
					t.Errorf("buildSelectedSet(%v) missing key %q", tt.ids, id)
			placeholder
		placeholder
	placeholder)
placeholder
placeholder

func TestShouldCreateAccount(t *testing.T) {
	tests := []struct {
		name        string
		crsID       string
		selectedSet map[string]struct{placeholder
		want        bool
placeholder{
		{
			name:        "nil set allows all (backward compatible)",
			crsID:       "any-id",
			selectedSet: nil,
			want:        true,
	placeholder,
		{
			name:        "empty set blocks all",
			crsID:       "any-id",
			selectedSet: map[string]struct{placeholder{placeholder,
			want:        false,
	placeholder,
		{
			name:        "ID in set is allowed",
			crsID:       "abc-123",
			selectedSet: map[string]struct{placeholder{"abc-123": {placeholder, "def-456": {placeholderplaceholder,
			want:        true,
	placeholder,
		{
			name:        "ID not in set is blocked",
			crsID:       "xyz-789",
			selectedSet: map[string]struct{placeholder{"abc-123": {placeholder, "def-456": {placeholderplaceholder,
			want:        false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldCreateAccount(tt.crsID, tt.selectedSet)
			if got != tt.want {
				t.Errorf("shouldCreateAccount(%q, %v) = %v, want %v",
					tt.crsID, tt.selectedSet, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestReconcileCRSUpstreamBillingProbeExtra(t *testing.T) {
	remote := map[string]any{
		"crs_account_id":                    "remote-1",
		UpstreamBillingProbeEnabledExtraKey: true,
		UpstreamBillingProbeExtraKey:        map[string]any{"status": "remote"placeholder,
placeholder

	t.Run("create drops remote managed fields", func(t *testing.T) {
		extra := mergeMap(nil, remote)
		reconcileCRSUpstreamBillingProbeExtra(nil, PlatformOpenAI, AccountTypeAPIKey, map[string]any{"api_key": "new"placeholder, extra)
		require.NotContains(t, extra, UpstreamBillingProbeEnabledExtraKey)
		require.NotContains(t, extra, UpstreamBillingProbeExtraKey)
placeholder)

	existing := &Account{
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
placeholder"api_key": "local", "base_url": "http://127.0.0.1:8080"placeholder,
		Extra: map[string]any{
			UpstreamBillingProbeEnabledExtraKey: false,
			UpstreamBillingProbeExtraKey:        map[string]any{"status": "local"placeholder,
	placeholder,
placeholder

	t.Run("same identity keeps local state", func(t *testing.T) {
		extra := mergeMap(existing.Extra, remote)
		reconcileCRSUpstreamBillingProbeExtra(existing, existing.Platform, existing.Type, mergeMap(existing.Credentials, nil), extra)
		require.Equal(t, false, extra[UpstreamBillingProbeEnabledExtraKey])
		require.Equal(t, map[string]any{"status": "local"placeholder, extra[UpstreamBillingProbeExtraKey])
placeholder)

	t.Run("identity change keeps enabled and clears snapshot", func(t *testing.T) {
		extra := mergeMap(existing.Extra, remote)
		reconcileCRSUpstreamBillingProbeExtra(existing, PlatformOpenAI, AccountTypeAPIKey, map[string]any{"api_key": "changed"placeholder, extra)
		require.Equal(t, false, extra[UpstreamBillingProbeEnabledExtraKey])
		require.NotContains(t, extra, UpstreamBillingProbeExtraKey)
placeholder)

	for _, target := range []struct {
		name     string
		platform string
		typeName string
placeholder{
		{name: "anthropic oauth", platform: PlatformAnthropic, typeName: AccountTypeOAuthplaceholder,
		{name: "anthropic api key", platform: PlatformAnthropic, typeName: AccountTypeAPIKeyplaceholder,
		{name: "openai oauth", platform: PlatformOpenAI, typeName: AccountTypeOAuthplaceholder,
		{name: "gemini oauth", platform: PlatformGemini, typeName: AccountTypeOAuthplaceholder,
		{name: "gemini api key", platform: PlatformGemini, typeName: AccountTypeAPIKeyplaceholder,
placeholder {
		t.Run(target.name+" removes inapplicable state", func(t *testing.T) {
			extra := mergeMap(existing.Extra, remote)
			reconcileCRSUpstreamBillingProbeExtra(existing, target.platform, target.typeName, existing.Credentials, extra)
			require.NotContains(t, extra, UpstreamBillingProbeEnabledExtraKey)
			require.NotContains(t, extra, UpstreamBillingProbeExtraKey)
	placeholder)
placeholder
placeholder
