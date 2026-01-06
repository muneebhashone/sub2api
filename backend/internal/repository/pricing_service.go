package repository

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/httpclient"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type pricingRemoteClient struct {
	httpClient *http.Client
placeholder

func NewPricingRemoteClient(cfg *config.Config) service.PricingRemoteClient {
	allowPrivate := false
	validateResolvedIP := true
	if cfg != nil {
		allowPrivate = cfg.Security.URLAllowlist.AllowPrivateHosts
		validateResolvedIP = cfg.Security.URLAllowlist.Enabled
placeholder
	sharedClient, err := httpclient.GetClient(httpclient.Options{
		Timeout:            30 * time.Second,
		ValidateResolvedIP: validateResolvedIP,
		AllowPrivateHosts:  allowPrivate,
placeholder)
	if err != nil {
		sharedClient = &http.Client{Timeout: 30 * time.Secondplaceholder
placeholder
	return &pricingRemoteClient{
		httpClient: sharedClient,
placeholder
placeholder

func (c *pricingRemoteClient) FetchPricingJSON(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
placeholder

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
placeholder

	return io.ReadAll(resp.Body)
placeholder

func (c *pricingRemoteClient) FetchHashText(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
placeholder

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
placeholder

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
placeholder

	// 哈希文件格式：hash  filename 或者纯 hash
	hash := strings.TrimSpace(string(body))
	parts := strings.Fields(hash)
	if len(parts) > 0 {
		return parts[0], nil
placeholder
	return hash, nil
placeholder
