// Package provider contains concrete payment provider implementations.
package provider

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Wei-Shaw/sub2api/internal/payment"
)

// EasyPay constants.
const (
	easypayCodeSuccess     = 1
	easypayStatusPaid      = 1
	easypayHTTPTimeout     = 10 * time.Second
	maxEasypayResponseSize = 1 << 20 // 1MB
	maxEasypayErrorSummary = 512
	tradeStatusSuccess     = "TRADE_SUCCESS"
	signTypeMD5            = "MD5"
	paymentModePopup       = "popup"
	deviceMobile           = "mobile"
)

// EasyPay implements payment.Provider for the EasyPay aggregation platform.
type EasyPay struct {
	instanceID string
	config     map[string]string
	httpClient *http.Client
placeholder

type easyPayCustomMethod struct {
	Type         string `json:"type"`
	UpstreamType string `json:"upstreamType"`
	DisplayName  string `json:"displayName"`
placeholder

// NewEasyPay creates a new EasyPay provider.
// config keys: pid, pkey, apiBase, notifyUrl, returnUrl, cid, cidAlipay, cidWxpay
func NewEasyPay(instanceID string, config map[string]string) (*EasyPay, error) {
	for _, k := range []string{"pid", "pkey", "apiBase", "notifyUrl", "returnUrl"placeholder {
		if strings.TrimSpace(config[k]) == "" {
			return nil, fmt.Errorf("easypay config missing required key: %s", k)
	placeholder
placeholder
	cfg := make(map[string]string, len(config))
	for k, v := range config {
		cfg[k] = v
placeholder
	cfg["apiBase"] = normalizeEasyPayAPIBase(cfg["apiBase"])
	return &EasyPay{
		instanceID: instanceID,
		config:     cfg,
		httpClient: &http.Client{Timeout: easypayHTTPTimeoutplaceholder,
placeholder, nil
placeholder

func normalizeEasyPayAPIBase(apiBase string) string {
	base := strings.TrimSpace(apiBase)
	if base == "" {
		return ""
placeholder
	if parsed, err := url.Parse(base); err == nil && parsed.Scheme != "" && parsed.Host != "" {
		parsed.RawQuery = ""
		parsed.Fragment = ""
		parsed.RawPath = ""
		parsed.Path = trimEasyPayEndpointPath(parsed.Path)
		return strings.TrimRight(parsed.String(), "/")
placeholder
	return strings.TrimRight(trimEasyPayEndpointPath(base), "/")
placeholder

func trimEasyPayEndpointPath(path string) string {
	path = strings.TrimRight(strings.TrimSpace(path), "/")
	lower := strings.ToLower(path)
	for _, endpoint := range []string{"/submit.php", "/mapi.php", "/api.php"placeholder {
		if strings.HasSuffix(lower, endpoint) {
			return strings.TrimRight(path[:len(path)-len(endpoint)], "/")
	placeholder
placeholder
	return path
placeholder

func (e *EasyPay) apiBase() string {
	if e == nil {
		return ""
placeholder
	return normalizeEasyPayAPIBase(e.config["apiBase"])
placeholder

func (e *EasyPay) Name() string        { return "EasyPay" placeholder
func (e *EasyPay) ProviderKey() string { return payment.TypeEasyPay placeholder
func (e *EasyPay) SupportedTypes() []payment.PaymentType {
	types := []payment.PaymentType{payment.TypeAlipay, payment.TypeWxpayplaceholder
	for _, method := range e.customMethods() {
		if method.Type != "" {
			types = append(types, method.Type)
	placeholder
placeholder
	return types
placeholder

func (e *EasyPay) MerchantIdentityMetadata() map[string]string {
	if e == nil {
		return nil
placeholder
	pid := strings.TrimSpace(e.config["pid"])
	if pid == "" {
		return nil
placeholder
	return map[string]string{"pid": pidplaceholder
placeholder

func (e *EasyPay) CreatePayment(ctx context.Context, req payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	// Payment mode determined by instance config, not payment type.
	// "popup" → hosted page (submit.php); "qrcode"/default → API call (mapi.php).
	mode := e.config["paymentMode"]
	if mode == paymentModePopup {
		return e.createRedirectPayment(req)
placeholder
	return e.createAPIPayment(ctx, req)
placeholder

// createRedirectPayment builds a submit.php URL for browser redirect.
// No server-side API call — the user is redirected to EasyPay's hosted page.
// TradeNo is empty; it arrives via the notify callback after payment.
func (e *EasyPay) createRedirectPayment(req payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	notifyURL, returnURL := e.resolveURLs(req)
	paymentType := e.upstreamPaymentType(req.PaymentType)
	params := map[string]string{
		"pid": e.config["pid"], "type": paymentType,
		"out_trade_no": req.OrderID, "notify_url": notifyURL,
		"return_url": returnURL, "name": req.Subject,
		"money": req.Amount,
placeholder
	if cid := e.resolveCID(paymentType); cid != "" {
		params["cid"] = cid
placeholder
	if req.IsMobile {
		params["device"] = deviceMobile
placeholder
	params["sign"] = easyPaySign(params, e.config["pkey"])
	params["sign_type"] = signTypeMD5

	q := url.Values{placeholder
	for k, v := range params {
		q.Set(k, v)
placeholder
	payURL := e.apiBase() + "/submit.php?" + q.Encode()
	return &payment.CreatePaymentResponse{PayURL: payURLplaceholder, nil
placeholder

// createAPIPayment calls mapi.php to get payurl/qrcode (existing behavior).
func (e *EasyPay) createAPIPayment(ctx context.Context, req payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	notifyURL, returnURL := e.resolveURLs(req)
	paymentType := e.upstreamPaymentType(req.PaymentType)
	params := map[string]string{
		"pid": e.config["pid"], "type": paymentType,
		"out_trade_no": req.OrderID, "notify_url": notifyURL,
		"return_url": returnURL, "name": req.Subject,
		"money": req.Amount, "clientip": req.ClientIP,
placeholder
	if cid := e.resolveCID(paymentType); cid != "" {
		params["cid"] = cid
placeholder
	if req.IsMobile {
		params["device"] = deviceMobile
placeholder
	params["sign"] = easyPaySign(params, e.config["pkey"])
	params["sign_type"] = signTypeMD5

	body, err := e.post(ctx, e.apiBase()+"/mapi.php", params)
	if err != nil {
		return nil, fmt.Errorf("easypay create: %w", err)
placeholder
	var resp struct {
		Code    int    `json:"code"`
		Msg     string `json:"msg"`
		TradeNo string `json:"trade_no"`
		PayURL  string `json:"payurl"`
		PayURL2 string `json:"payurl2"` // H5 mobile payment URL
		QRCode  string `json:"qrcode"`
placeholder
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("easypay parse: %w", err)
placeholder
	if resp.Code != easypayCodeSuccess {
		return nil, fmt.Errorf("easypay error: %s", resp.Msg)
placeholder
	payURL := resp.PayURL
	if req.IsMobile && resp.PayURL2 != "" {
		payURL = resp.PayURL2
placeholder
	return &payment.CreatePaymentResponse{TradeNo: resp.TradeNo, PayURL: payURL, QRCode: resp.QRCodeplaceholder, nil
placeholder

// resolveURLs returns (notifyURL, returnURL) preferring request values,
// falling back to instance config.
func (e *EasyPay) resolveURLs(req payment.CreatePaymentRequest) (string, string) {
	notifyURL := req.NotifyURL
	if notifyURL == "" {
		notifyURL = e.config["notifyUrl"]
placeholder
	returnURL := req.ReturnURL
	if returnURL == "" {
		returnURL = e.config["returnUrl"]
placeholder
	return notifyURL, returnURL
placeholder

func (e *EasyPay) customMethods() []easyPayCustomMethod {
	if e == nil {
		return nil
placeholder
	raw := strings.TrimSpace(e.config["customMethods"])
	if raw == "" {
		return nil
placeholder
	var methods []easyPayCustomMethod
	if err := json.Unmarshal([]byte(raw), &methods); err != nil {
		return nil
placeholder
	result := make([]easyPayCustomMethod, 0, len(methods))
	for _, method := range methods {
		method.Type = strings.TrimSpace(method.Type)
		method.UpstreamType = strings.TrimSpace(method.UpstreamType)
		method.DisplayName = strings.TrimSpace(method.DisplayName)
		if method.Type == "" || method.UpstreamType == "" {
			continue
	placeholder
		result = append(result, method)
placeholder
	return result
placeholder

func (e *EasyPay) upstreamPaymentType(paymentType string) string {
	paymentType = strings.TrimSpace(paymentType)
	for _, method := range e.customMethods() {
		if paymentType == method.Type {
			return method.UpstreamType
	placeholder
placeholder
	return paymentType
placeholder

func (e *EasyPay) QueryOrder(ctx context.Context, tradeNo string) (*payment.QueryOrderResponse, error) {
	params := map[string]string{
		"act": "order", "pid": e.config["pid"],
		"key": e.config["pkey"], "out_trade_no": tradeNo,
placeholder
	body, err := e.post(ctx, e.apiBase()+"/api.php", params)
	if err != nil {
		return nil, fmt.Errorf("easypay query: %w", err)
placeholder
	type easyPayQueryData struct {
		TradeStatus *string `json:"trade_status"`
		Status      *int    `json:"status"`
		Money       *string `json:"money"`
		TradeNo     *string `json:"trade_no"`
placeholder
	var resp struct {
		Code        int              `json:"code"`
		Msg         string           `json:"msg"`
		TradeStatus *string          `json:"trade_status"`
		Status      *int             `json:"status"`
		Money       *string          `json:"money"`
		TradeNo     *string          `json:"trade_no"`
		Data        easyPayQueryData `json:"data"`
placeholder
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("easypay parse query: %w", err)
placeholder
	status := payment.ProviderStatusPending
	if resp.TradeStatus != nil {
		if *resp.TradeStatus == tradeStatusSuccess {
			status = payment.ProviderStatusPaid
	placeholder
placeholder else if resp.Data.TradeStatus != nil {
		if *resp.Data.TradeStatus == tradeStatusSuccess {
			status = payment.ProviderStatusPaid
	placeholder
placeholder else if resp.Status != nil {
		if *resp.Status == easypayStatusPaid {
			status = payment.ProviderStatusPaid
	placeholder
placeholder else if resp.Data.Status != nil && *resp.Data.Status == easypayStatusPaid {
		status = payment.ProviderStatusPaid
placeholder

	money := ""
	if resp.Money != nil {
		money = *resp.Money
placeholder else if resp.Data.Money != nil {
		money = *resp.Data.Money
placeholder
	responseTradeNo := tradeNo
	if resp.TradeNo != nil {
		if *resp.TradeNo != "" {
			responseTradeNo = *resp.TradeNo
	placeholder
placeholder else if resp.Data.TradeNo != nil && *resp.Data.TradeNo != "" {
		responseTradeNo = *resp.Data.TradeNo
placeholder

	amount, _ := strconv.ParseFloat(money, 64)
	return &payment.QueryOrderResponse{
		TradeNo:  responseTradeNo,
		Status:   status,
		Amount:   amount,
		Metadata: e.MerchantIdentityMetadata(),
placeholder, nil
placeholder

func (e *EasyPay) VerifyNotification(_ context.Context, rawBody string, _ map[string]string) (*payment.PaymentNotification, error) {
	values, err := url.ParseQuery(rawBody)
	if err != nil {
		return nil, fmt.Errorf("parse notify: %w", err)
placeholder
	// url.ParseQuery already decodes values — no additional decode needed.
	params := make(map[string]string)
	for k := range values {
		params[k] = values.Get(k)
placeholder
	sign := params["sign"]
	if sign == "" {
		return nil, fmt.Errorf("missing sign")
placeholder
	if !easyPayVerifySign(params, e.config["pkey"], sign) {
		return nil, fmt.Errorf("invalid signature")
placeholder
	status := payment.ProviderStatusFailed
	if params["trade_status"] == tradeStatusSuccess {
		status = payment.ProviderStatusSuccess
placeholder
	amount, _ := strconv.ParseFloat(params["money"], 64)

	metadata := e.MerchantIdentityMetadata()
	if pid := strings.TrimSpace(params["pid"]); pid != "" {
		if metadata == nil {
			metadata = map[string]string{placeholder
	placeholder
		metadata["pid"] = pid
placeholder
	return &payment.PaymentNotification{
		TradeNo: params["trade_no"], OrderID: params["out_trade_no"],
		Amount: amount, Status: status, RawData: rawBody, Metadata: metadata,
placeholder, nil
placeholder

func (e *EasyPay) Refund(ctx context.Context, req payment.RefundRequest) (*payment.RefundResponse, error) {
	attempts := e.refundAttempts(req)
	if len(attempts) == 0 {
		return nil, fmt.Errorf("easypay refund missing order identifier")
placeholder
	var firstErr error
	for i, attempt := range attempts {
		body, status, err := e.postRaw(ctx, e.apiBase()+"/api.php?act=refund", attempt.params)
		if err != nil {
			return nil, fmt.Errorf("easypay refund request: %w", err)
	placeholder
		if err := parseEasyPayRefundResponse(status, body); err != nil {
			if firstErr == nil {
				firstErr = err
		placeholder
			if i+1 < len(attempts) && isEasyPayRefundOrderNotFound(err) {
				continue
		placeholder
			return nil, err
	placeholder
		return &payment.RefundResponse{RefundID: attempt.refundID, Status: payment.ProviderStatusSuccessplaceholder, nil
placeholder
	return nil, firstErr
placeholder

type easyPayRefundAttempt struct {
	params   map[string]string
	refundID string
placeholder

func (e *EasyPay) refundAttempts(req payment.RefundRequest) []easyPayRefundAttempt {
	base := map[string]string{
		"pid": e.config["pid"], "key": e.config["pkey"], "money": req.Amount,
placeholder
	var attempts []easyPayRefundAttempt
	if orderID := strings.TrimSpace(req.OrderID); orderID != "" {
		params := cloneStringMap(base)
		params["out_trade_no"] = orderID
		attempts = append(attempts, easyPayRefundAttempt{params: params, refundID: orderIDplaceholder)
placeholder
	if tradeNo := strings.TrimSpace(req.TradeNo); tradeNo != "" {
		params := cloneStringMap(base)
		params["trade_no"] = tradeNo
		attempts = append(attempts, easyPayRefundAttempt{params: params, refundID: tradeNoplaceholder)
placeholder
	return attempts
placeholder

func cloneStringMap(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
placeholder
	return out
placeholder

func isEasyPayRefundOrderNotFound(err error) bool {
	if err == nil {
		return false
placeholder
	msg := err.Error()
	lower := strings.ToLower(msg)
	return strings.Contains(msg, "订单编号不存在") ||
		strings.Contains(msg, "订单不存在") ||
		strings.Contains(lower, "order not found") ||
		strings.Contains(lower, "not exist")
placeholder

func parseEasyPayRefundResponse(status int, body []byte) error {
	summary := summarizeEasyPayResponse(body)
	if status < http.StatusOK || status >= http.StatusMultipleChoices {
		return fmt.Errorf("easypay refund HTTP %d: %s", status, summary)
placeholder

	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" {
		return fmt.Errorf("easypay refund empty response (HTTP %d): %s", status, summary)
placeholder

	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "<!doctype html") || strings.HasPrefix(lower, "<html") ||
		(strings.HasPrefix(lower, "<") && strings.Contains(lower, "html")) {
		return fmt.Errorf("easypay refund non-JSON response (HTTP %d): %s", status, summary)
placeholder

	var resp struct {
		Code any    `json:"code"`
		Msg  string `json:"msg"`
placeholder
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("easypay refund non-JSON response (HTTP %d): %s", status, summary)
placeholder
	if !easyPayResponseCodeIsSuccess(resp.Code) {
		msg := strings.TrimSpace(resp.Msg)
		if msg == "" {
			msg = summary
	placeholder
		return fmt.Errorf("easypay refund failed (HTTP %d): %s", status, msg)
placeholder
	return nil
placeholder

func easyPayResponseCodeIsSuccess(code any) bool {
	switch v := code.(type) {
	case float64:
		return int(v) == easypayCodeSuccess
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(v))
		return err == nil && n == easypayCodeSuccess
	default:
		return false
placeholder
placeholder

func summarizeEasyPayResponse(body []byte) string {
	summary := strings.Join(strings.Fields(string(body)), " ")
	if summary == "" {
		return "<empty>"
placeholder
	if len(summary) > maxEasypayErrorSummary {
		truncated := summary[:maxEasypayErrorSummary]
		for len(truncated) > 0 && !utf8.ValidString(truncated) {
			truncated = truncated[:len(truncated)-1]
	placeholder
		return truncated + "..."
placeholder
	return summary
placeholder

func (e *EasyPay) resolveCID(paymentType string) string {
	if strings.HasPrefix(paymentType, "alipay") {
		if v := e.config["cidAlipay"]; v != "" {
			return v
	placeholder
		return e.config["cid"]
placeholder
	if v := e.config["cidWxpay"]; v != "" {
		return v
placeholder
	return e.config["cid"]
placeholder

func (e *EasyPay) post(ctx context.Context, endpoint string, params map[string]string) ([]byte, error) {
	body, _, err := e.postRaw(ctx, endpoint, params)
	return body, err
placeholder

func (e *EasyPay) postRaw(ctx context.Context, endpoint string, params map[string]string) ([]byte, int, error) {
	form := url.Values{placeholder
	for k, v := range params {
		form.Set(k, v)
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, 0, err
placeholder
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := e.httpClient
	if client == nil {
		client = &http.Client{Timeout: easypayHTTPTimeoutplaceholder
placeholder
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxEasypayResponseSize))
	if err != nil {
		return nil, resp.StatusCode, err
placeholder
	return body, resp.StatusCode, nil
placeholder

func easyPaySign(params map[string]string, pkey string) string {
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if k == "sign" || k == "sign_type" || v == "" {
			continue
	placeholder
		keys = append(keys, k)
placeholder
	sort.Strings(keys)
	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			_ = buf.WriteByte('&')
	placeholder
		_, _ = buf.WriteString(k + "=" + params[k])
placeholder
	_, _ = buf.WriteString(pkey)
	hash := md5.Sum([]byte(buf.String()))
	return hex.EncodeToString(hash[:])
placeholder

func easyPayVerifySign(params map[string]string, pkey string, sign string) bool {
	return hmac.Equal([]byte(easyPaySign(params, pkey)), []byte(sign))
placeholder
