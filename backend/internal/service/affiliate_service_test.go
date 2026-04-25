//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type affiliateSettingRepoStub struct {
	value string
	err   error
placeholder

func (s *affiliateSettingRepoStub) Get(context.Context, string) (*Setting, error) { return nil, s.err placeholder
func (s *affiliateSettingRepoStub) GetValue(context.Context, string) (string, error) {
	if s.err != nil {
		return "", s.err
placeholder
	return s.value, nil
placeholder
func (s *affiliateSettingRepoStub) Set(context.Context, string, string) error { return s.err placeholder
func (s *affiliateSettingRepoStub) GetMultiple(context.Context, []string) (map[string]string, error) {
	if s.err != nil {
		return nil, s.err
placeholder
	return map[string]string{placeholder, nil
placeholder
func (s *affiliateSettingRepoStub) SetMultiple(context.Context, map[string]string) error {
	return s.err
placeholder
func (s *affiliateSettingRepoStub) GetAll(context.Context) (map[string]string, error) {
	if s.err != nil {
		return nil, s.err
placeholder
	return map[string]string{placeholder, nil
placeholder
func (s *affiliateSettingRepoStub) Delete(context.Context, string) error { return s.err placeholder

func TestAffiliateRebateRatePercentSemantics(t *testing.T) {
	t.Parallel()

	svc := &AffiliateService{settingRepo: &affiliateSettingRepoStub{value: "1"placeholderplaceholder
	rate := svc.loadAffiliateRebateRatePercent(context.Background())
	require.Equal(t, 1.0, rate)

	svc.settingRepo = &affiliateSettingRepoStub{value: "0.2"placeholder
	rate = svc.loadAffiliateRebateRatePercent(context.Background())
	require.Equal(t, 0.2, rate)
placeholder

func TestMaskEmail(t *testing.T) {
	t.Parallel()
	require.Equal(t, "a***@g***.com", maskEmail("alice@gmail.com"))
	require.Equal(t, "x***@d***", maskEmail("x@domain"))
	require.Equal(t, "", maskEmail(""))
placeholder

func TestIsValidAffiliateCodeFormat(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   string
		want bool
placeholder{
		{"valid canonical", "ABCDEFGHJKLM", trueplaceholder,
		{"valid all digits 2-9", "234567892345", trueplaceholder,
		{"valid mixed", "A2B3C4D5E6F7", trueplaceholder,
		{"too short", "ABCDEFGHJKL", falseplaceholder,
		{"too long", "ABCDEFGHJKLMN", falseplaceholder,
		{"contains excluded letter I", "IBCDEFGHJKLM", falseplaceholder,
		{"contains excluded letter O", "OBCDEFGHJKLM", falseplaceholder,
		{"contains excluded digit 0", "0BCDEFGHJKLM", falseplaceholder,
		{"contains excluded digit 1", "1BCDEFGHJKLM", falseplaceholder,
		{"lowercase rejected (caller must ToUpper first)", "abcdefghjklm", falseplaceholder,
		{"empty", "", falseplaceholder,
		{"12-byte utf8 non-ascii", "ÄÄÄÄÄÄ", falseplaceholder, // 6×2 bytes = 12 bytes, bytes out of charset
		{"ascii punctuation", "ABCDEFGHJK.M", falseplaceholder,
		{"whitespace", "ABCDEFGHJK M", falseplaceholder,
placeholder
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, isValidAffiliateCodeFormat(tc.in))
	placeholder)
placeholder
placeholder
