package service

import "context"

// HTTPUpstreamProfile marks HTTP upstream requests that need provider-specific
// transport policy.
type HTTPUpstreamProfile string

const (
	HTTPUpstreamProfileDefault HTTPUpstreamProfile = ""
	HTTPUpstreamProfileOpenAI  HTTPUpstreamProfile = "openai"
	HTTPUpstreamProfileGrok    HTTPUpstreamProfile = "grok"
)

type httpUpstreamProfileContextKey struct{placeholder
type httpUpstreamDisableRedirectsContextKey struct{placeholder

// WithHTTPUpstreamProfile injects an upstream transport profile into ctx.
func WithHTTPUpstreamProfile(ctx context.Context, profile HTTPUpstreamProfile) context.Context {
	if ctx == nil {
		ctx = context.Background()
placeholder
	if profile == HTTPUpstreamProfileDefault {
		return ctx
placeholder
	return context.WithValue(ctx, httpUpstreamProfileContextKey{placeholder, profile)
placeholder

// HTTPUpstreamProfileFromContext resolves the upstream transport profile from ctx.
func HTTPUpstreamProfileFromContext(ctx context.Context) HTTPUpstreamProfile {
	if ctx == nil {
		return HTTPUpstreamProfileDefault
placeholder
	profile, ok := ctx.Value(httpUpstreamProfileContextKey{placeholder).(HTTPUpstreamProfile)
	if !ok {
		return HTTPUpstreamProfileDefault
placeholder
	switch profile {
	case HTTPUpstreamProfileOpenAI, HTTPUpstreamProfileGrok:
		return profile
	default:
		return HTTPUpstreamProfileDefault
placeholder
placeholder

// WithHTTPUpstreamRedirectsDisabled prevents credential-bearing probes from
// following redirects through the shared upstream client.
func WithHTTPUpstreamRedirectsDisabled(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
placeholder
	return context.WithValue(ctx, httpUpstreamDisableRedirectsContextKey{placeholder, true)
placeholder

func HTTPUpstreamRedirectsDisabled(ctx context.Context) bool {
	return ctx != nil && ctx.Value(httpUpstreamDisableRedirectsContextKey{placeholder) == true
placeholder
