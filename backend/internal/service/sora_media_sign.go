package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

// SignSoraMediaURL 生成 Sora 媒体临时签名
func SignSoraMediaURL(path string, query string, expires int64, key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
placeholder
	mac := hmac.New(sha256.New, []byte(key))
	if _, err := mac.Write([]byte(buildSoraMediaSignPayload(path, query))); err != nil {
		return ""
placeholder
	if _, err := mac.Write([]byte("|")); err != nil {
		return ""
placeholder
	if _, err := mac.Write([]byte(strconv.FormatInt(expires, 10))); err != nil {
		return ""
placeholder
	return hex.EncodeToString(mac.Sum(nil))
placeholder

// VerifySoraMediaURL 校验 Sora 媒体签名
func VerifySoraMediaURL(path string, query string, expires int64, signature string, key string) bool {
	signature = strings.TrimSpace(signature)
	if signature == "" {
		return false
placeholder
	expected := SignSoraMediaURL(path, query, expires, key)
	if expected == "" {
		return false
placeholder
	return hmac.Equal([]byte(signature), []byte(expected))
placeholder

func buildSoraMediaSignPayload(path string, query string) string {
	if strings.TrimSpace(query) == "" {
		return path
placeholder
	return path + "?" + query
placeholder
