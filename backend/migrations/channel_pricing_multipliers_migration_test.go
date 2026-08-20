package migrations

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChannelPricingMultipliersMigration(t *testing.T) {
	content, err := FS.ReadFile("228_channel_pricing_multipliers.sql")
placeholder

	sql := strings.Join(strings.Fields(string(content)), " ")
	for _, column := range []string{
		"fast_multiplier NUMERIC(12,6)",
		"flex_multiplier NUMERIC(12,6)",
		"input_multiplier NUMERIC(12,6)",
		"output_multiplier NUMERIC(12,6)",
		"cache_write_multiplier NUMERIC(12,6)",
		"cache_read_multiplier NUMERIC(12,6)",
placeholder {
		require.Contains(t, sql, "ADD COLUMN IF NOT EXISTS "+column)
placeholder

	constraints := []struct {
		table  string
		name   string
		column string
placeholder{
		{"channel_model_pricing", "channel_model_pricing_fast_multiplier_positive", "fast_multiplier"placeholder,
		{"channel_model_pricing", "channel_model_pricing_flex_multiplier_positive", "flex_multiplier"placeholder,
		{"channel_pricing_intervals", "channel_pricing_intervals_input_multiplier_positive", "input_multiplier"placeholder,
		{"channel_pricing_intervals", "channel_pricing_intervals_output_multiplier_positive", "output_multiplier"placeholder,
		{"channel_pricing_intervals", "channel_pricing_intervals_cache_write_multiplier_positive", "cache_write_multiplier"placeholder,
		{"channel_pricing_intervals", "channel_pricing_intervals_cache_read_multiplier_positive", "cache_read_multiplier"placeholder,
placeholder
	for _, constraint := range constraints {
		require.Contains(t, sql, "conname = '"+constraint.name+"' AND conrelid = '"+constraint.table+"'::regclass")
		require.Contains(t, sql, "ALTER TABLE "+constraint.table+" ADD CONSTRAINT "+constraint.name+
			" CHECK ("+constraint.column+" IS NULL OR "+constraint.column+" > 0)")
placeholder
placeholder
