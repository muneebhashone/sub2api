// Package gemini provides minimal fallback model metadata for Gemini native endpoints.
// It is used when upstream model listing is unavailable (e.g. OAuth token missing AI Studio scopes).
package gemini

type Model struct {
	Name                       string   `json:"name"`
	DisplayName                string   `json:"displayName,omitempty"`
	Description                string   `json:"description,omitempty"`
	SupportedGenerationMethods []string `json:"supportedGenerationMethods,omitempty"`
placeholder

type ModelsListResponse struct {
	Models []Model `json:"models"`
placeholder

func DefaultModels() []Model {
	methods := []string{"generateContent", "streamGenerateContent"placeholder
	return []Model{
		{Name: "models/gemini-2.0-flash", SupportedGenerationMethods: methodsplaceholder,
		{Name: "models/gemini-2.5-flash", SupportedGenerationMethods: methodsplaceholder,
		{Name: "models/gemini-2.5-pro", SupportedGenerationMethods: methodsplaceholder,
		{Name: "models/gemini-3-flash-preview", SupportedGenerationMethods: methodsplaceholder,
		{Name: "models/gemini-3-pro-preview", SupportedGenerationMethods: methodsplaceholder,
		{Name: "models/gemini-3.1-pro-preview", SupportedGenerationMethods: methodsplaceholder,
placeholder
placeholder

func FallbackModelsList() ModelsListResponse {
	return ModelsListResponse{Models: DefaultModels()placeholder
placeholder

func FallbackModel(model string) Model {
	methods := []string{"generateContent", "streamGenerateContent"placeholder
	if model == "" {
		return Model{Name: "models/unknown", SupportedGenerationMethods: methodsplaceholder
placeholder
	if len(model) >= 7 && model[:7] == "models/" {
		return Model{Name: model, SupportedGenerationMethods: methodsplaceholder
placeholder
	return Model{Name: "models/" + model, SupportedGenerationMethods: methodsplaceholder
placeholder
