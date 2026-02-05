package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials     = infraerrors.Unauthorized("INVALID_CREDENTIALS", "invalid email or password")
	ErrUserNotActive          = infraerrors.Forbidden("USER_NOT_ACTIVE", "user is not active")
	ErrEmailExists            = infraerrors.Conflict("EMAIL_EXISTS", "email already exists")
	ErrEmailReserved          = infraerrors.BadRequest("EMAIL_RESERVED", "email is reserved")
	ErrInvalidToken           = infraerrors.Unauthorized("INVALID_TOKEN", "invalid token")
	ErrTokenExpired           = infraerrors.Unauthorized("TOKEN_EXPIRED", "token has expired")
	ErrAccessTokenExpired     = infraerrors.Unauthorized("ACCESS_TOKEN_EXPIRED", "access token has expired")
	ErrTokenTooLarge          = infraerrors.BadRequest("TOKEN_TOO_LARGE", "token too large")
	ErrTokenRevoked           = infraerrors.Unauthorized("TOKEN_REVOKED", "token has been revoked")
	ErrRefreshTokenInvalid    = infraerrors.Unauthorized("REFRESH_TOKEN_INVALID", "invalid refresh token")
	ErrRefreshTokenExpired    = infraerrors.Unauthorized("REFRESH_TOKEN_EXPIRED", "refresh token has expired")
	ErrRefreshTokenReused     = infraerrors.Unauthorized("REFRESH_TOKEN_REUSED", "refresh token has been reused")
	ErrEmailVerifyRequired    = infraerrors.BadRequest("EMAIL_VERIFY_REQUIRED", "email verification is required")
	ErrRegDisabled            = infraerrors.Forbidden("REGISTRATION_DISABLED", "registration is currently disabled")
	ErrServiceUnavailable     = infraerrors.ServiceUnavailable("SERVICE_UNAVAILABLE", "service temporarily unavailable")
	ErrInvitationCodeRequired = infraerrors.BadRequest("INVITATION_CODE_REQUIRED", "invitation code is required")
	ErrInvitationCodeInvalid  = infraerrors.BadRequest("INVITATION_CODE_INVALID", "invalid or used invitation code")
)

// maxTokenLength 限制 token 大小，避免超长 header 触发解析时的异常内存分配。
const maxTokenLength = 8192

// refreshTokenPrefix is the prefix for refresh tokens to distinguish them from access tokens.
const refreshTokenPrefix = "rt_"

// JWTClaims JWT载荷数据
type JWTClaims struct {
	UserID       int64  `json:"user_id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	TokenVersion int64  `json:"token_version"` // Used to invalidate tokens on password change
	jwt.RegisteredClaims
placeholder

// AuthService 认证服务
type AuthService struct {
	userRepo          UserRepository
	redeemRepo        RedeemCodeRepository
	refreshTokenCache RefreshTokenCache
	cfg               *config.Config
	settingService    *SettingService
	emailService      *EmailService
	turnstileService  *TurnstileService
	emailQueueService *EmailQueueService
	promoService      *PromoService
placeholder

// NewAuthService 创建认证服务实例
func NewAuthService(
	userRepo UserRepository,
	redeemRepo RedeemCodeRepository,
	refreshTokenCache RefreshTokenCache,
	cfg *config.Config,
	settingService *SettingService,
	emailService *EmailService,
	turnstileService *TurnstileService,
	emailQueueService *EmailQueueService,
	promoService *PromoService,
) *AuthService {
	return &AuthService{
		userRepo:          userRepo,
		redeemRepo:        redeemRepo,
		refreshTokenCache: refreshTokenCache,
		cfg:               cfg,
		settingService:    settingService,
		emailService:      emailService,
		turnstileService:  turnstileService,
		emailQueueService: emailQueueService,
		promoService:      promoService,
placeholder
placeholder

// Register 用户注册，返回token和用户
func (s *AuthService) Register(ctx context.Context, email, password string) (string, *User, error) {
	return s.RegisterWithVerification(ctx, email, password, "", "", "")
placeholder

// RegisterWithVerification 用户注册（支持邮件验证、优惠码和邀请码），返回token和用户
func (s *AuthService) RegisterWithVerification(ctx context.Context, email, password, verifyCode, promoCode, invitationCode string) (string, *User, error) {
	// 检查是否开放注册（默认关闭：settingService 未配置时不允许注册）
	if s.settingService == nil || !s.settingService.IsRegistrationEnabled(ctx) {
		return "", nil, ErrRegDisabled
placeholder

	// 防止用户注册 LinuxDo OAuth 合成邮箱，避免第三方登录与本地账号发生碰撞。
	if isReservedEmail(email) {
		return "", nil, ErrEmailReserved
placeholder

	// 检查是否需要邀请码
	var invitationRedeemCode *RedeemCode
	if s.settingService != nil && s.settingService.IsInvitationCodeEnabled(ctx) {
		if invitationCode == "" {
			return "", nil, ErrInvitationCodeRequired
	placeholder
		// 验证邀请码
		redeemCode, err := s.redeemRepo.GetByCode(ctx, invitationCode)
		if err != nil {
			log.Printf("[Auth] Invalid invitation code: %s, error: %v", invitationCode, err)
			return "", nil, ErrInvitationCodeInvalid
	placeholder
		// 检查类型和状态
		if redeemCode.Type != RedeemTypeInvitation || redeemCode.Status != StatusUnused {
			log.Printf("[Auth] Invitation code invalid: type=%s, status=%s", redeemCode.Type, redeemCode.Status)
			return "", nil, ErrInvitationCodeInvalid
	placeholder
		invitationRedeemCode = redeemCode
placeholder

	// 检查是否需要邮件验证
	if s.settingService != nil && s.settingService.IsEmailVerifyEnabled(ctx) {
		// 如果邮件验证已开启但邮件服务未配置，拒绝注册
		// 这是一个配置错误，不应该允许绕过验证
		if s.emailService == nil {
			log.Println("[Auth] Email verification enabled but email service not configured, rejecting registration")
			return "", nil, ErrServiceUnavailable
	placeholder
		if verifyCode == "" {
			return "", nil, ErrEmailVerifyRequired
	placeholder
		// 验证邮箱验证码
		if err := s.emailService.VerifyCode(ctx, email, verifyCode); err != nil {
			return "", nil, fmt.Errorf("verify code: %w", err)
	placeholder
placeholder

	// 检查邮箱是否已存在
	existsEmail, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		log.Printf("[Auth] Database error checking email exists: %v", err)
		return "", nil, ErrServiceUnavailable
placeholder
	if existsEmail {
		return "", nil, ErrEmailExists
placeholder

	// 密码哈希
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return "", nil, fmt.Errorf("hash password: %w", err)
placeholder

	// 获取默认配置
	defaultBalance := s.cfg.Default.UserBalance
	defaultConcurrency := s.cfg.Default.UserConcurrency
	if s.settingService != nil {
		defaultBalance = s.settingService.GetDefaultBalance(ctx)
		defaultConcurrency = s.settingService.GetDefaultConcurrency(ctx)
placeholder

	// 创建用户
	user := &User{
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         RoleUser,
		Balance:      defaultBalance,
		Concurrency:  defaultConcurrency,
		Status:       StatusActive,
placeholder

	if err := s.userRepo.Create(ctx, user); err != nil {
		// 优先检查邮箱冲突错误（竞态条件下可能发生）
		if errors.Is(err, ErrEmailExists) {
			return "", nil, ErrEmailExists
	placeholder
		log.Printf("[Auth] Database error creating user: %v", err)
		return "", nil, ErrServiceUnavailable
placeholder

	// 标记邀请码为已使用（如果使用了邀请码）
	if invitationRedeemCode != nil {
		if err := s.redeemRepo.Use(ctx, invitationRedeemCode.ID, user.ID); err != nil {
			// 邀请码标记失败不影响注册，只记录日志
			log.Printf("[Auth] Failed to mark invitation code as used for user %d: %v", user.ID, err)
	placeholder
placeholder
	// 应用优惠码（如果提供且功能已启用）
	if promoCode != "" && s.promoService != nil && s.settingService != nil && s.settingService.IsPromoCodeEnabled(ctx) {
		if err := s.promoService.ApplyPromoCode(ctx, user.ID, promoCode); err != nil {
			// 优惠码应用失败不影响注册，只记录日志
			log.Printf("[Auth] Failed to apply promo code for user %d: %v", user.ID, err)
	placeholder else {
			// 重新获取用户信息以获取更新后的余额
			if updatedUser, err := s.userRepo.GetByID(ctx, user.ID); err == nil {
				user = updatedUser
		placeholder
	placeholder
placeholder

	// 生成token
	token, err := s.GenerateToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("generate token: %w", err)
placeholder

	return token, user, nil
placeholder

// SendVerifyCodeResult 发送验证码返回结果
type SendVerifyCodeResult struct {
	Countdown int `json:"countdown"` // 倒计时秒数
placeholder

// SendVerifyCode 发送邮箱验证码（同步方式）
func (s *AuthService) SendVerifyCode(ctx context.Context, email string) error {
	// 检查是否开放注册（默认关闭）
	if s.settingService == nil || !s.settingService.IsRegistrationEnabled(ctx) {
		return ErrRegDisabled
placeholder

	if isReservedEmail(email) {
		return ErrEmailReserved
placeholder

	// 检查邮箱是否已存在
	existsEmail, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		log.Printf("[Auth] Database error checking email exists: %v", err)
		return ErrServiceUnavailable
placeholder
	if existsEmail {
		return ErrEmailExists
placeholder

	// 发送验证码
	if s.emailService == nil {
		return errors.New("email service not configured")
placeholder

	// 获取网站名称
	siteName := "Sub2API"
	if s.settingService != nil {
		siteName = s.settingService.GetSiteName(ctx)
placeholder

	return s.emailService.SendVerifyCode(ctx, email, siteName)
placeholder

// SendVerifyCodeAsync 异步发送邮箱验证码并返回倒计时
func (s *AuthService) SendVerifyCodeAsync(ctx context.Context, email string) (*SendVerifyCodeResult, error) {
	log.Printf("[Auth] SendVerifyCodeAsync called for email: %s", email)

	// 检查是否开放注册（默认关闭）
	if s.settingService == nil || !s.settingService.IsRegistrationEnabled(ctx) {
		log.Println("[Auth] Registration is disabled")
		return nil, ErrRegDisabled
placeholder

	if isReservedEmail(email) {
		return nil, ErrEmailReserved
placeholder

	// 检查邮箱是否已存在
	existsEmail, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		log.Printf("[Auth] Database error checking email exists: %v", err)
		return nil, ErrServiceUnavailable
placeholder
	if existsEmail {
		log.Printf("[Auth] Email already exists: %s", email)
		return nil, ErrEmailExists
placeholder

	// 检查邮件队列服务是否配置
	if s.emailQueueService == nil {
		log.Println("[Auth] Email queue service not configured")
		return nil, errors.New("email queue service not configured")
placeholder

	// 获取网站名称
	siteName := "Sub2API"
	if s.settingService != nil {
		siteName = s.settingService.GetSiteName(ctx)
placeholder

	// 异步发送
	log.Printf("[Auth] Enqueueing verify code for: %s", email)
	if err := s.emailQueueService.EnqueueVerifyCode(email, siteName); err != nil {
		log.Printf("[Auth] Failed to enqueue: %v", err)
		return nil, fmt.Errorf("enqueue verify code: %w", err)
placeholder

	log.Printf("[Auth] Verify code enqueued successfully for: %s", email)
	return &SendVerifyCodeResult{
		Countdown: 60, // 60秒倒计时
placeholder, nil
placeholder

// VerifyTurnstile 验证Turnstile token
func (s *AuthService) VerifyTurnstile(ctx context.Context, token string, remoteIP string) error {
	required := s.cfg != nil && s.cfg.Server.Mode == "release" && s.cfg.Turnstile.Required

	if required {
		if s.settingService == nil {
			log.Println("[Auth] Turnstile required but settings service is not configured")
			return ErrTurnstileNotConfigured
	placeholder
		enabled := s.settingService.IsTurnstileEnabled(ctx)
		secretConfigured := s.settingService.GetTurnstileSecretKey(ctx) != ""
		if !enabled || !secretConfigured {
			log.Printf("[Auth] Turnstile required but not configured (enabled=%v, secret_configured=%v)", enabled, secretConfigured)
			return ErrTurnstileNotConfigured
	placeholder
placeholder

	if s.turnstileService == nil {
		if required {
			log.Println("[Auth] Turnstile required but service not configured")
			return ErrTurnstileNotConfigured
	placeholder
		return nil // 服务未配置则跳过验证
placeholder

	if !required && s.settingService != nil && s.settingService.IsTurnstileEnabled(ctx) && s.settingService.GetTurnstileSecretKey(ctx) == "" {
		log.Println("[Auth] Turnstile enabled but secret key not configured")
placeholder

	return s.turnstileService.VerifyToken(ctx, token, remoteIP)
placeholder

// IsTurnstileEnabled 检查是否启用Turnstile验证
func (s *AuthService) IsTurnstileEnabled(ctx context.Context) bool {
	if s.turnstileService == nil {
		return false
placeholder
	return s.turnstileService.IsEnabled(ctx)
placeholder

// IsRegistrationEnabled 检查是否开放注册
func (s *AuthService) IsRegistrationEnabled(ctx context.Context) bool {
	if s.settingService == nil {
		return false // 安全默认：settingService 未配置时关闭注册
placeholder
	return s.settingService.IsRegistrationEnabled(ctx)
placeholder

// IsEmailVerifyEnabled 检查是否开启邮件验证
func (s *AuthService) IsEmailVerifyEnabled(ctx context.Context) bool {
	if s.settingService == nil {
		return false
placeholder
	return s.settingService.IsEmailVerifyEnabled(ctx)
placeholder

// Login 用户登录，返回JWT token
func (s *AuthService) Login(ctx context.Context, email, password string) (string, *User, error) {
	// 查找用户
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return "", nil, ErrInvalidCredentials
	placeholder
		// 记录数据库错误但不暴露给用户
		log.Printf("[Auth] Database error during login: %v", err)
		return "", nil, ErrServiceUnavailable
placeholder

	// 验证密码
	if !s.CheckPassword(password, user.PasswordHash) {
		return "", nil, ErrInvalidCredentials
placeholder

	// 检查用户状态
	if !user.IsActive() {
		return "", nil, ErrUserNotActive
placeholder

	// 生成JWT token
	token, err := s.GenerateToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("generate token: %w", err)
placeholder

	return token, user, nil
placeholder

// LoginOrRegisterOAuth 用于第三方 OAuth/SSO 登录：
// - 如果邮箱已存在：直接登录（不需要本地密码）
// - 如果邮箱不存在：创建新用户并登录
//
// 注意：该函数用于 LinuxDo OAuth 登录场景（不同于上游账号的 OAuth，例如 Claude/OpenAI/Gemini）。
// 为了满足现有数据库约束（需要密码哈希），新用户会生成随机密码并进行哈希保存。
func (s *AuthService) LoginOrRegisterOAuth(ctx context.Context, email, username string) (string, *User, error) {
	email = strings.TrimSpace(email)
	if email == "" || len(email) > 255 {
		return "", nil, infraerrors.BadRequest("INVALID_EMAIL", "invalid email")
placeholder
	if _, err := mail.ParseAddress(email); err != nil {
		return "", nil, infraerrors.BadRequest("INVALID_EMAIL", "invalid email")
placeholder

	username = strings.TrimSpace(username)
	if len([]rune(username)) > 100 {
		username = string([]rune(username)[:100])
placeholder

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// OAuth 首次登录视为注册（fail-close：settingService 未配置时不允许注册）
			if s.settingService == nil || !s.settingService.IsRegistrationEnabled(ctx) {
				return "", nil, ErrRegDisabled
		placeholder

			randomPassword, err := randomHexString(32)
			if err != nil {
				log.Printf("[Auth] Failed to generate random password for oauth signup: %v", err)
				return "", nil, ErrServiceUnavailable
		placeholder
			hashedPassword, err := s.HashPassword(randomPassword)
			if err != nil {
				return "", nil, fmt.Errorf("hash password: %w", err)
		placeholder

			// 新用户默认值。
			defaultBalance := s.cfg.Default.UserBalance
			defaultConcurrency := s.cfg.Default.UserConcurrency
			if s.settingService != nil {
				defaultBalance = s.settingService.GetDefaultBalance(ctx)
				defaultConcurrency = s.settingService.GetDefaultConcurrency(ctx)
		placeholder

			newUser := &User{
				Email:        email,
				Username:     username,
				PasswordHash: hashedPassword,
				Role:         RoleUser,
				Balance:      defaultBalance,
				Concurrency:  defaultConcurrency,
				Status:       StatusActive,
		placeholder

			if err := s.userRepo.Create(ctx, newUser); err != nil {
				if errors.Is(err, ErrEmailExists) {
					// 并发场景：GetByEmail 与 Create 之间用户被创建。
					user, err = s.userRepo.GetByEmail(ctx, email)
					if err != nil {
						log.Printf("[Auth] Database error getting user after conflict: %v", err)
						return "", nil, ErrServiceUnavailable
				placeholder
			placeholder else {
					log.Printf("[Auth] Database error creating oauth user: %v", err)
					return "", nil, ErrServiceUnavailable
			placeholder
		placeholder else {
				user = newUser
		placeholder
	placeholder else {
			log.Printf("[Auth] Database error during oauth login: %v", err)
			return "", nil, ErrServiceUnavailable
	placeholder
placeholder

	if !user.IsActive() {
		return "", nil, ErrUserNotActive
placeholder

	// 尽力补全：当用户名为空时，使用第三方返回的用户名回填。
	if user.Username == "" && username != "" {
		user.Username = username
		if err := s.userRepo.Update(ctx, user); err != nil {
			log.Printf("[Auth] Failed to update username after oauth login: %v", err)
	placeholder
placeholder

	token, err := s.GenerateToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("generate token: %w", err)
placeholder
	return token, user, nil
placeholder

// LoginOrRegisterOAuthWithTokenPair 用于第三方 OAuth/SSO 登录，返回完整的 TokenPair
// 与 LoginOrRegisterOAuth 功能相同，但返回 TokenPair 而非单个 token
func (s *AuthService) LoginOrRegisterOAuthWithTokenPair(ctx context.Context, email, username string) (*TokenPair, *User, error) {
	// 检查 refreshTokenCache 是否可用
	if s.refreshTokenCache == nil {
		return nil, nil, errors.New("refresh token cache not configured")
placeholder

	email = strings.TrimSpace(email)
	if email == "" || len(email) > 255 {
		return nil, nil, infraerrors.BadRequest("INVALID_EMAIL", "invalid email")
placeholder
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, nil, infraerrors.BadRequest("INVALID_EMAIL", "invalid email")
placeholder

	username = strings.TrimSpace(username)
	if len([]rune(username)) > 100 {
		username = string([]rune(username)[:100])
placeholder

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// OAuth 首次登录视为注册
			if s.settingService == nil || !s.settingService.IsRegistrationEnabled(ctx) {
				return nil, nil, ErrRegDisabled
		placeholder

			randomPassword, err := randomHexString(32)
			if err != nil {
				log.Printf("[Auth] Failed to generate random password for oauth signup: %v", err)
				return nil, nil, ErrServiceUnavailable
		placeholder
			hashedPassword, err := s.HashPassword(randomPassword)
			if err != nil {
				return nil, nil, fmt.Errorf("hash password: %w", err)
		placeholder

			defaultBalance := s.cfg.Default.UserBalance
			defaultConcurrency := s.cfg.Default.UserConcurrency
			if s.settingService != nil {
				defaultBalance = s.settingService.GetDefaultBalance(ctx)
				defaultConcurrency = s.settingService.GetDefaultConcurrency(ctx)
		placeholder

			newUser := &User{
				Email:        email,
				Username:     username,
				PasswordHash: hashedPassword,
				Role:         RoleUser,
				Balance:      defaultBalance,
				Concurrency:  defaultConcurrency,
				Status:       StatusActive,
		placeholder

			if err := s.userRepo.Create(ctx, newUser); err != nil {
				if errors.Is(err, ErrEmailExists) {
					user, err = s.userRepo.GetByEmail(ctx, email)
					if err != nil {
						log.Printf("[Auth] Database error getting user after conflict: %v", err)
						return nil, nil, ErrServiceUnavailable
				placeholder
			placeholder else {
					log.Printf("[Auth] Database error creating oauth user: %v", err)
					return nil, nil, ErrServiceUnavailable
			placeholder
		placeholder else {
				user = newUser
		placeholder
	placeholder else {
			log.Printf("[Auth] Database error during oauth login: %v", err)
			return nil, nil, ErrServiceUnavailable
	placeholder
placeholder

	if !user.IsActive() {
		return nil, nil, ErrUserNotActive
placeholder

	if user.Username == "" && username != "" {
		user.Username = username
		if err := s.userRepo.Update(ctx, user); err != nil {
			log.Printf("[Auth] Failed to update username after oauth login: %v", err)
	placeholder
placeholder

	tokenPair, err := s.GenerateTokenPair(ctx, user, "")
	if err != nil {
		return nil, nil, fmt.Errorf("generate token pair: %w", err)
placeholder
	return tokenPair, user, nil
placeholder

// ValidateToken 验证JWT token并返回用户声明
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	// 先做长度校验，尽早拒绝异常超长 token，降低 DoS 风险。
	if len(tokenString) > maxTokenLength {
		return nil, ErrTokenTooLarge
placeholder

	// 使用解析器并限制可接受的签名算法，防止算法混淆。
	parser := jwt.NewParser(jwt.WithValidMethods([]string{
		jwt.SigningMethodHS256.Name,
		jwt.SigningMethodHS384.Name,
		jwt.SigningMethodHS512.Name,
placeholder))

	// 保留默认 claims 校验（exp/nbf），避免放行过期或未生效的 token。
	token, err := parser.ParseWithClaims(tokenString, &JWTClaims{placeholder, func(token *jwt.Token) (any, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	placeholder
		return []byte(s.cfg.JWT.Secret), nil
placeholder)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			// token 过期但仍返回 claims（用于 RefreshToken 等场景）
			// jwt-go 在解析时即使遇到过期错误，token.Claims 仍会被填充
			if claims, ok := token.Claims.(*JWTClaims); ok {
				return claims, ErrTokenExpired
		placeholder
			return nil, ErrTokenExpired
	placeholder
		return nil, ErrInvalidToken
placeholder

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
placeholder

	return nil, ErrInvalidToken
placeholder

func randomHexString(byteLength int) (string, error) {
	if byteLength <= 0 {
		byteLength = 16
placeholder
	buf := make([]byte, byteLength)
	if _, err := rand.Read(buf); err != nil {
		return "", err
placeholder
	return hex.EncodeToString(buf), nil
placeholder

func isReservedEmail(email string) bool {
	normalized := strings.ToLower(strings.TrimSpace(email))
	return strings.HasSuffix(normalized, LinuxDoConnectSyntheticEmailDomain)
placeholder

// GenerateToken 生成JWT access token
// 使用新的access_token_expire_minutes配置项（如果配置了），否则回退到expire_hour
func (s *AuthService) GenerateToken(user *User) (string, error) {
	now := time.Now()
	var expiresAt time.Time
	if s.cfg.JWT.AccessTokenExpireMinutes > 0 {
		expiresAt = now.Add(time.Duration(s.cfg.JWT.AccessTokenExpireMinutes) * time.Minute)
placeholder else {
		// 向后兼容：使用旧的expire_hour配置
		expiresAt = now.Add(time.Duration(s.cfg.JWT.ExpireHour) * time.Hour)
placeholder

	claims := &JWTClaims{
		UserID:       user.ID,
		Email:        user.Email,
		Role:         user.Role,
		TokenVersion: user.TokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
	placeholder,
placeholder

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
placeholder

	return tokenString, nil
placeholder

// GetAccessTokenExpiresIn 返回Access Token的有效期（秒）
// 用于前端设置刷新定时器
func (s *AuthService) GetAccessTokenExpiresIn() int {
	if s.cfg.JWT.AccessTokenExpireMinutes > 0 {
		return s.cfg.JWT.AccessTokenExpireMinutes * 60
placeholder
	return s.cfg.JWT.ExpireHour * 3600
placeholder

// HashPassword 使用bcrypt加密密码
func (s *AuthService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
placeholder
	return string(hashedBytes), nil
placeholder

// CheckPassword 验证密码是否匹配
func (s *AuthService) CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
placeholder

// RefreshToken 刷新token
func (s *AuthService) RefreshToken(ctx context.Context, oldTokenString string) (string, error) {
	// 验证旧token（即使过期也允许，用于刷新）
	claims, err := s.ValidateToken(oldTokenString)
	if err != nil && !errors.Is(err, ErrTokenExpired) {
		return "", err
placeholder

	// 获取最新的用户信息
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return "", ErrInvalidToken
	placeholder
		log.Printf("[Auth] Database error refreshing token: %v", err)
		return "", ErrServiceUnavailable
placeholder

	// 检查用户状态
	if !user.IsActive() {
		return "", ErrUserNotActive
placeholder

	// Security: Check TokenVersion to prevent refreshing revoked tokens
	// This ensures tokens issued before a password change cannot be refreshed
	if claims.TokenVersion != user.TokenVersion {
		return "", ErrTokenRevoked
placeholder

	// 生成新token
	return s.GenerateToken(user)
placeholder

// IsPasswordResetEnabled 检查是否启用密码重置功能
// 要求：必须同时开启邮件验证且 SMTP 配置正确
func (s *AuthService) IsPasswordResetEnabled(ctx context.Context) bool {
	if s.settingService == nil {
		return false
placeholder
	// Must have email verification enabled and SMTP configured
	if !s.settingService.IsEmailVerifyEnabled(ctx) {
		return false
placeholder
	return s.settingService.IsPasswordResetEnabled(ctx)
placeholder

// preparePasswordReset validates the password reset request and returns necessary data
// Returns (siteName, resetURL, shouldProceed)
// shouldProceed is false when we should silently return success (to prevent enumeration)
func (s *AuthService) preparePasswordReset(ctx context.Context, email, frontendBaseURL string) (string, string, bool) {
	// Check if user exists (but don't reveal this to the caller)
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// Security: Log but don't reveal that user doesn't exist
			log.Printf("[Auth] Password reset requested for non-existent email: %s", email)
			return "", "", false
	placeholder
		log.Printf("[Auth] Database error checking email for password reset: %v", err)
		return "", "", false
placeholder

	// Check if user is active
	if !user.IsActive() {
		log.Printf("[Auth] Password reset requested for inactive user: %s", email)
		return "", "", false
placeholder

	// Get site name
	siteName := "Sub2API"
	if s.settingService != nil {
		siteName = s.settingService.GetSiteName(ctx)
placeholder

	// Build reset URL base
	resetURL := fmt.Sprintf("%s/reset-password", strings.TrimSuffix(frontendBaseURL, "/"))

	return siteName, resetURL, true
placeholder

// RequestPasswordReset 请求密码重置（同步发送）
// Security: Returns the same response regardless of whether the email exists (prevent user enumeration)
func (s *AuthService) RequestPasswordReset(ctx context.Context, email, frontendBaseURL string) error {
	if !s.IsPasswordResetEnabled(ctx) {
		return infraerrors.Forbidden("PASSWORD_RESET_DISABLED", "password reset is not enabled")
placeholder
	if s.emailService == nil {
		return ErrServiceUnavailable
placeholder

	siteName, resetURL, shouldProceed := s.preparePasswordReset(ctx, email, frontendBaseURL)
	if !shouldProceed {
		return nil // Silent success to prevent enumeration
placeholder

	if err := s.emailService.SendPasswordResetEmail(ctx, email, siteName, resetURL); err != nil {
		log.Printf("[Auth] Failed to send password reset email to %s: %v", email, err)
		return nil // Silent success to prevent enumeration
placeholder

	log.Printf("[Auth] Password reset email sent to: %s", email)
	return nil
placeholder

// RequestPasswordResetAsync 异步请求密码重置（队列发送）
// Security: Returns the same response regardless of whether the email exists (prevent user enumeration)
func (s *AuthService) RequestPasswordResetAsync(ctx context.Context, email, frontendBaseURL string) error {
	if !s.IsPasswordResetEnabled(ctx) {
		return infraerrors.Forbidden("PASSWORD_RESET_DISABLED", "password reset is not enabled")
placeholder
	if s.emailQueueService == nil {
		return ErrServiceUnavailable
placeholder

	siteName, resetURL, shouldProceed := s.preparePasswordReset(ctx, email, frontendBaseURL)
	if !shouldProceed {
		return nil // Silent success to prevent enumeration
placeholder

	if err := s.emailQueueService.EnqueuePasswordReset(email, siteName, resetURL); err != nil {
		log.Printf("[Auth] Failed to enqueue password reset email for %s: %v", email, err)
		return nil // Silent success to prevent enumeration
placeholder

	log.Printf("[Auth] Password reset email enqueued for: %s", email)
	return nil
placeholder

// ResetPassword 重置密码
// Security: Increments TokenVersion to invalidate all existing JWT tokens
func (s *AuthService) ResetPassword(ctx context.Context, email, token, newPassword string) error {
	// Check if password reset is enabled
	if !s.IsPasswordResetEnabled(ctx) {
		return infraerrors.Forbidden("PASSWORD_RESET_DISABLED", "password reset is not enabled")
placeholder

	if s.emailService == nil {
		return ErrServiceUnavailable
placeholder

	// Verify and consume the reset token (one-time use)
	if err := s.emailService.ConsumePasswordResetToken(ctx, email, token); err != nil {
		return err
placeholder

	// Get user
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return ErrInvalidResetToken // Token was valid but user was deleted
	placeholder
		log.Printf("[Auth] Database error getting user for password reset: %v", err)
		return ErrServiceUnavailable
placeholder

	// Check if user is active
	if !user.IsActive() {
		return ErrUserNotActive
placeholder

	// Hash new password
	hashedPassword, err := s.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
placeholder

	// Update password and increment TokenVersion
	user.PasswordHash = hashedPassword
	user.TokenVersion++ // Invalidate all existing tokens

	if err := s.userRepo.Update(ctx, user); err != nil {
		log.Printf("[Auth] Database error updating password for user %d: %v", user.ID, err)
		return ErrServiceUnavailable
placeholder

	// Also revoke all refresh tokens for this user
	if err := s.RevokeAllUserSessions(ctx, user.ID); err != nil {
		log.Printf("[Auth] Failed to revoke refresh tokens for user %d: %v", user.ID, err)
		// Don't return error - password was already changed successfully
placeholder

	log.Printf("[Auth] Password reset successful for user: %s", email)
	return nil
placeholder

// ==================== Refresh Token Methods ====================

// TokenPair 包含Access Token和Refresh Token
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // Access Token有效期（秒）
placeholder

// GenerateTokenPair 生成Access Token和Refresh Token对
// familyID: 可选的Token家族ID，用于Token轮转时保持家族关系
func (s *AuthService) GenerateTokenPair(ctx context.Context, user *User, familyID string) (*TokenPair, error) {
	// 检查 refreshTokenCache 是否可用
	if s.refreshTokenCache == nil {
		return nil, errors.New("refresh token cache not configured")
placeholder

	// 生成Access Token
	accessToken, err := s.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
placeholder

	// 生成Refresh Token
	refreshToken, err := s.generateRefreshToken(ctx, user, familyID)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
placeholder

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.GetAccessTokenExpiresIn(),
placeholder, nil
placeholder

// generateRefreshToken 生成并存储Refresh Token
func (s *AuthService) generateRefreshToken(ctx context.Context, user *User, familyID string) (string, error) {
	// 生成随机Token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
placeholder
	rawToken := refreshTokenPrefix + hex.EncodeToString(tokenBytes)

	// 计算Token哈希（存储哈希而非原始Token）
	tokenHash := hashToken(rawToken)

	// 如果没有提供familyID，生成新的
	if familyID == "" {
		familyBytes := make([]byte, 16)
		if _, err := rand.Read(familyBytes); err != nil {
			return "", fmt.Errorf("generate family id: %w", err)
	placeholder
		familyID = hex.EncodeToString(familyBytes)
placeholder

	now := time.Now()
	ttl := time.Duration(s.cfg.JWT.RefreshTokenExpireDays) * 24 * time.Hour

	data := &RefreshTokenData{
		UserID:       user.ID,
		TokenVersion: user.TokenVersion,
		FamilyID:     familyID,
		CreatedAt:    now,
		ExpiresAt:    now.Add(ttl),
placeholder

	// 存储Token数据
	if err := s.refreshTokenCache.StoreRefreshToken(ctx, tokenHash, data, ttl); err != nil {
		return "", fmt.Errorf("store refresh token: %w", err)
placeholder

	// 添加到用户Token集合
	if err := s.refreshTokenCache.AddToUserTokenSet(ctx, user.ID, tokenHash, ttl); err != nil {
		log.Printf("[Auth] Failed to add token to user set: %v", err)
		// 不影响主流程
placeholder

	// 添加到家族Token集合
	if err := s.refreshTokenCache.AddToFamilyTokenSet(ctx, familyID, tokenHash, ttl); err != nil {
		log.Printf("[Auth] Failed to add token to family set: %v", err)
		// 不影响主流程
placeholder

	return rawToken, nil
placeholder

// RefreshTokenPair 使用Refresh Token刷新Token对
// 实现Token轮转：每次刷新都会生成新的Refresh Token，旧Token立即失效
func (s *AuthService) RefreshTokenPair(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// 检查 refreshTokenCache 是否可用
	if s.refreshTokenCache == nil {
		return nil, ErrRefreshTokenInvalid
placeholder

	// 验证Token格式
	if !strings.HasPrefix(refreshToken, refreshTokenPrefix) {
		return nil, ErrRefreshTokenInvalid
placeholder

	tokenHash := hashToken(refreshToken)

	// 获取Token数据
	data, err := s.refreshTokenCache.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, ErrRefreshTokenNotFound) {
			// Token不存在，可能是已被使用（Token轮转）或已过期
			log.Printf("[Auth] Refresh token not found, possible reuse attack")
			return nil, ErrRefreshTokenInvalid
	placeholder
		log.Printf("[Auth] Error getting refresh token: %v", err)
		return nil, ErrServiceUnavailable
placeholder

	// 检查Token是否过期
	if time.Now().After(data.ExpiresAt) {
		// 删除过期Token
		_ = s.refreshTokenCache.DeleteRefreshToken(ctx, tokenHash)
		return nil, ErrRefreshTokenExpired
placeholder

	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, data.UserID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// 用户已删除，撤销整个Token家族
			_ = s.refreshTokenCache.DeleteTokenFamily(ctx, data.FamilyID)
			return nil, ErrRefreshTokenInvalid
	placeholder
		log.Printf("[Auth] Database error getting user for token refresh: %v", err)
		return nil, ErrServiceUnavailable
placeholder

	// 检查用户状态
	if !user.IsActive() {
		// 用户被禁用，撤销整个Token家族
		_ = s.refreshTokenCache.DeleteTokenFamily(ctx, data.FamilyID)
		return nil, ErrUserNotActive
placeholder

	// 检查TokenVersion（密码更改后所有Token失效）
	if data.TokenVersion != user.TokenVersion {
		// TokenVersion不匹配，撤销整个Token家族
		_ = s.refreshTokenCache.DeleteTokenFamily(ctx, data.FamilyID)
		return nil, ErrTokenRevoked
placeholder

	// Token轮转：立即使旧Token失效
	if err := s.refreshTokenCache.DeleteRefreshToken(ctx, tokenHash); err != nil {
		log.Printf("[Auth] Failed to delete old refresh token: %v", err)
		// 继续处理，不影响主流程
placeholder

	// 生成新的Token对，保持同一个家族ID
	return s.GenerateTokenPair(ctx, user, data.FamilyID)
placeholder

// RevokeRefreshToken 撤销单个Refresh Token
func (s *AuthService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	if s.refreshTokenCache == nil {
		return nil // No-op if cache not configured
placeholder
	if !strings.HasPrefix(refreshToken, refreshTokenPrefix) {
		return ErrRefreshTokenInvalid
placeholder

	tokenHash := hashToken(refreshToken)
	return s.refreshTokenCache.DeleteRefreshToken(ctx, tokenHash)
placeholder

// RevokeAllUserSessions 撤销用户的所有会话（所有Refresh Token）
// 用于密码更改或用户主动登出所有设备
func (s *AuthService) RevokeAllUserSessions(ctx context.Context, userID int64) error {
	if s.refreshTokenCache == nil {
		return nil // No-op if cache not configured
placeholder
	return s.refreshTokenCache.DeleteUserRefreshTokens(ctx, userID)
placeholder

// hashToken 计算Token的SHA256哈希
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
placeholder
