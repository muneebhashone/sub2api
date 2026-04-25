package service

import "testing"

func TestAccountGetOpenAICompactMode(t *testing.T) {
	tests := []struct {
		name    string
		account *Account
		want    string
placeholder{
		{
			name: "nil account defaults to auto",
			want: OpenAICompactModeAuto,
	placeholder,
		{
			name: "non openai account defaults to auto",
			account: &Account{
				Platform: PlatformAnthropic,
				Extra:    map[string]any{"openai_compact_mode": OpenAICompactModeForceOnplaceholder,
		placeholder,
			want: OpenAICompactModeAuto,
	placeholder,
		{
			name: "missing extra defaults to auto",
			account: &Account{
				Platform: PlatformOpenAI,
		placeholder,
			want: OpenAICompactModeAuto,
	placeholder,
		{
			name: "invalid mode falls back to auto",
			account: &Account{
				Platform: PlatformOpenAI,
				Extra:    map[string]any{"openai_compact_mode": "  invalid  "placeholder,
		placeholder,
			want: OpenAICompactModeAuto,
	placeholder,
		{
			name: "force on is normalized",
			account: &Account{
				Platform: PlatformOpenAI,
				Extra:    map[string]any{"openai_compact_mode": " FORCE_ON "placeholder,
		placeholder,
			want: OpenAICompactModeForceOn,
	placeholder,
		{
			name: "force off is normalized",
			account: &Account{
				Platform: PlatformOpenAI,
				Extra:    map[string]any{"openai_compact_mode": "force_off"placeholder,
		placeholder,
			want: OpenAICompactModeForceOff,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.account.GetOpenAICompactMode(); got != tt.want {
				t.Fatalf("GetOpenAICompactMode() = %q, want %q", got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestAccountOpenAICompactSupportKnown(t *testing.T) {
	tests := []struct {
		name          string
		account       *Account
		wantSupported bool
		wantKnown     bool
placeholder{
		{
			name:          "nil account is unknown",
			wantSupported: false,
			wantKnown:     false,
	placeholder,
		{
			name: "non openai account is unknown",
			account: &Account{
				Platform: PlatformAnthropic,
				Extra:    map[string]any{"openai_compact_supported": trueplaceholder,
		placeholder,
			wantSupported: false,
			wantKnown:     false,
	placeholder,
		{
			name: "force on overrides probe state",
			account: &Account{
				Platform: PlatformOpenAI,
				Extra: map[string]any{
					"openai_compact_mode":      OpenAICompactModeForceOn,
					"openai_compact_supported": false,
			placeholder,
		placeholder,
			wantSupported: true,
			wantKnown:     true,
	placeholder,
		{
			name: "force off overrides probe state",
			account: &Account{
				Platform: PlatformOpenAI,
				Extra: map[string]any{
					"openai_compact_mode":      OpenAICompactModeForceOff,
					"openai_compact_supported": true,
			placeholder,
		placeholder,
			wantSupported: false,
			wantKnown:     true,
	placeholder,
		{
			name: "auto true is known supported",
			account: &Account{
				Platform: PlatformOpenAI,
				Extra:    map[string]any{"openai_compact_supported": trueplaceholder,
		placeholder,
			wantSupported: true,
			wantKnown:     true,
	placeholder,
		{
			name: "auto false is known unsupported",
			account: &Account{
				Platform: PlatformOpenAI,
				Extra:    map[string]any{"openai_compact_supported": falseplaceholder,
		placeholder,
			wantSupported: false,
			wantKnown:     true,
	placeholder,
		{
			name: "auto without probe state remains unknown",
			account: &Account{
				Platform: PlatformOpenAI,
				Extra:    map[string]any{placeholder,
		placeholder,
			wantSupported: false,
			wantKnown:     false,
	placeholder,
		{
			name: "invalid probe field remains unknown",
			account: &Account{
				Platform: PlatformOpenAI,
				Extra:    map[string]any{"openai_compact_supported": "true"placeholder,
		placeholder,
			wantSupported: false,
			wantKnown:     false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSupported, gotKnown := tt.account.OpenAICompactSupportKnown()
			if gotSupported != tt.wantSupported || gotKnown != tt.wantKnown {
				t.Fatalf("OpenAICompactSupportKnown() = (%v, %v), want (%v, %v)", gotSupported, gotKnown, tt.wantSupported, tt.wantKnown)
		placeholder
	placeholder)
placeholder
placeholder

func TestAccountAllowsOpenAICompact(t *testing.T) {
	tests := []struct {
		name    string
		account *Account
		want    bool
placeholder{
		{
			name: "nil account does not allow compact",
			want: false,
	placeholder,
		{
			name: "non openai account does not allow compact",
			account: &Account{
				Platform: PlatformAnthropic,
		placeholder,
			want: false,
	placeholder,
		{
			name: "unknown openai account remains allowed",
			account: &Account{
				Platform: PlatformOpenAI,
				Extra:    map[string]any{placeholder,
		placeholder,
			want: true,
	placeholder,
		{
			name: "supported openai account is allowed",
			account: &Account{
				Platform: PlatformOpenAI,
				Extra:    map[string]any{"openai_compact_supported": trueplaceholder,
		placeholder,
			want: true,
	placeholder,
		{
			name: "unsupported openai account is rejected",
			account: &Account{
				Platform: PlatformOpenAI,
				Extra:    map[string]any{"openai_compact_supported": falseplaceholder,
		placeholder,
			want: false,
	placeholder,
		{
			name: "force on is allowed",
			account: &Account{
				Platform: PlatformOpenAI,
				Extra:    map[string]any{"openai_compact_mode": OpenAICompactModeForceOnplaceholder,
		placeholder,
			want: true,
	placeholder,
		{
			name: "force off is rejected",
			account: &Account{
				Platform: PlatformOpenAI,
				Extra:    map[string]any{"openai_compact_mode": OpenAICompactModeForceOffplaceholder,
		placeholder,
			want: false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.account.AllowsOpenAICompact(); got != tt.want {
				t.Fatalf("AllowsOpenAICompact() = %v, want %v", got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestAccountGetCompactModelMapping(t *testing.T) {
	tests := []struct {
		name    string
		account *Account
		want    map[string]string
placeholder{
		{
			name: "nil account returns nil",
			want: nil,
	placeholder,
		{
			name: "missing credentials returns nil",
			account: &Account{
				Platform: PlatformOpenAI,
		placeholder,
			want: nil,
	placeholder,
		{
			name: "map any is converted",
			account: &Account{
		placeholder
					"compact_model_mapping": map[string]any{
						"gpt-5.4": "gpt-5.4-openai-compact",
						"invalid": 1,
				placeholder,
			placeholder,
		placeholder,
			want: map[string]string{
				"gpt-5.4": "gpt-5.4-openai-compact",
		placeholder,
	placeholder,
		{
			name: "map string string is copied",
			account: &Account{
		placeholder
					"compact_model_mapping": map[string]string{
						"gpt-*": "compact-*",
				placeholder,
			placeholder,
		placeholder,
			want: map[string]string{
				"gpt-*": "compact-*",
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.account.GetCompactModelMapping()
			if !equalStringMap(got, tt.want) {
				t.Fatalf("GetCompactModelMapping() = %#v, want %#v", got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestAccountResolveCompactMappedModel(t *testing.T) {
	tests := []struct {
		name           string
		credentials    map[string]any
		requestedModel string
		expectedModel  string
		expectedMatch  bool
placeholder{
		{
			name:           "no compact mapping reports unmatched",
			credentials:    nil,
			requestedModel: "gpt-5.4",
			expectedModel:  "gpt-5.4",
			expectedMatch:  false,
	placeholder,
		{
			name: "exact compact mapping matches",
			credentials: map[string]any{
				"compact_model_mapping": map[string]any{
					"gpt-5.4": "gpt-5.4-openai-compact",
			placeholder,
		placeholder,
			requestedModel: "gpt-5.4",
			expectedModel:  "gpt-5.4-openai-compact",
			expectedMatch:  true,
	placeholder,
		{
			name: "exact passthrough counts as match",
			credentials: map[string]any{
				"compact_model_mapping": map[string]any{
					"gpt-5.4": "gpt-5.4",
			placeholder,
		placeholder,
			requestedModel: "gpt-5.4",
			expectedModel:  "gpt-5.4",
			expectedMatch:  true,
	placeholder,
		{
			name: "longest wildcard wins",
			credentials: map[string]any{
				"compact_model_mapping": map[string]any{
					"gpt-*":         "fallback-compact",
					"gpt-5.4*":      "gpt-5.4-openai-compact",
					"gpt-5.4-mini*": "gpt-5.4-mini-openai-compact",
			placeholder,
		placeholder,
			requestedModel: "gpt-5.4-mini",
			expectedModel:  "gpt-5.4-mini-openai-compact",
			expectedMatch:  true,
	placeholder,
		{
			name: "missing compact mapping reports unmatched",
			credentials: map[string]any{
				"compact_model_mapping": map[string]any{
					"gpt-5.3": "gpt-5.3-openai-compact",
			placeholder,
		placeholder,
			requestedModel: "gpt-5.4",
			expectedModel:  "gpt-5.4",
			expectedMatch:  false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &Account{
				Platform:    PlatformOpenAI,
				Credentials: tt.credentials,
		placeholder
			gotModel, gotMatch := account.ResolveCompactMappedModel(tt.requestedModel)
			if gotModel != tt.expectedModel || gotMatch != tt.expectedMatch {
				t.Fatalf("ResolveCompactMappedModel(%q) = (%q, %v), want (%q, %v)", tt.requestedModel, gotModel, gotMatch, tt.expectedModel, tt.expectedMatch)
		placeholder
	placeholder)
placeholder
placeholder

func equalStringMap(left, right map[string]string) bool {
	if len(left) != len(right) {
		return false
placeholder
	for key, want := range right {
		if got, ok := left[key]; !ok || got != want {
			return false
	placeholder
placeholder
	return true
placeholder
