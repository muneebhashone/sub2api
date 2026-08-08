package service

import (
	"log/slog"
	"sort"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
)

// Canonical video price family keys used in groups.video_model_prices JSONB.
const (
	VideoPriceFamilyGrokImagineVideo   = "grok-imagine-video"
	VideoPriceFamilyGrokImagineVideo15 = "grok-imagine-video-1.5"
)

// CanonicalGrokImagineVideoPriceFamily normalizes model aliases / preview / legacy
// IDs onto the price-family keys stored in video_model_prices.
func CanonicalGrokImagineVideoPriceFamily(model string) string {
	if model == "" {
		return ""
placeholder
	// Prefer shared xAI helper for known aliases. Keep future native Imagine
	// models distinct so operators can assign them independent prices.
	if c := xai.CanonicalImagineVideoModel(model); c != "" {
		switch c {
		case xai.DefaultImagineVideo15Model:
			return VideoPriceFamilyGrokImagineVideo15
		case xai.DefaultImagineVideoModel:
			return VideoPriceFamilyGrokImagineVideo
	placeholder
		if strings.HasPrefix(c, "grok-imagine-video-") {
			return c
	placeholder
placeholder
	m := strings.ToLower(strings.TrimSpace(model))
	for _, prefix := range []string{"xai/", "x-ai/", "grok/"placeholder {
		if strings.HasPrefix(m, prefix) {
			m = strings.TrimPrefix(m, prefix)
			break
	placeholder
placeholder
	switch {
	case m == "grok-imagine-video-1.5" || m == "grok-imagine-video-1.5-preview" ||
		m == "grok-video-1.5" || strings.Contains(m, "video-1.5"):
		return VideoPriceFamilyGrokImagineVideo15
	case m == "grok-imagine-video" || m == "grok-imagine-video-preview" ||
		m == "grok-video" || m == "grok-video-latest":
		return VideoPriceFamilyGrokImagineVideo
	default:
		return ""
placeholder
placeholder

// NormalizeVideoModelPrices cleans and canonicalizes a per-model resolution map.
// Keys become price families; tiers use 480p/720p/1080p. Negative prices dropped.
//
// Model keys are walked in sorted order rather than in Go map order: several
// aliases can canonicalize onto the same family, and an unordered walk would
// make the winning price for a conflicting tier vary between processes.
// Unrecognized tiers are dropped with a warning instead of silently collapsing
// into the 480p bucket.
func NormalizeVideoModelPrices(in map[string]map[string]float64) map[string]map[string]float64 {
	if len(in) == 0 {
		return nil
placeholder
	modelKeys := make([]string, 0, len(in))
	for modelKey := range in {
		modelKeys = append(modelKeys, modelKey)
placeholder
	sort.Strings(modelKeys)
	out := make(map[string]map[string]float64)
	for _, modelKey := range modelKeys {
		tierPrices := in[modelKey]
		if len(tierPrices) == 0 {
			continue
	placeholder
		family := CanonicalGrokImagineVideoPriceFamily(modelKey)
		if family == "" {
			key := strings.ToLower(strings.TrimSpace(modelKey))
			switch key {
			case VideoPriceFamilyGrokImagineVideo, VideoPriceFamilyGrokImagineVideo15:
				family = key
			default:
				if key == "" {
					continue
			placeholder
				family = key
		placeholder
	placeholder
		normalizedTiers := out[family]
		if normalizedTiers == nil {
			normalizedTiers = make(map[string]float64)
	placeholder
		tierKeys := make([]string, 0, len(tierPrices))
		for tierKey := range tierPrices {
			tierKeys = append(tierKeys, tierKey)
	placeholder
		sort.Strings(tierKeys)
		for _, tierKey := range tierKeys {
			price := tierPrices[tierKey]
			if price < 0 {
				continue
		placeholder
			tier, ok := LookupVideoBillingResolution(tierKey)
			if !ok {
				slog.Warn("video_model_prices_unknown_resolution_dropped",
					"model_key", modelKey,
					"family", family,
					"resolution", tierKey)
				continue
		placeholder
			if existing, exists := normalizedTiers[tier]; exists && existing != price {
				slog.Warn("video_model_prices_conflicting_tier_price",
					"model_key", modelKey,
					"family", family,
					"resolution", tier,
					"previous_price", existing,
					"price", price)
		placeholder
			normalizedTiers[tier] = price
	placeholder
		if len(normalizedTiers) > 0 {
			out[family] = normalizedTiers
	placeholder
placeholder
	if len(out) == 0 {
		return nil
placeholder
	return out
placeholder

// LookupVideoModelPrice returns a per-second price from a model×resolution map, or nil.
func LookupVideoModelPrice(prices map[string]map[string]float64, model, resolution string) *float64 {
	if len(prices) == 0 {
		return nil
placeholder
	family := CanonicalGrokImagineVideoPriceFamily(model)
	if family == "" {
		family = strings.ToLower(strings.TrimSpace(model))
placeholder
	if family == "" {
		return nil
placeholder
	tierPrices, ok := prices[family]
	if !ok || len(tierPrices) == 0 {
		return nil
placeholder
	tier := NormalizeVideoBillingResolutionOrDefault(resolution)
	if price, ok := tierPrices[tier]; ok {
		p := price
		return &p
placeholder
	return nil
placeholder
