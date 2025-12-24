package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

type githubReleaseClient struct {
	httpClient *http.Client
placeholder

func NewGitHubReleaseClient() service.GitHubReleaseClient {
	return &githubReleaseClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
	placeholder,
placeholder
placeholder

func (c *githubReleaseClient) FetchLatestRelease(ctx context.Context, repo string) (*service.GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Sub2API-Updater")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
placeholder

	var release service.GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
placeholder

	return &release, nil
placeholder

func (c *githubReleaseClient) DownloadFile(ctx context.Context, url, dest string, maxSize int64) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
placeholder

	client := &http.Client{Timeout: 10 * time.Minuteplaceholder
	resp, err := client.Do(req)
	if err != nil {
		return err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned %d", resp.StatusCode)
placeholder

	// SECURITY: Check Content-Length if available
	if resp.ContentLength > maxSize {
		return fmt.Errorf("file too large: %d bytes (max %d)", resp.ContentLength, maxSize)
placeholder

	out, err := os.Create(dest)
	if err != nil {
		return err
placeholder
	defer func() { _ = out.Close() placeholder()

	// SECURITY: Use LimitReader to enforce max download size even if Content-Length is missing/wrong
	limited := io.LimitReader(resp.Body, maxSize+1)
	written, err := io.Copy(out, limited)
	if err != nil {
		return err
placeholder

	// Check if we hit the limit (downloaded more than maxSize)
	if written > maxSize {
		_ = os.Remove(dest) // Clean up partial file (best-effort)
		return fmt.Errorf("download exceeded maximum size of %d bytes", maxSize)
placeholder

	return nil
placeholder

func (c *githubReleaseClient) FetchChecksumFile(ctx context.Context, url string) ([]byte, error) {
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
