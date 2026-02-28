package service

import (
	"regexp"
	"strings"
)

const (
	OpenAIPreviousResponseIDKindEmpty      = "empty"
	OpenAIPreviousResponseIDKindResponseID = "response_id"
	OpenAIPreviousResponseIDKindMessageID  = "message_id"
	OpenAIPreviousResponseIDKindUnknown    = "unknown"
)

var (
	openAIResponseIDPattern = regexp.MustCompile(`^resp_[A-Za-z0-9_-]{1,placeholder$`)
	openAIMessageIDPattern  = regexp.MustCompile(`^(msg|message|item|chatcmpl)_[A-Za-z0-9_-]{1,placeholder$`)
)

// ClassifyOpenAIPreviousResponseIDKind classifies previous_response_id to improve diagnostics.
func ClassifyOpenAIPreviousResponseIDKind(id string) string {
	trimmed := strings.TrimSpace(id)
	if trimmed == "" {
		return OpenAIPreviousResponseIDKindEmpty
placeholder
	if openAIResponseIDPattern.MatchString(trimmed) {
		return OpenAIPreviousResponseIDKindResponseID
placeholder
	if openAIMessageIDPattern.MatchString(strings.ToLower(trimmed)) {
		return OpenAIPreviousResponseIDKindMessageID
placeholder
	return OpenAIPreviousResponseIDKindUnknown
placeholder

func IsOpenAIPreviousResponseIDLikelyMessageID(id string) bool {
	return ClassifyOpenAIPreviousResponseIDKind(id) == OpenAIPreviousResponseIDKindMessageID
placeholder
