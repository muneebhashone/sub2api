package admin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

func TestParseCodexSessionImportEntriesSupportsRawTokenJSONAndArray(t *testing.T) {
	token1 := "raw-access-token-1"
	token2 := buildCodexImportTestJWT(t, time.Now().Add(time.Hour), map[string]any{
		"email": "json@example.com",
placeholder)
	token3 := "raw-access-token-3"

	req := CodexSessionImportRequest{
		Content: fmt.Sprintf("%s\n{\"accessToken\":%qplaceholder\n[%q]", token1, token2, token3),
placeholder

	entries, err := parseCodexSessionImportEntries(req)
	if err != nil {
		t.Fatalf("parseCodexSessionImportEntries error = %v", err)
placeholder
	if len(entries) != 3 {
		t.Fatalf("len(entries) = %d, want 3", len(entries))
placeholder

	first, err := normalizeCodexImportEntry(entries[0])
	if err != nil {
		t.Fatalf("normalize raw token error = %v", err)
placeholder
	if first.Credentials["access_token"] != token1 {
		t.Fatalf("raw token access_token = %v, want %s", first.Credentials["access_token"], token1)
placeholder

	second, err := normalizeCodexImportEntry(entries[1])
	if err != nil {
		t.Fatalf("normalize json token error = %v", err)
placeholder
	if second.Email != "json@example.com" {
		t.Fatalf("email = %q, want json@example.com", second.Email)
placeholder

	third, err := normalizeCodexImportEntry(entries[2])
	if err != nil {
		t.Fatalf("normalize array token error = %v", err)
placeholder
	if third.Credentials["access_token"] != token3 {
		t.Fatalf("array token access_token = %v, want %s", third.Credentials["access_token"], token3)
placeholder
placeholder

func TestParseCodexSessionImportEntriesFallsBackToLineModeForMixedJSONAndToken(t *testing.T) {
	req := CodexSessionImportRequest{
		Content: "{\"accessToken\":\"json-line-token\"placeholder\nraw-line-token",
placeholder

	entries, err := parseCodexSessionImportEntries(req)
	if err != nil {
		t.Fatalf("parseCodexSessionImportEntries error = %v", err)
placeholder
	if len(entries) != 2 {
		t.Fatalf("len(entries) = %d, want 2", len(entries))
placeholder

	first, err := normalizeCodexImportEntry(entries[0])
	if err != nil {
		t.Fatalf("normalize json line error = %v", err)
placeholder
	if first.Credentials["access_token"] != "json-line-token" {
		t.Fatalf("json line access_token = %v, want json-line-token", first.Credentials["access_token"])
placeholder

	second, err := normalizeCodexImportEntry(entries[1])
	if err != nil {
		t.Fatalf("normalize raw line error = %v", err)
placeholder
	if second.Credentials["access_token"] != "raw-line-token" {
		t.Fatalf("raw line access_token = %v, want raw-line-token", second.Credentials["access_token"])
placeholder
placeholder

func TestNormalizeCodexSessionJSONExtractsCredentialsAndIgnoresSessionToken(t *testing.T) {
	accessToken := buildCodexImportTestJWT(t, time.Now().Add(time.Hour), map[string]any{
		"email": "claim@example.com",
		"https://api.openai.com/auth": map[string]any{
			"chatgpt_account_id": "acct-from-claim",
			"chatgpt_user_id":    "user-from-claim",
			"chatgpt_plan_type":  "plus",
			"poid":               "org-from-claim",
	placeholder,
placeholder)
	raw := map[string]any{
		"user": map[string]any{
			"id":    "user-from-json",
			"name":  "Sup OO",
			"email": "json@example.com",
			"image": "https://example.com/avatar.png",
	placeholder,
		"account": map[string]any{
			"id":       "acct-from-json",
			"planType": "free",
	placeholder,
		"accessToken":  accessToken,
		"sessionToken": "secret-session-token",
		"expires":      "2026-08-05T13:40:42.836Z",
placeholder

	item, err := normalizeCodexImportEntry(codexImportEntry{Index: 1, Value: rawplaceholder)
	if err != nil {
		t.Fatalf("normalizeCodexImportEntry error = %v", err)
placeholder
	if item.Credentials["access_token"] != accessToken {
		t.Fatalf("access_token not stored")
placeholder
	if item.Credentials["email"] != "json@example.com" {
		t.Fatalf("email = %v, want json@example.com", item.Credentials["email"])
placeholder
	if item.Credentials["chatgpt_account_id"] != "acct-from-json" {
		t.Fatalf("chatgpt_account_id = %v, want acct-from-json", item.Credentials["chatgpt_account_id"])
placeholder
	if item.Credentials["chatgpt_user_id"] != "user-from-json" {
		t.Fatalf("chatgpt_user_id = %v, want user-from-json", item.Credentials["chatgpt_user_id"])
placeholder
	if item.Credentials["plan_type"] != "free" {
		t.Fatalf("plan_type = %v, want free", item.Credentials["plan_type"])
placeholder
	if _, ok := item.Credentials["session_token"]; ok {
		t.Fatalf("session_token should not be written to credentials")
placeholder
	if item.Extra["session_token_present"] != true {
		t.Fatalf("session_token_present = %v, want true", item.Extra["session_token_present"])
placeholder
	if item.Extra["session_expires_at"] != "2026-08-05T13:40:42Z" {
		t.Fatalf("session_expires_at = %v", item.Extra["session_expires_at"])
placeholder
	if item.TokenExpiresAt == nil {
		t.Fatalf("TokenExpiresAt should be parsed from accessToken")
placeholder
placeholder

func TestMergeCodexImportCredentialsClearsStaleRefreshFieldsWhenIncomingHasNoRefreshToken(t *testing.T) {
	existing := map[string]any{
		"access_token":       "old-access-token",
		"refresh_token":      "old-refresh-token",
		"client_id":          "old-client-id",
		"id_token":           "old-id-token",
		"model_mapping":      map[string]any{"from": "existing"placeholder,
		"chatgpt_account_id": "acct-old",
		"unrelated_existing": "keep",
placeholder
	incoming := map[string]any{
		"access_token":       "new-access-token",
		"expires_at":         "2026-08-05T13:40:42Z",
		"chatgpt_account_id": "acct-new",
placeholder
	item := &codexImportAccount{
		AccessToken: "new-access-token",
placeholder

	merged := mergeCodexImportCredentials(existing, incoming, item)

	if merged["access_token"] != "new-access-token" {
		t.Fatalf("access_token = %v, want new-access-token", merged["access_token"])
placeholder
	if merged["chatgpt_account_id"] != "acct-new" {
		t.Fatalf("chatgpt_account_id = %v, want acct-new", merged["chatgpt_account_id"])
placeholder
	if _, ok := merged["refresh_token"]; ok {
		t.Fatalf("refresh_token should be cleared")
placeholder
	if _, ok := merged["client_id"]; ok {
		t.Fatalf("client_id should be cleared")
placeholder
	if _, ok := merged["id_token"]; ok {
		t.Fatalf("id_token should be cleared")
placeholder
	if merged["unrelated_existing"] != "keep" {
		t.Fatalf("unrelated_existing = %v, want keep", merged["unrelated_existing"])
placeholder
	if _, ok := merged["model_mapping"]; !ok {
		t.Fatalf("model_mapping should be preserved")
placeholder
placeholder

func TestMergeCodexImportCredentialsKeepsRefreshFieldsWhenIncomingHasRefreshToken(t *testing.T) {
	existing := map[string]any{
		"refresh_token": "old-refresh-token",
		"client_id":     "old-client-id",
		"id_token":      "old-id-token",
placeholder
	incoming := map[string]any{
		"access_token":  "new-access-token",
		"refresh_token": "new-refresh-token",
		"client_id":     "new-client-id",
		"id_token":      "new-id-token",
placeholder
	item := &codexImportAccount{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		IDToken:      "new-id-token",
placeholder

	merged := mergeCodexImportCredentials(existing, incoming, item)

	if merged["refresh_token"] != "new-refresh-token" {
		t.Fatalf("refresh_token = %v, want new-refresh-token", merged["refresh_token"])
placeholder
	if merged["client_id"] != "new-client-id" {
		t.Fatalf("client_id = %v, want new-client-id", merged["client_id"])
placeholder
	if merged["id_token"] != "new-id-token" {
		t.Fatalf("id_token = %v, want new-id-token", merged["id_token"])
placeholder
placeholder

func TestNormalizeCodexImportRejectsExpiredAccessToken(t *testing.T) {
	expiredToken := buildCodexImportTestJWT(t, time.Now().Add(-time.Hour), map[string]any{placeholder)

	_, err := normalizeCodexImportEntry(codexImportEntry{Index: 1, Value: expiredTokenplaceholder)
	if err == nil {
		t.Fatal("normalizeCodexImportEntry error = nil, want expired token error")
placeholder
	if !strings.Contains(err.Error(), "已过期") {
		t.Fatalf("error = %v, want expired token message", err)
placeholder
placeholder

func TestResolveCodexImportExpiryForNoRefreshTokenUsesTokenExpiry(t *testing.T) {
	tokenExpiresAt := time.Now().Add(time.Hour).UTC()
	item := &codexImportAccount{
		AccessToken:    "access-token",
		Credentials:    map[string]any{"access_token": "access-token"placeholder,
		TokenExpiresAt: &tokenExpiresAt,
		WarningTexts:   []string{placeholder,
placeholder
	disabled := false
	req := CodexSessionImportRequest{AutoPauseOnExpired: &disabledplaceholder

	accountExpiresAt, credentialExpiresAt, autoPause, warnings, err := resolveCodexImportExpiry(req, item)
	if err != nil {
		t.Fatalf("resolveCodexImportExpiry error = %v", err)
placeholder
	if accountExpiresAt == nil || *accountExpiresAt != tokenExpiresAt.Unix() {
		t.Fatalf("account expires_at = %v, want %d", accountExpiresAt, tokenExpiresAt.Unix())
placeholder
	if credentialExpiresAt == nil || credentialExpiresAt.Unix() != tokenExpiresAt.Unix() {
		t.Fatalf("credential expires_at = %v, want %s", credentialExpiresAt, tokenExpiresAt)
placeholder
	if autoPause == nil || !*autoPause {
		t.Fatalf("autoPause = %v, want true", autoPause)
placeholder
	if len(warnings) == 0 {
		t.Fatalf("warnings should not be empty")
placeholder
placeholder

func TestResolveCodexImportExpiryForNoRefreshTokenRequiresExpiry(t *testing.T) {
	item := &codexImportAccount{
		AccessToken:  "opaque-access-token",
		Credentials:  map[string]any{"access_token": "opaque-access-token"placeholder,
		WarningTexts: []string{placeholder,
placeholder

	_, _, _, _, err := resolveCodexImportExpiry(CodexSessionImportRequest{placeholder, item)
	if err == nil {
		t.Fatal("resolveCodexImportExpiry error = nil, want missing expiry error")
placeholder
	if !strings.Contains(err.Error(), "无法解析 accessToken 过期时间") {
		t.Fatalf("error = %v, want missing expiry message", err)
placeholder
placeholder

func TestResolveCodexImportExpiryForNoRefreshTokenUsesEarlierRequestExpiry(t *testing.T) {
	tokenExpiresAt := time.Now().Add(2 * time.Hour).UTC()
	requestExpiresAt := time.Now().Add(time.Hour).UTC()
	item := &codexImportAccount{
		AccessToken:    "access-token",
		Credentials:    map[string]any{"access_token": "access-token"placeholder,
		TokenExpiresAt: &tokenExpiresAt,
		WarningTexts:   []string{placeholder,
placeholder
	reqUnix := requestExpiresAt.Unix()
	req := CodexSessionImportRequest{ExpiresAt: &reqUnixplaceholder

	accountExpiresAt, credentialExpiresAt, _, _, err := resolveCodexImportExpiry(req, item)
	if err != nil {
		t.Fatalf("resolveCodexImportExpiry error = %v", err)
placeholder
	if accountExpiresAt == nil || *accountExpiresAt != requestExpiresAt.Unix() {
		t.Fatalf("account expires_at = %v, want %d", accountExpiresAt, requestExpiresAt.Unix())
placeholder
	if credentialExpiresAt == nil || credentialExpiresAt.Unix() != requestExpiresAt.Unix() {
		t.Fatalf("credential expires_at = %v, want %s", credentialExpiresAt, requestExpiresAt)
placeholder
placeholder

func TestCodexIdentityKeysPreferStrongIdentifiers(t *testing.T) {
	keys := buildCodexIdentityKeys("acct-1", "user-1", "same@example.com", "token")
	if len(keys) == 0 || keys[0] != "user:user-1" {
		t.Fatalf("user key should have highest priority: %v", keys)
placeholder
	if keys[len(keys)-1] != "account:acct-1" {
		t.Fatalf("shared account key should be the last fallback: %v", keys)
placeholder
	for _, key := range keys {
		if strings.HasPrefix(key, "email:") {
			t.Fatalf("strong identity should not include email fallback: %v", keys)
	placeholder
placeholder

	keys = buildCodexIdentityKeys("", "", "same@example.com", "token")
	hasEmail := false
	for _, key := range keys {
		if key == "email:same@example.com" {
			hasEmail = true
	placeholder
placeholder
	if !hasEmail {
		t.Fatalf("weak identity should include email fallback: %v", keys)
placeholder
placeholder

func TestCodexAccountIndexDoesNotMatchDifferentUsersInSameChatGPTAccount(t *testing.T) {
	existing := service.Account{
		ID: 10,
placeholder
			"chatgpt_account_id": "team-1",
			"chatgpt_user_id":    "user-1",
			"access_token":       "token-1",
	placeholder,
placeholder
	index := buildCodexAccountIndex([]service.Account{existingplaceholder)

	keys := buildCodexIdentityKeys("team-1", "user-2", "", "token-2")
	if got, _ := index.Find(keys, "user-2"); got != nil {
		t.Fatalf("Find matched account ID %d for a different chatgpt_user_id in the same team", got.ID)
placeholder

	keys = buildCodexIdentityKeys("team-1", "user-1", "", "token-2")
	got, _ := index.Find(keys, "user-1")
	if got == nil || got.ID != existing.ID {
		t.Fatalf("Find by same chatgpt_user_id = %v, want account ID %d", got, existing.ID)
placeholder
placeholder

func TestCodexAccountIndexFallsBackToAccountKeyWhenUserIDMissing(t *testing.T) {
	// 存量账号缺少 chatgpt_user_id：携带 user id 的重新导入应命中并更新（回填），
	// 而不是创建重复账号。
	legacy := service.Account{
		ID: 20,
placeholder
			"chatgpt_account_id": "team-1",
			"access_token":       "token-old",
	placeholder,
placeholder
	index := buildCodexAccountIndex([]service.Account{legacyplaceholder)

	keys := buildCodexIdentityKeys("team-1", "user-1", "", "token-new")
	got, matchedKey := index.Find(keys, "user-1")
	if got == nil || got.ID != legacy.ID {
		t.Fatalf("Find legacy account without stored user id = %v, want account ID %d", got, legacy.ID)
placeholder
	if matchedKey != "account:team-1" {
		t.Fatalf("matched key = %q, want account:team-1", matchedKey)
placeholder

	// 反向：导入条目无法解析出 user id 时，仍应通过 account 键命中已有账号。
	full := service.Account{
		ID: 21,
placeholder
			"chatgpt_account_id": "team-2",
			"chatgpt_user_id":    "user-9",
			"access_token":       "token-old",
	placeholder,
placeholder
	index = buildCodexAccountIndex([]service.Account{fullplaceholder)

	keys = buildCodexIdentityKeys("team-2", "", "", "token-opaque")
	got, _ = index.Find(keys, "")
	if got == nil || got.ID != full.ID {
		t.Fatalf("Find by account key without entry user id = %v, want account ID %d", got, full.ID)
placeholder
placeholder

func TestCodexAccountIndexKeepsAllCandidatesForSharedAccountKey(t *testing.T) {
	legacy := service.Account{
		ID: 30,
placeholder
			"chatgpt_account_id": "team-1",
			"access_token":       "token-legacy",
	placeholder,
placeholder
	member := service.Account{
		ID: 31,
placeholder
			"chatgpt_account_id": "team-1",
			"chatgpt_user_id":    "user-2",
			"access_token":       "token-member",
	placeholder,
placeholder

	// 无论索引构建顺序如何，携带新 user id 的条目都应跳过 user-2 的账号、
	// 命中缺少 user id 的存量账号，而不是因单一候选被遮蔽而落空。
	for _, accounts := range [][]service.Account{
		{member, legacyplaceholder,
		{legacy, memberplaceholder,
placeholder {
		index := buildCodexAccountIndex(accounts)

		keys := buildCodexIdentityKeys("team-1", "user-1", "", "token-new")
		got, matchedKey := index.Find(keys, "user-1")
		if got == nil || got.ID != legacy.ID {
			t.Fatalf("Find with shared account key = %v, want legacy account ID %d", got, legacy.ID)
	placeholder
		if matchedKey != "account:team-1" {
			t.Fatalf("matched key = %q, want account:team-1", matchedKey)
	placeholder

		keys = buildCodexIdentityKeys("team-1", "user-2", "", "token-new")
		got, matchedKey = index.Find(keys, "user-2")
		if got == nil || got.ID != member.ID {
			t.Fatalf("Find by user key = %v, want member account ID %d", got, member.ID)
	placeholder
		if matchedKey != "user:user-2" {
			t.Fatalf("matched key = %q, want user:user-2", matchedKey)
	placeholder
placeholder
placeholder

func TestCodexAccountIndexUpsertReplacesSameAccount(t *testing.T) {
	legacy := service.Account{
		ID: 40,
placeholder
			"chatgpt_account_id": "team-1",
			"access_token":       "token-old",
	placeholder,
placeholder
	index := buildCodexAccountIndex([]service.Account{legacyplaceholder)

	backfilled := service.Account{
		ID: 40,
placeholder
			"chatgpt_account_id": "team-1",
			"chatgpt_user_id":    "user-1",
			"access_token":       "token-new",
	placeholder,
placeholder
	index.Add(backfilled)

	// 回填后同一账号在 account 键下应被原位替换而非残留旧副本：
	// 其他成员的条目不应再通过旧副本（无 user id）命中该账号。
	keys := buildCodexIdentityKeys("team-1", "user-2", "", "token-other")
	if got, _ := index.Find(keys, "user-2"); got != nil {
		t.Fatalf("stale candidate matched after upsert: account ID %d", got.ID)
placeholder

	keys = buildCodexIdentityKeys("team-1", "user-1", "", "token-other")
	got, _ := index.Find(keys, "user-1")
	if got == nil || got.ID != backfilled.ID {
		t.Fatalf("Find after upsert = %v, want account ID %d", got, backfilled.ID)
placeholder
	if uid := codexCredentialString(got.Credentials, "chatgpt_user_id"); uid != "user-1" {
		t.Fatalf("upsert did not replace credentials, chatgpt_user_id = %q", uid)
placeholder
placeholder

func TestCodexIdentitySeenDistinguishesTeamMembers(t *testing.T) {
	seen := map[string]codexSeenIdentity{placeholder
	member1 := buildCodexIdentityKeys("team-1", "user-1", "", "token-1")
	markCodexIdentitySeen(seen, member1, 1, "user-1")

	member2 := buildCodexIdentityKeys("team-1", "user-2", "", "token-2")
	if index, ok := firstSeenCodexIdentity(seen, member2, "user-2"); ok {
		t.Fatalf("different team member treated as duplicate of entry %d", index)
placeholder

	again := buildCodexIdentityKeys("team-1", "user-1", "", "token-3")
	index, ok := firstSeenCodexIdentity(seen, again, "user-1")
	if !ok || index != 1 {
		t.Fatalf("same user re-entry dedup = (%d, %v), want (1, true)", index, ok)
placeholder

	// 无 user id 的条目与已见同 account 条目视为重复（保守跳过，与既有行为一致）。
	opaque := buildCodexIdentityKeys("team-1", "", "", "token-4")
	index, ok = firstSeenCodexIdentity(seen, opaque, "")
	if !ok || index != 1 {
		t.Fatalf("entry without user id dedup = (%d, %v), want (1, true)", index, ok)
placeholder
placeholder

func buildCodexImportTestJWT(t *testing.T, exp time.Time, extraClaims map[string]any) string {
placeholder
	header := map[string]any{
		"alg": "none",
		"typ": "JWT",
placeholder
	claims := map[string]any{
		"sub": "user-from-sub",
		"exp": exp.Unix(),
		"iat": time.Now().Unix(),
placeholder
	for k, v := range extraClaims {
		claims[k] = v
placeholder
	headerBytes, err := json.Marshal(header)
	if err != nil {
		t.Fatalf("marshal header: %v", err)
placeholder
	claimBytes, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("marshal claims: %v", err)
placeholder
	return base64.RawURLEncoding.EncodeToString(headerBytes) + "." + base64.RawURLEncoding.EncodeToString(claimBytes) + "."
placeholder
