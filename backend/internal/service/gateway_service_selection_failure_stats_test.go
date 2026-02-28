package service

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestCollectSelectionFailureStats(t *testing.T) {
	svc := &GatewayService{placeholder
	model := "sora2-landscape-10s"
	resetAt := time.Now().Add(2 * time.Minute).Format(time.RFC3339)

	accounts := []Account{
		// excluded
		{
			ID:          1,
			Platform:    PlatformSora,
			Status:      StatusActive,
			Schedulable: true,
	placeholder,
		// unschedulable
		{
			ID:          2,
			Platform:    PlatformSora,
			Status:      StatusActive,
			Schedulable: false,
	placeholder,
		// platform filtered
		{
			ID:          3,
			Platform:    PlatformOpenAI,
			Status:      StatusActive,
			Schedulable: true,
	placeholder,
		// model unsupported
		{
			ID:          4,
			Platform:    PlatformSora,
			Status:      StatusActive,
			Schedulable: true,
	placeholder
				"model_mapping": map[string]any{
					"gpt-image": "gpt-image",
			placeholder,
		placeholder,
	placeholder,
		// model rate limited
		{
			ID:          5,
			Platform:    PlatformSora,
			Status:      StatusActive,
			Schedulable: true,
			Extra: map[string]any{
				"model_rate_limits": map[string]any{
					model: map[string]any{
						"rate_limit_reset_at": resetAt,
				placeholder,
			placeholder,
		placeholder,
	placeholder,
		// eligible
		{
			ID:          6,
			Platform:    PlatformSora,
			Status:      StatusActive,
			Schedulable: true,
	placeholder,
placeholder

	excluded := map[int64]struct{placeholder{1: {placeholderplaceholder
	stats := svc.collectSelectionFailureStats(context.Background(), accounts, model, PlatformSora, excluded, false)

	if stats.Total != 6 {
		t.Fatalf("total=%d want=6", stats.Total)
placeholder
	if stats.Excluded != 1 {
		t.Fatalf("excluded=%d want=1", stats.Excluded)
placeholder
	if stats.Unschedulable != 1 {
		t.Fatalf("unschedulable=%d want=1", stats.Unschedulable)
placeholder
	if stats.PlatformFiltered != 1 {
		t.Fatalf("platform_filtered=%d want=1", stats.PlatformFiltered)
placeholder
	if stats.ModelUnsupported != 1 {
		t.Fatalf("model_unsupported=%d want=1", stats.ModelUnsupported)
placeholder
	if stats.ModelRateLimited != 1 {
		t.Fatalf("model_rate_limited=%d want=1", stats.ModelRateLimited)
placeholder
	if stats.Eligible != 1 {
		t.Fatalf("eligible=%d want=1", stats.Eligible)
placeholder
placeholder

func TestDiagnoseSelectionFailure_SoraUnschedulableDetail(t *testing.T) {
	svc := &GatewayService{placeholder
	acc := &Account{
		ID:          7,
		Platform:    PlatformSora,
		Status:      StatusActive,
		Schedulable: false,
placeholder

	diagnosis := svc.diagnoseSelectionFailure(context.Background(), acc, "sora2-landscape-10s", PlatformSora, map[int64]struct{placeholder{placeholder, false)
	if diagnosis.Category != "unschedulable" {
		t.Fatalf("category=%s want=unschedulable", diagnosis.Category)
placeholder
	if diagnosis.Detail != "schedulable=false" {
		t.Fatalf("detail=%s want=schedulable=false", diagnosis.Detail)
placeholder
placeholder

func TestDiagnoseSelectionFailure_SoraModelRateLimitedDetail(t *testing.T) {
	svc := &GatewayService{placeholder
	model := "sora2-landscape-10s"
	resetAt := time.Now().Add(2 * time.Minute).UTC().Format(time.RFC3339)
	acc := &Account{
		ID:          8,
		Platform:    PlatformSora,
		Status:      StatusActive,
		Schedulable: true,
		Extra: map[string]any{
			"model_rate_limits": map[string]any{
				model: map[string]any{
					"rate_limit_reset_at": resetAt,
			placeholder,
		placeholder,
	placeholder,
placeholder

	diagnosis := svc.diagnoseSelectionFailure(context.Background(), acc, model, PlatformSora, map[int64]struct{placeholder{placeholder, false)
	if diagnosis.Category != "model_rate_limited" {
		t.Fatalf("category=%s want=model_rate_limited", diagnosis.Category)
placeholder
	if !strings.Contains(diagnosis.Detail, "remaining=") {
		t.Fatalf("detail=%s want contains remaining=", diagnosis.Detail)
placeholder
placeholder
