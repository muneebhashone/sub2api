package service

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/antigravity"
)

func TestShortHash(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
placeholder{
		{"empty", []byte{placeholderplaceholder,
		{"simple", []byte("hello world")placeholder,
		{"json", []byte(`{"role":"user","parts":[{"text":"hello"placeholder]placeholder`)placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shortHash(tt.input)
			// Base36 编码的 uint64 最长 13 个字符
			if len(result) > 13 {
				t.Errorf("shortHash result too long: %d characters", len(result))
		placeholder
			// 相同输入应该产生相同输出
			result2 := shortHash(tt.input)
			if result != result2 {
				t.Errorf("shortHash not deterministic: %s vs %s", result, result2)
		placeholder
	placeholder)
placeholder
placeholder

func TestBuildGeminiDigestChain(t *testing.T) {
	tests := []struct {
		name     string
		req      *antigravity.GeminiRequest
		wantLen  int  // 预期的分段数量
		hasEmpty bool // 是否应该是空字符串
placeholder{
		{
			name:     "nil request",
			req:      nil,
			hasEmpty: true,
	placeholder,
		{
			name: "empty contents",
			req: &antigravity.GeminiRequest{
				Contents: []antigravity.GeminiContent{placeholder,
		placeholder,
			hasEmpty: true,
	placeholder,
		{
			name: "single user message",
			req: &antigravity.GeminiRequest{
				Contents: []antigravity.GeminiContent{
					{Role: "user", Parts: []antigravity.GeminiPart{{Text: "hello"placeholderplaceholderplaceholder,
			placeholder,
		placeholder,
			wantLen: 1, // u:<hash>
	placeholder,
		{
			name: "user and model messages",
			req: &antigravity.GeminiRequest{
				Contents: []antigravity.GeminiContent{
					{Role: "user", Parts: []antigravity.GeminiPart{{Text: "hello"placeholderplaceholderplaceholder,
					{Role: "model", Parts: []antigravity.GeminiPart{{Text: "hi there"placeholderplaceholderplaceholder,
			placeholder,
		placeholder,
			wantLen: 2, // u:<hash>-m:<hash>
	placeholder,
		{
			name: "with system instruction",
			req: &antigravity.GeminiRequest{
				SystemInstruction: &antigravity.GeminiContent{
					Role:  "user",
					Parts: []antigravity.GeminiPart{{Text: "You are a helpful assistant"placeholderplaceholder,
			placeholder,
				Contents: []antigravity.GeminiContent{
					{Role: "user", Parts: []antigravity.GeminiPart{{Text: "hello"placeholderplaceholderplaceholder,
			placeholder,
		placeholder,
			wantLen: 2, // s:<hash>-u:<hash>
	placeholder,
		{
			name: "conversation with system",
			req: &antigravity.GeminiRequest{
				SystemInstruction: &antigravity.GeminiContent{
					Role:  "user",
					Parts: []antigravity.GeminiPart{{Text: "System prompt"placeholderplaceholder,
			placeholder,
				Contents: []antigravity.GeminiContent{
					{Role: "user", Parts: []antigravity.GeminiPart{{Text: "hello"placeholderplaceholderplaceholder,
					{Role: "model", Parts: []antigravity.GeminiPart{{Text: "hi"placeholderplaceholderplaceholder,
					{Role: "user", Parts: []antigravity.GeminiPart{{Text: "how are you?"placeholderplaceholderplaceholder,
			placeholder,
		placeholder,
			wantLen: 4, // s:<hash>-u:<hash>-m:<hash>-u:<hash>
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildGeminiDigestChain(tt.req)

			if tt.hasEmpty {
				if result != "" {
					t.Errorf("expected empty string, got: %s", result)
			placeholder
				return
		placeholder

			// 检查分段数量
			parts := splitChain(result)
			if len(parts) != tt.wantLen {
				t.Errorf("expected %d parts, got %d: %s", tt.wantLen, len(parts), result)
		placeholder

			// 验证每个分段的格式
			for _, part := range parts {
				if len(part) < 3 || part[1] != ':' {
					t.Errorf("invalid part format: %s", part)
			placeholder
				prefix := part[0]
				if prefix != 's' && prefix != 'u' && prefix != 'm' {
					t.Errorf("invalid prefix: %c", prefix)
			placeholder
		placeholder
	placeholder)
placeholder
placeholder

func TestGenerateGeminiPrefixHash(t *testing.T) {
	hash1 := GenerateGeminiPrefixHash(1, 100, "192.168.1.1", "Mozilla/5.0", "antigravity", "gemini-2.5-pro")
	hash2 := GenerateGeminiPrefixHash(1, 100, "192.168.1.1", "Mozilla/5.0", "antigravity", "gemini-2.5-pro")
	hash3 := GenerateGeminiPrefixHash(2, 100, "192.168.1.1", "Mozilla/5.0", "antigravity", "gemini-2.5-pro")

	// 相同输入应该产生相同输出
	if hash1 != hash2 {
		t.Errorf("GenerateGeminiPrefixHash not deterministic: %s vs %s", hash1, hash2)
placeholder

	// 不同输入应该产生不同输出
	if hash1 == hash3 {
		t.Errorf("GenerateGeminiPrefixHash collision for different inputs")
placeholder

	// Base64 URL 编码的 12 字节正好是 16 字符
	if len(hash1) != 16 {
		t.Errorf("expected 16 characters, got %d: %s", len(hash1), hash1)
placeholder
placeholder

func TestGenerateDigestChainPrefixes(t *testing.T) {
	tests := []struct {
		name    string
		chain   string
		want    []string
		wantLen int
placeholder{
		{
			name:    "empty",
			chain:   "",
			wantLen: 0,
	placeholder,
		{
			name:    "single part",
			chain:   "u:abc123",
			want:    []string{"u:abc123"placeholder,
			wantLen: 1,
	placeholder,
		{
			name:    "two parts",
			chain:   "s:xyz-u:abc",
			want:    []string{"s:xyz-u:abc", "s:xyz"placeholder,
			wantLen: 2,
	placeholder,
		{
			name:    "four parts",
			chain:   "s:a-u:b-m:c-u:d",
			want:    []string{"s:a-u:b-m:c-u:d", "s:a-u:b-m:c", "s:a-u:b", "s:a"placeholder,
			wantLen: 4,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateDigestChainPrefixes(tt.chain)

			if len(result) != tt.wantLen {
				t.Errorf("expected %d prefixes, got %d: %v", tt.wantLen, len(result), result)
		placeholder

			if tt.want != nil {
				for i, want := range tt.want {
					if i >= len(result) {
						t.Errorf("missing prefix at index %d", i)
						continue
				placeholder
					if result[i] != want {
						t.Errorf("prefix[%d]: expected %s, got %s", i, want, result[i])
				placeholder
			placeholder
		placeholder
	placeholder)
placeholder
placeholder

func TestParseGeminiSessionValue(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantUUID  string
		wantAccID int64
		wantOK    bool
placeholder{
		{
			name:   "empty",
			value:  "",
			wantOK: false,
	placeholder,
		{
			name:   "no colon",
			value:  "abc123",
			wantOK: false,
	placeholder,
		{
			name:      "valid",
			value:     "uuid-1234:100",
			wantUUID:  "uuid-1234",
			wantAccID: 100,
			wantOK:    true,
	placeholder,
		{
			name:      "uuid with colon",
			value:     "a:b:c:123",
			wantUUID:  "a:b:c",
			wantAccID: 123,
			wantOK:    true,
	placeholder,
		{
			name:   "invalid account id",
			value:  "uuid:abc",
			wantOK: false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuid, accID, ok := ParseGeminiSessionValue(tt.value)

			if ok != tt.wantOK {
				t.Errorf("ok: expected %v, got %v", tt.wantOK, ok)
		placeholder

			if tt.wantOK {
				if uuid != tt.wantUUID {
					t.Errorf("uuid: expected %s, got %s", tt.wantUUID, uuid)
			placeholder
				if accID != tt.wantAccID {
					t.Errorf("accountID: expected %d, got %d", tt.wantAccID, accID)
			placeholder
		placeholder
	placeholder)
placeholder
placeholder

func TestFormatGeminiSessionValue(t *testing.T) {
	result := FormatGeminiSessionValue("test-uuid", 123)
	expected := "test-uuid:123"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
placeholder

	// 验证往返一致性
	uuid, accID, ok := ParseGeminiSessionValue(result)
	if !ok {
		t.Error("ParseGeminiSessionValue failed on formatted value")
placeholder
	if uuid != "test-uuid" || accID != 123 {
		t.Errorf("round-trip failed: uuid=%s, accID=%d", uuid, accID)
placeholder
placeholder

// splitChain 辅助函数：按 "-" 分割摘要链
func splitChain(chain string) []string {
	if chain == "" {
		return nil
placeholder
	var parts []string
	start := 0
	for i := 0; i < len(chain); i++ {
		if chain[i] == '-' {
			parts = append(parts, chain[start:i])
			start = i + 1
	placeholder
placeholder
	if start < len(chain) {
		parts = append(parts, chain[start:])
placeholder
	return parts
placeholder

func TestDigestChainDifferentSysInstruction(t *testing.T) {
	req1 := &antigravity.GeminiRequest{
		SystemInstruction: &antigravity.GeminiContent{
			Parts: []antigravity.GeminiPart{{Text: "SYS_ORIGINAL"placeholderplaceholder,
	placeholder,
		Contents: []antigravity.GeminiContent{
			{Role: "user", Parts: []antigravity.GeminiPart{{Text: "hello"placeholderplaceholderplaceholder,
	placeholder,
placeholder

	req2 := &antigravity.GeminiRequest{
		SystemInstruction: &antigravity.GeminiContent{
			Parts: []antigravity.GeminiPart{{Text: "SYS_MODIFIED"placeholderplaceholder,
	placeholder,
		Contents: []antigravity.GeminiContent{
			{Role: "user", Parts: []antigravity.GeminiPart{{Text: "hello"placeholderplaceholderplaceholder,
	placeholder,
placeholder

	chain1 := BuildGeminiDigestChain(req1)
	chain2 := BuildGeminiDigestChain(req2)

	t.Logf("Chain1: %s", chain1)
	t.Logf("Chain2: %s", chain2)

	if chain1 == chain2 {
		t.Error("Different systemInstruction should produce different chains")
placeholder
placeholder

func TestDigestChainTamperedMiddleContent(t *testing.T) {
	req1 := &antigravity.GeminiRequest{
		Contents: []antigravity.GeminiContent{
			{Role: "user", Parts: []antigravity.GeminiPart{{Text: "hello"placeholderplaceholderplaceholder,
			{Role: "model", Parts: []antigravity.GeminiPart{{Text: "ORIGINAL_REPLY"placeholderplaceholderplaceholder,
			{Role: "user", Parts: []antigravity.GeminiPart{{Text: "next"placeholderplaceholderplaceholder,
	placeholder,
placeholder

	req2 := &antigravity.GeminiRequest{
		Contents: []antigravity.GeminiContent{
			{Role: "user", Parts: []antigravity.GeminiPart{{Text: "hello"placeholderplaceholderplaceholder,
			{Role: "model", Parts: []antigravity.GeminiPart{{Text: "TAMPERED_REPLY"placeholderplaceholderplaceholder,
			{Role: "user", Parts: []antigravity.GeminiPart{{Text: "next"placeholderplaceholderplaceholder,
	placeholder,
placeholder

	chain1 := BuildGeminiDigestChain(req1)
	chain2 := BuildGeminiDigestChain(req2)

	t.Logf("Chain1: %s", chain1)
	t.Logf("Chain2: %s", chain2)

	if chain1 == chain2 {
		t.Error("Tampered middle content should produce different chains")
placeholder

	// 验证第一个 user 的 hash 相同
	parts1 := splitChain(chain1)
	parts2 := splitChain(chain2)

	if parts1[0] != parts2[0] {
		t.Error("First user message hash should be the same")
placeholder
	if parts1[1] == parts2[1] {
		t.Error("Model reply hash should be different")
placeholder
placeholder

func TestGenerateGeminiDigestSessionKey(t *testing.T) {
	tests := []struct {
		name       string
		prefixHash string
		uuid       string
		want       string
placeholder{
		{
			name:       "normal 16 char hash with uuid",
			prefixHash: "abcdefgh12345678",
			uuid:       "550e8400-e29b-41d4-a716-446655440000",
			want:       "gemini:digest:abcdefgh:550e8400",
	placeholder,
		{
			name:       "exactly 8 chars prefix and uuid",
			prefixHash: "12345678",
			uuid:       "abcdefgh",
			want:       "gemini:digest:12345678:abcdefgh",
	placeholder,
		{
			name:       "short hash and short uuid (less than 8)",
			prefixHash: "abc",
			uuid:       "xyz",
			want:       "gemini:digest:abc:xyz",
	placeholder,
		{
			name:       "empty hash and uuid",
			prefixHash: "",
			uuid:       "",
			want:       "gemini:digest::",
	placeholder,
		{
			name:       "normal prefix with short uuid",
			prefixHash: "abcdefgh12345678",
			uuid:       "short",
			want:       "gemini:digest:abcdefgh:short",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateGeminiDigestSessionKey(tt.prefixHash, tt.uuid)
			if got != tt.want {
				t.Errorf("GenerateGeminiDigestSessionKey(%q, %q) = %q, want %q", tt.prefixHash, tt.uuid, got, tt.want)
		placeholder
	placeholder)
placeholder

	// 验证确定性：相同输入产生相同输出
	t.Run("deterministic", func(t *testing.T) {
		hash := "testprefix123456"
		uuid := "test-uuid-12345"
		result1 := GenerateGeminiDigestSessionKey(hash, uuid)
		result2 := GenerateGeminiDigestSessionKey(hash, uuid)
		if result1 != result2 {
			t.Errorf("GenerateGeminiDigestSessionKey not deterministic: %s vs %s", result1, result2)
	placeholder
placeholder)

	// 验证不同 uuid 产生不同 sessionKey（负载均衡核心逻辑）
	t.Run("different uuid different key", func(t *testing.T) {
		hash := "sameprefix123456"
		uuid1 := "uuid0001-session-a"
		uuid2 := "uuid0002-session-b"
		result1 := GenerateGeminiDigestSessionKey(hash, uuid1)
		result2 := GenerateGeminiDigestSessionKey(hash, uuid2)
		if result1 == result2 {
			t.Errorf("Different UUIDs should produce different session keys: %s vs %s", result1, result2)
	placeholder
placeholder)
placeholder

func TestBuildGeminiTrieKey(t *testing.T) {
	tests := []struct {
		name       string
		groupID    int64
		prefixHash string
		want       string
placeholder{
		{
			name:       "normal",
			groupID:    123,
			prefixHash: "abcdef12",
			want:       "gemini:trie:123:abcdef12",
	placeholder,
		{
			name:       "zero group",
			groupID:    0,
			prefixHash: "xyz",
			want:       "gemini:trie:0:xyz",
	placeholder,
		{
			name:       "empty prefix",
			groupID:    1,
			prefixHash: "",
			want:       "gemini:trie:1:",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildGeminiTrieKey(tt.groupID, tt.prefixHash)
			if got != tt.want {
				t.Errorf("BuildGeminiTrieKey(%d, %q) = %q, want %q", tt.groupID, tt.prefixHash, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder
