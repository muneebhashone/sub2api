package service

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

var (
	ErrAffiliateProfileNotFound = infraerrors.NotFound("AFFILIATE_PROFILE_NOT_FOUND", "affiliate profile not found")
	ErrAffiliateCodeInvalid     = infraerrors.BadRequest("AFFILIATE_CODE_INVALID", "invalid affiliate code")
	ErrAffiliateAlreadyBound    = infraerrors.Conflict("AFFILIATE_ALREADY_BOUND", "affiliate inviter already bound")
	ErrAffiliateQuotaEmpty      = infraerrors.BadRequest("AFFILIATE_QUOTA_EMPTY", "no affiliate quota available to transfer")
)

const (
	affiliateInviteesLimit = 100
)

type AffiliateSummary struct {
	UserID          int64     `json:"user_id"`
	AffCode         string    `json:"aff_code"`
	InviterID       *int64    `json:"inviter_id,omitempty"`
	AffCount        int       `json:"aff_count"`
	AffQuota        float64   `json:"aff_quota"`
	AffHistoryQuota float64   `json:"aff_history_quota"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
placeholder

type AffiliateInvitee struct {
	UserID    int64      `json:"user_id"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
placeholder

type AffiliateDetail struct {
	UserID          int64              `json:"user_id"`
	AffCode         string             `json:"aff_code"`
	InviterID       *int64             `json:"inviter_id,omitempty"`
	AffCount        int                `json:"aff_count"`
	AffQuota        float64            `json:"aff_quota"`
	AffHistoryQuota float64            `json:"aff_history_quota"`
	Invitees        []AffiliateInvitee `json:"invitees"`
placeholder

type AffiliateRepository interface {
	EnsureUserAffiliate(ctx context.Context, userID int64) (*AffiliateSummary, error)
	GetAffiliateByCode(ctx context.Context, code string) (*AffiliateSummary, error)
	BindInviter(ctx context.Context, userID, inviterID int64) (bool, error)
	AccrueQuota(ctx context.Context, inviterID, inviteeUserID int64, amount float64) (bool, error)
	TransferQuotaToBalance(ctx context.Context, userID int64) (float64, float64, error)
	ListInvitees(ctx context.Context, inviterID int64, limit int) ([]AffiliateInvitee, error)
placeholder

type AffiliateService struct {
	repo                 AffiliateRepository
	settingRepo          SettingRepository
	authCacheInvalidator APIKeyAuthCacheInvalidator
	billingCacheService  *BillingCacheService
placeholder

func NewAffiliateService(repo AffiliateRepository, settingRepo SettingRepository, authCacheInvalidator APIKeyAuthCacheInvalidator, billingCacheService *BillingCacheService) *AffiliateService {
	return &AffiliateService{
		repo:                 repo,
		settingRepo:          settingRepo,
		authCacheInvalidator: authCacheInvalidator,
		billingCacheService:  billingCacheService,
placeholder
placeholder

func (s *AffiliateService) EnsureUserAffiliate(ctx context.Context, userID int64) (*AffiliateSummary, error) {
	if userID <= 0 {
		return nil, infraerrors.BadRequest("INVALID_USER", "invalid user")
placeholder
	if s == nil || s.repo == nil {
		return nil, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
placeholder
	return s.repo.EnsureUserAffiliate(ctx, userID)
placeholder

func (s *AffiliateService) GetAffiliateDetail(ctx context.Context, userID int64) (*AffiliateDetail, error) {
	summary, err := s.EnsureUserAffiliate(ctx, userID)
	if err != nil {
		return nil, err
placeholder
	invitees, err := s.listInvitees(ctx, userID)
	if err != nil {
		return nil, err
placeholder
	return &AffiliateDetail{
		UserID:          summary.UserID,
		AffCode:         summary.AffCode,
		InviterID:       summary.InviterID,
		AffCount:        summary.AffCount,
		AffQuota:        summary.AffQuota,
		AffHistoryQuota: summary.AffHistoryQuota,
		Invitees:        invitees,
placeholder, nil
placeholder

func (s *AffiliateService) BindInviterByCode(ctx context.Context, userID int64, rawCode string) error {
	code := strings.ToUpper(strings.TrimSpace(rawCode))
	if code == "" {
		return nil
placeholder
	if s == nil || s.repo == nil {
		return infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
placeholder

	selfSummary, err := s.repo.EnsureUserAffiliate(ctx, userID)
	if err != nil {
		return err
placeholder
	if selfSummary.InviterID != nil {
		return nil
placeholder

	inviterSummary, err := s.repo.GetAffiliateByCode(ctx, code)
	if err != nil {
		if errors.Is(err, ErrAffiliateProfileNotFound) {
			return ErrAffiliateCodeInvalid
	placeholder
		return err
placeholder
	if inviterSummary == nil || inviterSummary.UserID <= 0 || inviterSummary.UserID == userID {
		return ErrAffiliateCodeInvalid
placeholder

	bound, err := s.repo.BindInviter(ctx, userID, inviterSummary.UserID)
	if err != nil {
		return err
placeholder
	if !bound {
		return ErrAffiliateAlreadyBound
placeholder
	return nil
placeholder

func (s *AffiliateService) AccrueInviteRebate(ctx context.Context, inviteeUserID int64, baseRechargeAmount float64) (float64, error) {
	if s == nil || s.repo == nil {
		return 0, nil
placeholder
	if inviteeUserID <= 0 || baseRechargeAmount <= 0 || math.IsNaN(baseRechargeAmount) || math.IsInf(baseRechargeAmount, 0) {
		return 0, nil
placeholder

	inviteeSummary, err := s.repo.EnsureUserAffiliate(ctx, inviteeUserID)
	if err != nil {
		return 0, err
placeholder
	if inviteeSummary.InviterID == nil || *inviteeSummary.InviterID <= 0 {
		return 0, nil
placeholder

	rebateRatePercent := s.loadAffiliateRebateRatePercent(ctx)
	rebate := roundTo(baseRechargeAmount*(rebateRatePercent/100), 8)
	if rebate <= 0 {
		return 0, nil
placeholder

	if _, err := s.repo.EnsureUserAffiliate(ctx, *inviteeSummary.InviterID); err != nil {
		return 0, err
placeholder

	applied, err := s.repo.AccrueQuota(ctx, *inviteeSummary.InviterID, inviteeUserID, rebate)
	if err != nil {
		return 0, err
placeholder
	if !applied {
		return 0, nil
placeholder
	return rebate, nil
placeholder

func (s *AffiliateService) TransferAffiliateQuota(ctx context.Context, userID int64) (float64, float64, error) {
	if s == nil || s.repo == nil {
		return 0, 0, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
placeholder

	transferred, balance, err := s.repo.TransferQuotaToBalance(ctx, userID)
	if err != nil {
		return 0, 0, err
placeholder
	if transferred > 0 {
		s.invalidateAffiliateCaches(ctx, userID)
placeholder
	return transferred, balance, nil
placeholder

func (s *AffiliateService) listInvitees(ctx context.Context, inviterID int64) ([]AffiliateInvitee, error) {
	if s == nil || s.repo == nil {
		return nil, infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "affiliate service unavailable")
placeholder
	invitees, err := s.repo.ListInvitees(ctx, inviterID, affiliateInviteesLimit)
	if err != nil {
		return nil, err
placeholder
	for i := range invitees {
		invitees[i].Email = maskEmail(invitees[i].Email)
placeholder
	return invitees, nil
placeholder

func (s *AffiliateService) loadAffiliateRebateRatePercent(ctx context.Context) float64 {
	if s == nil || s.settingRepo == nil {
		return AffiliateRebateRateDefault
placeholder

	raw, err := s.settingRepo.GetValue(ctx, SettingKeyAffiliateRebateRate)
	if err != nil {
		return AffiliateRebateRateDefault
placeholder

	rate, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return AffiliateRebateRateDefault
placeholder
	if math.IsNaN(rate) || math.IsInf(rate, 0) {
		return AffiliateRebateRateDefault
placeholder
	if rate < AffiliateRebateRateMin {
		return AffiliateRebateRateMin
placeholder
	if rate > AffiliateRebateRateMax {
		return AffiliateRebateRateMax
placeholder
	return rate
placeholder

func roundTo(v float64, scale int) float64 {
	factor := math.Pow10(scale)
	return math.Round(v*factor) / factor
placeholder

func maskEmail(email string) string {
	email = strings.TrimSpace(email)
	if email == "" {
		return ""
placeholder
	at := strings.Index(email, "@")
	if at <= 0 || at >= len(email)-1 {
		return "***"
placeholder

	local := email[:at]
	domain := email[at+1:]
	dot := strings.LastIndex(domain, ".")

	maskedLocal := maskSegment(local)
	if dot <= 0 || dot >= len(domain)-1 {
		return maskedLocal + "@" + maskSegment(domain)
placeholder

	domainName := domain[:dot]
	tld := domain[dot:]
	return maskedLocal + "@" + maskSegment(domainName) + tld
placeholder

func maskSegment(s string) string {
	r := []rune(s)
	if len(r) == 0 {
		return "***"
placeholder
	if len(r) == 1 {
		return string(r[0]) + "***"
placeholder
	return string(r[0]) + "***"
placeholder

func (s *AffiliateService) invalidateAffiliateCaches(ctx context.Context, userID int64) {
	if s.authCacheInvalidator != nil {
		s.authCacheInvalidator.InvalidateAuthCacheByUserID(ctx, userID)
placeholder
	if s.billingCacheService != nil {
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = s.billingCacheService.InvalidateUserBalance(cacheCtx, userID)
	placeholder()
placeholder
placeholder
