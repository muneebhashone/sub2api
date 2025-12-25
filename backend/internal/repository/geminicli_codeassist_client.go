package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/geminicli"
	"github.com/Wei-Shaw/sub2api/internal/service/ports"

	"github.com/imroc/req/v3"
)

type geminiCliCodeAssistClient struct {
	baseURL string
placeholder

func NewGeminiCliCodeAssistClient() ports.GeminiCliCodeAssistClient {
	return &geminiCliCodeAssistClient{baseURL: geminicli.GeminiCliBaseURLplaceholder
placeholder

func (c *geminiCliCodeAssistClient) LoadCodeAssist(ctx context.Context, accessToken, proxyURL string, reqBody *geminicli.LoadCodeAssistRequest) (*geminicli.LoadCodeAssistResponse, error) {
	if reqBody == nil {
		reqBody = defaultLoadCodeAssistRequest()
placeholder

	var out geminicli.LoadCodeAssistResponse
	resp, err := createGeminiCliReqClient(proxyURL).R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", geminicli.GeminiCLIUserAgent).
		SetBody(reqBody).
		SetSuccessResult(&out).
		Post(c.baseURL + "/v1internal:loadCodeAssist")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
placeholder
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("loadCodeAssist failed: status %d, body: %s", resp.StatusCode, geminicli.SanitizeBodyForLogs(resp.String()))
placeholder
	return &out, nil
placeholder

func (c *geminiCliCodeAssistClient) OnboardUser(ctx context.Context, accessToken, proxyURL string, reqBody *geminicli.OnboardUserRequest) (*geminicli.OnboardUserResponse, error) {
	if reqBody == nil {
		reqBody = defaultOnboardUserRequest()
placeholder

	var out geminicli.OnboardUserResponse
	resp, err := createGeminiCliReqClient(proxyURL).R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", geminicli.GeminiCLIUserAgent).
		SetBody(reqBody).
		SetSuccessResult(&out).
		Post(c.baseURL + "/v1internal:onboardUser")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
placeholder
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("onboardUser failed: status %d, body: %s", resp.StatusCode, geminicli.SanitizeBodyForLogs(resp.String()))
placeholder
	return &out, nil
placeholder

func createGeminiCliReqClient(proxyURL string) *req.Client {
	client := req.C().SetTimeout(30 * time.Second)
	if proxyURL != "" {
		client.SetProxyURL(proxyURL)
placeholder
	return client
placeholder

func defaultLoadCodeAssistRequest() *geminicli.LoadCodeAssistRequest {
	return &geminicli.LoadCodeAssistRequest{
		Metadata: geminicli.LoadCodeAssistMetadata{
			IDEType:    "ANTIGRAVITY",
			Platform:   "PLATFORM_UNSPECIFIED",
			PluginType: "GEMINI",
	placeholder,
placeholder
placeholder

func defaultOnboardUserRequest() *geminicli.OnboardUserRequest {
	return &geminicli.OnboardUserRequest{
		TierID: "LEGACY",
		Metadata: geminicli.LoadCodeAssistMetadata{
			IDEType:    "ANTIGRAVITY",
			Platform:   "PLATFORM_UNSPECIFIED",
			PluginType: "GEMINI",
	placeholder,
placeholder
placeholder
