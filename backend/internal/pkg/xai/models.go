package xai

// Model describes an xAI model in OpenAI-compatible /models shape.
type Model struct {
	ID          string `json:"id"`
	Object      string `json:"object"`
	Created     int64  `json:"created,omitempty"`
	OwnedBy     string `json:"owned_by"`
	DisplayName string `json:"display_name,omitempty"`
placeholder

var defaultModels = []Model{
	{ID: "grok-4.3", Object: "model", OwnedBy: "xai", DisplayName: "Grok 4.3"placeholder,
	{ID: "grok-build-0.1", Object: "model", OwnedBy: "xai", DisplayName: "Grok Build 0.1"placeholder,
	{ID: "grok-4.20-0309-reasoning", Object: "model", OwnedBy: "xai", DisplayName: "Grok 4.20 Reasoning"placeholder,
	{ID: "grok-4.20-0309-non-reasoning", Object: "model", OwnedBy: "xai", DisplayName: "Grok 4.20 Non Reasoning"placeholder,
	{ID: "grok-4.20-multi-agent-0309", Object: "model", OwnedBy: "xai", DisplayName: "Grok 4.20 Multi Agent"placeholder,
	{ID: "grok-imagine", Object: "model", OwnedBy: "xai", DisplayName: "Grok Imagine"placeholder,
	{ID: "grok-imagine-edit", Object: "model", OwnedBy: "xai", DisplayName: "Grok Imagine Edit"placeholder,
	{ID: "grok-imagine-video-1.5", Object: "model", OwnedBy: "xai", DisplayName: "Grok Imagine Video 1.5"placeholder,
placeholder

func DefaultModels() []Model {
	out := make([]Model, len(defaultModels))
	copy(out, defaultModels)
	return out
placeholder

func DefaultModelIDs() []string {
	models := DefaultModels()
	ids := make([]string, 0, len(models))
	for _, model := range models {
		ids = append(ids, model.ID)
placeholder
	return ids
placeholder

func DefaultModelMapping() map[string]string {
	mapping := make(map[string]string, len(defaultModels)+3)
	for _, model := range defaultModels {
		mapping[model.ID] = model.ID
placeholder
	mapping["grok"] = "grok-4.3"
	mapping["grok-latest"] = "grok-4.3"
	mapping["grok-build"] = "grok-build-0.1"
	mapping["grok-4.20-reasoning"] = "grok-4.20-0309-reasoning"
	mapping["grok-4.20-non-reasoning"] = "grok-4.20-0309-non-reasoning"
	return mapping
placeholder
