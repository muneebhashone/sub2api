package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/payment"
)

func TestUsesOfficialWxpayVisibleMethodDerivesFromEnabledProviderInstance(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)

	_, err := client.PaymentProviderInstance.Create().
		SetProviderKey(payment.TypeWxpay).
		SetName("Official WeChat").
		SetConfig("{placeholder").
		SetSupportedTypes("wxpay").
		SetEnabled(true).
		SetSortOrder(1).
		Save(ctx)
	if err != nil {
		t.Fatalf("create official wxpay instance: %v", err)
placeholder

	svc := &PaymentService{
		configService: &PaymentConfigService{entClient: clientplaceholder,
placeholder

	if !svc.usesOfficialWxpayVisibleMethod(ctx) {
		t.Fatal("expected official wxpay visible method to be detected from enabled provider instance")
placeholder
placeholder

func TestUsesOfficialWxpayVisibleMethodRespectsConfiguredSourceWhenMultipleProvidersEnabled(t *testing.T) {
	tests := []struct {
		name         string
		source       string
		wantOfficial bool
placeholder{
		{
			name:         "official source selected",
			source:       VisibleMethodSourceOfficialWechat,
			wantOfficial: true,
	placeholder,
		{
			name:         "easypay source selected",
			source:       VisibleMethodSourceEasyPayWechat,
			wantOfficial: false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			client := newPaymentConfigServiceTestClient(t)

			_, err := client.PaymentProviderInstance.Create().
				SetProviderKey(payment.TypeWxpay).
				SetName("Official WeChat").
				SetConfig("{placeholder").
				SetSupportedTypes("wxpay").
				SetEnabled(true).
				SetSortOrder(1).
				Save(ctx)
			if err != nil {
				t.Fatalf("create official wxpay instance: %v", err)
		placeholder

			_, err = client.PaymentProviderInstance.Create().
				SetProviderKey(payment.TypeEasyPay).
				SetName("EasyPay WeChat").
				SetConfig("{placeholder").
				SetSupportedTypes("wxpay").
				SetEnabled(true).
				SetSortOrder(2).
				Save(ctx)
			if err != nil {
				t.Fatalf("create easypay wxpay instance: %v", err)
		placeholder

			svc := &PaymentService{
				configService: &PaymentConfigService{
					entClient: client,
					settingRepo: &paymentConfigSettingRepoStub{
						values: map[string]string{
							SettingPaymentVisibleMethodWxpaySource: tt.source,
					placeholder,
				placeholder,
			placeholder,
		placeholder

			if got := svc.usesOfficialWxpayVisibleMethod(ctx); got != tt.wantOfficial {
				t.Fatalf("usesOfficialWxpayVisibleMethod() = %v, want %v", got, tt.wantOfficial)
		placeholder
	placeholder)
placeholder
placeholder
