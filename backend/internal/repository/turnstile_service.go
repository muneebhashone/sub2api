package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sub2api/internal/service"
)

const turnstileVerifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

type turnstileVerifier struct {
	httpClient *http.Client
placeholder

func NewTurnstileVerifier() service.TurnstileVerifier {
	return &turnstileVerifier{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
	placeholder,
placeholder
placeholder

func (v *turnstileVerifier) VerifyToken(ctx context.Context, secretKey, token, remoteIP string) (*service.TurnstileVerifyResponse, error) {
	formData := url.Values{placeholder
	formData.Set("secret", secretKey)
	formData.Set("response", token)
	if remoteIP != "" {
		formData.Set("remoteip", remoteIP)
placeholder

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, turnstileVerifyURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
placeholder
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	var result service.TurnstileVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
placeholder

	return &result, nil
placeholder
