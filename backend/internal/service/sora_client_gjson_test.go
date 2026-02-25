//go:build unit

package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// ---------- 辅助解析函数（复制生产代码中的 gjson 解析逻辑，用于单元测试） ----------

// testParseUploadOrCreateTaskID 模拟 UploadImage / CreateImageTask / CreateVideoTask 中
// 用 gjson.GetBytes(respBody, "id") 提取 id 的逻辑。
func testParseUploadOrCreateTaskID(respBody []byte) (string, error) {
	id := strings.TrimSpace(gjson.GetBytes(respBody, "id").String())
	if id == "" {
		return "", assert.AnError // 占位错误，表示 "missing id"
placeholder
	return id, nil
placeholder

// testParseFetchRecentImageTask 模拟 fetchRecentImageTask 中的 gjson.ForEach 解析逻辑。
func testParseFetchRecentImageTask(respBody []byte, taskID string) (*SoraImageTaskStatus, bool) {
	var found *SoraImageTaskStatus
	gjson.GetBytes(respBody, "task_responses").ForEach(func(_, item gjson.Result) bool {
		if item.Get("id").String() != taskID {
			return true // continue
	placeholder
		status := strings.TrimSpace(item.Get("status").String())
		progress := item.Get("progress_pct").Float()
		var urls []string
		item.Get("generations").ForEach(func(_, gen gjson.Result) bool {
			if u := strings.TrimSpace(gen.Get("url").String()); u != "" {
				urls = append(urls, u)
		placeholder
			return true
	placeholder)
		found = &SoraImageTaskStatus{
			ID:          taskID,
			Status:      status,
			ProgressPct: progress,
			URLs:        urls,
	placeholder
		return false // break
placeholder)
	if found != nil {
		return found, true
placeholder
	return &SoraImageTaskStatus{ID: taskID, Status: "processing"placeholder, false
placeholder

// testParseGetVideoTaskPending 模拟 GetVideoTask 中解析 pending 列表的逻辑。
func testParseGetVideoTaskPending(respBody []byte, taskID string) (*SoraVideoTaskStatus, bool) {
	pendingResult := gjson.ParseBytes(respBody)
	if !pendingResult.IsArray() {
		return nil, false
placeholder
	var pendingFound *SoraVideoTaskStatus
	pendingResult.ForEach(func(_, task gjson.Result) bool {
		if task.Get("id").String() != taskID {
			return true
	placeholder
		progress := 0
		if v := task.Get("progress_pct"); v.Exists() {
			progress = int(v.Float() * 100)
	placeholder
		status := strings.TrimSpace(task.Get("status").String())
		pendingFound = &SoraVideoTaskStatus{
			ID:          taskID,
			Status:      status,
			ProgressPct: progress,
	placeholder
		return false
placeholder)
	if pendingFound != nil {
		return pendingFound, true
placeholder
	return nil, false
placeholder

// testParseGetVideoTaskDrafts 模拟 GetVideoTask 中解析 drafts 列表的逻辑。
func testParseGetVideoTaskDrafts(respBody []byte, taskID string) (*SoraVideoTaskStatus, bool) {
	var draftFound *SoraVideoTaskStatus
	gjson.GetBytes(respBody, "items").ForEach(func(_, draft gjson.Result) bool {
		if draft.Get("task_id").String() != taskID {
			return true
	placeholder
		kind := strings.TrimSpace(draft.Get("kind").String())
		reason := strings.TrimSpace(draft.Get("reason_str").String())
		if reason == "" {
			reason = strings.TrimSpace(draft.Get("markdown_reason_str").String())
	placeholder
		urlStr := strings.TrimSpace(draft.Get("downloadable_url").String())
		if urlStr == "" {
			urlStr = strings.TrimSpace(draft.Get("url").String())
	placeholder

		if kind == "sora_content_violation" || reason != "" || urlStr == "" {
			msg := reason
			if msg == "" {
				msg = "Content violates guardrails"
		placeholder
			draftFound = &SoraVideoTaskStatus{
				ID:       taskID,
				Status:   "failed",
				ErrorMsg: msg,
		placeholder
	placeholder else {
			draftFound = &SoraVideoTaskStatus{
				ID:     taskID,
				Status: "completed",
				URLs:   []string{urlStrplaceholder,
		placeholder
	placeholder
		return false
placeholder)
	if draftFound != nil {
		return draftFound, true
placeholder
	return nil, false
placeholder

// ===================== Test 1: TestSoraParseUploadResponse =====================

func TestSoraParseUploadResponse(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantID  string
		wantErr bool
placeholder{
		{
			name:   "正常 id",
			body:   `{"id":"file-abc123","status":"uploaded"placeholder`,
			wantID: "file-abc123",
	placeholder,
		{
			name:    "空 id",
			body:    `{"id":"","status":"uploaded"placeholder`,
			wantErr: true,
	placeholder,
		{
			name:    "无 id 字段",
			body:    `{"status":"uploaded"placeholder`,
			wantErr: true,
	placeholder,
		{
			name:    "id 全为空白",
			body:    `{"id":"   ","status":"uploaded"placeholder`,
			wantErr: true,
	placeholder,
		{
			name:   "id 前后有空白",
			body:   `{"id":"  file-trimmed  ","status":"uploaded"placeholder`,
			wantID: "file-trimmed",
	placeholder,
		{
			name:    "空 JSON 对象",
			body:    `{placeholder`,
			wantErr: true,
	placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := testParseUploadOrCreateTaskID([]byte(tt.body))
			if tt.wantErr {
				require.Error(t, err, "应返回错误")
				return
		placeholder
		placeholder
			require.Equal(t, tt.wantID, id)
	placeholder)
placeholder
placeholder

// ===================== Test 2: TestSoraParseCreateTaskResponse =====================

func TestSoraParseCreateTaskResponse(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantID  string
		wantErr bool
placeholder{
		{
			name:   "正常任务 id",
			body:   `{"id":"task-123"placeholder`,
			wantID: "task-123",
	placeholder,
		{
			name:    "缺失 id",
			body:    `{"status":"created"placeholder`,
			wantErr: true,
	placeholder,
		{
			name:    "空 id",
			body:    `{"id":"  "placeholder`,
			wantErr: true,
	placeholder,
		{
			name:   "id 为数字（gjson 转字符串）",
			body:   `{"id":placeholder`,
			wantID: "123",
	placeholder,
		{
			name:   "id 含特殊字符",
			body:   `{"id":"task-abc-def-456-ghi"placeholder`,
			wantID: "task-abc-def-456-ghi",
	placeholder,
		{
			name:   "额外字段不影响解析",
			body:   `{"id":"task-999","type":"image_gen","extra":"data"placeholder`,
			wantID: "task-999",
	placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := testParseUploadOrCreateTaskID([]byte(tt.body))
			if tt.wantErr {
				require.Error(t, err, "应返回错误")
				return
		placeholder
		placeholder
			require.Equal(t, tt.wantID, id)
	placeholder)
placeholder
placeholder

// ===================== Test 3: TestSoraParseFetchRecentImageTask =====================

func TestSoraParseFetchRecentImageTask(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		taskID       string
		wantFound    bool
		wantStatus   string
		wantProgress float64
		wantURLs     []string
placeholder{
		{
			name:         "匹配已完成任务",
			body:         `{"task_responses":[{"id":"task-1","status":"completed","progress_pct":1.0,"generations":[{"url":"https://example.com/img.png"placeholder]placeholder]placeholder`,
			taskID:       "task-1",
			wantFound:    true,
			wantStatus:   "completed",
			wantProgress: 1.0,
			wantURLs:     []string{"https://example.com/img.png"placeholder,
	placeholder,
		{
			name:         "匹配处理中任务",
			body:         `{"task_responses":[{"id":"task-2","status":"processing","progress_pct":0.5,"generations":[]placeholder]placeholder`,
			taskID:       "task-2",
			wantFound:    true,
			wantStatus:   "processing",
			wantProgress: 0.5,
			wantURLs:     nil,
	placeholder,
		{
			name:       "无匹配任务",
			body:       `{"task_responses":[{"id":"other","status":"completed"placeholder]placeholder`,
			taskID:     "task-1",
			wantFound:  false,
			wantStatus: "processing",
	placeholder,
		{
			name:       "空 task_responses",
			body:       `{"task_responses":[]placeholder`,
			taskID:     "task-1",
			wantFound:  false,
			wantStatus: "processing",
	placeholder,
		{
			name:       "缺少 task_responses 字段",
			body:       `{"other":"data"placeholder`,
			taskID:     "task-1",
			wantFound:  false,
			wantStatus: "processing",
	placeholder,
		{
			name:         "多个任务中精准匹配",
			body:         `{"task_responses":[{"id":"task-a","status":"completed","progress_pct":1.0,"generations":[{"url":"https://a.com/1.png"placeholder]placeholder,{"id":"task-b","status":"processing","progress_pct":0.3,"generations":[]placeholder,{"id":"task-c","status":"failed","progress_pct":0placeholder]placeholder`,
			taskID:       "task-b",
			wantFound:    true,
			wantStatus:   "processing",
			wantProgress: 0.3,
			wantURLs:     nil,
	placeholder,
		{
			name:         "多个 generations",
			body:         `{"task_responses":[{"id":"task-m","status":"completed","progress_pct":1.0,"generations":[{"url":"https://a.com/1.png"placeholder,{"url":"https://a.com/2.png"placeholder,{"url":""placeholder]placeholder]placeholder`,
			taskID:       "task-m",
			wantFound:    true,
			wantStatus:   "completed",
			wantProgress: 1.0,
			wantURLs:     []string{"https://a.com/1.png", "https://a.com/2.png"placeholder,
	placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, found := testParseFetchRecentImageTask([]byte(tt.body), tt.taskID)
			require.Equal(t, tt.wantFound, found, "found 不匹配")
			require.NotNil(t, status)
			require.Equal(t, tt.taskID, status.ID)
			require.Equal(t, tt.wantStatus, status.Status)
			if tt.wantFound {
				require.InDelta(t, tt.wantProgress, status.ProgressPct, 0.001, "进度不匹配")
				require.Equal(t, tt.wantURLs, status.URLs)
		placeholder
	placeholder)
placeholder
placeholder

// ===================== Test 4: TestSoraParseGetVideoTaskPending =====================

func TestSoraParseGetVideoTaskPending(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		taskID       string
		wantFound    bool
		wantStatus   string
		wantProgress int
placeholder{
		{
			name:         "匹配 pending 任务",
			body:         `[{"id":"task-1","status":"processing","progress_pct":0.5placeholder]`,
			taskID:       "task-1",
			wantFound:    true,
			wantStatus:   "processing",
			wantProgress: 50,
	placeholder,
		{
			name:         "进度为 0",
			body:         `[{"id":"task-2","status":"queued","progress_pct":0placeholder]`,
			taskID:       "task-2",
			wantFound:    true,
			wantStatus:   "queued",
			wantProgress: 0,
	placeholder,
		{
			name:         "进度为 1（100%）",
			body:         `[{"id":"task-3","status":"completing","progress_pct":1.0placeholder]`,
			taskID:       "task-3",
			wantFound:    true,
			wantStatus:   "completing",
			wantProgress: 100,
	placeholder,
		{
			name:      "空数组",
			body:      `[]`,
			taskID:    "task-1",
			wantFound: false,
	placeholder,
		{
			name:      "无匹配 id",
			body:      `[{"id":"task-other","status":"processing","progress_pct":0.3placeholder]`,
			taskID:    "task-1",
			wantFound: false,
	placeholder,
		{
			name:         "多个任务精准匹配",
			body:         `[{"id":"task-a","status":"processing","progress_pct":0.2placeholder,{"id":"task-b","status":"queued","progress_pct":0placeholder,{"id":"task-c","status":"processing","progress_pct":0.8placeholder]`,
			taskID:       "task-c",
			wantFound:    true,
			wantStatus:   "processing",
			wantProgress: 80,
	placeholder,
		{
			name:      "非数组 JSON",
			body:      `{"id":"task-1","status":"processing"placeholder`,
			taskID:    "task-1",
			wantFound: false,
	placeholder,
		{
			name:         "无 progress_pct 字段",
			body:         `[{"id":"task-4","status":"pending"placeholder]`,
			taskID:       "task-4",
			wantFound:    true,
			wantStatus:   "pending",
			wantProgress: 0,
	placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, found := testParseGetVideoTaskPending([]byte(tt.body), tt.taskID)
			require.Equal(t, tt.wantFound, found, "found 不匹配")
			if tt.wantFound {
				require.NotNil(t, status)
				require.Equal(t, tt.taskID, status.ID)
				require.Equal(t, tt.wantStatus, status.Status)
				require.Equal(t, tt.wantProgress, status.ProgressPct)
		placeholder
	placeholder)
placeholder
placeholder

// ===================== Test 5: TestSoraParseGetVideoTaskDrafts =====================

func TestSoraParseGetVideoTaskDrafts(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		taskID     string
		wantFound  bool
		wantStatus string
		wantURLs   []string
		wantErr    string
placeholder{
		{
			name:       "正常完成的视频",
			body:       `{"items":[{"task_id":"task-1","kind":"video","downloadable_url":"https://example.com/video.mp4"placeholder]placeholder`,
			taskID:     "task-1",
			wantFound:  true,
			wantStatus: "completed",
			wantURLs:   []string{"https://example.com/video.mp4"placeholder,
	placeholder,
		{
			name:       "使用 url 字段回退",
			body:       `{"items":[{"task_id":"task-2","kind":"video","url":"https://example.com/fallback.mp4"placeholder]placeholder`,
			taskID:     "task-2",
			wantFound:  true,
			wantStatus: "completed",
			wantURLs:   []string{"https://example.com/fallback.mp4"placeholder,
	placeholder,
		{
			name:       "内容违规",
			body:       `{"items":[{"task_id":"task-3","kind":"sora_content_violation","reason_str":"Content policy violation"placeholder]placeholder`,
			taskID:     "task-3",
			wantFound:  true,
			wantStatus: "failed",
			wantErr:    "Content policy violation",
	placeholder,
		{
			name:       "内容违规 - markdown_reason_str 回退",
			body:       `{"items":[{"task_id":"task-4","kind":"sora_content_violation","markdown_reason_str":"Markdown reason"placeholder]placeholder`,
			taskID:     "task-4",
			wantFound:  true,
			wantStatus: "failed",
			wantErr:    "Markdown reason",
	placeholder,
		{
			name:       "内容违规 - 无 reason 使用默认消息",
			body:       `{"items":[{"task_id":"task-5","kind":"sora_content_violation"placeholder]placeholder`,
			taskID:     "task-5",
			wantFound:  true,
			wantStatus: "failed",
			wantErr:    "Content violates guardrails",
	placeholder,
		{
			name:       "有 reason_str 但非 violation kind（仍判定失败）",
			body:       `{"items":[{"task_id":"task-6","kind":"video","reason_str":"Some error occurred"placeholder]placeholder`,
			taskID:     "task-6",
			wantFound:  true,
			wantStatus: "failed",
			wantErr:    "Some error occurred",
	placeholder,
		{
			name:       "空 URL 判定为失败",
			body:       `{"items":[{"task_id":"task-7","kind":"video","downloadable_url":"","url":""placeholder]placeholder`,
			taskID:     "task-7",
			wantFound:  true,
			wantStatus: "failed",
			wantErr:    "Content violates guardrails",
	placeholder,
		{
			name:      "无匹配 task_id",
			body:      `{"items":[{"task_id":"task-other","kind":"video","downloadable_url":"https://example.com/video.mp4"placeholder]placeholder`,
			taskID:    "task-1",
			wantFound: false,
	placeholder,
		{
			name:      "空 items",
			body:      `{"items":[]placeholder`,
			taskID:    "task-1",
			wantFound: false,
	placeholder,
		{
			name:      "缺少 items 字段",
			body:      `{"other":"data"placeholder`,
			taskID:    "task-1",
			wantFound: false,
	placeholder,
		{
			name:       "多个 items 精准匹配",
			body:       `{"items":[{"task_id":"task-a","kind":"video","downloadable_url":"https://a.com/a.mp4"placeholder,{"task_id":"task-b","kind":"sora_content_violation","reason_str":"Bad content"placeholder,{"task_id":"task-c","kind":"video","downloadable_url":"https://c.com/c.mp4"placeholder]placeholder`,
			taskID:     "task-b",
			wantFound:  true,
			wantStatus: "failed",
			wantErr:    "Bad content",
	placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, found := testParseGetVideoTaskDrafts([]byte(tt.body), tt.taskID)
			require.Equal(t, tt.wantFound, found, "found 不匹配")
			if !tt.wantFound {
				return
		placeholder
			require.NotNil(t, status)
			require.Equal(t, tt.taskID, status.ID)
			require.Equal(t, tt.wantStatus, status.Status)
			if tt.wantErr != "" {
				require.Equal(t, tt.wantErr, status.ErrorMsg)
		placeholder
			if tt.wantURLs != nil {
				require.Equal(t, tt.wantURLs, status.URLs)
		placeholder
	placeholder)
placeholder
placeholder
