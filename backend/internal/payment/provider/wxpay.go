package provider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

// WeChat Pay constants.
const (
	wxpayCurrency = "CNY"
	wxpayH5Type   = "Wap"
)

// WeChat Pay trade states.
const (
	wxpayTradeStateSuccess  = "SUCCESS"
	wxpayTradeStateRefund   = "REFUND"
	wxpayTradeStateClosed   = "CLOSED"
	wxpayTradeStatePayError = "PAYERROR"
)

// WeChat Pay notification event types.
const (
	wxpayEventTransactionSuccess = "TRANSACTION.SUCCESS"
)

// WeChat Pay error codes.
const (
	wxpayErrNoAuth = "NO_AUTH"
)

type Wxpay struct {
	instanceID    string
	config        map[string]string
	mu            sync.Mutex
	coreClient    *core.Client
	notifyHandler *notify.Handler
placeholder

const wxpayAPIv3KeyLength = 32

func NewWxpay(instanceID string, config map[string]string) (*Wxpay, error) {
	// All fields are required. Platform-certificate mode is intentionally unsupported —
	// WeChat has been migrating all merchants to the pubkey verifier since 2024-10,
	// and newly-provisioned merchants cannot download platform certificates at all.
	required := []string{"appId", "mchId", "privateKey", "apiV3Key", "certSerial", "publicKey", "publicKeyId"placeholder
	for _, k := range required {
		if config[k] == "" {
			return nil, infraerrors.BadRequest("WXPAY_CONFIG_MISSING_KEY", "missing_required_key").
				WithMetadata(map[string]string{"key": kplaceholder)
	placeholder
placeholder
	if len(config["apiV3Key"]) != wxpayAPIv3KeyLength {
		return nil, infraerrors.BadRequest("WXPAY_CONFIG_INVALID_KEY_LENGTH", "invalid_key_length").
			WithMetadata(map[string]string{
				"key":      "apiV3Key",
				"expected": strconv.Itoa(wxpayAPIv3KeyLength),
				"actual":   strconv.Itoa(len(config["apiV3Key"])),
		placeholder)
placeholder
	// Parse PEMs eagerly so malformed keys surface at save time, not at order creation.
	if _, err := utils.LoadPrivateKey(formatPEM(config["privateKey"], "PRIVATE KEY")); err != nil {
		return nil, infraerrors.BadRequest("WXPAY_CONFIG_INVALID_KEY", "invalid_key").
			WithMetadata(map[string]string{"key": "privateKey"placeholder)
placeholder
	if _, err := utils.LoadPublicKey(formatPEM(config["publicKey"], "PUBLIC KEY")); err != nil {
		return nil, infraerrors.BadRequest("WXPAY_CONFIG_INVALID_KEY", "invalid_key").
			WithMetadata(map[string]string{"key": "publicKey"placeholder)
placeholder
	return &Wxpay{instanceID: instanceID, config: configplaceholder, nil
placeholder

func (w *Wxpay) Name() string        { return "Wxpay" placeholder
func (w *Wxpay) ProviderKey() string { return payment.TypeWxpay placeholder
func (w *Wxpay) SupportedTypes() []payment.PaymentType {
	return []payment.PaymentType{payment.TypeWxpayplaceholder
placeholder

func formatPEM(key, keyType string) string {
	key = strings.TrimSpace(key)
	if strings.HasPrefix(key, "-----BEGIN") {
		return key
placeholder
	return fmt.Sprintf("-----BEGIN %s-----\n%s\n-----END %s-----", keyType, key, keyType)
placeholder

func (w *Wxpay) ensureClient() (*core.Client, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.coreClient != nil {
		return w.coreClient, nil
placeholder
	privateKey, err := utils.LoadPrivateKey(formatPEM(w.config["privateKey"], "PRIVATE KEY"))
	if err != nil {
		return nil, infraerrors.BadRequest("WXPAY_CONFIG_INVALID_KEY", "invalid_key").
			WithMetadata(map[string]string{"key": "privateKey"placeholder)
placeholder
	publicKey, err := utils.LoadPublicKey(formatPEM(w.config["publicKey"], "PUBLIC KEY"))
	if err != nil {
		return nil, infraerrors.BadRequest("WXPAY_CONFIG_INVALID_KEY", "invalid_key").
			WithMetadata(map[string]string{"key": "publicKey"placeholder)
placeholder
	verifier := verifiers.NewSHA256WithRSAPubkeyVerifier(w.config["publicKeyId"], *publicKey)
	client, err := core.NewClient(context.Background(),
		option.WithMerchantCredential(w.config["mchId"], w.config["certSerial"], privateKey),
		option.WithVerifier(verifier))
	if err != nil {
		return nil, fmt.Errorf("wxpay init client: %w", err)
placeholder
	handler, err := notify.NewRSANotifyHandler(w.config["apiV3Key"], verifier)
	if err != nil {
		return nil, fmt.Errorf("wxpay init notify handler: %w", err)
placeholder
	w.notifyHandler = handler
	w.coreClient = client
	return w.coreClient, nil
placeholder

func (w *Wxpay) CreatePayment(ctx context.Context, req payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	client, err := w.ensureClient()
	if err != nil {
		return nil, err
placeholder
	// Request-first, config-fallback (consistent with EasyPay/Alipay)
	notifyURL := req.NotifyURL
	if notifyURL == "" {
		notifyURL = w.config["notifyUrl"]
placeholder
	if notifyURL == "" {
		return nil, fmt.Errorf("wxpay notifyUrl is required")
placeholder
	totalFen, err := payment.YuanToFen(req.Amount)
	if err != nil {
		return nil, fmt.Errorf("wxpay create payment: %w", err)
placeholder
	if req.IsMobile && req.ClientIP != "" {
		resp, err := w.createOrder(ctx, client, req, notifyURL, totalFen, true)
		if err == nil {
			return resp, nil
	placeholder
		if !strings.Contains(err.Error(), wxpayErrNoAuth) {
			return nil, err
	placeholder
		slog.Warn("wxpay H5 payment not authorized, falling back to native", "order", req.OrderID)
placeholder
	return w.createOrder(ctx, client, req, notifyURL, totalFen, false)
placeholder

func (w *Wxpay) createOrder(ctx context.Context, c *core.Client, req payment.CreatePaymentRequest, notifyURL string, totalFen int64, useH5 bool) (*payment.CreatePaymentResponse, error) {
	if useH5 {
		return w.prepayH5(ctx, c, req, notifyURL, totalFen)
placeholder
	return w.prepayNative(ctx, c, req, notifyURL, totalFen)
placeholder

func (w *Wxpay) prepayNative(ctx context.Context, c *core.Client, req payment.CreatePaymentRequest, notifyURL string, totalFen int64) (*payment.CreatePaymentResponse, error) {
	svc := native.NativeApiService{Client: cplaceholder
	cur := wxpayCurrency
	resp, _, err := svc.Prepay(ctx, native.PrepayRequest{
		Appid: core.String(w.config["appId"]), Mchid: core.String(w.config["mchId"]),
		Description: core.String(req.Subject), OutTradeNo: core.String(req.OrderID),
		NotifyUrl: core.String(notifyURL),
		Amount:    &native.Amount{Total: core.Int64(totalFen), Currency: &curplaceholder,
placeholder)
	if err != nil {
		return nil, fmt.Errorf("wxpay native prepay: %w", err)
placeholder
	codeURL := ""
	if resp.CodeUrl != nil {
		codeURL = *resp.CodeUrl
placeholder
	return &payment.CreatePaymentResponse{TradeNo: req.OrderID, QRCode: codeURLplaceholder, nil
placeholder

func (w *Wxpay) prepayH5(ctx context.Context, c *core.Client, req payment.CreatePaymentRequest, notifyURL string, totalFen int64) (*payment.CreatePaymentResponse, error) {
	svc := h5.H5ApiService{Client: cplaceholder
	cur := wxpayCurrency
	tp := wxpayH5Type
	resp, _, err := svc.Prepay(ctx, h5.PrepayRequest{
		Appid: core.String(w.config["appId"]), Mchid: core.String(w.config["mchId"]),
		Description: core.String(req.Subject), OutTradeNo: core.String(req.OrderID),
		NotifyUrl: core.String(notifyURL),
		Amount:    &h5.Amount{Total: core.Int64(totalFen), Currency: &curplaceholder,
		SceneInfo: &h5.SceneInfo{PayerClientIp: core.String(req.ClientIP), H5Info: &h5.H5Info{Type: &tpplaceholderplaceholder,
placeholder)
	if err != nil {
		return nil, fmt.Errorf("wxpay h5 prepay: %w", err)
placeholder
	h5URL := ""
	if resp.H5Url != nil {
		h5URL = *resp.H5Url
placeholder
	return &payment.CreatePaymentResponse{TradeNo: req.OrderID, PayURL: h5URLplaceholder, nil
placeholder

func wxSV(s *string) string {
	if s == nil {
		return ""
placeholder
	return *s
placeholder

func mapWxState(s string) string {
	switch s {
	case wxpayTradeStateSuccess:
		return payment.ProviderStatusPaid
	case wxpayTradeStateRefund:
		return payment.ProviderStatusRefunded
	case wxpayTradeStateClosed, wxpayTradeStatePayError:
		return payment.ProviderStatusFailed
	default:
		return payment.ProviderStatusPending
placeholder
placeholder

func (w *Wxpay) QueryOrder(ctx context.Context, tradeNo string) (*payment.QueryOrderResponse, error) {
	c, err := w.ensureClient()
	if err != nil {
		return nil, err
placeholder
	svc := native.NativeApiService{Client: cplaceholder
	tx, _, err := svc.QueryOrderByOutTradeNo(ctx, native.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(tradeNo), Mchid: core.String(w.config["mchId"]),
placeholder)
	if err != nil {
		return nil, fmt.Errorf("wxpay query order: %w", err)
placeholder
	var amt float64
	if tx.Amount != nil && tx.Amount.Total != nil {
		amt = payment.FenToYuan(*tx.Amount.Total)
placeholder
	id := tradeNo
	if tx.TransactionId != nil {
		id = *tx.TransactionId
placeholder
	pa := ""
	if tx.SuccessTime != nil {
		pa = *tx.SuccessTime
placeholder
	return &payment.QueryOrderResponse{TradeNo: id, Status: mapWxState(wxSV(tx.TradeState)), Amount: amt, PaidAt: paplaceholder, nil
placeholder

func (w *Wxpay) VerifyNotification(ctx context.Context, rawBody string, headers map[string]string) (*payment.PaymentNotification, error) {
	if _, err := w.ensureClient(); err != nil {
		return nil, err
placeholder
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", io.NopCloser(bytes.NewBufferString(rawBody)))
	if err != nil {
		return nil, fmt.Errorf("wxpay construct request: %w", err)
placeholder
	for k, v := range headers {
		r.Header.Set(k, v)
placeholder
	var tx payments.Transaction
	nr, err := w.notifyHandler.ParseNotifyRequest(ctx, r, &tx)
	if err != nil {
		return nil, fmt.Errorf("wxpay verify notification: %w", err)
placeholder
	if nr.EventType != wxpayEventTransactionSuccess {
		return nil, nil
placeholder
	var amt float64
	if tx.Amount != nil && tx.Amount.Total != nil {
		amt = payment.FenToYuan(*tx.Amount.Total)
placeholder
	st := payment.ProviderStatusFailed
	if wxSV(tx.TradeState) == wxpayTradeStateSuccess {
		st = payment.ProviderStatusSuccess
placeholder
	return &payment.PaymentNotification{
		TradeNo: wxSV(tx.TransactionId), OrderID: wxSV(tx.OutTradeNo),
		Amount: amt, Status: st, RawData: rawBody,
placeholder, nil
placeholder

func (w *Wxpay) Refund(ctx context.Context, req payment.RefundRequest) (*payment.RefundResponse, error) {
	c, err := w.ensureClient()
	if err != nil {
		return nil, err
placeholder
	rf, err := payment.YuanToFen(req.Amount)
	if err != nil {
		return nil, fmt.Errorf("wxpay refund amount: %w", err)
placeholder
	tf, err := w.queryOrderTotalFen(ctx, c, req.OrderID)
	if err != nil {
		return nil, err
placeholder
	rs := refunddomestic.RefundsApiService{Client: cplaceholder
	cur := wxpayCurrency
	res, _, err := rs.Create(ctx, refunddomestic.CreateRequest{
		OutTradeNo:  core.String(req.OrderID),
		OutRefundNo: core.String(fmt.Sprintf("%s-refund-%d", req.OrderID, time.Now().UnixNano())),
		Reason:      core.String(req.Reason),
		Amount:      &refunddomestic.AmountReq{Refund: core.Int64(rf), Total: core.Int64(tf), Currency: &curplaceholder,
placeholder)
	if err != nil {
		return nil, fmt.Errorf("wxpay refund: %w", err)
placeholder
	rid := wxSV(res.RefundId)
	if rid == "" {
		rid = fmt.Sprintf("%s-refund", req.OrderID)
placeholder
	st := payment.ProviderStatusPending
	if res.Status != nil && *res.Status == refunddomestic.STATUS_SUCCESS {
		st = payment.ProviderStatusSuccess
placeholder
	return &payment.RefundResponse{RefundID: rid, Status: stplaceholder, nil
placeholder

func (w *Wxpay) queryOrderTotalFen(ctx context.Context, c *core.Client, orderID string) (int64, error) {
	svc := native.NativeApiService{Client: cplaceholder
	tx, _, err := svc.QueryOrderByOutTradeNo(ctx, native.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(orderID), Mchid: core.String(w.config["mchId"]),
placeholder)
	if err != nil {
		return 0, fmt.Errorf("wxpay refund query order: %w", err)
placeholder
	var tf int64
	if tx.Amount != nil && tx.Amount.Total != nil {
		tf = *tx.Amount.Total
placeholder
	return tf, nil
placeholder

func (w *Wxpay) CancelPayment(ctx context.Context, tradeNo string) error {
	c, err := w.ensureClient()
	if err != nil {
		return err
placeholder
	svc := native.NativeApiService{Client: cplaceholder
	_, err = svc.CloseOrder(ctx, native.CloseOrderRequest{
		OutTradeNo: core.String(tradeNo), Mchid: core.String(w.config["mchId"]),
placeholder)
	if err != nil {
		return fmt.Errorf("wxpay cancel payment: %w", err)
placeholder
	return nil
placeholder

var (
	_ payment.Provider           = (*Wxpay)(nil)
	_ payment.CancelableProvider = (*Wxpay)(nil)
)
