package service

import "context"

type openAIForwardModelContextKey struct{placeholder

type openAIForwardModel struct {
	model                  string
	useCompactModelMapping bool
placeholder

// WithOpenAIForwardModel records the model present in the forwarded request
// body after channel mapping and whether the legacy compact-only model mapping
// applies. Channel restriction checks then follow the same model chain used by
// Forward.
func WithOpenAIForwardModel(ctx context.Context, forwardModel string, useCompactModelMapping bool) context.Context {
	if ctx == nil {
		ctx = context.Background()
placeholder
	return context.WithValue(ctx, openAIForwardModelContextKey{placeholder, openAIForwardModel{
		model:                  forwardModel,
		useCompactModelMapping: useCompactModelMapping,
placeholder)
placeholder

func openAIForwardModelFromContext(ctx context.Context) (openAIForwardModel, bool) {
	if ctx == nil {
		return openAIForwardModel{placeholder, false
placeholder
	forwardModel, ok := ctx.Value(openAIForwardModelContextKey{placeholder).(openAIForwardModel)
	return forwardModel, ok
placeholder
