//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/stretchr/testify/assert"
)

type paymentFulfillmentTestProvider struct {
	key            string
	supportedTypes []payment.PaymentType
placeholder

func (p paymentFulfillmentTestProvider) Name() string        { return p.key placeholder
func (p paymentFulfillmentTestProvider) ProviderKey() string { return p.key placeholder
func (p paymentFulfillmentTestProvider) SupportedTypes() []payment.PaymentType {
	return p.supportedTypes
placeholder
func (p paymentFulfillmentTestProvider) CreatePayment(ctx context.Context, req payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	panic("unexpected call")
placeholder
func (p paymentFulfillmentTestProvider) QueryOrder(ctx context.Context, tradeNo string) (*payment.QueryOrderResponse, error) {
	panic("unexpected call")
placeholder
func (p paymentFulfillmentTestProvider) VerifyNotification(ctx context.Context, rawBody string, headers map[string]string) (*payment.PaymentNotification, error) {
	panic("unexpected call")
placeholder
func (p paymentFulfillmentTestProvider) Refund(ctx context.Context, req payment.RefundRequest) (*payment.RefundResponse, error) {
	panic("unexpected call")
placeholder

// ---------------------------------------------------------------------------
// resolveRedeemAction — pure idempotency decision logic
// ---------------------------------------------------------------------------

func TestResolveRedeemAction_CodeNotFound(t *testing.T) {
	t.Parallel()
	action := resolveRedeemAction(nil, nil)
	assert.Equal(t, redeemActionCreate, action, "nil code with nil error should create")
placeholder

func TestResolveRedeemAction_LookupError(t *testing.T) {
	t.Parallel()
	action := resolveRedeemAction(nil, errors.New("db connection lost"))
	assert.Equal(t, redeemActionCreate, action, "lookup error should fall back to create")
placeholder

func TestResolveRedeemAction_LookupErrorWithNonNilCode(t *testing.T) {
	t.Parallel()
	// Edge case: both code and error are non-nil (shouldn't happen in practice,
	// but the function should still treat error as authoritative)
	code := &RedeemCode{Status: StatusUnusedplaceholder
	action := resolveRedeemAction(code, errors.New("partial error"))
	assert.Equal(t, redeemActionCreate, action, "non-nil error should always result in create regardless of code")
placeholder

func TestResolveRedeemAction_CodeExistsAndUsed(t *testing.T) {
	t.Parallel()
	code := &RedeemCode{
		Code:   "test-code-123",
		Status: StatusUsed,
		Type:   RedeemTypeBalance,
		Value:  10.0,
placeholder
	action := resolveRedeemAction(code, nil)
	assert.Equal(t, redeemActionSkipCompleted, action, "used code should skip to completed")
placeholder

func TestResolveRedeemAction_CodeExistsAndUnused(t *testing.T) {
	t.Parallel()
	code := &RedeemCode{
		Code:   "test-code-456",
		Status: StatusUnused,
		Type:   RedeemTypeBalance,
		Value:  25.0,
placeholder
	action := resolveRedeemAction(code, nil)
	assert.Equal(t, redeemActionRedeem, action, "unused code should skip creation and proceed to redeem")
placeholder

func TestResolveRedeemAction_CodeExistsWithExpiredStatus(t *testing.T) {
	t.Parallel()
	// A code with a non-standard status (neither "unused" nor "used")
	// should NOT be treated as used, so it falls through to redeemActionRedeem.
	code := &RedeemCode{
		Code:   "expired-code",
		Status: StatusExpired,
placeholder
	action := resolveRedeemAction(code, nil)
	assert.Equal(t, redeemActionRedeem, action, "expired-status code is not IsUsed(), should redeem")
placeholder

// ---------------------------------------------------------------------------
// Table-driven comprehensive test
// ---------------------------------------------------------------------------

func TestResolveRedeemAction_Table(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		code     *RedeemCode
		err      error
		expected redeemAction
placeholder{
		{
			name:     "nil code, nil error — first run",
			code:     nil,
			err:      nil,
			expected: redeemActionCreate,
	placeholder,
		{
			name:     "nil code, lookup error — treat as not found",
			code:     nil,
			err:      ErrRedeemCodeNotFound,
			expected: redeemActionCreate,
	placeholder,
		{
			name:     "nil code, generic DB error — treat as not found",
			code:     nil,
			err:      errors.New("connection refused"),
			expected: redeemActionCreate,
	placeholder,
		{
			name:     "code exists, used — previous run completed redeem",
			code:     &RedeemCode{Status: StatusUsedplaceholder,
			err:      nil,
			expected: redeemActionSkipCompleted,
	placeholder,
		{
			name:     "code exists, unused — previous run created code but crashed before redeem",
			code:     &RedeemCode{Status: StatusUnusedplaceholder,
			err:      nil,
			expected: redeemActionRedeem,
	placeholder,
		{
			name:     "code exists but error also set — error takes precedence",
			code:     &RedeemCode{Status: StatusUsedplaceholder,
			err:      errors.New("unexpected"),
			expected: redeemActionCreate,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := resolveRedeemAction(tt.code, tt.err)
			assert.Equal(t, tt.expected, got)
	placeholder)
placeholder
placeholder

// ---------------------------------------------------------------------------
// redeemAction enum value sanity
// ---------------------------------------------------------------------------

func TestRedeemAction_DistinctValues(t *testing.T) {
	t.Parallel()
	// Ensure the three actions have distinct values (iota correctness)
	assert.NotEqual(t, redeemActionCreate, redeemActionRedeem)
	assert.NotEqual(t, redeemActionCreate, redeemActionSkipCompleted)
	assert.NotEqual(t, redeemActionRedeem, redeemActionSkipCompleted)
placeholder

// ---------------------------------------------------------------------------
// RedeemCode.IsUsed / CanUse interaction with resolveRedeemAction
// ---------------------------------------------------------------------------

func TestResolveRedeemAction_IsUsedCanUseConsistency(t *testing.T) {
	t.Parallel()

	usedCode := &RedeemCode{Status: StatusUsedplaceholder
	unusedCode := &RedeemCode{Status: StatusUnusedplaceholder

	// Verify our decision function is consistent with the domain model methods
	assert.True(t, usedCode.IsUsed())
	assert.False(t, usedCode.CanUse())
	assert.Equal(t, redeemActionSkipCompleted, resolveRedeemAction(usedCode, nil))

	assert.False(t, unusedCode.IsUsed())
	assert.True(t, unusedCode.CanUse())
	assert.Equal(t, redeemActionRedeem, resolveRedeemAction(unusedCode, nil))
placeholder

func TestExpectedNotificationProviderKeyPrefersOrderInstanceProvider(t *testing.T) {
	t.Parallel()

	registry := payment.NewRegistry()
	registry.Register(paymentFulfillmentTestProvider{
		key:            payment.TypeAlipay,
		supportedTypes: []payment.PaymentType{payment.TypeAlipayplaceholder,
placeholder)

	assert.Equal(t,
		payment.TypeEasyPay,
		expectedNotificationProviderKey(registry, payment.TypeAlipay, payment.TypeEasyPay),
	)
placeholder

func TestExpectedNotificationProviderKeyUsesRegistryMappingForLegacyOrders(t *testing.T) {
	t.Parallel()

	registry := payment.NewRegistry()
	registry.Register(paymentFulfillmentTestProvider{
		key:            payment.TypeEasyPay,
		supportedTypes: []payment.PaymentType{payment.TypeAlipayplaceholder,
placeholder)

	assert.Equal(t,
		payment.TypeEasyPay,
		expectedNotificationProviderKey(registry, payment.TypeAlipay, ""),
	)
placeholder

func TestExpectedNotificationProviderKeyFallsBackToPaymentType(t *testing.T) {
	t.Parallel()

	assert.Equal(t,
		payment.TypeWxpay,
		expectedNotificationProviderKey(nil, payment.TypeWxpay, ""),
	)
placeholder
