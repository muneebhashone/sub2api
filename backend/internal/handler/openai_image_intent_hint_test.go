package handler

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSeedOpenAIForwardImageIntentHint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name          string
		channelMapped bool
		imageIntent   bool
		wantHint      bool
placeholder{
		{name: "seed true", imageIntent: true, wantHint: trueplaceholder,
		{name: "seed false", imageIntent: false, wantHint: trueplaceholder,
		{name: "mapped body stays unknown", channelMapped: true, imageIntent: trueplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &gin.Context{placeholder
			service.SetOpenAIClientTransport(c, service.OpenAIClientTransportHTTP)

			seedOpenAIForwardImageIntentHint(c, tt.channelMapped, tt.imageIntent)

			var hintValues []bool
			for _, value := range c.Keys {
				if hint, ok := value.(bool); ok {
					hintValues = append(hintValues, hint)
			placeholder
		placeholder
			if !tt.wantHint {
				require.Empty(t, hintValues)
				return
		placeholder
			require.Equal(t, []bool{tt.imageIntentplaceholder, hintValues)
	placeholder)
placeholder
placeholder
