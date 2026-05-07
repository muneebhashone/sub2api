package service

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

func ExtractContentModerationText(protocol string, body []byte) string {
	return ExtractContentModerationInput(protocol, body).Text
placeholder

func ExtractContentModerationInput(protocol string, body []byte) ContentModerationInput {
	if len(body) == 0 || !gjson.ValidBytes(body) {
		return ContentModerationInput{placeholder
placeholder
	var parts []string
	var images []string
	switch protocol {
	case ContentModerationProtocolAnthropicMessages:
		collectLastAnthropicUserMessage(gjson.GetBytes(body, "messages"), &parts, &images)
	case ContentModerationProtocolOpenAIChat:
		collectLastRoleMessage(gjson.GetBytes(body, "messages"), "user", &parts, &images)
	case ContentModerationProtocolOpenAIResponses:
		collectLastResponsesInput(gjson.GetBytes(body, "input"), &parts, &images)
	case ContentModerationProtocolGemini:
		collectLastGeminiContent(gjson.GetBytes(body, "contents"), &parts, &images)
	case ContentModerationProtocolOpenAIImages:
		addModerationText(&parts, gjson.GetBytes(body, "prompt").String())
		collectContentValue(gjson.GetBytes(body, "images"), &parts, &images)
	default:
		collectLastResponsesInput(gjson.GetBytes(body, "input"), &parts, &images)
		collectLastRoleMessage(gjson.GetBytes(body, "messages"), "user", &parts, &images)
		collectLastGeminiContent(gjson.GetBytes(body, "contents"), &parts, &images)
placeholder
	out := ContentModerationInput{
		Text:   normalizeContentModerationText(strings.Join(parts, "\n")),
		Images: normalizeModerationImages(images),
placeholder
	out.Normalize()
	return out
placeholder

func collectLastRoleMessage(messages gjson.Result, role string, parts *[]string, images *[]string) {
	if !messages.IsArray() {
		return
placeholder
	var lastParts []string
	var lastImages []string
	messages.ForEach(func(_, msg gjson.Result) bool {
		if strings.ToLower(strings.TrimSpace(msg.Get("role").String())) == role {
			var candidate []string
			var candidateImages []string
			collectContentValue(msg.Get("content"), &candidate, &candidateImages)
			if normalizeContentModerationText(strings.Join(candidate, "\n")) != "" || len(candidateImages) > 0 {
				lastParts = candidate
				lastImages = candidateImages
		placeholder
	placeholder
		return true
placeholder)
	*parts = append(*parts, lastParts...)
	*images = append(*images, lastImages...)
placeholder

func collectLastAnthropicUserMessage(messages gjson.Result, parts *[]string, images *[]string) {
	if !messages.IsArray() {
		return
placeholder
	var lastParts []string
	var lastImages []string
	messages.ForEach(func(_, msg gjson.Result) bool {
		if strings.ToLower(strings.TrimSpace(msg.Get("role").String())) == "user" {
			var candidate []string
			var candidateImages []string
			collectAnthropicUserContentValue(msg.Get("content"), &candidate, &candidateImages)
			if normalizeContentModerationText(strings.Join(candidate, "\n")) != "" || len(candidateImages) > 0 {
				lastParts = candidate
				lastImages = candidateImages
		placeholder
	placeholder
		return true
placeholder)
	*parts = append(*parts, lastParts...)
	*images = append(*images, lastImages...)
placeholder

func collectAnthropicUserContentValue(value gjson.Result, parts *[]string, images *[]string) {
	switch {
	case !value.Exists():
		return
	case value.Type == gjson.String:
		if !isAnthropicSystemReminderText(value.String()) {
			addModerationText(parts, value.String())
	placeholder
	case value.IsArray():
		value.ForEach(func(_, item gjson.Result) bool {
			collectAnthropicUserContentValue(item, parts, images)
			return true
	placeholder)
	case value.IsObject():
		typ := strings.ToLower(strings.TrimSpace(value.Get("type").String()))
		switch typ {
		case "", "text", "input_text", "message":
			if value.Get("text").Exists() && !isAnthropicSystemReminderText(value.Get("text").String()) {
				addModerationText(parts, value.Get("text").String())
		placeholder
			if value.Get("content").Exists() {
				collectAnthropicUserContentValue(value.Get("content"), parts, images)
		placeholder
		case "image_url", "input_image", "image":
			collectContentValue(value, parts, images)
	placeholder
placeholder
placeholder

func isAnthropicSystemReminderText(text string) bool {
	return strings.HasPrefix(strings.TrimSpace(text), "<system-reminder>")
placeholder

func collectLastResponsesInput(input gjson.Result, parts *[]string, images *[]string) {
	switch {
	case !input.Exists():
		return
	case input.Type == gjson.String:
		addModerationText(parts, input.String())
	case input.IsArray():
		var last gjson.Result
		input.ForEach(func(_, item gjson.Result) bool {
			if isResponsesUserTextItem(item) {
				last = item
		placeholder
			return true
	placeholder)
		if last.Exists() {
			collectContentValue(last.Get("content"), parts, images)
			if last.Get("type").String() == "input_text" || last.Get("text").Exists() {
				collectContentValue(last, parts, images)
		placeholder
	placeholder
	case input.IsObject():
		if isResponsesUserTextItem(input) {
			collectContentValue(input.Get("content"), parts, images)
			if input.Get("type").String() == "input_text" || input.Get("text").Exists() {
				collectContentValue(input, parts, images)
		placeholder
	placeholder
placeholder
placeholder

func isResponsesUserTextItem(item gjson.Result) bool {
	role := strings.ToLower(strings.TrimSpace(item.Get("role").String()))
	if role == "user" {
		return responseItemHasModerationText(item)
placeholder
	if role != "" {
		return false
placeholder
	return responseItemHasModerationText(item)
placeholder

func responseItemHasModerationText(item gjson.Result) bool {
	var parts []string
	var images []string
	collectContentValue(item.Get("content"), &parts, &images)
	if item.Get("type").String() == "input_text" || item.Get("text").Exists() {
		collectContentValue(item, &parts, &images)
placeholder
	return normalizeContentModerationText(strings.Join(parts, "\n")) != "" || len(images) > 0
placeholder

func collectLastGeminiContent(contents gjson.Result, parts *[]string, images *[]string) {
	if !contents.IsArray() {
		return
placeholder
	var lastParts []string
	var lastImages []string
	contents.ForEach(func(_, content gjson.Result) bool {
		role := strings.ToLower(strings.TrimSpace(content.Get("role").String()))
		if role == "" || role == "user" {
			var candidate []string
			var candidateImages []string
			if arr := content.Get("parts"); arr.IsArray() {
				arr.ForEach(func(_, part gjson.Result) bool {
					addModerationText(&candidate, part.Get("text").String())
					addGeminiModerationImage(&candidateImages, part)
					return true
			placeholder)
		placeholder
			if normalizeContentModerationText(strings.Join(candidate, "\n")) != "" || len(candidateImages) > 0 {
				lastParts = candidate
				lastImages = candidateImages
		placeholder
	placeholder
		return true
placeholder)
	*parts = append(*parts, lastParts...)
	*images = append(*images, lastImages...)
placeholder

func collectContentValue(value gjson.Result, parts *[]string, images *[]string) {
	switch {
	case !value.Exists():
		return
	case value.Type == gjson.String:
		addModerationText(parts, value.String())
	case value.IsArray():
		value.ForEach(func(_, item gjson.Result) bool {
			collectContentValue(item, parts, images)
			return true
	placeholder)
	case value.IsObject():
		typ := strings.ToLower(strings.TrimSpace(value.Get("type").String()))
		addModerationImage(images, value.Get("image_url.url").String())
		addModerationImage(images, value.Get("image_url").String())
		addModerationImage(images, value.Get("url").String())
		addModerationImageData(images, value.Get("source.media_type").String(), value.Get("source.data").String())
		addModerationImageData(images, value.Get("source.mediaType").String(), value.Get("source.data").String())
		addModerationImageData(images, value.Get("media_type").String(), value.Get("data").String())
		addModerationImageData(images, value.Get("mime_type").String(), value.Get("data").String())
		addModerationImageData(images, value.Get("mimeType").String(), value.Get("data").String())
		addModerationImage(images, value.Get("source.data").String())
		addModerationImage(images, value.Get("data").String())
		addModerationImage(images, value.Get("base64").String())
		switch typ {
		case "", "text", "input_text", "message":
			if value.Get("text").Exists() {
				addModerationText(parts, value.Get("text").String())
		placeholder
			if value.Get("content").Exists() {
				collectContentValue(value.Get("content"), parts, images)
		placeholder
		case "image_url", "input_image", "image":
	placeholder
placeholder
placeholder

func addGeminiModerationImage(images *[]string, part gjson.Result) {
	if inlineData := part.Get("inline_data"); inlineData.IsObject() {
		mimeType := strings.TrimSpace(inlineData.Get("mime_type").String())
		data := strings.TrimSpace(inlineData.Get("data").String())
		if mimeType != "" && data != "" {
			addModerationImage(images, fmt.Sprintf("data:%s;base64,%s", mimeType, data))
	placeholder
placeholder
	if inlineData := part.Get("inlineData"); inlineData.IsObject() {
		mimeType := strings.TrimSpace(inlineData.Get("mimeType").String())
		data := strings.TrimSpace(inlineData.Get("data").String())
		if mimeType != "" && data != "" {
			addModerationImage(images, fmt.Sprintf("data:%s;base64,%s", mimeType, data))
	placeholder
placeholder
	addModerationImage(images, part.Get("file_data.file_uri").String())
	addModerationImage(images, part.Get("fileData.fileUri").String())
placeholder

func addModerationImageData(images *[]string, mimeType string, data string) {
	mimeType = strings.TrimSpace(mimeType)
	data = strings.TrimSpace(data)
	if mimeType == "" || data == "" {
		return
placeholder
	addModerationImage(images, fmt.Sprintf("data:%s;base64,%s", mimeType, data))
placeholder

func addModerationImage(images *[]string, image string) {
	image = strings.TrimSpace(image)
	if image == "" {
		return
placeholder
	if strings.HasPrefix(image, "data:") || strings.HasPrefix(image, "http://") || strings.HasPrefix(image, "https://") {
		*images = append(*images, image)
placeholder
placeholder

func normalizeModerationImages(images []string) []string {
	out := make([]string, 0, len(images))
	seen := make(map[string]struct{placeholder, len(images))
	for _, image := range images {
		image = strings.TrimSpace(image)
		if image == "" {
			continue
	placeholder
		if _, ok := seen[image]; ok {
			continue
	placeholder
		seen[image] = struct{placeholder{placeholder
		out = append(out, image)
placeholder
	return out
placeholder

func addModerationText(parts *[]string, text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
placeholder
	if strings.Contains(text, "<system-reminder>") {
		return
placeholder
	*parts = append(*parts, text)
placeholder

func normalizeContentModerationText(text string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(text)), " ")
placeholder
