//go:build unit

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestIsImageGenerationModel_GeminiProImage 测试 gemini-3-pro-image 识别
func TestIsImageGenerationModel_GeminiProImage(t *testing.T) {
	require.True(t, isImageGenerationModel("gemini-3-pro-image"))
	require.True(t, isImageGenerationModel("gemini-3-pro-image-preview"))
	require.True(t, isImageGenerationModel("models/gemini-3-pro-image"))
placeholder

// TestIsImageGenerationModel_GeminiFlashImage 测试 gemini-2.5-flash-image 识别
func TestIsImageGenerationModel_GeminiFlashImage(t *testing.T) {
	require.True(t, isImageGenerationModel("gemini-2.5-flash-image"))
	require.True(t, isImageGenerationModel("gemini-2.5-flash-image-preview"))
placeholder

// TestIsImageGenerationModel_RegularModel 测试普通模型不被识别为图片模型
func TestIsImageGenerationModel_RegularModel(t *testing.T) {
	require.False(t, isImageGenerationModel("claude-3-opus"))
	require.False(t, isImageGenerationModel("claude-sonnet-4-20250514"))
	require.False(t, isImageGenerationModel("gpt-4o"))
	require.False(t, isImageGenerationModel("gemini-2.5-pro")) // 非图片模型
	require.False(t, isImageGenerationModel("gemini-2.5-flash"))
placeholder

// TestIsImageGenerationModel_CaseInsensitive 测试大小写不敏感
func TestIsImageGenerationModel_CaseInsensitive(t *testing.T) {
	require.True(t, isImageGenerationModel("GEMINI-3-PRO-IMAGE"))
	require.True(t, isImageGenerationModel("Gemini-3-Pro-Image"))
	require.True(t, isImageGenerationModel("GEMINI-2.5-FLASH-IMAGE"))
placeholder

// TestExtractImageSize_ValidSizes 测试有效尺寸解析
func TestExtractImageSize_ValidSizes(t *testing.T) {
	svc := &AntigravityGatewayService{placeholder

	// 1K
	body := []byte(`{"generationConfig":{"imageConfig":{"imageSize":"1K"placeholderplaceholderplaceholder`)
	require.Equal(t, "1K", svc.extractImageSize(body))

	// 2K
	body = []byte(`{"generationConfig":{"imageConfig":{"imageSize":"2K"placeholderplaceholderplaceholder`)
	require.Equal(t, "2K", svc.extractImageSize(body))

	// 4K
	body = []byte(`{"generationConfig":{"imageConfig":{"imageSize":"4K"placeholderplaceholderplaceholder`)
	require.Equal(t, "4K", svc.extractImageSize(body))
placeholder

// TestExtractImageSize_CaseInsensitive 测试大小写不敏感
func TestExtractImageSize_CaseInsensitive(t *testing.T) {
	svc := &AntigravityGatewayService{placeholder

	body := []byte(`{"generationConfig":{"imageConfig":{"imageSize":"1k"placeholderplaceholderplaceholder`)
	require.Equal(t, "1K", svc.extractImageSize(body))

	body = []byte(`{"generationConfig":{"imageConfig":{"imageSize":"4k"placeholderplaceholderplaceholder`)
	require.Equal(t, "4K", svc.extractImageSize(body))
placeholder

// TestExtractImageSize_Default 测试无 imageConfig 返回默认 2K
func TestExtractImageSize_Default(t *testing.T) {
	svc := &AntigravityGatewayService{placeholder

	// 无 generationConfig
	body := []byte(`{"contents":[]placeholder`)
	require.Equal(t, "2K", svc.extractImageSize(body))

	// 有 generationConfig 但无 imageConfig
	body = []byte(`{"generationConfig":{"temperature":0.7placeholderplaceholder`)
	require.Equal(t, "2K", svc.extractImageSize(body))

	// 有 imageConfig 但无 imageSize
	body = []byte(`{"generationConfig":{"imageConfig":{placeholderplaceholderplaceholder`)
	require.Equal(t, "2K", svc.extractImageSize(body))
placeholder

// TestExtractImageSize_InvalidJSON 测试非法 JSON 返回默认 2K
func TestExtractImageSize_InvalidJSON(t *testing.T) {
	svc := &AntigravityGatewayService{placeholder

	body := []byte(`not valid json`)
	require.Equal(t, "2K", svc.extractImageSize(body))

	body = []byte(`{"broken":`)
	require.Equal(t, "2K", svc.extractImageSize(body))
placeholder

// TestExtractImageSize_EmptySize 测试空 imageSize 返回默认 2K
func TestExtractImageSize_EmptySize(t *testing.T) {
	svc := &AntigravityGatewayService{placeholder

	body := []byte(`{"generationConfig":{"imageConfig":{"imageSize":""placeholderplaceholderplaceholder`)
	require.Equal(t, "2K", svc.extractImageSize(body))

	// 空格
	body = []byte(`{"generationConfig":{"imageConfig":{"imageSize":"   "placeholderplaceholderplaceholder`)
	require.Equal(t, "2K", svc.extractImageSize(body))
placeholder

// TestExtractImageSize_InvalidSize 测试无效尺寸返回默认 2K
func TestExtractImageSize_InvalidSize(t *testing.T) {
	svc := &AntigravityGatewayService{placeholder

	body := []byte(`{"generationConfig":{"imageConfig":{"imageSize":"3K"placeholderplaceholderplaceholder`)
	require.Equal(t, "2K", svc.extractImageSize(body))

	body = []byte(`{"generationConfig":{"imageConfig":{"imageSize":"8K"placeholderplaceholderplaceholder`)
	require.Equal(t, "2K", svc.extractImageSize(body))

	body = []byte(`{"generationConfig":{"imageConfig":{"imageSize":"invalid"placeholderplaceholderplaceholder`)
	require.Equal(t, "2K", svc.extractImageSize(body))
placeholder
