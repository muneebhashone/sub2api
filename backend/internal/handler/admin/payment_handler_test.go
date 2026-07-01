package admin

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

func TestSanitizeAdminPaymentOrderForResponseAddsCurrency(t *testing.T) {
	now := time.Now()
	order := &dbent.PaymentOrder{
		ID:          1,
		UserID:      2,
		Amount:      100,
		PayAmount:   108,
		FeeRate:     8,
		OutTradeNo:  "sub2_202606250001",
		PaymentType: "stripe",
		OrderType:   "subscription",
		Status:      "COMPLETED",
		ExpiresAt:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
		ProviderSnapshot: map[string]any{
			"schema_version": 2,
			"currency":       "USD",
	placeholder,
placeholder

	got := sanitizeAdminPaymentOrderForResponse(order)
	if got == nil {
		t.Fatal("expected sanitized order")
placeholder
	if got.Currency != "USD" {
		t.Fatalf("expected currency USD, got %q", got.Currency)
placeholder

	body, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal sanitized order: %v", err)
placeholder
	if strings.Contains(string(body), "provider_snapshot") {
		t.Fatalf("expected provider_snapshot to be omitted, got %s", string(body))
placeholder
placeholder

func TestAdminSubscriptionPlansForResponseIncludesCompositeGroupInfo(t *testing.T) {
	weekly := 25.0
	plans := []*dbent.SubscriptionPlan{
		{
			ID:           11,
			GroupID:      7,
			Name:         "All models",
			Description:  "Composite access",
			Price:        19.99,
			ValidityDays: 30,
			ValidityUnit: "days",
			Features:     "OpenAI\nClaude\nGemini\nGrok",
			ProductName:  "Sub2API",
			ForSale:      true,
			SortOrder:    1,
	placeholder,
placeholder
	groupInfo := map[int64]service.PlanGroupInfo{
		7: {
			Platform:       service.PlatformComposite,
			Name:           "Bucket 2 composite",
			RateMultiplier: 1.5,
			WeeklyLimitUSD: &weekly,
			ModelScopes:    []string{"openai", "claude", "gemini", "grok"placeholder,
	placeholder,
placeholder

	got := adminSubscriptionPlansForResponse(plans, groupInfo)

	if len(got) != 1 {
		t.Fatalf("expected one plan, got %d", len(got))
placeholder
	if got[0].GroupPlatform != service.PlatformComposite {
		t.Fatalf("expected composite group platform, got %q", got[0].GroupPlatform)
placeholder
	if got[0].GroupName != "Bucket 2 composite" {
		t.Fatalf("expected group name to be included, got %q", got[0].GroupName)
placeholder
	if got[0].WeeklyLimitUSD == nil || *got[0].WeeklyLimitUSD != weekly {
		t.Fatalf("expected weekly limit to be included, got %#v", got[0].WeeklyLimitUSD)
placeholder
	if strings.Join(got[0].ModelScopes, ",") != "openai,claude,gemini,grok" {
		t.Fatalf("expected model scopes to be preserved, got %#v", got[0].ModelScopes)
placeholder
placeholder
