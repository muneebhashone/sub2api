package service

import (
	"math"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Official Imagine image geometry: https://docs.x.ai/developers/model-capabilities/images/generation
var grokImagineAspectRatioValues = []struct {
	label string
	ratio float64
placeholder{
	{"1:1", 1placeholder,
	{"16:9", 16.0 / 9.0placeholder,
	{"9:16", 9.0 / 16.0placeholder,
	{"4:3", 4.0 / 3.0placeholder,
	{"3:4", 3.0 / 4.0placeholder,
	{"3:2", placeholder,
	{"2:3", 2.0 / 3.0placeholder,
	{"2:1", 2placeholder,
	{"1:2", 0.5placeholder,
	{"19.5:9", 19.5 / 9.0placeholder,
	{"9:19.5", 9.0 / 19.5placeholder,
	{"20:9", 20.0 / 9.0placeholder,
	{"9:20", 9.0 / 20.0placeholder,
placeholder

func applyGrokImagineImageGeometry(body []byte) ([]byte, error) {
	size := strings.TrimSpace(gjson.GetBytes(body, "size").String())
	resolution := grokImagineImageResolution(gjson.GetBytes(body, "resolution").String())
	aspect := strings.TrimSpace(gjson.GetBytes(body, "aspect_ratio").String())
	out := append([]byte(nil), body...)

	if resolution == "" {
		if derived := grokImagineImageResolutionFromSize(size); derived != "" {
			next, err := sjson.SetBytes(out, "resolution", derived)
			if err != nil {
				return nil, err
		placeholder
			out = next
	placeholder
placeholder else if gjson.GetBytes(body, "resolution").String() != resolution {
		next, err := sjson.SetBytes(out, "resolution", resolution)
		if err != nil {
			return nil, err
	placeholder
		out = next
placeholder

	if aspect == "" {
		if derived := grokImagineAspectRatioFromSize(size); derived != "" {
			next, err := sjson.SetBytes(out, "aspect_ratio", derived)
			if err != nil {
				return nil, err
		placeholder
			out = next
	placeholder
placeholder

	if !gjson.GetBytes(out, "size").Exists() {
		return out, nil
placeholder
	return sjson.DeleteBytes(out, "size")
placeholder

func assignGrokMediaResolution(value string, info *GrokMediaRequestInfo) {
	if info == nil {
		return
placeholder
	value = strings.TrimSpace(value)
	if value == "" {
		return
placeholder
	if img := grokImagineImageResolution(value); img != "" {
		info.ImageResolution = img
		return
placeholder
	info.Resolution = value
placeholder

func grokImagineImageResolution(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1k":
		return "1k"
	case "2k":
		return "2k"
	default:
		return ""
placeholder
placeholder

func grokImagineImageResolutionFromSize(size string) string {
	if explicit := grokImagineImageResolution(size); explicit != "" {
		return explicit
placeholder
	tier, ok := ClassifyImageBillingTier(size)
	if !ok {
		return ""
placeholder
	if tier == ImageBillingSize1K {
		return "1k"
placeholder
	return "2k"
placeholder

func grokImagineAspectRatioFromSize(size string) string {
	width, height, ok := parseImageBillingDimensions(strings.TrimSpace(size))
	if !ok || width <= 0 || height <= 0 {
		return ""
placeholder
	div := grokImagineGCD(width, height)
	exact := strconv.Itoa(width/div) + ":" + strconv.Itoa(height/div)
	for _, candidate := range grokImagineAspectRatioValues {
		if candidate.label == exact {
			return exact
	placeholder
placeholder
	ratio := float64(width) / float64(height)
	bestLabel := ""
	bestDelta := math.MaxFloat64
	for _, candidate := range grokImagineAspectRatioValues {
		delta := math.Abs(ratio - candidate.ratio)
		if delta < bestDelta {
			bestDelta = delta
			bestLabel = candidate.label
	placeholder
placeholder
	return bestLabel
placeholder

func grokImagineGCD(a, b int) int {
	if a < 0 {
		a = -a
placeholder
	if b < 0 {
		b = -b
placeholder
	for b != 0 {
		a, b = b, a%b
placeholder
	if a == 0 {
		return 1
placeholder
	return a
placeholder
