package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const AccountSchedulingThresholdReasonSource = "account_scheduling_threshold"

const (
	defaultTempUnschedReasonErrorMessage          = "temporary scheduling block reason unavailable"
	defaultAccountSchedulingThresholdErrorMessage = "account scheduling threshold reached"
)

type tempUnschedReasonPayload struct {
	Source           string  `json:"source,omitempty"`
	Platform         string  `json:"platform,omitempty"`
	Window           string  `json:"window,omitempty"`
	Scope            string  `json:"scope,omitempty"`
	ThresholdPercent int     `json:"threshold_percent,omitempty"`
	UsedPercent      float64 `json:"used_percent,omitempty"`
	UntilUnix        int64   `json:"until_unix,omitempty"`
	TriggeredAtUnix  int64   `json:"triggered_at_unix,omitempty"`
	ErrorMessage     string  `json:"error_message"`
placeholder

type AccountSchedulingThresholdReasonInput struct {
	Platform         string
	Window           string
	Scope            string
	ThresholdPercent int
	UsedPercent      float64
	Until            time.Time
	Now              time.Time
placeholder

func BuildTempUnschedReasonPayload(source string, errorMessage string) string {
	payload := tempUnschedReasonPayload{
		Source:       strings.TrimSpace(source),
		ErrorMessage: normalizeTempUnschedReasonErrorMessage(errorMessage, defaultTempUnschedReasonErrorMessage),
placeholder

	raw, err := json.Marshal(payload)
	if err != nil {
		return payload.ErrorMessage
placeholder
	return string(raw)
placeholder

func BuildAccountSchedulingThresholdReason(errorMessage string) string {
	return BuildTempUnschedReasonPayload(
		AccountSchedulingThresholdReasonSource,
		normalizeTempUnschedReasonErrorMessage(errorMessage, defaultAccountSchedulingThresholdErrorMessage),
	)
placeholder

func BuildDetailedAccountSchedulingThresholdReason(input AccountSchedulingThresholdReasonInput) string {
	triggeredAt := input.Now
	if triggeredAt.IsZero() {
		triggeredAt = time.Now().UTC()
placeholder
	payload := tempUnschedReasonPayload{
		Source:           AccountSchedulingThresholdReasonSource,
		Platform:         strings.TrimSpace(input.Platform),
		Window:           strings.TrimSpace(input.Window),
		Scope:            strings.TrimSpace(input.Scope),
		ThresholdPercent: input.ThresholdPercent,
		UsedPercent:      input.UsedPercent,
		TriggeredAtUnix:  triggeredAt.Unix(),
		ErrorMessage:     buildAccountSchedulingThresholdErrorMessage(input),
placeholder
	if !input.Until.IsZero() {
		payload.UntilUnix = input.Until.UTC().Unix()
placeholder

	raw, err := json.Marshal(payload)
	if err != nil {
		return payload.ErrorMessage
placeholder
	return string(raw)
placeholder

func IsAccountSchedulingThresholdReason(rawReason string) bool {
	payload, ok := parseTempUnschedReasonPayload(rawReason)
	if !ok {
		return false
placeholder
	return payload.Source == AccountSchedulingThresholdReasonSource
placeholder

func parseTempUnschedReasonPayload(rawReason string) (tempUnschedReasonPayload, bool) {
	rawReason = strings.TrimSpace(rawReason)
	if rawReason == "" {
		return tempUnschedReasonPayload{placeholder, false
placeholder

	var payload tempUnschedReasonPayload
	if err := json.Unmarshal([]byte(rawReason), &payload); err != nil {
		return tempUnschedReasonPayload{placeholder, false
placeholder
	payload.Source = strings.TrimSpace(payload.Source)
	payload.ErrorMessage = strings.TrimSpace(payload.ErrorMessage)
	return payload, true
placeholder

func normalizeTempUnschedReasonErrorMessage(errorMessage string, fallback string) string {
	errorMessage = strings.TrimSpace(errorMessage)
	if errorMessage != "" {
		return errorMessage
placeholder

	fallback = strings.TrimSpace(fallback)
	if fallback != "" {
		return fallback
placeholder
	return defaultTempUnschedReasonErrorMessage
placeholder

func buildAccountSchedulingThresholdErrorMessage(input AccountSchedulingThresholdReasonInput) string {
	platform := strings.TrimSpace(input.Platform)
	if platform == "" {
		platform = "account"
placeholder

	target := strings.TrimSpace(input.Window)
	if scope := strings.TrimSpace(input.Scope); scope != "" {
		if target == "" {
			target = scope
	placeholder else {
			target = target + "/" + scope
	placeholder
placeholder
	if target == "" {
		target = "usage window"
placeholder

	threshold := input.ThresholdPercent
	if threshold <= 0 {
		threshold = 100
placeholder

	untilText := "the window reset"
	if !input.Until.IsZero() {
		untilText = input.Until.UTC().Format(time.RFC3339)
placeholder

	return fmt.Sprintf(
		"%s scheduling threshold reached for %s: %.1f%% used >= %d%%; paused until %s",
		platform,
		target,
		input.UsedPercent,
		threshold,
		untilText,
	)
placeholder

func tempUnschedStateFromStoredReason(rawReason string, fallbackUntilUnix int64) *TempUnschedState {
	state := &TempUnschedState{
		UntilUnix: fallbackUntilUnix,
		RuleIndex: -1,
placeholder

	rawReason = strings.TrimSpace(rawReason)
	if rawReason == "" {
		state.ErrorMessage = defaultTempUnschedReasonErrorMessage
		return state
placeholder

	parsed := TempUnschedState{RuleIndex: -1placeholder
	if err := json.Unmarshal([]byte(rawReason), &parsed); err == nil {
		if fallbackUntilUnix > parsed.UntilUnix {
			parsed.UntilUnix = fallbackUntilUnix
	placeholder
		if strings.TrimSpace(parsed.ErrorMessage) == "" {
			if IsAccountSchedulingThresholdReason(rawReason) {
				parsed.ErrorMessage = defaultAccountSchedulingThresholdErrorMessage
		placeholder else {
				parsed.ErrorMessage = rawReason
		placeholder
	placeholder
		return &parsed
placeholder

	state.ErrorMessage = rawReason
	return state
placeholder
