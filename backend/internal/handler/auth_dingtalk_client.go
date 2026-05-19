package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// dingTalkClientConfig 是 DingTalkClient 需要的最小配置子集
type dingTalkClientConfig struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
	UserInfoURL  string
placeholder

type DingTalkClient struct {
	cfg         dingTalkClientConfig
	appToken    string
	appTokenExp time.Time // 钉钉 7200s，留 200s 余量 → 7000s
	mu          sync.Mutex
	httpClient  *http.Client
	// TODO(multi-instance): Redis 集中缓存 appToken
placeholder

type DingTalkUserTokenResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpireIn     int64  `json:"expireIn"`
	CorpID       string `json:"corpId"`
placeholder

func (c *DingTalkClient) ExchangeCodeForUserToken(ctx context.Context, code string) (*DingTalkUserTokenResp, error) {
	body := map[string]string{
		"clientId":     c.cfg.ClientID,
		"clientSecret": c.cfg.ClientSecret,
		"code":         code,
		"grantType":    "authorization_code",
placeholder
	payload, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.TokenURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, parseDingTalkErr(raw, resp.StatusCode)
placeholder
	var out DingTalkUserTokenResp
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
placeholder
	if strings.TrimSpace(out.AccessToken) == "" {
		return nil, parseDingTalkErr(raw, resp.StatusCode)
placeholder
	return &out, nil
placeholder

type DingTalkAPIError struct {
	Code    string
	Message string
	HTTP    int
placeholder

func (e *DingTalkAPIError) Error() string {
	return fmt.Sprintf("dingtalk api error code=%s msg=%s http=%d", e.Code, e.Message, e.HTTP)
placeholder

func parseDingTalkErr(raw []byte, status int) error {
	var v struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
placeholder
	_ = json.Unmarshal(raw, &v)
	code := v.Code
	if code == "" && v.ErrCode != 0 {
		code = fmt.Sprintf("%d", v.ErrCode)
placeholder
	msg := v.Message
	if msg == "" {
		msg = v.ErrMsg
placeholder
	return &DingTalkAPIError{Code: code, Message: msg, HTTP: statusplaceholder
placeholder

// GetUnionIdByUserToken 调用 /v1.0/contact/users/me 返回 unionId 与用户自设昵称 nick。
// nick 来自钉钉新版 OIDC 接口（用户在 App 个人资料填的昵称），与旧版 user/get.nickname 不同源。
func (c *DingTalkClient) GetUnionIdByUserToken(ctx context.Context, userToken string) (unionID string, nick string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.cfg.UserInfoURL, nil)
	if err != nil {
		return "", "", err
placeholder
	req.Header.Set("x-acs-dingtalk-access-token", userToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", "", parseDingTalkErr(raw, resp.StatusCode)
placeholder
	var v struct {
		UnionID string `json:"unionId"`
		Nick    string `json:"nick"`
placeholder
	if err := json.Unmarshal(raw, &v); err != nil {
		return "", "", err
placeholder
	if strings.TrimSpace(v.UnionID) == "" {
		return "", "", parseDingTalkErr(raw, resp.StatusCode)
placeholder
	return v.UnionID, v.Nick, nil
placeholder

type DingTalkStaffInfo struct {
	UserID   string
	Name     string // 企业内真实姓名（钉钉企业管理后台配置）
	Nickname string // 钉钉个人昵称（用户自己设置）
	Email    string
	DeptIDs  []int64
	// CorpID 不来自 staff 接口，来自 userToken；不在此 struct
placeholder

// dingTalkOAPIBase 推导钉钉旧版 OAPI base URL（host: api.dingtalk.com → oapi.dingtalk.com）。
// getbyunionid 与 topapi/v2/user/get 仅在旧版 OAPI 提供，不在 v1.0 OpenAPI。
func (c *DingTalkClient) dingTalkOAPIBase() string {
	u, err := url.Parse(c.cfg.UserInfoURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "https://oapi.dingtalk.com"
placeholder
	host := u.Host
	if strings.HasPrefix(host, "api.") {
		host = "oapi." + strings.TrimPrefix(host, "api.")
placeholder
	return u.Scheme + "://" + host
placeholder

func (c *DingTalkClient) GetAppToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.appToken != "" && time.Now().Before(c.appTokenExp) {
		return c.appToken, nil
placeholder
	body := map[string]string{"appKey": c.cfg.ClientID, "appSecret": c.cfg.ClientSecretplaceholder
	payload, _ := json.Marshal(body)
	// 钉钉新版 v1.0 企业内部应用 access_token: POST /v1.0/oauth2/accessToken
	// 此 token 也可作为旧版 OAPI 的 access_token 使用（钉钉文档已说明）
	appTokenURL := strings.Replace(c.cfg.TokenURL, "/oauth2/userAccessToken", "/oauth2/accessToken", 1)
	if !strings.Contains(appTokenURL, "accessToken") && !strings.Contains(appTokenURL, "gettoken") {
		appTokenURL = c.cfg.TokenURL // fallback for test stub
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, appTokenURL, bytes.NewReader(payload))
	if err != nil {
		return "", err
placeholder
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", parseDingTalkErr(raw, resp.StatusCode)
placeholder
	var v struct {
		AccessToken string `json:"accessToken"`
		ExpireIn    int64  `json:"expireIn"`
placeholder
	if err := json.Unmarshal(raw, &v); err != nil {
		return "", err
placeholder
	if v.AccessToken == "" {
		return "", parseDingTalkErr(raw, resp.StatusCode)
placeholder
	c.appToken = v.AccessToken
	ttl := v.ExpireIn
	if ttl > 200 {
		ttl -= 200
placeholder
	c.appTokenExp = time.Now().Add(time.Duration(ttl) * time.Second)
	return c.appToken, nil
placeholder

func (c *DingTalkClient) GetUserIdByUnionId(ctx context.Context, unionID string) (string, error) {
	appToken, err := c.GetAppToken(ctx)
	if err != nil {
		return "", err
placeholder
	body := map[string]string{"unionid": unionIDplaceholder
	payload, _ := json.Marshal(body)
	// 钉钉旧版 OAPI: POST https://oapi.dingtalk.com/topapi/user/getbyunionid?access_token=XXX
	// access_token 通过 query string 传递（不是 header）
	var targetURL string
	if strings.Contains(c.cfg.UserInfoURL, "/contact/users/me") {
		targetURL = c.dingTalkOAPIBase() + "/topapi/user/getbyunionid?access_token=" + url.QueryEscape(appToken)
placeholder else {
		targetURL = c.cfg.UserInfoURL // fallback for test stub
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(payload))
	if err != nil {
		return "", err
placeholder
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", parseDingTalkErr(raw, resp.StatusCode)
placeholder
	var v struct {
		Result struct {
			UserID string `json:"userid"`
	placeholder `json:"result"`
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
placeholder
	if err := json.Unmarshal(raw, &v); err != nil {
		return "", err
placeholder
	if v.ErrCode != 0 {
		return "", parseDingTalkErr(raw, resp.StatusCode)
placeholder
	if strings.TrimSpace(v.Result.UserID) == "" {
		return "", parseDingTalkErr(raw, resp.StatusCode)
placeholder
	return v.Result.UserID, nil
placeholder

// DingTalkDeptInfo 部门信息（topapi/v2/department/get 返回子集）
type DingTalkDeptInfo struct {
	DeptID   int64
	Name     string
	ParentID int64
placeholder

// GetDeptInfo 查询单个部门信息（用于递归拼部门路径）。
// 调用钉钉旧版 OAPI: POST /topapi/v2/department/get?access_token=XXX
func (c *DingTalkClient) GetDeptInfo(ctx context.Context, deptID int64) (*DingTalkDeptInfo, error) {
	appToken, err := c.GetAppToken(ctx)
	if err != nil {
		return nil, err
placeholder
	body := map[string]any{"dept_id": deptID, "language": "zh_CN"placeholder
	payload, _ := json.Marshal(body)
	var targetURL string
	if strings.Contains(c.cfg.UserInfoURL, "/contact/users/me") {
		targetURL = c.dingTalkOAPIBase() + "/topapi/v2/department/get?access_token=" + url.QueryEscape(appToken)
placeholder else {
		targetURL = c.cfg.UserInfoURL // test stub fallback
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, parseDingTalkErr(raw, resp.StatusCode)
placeholder
	var v struct {
		Result struct {
			DeptID   int64  `json:"dept_id"`
			Name     string `json:"name"`
			ParentID int64  `json:"parent_id"`
	placeholder `json:"result"`
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
placeholder
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, err
placeholder
	if v.ErrCode != 0 {
		return nil, parseDingTalkErr(raw, resp.StatusCode)
placeholder
	return &DingTalkDeptInfo{
		DeptID:   v.Result.DeptID,
		Name:     v.Result.Name,
		ParentID: v.Result.ParentID,
placeholder, nil
placeholder

func (c *DingTalkClient) GetStaffInfoByUserId(ctx context.Context, userID string) (*DingTalkStaffInfo, error) {
	appToken, err := c.GetAppToken(ctx)
	if err != nil {
		return nil, err
placeholder
	body := map[string]string{"userid": userIDplaceholder
	payload, _ := json.Marshal(body)
	// 钉钉旧版 OAPI: POST https://oapi.dingtalk.com/topapi/v2/user/get?access_token=XXX
	var targetURL string
	if strings.Contains(c.cfg.UserInfoURL, "/contact/users/me") {
		targetURL = c.dingTalkOAPIBase() + "/topapi/v2/user/get?access_token=" + url.QueryEscape(appToken)
placeholder else {
		targetURL = c.cfg.UserInfoURL // fallback for test stub
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, parseDingTalkErr(raw, resp.StatusCode)
placeholder
	var v struct {
		Result struct {
			UserID    string  `json:"userid"`
			Name      string  `json:"name"`
			Nickname  string  `json:"nickname"`
			Email     string  `json:"email"`
			OrgEmail  string  `json:"org_email"`
			Extension string  `json:"extension"`
			DeptID    []int64 `json:"dept_id_list"`
	placeholder `json:"result"`
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
placeholder
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, err
placeholder
	if v.ErrCode != 0 {
		return nil, parseDingTalkErr(raw, resp.StatusCode)
placeholder
	if strings.TrimSpace(v.Result.UserID) == "" {
		return nil, parseDingTalkErr(raw, resp.StatusCode)
placeholder
	// 邮箱三级 fallback：org_email > email > extension["企业邮箱"]（钉钉自定义扩展字段，JSON string）
	email := strings.TrimSpace(v.Result.OrgEmail)
	emailSource := "org_email"
	if email == "" {
		email = strings.TrimSpace(v.Result.Email)
		emailSource = "email"
placeholder
	extensionParsed := false
	if email == "" && strings.TrimSpace(v.Result.Extension) != "" {
		var ext map[string]string
		if err := json.Unmarshal([]byte(v.Result.Extension), &ext); err == nil {
			extensionParsed = true
			if v, ok := ext["企业邮箱"]; ok {
				email = strings.TrimSpace(v)
				emailSource = "extension.企业邮箱"
		placeholder
	placeholder
placeholder
	if email == "" {
		emailSource = "none"
placeholder
	slog.Info("dingtalk staff fetched",
		"userid", v.Result.UserID,
		"name_present", v.Result.Name != "",
		"nickname_present", v.Result.Nickname != "",
		"name_eq_nickname", v.Result.Name != "" && v.Result.Name == v.Result.Nickname,
		"email_present", v.Result.Email != "",
		"org_email_present", v.Result.OrgEmail != "",
		"extension_present", v.Result.Extension != "",
		"extension_parsed", extensionParsed,
		"email_source", emailSource,
		"dept_count", len(v.Result.DeptID),
	)
	return &DingTalkStaffInfo{
		UserID:   v.Result.UserID,
		Name:     v.Result.Name,
		Nickname: v.Result.Nickname,
		Email:    email,
		DeptIDs:  v.Result.DeptID,
placeholder, nil
placeholder
