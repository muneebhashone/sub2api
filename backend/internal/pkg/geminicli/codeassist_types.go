package geminicli

// LoadCodeAssistRequest matches done-hub's internal Code Assist call.
type LoadCodeAssistRequest struct {
	Metadata LoadCodeAssistMetadata `json:"metadata"`
placeholder

type LoadCodeAssistMetadata struct {
	IDEType    string `json:"ideType"`
	Platform   string `json:"platform"`
	PluginType string `json:"pluginType"`
placeholder

type LoadCodeAssistResponse struct {
	CurrentTier             string        `json:"currentTier,omitempty"`
	CloudAICompanionProject string        `json:"cloudaicompanionProject,omitempty"`
	AllowedTiers            []AllowedTier `json:"allowedTiers,omitempty"`
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
