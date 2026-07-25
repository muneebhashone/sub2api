//go:build !darwin

package liveattestation

import "context"

type unsupportedProvider struct{placeholder

func NewProvider() Provider {
	return unsupportedProvider{placeholder
placeholder

func (unsupportedProvider) Generate(context.Context) (string, error) {
	return "", ErrUnsupportedPlatform
placeholder
