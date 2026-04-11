package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccount_IsWebSearchEmulationEnabled_Enabled(t *testing.T) {
	a := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
		Extra:    map[string]any{featureKeyWebSearchEmulation: trueplaceholder,
placeholder
	require.True(t, a.IsWebSearchEmulationEnabled())
placeholder

func TestAccount_IsWebSearchEmulationEnabled_Disabled(t *testing.T) {
	a := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
		Extra:    map[string]any{featureKeyWebSearchEmulation: falseplaceholder,
placeholder
	require.False(t, a.IsWebSearchEmulationEnabled())
placeholder

func TestAccount_IsWebSearchEmulationEnabled_MissingField(t *testing.T) {
	a := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
		Extra:    map[string]any{placeholder,
placeholder
	require.False(t, a.IsWebSearchEmulationEnabled())
placeholder

func TestAccount_IsWebSearchEmulationEnabled_WrongType(t *testing.T) {
	a := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeAPIKey,
		Extra:    map[string]any{featureKeyWebSearchEmulation: "true"placeholder,
placeholder
	require.False(t, a.IsWebSearchEmulationEnabled())
placeholder

func TestAccount_IsWebSearchEmulationEnabled_NilExtra(t *testing.T) {
	a := &Account{Platform: PlatformAnthropic, Type: AccountTypeAPIKey, Extra: nilplaceholder
	require.False(t, a.IsWebSearchEmulationEnabled())
placeholder

func TestAccount_IsWebSearchEmulationEnabled_NilAccount(t *testing.T) {
	var a *Account
	require.False(t, a.IsWebSearchEmulationEnabled())
placeholder

func TestAccount_IsWebSearchEmulationEnabled_NonAnthropicPlatform(t *testing.T) {
	a := &Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Extra:    map[string]any{featureKeyWebSearchEmulation: trueplaceholder,
placeholder
	require.False(t, a.IsWebSearchEmulationEnabled())
placeholder

func TestAccount_IsWebSearchEmulationEnabled_NonAPIKeyType(t *testing.T) {
	a := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeOAuth,
		Extra:    map[string]any{featureKeyWebSearchEmulation: trueplaceholder,
placeholder
	require.False(t, a.IsWebSearchEmulationEnabled())
placeholder
