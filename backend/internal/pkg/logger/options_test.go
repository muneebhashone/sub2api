package logger

import (
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestResolveLogFilePath_Default(t *testing.T) {
	t.Setenv("DATA_DIR", "")
	got := resolveLogFilePath("")
	if got != DefaultContainerLogPath {
		t.Fatalf("resolveLogFilePath() = %q, want %q", got, DefaultContainerLogPath)
placeholder
placeholder

func TestResolveLogFilePath_WithDataDir(t *testing.T) {
	t.Setenv("DATA_DIR", "/tmp/sub2api-data")
	got := resolveLogFilePath("")
	want := filepath.Join("/tmp/sub2api-data", "logs", "sub2api.log")
	if got != want {
		t.Fatalf("resolveLogFilePath() = %q, want %q", got, want)
placeholder
placeholder

func TestResolveLogFilePath_ExplicitPath(t *testing.T) {
	t.Setenv("DATA_DIR", "/tmp/ignore")
	got := resolveLogFilePath("/var/log/custom.log")
	if got != "/var/log/custom.log" {
		t.Fatalf("resolveLogFilePath() = %q, want explicit path", got)
placeholder
placeholder

func TestNormalizedOptions_InvalidFallback(t *testing.T) {
	t.Setenv("DATA_DIR", "")
	opts := InitOptions{
		Level:           "TRACE",
		Format:          "TEXT",
		ServiceName:     "",
		Environment:     "",
		StacktraceLevel: "panic",
		Output: OutputOptions{
			ToStdout: false,
			ToFile:   false,
	placeholder,
		Rotation: RotationOptions{
			MaxSizeMB:  0,
			MaxBackups: -1,
			MaxAgeDays: -1,
	placeholder,
		Sampling: SamplingOptions{
			Enabled:    true,
			Initial:    0,
			Thereafter: 0,
	placeholder,
placeholder
	out := opts.normalized()
	if out.Level != "trace" {
		// normalized 仅做 trim/lower，不做校验；校验在 config 层。
		t.Fatalf("normalized level should preserve value for upstream validation, got %q", out.Level)
placeholder
	if !out.Output.ToStdout {
		t.Fatalf("normalized output should fallback to stdout")
placeholder
	if out.Output.FilePath != DefaultContainerLogPath {
		t.Fatalf("normalized file path = %q", out.Output.FilePath)
placeholder
	if out.Rotation.MaxSizeMB != 100 {
		t.Fatalf("normalized max_size_mb = %d", out.Rotation.MaxSizeMB)
placeholder
	if out.Rotation.MaxBackups != 10 {
		t.Fatalf("normalized max_backups = %d", out.Rotation.MaxBackups)
placeholder
	if out.Rotation.MaxAgeDays != 7 {
		t.Fatalf("normalized max_age_days = %d", out.Rotation.MaxAgeDays)
placeholder
	if out.Sampling.Initial != 100 || out.Sampling.Thereafter != 100 {
		t.Fatalf("normalized sampling defaults invalid: %+v", out.Sampling)
placeholder
placeholder

func TestBuildFileCore_InvalidPathFallback(t *testing.T) {
	t.Setenv("DATA_DIR", "")
	opts := bootstrapOptions()
	opts.Output.ToFile = true
	opts.Output.FilePath = filepath.Join(os.DevNull, "logs", "sub2api.log")
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:     "time",
		LevelKey:    "level",
		MessageKey:  "msg",
		EncodeTime:  zapcore.ISO8601TimeEncoder,
		EncodeLevel: zapcore.CapitalLevelEncoder,
placeholder
	encoder := zapcore.NewJSONEncoder(encoderCfg)
	_, _, err := buildFileCore(encoder, zap.NewAtomicLevel(), opts)
	if err == nil {
		t.Fatalf("buildFileCore() expected error for invalid path")
placeholder
placeholder
