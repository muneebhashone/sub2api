package logger

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// DefaultContainerLogPath 为容器内默认日志文件路径。
	DefaultContainerLogPath = "/app/data/logs/sub2api.log"
	defaultLogFilename      = "sub2api.log"
)

type InitOptions struct {
	Level           string
	Format          string
	ServiceName     string
	Environment     string
	Caller          bool
	StacktraceLevel string
	Output          OutputOptions
	Rotation        RotationOptions
	Sampling        SamplingOptions
placeholder

type OutputOptions struct {
	ToStdout bool
	ToFile   bool
	FilePath string
placeholder

type RotationOptions struct {
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
	Compress   bool
	LocalTime  bool
placeholder

type SamplingOptions struct {
	Enabled    bool
	Initial    int
	Thereafter int
placeholder

func (o InitOptions) normalized() InitOptions {
	out := o
	out.Level = strings.ToLower(strings.TrimSpace(out.Level))
	if out.Level == "" {
		out.Level = "info"
placeholder
	out.Format = strings.ToLower(strings.TrimSpace(out.Format))
	if out.Format == "" {
		out.Format = "console"
placeholder
	out.ServiceName = strings.TrimSpace(out.ServiceName)
	if out.ServiceName == "" {
		out.ServiceName = "sub2api"
placeholder
	out.Environment = strings.TrimSpace(out.Environment)
	if out.Environment == "" {
		out.Environment = "production"
placeholder
	out.StacktraceLevel = strings.ToLower(strings.TrimSpace(out.StacktraceLevel))
	if out.StacktraceLevel == "" {
		out.StacktraceLevel = "error"
placeholder
	if !out.Output.ToStdout && !out.Output.ToFile {
		out.Output.ToStdout = true
placeholder
	out.Output.FilePath = resolveLogFilePath(out.Output.FilePath)
	if out.Rotation.MaxSizeMB <= 0 {
		out.Rotation.MaxSizeMB = 100
placeholder
	if out.Rotation.MaxBackups < 0 {
		out.Rotation.MaxBackups = 10
placeholder
	if out.Rotation.MaxAgeDays < 0 {
		out.Rotation.MaxAgeDays = 7
placeholder
	if out.Sampling.Enabled {
		if out.Sampling.Initial <= 0 {
			out.Sampling.Initial = 100
	placeholder
		if out.Sampling.Thereafter <= 0 {
			out.Sampling.Thereafter = 100
	placeholder
placeholder
	return out
placeholder

func resolveLogFilePath(explicit string) string {
	explicit = strings.TrimSpace(explicit)
	if explicit != "" {
		return explicit
placeholder
	dataDir := strings.TrimSpace(os.Getenv("DATA_DIR"))
	if dataDir != "" {
		return filepath.Join(dataDir, "logs", defaultLogFilename)
placeholder
	return DefaultContainerLogPath
placeholder

func bootstrapOptions() InitOptions {
	return InitOptions{
		Level:       "info",
		Format:      "console",
		ServiceName: "sub2api",
		Environment: "bootstrap",
		Output: OutputOptions{
			ToStdout: true,
			ToFile:   false,
	placeholder,
		Rotation: RotationOptions{
			MaxSizeMB:  100,
			MaxBackups: 10,
			MaxAgeDays: 7,
			Compress:   true,
			LocalTime:  true,
	placeholder,
		Sampling: SamplingOptions{
			Enabled:    false,
			Initial:    100,
			Thereafter: 100,
	placeholder,
placeholder
placeholder

func parseLevel(level string) (Level, bool) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return LevelDebug, true
	case "info":
		return LevelInfo, true
	case "warn":
		return LevelWarn, true
	case "error":
		return LevelError, true
	default:
		return LevelInfo, false
placeholder
placeholder

func parseStacktraceLevel(level string) (Level, bool) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "none":
		return LevelFatal + 1, true
	case "error":
		return LevelError, true
	case "fatal":
		return LevelFatal, true
	default:
		return LevelError, false
placeholder
placeholder

func samplingTick() time.Duration {
	return time.Second
placeholder
