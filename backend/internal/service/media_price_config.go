package service

func imagePriceConfigFromAPIKey(apiKey *APIKey) *ImagePriceConfig {
	if apiKey == nil || apiKey.Group == nil {
		return nil
placeholder
	return &ImagePriceConfig{
		Price1K: apiKey.Group.ImagePrice1K,
		Price2K: apiKey.Group.ImagePrice2K,
		Price4K: apiKey.Group.ImagePrice4K,
placeholder
placeholder

func apiKeyHasConfiguredImagePrice(apiKey *APIKey, imageSize string) bool {
	return apiKey != nil && apiKey.Group != nil && apiKey.Group.GetImagePrice(imageSize) != nil
placeholder

func videoPriceConfigFromAPIKey(apiKey *APIKey) *VideoPriceConfig {
	if apiKey == nil || apiKey.Group == nil {
		return nil
placeholder
	return &VideoPriceConfig{
		Price480P:  apiKey.Group.VideoPrice480P,
		Price720P:  apiKey.Group.VideoPrice720P,
		Price1080P: apiKey.Group.VideoPrice1080P,
placeholder
placeholder

func apiKeyHasConfiguredVideoPrice(apiKey *APIKey, resolution string) bool {
	return apiKey != nil && apiKey.Group != nil && apiKey.Group.GetVideoPrice(resolution) != nil
placeholder
