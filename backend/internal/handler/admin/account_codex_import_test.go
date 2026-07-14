package admin

import (
	"context"
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

func TestMergeCodexImportCredentialsPreservesExistingRefreshFieldsWhenIncomingHasNoRefreshToken(t *testing.T) {
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
	if merged["refresh_token"] != "old-refresh-token" {
		t.Fatalf("refresh_token = %v, want old-refresh-token", merged["refresh_token"])
placeholder
	if merged["client_id"] != "old-client-id" {
		t.Fatalf("client_id = %v, want old-client-id", merged["client_id"])
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
	keys := buildCodexImportIdentityKeys("acct-1", "user-1", "same@example.com", "token", "refresh")
	if len(keys) == 0 || keys[0] != "user:user-1" {
		t.Fatalf("user key should have highest priority when refresh token exists: %v", keys)
placeholder
	if keys[len(keys)-1] != "account:acct-1" {
		t.Fatalf("shared account key should be the last fallback: %v", keys)
placeholder
	for _, key := range keys {
		if strings.HasPrefix(key, "email:") {
			t.Fatalf("strong identity should not include email fallback: %v", keys)
	placeholder
placeholder

	keys = buildCodexImportIdentityKeys("", "", "same@example.com", "token", "refresh")
	hasEmail := false
	for _, key := range keys {
		if key == "email:same@example.com" {
			hasEmail = true
	placeholder
placeholder
	if !hasEmail {
		t.Fatalf("weak identity should include email fallback: %v", keys)
placeholder

	keys = buildCodexImportIdentityKeys("acct-1", "user-1", "same@example.com", "token", "")
	if len(keys) != 1 || !strings.HasPrefix(keys[0], "access:") {
		t.Fatalf("accessToken-only identity should use only access fingerprint: %v", keys)
placeholder
placeholder

func TestCodexAccountIndexDoesNotMatchDifferentUsersInSameChatGPTAccount(t *testing.T) {
	existing := service.Account{
		ID: 10,
placeholder
			"chatgpt_account_id": "team-1",
			"chatgpt_user_id":    "user-1",
			"access_token":       "token-1",
			"refresh_token":      "refresh-1",
	placeholder,
placeholder
	index := buildCodexAccountIndex([]service.Account{existingplaceholder)

	keys := buildCodexImportIdentityKeys("team-1", "user-2", "", "token-2", "refresh-2")
	if got, _ := index.Find(keys, "user-2"); got != nil {
		t.Fatalf("Find matched account ID %d for a different chatgpt_user_id in the same team", got.ID)
placeholder

	keys = buildCodexImportIdentityKeys("team-1", "user-1", "", "token-2", "refresh-2")
	got, _ := index.Find(keys, "user-1")
	if got == nil || got.ID != existing.ID {
		t.Fatalf("Find by same chatgpt_user_id = %v, want account ID %d", got, existing.ID)
placeholder
placeholder

func TestCodexAccountIndexFallsBackToAccountKeyWhenRefreshTokenExistsAndUserIDMissing(t *testing.T) {
	// 含 refresh_token 的常规导入沿用 a5638a4e 的兼容逻辑：存量账号缺少
	// chatgpt_user_id 时，携带 user id 的重新导入仍可命中并回填。
	legacy := service.Account{
		ID: 20,
placeholder
			"chatgpt_account_id": "team-1",
			"access_token":       "token-old",
			"refresh_token":      "refresh-old",
	placeholder,
placeholder
	index := buildCodexAccountIndex([]service.Account{legacyplaceholder)

	keys := buildCodexImportIdentityKeys("team-1", "user-1", "", "token-new", "refresh-new")
	got, matchedKey := index.Find(keys, "user-1")
	if got == nil || got.ID != legacy.ID {
		t.Fatalf("Find legacy account without stored user id = %v, want account ID %d", got, legacy.ID)
placeholder
	if matchedKey != "account:team-1" {
		t.Fatalf("matched key = %q, want account:team-1", matchedKey)
placeholder

	// 反向：含 refresh_token 的导入条目无法解析出 user id 时，仍应通过
	// account 键命中已有账号，保持常规导入去重行为。
	full := service.Account{
		ID: 21,
placeholder
			"chatgpt_account_id": "team-2",
			"chatgpt_user_id":    "user-9",
			"access_token":       "token-old",
			"refresh_token":      "refresh-old",
	placeholder,
placeholder
	index = buildCodexAccountIndex([]service.Account{fullplaceholder)

	keys = buildCodexImportIdentityKeys("team-2", "", "", "token-opaque", "refresh-new")
	got, _ = index.Find(keys, "")
	if got == nil || got.ID != full.ID {
		t.Fatalf("Find by account key without entry user id = %v, want account ID %d", got, full.ID)
placeholder
placeholder

func TestCodexAccountIndexAccessTokenOnlyUsesTokenFingerprint(t *testing.T) {
	existing := service.Account{
		ID: 22,
placeholder
			"chatgpt_account_id": "team-1",
			"chatgpt_user_id":    "user-1",
			"access_token":       "token-old",
	placeholder,
placeholder
	index := buildCodexAccountIndex([]service.Account{existingplaceholder)

	keys := buildCodexImportIdentityKeys("team-1", "user-1", "", "token-new", "")
	if got, matchedKey := index.Find(keys, "user-1"); got != nil {
		t.Fatalf("accessToken-only import matched by %q despite different token: account ID %d", matchedKey, got.ID)
placeholder

	keys = buildCodexImportIdentityKeys("team-1", "user-1", "", "token-old", "")
	got, matchedKey := index.Find(keys, "user-1")
	if got == nil || got.ID != existing.ID {
		t.Fatalf("Find accessToken-only duplicate by fingerprint = %v, want account ID %d", got, existing.ID)
placeholder
	if !strings.HasPrefix(matchedKey, "access:") {
		t.Fatalf("matched key = %q, want access fingerprint", matchedKey)
placeholder
placeholder

func TestCodexAccountIndexKeepsAllCandidatesForSharedAccountKey(t *testing.T) {
	legacy := service.Account{
		ID: 30,
placeholder
			"chatgpt_account_id": "team-1",
			"access_token":       "token-legacy",
			"refresh_token":      "refresh-legacy",
	placeholder,
placeholder
	member := service.Account{
		ID: 31,
placeholder
			"chatgpt_account_id": "team-1",
			"chatgpt_user_id":    "user-2",
			"access_token":       "token-member",
			"refresh_token":      "refresh-member",
	placeholder,
placeholder

	// 无论索引构建顺序如何，携带新 user id 的条目都应跳过 user-2 的账号，
	// 命中缺少 user id 的存量账号，而不是因单一候选被遮蔽而落空。
	for _, accounts := range [][]service.Account{
		{member, legacyplaceholder,
		{legacy, memberplaceholder,
placeholder {
		index := buildCodexAccountIndex(accounts)

		keys := buildCodexImportIdentityKeys("team-1", "user-1", "", "token-new", "refresh-new")
		got, matchedKey := index.Find(keys, "user-1")
		if got == nil || got.ID != legacy.ID {
			t.Fatalf("Find with shared account key = %v, want legacy account ID %d", got, legacy.ID)
	placeholder
		if matchedKey != "account:team-1" {
			t.Fatalf("matched key = %q, want account:team-1", matchedKey)
	placeholder

		keys = buildCodexImportIdentityKeys("team-1", "user-2", "", "token-new", "refresh-new")
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
			"refresh_token":      "refresh-new",
	placeholder,
placeholder
	index.Add(backfilled)

	// 回填后同一账号在 account 键下应被原位替换而非残留旧副本：
	// 其他成员的条目不应再通过旧副本（无 user id）命中该账号。
	keys := buildCodexImportIdentityKeys("team-1", "user-2", "", "token-other", "refresh-other")
	if got, matchedKey := index.Find(keys, "user-2"); got != nil {
		t.Fatalf("stale candidate matched after upsert by %q: account ID %d", matchedKey, got.ID)
placeholder

	keys = buildCodexImportIdentityKeys("team-1", "user-1", "", "token-other", "refresh-other")
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
	member1 := buildCodexImportIdentityKeys("team-1", "user-1", "", "token-1", "refresh-1")
	markCodexIdentitySeen(seen, member1, 1, "user-1")

	member2 := buildCodexImportIdentityKeys("team-1", "user-2", "", "token-2", "refresh-2")
	if index, ok := firstSeenCodexIdentity(seen, member2, "user-2"); ok {
		t.Fatalf("different team member treated as duplicate of entry %d", index)
placeholder

	again := buildCodexImportIdentityKeys("team-1", "user-1", "", "token-3", "refresh-3")
	index, ok := firstSeenCodexIdentity(seen, again, "user-1")
	if !ok || index != 1 {
		t.Fatalf("same user re-entry dedup = (%d, %v), want (1, true)", index, ok)
placeholder

	// 无 user id 的条目不应因共享 account id 与已见团队成员互相去重；
	// 只有相同 access token 指纹才视为重复。
	opaque := buildCodexImportIdentityKeys("team-1", "", "", "token-4", "")
	index, ok = firstSeenCodexIdentity(seen, opaque, "")
	if ok {
		t.Fatalf("entry without user id dedup = (%d, %v), want no match", index, ok)
placeholder
placeholder

func TestNormalizeCodexImportUsesJWTSubForAccessTokenOnlyIdentity(t *testing.T) {
	accessToken := buildCodexImportTestJWT(t, time.Now().Add(time.Hour), map[string]any{
		"sub": "user-from-access-token",
		"https://api.openai.com/auth": map[string]any{
			"chatgpt_account_id": "workspace-1",
	placeholder,
placeholder)

	item, err := normalizeCodexImportEntry(codexImportEntry{Index: 1, Value: accessTokenplaceholder)
	if err != nil {
		t.Fatalf("normalizeCodexImportEntry error = %v", err)
placeholder
	if item.UserID != "user-from-access-token" {
		t.Fatalf("UserID = %q, want JWT sub", item.UserID)
placeholder
	if len(item.IdentityKeys) != 1 || !strings.HasPrefix(item.IdentityKeys[0], "access:") {
		t.Fatalf("IdentityKeys = %v, want access fingerprint only for accessToken-only import", item.IdentityKeys)
placeholder
	if got := item.Credentials["chatgpt_user_id"]; got != "user-from-access-token" {
		t.Fatalf("credential chatgpt_user_id = %v, want JWT sub", got)
placeholder
placeholder

func TestImportCodexSessionsAccessTokenOnlySameWorkspaceDifferentUsersCreatesTwoAccounts(t *testing.T) {
	svc := newCodexImportMemoryAdminService(nil)
	handler := NewAccountHandler(svc, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := CodexSessionImportRequest{SkipDefaultGroupBind: boolPtr(true)placeholder
	entries := []codexImportEntry{
		{Index: 1, Value: buildCodexAccessOnlyImportValue(t, "workspace-1", "user-1")placeholder,
		{Index: 2, Value: buildCodexAccessOnlyImportValue(t, "workspace-1", "user-2")placeholder,
placeholder

	result, err := handler.importCodexSessions(context.Background(), req, entries)
	if err != nil {
		t.Fatalf("importCodexSessions error = %v", err)
placeholder
	if result.Created != 2 || result.Updated != 0 || result.Skipped != 0 || result.Failed != 0 {
		t.Fatalf("result = %+v, want two created accounts", result)
placeholder
	if len(svc.createdAccounts) != 2 {
		t.Fatalf("created accounts = %d, want 2", len(svc.createdAccounts))
placeholder
	if svc.createdAccounts[0].Credentials["chatgpt_user_id"] == svc.createdAccounts[1].Credentials["chatgpt_user_id"] {
		t.Fatalf("created accounts share user id: %v", svc.createdAccounts)
placeholder
placeholder

func TestImportCodexSessionsAccessTokenOnlySameWorkspaceAndUserDifferentTokensCreatesTwoAccounts(t *testing.T) {
	svc := newCodexImportMemoryAdminService(nil)
	handler := NewAccountHandler(svc, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := CodexSessionImportRequest{SkipDefaultGroupBind: boolPtr(true)placeholder
	entries := []codexImportEntry{
		{Index: 1, Value: map[string]any{
			"access_token": buildCodexImportTestJWT(t, time.Now().Add(time.Hour), map[string]any{
				"sub": "shared-user",
				"jti": "token-1",
				"https://api.openai.com/auth": map[string]any{
					"chatgpt_account_id": "workspace-1",
			placeholder,
		placeholder),
placeholder
		{Index: 2, Value: map[string]any{
			"access_token": buildCodexImportTestJWT(t, time.Now().Add(time.Hour), map[string]any{
				"sub": "shared-user",
				"jti": "token-2",
				"https://api.openai.com/auth": map[string]any{
					"chatgpt_account_id": "workspace-1",
			placeholder,
		placeholder),
placeholder
placeholder

	result, err := handler.importCodexSessions(context.Background(), req, entries)
	if err != nil {
		t.Fatalf("importCodexSessions error = %v", err)
placeholder
	if result.Created != 2 || result.Updated != 0 || result.Skipped != 0 || result.Failed != 0 {
		t.Fatalf("result = %+v, want two created accounts", result)
placeholder
	if len(svc.createdAccounts) != 2 {
		t.Fatalf("created accounts = %d, want 2", len(svc.createdAccounts))
placeholder
placeholder

func TestImportCodexSessionsAccessTokenOnlySameUserUpdatesExisting(t *testing.T) {
	existingToken := buildCodexAccessToken(t, "workspace-1", "user-1", time.Now().Add(time.Hour))
	svc := newCodexImportMemoryAdminService([]service.Account{{
		ID:       10,
		Name:     "existing",
		Platform: service.PlatformOpenAI,
		Type:     service.AccountTypeOAuth,
placeholder
			"chatgpt_account_id": "workspace-1",
			"chatgpt_user_id":    "user-1",
			"access_token":       existingToken,
	placeholder,
		Extra: map[string]any{"openai_long_context_billing_enabled": falseplaceholder,
placeholderplaceholder)
	handler := NewAccountHandler(svc, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := CodexSessionImportRequest{SkipDefaultGroupBind: boolPtr(true)placeholder
	entries := []codexImportEntry{
		{Index: 1, Value: map[string]any{"access_token": existingTokenplaceholderplaceholder,
placeholder

	result, err := handler.importCodexSessions(context.Background(), req, entries)
	if err != nil {
		t.Fatalf("importCodexSessions error = %v", err)
placeholder
	if result.Created != 0 || result.Updated != 1 || result.Failed != 0 {
		t.Fatalf("result = %+v, want one updated account", result)
placeholder
	if len(svc.createdAccounts) != 0 {
		t.Fatalf("created accounts = %d, want 0", len(svc.createdAccounts))
placeholder
	if len(svc.updatedAccounts) != 1 || svc.updatedAccounts[0].id != 10 {
		t.Fatalf("updated accounts = %+v, want account 10", svc.updatedAccounts)
placeholder
	if got := svc.updatedAccounts[0].input.Extra["openai_long_context_billing_enabled"]; got != false {
		t.Fatalf("openai_long_context_billing_enabled = %v, want false", got)
placeholder
placeholder

func TestImportCodexSessionsUpgradesAccessTokenOnlyAccountWithRefreshToken(t *testing.T) {
	oldToken := buildCodexAccessTokenWithJTI(t, "workspace-1", "user-1", "old-token", time.Now().Add(time.Hour))
	newToken := buildCodexAccessTokenWithJTI(t, "workspace-1", "user-1", "new-token", time.Now().Add(time.Hour))
	svc := newCodexImportMemoryAdminService([]service.Account{{
		ID:       12,
		Name:     "existing",
		Platform: service.PlatformOpenAI,
		Type:     service.AccountTypeOAuth,
placeholder
			"chatgpt_account_id": "workspace-1",
			"chatgpt_user_id":    "user-1",
			"access_token":       oldToken,
	placeholder,
placeholderplaceholder)
	handler := NewAccountHandler(svc, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := CodexSessionImportRequest{SkipDefaultGroupBind: boolPtr(true)placeholder
	entries := []codexImportEntry{
		{Index: 1, Value: map[string]any{
			"access_token":  newToken,
			"refresh_token": "refresh-new",
placeholder
placeholder

	result, err := handler.importCodexSessions(context.Background(), req, entries)
	if err != nil {
		t.Fatalf("importCodexSessions error = %v", err)
placeholder
	if result.Created != 0 || result.Updated != 1 || result.Failed != 0 {
		t.Fatalf("result = %+v, want one updated account", result)
placeholder
	if len(svc.updatedAccounts) != 1 || svc.updatedAccounts[0].id != 12 {
		t.Fatalf("updated accounts = %+v, want account 12", svc.updatedAccounts)
placeholder
	if got := svc.updatedAccounts[0].input.Credentials["refresh_token"]; got != "refresh-new" {
		t.Fatalf("updated refresh_token = %v, want refresh-new", got)
placeholder
placeholder

func TestImportCodexSessionsAccessTokenOnlyPreservesExistingRefreshToken(t *testing.T) {
	existingToken := buildCodexAccessToken(t, "workspace-1", "user-1", time.Now().Add(time.Hour))
	svc := newCodexImportMemoryAdminService([]service.Account{{
		ID:       13,
		Name:     "existing",
		Platform: service.PlatformOpenAI,
		Type:     service.AccountTypeOAuth,
placeholder
			"chatgpt_account_id": "workspace-1",
			"chatgpt_user_id":    "user-1",
			"access_token":       existingToken,
			"refresh_token":      "refresh-old",
			"client_id":          "client-old",
	placeholder,
placeholderplaceholder)
	handler := NewAccountHandler(svc, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := CodexSessionImportRequest{SkipDefaultGroupBind: boolPtr(true)placeholder
	entries := []codexImportEntry{
		{Index: 1, Value: map[string]any{"access_token": existingTokenplaceholderplaceholder,
placeholder

	result, err := handler.importCodexSessions(context.Background(), req, entries)
	if err != nil {
		t.Fatalf("importCodexSessions error = %v", err)
placeholder
	if result.Created != 0 || result.Updated != 1 || result.Failed != 0 {
		t.Fatalf("result = %+v, want one updated account", result)
placeholder
	update := svc.updatedAccounts[0].input
	if got := update.Credentials["refresh_token"]; got != "refresh-old" {
		t.Fatalf("refresh_token = %v, want refresh-old", got)
placeholder
	if got := update.Credentials["client_id"]; got != "client-old" {
		t.Fatalf("client_id = %v, want client-old", got)
placeholder
	if update.ExpiresAt != nil {
		t.Fatalf("ExpiresAt = %v, want nil to preserve OAuth account expiry", *update.ExpiresAt)
placeholder
	if update.AutoPauseOnExpired != nil {
		t.Fatalf("AutoPauseOnExpired = %v, want nil to preserve OAuth account scheduling", *update.AutoPauseOnExpired)
placeholder
placeholder

func TestImportCodexSessionsBatchOldAccessTokenDoesNotRollbackRefreshToken(t *testing.T) {
	oldToken := buildCodexAccessTokenWithJTI(t, "workspace-1", "user-1", "old-token", time.Now().Add(time.Hour))
	newToken := buildCodexAccessTokenWithJTI(t, "workspace-1", "user-1", "new-token", time.Now().Add(time.Hour))
	svc := newCodexImportMemoryAdminService([]service.Account{{
		ID:       14,
		Name:     "existing",
		Platform: service.PlatformOpenAI,
		Type:     service.AccountTypeOAuth,
placeholder
			"chatgpt_account_id": "workspace-1",
			"chatgpt_user_id":    "user-1",
			"access_token":       oldToken,
			"refresh_token":      "refresh-old",
	placeholder,
placeholderplaceholder)
	handler := NewAccountHandler(svc, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := CodexSessionImportRequest{SkipDefaultGroupBind: boolPtr(true)placeholder
	entries := []codexImportEntry{
		{Index: 1, Value: map[string]any{
			"access_token":  newToken,
			"refresh_token": "refresh-new",
placeholder
		{Index: 2, Value: map[string]any{"access_token": oldTokenplaceholderplaceholder,
placeholder

	result, err := handler.importCodexSessions(context.Background(), req, entries)
	if err != nil {
		t.Fatalf("importCodexSessions error = %v", err)
placeholder
	if result.Updated != 1 || result.Created != 1 || result.Failed != 0 {
		t.Fatalf("result = %+v, want first item updated and stale access token created separately", result)
placeholder
	if len(svc.updatedAccounts) != 1 || svc.updatedAccounts[0].id != 14 {
		t.Fatalf("updated accounts = %+v, want account 14 updated once", svc.updatedAccounts)
placeholder
	stored, err := svc.GetAccount(context.Background(), 14)
	if err != nil {
		t.Fatalf("GetAccount error = %v", err)
placeholder
	if got := stored.Credentials["access_token"]; got != newToken {
		t.Fatalf("stored access_token rolled back = %v, want new token", got)
placeholder
	if got := stored.Credentials["refresh_token"]; got != "refresh-new" {
		t.Fatalf("stored refresh_token = %v, want refresh-new", got)
placeholder
placeholder

func TestImportCodexSessionsWithRefreshTokenKeepsExistingDedup(t *testing.T) {
	existingToken := buildCodexAccessToken(t, "workspace-1", "user-1", time.Now().Add(time.Hour))
	svc := newCodexImportMemoryAdminService([]service.Account{{
		ID:       11,
		Name:     "existing",
		Platform: service.PlatformOpenAI,
		Type:     service.AccountTypeOAuth,
placeholder
			"chatgpt_account_id": "workspace-1",
			"chatgpt_user_id":    "user-1",
			"access_token":       existingToken,
			"refresh_token":      "refresh-old",
	placeholder,
placeholderplaceholder)
	handler := NewAccountHandler(svc, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := CodexSessionImportRequest{SkipDefaultGroupBind: boolPtr(true)placeholder
	entries := []codexImportEntry{
		{Index: 1, Value: buildCodexRefreshImportValue(t, "workspace-1", "user-1", "refresh-new")placeholder,
placeholder

	result, err := handler.importCodexSessions(context.Background(), req, entries)
	if err != nil {
		t.Fatalf("importCodexSessions error = %v", err)
placeholder
	if result.Created != 0 || result.Updated != 1 || result.Failed != 0 {
		t.Fatalf("result = %+v, want one updated account", result)
placeholder
	if got := svc.updatedAccounts[0].input.Credentials["refresh_token"]; got != "refresh-new" {
		t.Fatalf("updated refresh_token = %v, want refresh-new", got)
placeholder
placeholder

type codexImportMemoryAdminService struct {
	*stubAdminService
	nextID          int64
	updatedAccounts []struct {
		id    int64
		input *service.UpdateAccountInput
placeholder
placeholder

func newCodexImportMemoryAdminService(accounts []service.Account) *codexImportMemoryAdminService {
	stub := newStubAdminService()
	stub.accounts = append([]service.Account(nil), accounts...)
	return &codexImportMemoryAdminService{
		stubAdminService: stub,
		nextID:           100,
placeholder
placeholder

func (s *codexImportMemoryAdminService) CreateAccount(ctx context.Context, input *service.CreateAccountInput) (*service.Account, error) {
	s.createdAccounts = append(s.createdAccounts, input)
	if s.createAccountErr != nil {
		return nil, s.createAccountErr
placeholder
	account := service.Account{
		ID:          s.nextID,
		Name:        input.Name,
		Platform:    input.Platform,
		Type:        input.Type,
		Status:      service.StatusActive,
		Credentials: cloneCodexImportTestMap(input.Credentials),
		Extra:       cloneCodexImportTestMap(input.Extra),
placeholder
	s.nextID++
	s.accounts = append(s.accounts, account)
	return &account, nil
placeholder

func (s *codexImportMemoryAdminService) UpdateAccount(ctx context.Context, id int64, input *service.UpdateAccountInput) (*service.Account, error) {
	s.updatedAccounts = append(s.updatedAccounts, struct {
		id    int64
		input *service.UpdateAccountInput
placeholder{id: id, input: inputplaceholder)
	if s.updateAccountErr != nil {
		return nil, s.updateAccountErr
placeholder
	for idx := range s.accounts {
		if s.accounts[idx].ID == id {
			s.accounts[idx].Credentials = cloneCodexImportTestMap(input.Credentials)
			s.accounts[idx].Extra = cloneCodexImportTestMap(input.Extra)
			return &s.accounts[idx], nil
	placeholder
placeholder
	account := service.Account{ID: id, Status: service.StatusActive, Credentials: cloneCodexImportTestMap(input.Credentials)placeholder
	return &account, nil
placeholder

func (s *codexImportMemoryAdminService) GetAccount(ctx context.Context, id int64) (*service.Account, error) {
	for idx := range s.accounts {
		if s.accounts[idx].ID == id {
			return &s.accounts[idx], nil
	placeholder
placeholder
	return s.stubAdminService.GetAccount(ctx, id)
placeholder

func buildCodexAccessOnlyImportValue(t *testing.T, accountID, userID string) map[string]any {
placeholder
	return map[string]any{
		"access_token": buildCodexAccessToken(t, accountID, userID, time.Now().Add(time.Hour)),
placeholder
placeholder

func buildCodexRefreshImportValue(t *testing.T, accountID, userID, refreshToken string) map[string]any {
placeholder
	return map[string]any{
		"access_token":  buildCodexAccessToken(t, accountID, userID, time.Now().Add(time.Hour)),
		"refresh_token": refreshToken,
placeholder
placeholder

func buildCodexAccessToken(t *testing.T, accountID, userID string, exp time.Time) string {
placeholder
	return buildCodexAccessTokenWithJTI(t, accountID, userID, "", exp)
placeholder

func buildCodexAccessTokenWithJTI(t *testing.T, accountID, userID, jti string, exp time.Time) string {
placeholder
	claims := map[string]any{
		"sub": userID,
		"https://api.openai.com/auth": map[string]any{
			"chatgpt_account_id": accountID,
	placeholder,
placeholder
	if jti != "" {
		claims["jti"] = jti
placeholder
	return buildCodexImportTestJWT(t, exp, claims)
placeholder

func cloneCodexImportTestMap(input map[string]any) map[string]any {
	if input == nil {
		return nil
placeholder
	out := make(map[string]any, len(input))
	for key, value := range input {
		out[key] = value
placeholder
	return out
placeholder

func boolPtr(v bool) *bool {
	return &v
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
