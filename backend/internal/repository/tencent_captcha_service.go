package repository

import (
	"context"
	"fmt"

	"github.com/Wei-Shaw/sub2api/internal/service"
	captcha "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/captcha/v20190722"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const tencentCaptchaEndpoint = "captcha.tencentcloudapi.com"

type tencentCaptchaAPI interface {
	DescribeCaptchaResultWithContext(context.Context, *captcha.DescribeCaptchaResultRequest) (*captcha.DescribeCaptchaResultResponse, error)
placeholder

type tencentCaptchaClientFactory func(secretID, secretKey string) (tencentCaptchaAPI, error)

type tencentCaptchaVerifier struct {
	newClient tencentCaptchaClientFactory
placeholder

func NewTencentCaptchaVerifier() service.TencentCaptchaVerifier {
	return &tencentCaptchaVerifier{newClient: newTencentCaptchaSDKClientplaceholder
placeholder

func newTencentCaptchaSDKClient(secretID, secretKey string) (tencentCaptchaAPI, error) {
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = tencentCaptchaEndpoint
	clientProfile.HttpProfile.ReqMethod = "POST"
	clientProfile.HttpProfile.ReqTimeout = 5
	return captcha.NewClient(common.NewCredential(secretID, secretKey), "", clientProfile)
placeholder

func (v *tencentCaptchaVerifier) VerifyTicket(ctx context.Context, credentials service.TencentCaptchaCredentials, proof service.TencentCaptchaProof, remoteIP string) (*service.TencentCaptchaVerifyResponse, error) {
	client, err := v.newClient(credentials.CloudSecretID, credentials.CloudSecretKey)
	if err != nil {
		return nil, fmt.Errorf("create tencent captcha client: %w", err)
placeholder
	request := captcha.NewDescribeCaptchaResultRequest()
	request.CaptchaType = common.Uint64Ptr(9)
	request.Ticket = common.StringPtr(proof.Ticket)
	request.UserIp = common.StringPtr(remoteIP)
	request.Randstr = common.StringPtr(proof.Randstr)
	request.CaptchaAppId = common.Uint64Ptr(credentials.AppID)
	request.AppSecretKey = common.StringPtr(credentials.AppSecretKey)

	response, err := client.DescribeCaptchaResultWithContext(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("describe captcha result: %w", err)
placeholder
	if response == nil || response.Response == nil {
		return nil, fmt.Errorf("describe captcha result: empty response")
placeholder
	return &service.TencentCaptchaVerifyResponse{
		CaptchaCode: valueOrZero(response.Response.CaptchaCode),
		CaptchaMsg:  valueOrEmpty(response.Response.CaptchaMsg),
		RequestID:   valueOrEmpty(response.Response.RequestId),
placeholder, nil
placeholder

func valueOrZero(value *int64) int64 {
	if value == nil {
		return 0
placeholder
	return *value
placeholder

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
placeholder
	return *value
placeholder
