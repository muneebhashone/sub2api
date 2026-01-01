package service

import "testing"

func TestInferGoogleOneTier(t *testing.T) {
	tests := []struct {
		name         string
		storageBytes int64
		expectedTier string
placeholder{
		{"Negative storage", -1, TierGoogleOneUnknownplaceholder,
		{"Zero storage", 0, TierGoogleOneUnknownplaceholder,

		// Free tier boundary (15GB)
		{"Below free tier", 10 * GB, TierGoogleOneUnknownplaceholder,
		{"Just below free tier", StorageTierFree - 1, TierGoogleOneUnknownplaceholder,
		{"Free tier (15GB)", StorageTierFree, TierFreeplaceholder,

		// Basic tier boundary (100GB)
		{"Between free and basic", 50 * GB, TierFreeplaceholder,
		{"Just below basic tier", StorageTierBasic - 1, TierFreeplaceholder,
		{"Basic tier (100GB)", StorageTierBasic, TierGoogleOneBasicplaceholder,

		// Standard tier boundary (200GB)
		{"Between basic and standard", 150 * GB, TierGoogleOneBasicplaceholder,
		{"Just below standard tier", StorageTierStandard - 1, TierGoogleOneBasicplaceholder,
		{"Standard tier (200GB)", StorageTierStandard, TierGoogleOneStandardplaceholder,

		// AI Premium tier boundary (2TB)
		{"Between standard and premium", 1 * TB, TierGoogleOneStandardplaceholder,
		{"Just below AI Premium tier", StorageTierAIPremium - 1, TierGoogleOneStandardplaceholder,
		{"AI Premium tier (2TB)", StorageTierAIPremium, TierAIPremiumplaceholder,

		// Unlimited tier boundary (> 100TB)
		{"Between premium and unlimited", 50 * TB, TierAIPremiumplaceholder,
		{"At unlimited threshold (100TB)", StorageTierUnlimited, TierAIPremiumplaceholder,
		{"Unlimited tier (100TB+)", StorageTierUnlimited + 1, TierGoogleOneUnlimitedplaceholder,
		{"Unlimited tier (101TB+)", 101 * TB, TierGoogleOneUnlimitedplaceholder,
		{"Very large storage", 1000 * TB, TierGoogleOneUnlimitedplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inferGoogleOneTier(tt.storageBytes)
			if result != tt.expectedTier {
				t.Errorf("inferGoogleOneTier(%d) = %s, want %s",
					tt.storageBytes, result, tt.expectedTier)
		placeholder
	placeholder)
placeholder
placeholder
