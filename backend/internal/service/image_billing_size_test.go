package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClassifyImageBillingTier(t *testing.T) {
	tests := []struct {
		name     string
		size     string
		wantTier string
		wantOK   bool
placeholder{
		{name: "explicit 2k square", size: "2048x2048", wantTier: "2K", wantOK: trueplaceholder,
		{name: "explicit 2k landscape", size: "2048x1152", wantTier: "2K", wantOK: trueplaceholder,
		{name: "explicit 4k landscape", size: "3840x2160", wantTier: "4K", wantOK: trueplaceholder,
		{name: "explicit 4k portrait", size: "2160x3840", wantTier: "4K", wantOK: trueplaceholder,
		{name: "long edge 1k", size: "1024X768", wantTier: "1K", wantOK: trueplaceholder,
		{name: "long edge 2k", size: "1280x768", wantTier: "2K", wantOK: trueplaceholder,
		{name: "long edge 4k", size: "2560x1600", wantTier: "4K", wantOK: trueplaceholder,
		{name: "tier string 1k", size: "1k", wantTier: "1K", wantOK: trueplaceholder,
		{name: "empty", size: "", wantOK: falseplaceholder,
		{name: "auto", size: "auto", wantOK: falseplaceholder,
		{name: "invalid", size: "not-a-size", wantOK: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTier, gotOK := ClassifyImageBillingTier(tt.size)
			require.Equal(t, tt.wantOK, gotOK)
			require.Equal(t, tt.wantTier, gotTier)
	placeholder)
placeholder
placeholder

func TestResolveImageBillingSize(t *testing.T) {
	tests := []struct {
		name          string
		inputSize     string
		outputSizes   []string
		wantBilling   string
		wantOutput    string
		wantSource    string
		wantBreakdown map[string]int
placeholder{
		{
			name:          "output wins over input",
			inputSize:     "1024x1024",
			outputSizes:   []string{"3840x2160"placeholder,
			wantBilling:   "4K",
			wantOutput:    "3840x2160",
			wantSource:    ImageSizeSourceOutput,
			wantBreakdown: map[string]int{"4K": 1placeholder,
	placeholder,
		{
			name:        "input fallback",
			inputSize:   "1024x1024",
			wantBilling: "1K",
			wantSource:  ImageSizeSourceInput,
	placeholder,
		{
			name:        "auto defaults",
			inputSize:   "auto",
			wantBilling: "2K",
			wantSource:  ImageSizeSourceDefault,
	placeholder,
		{
			name:        "empty defaults",
			inputSize:   "",
			wantBilling: "2K",
			wantSource:  ImageSizeSourceDefault,
	placeholder,
		{
			name:        "invalid defaults",
			inputSize:   "largest",
			wantBilling: "2K",
			wantSource:  ImageSizeSourceDefault,
	placeholder,
		{
			name:          "mixed output chooses highest tier",
			inputSize:     "1024x1024",
			outputSizes:   []string{"1024x1024", "3840x2160", "1280x720"placeholder,
			wantBilling:   "4K",
			wantOutput:    "1024x1024",
			wantSource:    ImageSizeSourceOutput,
			wantBreakdown: map[string]int{"1K": 1, "2K": 1, "4K": 1placeholder,
	placeholder,
		{
			name:        "unparseable output falls back to parseable input",
			inputSize:   "2048x1152",
			outputSizes: []string{"auto"placeholder,
			wantBilling: "2K",
			wantOutput:  "auto",
			wantSource:  ImageSizeSourceInput,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveImageBillingSize(tt.inputSize, tt.outputSizes)
			require.Equal(t, tt.wantBilling, got.BillingSize)
			require.Equal(t, tt.inputSize, got.InputSize)
			require.Equal(t, tt.wantOutput, got.OutputSize)
			require.Equal(t, tt.wantSource, got.Source)
			require.Equal(t, tt.wantBreakdown, got.Breakdown)
	placeholder)
placeholder
placeholder
