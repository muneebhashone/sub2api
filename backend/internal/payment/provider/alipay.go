package provider

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/smartwalle/alipay/v3"
)

// Alipay product codes.
const (
	alipayProductCodePreCreate = "FACE_TO_FACE_PAYMENT"
	alipayProductCodeWapPay    = "QUICK_WAP_WAY"
	alipayProductCodePagePay   = "FAST_INSTANT_TRADE_PAY"
)

// Alipay response constants.
const (
	alipayFundChangeYes    = "Y"
	alipayErrTradeNotExist = "ACQ.TRADE_NOT_EXIST"
	alipayRefundSuffix     = "-refund"
)

var (
	alipayTradeWapPay = func(client *alipay.Client, param alipay.TradeWapPay) (*url.URL, error) {
		return client.TradeWapPay(param)
placeholder
	alipayTradePreCreate = func(ctx context.Context, client *alipay.Client, param alipay.TradePreCreate) (*alipay.TradePreCreateRsp, error) {
		return client.TradePreCreate(ctx, param)
placeholder
	alipayTradePagePay = func(client *alipay.Client, param alipay.TradePagePay) (*url.URL, error) {
		return client.TradePagePay(param)
placeholder
)

// Alipay implements payment.Provider and payment.CancelableProvider using the smartwalle/alipay SDK.
type Alipay struct {
	instanceID string
	config     map[string]string // appId, privateKey, publicKey (or alipayPublicKey), notifyUrl, returnUrl

	mu     sync.Mutex
	client *alipay.Client
placeholder

// NewAlipay creates a new Alipay provider instance.
func NewAlipay(instanceID string, config map[string]string) (*Alipay, error) {
	required := []string{"appId", "privateKey"placeholder
	for _, k := range required {
		if config[k] == "" {
			return nil, fmt.Errorf("alipay config missing required key: %s", k)
	placeholder
placeholder
	return &Alipay{
		instanceID: instanceID,
		config:     config,
placeholder, nil
placeholder

func (a *Alipay) getClient() (*alipay.Client, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.client != nil {
		return a.client, nil
placeholder
	client, err := alipay.New(a.config["appId"], a.config["privateKey"], true)
	if err != nil {
		return nil, fmt.Errorf("alipay init client: %w", err)
placeholder
	pubKey := a.config["publicKey"]
	if pubKey == "" {
		pubKey = a.config["alipayPublicKey"]
placeholder
	if pubKey == "" {
		return nil, fmt.Errorf("alipay config missing required key: publicKey (or alipayPublicKey)")
placeholder
	if err := client.LoadAliPayPublicKey(pubKey); err != nil {
		return nil, fmt.Errorf("alipay load public key: %w", err)
placeholder
	a.client = client
	return a.client, nil
placeholder

func (a *Alipay) Name() string        { return "Alipay" placeholder
func (a *Alipay) ProviderKey() string { return payment.TypeAlipay placeholder
func (a *Alipay) SupportedTypes() []payment.PaymentType {
	return []payment.PaymentType{payment.TypeAlipayplaceholder
placeholder

func (a *Alipay) MerchantIdentityMetadata() map[string]string {
	if a == nil {
		return nil
placeholder
	appID := strings.TrimSpace(a.config["appId"])
	if appID == "" {
		return nil
placeholder
	return map[string]string{"app_id": appIDplaceholder
placeholder

// CreatePayment creates an Alipay payment using the following routing:
//   - Mobile (H5), default: alipay.trade.wap.pay — browser redirect into Alipay.
//   - Mobile with AlipayMobilePrecreate: alipay.trade.precreate — return the
//     dynamic QR payload so the frontend can open it through the Alipay app.
//   - Desktop, default: prefer alipay.trade.precreate (FACE_TO_FACE_PAYMENT) to
//     get a scannable QR payload. If precreate is unavailable for the merchant,
//     fall back to alipay.trade.page.pay and expose pay_url only — the frontend
//     opens the Alipay checkout in a new tab.
//   - Desktop, paymentMode == "redirect": skip precreate and go straight to
//     alipay.trade.page.pay so the frontend always opens the Alipay checkout
//     in a new tab. Use this when the merchant has not enabled FACE_TO_FACE_PAYMENT.
//
// Note: alipay.trade.page.pay returns a checkout page URL, not a scannable
// payment QR. Never expose it via the QRCode field.
func (a *Alipay) CreatePayment(ctx context.Context, req payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	client, err := a.getClient()
	if err != nil {
		return nil, err
placeholder

	notifyURL := a.config["notifyUrl"]
	if req.NotifyURL != "" {
		notifyURL = req.NotifyURL
placeholder
	returnURL := a.config["returnUrl"]
	if req.ReturnURL != "" {
		returnURL = req.ReturnURL
placeholder

	if req.IsMobile {
		if req.AlipayMobilePrecreate {
			return a.createPrecreateTrade(ctx, client, req, notifyURL)
	placeholder
		return a.createWapTrade(client, req, notifyURL, returnURL)
placeholder
	return a.createDesktopTrade(ctx, client, req, notifyURL, returnURL)
placeholder

func (a *Alipay) createWapTrade(client *alipay.Client, req payment.CreatePaymentRequest, notifyURL, returnURL string) (*payment.CreatePaymentResponse, error) {
	param := alipay.TradeWapPay{placeholder
	param.OutTradeNo = req.OrderID
	param.TotalAmount = req.Amount
	param.Subject = req.Subject
	param.ProductCode = alipayProductCodeWapPay
	param.NotifyURL = notifyURL
	param.ReturnURL = returnURL

	payURL, err := alipayTradeWapPay(client, param)
	if err != nil {
		return nil, fmt.Errorf("alipay TradeWapPay: %w", err)
placeholder
	return &payment.CreatePaymentResponse{
		TradeNo: req.OrderID,
		PayURL:  payURL.String(),
placeholder, nil
placeholder

func (a *Alipay) createDesktopTrade(ctx context.Context, client *alipay.Client, req payment.CreatePaymentRequest, notifyURL, returnURL string) (*payment.CreatePaymentResponse, error) {
	// Explicit redirect mode: merchant opted into "always open the Alipay
	// checkout page in a new tab" via the provider instance's payment_mode.
	// Skip precreate to avoid a wasted API call.
	if strings.EqualFold(strings.TrimSpace(a.config["paymentMode"]), "redirect") {
		return a.createPagePayTrade(client, req, notifyURL, returnURL)
placeholder

	resp, precreateErr := a.createPrecreateTrade(ctx, client, req, notifyURL)
	if precreateErr == nil {
		return resp, nil
placeholder

	resp, pagePayErr := a.createPagePayTrade(client, req, notifyURL, returnURL)
	if pagePayErr == nil {
		return resp, nil
placeholder

	return nil, fmt.Errorf("alipay desktop payment failed: precreate=%v; pagepay=%w", precreateErr, pagePayErr)
placeholder

func (a *Alipay) createPrecreateTrade(ctx context.Context, client *alipay.Client, req payment.CreatePaymentRequest, notifyURL string) (*payment.CreatePaymentResponse, error) {
	param := alipay.TradePreCreate{placeholder
	param.OutTradeNo = req.OrderID
	param.TotalAmount = req.Amount
	param.Subject = req.Subject
	param.ProductCode = alipayProductCodePreCreate
	param.NotifyURL = notifyURL

	rsp, err := alipayTradePreCreate(ctx, client, param)
	if err != nil {
		return nil, fmt.Errorf("alipay TradePreCreate: %w", err)
placeholder
	if rsp == nil {
		return nil, fmt.Errorf("alipay TradePreCreate: empty response")
placeholder
	if rsp.IsFailure() {
		return nil, fmt.Errorf("alipay TradePreCreate failed: %s", rsp.Error.Error())
placeholder
	if strings.TrimSpace(rsp.QRCode) == "" {
		return nil, fmt.Errorf("alipay TradePreCreate: empty qr_code")
placeholder

	return &payment.CreatePaymentResponse{
		TradeNo: req.OrderID,
		QRCode:  rsp.QRCode,
placeholder, nil
placeholder

func (a *Alipay) createPagePayTrade(client *alipay.Client, req payment.CreatePaymentRequest, notifyURL, returnURL string) (*payment.CreatePaymentResponse, error) {
	param := alipay.TradePagePay{placeholder
	param.OutTradeNo = req.OrderID
	param.TotalAmount = req.Amount
	param.Subject = req.Subject
	param.ProductCode = alipayProductCodePagePay
	param.NotifyURL = notifyURL
	param.ReturnURL = returnURL

	payURL, err := alipayTradePagePay(client, param)
	if err != nil {
		return nil, fmt.Errorf("alipay TradePagePay: %w", err)
placeholder
	// Only PayURL is exposed: alipay.trade.page.pay returns a checkout page URL
	// that must be opened in a browser, not a scannable payment QR. Setting it
	// as QRCode would let the frontend render an unscannable image.
	return &payment.CreatePaymentResponse{
		TradeNo: req.OrderID,
		PayURL:  payURL.String(),
placeholder, nil
placeholder

// QueryOrder queries the trade status via Alipay.
func (a *Alipay) QueryOrder(ctx context.Context, tradeNo string) (*payment.QueryOrderResponse, error) {
	client, err := a.getClient()
	if err != nil {
		return nil, err
placeholder

	result, err := client.TradeQuery(ctx, alipay.TradeQuery{OutTradeNo: tradeNoplaceholder)
	if err != nil {
		if isTradeNotExist(err) {
			return &payment.QueryOrderResponse{
				TradeNo: tradeNo,
				Status:  payment.ProviderStatusPending,
		placeholder, nil
	placeholder
		return nil, fmt.Errorf("alipay TradeQuery: %w", err)
placeholder

	status := payment.ProviderStatusPending
	switch result.TradeStatus {
	case alipay.TradeStatusSuccess, alipay.TradeStatusFinished:
		status = payment.ProviderStatusPaid
	case alipay.TradeStatusClosed:
		status = payment.ProviderStatusFailed
placeholder

	amount, err := strconv.ParseFloat(result.TotalAmount, 64)
	if err != nil {
		amount, err = parseAlipayAmount(
			result.TotalAmount,
			result.ReceiptAmount,
			result.BuyerPayAmount,
			result.InvoiceAmount,
		)
		if err != nil {
			return nil, fmt.Errorf("alipay parse amount: %w", err)
	placeholder
placeholder

	return &payment.QueryOrderResponse{
		TradeNo:  result.TradeNo,
		Status:   status,
		Amount:   amount,
		PaidAt:   result.SendPayDate,
		Metadata: a.MerchantIdentityMetadata(),
placeholder, nil
placeholder

// VerifyNotification decodes and verifies an Alipay async notification.
func (a *Alipay) VerifyNotification(ctx context.Context, rawBody string, _ map[string]string) (*payment.PaymentNotification, error) {
	client, err := a.getClient()
	if err != nil {
		return nil, err
placeholder

	values, err := url.ParseQuery(rawBody)
	if err != nil {
		return nil, fmt.Errorf("alipay parse notification: %w", err)
placeholder

	notification, err := client.DecodeNotification(ctx, values)
	if err != nil {
		return nil, fmt.Errorf("alipay verify notification: %w", err)
placeholder

	status := payment.ProviderStatusFailed
	if notification.TradeStatus == alipay.TradeStatusSuccess || notification.TradeStatus == alipay.TradeStatusFinished {
		status = payment.ProviderStatusSuccess
placeholder

	amount, err := strconv.ParseFloat(notification.TotalAmount, 64)
	if err != nil {
		amount, err = parseAlipayAmount(
			notification.TotalAmount,
			notification.ReceiptAmount,
			notification.BuyerPayAmount,
		)
		if err != nil {
			return nil, fmt.Errorf("alipay parse notification amount: %w", err)
	placeholder
placeholder

	metadata := a.MerchantIdentityMetadata()
	if appID := strings.TrimSpace(notification.AppId); appID != "" {
		if metadata == nil {
			metadata = map[string]string{placeholder
	placeholder
		metadata["app_id"] = appID
placeholder

	return &payment.PaymentNotification{
		TradeNo:  notification.TradeNo,
		OrderID:  notification.OutTradeNo,
		Amount:   amount,
		Status:   status,
		RawData:  rawBody,
		Metadata: metadata,
placeholder, nil
placeholder

// Refund requests a refund through Alipay.
func (a *Alipay) Refund(ctx context.Context, req payment.RefundRequest) (*payment.RefundResponse, error) {
	client, err := a.getClient()
	if err != nil {
		return nil, err
placeholder

	result, err := client.TradeRefund(ctx, alipay.TradeRefund{
		OutTradeNo:   req.OrderID,
		RefundAmount: req.Amount,
		RefundReason: req.Reason,
		OutRequestNo: fmt.Sprintf("%s-refund-%d", req.OrderID, time.Now().UnixNano()),
placeholder)
	if err != nil {
		return nil, fmt.Errorf("alipay TradeRefund: %w", err)
placeholder

	refundStatus := payment.ProviderStatusPending
	if result.FundChange == alipayFundChangeYes {
		refundStatus = payment.ProviderStatusSuccess
placeholder

	refundID := result.TradeNo
	if refundID == "" {
		refundID = req.OrderID + alipayRefundSuffix
placeholder

	return &payment.RefundResponse{
		RefundID: refundID,
		Status:   refundStatus,
placeholder, nil
placeholder

// CancelPayment closes a pending trade on Alipay.
func (a *Alipay) CancelPayment(ctx context.Context, tradeNo string) error {
	client, err := a.getClient()
	if err != nil {
		return err
placeholder

	_, err = client.TradeClose(ctx, alipay.TradeClose{OutTradeNo: tradeNoplaceholder)
	if err != nil {
		if isTradeNotExist(err) {
			return nil
	placeholder
		return fmt.Errorf("alipay TradeClose: %w", err)
placeholder
	return nil
placeholder

func isTradeNotExist(err error) bool {
	if err == nil {
		return false
placeholder
	return strings.Contains(err.Error(), alipayErrTradeNotExist)
placeholder

func parseAlipayAmount(values ...string) (float64, error) {
	for _, raw := range values {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
	placeholder
		amount, err := strconv.ParseFloat(raw, 64)
		if err == nil {
			return amount, nil
	placeholder
placeholder
	return 0, fmt.Errorf("no valid amount field")
placeholder

// Ensure interface compliance.
var (
	_ payment.Provider                 = (*Alipay)(nil)
	_ payment.CancelableProvider       = (*Alipay)(nil)
	_ payment.MerchantIdentityProvider = (*Alipay)(nil)
)
