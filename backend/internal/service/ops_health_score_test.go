//go:build unit

package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestComputeDashboardHealthScore_IdleReturns100(t *testing.T) {
	t.Parallel()

	score := computeDashboardHealthScore(time.Now().UTC(), &OpsDashboardOverview{placeholder)
	require.Equal(t, 100, score)
placeholder

func TestComputeDashboardHealthScore_DegradesOnBadSignals(t *testing.T) {
	t.Parallel()

	ov := &OpsDashboardOverview{
		RequestCountTotal: 100,
		RequestCountSLA:   100,
		SuccessCount:      90,
		ErrorCountTotal:   10,
		ErrorCountSLA:     10,

		SLA:               0.90,
		ErrorRate:         0.10,
		UpstreamErrorRate: 0.08,

		Duration: OpsPercentiles{P99: intPtr(20_000)placeholder,
		TTFT:     OpsPercentiles{P99: intPtr(2_000)placeholder,

		SystemMetrics: &OpsSystemMetricsSnapshot{
			DBOK:                  boolPtr(false),
			RedisOK:               boolPtr(false),
			CPUUsagePercent:       float64Ptr(98.0),
			MemoryUsagePercent:    float64Ptr(97.0),
			DBConnWaiting:         intPtr(3),
			ConcurrencyQueueDepth: intPtr(10),
	placeholder,
		JobHeartbeats: []*OpsJobHeartbeat{
			{
				JobName:     "job-a",
				LastErrorAt: timePtr(time.Now().UTC().Add(-1 * time.Minute)),
				LastError:   stringPtr("boom"),
		placeholder,
	placeholder,
placeholder

	score := computeDashboardHealthScore(time.Now().UTC(), ov)
	require.Less(t, score, 80)
	require.GreaterOrEqual(t, score, 0)
placeholder

func TestComputeDashboardHealthScore_Comprehensive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		overview *OpsDashboardOverview
		wantMin  int
		wantMax  int
placeholder{
		{
			name:     "nil overview returns 0",
			overview: nil,
			wantMin:  0,
			wantMax:  0,
	placeholder,
		{
			name: "perfect health",
			overview: &OpsDashboardOverview{
				RequestCountTotal: 1000,
				RequestCountSLA:   1000,
				SLA:               1.0,
				ErrorRate:         0,
				UpstreamErrorRate: 0,
				Duration:          OpsPercentiles{P99: intPtr(500)placeholder,
				TTFT:              OpsPercentiles{P99: intPtr(100)placeholder,
				SystemMetrics: &OpsSystemMetricsSnapshot{
					DBOK:               boolPtr(true),
					RedisOK:            boolPtr(true),
					CPUUsagePercent:    float64Ptr(30),
					MemoryUsagePercent: float64Ptr(40),
			placeholder,
		placeholder,
			wantMin: 100,
			wantMax: 100,
	placeholder,
		{
			name: "good health - SLA 99.8%",
			overview: &OpsDashboardOverview{
				RequestCountTotal: 1000,
				RequestCountSLA:   1000,
				SLA:               0.998,
				ErrorRate:         0.003,
				UpstreamErrorRate: 0.001,
				Duration:          OpsPercentiles{P99: intPtr(800)placeholder,
				TTFT:              OpsPercentiles{P99: intPtr(200)placeholder,
				SystemMetrics: &OpsSystemMetricsSnapshot{
					DBOK:               boolPtr(true),
					RedisOK:            boolPtr(true),
					CPUUsagePercent:    float64Ptr(50),
					MemoryUsagePercent: float64Ptr(60),
			placeholder,
		placeholder,
			wantMin: 95,
			wantMax: 100,
	placeholder,
		{
			name: "medium health - SLA 96%",
			overview: &OpsDashboardOverview{
				RequestCountTotal: 1000,
				RequestCountSLA:   1000,
				SLA:               0.96,
				ErrorRate:         0.02,
				UpstreamErrorRate: 0.01,
				Duration:          OpsPercentiles{P99: intPtr(3000)placeholder,
				TTFT:              OpsPercentiles{P99: intPtr(600)placeholder,
				SystemMetrics: &OpsSystemMetricsSnapshot{
					DBOK:               boolPtr(true),
					RedisOK:            boolPtr(true),
					CPUUsagePercent:    float64Ptr(70),
					MemoryUsagePercent: float64Ptr(75),
			placeholder,
		placeholder,
			wantMin: 96,
			wantMax: 97,
	placeholder,
		{
			name: "DB failure",
			overview: &OpsDashboardOverview{
				RequestCountTotal: 1000,
				RequestCountSLA:   1000,
				SLA:               0.995,
				ErrorRate:         0,
				UpstreamErrorRate: 0,
				Duration:          OpsPercentiles{P99: intPtr(500)placeholder,
				SystemMetrics: &OpsSystemMetricsSnapshot{
					DBOK:               boolPtr(false),
					RedisOK:            boolPtr(true),
					CPUUsagePercent:    float64Ptr(30),
					MemoryUsagePercent: float64Ptr(40),
			placeholder,
		placeholder,
			wantMin: 70,
			wantMax: 90,
	placeholder,
		{
			name: "Redis failure",
			overview: &OpsDashboardOverview{
				RequestCountTotal: 1000,
				RequestCountSLA:   1000,
				SLA:               0.995,
				ErrorRate:         0,
				UpstreamErrorRate: 0,
				Duration:          OpsPercentiles{P99: intPtr(500)placeholder,
				SystemMetrics: &OpsSystemMetricsSnapshot{
					DBOK:               boolPtr(true),
					RedisOK:            boolPtr(false),
					CPUUsagePercent:    float64Ptr(30),
					MemoryUsagePercent: float64Ptr(40),
			placeholder,
		placeholder,
			wantMin: 85,
			wantMax: 95,
	placeholder,
		{
			name: "high CPU usage",
			overview: &OpsDashboardOverview{
				RequestCountTotal: 1000,
				RequestCountSLA:   1000,
				SLA:               0.995,
				ErrorRate:         0,
				UpstreamErrorRate: 0,
				Duration:          OpsPercentiles{P99: intPtr(500)placeholder,
				SystemMetrics: &OpsSystemMetricsSnapshot{
					DBOK:               boolPtr(true),
					RedisOK:            boolPtr(true),
					CPUUsagePercent:    float64Ptr(95),
					MemoryUsagePercent: float64Ptr(40),
			placeholder,
		placeholder,
			wantMin: 85,
			wantMax: 100,
	placeholder,
		{
			name: "combined failures - business degraded + infra healthy",
			overview: &OpsDashboardOverview{
				RequestCountTotal: 1000,
				RequestCountSLA:   1000,
				SLA:               0.90,
				ErrorRate:         0.05,
				UpstreamErrorRate: 0.02,
				Duration:          OpsPercentiles{P99: intPtr(10000)placeholder,
				SystemMetrics: &OpsSystemMetricsSnapshot{
					DBOK:               boolPtr(true),
					RedisOK:            boolPtr(true),
					CPUUsagePercent:    float64Ptr(20),
					MemoryUsagePercent: float64Ptr(30),
			placeholder,
		placeholder,
			wantMin: 84,
			wantMax: 85,
	placeholder,
		{
			name: "combined failures - business healthy + infra degraded",
			overview: &OpsDashboardOverview{
				RequestCountTotal: 1000,
				RequestCountSLA:   1000,
				SLA:               0.998,
				ErrorRate:         0.001,
				UpstreamErrorRate: 0,
				Duration:          OpsPercentiles{P99: intPtr(600)placeholder,
				SystemMetrics: &OpsSystemMetricsSnapshot{
					DBOK:               boolPtr(false),
					RedisOK:            boolPtr(false),
					CPUUsagePercent:    float64Ptr(95),
					MemoryUsagePercent: float64Ptr(95),
			placeholder,
		placeholder,
			wantMin: 70,
			wantMax: 90,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := computeDashboardHealthScore(time.Now().UTC(), tt.overview)
			require.GreaterOrEqual(t, score, tt.wantMin, "score should be >= %d", tt.wantMin)
			require.LessOrEqual(t, score, tt.wantMax, "score should be <= %d", tt.wantMax)
			require.GreaterOrEqual(t, score, 0, "score must be >= 0")
			require.LessOrEqual(t, score, 100, "score must be <= 100")
	placeholder)
placeholder
placeholder

func TestComputeBusinessHealth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		overview *OpsDashboardOverview
		wantMin  float64
		wantMax  float64
placeholder{
		{
			name: "perfect metrics",
			overview: &OpsDashboardOverview{
				SLA:               1.0,
				ErrorRate:         0,
				UpstreamErrorRate: 0,
				Duration:          OpsPercentiles{P99: intPtr(500)placeholder,
		placeholder,
			wantMin: 100,
			wantMax: 100,
	placeholder,
		{
			name: "SLA boundary 99.5%",
			overview: &OpsDashboardOverview{
				SLA:               0.995,
				ErrorRate:         0,
				UpstreamErrorRate: 0,
				Duration:          OpsPercentiles{P99: intPtr(500)placeholder,
		placeholder,
			wantMin: 100,
			wantMax: 100,
	placeholder,
		{
			name: "SLA boundary 95%",
			overview: &OpsDashboardOverview{
				SLA:               0.95,
				ErrorRate:         0,
				UpstreamErrorRate: 0,
				Duration:          OpsPercentiles{P99: intPtr(500)placeholder,
		placeholder,
			wantMin: 100,
			wantMax: 100,
	placeholder,
		{
			name: "error rate boundary 1%",
			overview: &OpsDashboardOverview{
				SLA:               0.99,
				ErrorRate:         0.01,
				UpstreamErrorRate: 0,
				Duration:          OpsPercentiles{P99: intPtr(500)placeholder,
		placeholder,
			wantMin: 100,
			wantMax: 100,
	placeholder,
		{
			name: "error rate 5%",
			overview: &OpsDashboardOverview{
				SLA:               0.95,
				ErrorRate:         0.05,
				UpstreamErrorRate: 0,
				Duration:          OpsPercentiles{P99: intPtr(500)placeholder,
		placeholder,
			wantMin: 77,
			wantMax: 78,
	placeholder,
		{
			name: "TTFT boundary 2s",
			overview: &OpsDashboardOverview{
				SLA:               0.99,
				ErrorRate:         0,
				UpstreamErrorRate: 0,
				TTFT:              OpsPercentiles{P99: intPtr(2000)placeholder,
		placeholder,
			wantMin: 75,
			wantMax: 75,
	placeholder,
		{
			name: "upstream error dominates",
			overview: &OpsDashboardOverview{
				SLA:               0.995,
				ErrorRate:         0.001,
				UpstreamErrorRate: 0.03,
				Duration:          OpsPercentiles{P99: intPtr(500)placeholder,
		placeholder,
			wantMin: 88,
			wantMax: 90,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := computeBusinessHealth(tt.overview)
			require.GreaterOrEqual(t, score, tt.wantMin, "score should be >= %.1f", tt.wantMin)
			require.LessOrEqual(t, score, tt.wantMax, "score should be <= %.1f", tt.wantMax)
			require.GreaterOrEqual(t, score, 0.0, "score must be >= 0")
			require.LessOrEqual(t, score, 100.0, "score must be <= 100")
	placeholder)
placeholder
placeholder

func TestComputeInfraHealth(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()

	tests := []struct {
		name     string
		overview *OpsDashboardOverview
		wantMin  float64
		wantMax  float64
placeholder{
		{
			name: "all infrastructure healthy",
			overview: &OpsDashboardOverview{
				RequestCountTotal: 1000,
				SystemMetrics: &OpsSystemMetricsSnapshot{
					DBOK:               boolPtr(true),
					RedisOK:            boolPtr(true),
					CPUUsagePercent:    float64Ptr(30),
					MemoryUsagePercent: float64Ptr(40),
			placeholder,
		placeholder,
			wantMin: 100,
			wantMax: 100,
	placeholder,
		{
			name: "DB down",
			overview: &OpsDashboardOverview{
				RequestCountTotal: 1000,
				SystemMetrics: &OpsSystemMetricsSnapshot{
					DBOK:               boolPtr(false),
					RedisOK:            boolPtr(true),
					CPUUsagePercent:    float64Ptr(30),
					MemoryUsagePercent: float64Ptr(40),
			placeholder,
		placeholder,
			wantMin: 50,
			wantMax: 70,
	placeholder,
		{
			name: "Redis down",
			overview: &OpsDashboardOverview{
				RequestCountTotal: 1000,
				SystemMetrics: &OpsSystemMetricsSnapshot{
					DBOK:               boolPtr(true),
					RedisOK:            boolPtr(false),
					CPUUsagePercent:    float64Ptr(30),
					MemoryUsagePercent: float64Ptr(40),
			placeholder,
		placeholder,
			wantMin: 80,
			wantMax: 95,
	placeholder,
		{
			name: "CPU at 90%",
			overview: &OpsDashboardOverview{
				RequestCountTotal: 1000,
				SystemMetrics: &OpsSystemMetricsSnapshot{
					DBOK:               boolPtr(true),
					RedisOK:            boolPtr(true),
					CPUUsagePercent:    float64Ptr(90),
					MemoryUsagePercent: float64Ptr(40),
			placeholder,
		placeholder,
			wantMin: 85,
			wantMax: 95,
	placeholder,
		{
			name: "failed background job",
			overview: &OpsDashboardOverview{
				RequestCountTotal: 1000,
				SystemMetrics: &OpsSystemMetricsSnapshot{
					DBOK:               boolPtr(true),
					RedisOK:            boolPtr(true),
					CPUUsagePercent:    float64Ptr(30),
					MemoryUsagePercent: float64Ptr(40),
			placeholder,
				JobHeartbeats: []*OpsJobHeartbeat{
					{
						JobName:     "test-job",
						LastErrorAt: &now,
				placeholder,
			placeholder,
		placeholder,
			wantMin: 70,
			wantMax: 90,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := computeInfraHealth(now, tt.overview)
			require.GreaterOrEqual(t, score, tt.wantMin, "score should be >= %.1f", tt.wantMin)
			require.LessOrEqual(t, score, tt.wantMax, "score should be <= %.1f", tt.wantMax)
			require.GreaterOrEqual(t, score, 0.0, "score must be >= 0")
			require.LessOrEqual(t, score, 100.0, "score must be <= 100")
	placeholder)
placeholder
placeholder

func timePtr(v time.Time) *time.Time { return &v placeholder

func stringPtr(v string) *string { return &v placeholder
