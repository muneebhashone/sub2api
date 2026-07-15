package service

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/httpclient"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"
)

const (
	OpenAIAuthModeAgentIdentity          = "agentIdentity"
	agentIdentityAuthAPIBaseURL          = "https://auth.openai.com/api/accounts"
	agentIdentityTaskRegistrationTimeout = 30 * time.Second
)

var openAIAgentIdentityAuthAPIBaseURL = agentIdentityAuthAPIBaseURL

var agentIdentityTaskLocks sync.Map // map[int64]*sync.Mutex

type agentIdentityWSConnectionInvalidator interface {
	InvalidateAgentIdentityWSConnections(accountID int64)
placeholder

type agentIdentityKey struct {
	runtimeID  string
	privateKey ed25519.PrivateKey
	taskID     string
placeholder

type agentIdentityTaskRegistrationResponse struct {
	TaskID               string `json:"task_id"`
	TaskIDCamel          string `json:"taskId"`
	EncryptedTaskID      string `json:"encrypted_task_id"`
	EncryptedTaskIDCamel string `json:"encryptedTaskId"`
placeholder

type agentIdentityTaskRecoveredError struct{placeholder

func (e *agentIdentityTaskRecoveredError) Error() string {
	return "agent identity task recovered"
placeholder

func (a *Account) IsOpenAIAgentIdentity() bool {
	if a == nil || !a.IsOpenAIOAuth() {
		return false
placeholder
	return strings.EqualFold(strings.TrimSpace(a.GetCredential(openAIAuthModeCredentialKey)), OpenAIAuthModeAgentIdentity)
placeholder

func agentIdentityPrivateKey(account *Account) (ed25519.PrivateKey, error) {
	if account == nil {
		return nil, errors.New("agent identity account is nil")
placeholder
	raw := strings.TrimSpace(account.GetCredential("agent_private_key"))
	if raw == "" {
		return nil, errors.New("agent identity private key is missing")
placeholder
	der, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, errors.New("agent identity private key is not valid base64")
placeholder
	key, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, errors.New("agent identity private key is not valid PKCS#8")
placeholder
	privateKey, ok := key.(ed25519.PrivateKey)
	if !ok || len(privateKey) != ed25519.PrivateKeySize {
		return nil, errors.New("agent identity private key is not Ed25519")
placeholder
	return privateKey, nil
placeholder

// ValidateOpenAIAgentIdentityPrivateKey validates the stored PKCS#8 Ed25519
// form without returning or logging the key material.
func ValidateOpenAIAgentIdentityPrivateKey(encoded string) error {
	account := &Account{Credentials: map[string]any{"agent_private_key": encodedplaceholderplaceholder
	_, err := agentIdentityPrivateKey(account)
	return err
placeholder

func agentIdentityKeyFromAccount(account *Account) (agentIdentityKey, error) {
	privateKey, err := agentIdentityPrivateKey(account)
	if err != nil {
		return agentIdentityKey{placeholder, err
placeholder
	runtimeID := strings.TrimSpace(account.GetCredential("agent_runtime_id"))
	if runtimeID == "" {
		return agentIdentityKey{placeholder, errors.New("agent identity runtime id is missing")
placeholder
	return agentIdentityKey{
		runtimeID:  runtimeID,
		privateKey: privateKey,
		taskID:     strings.TrimSpace(account.GetCredential("task_id")),
placeholder, nil
placeholder

func buildAgentAssertion(key agentIdentityKey, now time.Time) (string, error) {
	if key.runtimeID == "" || key.taskID == "" {
		return "", errors.New("agent identity runtime or task id is missing")
placeholder
	timestamp := now.UTC().Format(time.RFC3339)
	payload := []byte(key.runtimeID + ":" + key.taskID + ":" + timestamp)
	signature, err := key.privateKey.Sign(nil, payload, crypto.Hash(0))
	if err != nil {
		return "", errors.New("failed to sign agent assertion")
placeholder
	envelope := map[string]string{
		"agent_runtime_id": key.runtimeID,
		"task_id":          key.taskID,
		"timestamp":        timestamp,
		"signature":        base64.StdEncoding.EncodeToString(signature),
placeholder
	encoded, err := json.Marshal(envelope)
	if err != nil {
		return "", errors.New("failed to serialize agent assertion")
placeholder
	return "AgentAssertion " + base64.RawURLEncoding.EncodeToString(encoded), nil
placeholder

func signAgentTaskRegistration(key agentIdentityKey, timestamp time.Time) (string, string, error) {
	if key.runtimeID == "" {
		return "", "", errors.New("agent identity runtime id is missing")
placeholder
	formatted := timestamp.UTC().Format(time.RFC3339)
	signature, err := key.privateKey.Sign(nil, []byte(key.runtimeID+":"+formatted), crypto.Hash(0))
	if err != nil {
		return "", "", errors.New("failed to sign agent task registration")
placeholder
	return formatted, base64.StdEncoding.EncodeToString(signature), nil
placeholder

func decryptAgentTaskID(key agentIdentityKey, encoded string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(strings.TrimSpace(encoded))
	if err != nil {
		return "", errors.New("encrypted agent task id is not valid base64")
placeholder
	seed := key.privateKey.Seed()
	digest := sha512.Sum512(seed)
	var curvePrivate [32]byte
	copy(curvePrivate[:], digest[:32])
	curvePrivate[0] &= 248
	curvePrivate[31] &= 127
	curvePrivate[31] |= 64
	curvePublicBytes, err := curve25519.X25519(curvePrivate[:], curve25519.Basepoint)
	if err != nil {
		return "", errors.New("failed to derive agent identity decryption key")
placeholder
	var curvePublic [32]byte
	copy(curvePublic[:], curvePublicBytes)
	plaintext, ok := box.OpenAnonymous(nil, ciphertext, &curvePublic, &curvePrivate)
	if !ok {
		return "", errors.New("failed to decrypt encrypted agent task id")
placeholder
	taskID := strings.TrimSpace(string(plaintext))
	if taskID == "" {
		return "", errors.New("decrypted agent task id is empty")
placeholder
	return taskID, nil
placeholder

func registerAgentIdentityTask(ctx context.Context, account *Account) (string, error) {
	key, err := agentIdentityKeyFromAccount(account)
	if err != nil {
		return "", err
placeholder
	timestamp, signature, err := signAgentTaskRegistration(key, time.Now())
	if err != nil {
		return "", err
placeholder
	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder
	client, err := httpclient.GetClient(httpclient.Options{
		ProxyURL:              proxyURL,
		Timeout:               agentIdentityTaskRegistrationTimeout,
		ResponseHeaderTimeout: 15 * time.Second,
placeholder)
	if err != nil {
		return "", errors.New("invalid proxy configuration for agent task registration")
placeholder
	body, err := json.Marshal(map[string]string{
		"timestamp": timestamp,
		"signature": signature,
placeholder)
	if err != nil {
		return "", errors.New("failed to serialize agent task registration")
placeholder
	url := strings.TrimRight(strings.TrimSpace(openAIAgentIdentityAuthAPIBaseURL), "/") + "/v1/agent/" + key.runtimeID + "/task/register"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return "", errors.New("failed to build agent task registration request")
placeholder
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("agent task registration request failed")
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("agent task registration returned status %d", resp.StatusCode)
placeholder
	var result agentIdentityTaskRegistrationResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 64*1024)).Decode(&result); err != nil {
		return "", errors.New("agent task registration response is invalid")
placeholder
	if taskID := strings.TrimSpace(result.TaskID); taskID != "" {
		return taskID, nil
placeholder
	if taskID := strings.TrimSpace(result.TaskIDCamel); taskID != "" {
		return taskID, nil
placeholder
	encrypted := strings.TrimSpace(result.EncryptedTaskID)
	if encrypted == "" {
		encrypted = strings.TrimSpace(result.EncryptedTaskIDCamel)
placeholder
	if encrypted == "" {
		return "", errors.New("agent task registration response omitted task id")
placeholder
	return decryptAgentTaskID(key, encrypted)
placeholder

func ensureAgentIdentityTaskForAccount(ctx context.Context, repo AccountRepository, wsInvalidator agentIdentityWSConnectionInvalidator, taskMu *sync.Mutex, account *Account, expectedTaskID string) error {
	if account == nil || !account.IsOpenAIAgentIdentity() {
		return nil
placeholder
	credAccount := account
	if account.IsShadow() {
		resolved, err := resolveCredentialAccount(ctx, repo, account)
		if err != nil {
			return err
	placeholder
		credAccount = resolved
placeholder
	if credAccount == nil || !credAccount.IsOpenAIAgentIdentity() {
		return errors.New("agent identity credentials are unavailable")
placeholder
	currentTaskID := strings.TrimSpace(credAccount.GetCredential("task_id"))
	if currentTaskID != "" && (expectedTaskID == "" || currentTaskID != expectedTaskID) {
		return nil
placeholder
	if taskMu == nil {
		return errors.New("agent identity task lock is unavailable")
placeholder
	sharedTaskMu := taskMu
	if credAccount.ID > 0 {
		candidate := &sync.Mutex{placeholder
		actual, _ := agentIdentityTaskLocks.LoadOrStore(credAccount.ID, candidate)
		loadedTaskMu, ok := actual.(*sync.Mutex)
		if !ok {
			return errors.New("agent identity task lock has invalid type")
	placeholder
		sharedTaskMu = loadedTaskMu
placeholder
	sharedTaskMu.Lock()
	defer sharedTaskMu.Unlock()
	// Re-read inside the shared lock. Different request paths often receive
	// independent repository snapshots; checking only the caller's snapshot
	// would allow sequential duplicate registrations after the first writer
	// has already persisted a new task.
	if repo != nil && credAccount.ID > 0 {
		if refreshed, refreshErr := repo.GetByID(ctx, credAccount.ID); refreshErr == nil && refreshed != nil {
			if refreshed.IsShadow() {
				if resolved, resolveErr := resolveCredentialAccount(ctx, repo, refreshed); resolveErr == nil && resolved != nil {
					refreshed = resolved
			placeholder
		placeholder
			if refreshed.IsOpenAIAgentIdentity() {
				credAccount = refreshed
				if !account.IsShadow() {
					account.Credentials = shallowCopyMap(credAccount.Credentials)
			placeholder
		placeholder
	placeholder
placeholder
	currentTaskID = strings.TrimSpace(credAccount.GetCredential("task_id"))
	if currentTaskID != "" && (expectedTaskID == "" || currentTaskID != expectedTaskID) {
		return nil
placeholder
	newTaskID, err := registerAgentIdentityTask(ctx, credAccount)
	if err != nil {
		return err
placeholder
	credentials := make(map[string]any, len(credAccount.Credentials)+1)
	for key, value := range credAccount.Credentials {
		credentials[key] = value
placeholder
	credentials["task_id"] = newTaskID
	if err := persistAccountCredentials(ctx, repo, credAccount, credentials); err != nil {
		return err
placeholder
	if !account.IsShadow() && account != credAccount {
		account.Credentials = shallowCopyMap(credAccount.Credentials)
placeholder
	if wsInvalidator != nil {
		wsInvalidator.InvalidateAgentIdentityWSConnections(credAccount.ID)
placeholder
	return nil
placeholder

func (s *OpenAIGatewayService) ensureAgentIdentityTask(ctx context.Context, account *Account, expectedTaskID string) error {
	if s == nil {
		return errors.New("openai gateway service is nil")
placeholder
	return ensureAgentIdentityTaskForAccount(ctx, s.accountRepo, s, &s.agentIdentityTaskMu, account, expectedTaskID)
placeholder

func isAgentIdentityTaskInvalidHTTPResponse(statusCode int, body []byte) bool {
	if statusCode != http.StatusUnauthorized {
		return false
placeholder
	lower := strings.ToLower(string(body))
	compact := strings.NewReplacer(" ", "", "\t", "", "\r", "", "\n", "").Replace(lower)
	for _, marker := range []string{
		`"code":"invalid_task_id"`,
		`"code":"task_not_found"`,
		`"code":"task_expired"`,
		`"error":"invalid_task_id"`,
placeholder {
		if strings.Contains(compact, marker) {
			return true
	placeholder
placeholder
	for _, marker := range []string{
		"invalid task_id",
		"invalid task id",
		"task_id is invalid",
		"task id is invalid",
		"task not found",
		"task expired",
		"unknown task_id",
		"unknown task id",
placeholder {
		if strings.Contains(lower, marker) {
			return true
	placeholder
placeholder
	return false
placeholder

type agentIdentityTaskRecoveryContextKey struct{placeholder

func markAgentIdentityTaskRecoveryTried(ctx context.Context) context.Context {
	return context.WithValue(ctx, agentIdentityTaskRecoveryContextKey{placeholder, true)
placeholder

func agentIdentityTaskRecoveryWasTried(ctx context.Context) bool {
	tried, _ := ctx.Value(agentIdentityTaskRecoveryContextKey{placeholder).(bool)
	return tried
placeholder

func isAgentIdentityTaskInvalidWSDialError(err *openAIWSDialError) bool {
	return err != nil && isAgentIdentityTaskInvalidHTTPResponse(err.StatusCode, err.ResponseBody)
placeholder

func (s *OpenAIGatewayService) buildOpenAIAuthenticationHeaders(ctx context.Context, account *Account, token string) (http.Header, error) {
	if account == nil {
		return nil, errors.New("account is nil")
placeholder
	credAccount := account
	if account.IsShadow() {
		resolved, err := resolveCredentialAccount(ctx, s.accountRepo, account)
		if err != nil {
			return nil, err
	placeholder
		credAccount = resolved
placeholder
	headers := make(http.Header)
	if credAccount != nil && credAccount.IsOpenAIAgentIdentity() {
		agentHeaders, err := buildAgentIdentityAuthenticationHeaders(ctx, s.accountRepo, s, &s.agentIdentityTaskMu, credAccount)
		if err != nil {
			return nil, err
	placeholder
		return agentHeaders, nil
placeholder
	headers.Set("Authorization", "Bearer "+token)
	return headers, nil
placeholder

func buildAgentIdentityAuthenticationHeaders(ctx context.Context, repo AccountRepository, wsInvalidator agentIdentityWSConnectionInvalidator, taskMu *sync.Mutex, account *Account) (http.Header, error) {
	if account == nil || !account.IsOpenAIAgentIdentity() {
		return nil, errors.New("agent identity account is required")
placeholder
	if err := ensureAgentIdentityTaskForAccount(ctx, repo, wsInvalidator, taskMu, account, ""); err != nil {
		return nil, err
placeholder
	key, err := agentIdentityKeyFromAccount(account)
	if err != nil {
		return nil, err
placeholder
	assertion, err := buildAgentAssertion(key, time.Now())
	if err != nil {
		return nil, err
placeholder
	headers := make(http.Header)
	headers.Set("Authorization", assertion)
	return headers, nil
placeholder

func (s *OpenAIGatewayService) refreshOpenAIAgentIdentityHeaders(ctx context.Context, account *Account, headers http.Header) (http.Header, error) {
	if account == nil {
		return cloneHeader(headers), nil
placeholder
	credAccount := account
	if account.IsShadow() {
		resolved, err := resolveCredentialAccount(ctx, s.accountRepo, account)
		if err != nil {
			return nil, err
	placeholder
		credAccount = resolved
placeholder
	if !credAccount.IsOpenAIAgentIdentity() {
		return cloneHeader(headers), nil
placeholder
	refreshed := cloneHeader(headers)
	if refreshed == nil {
		refreshed = make(http.Header)
placeholder
	authHeaders, err := buildAgentIdentityAuthenticationHeaders(ctx, s.accountRepo, s, &s.agentIdentityTaskMu, credAccount)
	if err != nil {
		return nil, err
placeholder
	refreshed.Set("Authorization", authHeaders.Get("Authorization"))
	return refreshed, nil
placeholder

func (s *OpenAIGatewayService) recoverAgentIdentityTask(ctx context.Context, account *Account, expectedTaskID string) error {
	if account != nil && account.IsShadow() {
		if resolved, err := resolveCredentialAccount(ctx, s.accountRepo, account); err == nil && resolved != nil && strings.TrimSpace(expectedTaskID) == "" {
			expectedTaskID = strings.TrimSpace(resolved.GetCredential("task_id"))
	placeholder
placeholder
	return s.ensureAgentIdentityTask(ctx, account, expectedTaskID)
placeholder

func (s *OpenAIGatewayService) isAgentIdentityAccount(ctx context.Context, account *Account) bool {
	if account == nil {
		return false
placeholder
	credAccount := account
	if account.IsShadow() {
		resolved, err := resolveCredentialAccount(ctx, s.accountRepo, account)
		if err != nil {
			return false
	placeholder
		credAccount = resolved
placeholder
	return credAccount != nil && credAccount.IsOpenAIAgentIdentity()
placeholder

// redactAgentIdentitySensitiveBody removes credential values before an
// upstream error can reach logs, ops events, or returned error text. Agent
// Identity responses should not echo these values, but keeping this boundary
// defensive prevents accidental disclosure if an upstream error does.
func redactAgentIdentitySensitiveBodyForAccount(ctx context.Context, repo AccountRepository, account *Account, body []byte) []byte {
	if account == nil || len(body) == 0 {
		return body
placeholder
	credAccount := account
	if account != nil && account.IsShadow() {
		if resolved, err := resolveCredentialAccount(ctx, repo, account); err == nil && resolved != nil {
			credAccount = resolved
	placeholder
placeholder
	if credAccount == nil || !credAccount.IsOpenAIAgentIdentity() {
		return body
placeholder
	redacted := string(body)
	for _, key := range []string{
		"agent_private_key",
		"agent_runtime_id",
		"task_id",
		"access_token",
		"refresh_token",
		"id_token",
		"api_key",
		"session_key",
		"cookie",
placeholder {
		if value := strings.TrimSpace(credAccount.GetCredential(key)); value != "" {
			redacted = strings.ReplaceAll(redacted, value, "[redacted]")
	placeholder
placeholder
	const assertionPrefix = "AgentAssertion "
	for offset := 0; offset < len(redacted); {
		relativeStart := strings.Index(redacted[offset:], assertionPrefix)
		if relativeStart < 0 {
			break
	placeholder
		start := offset + relativeStart
		valueStart := start + len(assertionPrefix)
		end := valueStart
		for end < len(redacted) && !strings.ContainsRune(" \t\r\n\"',placeholder", rune(redacted[end])) {
			end++
	placeholder
		redacted = redacted[:valueStart] + "[redacted]" + redacted[end:]
		offset = valueStart + len("[redacted]")
placeholder
	return []byte(redacted)
placeholder

func (s *OpenAIGatewayService) redactAgentIdentitySensitiveBody(ctx context.Context, account *Account, body []byte) []byte {
	if !s.isAgentIdentityAccount(ctx, account) {
		return body
placeholder
	return redactAgentIdentitySensitiveBodyForAccount(ctx, s.accountRepo, account, body)
placeholder
