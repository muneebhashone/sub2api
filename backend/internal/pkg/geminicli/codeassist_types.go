package geminicli

import (
	"bytes"
	"encoding/json"
)

// LoadCodeAssistRequest matches done-hub's internal Code Assist call.
type LoadCodeAssistRequest struct {
	Metadata LoadCodeAssistMetadata `json:"metadata"`
placeholder

type LoadCodeAssistMetadata struct {
	IDEType    string `json:"ideType"`
	Platform   string `json:"platform"`
	PluginType string `json:"pluginType"`
placeholder

type TierInfo struct {
	ID string `json:"id"`
placeholder

// UnmarshalJSON supports both legacy string tiers and object tiers.
func (t *TierInfo) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		return nil
placeholder
	if data[0] == '"' {
		var id string
		if err := json.Unmarshal(data, &id); err != nil {
			return err
	placeholder
		t.ID = id
		return nil
placeholder
	type alias TierInfo
	var decoded alias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
placeholder
	*t = TierInfo(decoded)
	return nil
placeholder

type LoadCodeAssistResponse struct {
	CurrentTier             *TierInfo     `json:"currentTier,omitempty"`
	PaidTier                *TierInfo     `json:"paidTier,omitempty"`
	CloudAICompanionProject string        `json:"cloudaicompanionProject,omitempty"`
	AllowedTiers            []AllowedTier `json:"allowedTiers,omitempty"`
placeholder

// GetTier extracts tier ID, prioritizing paidTier over currentTier
func (r *LoadCodeAssistResponse) GetTier() string {
	if r.PaidTier != nil && r.PaidTier.ID != "" {
		return r.PaidTier.ID
placeholder
	if r.CurrentTier != nil {
		return r.CurrentTier.ID
placeholder
	return ""
placeholder

type AllowedTier struct {
	ID        string `json:"id"`
	IsDefault bool   `json:"isDefault,omitempty"`
placeholder

type OnboardUserRequest struct {
	TierID   string                 `json:"tierId"`
	Metadata LoadCodeAssistMetadata `json:"metadata"`
placeholder

type OnboardUserResponse struct {
	Done     bool                   `json:"done"`
	Response *OnboardUserResultData `json:"response,omitempty"`
	Name     string                 `json:"name,omitempty"`
placeholder

type OnboardUserResultData struct {
	CloudAICompanionProject any `json:"cloudaicompanionProject,omitempty"`
placeholder
