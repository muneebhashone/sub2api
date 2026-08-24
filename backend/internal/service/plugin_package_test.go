package service

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	pluginv1 "github.com/Wei-Shaw/sub2api/pkg/pluginapi/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginPackageInstallerInstallUnsignedDevelopmentPackage(t *testing.T) {
	cfg := testPluginConfig(t.TempDir(), true)
	installer := NewPluginPackageInstaller(cfg, PluginHostInfo{Version: "0.1.179", BuildType: "release"placeholder)
	archive := buildTestPluginArchive(t, nil, "")

	installation, err := installer.Install(context.Background(), bytes.NewReader(archive), nil)

placeholder
	assert.Equal(t, PluginStateDisabled, installation.State)
	assert.Equal(t, PluginSignatureUnsigned, installation.SignatureStatus)
	assert.FileExists(t, installation.ArtifactPath)
	info, statErr := os.Stat(installation.BinaryPath)
	require.NoError(t, statErr)
	assert.NotZero(t, info.Mode()&0o100)
	assert.Contains(t, installation.InstallPath, filepath.Join("installed", "com.example.openai-transport"))
placeholder

func TestPluginPackageInstallerAllowsRepeatedIdenticalUpload(t *testing.T) {
	cfg := testPluginConfig(t.TempDir(), true)
	installer := NewPluginPackageInstaller(cfg, PluginHostInfo{Version: "0.1.179"placeholder)
	archive := buildTestPluginArchive(t, nil, "")

	first, err := installer.Install(context.Background(), bytes.NewReader(archive), nil)
placeholder
	second, err := installer.Install(context.Background(), bytes.NewReader(archive), nil)
placeholder

	assert.NotEqual(t, first.InstallPath, second.InstallPath)
	assert.NotEqual(t, first.ArtifactPath, second.ArtifactPath)
	assert.FileExists(t, first.BinaryPath)
	assert.FileExists(t, second.BinaryPath)
placeholder

func TestPluginPackageInstallerRejectsUnsignedPackageByDefault(t *testing.T) {
	cfg := testPluginConfig(t.TempDir(), false)
	installer := NewPluginPackageInstaller(cfg, PluginHostInfo{Version: "0.1.179"placeholder)

	_, err := installer.Install(context.Background(), bytes.NewReader(buildTestPluginArchive(t, nil, "")), nil)

placeholder
	assert.Contains(t, err.Error(), "不允许安装未签名插件")
placeholder

func TestPluginPackageInstallerVerifiesTrustedSignature(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
placeholder
	cfg := testPluginConfig(t.TempDir(), false)
	cfg.Plugins.TrustedPublishers["local-test"] = base64.StdEncoding.EncodeToString(publicKey)
	installer := NewPluginPackageInstaller(cfg, PluginHostInfo{Version: "0.1.179"placeholder)

	installation, installErr := installer.Install(
		context.Background(),
		bytes.NewReader(buildTestPluginArchive(t, privateKey, "local-test")),
		nil,
	)

	require.NoError(t, installErr)
	assert.Equal(t, PluginSignatureTrusted, installation.SignatureStatus)
placeholder

func TestBuiltInOpenAITransportPublisherDoesNotRequireConfiguration(t *testing.T) {
	cfg := testPluginConfig(t.TempDir(), false)
	cfg.Plugins.TrustedPublishers[builtInOpenAITransportPublisherKeyID] = "不能覆盖内置公钥"

	encodedKey := trustedPluginPublisherKey(cfg, builtInOpenAITransportPublisherKeyID, builtInOpenAITransportPluginID)
	publicKey, err := base64.StdEncoding.DecodeString(encodedKey)

placeholder
	assert.Len(t, publicKey, ed25519.PublicKeySize)
	assert.Equal(t, builtInOpenAITransportPublisherKeyBase64, encodedKey)
	assert.Empty(t, trustedPluginPublisherKey(cfg, builtInOpenAITransportPublisherKeyID, "com.example.other-plugin"))
placeholder

func TestPluginPackageInstallerRejectsPathTraversal(t *testing.T) {
	cfg := testPluginConfig(t.TempDir(), true)
	installer := NewPluginPackageInstaller(cfg, PluginHostInfo{Version: "0.1.179"placeholder)
	archive := buildTestPluginArchiveWithExtra(t, nil, "", "../escape", []byte("escape"))

	_, err := installer.Install(context.Background(), bytes.NewReader(archive), nil)

placeholder
	assert.Contains(t, err.Error(), "不安全路径")
placeholder

func TestPluginPackageInstallerKeepsHostVersionMismatchDisabled(t *testing.T) {
	cfg := testPluginConfig(t.TempDir(), true)
	installer := NewPluginPackageInstaller(cfg, PluginHostInfo{Version: "0.2.0"placeholder)

	installation, err := installer.Install(context.Background(), bytes.NewReader(buildTestPluginArchive(t, nil, "")), nil)

placeholder
	assert.Equal(t, PluginStateIncompatible, installation.State)
	assert.False(t, installation.Compatibility.Compatible)
placeholder

func TestPluginPackageInstallerRejectsHashMismatch(t *testing.T) {
	cfg := testPluginConfig(t.TempDir(), true)
	installer := NewPluginPackageInstaller(cfg, PluginHostInfo{Version: "0.1.179"placeholder)
	manifest := testPluginManifest(map[string][]byte{
		"bin/plugin":    []byte("binary"),
		"ui/index.html": []byte("<html></html>"),
placeholder)
	manifest.Files["bin/plugin"] = string(bytes.Repeat([]byte("0"), 64))

	_, err := installer.Install(context.Background(), bytes.NewReader(buildPluginArchive(t, manifest, nil, "", nil)), nil)

placeholder
	assert.Contains(t, err.Error(), "哈希不匹配")
placeholder

func TestPluginPackageInstallerEnforcesActualExtractionLimit(t *testing.T) {
	archiveData := buildTestPluginArchive(t, nil, "")
	archive, err := zip.NewReader(bytes.NewReader(archiveData), int64(len(archiveData)))
placeholder
	cfg := testPluginConfig(t.TempDir(), true)
	cfg.Plugins.MaxUncompressedBytes = 8
	installer := NewPluginPackageInstaller(cfg, PluginHostInfo{Version: "0.1.179"placeholder)

	err = installer.extractArchive(context.Background(), archive, testPluginManifest(nil), t.TempDir())

	require.ErrorContains(t, err, "实际解压体积超过限制")
placeholder

func testPluginConfig(root string, allowUnsigned bool) *config.Config {
	return &config.Config{Plugins: config.PluginConfig{
		DataDir:              root,
		AllowUnsigned:        allowUnsigned,
		TrustedPublishers:    map[string]string{placeholder,
		MaxUploadBytes:       64 * 1024 * 1024,
		MaxUncompressedBytes: 128 * 1024 * 1024,
		StartTimeoutSeconds:  5,
placeholderplaceholder
placeholder

func testPluginManifest(files map[string][]byte) PluginManifest {
	if files == nil {
		files = map[string][]byte{
			"bin/plugin":    []byte("binary"),
			"ui/index.html": []byte("<html></html>"),
	placeholder
placeholder
	hashes := make(map[string]string, len(files))
	for path, data := range files {
		digest := sha256.Sum256(data)
		hashes[path] = hex.EncodeToString(digest[:])
placeholder
	return PluginManifest{
		SchemaVersion: 1,
		ID:            "com.example.openai-transport",
		Name:          "测试 OpenAI Transport",
		Version:       "0.1.0",
		Requires: PluginRequirements{
			Sub2API:                   ">=0.1.170 <0.2.0",
			RecommendedSub2APIVersion: "0.1.179",
			TestedSub2APIVersions:     []string{"0.1.179"placeholder,
			PluginProtocol:            pluginv1.ProtocolVersion,
			TransportAPI:              pluginv1.TransportAPIVersion,
			UIBridge:                  pluginv1.UIBridgeVersion,
	placeholder,
		Capabilities: []PluginCapability{{
			ID:          PluginCapabilityOpenAIOAuthOutbound,
			Platform:    PlatformOpenAI,
			AccountType: AccountTypeOAuth,
placeholder
		Runtimes: map[string]PluginRuntime{
			PluginManifest{placeholder.RuntimeKey(): {Path: "bin/plugin"placeholder,
	placeholder,
		UI:    PluginUIManifest{Entrypoint: "ui/index.html"placeholder,
		Files: hashes,
placeholder
placeholder

func buildTestPluginArchive(t *testing.T, privateKey ed25519.PrivateKey, keyID string) []byte {
	return buildPluginArchive(t, testPluginManifest(nil), privateKey, keyID, nil)
placeholder

func buildTestPluginArchiveWithExtra(t *testing.T, privateKey ed25519.PrivateKey, keyID, path string, data []byte) []byte {
	return buildPluginArchive(t, testPluginManifest(nil), privateKey, keyID, map[string][]byte{path: dataplaceholder)
placeholder

func buildPluginArchive(
	t *testing.T,
	manifest PluginManifest,
	privateKey ed25519.PrivateKey,
	keyID string,
	extra map[string][]byte,
) []byte {
placeholder
	manifestRaw, err := json.Marshal(manifest)
placeholder
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	writeZipEntry(t, writer, pluginManifestFilename, manifestRaw)
	for path := range manifest.Files {
		data := []byte("binary")
		if path == "ui/index.html" {
			data = []byte("<html></html>")
	placeholder
		writeZipEntry(t, writer, path, data)
placeholder
	for path, data := range extra {
		writeZipEntry(t, writer, path, data)
placeholder
	if len(privateKey) > 0 {
		signatureRaw, marshalErr := json.Marshal(PluginSignature{
			Algorithm: "ed25519",
			KeyID:     keyID,
			Signature: base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, manifestRaw)),
	placeholder)
		require.NoError(t, marshalErr)
		writeZipEntry(t, writer, pluginSignatureFilename, signatureRaw)
placeholder
	require.NoError(t, writer.Close())
	return buffer.Bytes()
placeholder

func writeZipEntry(t *testing.T, writer *zip.Writer, path string, data []byte) {
placeholder
	entry, err := writer.Create(path)
placeholder
	_, err = entry.Write(data)
placeholder
placeholder
