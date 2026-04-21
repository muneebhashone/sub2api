package service

import (
	"context"
	"testing"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/payment"
)

func TestUnionFloat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		agg         float64
		limited     bool
		val         float64
		wantMin     bool
		wantAgg     float64
		wantLimited bool
placeholder{
		{"first non-zero value", 0, true, 5, true, 5, trueplaceholder,
		{"lower min replaces", 10, true, 3, true, 3, trueplaceholder,
		{"higher min does not replace", 3, true, 10, true, 3, trueplaceholder,
		{"higher max replaces", 10, true, 20, false, 20, trueplaceholder,
		{"lower max does not replace", 20, true, 10, false, 20, trueplaceholder,
		{"zero value makes unlimited", 5, true, 0, true, 5, falseplaceholder,
		{"already unlimited stays unlimited", 5, false, 10, true, 5, falseplaceholder,
		{"zero on first call", 0, true, 0, true, 0, falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotAgg, gotLimited := unionFloat(tt.agg, tt.limited, tt.val, tt.wantMin)
			if gotAgg != tt.wantAgg || gotLimited != tt.wantLimited {
				t.Fatalf("unionFloat(%v, %v, %v, %v) = (%v, %v), want (%v, %v)",
					tt.agg, tt.limited, tt.val, tt.wantMin,
					gotAgg, gotLimited, tt.wantAgg, tt.wantLimited)
		placeholder
	placeholder)
placeholder
placeholder

func makeInstance(id int64, providerKey, supportedTypes, limits string) *dbent.PaymentProviderInstance {
	return &dbent.PaymentProviderInstance{
		ID:             id,
		ProviderKey:    providerKey,
		SupportedTypes: supportedTypes,
		Limits:         limits,
		Enabled:        true,
placeholder
placeholder

func TestPcAggregateMethodLimits(t *testing.T) {
	t.Parallel()

	t.Run("single instance with limits", func(t *testing.T) {
		t.Parallel()
		inst := makeInstance(1, "easypay", "alipay,wxpay",
			`{"alipay":{"singleMin":2,"singleMax":14placeholder,"wxpay":{"singleMin":1,"singleMax":12placeholderplaceholder`)
		ml := pcAggregateMethodLimits("alipay", []*dbent.PaymentProviderInstance{instplaceholder)
		if ml.SingleMin != 2 || ml.SingleMax != 14 {
			t.Fatalf("alipay limits = min:%v max:%v, want min:2 max:14", ml.SingleMin, ml.SingleMax)
	placeholder
placeholder)

	t.Run("two instances union takes widest range", func(t *testing.T) {
		t.Parallel()
		inst1 := makeInstance(1, "easypay", "alipay,wxpay",
			`{"alipay":{"singleMin":5,"singleMax":100placeholderplaceholder`)
		inst2 := makeInstance(2, "easypay", "alipay,wxpay",
			`{"alipay":{"singleMin":2,"singleMax":200placeholderplaceholder`)
		ml := pcAggregateMethodLimits("alipay", []*dbent.PaymentProviderInstance{inst1, inst2placeholder)
		if ml.SingleMin != 2 {
			t.Fatalf("SingleMin = %v, want 2 (lowest floor)", ml.SingleMin)
	placeholder
		if ml.SingleMax != 200 {
			t.Fatalf("SingleMax = %v, want 200 (highest ceiling)", ml.SingleMax)
	placeholder
placeholder)

	t.Run("one instance unlimited makes aggregate unlimited", func(t *testing.T) {
		t.Parallel()
		inst1 := makeInstance(1, "easypay", "wxpay",
			`{"wxpay":{"singleMin":3,"singleMax":10placeholderplaceholder`)
		inst2 := makeInstance(2, "easypay", "wxpay", "") // no limits = unlimited
		ml := pcAggregateMethodLimits("wxpay", []*dbent.PaymentProviderInstance{inst1, inst2placeholder)
		if ml.SingleMin != 0 || ml.SingleMax != 0 {
			t.Fatalf("limits = min:%v max:%v, want min:0 max:0 (unlimited)", ml.SingleMin, ml.SingleMax)
	placeholder
placeholder)

	t.Run("one field unlimited others limited", func(t *testing.T) {
		t.Parallel()
		inst1 := makeInstance(1, "easypay", "alipay",
			`{"alipay":{"singleMin":5,"singleMax":100placeholderplaceholder`)
		inst2 := makeInstance(2, "easypay", "alipay",
			`{"alipay":{"singleMin":3,"singleMax":0placeholderplaceholder`) // singleMax=0 = unlimited
		ml := pcAggregateMethodLimits("alipay", []*dbent.PaymentProviderInstance{inst1, inst2placeholder)
		if ml.SingleMin != 3 {
			t.Fatalf("SingleMin = %v, want 3 (lowest floor)", ml.SingleMin)
	placeholder
		if ml.SingleMax != 0 {
			t.Fatalf("SingleMax = %v, want 0 (unlimited)", ml.SingleMax)
	placeholder
placeholder)

	t.Run("empty instances returns zeros", func(t *testing.T) {
		t.Parallel()
		ml := pcAggregateMethodLimits("alipay", nil)
		if ml.SingleMin != 0 || ml.SingleMax != 0 || ml.DailyLimit != 0 {
			t.Fatalf("empty instances should return all zeros, got %+v", ml)
	placeholder
placeholder)

	t.Run("invalid JSON treated as unlimited", func(t *testing.T) {
		t.Parallel()
		inst := makeInstance(1, "easypay", "alipay", `{invalid jsonplaceholder`)
		ml := pcAggregateMethodLimits("alipay", []*dbent.PaymentProviderInstance{instplaceholder)
		if ml.SingleMin != 0 || ml.SingleMax != 0 {
			t.Fatalf("invalid JSON should be treated as unlimited, got %+v", ml)
	placeholder
placeholder)

	t.Run("type not in limits JSON treated as unlimited", func(t *testing.T) {
		t.Parallel()
		inst := makeInstance(1, "easypay", "alipay,wxpay",
			`{"wxpay":{"singleMin":1,"singleMax":10placeholderplaceholder`) // only wxpay, no alipay
		ml := pcAggregateMethodLimits("alipay", []*dbent.PaymentProviderInstance{instplaceholder)
		if ml.SingleMin != 0 || ml.SingleMax != 0 {
			t.Fatalf("missing type should be treated as unlimited, got %+v", ml)
	placeholder
placeholder)

	t.Run("daily limit aggregation", func(t *testing.T) {
		t.Parallel()
		inst1 := makeInstance(1, "easypay", "alipay",
			`{"alipay":{"singleMin":1,"singleMax":100,"dailyLimit":500placeholderplaceholder`)
		inst2 := makeInstance(2, "easypay", "alipay",
			`{"alipay":{"singleMin":2,"singleMax":200,"dailyLimit":1000placeholderplaceholder`)
		ml := pcAggregateMethodLimits("alipay", []*dbent.PaymentProviderInstance{inst1, inst2placeholder)
		if ml.DailyLimit != 1000 {
			t.Fatalf("DailyLimit = %v, want 1000 (highest cap)", ml.DailyLimit)
	placeholder
placeholder)
placeholder

func TestPcGroupByPaymentType(t *testing.T) {
	t.Parallel()

	t.Run("stripe instance maps all types to stripe group", func(t *testing.T) {
		t.Parallel()
		stripe := makeInstance(1, payment.TypeStripe, "card,alipay,link,wxpay", "")
		easypay := makeInstance(2, payment.TypeEasyPay, "alipay,wxpay", "")

		groups := pcGroupByPaymentType([]*dbent.PaymentProviderInstance{stripe, easypayplaceholder)

		// Stripe instance should only be in "stripe" group
		if len(groups[payment.TypeStripe]) != 1 || groups[payment.TypeStripe][0].ID != 1 {
			t.Fatalf("stripe group should contain only stripe instance, got %v", groups[payment.TypeStripe])
	placeholder
		// alipay group should only contain easypay, NOT stripe
		if len(groups[payment.TypeAlipay]) != 1 || groups[payment.TypeAlipay][0].ID != 2 {
			t.Fatalf("alipay group should contain only easypay instance, got %v", groups[payment.TypeAlipay])
	placeholder
		// wxpay group should only contain easypay, NOT stripe
		if len(groups[payment.TypeWxpay]) != 1 || groups[payment.TypeWxpay][0].ID != 2 {
			t.Fatalf("wxpay group should contain only easypay instance, got %v", groups[payment.TypeWxpay])
	placeholder
placeholder)

	t.Run("multiple easypay instances in same groups", func(t *testing.T) {
		t.Parallel()
		ep1 := makeInstance(1, payment.TypeEasyPay, "alipay,wxpay", "")
		ep2 := makeInstance(2, payment.TypeEasyPay, "alipay,wxpay", "")

		groups := pcGroupByPaymentType([]*dbent.PaymentProviderInstance{ep1, ep2placeholder)

		if len(groups[payment.TypeAlipay]) != 2 {
			t.Fatalf("alipay group should have 2 instances, got %d", len(groups[payment.TypeAlipay]))
	placeholder
		if len(groups[payment.TypeWxpay]) != 2 {
			t.Fatalf("wxpay group should have 2 instances, got %d", len(groups[payment.TypeWxpay]))
	placeholder
placeholder)

	t.Run("stripe with no supported types still in stripe group", func(t *testing.T) {
		t.Parallel()
		stripe := makeInstance(1, payment.TypeStripe, "", "")

		groups := pcGroupByPaymentType([]*dbent.PaymentProviderInstance{stripeplaceholder)

		if len(groups[payment.TypeStripe]) != 1 {
			t.Fatalf("stripe with empty types should still be in stripe group, got %v", groups)
	placeholder
placeholder)
placeholder

func TestPcComputeGlobalRange(t *testing.T) {
	t.Parallel()

	t.Run("all methods have limits", func(t *testing.T) {
		t.Parallel()
		methods := map[string]MethodLimits{
			"alipay": {SingleMin: 2, SingleMax: 14placeholder,
			"wxpay":  {SingleMin: 1, SingleMax: 12placeholder,
			"stripe": {SingleMin: 5, SingleMax: 100placeholder,
	placeholder
		gMin, gMax := pcComputeGlobalRange(methods)
		if gMin != 1 {
			t.Fatalf("global min = %v, want 1 (lowest floor)", gMin)
	placeholder
		if gMax != 100 {
			t.Fatalf("global max = %v, want 100 (highest ceiling)", gMax)
	placeholder
placeholder)

	t.Run("one method unlimited makes global unlimited", func(t *testing.T) {
		t.Parallel()
		methods := map[string]MethodLimits{
			"alipay": {SingleMin: 2, SingleMax: 14placeholder,
			"stripe": {SingleMin: 0, SingleMax: 0placeholder, // unlimited
	placeholder
		gMin, gMax := pcComputeGlobalRange(methods)
		if gMin != 0 {
			t.Fatalf("global min = %v, want 0 (unlimited)", gMin)
	placeholder
		if gMax != 0 {
			t.Fatalf("global max = %v, want 0 (unlimited)", gMax)
	placeholder
placeholder)

	t.Run("empty methods returns zeros", func(t *testing.T) {
		t.Parallel()
		gMin, gMax := pcComputeGlobalRange(map[string]MethodLimits{placeholder)
		if gMin != 0 || gMax != 0 {
			t.Fatalf("empty methods should return (0, 0), got (%v, %v)", gMin, gMax)
	placeholder
placeholder)

	t.Run("only min unlimited", func(t *testing.T) {
		t.Parallel()
		methods := map[string]MethodLimits{
			"alipay": {SingleMin: 0, SingleMax: 100placeholder,
			"wxpay":  {SingleMin: 5, SingleMax: 50placeholder,
	placeholder
		gMin, gMax := pcComputeGlobalRange(methods)
		if gMin != 0 {
			t.Fatalf("global min = %v, want 0 (unlimited)", gMin)
	placeholder
		if gMax != 100 {
			t.Fatalf("global max = %v, want 100", gMax)
	placeholder
placeholder)
placeholder

func TestPcInstanceTypeLimits(t *testing.T) {
	t.Parallel()

	t.Run("empty limits string returns false", func(t *testing.T) {
		t.Parallel()
		inst := makeInstance(1, "easypay", "alipay", "")
		_, ok := pcInstanceTypeLimits(inst, "alipay")
		if ok {
			t.Fatal("expected ok=false for empty limits")
	placeholder
placeholder)

	t.Run("type found returns correct values", func(t *testing.T) {
		t.Parallel()
		inst := makeInstance(1, "easypay", "alipay",
			`{"alipay":{"singleMin":2,"singleMax":14,"dailyLimit":500placeholderplaceholder`)
		cl, ok := pcInstanceTypeLimits(inst, "alipay")
		if !ok {
			t.Fatal("expected ok=true")
	placeholder
		if cl.SingleMin != 2 || cl.SingleMax != 14 || cl.DailyLimit != 500 {
			t.Fatalf("limits = %+v, want min:2 max:14 daily:500", cl)
	placeholder
placeholder)

	t.Run("type not found returns false", func(t *testing.T) {
		t.Parallel()
		inst := makeInstance(1, "easypay", "alipay",
			`{"wxpay":{"singleMin":1placeholderplaceholder`)
		_, ok := pcInstanceTypeLimits(inst, "alipay")
		if ok {
			t.Fatal("expected ok=false for missing type")
	placeholder
placeholder)

	t.Run("invalid JSON returns false", func(t *testing.T) {
		t.Parallel()
		inst := makeInstance(1, "easypay", "alipay", `{bad jsonplaceholder`)
		_, ok := pcInstanceTypeLimits(inst, "alipay")
		if ok {
			t.Fatal("expected ok=false for invalid JSON")
	placeholder
placeholder)
placeholder

func TestGetAvailableMethodLimitsHidesConflictingVisibleMethodProviders(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)

	_, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeAlipay).
		SetName("Official Alipay").
		SetConfig("{placeholder").
		SetSupportedTypes("alipay").
		SetLimits(`{"alipay":{"singleMin":10,"singleMax":100placeholderplaceholder`).
		SetEnabled(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("create official alipay instance: %v", err)
placeholder
	_, err = client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeEasyPay).
		SetName("EasyPay Alipay").
		SetConfig("{placeholder").
		SetSupportedTypes("alipay").
		SetLimits(`{"alipay":{"singleMin":20,"singleMax":200placeholderplaceholder`).
		SetEnabled(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("create easypay alipay instance: %v", err)
placeholder
	_, err = client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeWxpay).
		SetName("Official WeChat").
		SetConfig("{placeholder").
		SetSupportedTypes("wxpay").
		SetLimits(`{"wxpay":{"singleMin":30,"singleMax":300placeholderplaceholder`).
		SetEnabled(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("create official wxpay instance: %v", err)
placeholder

	svc := &PaymentConfigService{
		entClient: client,
placeholder

	resp, err := svc.GetAvailableMethodLimits(ctx)
	if err != nil {
		t.Fatalf("GetAvailableMethodLimits returned error: %v", err)
placeholder

	if _, ok := resp.Methods[payment.TypeAlipay]; ok {
		t.Fatalf("alipay should be hidden when multiple enabled providers claim it, got %v", resp.Methods[payment.TypeAlipay])
placeholder

	wxpayLimits, ok := resp.Methods[payment.TypeWxpay]
	if !ok {
		t.Fatalf("expected wxpay limits to remain visible, got %v", resp.Methods)
placeholder
	if wxpayLimits.SingleMin != 30 || wxpayLimits.SingleMax != 300 {
		t.Fatalf("wxpay limits = %+v, want official-only min=30 max=300", wxpayLimits)
placeholder
	if resp.GlobalMin != 30 || resp.GlobalMax != 300 {
		t.Fatalf("global range = (%v, %v), want (30, 300)", resp.GlobalMin, resp.GlobalMax)
placeholder
placeholder
