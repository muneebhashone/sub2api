package repository

import (
	"context"

	"github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/errorpassthroughrule"
	"github.com/Wei-Shaw/sub2api/internal/model"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type errorPassthroughRepository struct {
	client *ent.Client
placeholder

// NewErrorPassthroughRepository 创建错误透传规则仓库
func NewErrorPassthroughRepository(client *ent.Client) service.ErrorPassthroughRepository {
	return &errorPassthroughRepository{client: clientplaceholder
placeholder

// List 获取所有规则
func (r *errorPassthroughRepository) List(ctx context.Context) ([]*model.ErrorPassthroughRule, error) {
	rules, err := r.client.ErrorPassthroughRule.Query().
		Order(ent.Asc(errorpassthroughrule.FieldPriority)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder

	result := make([]*model.ErrorPassthroughRule, len(rules))
	for i, rule := range rules {
		result[i] = r.toModel(rule)
placeholder
	return result, nil
placeholder

// GetByID 根据 ID 获取规则
func (r *errorPassthroughRepository) GetByID(ctx context.Context, id int64) (*model.ErrorPassthroughRule, error) {
	rule, err := r.client.ErrorPassthroughRule.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
	placeholder
		return nil, err
placeholder
	return r.toModel(rule), nil
placeholder

// Create 创建规则
func (r *errorPassthroughRepository) Create(ctx context.Context, rule *model.ErrorPassthroughRule) (*model.ErrorPassthroughRule, error) {
	builder := r.client.ErrorPassthroughRule.Create().
		SetName(rule.Name).
		SetEnabled(rule.Enabled).
		SetPriority(rule.Priority).
		SetMatchMode(rule.MatchMode).
		SetPassthroughCode(rule.PassthroughCode).
		SetPassthroughBody(rule.PassthroughBody).
		SetSkipMonitoring(rule.SkipMonitoring)

	if len(rule.ErrorCodes) > 0 {
		builder.SetErrorCodes(rule.ErrorCodes)
placeholder
	if len(rule.Keywords) > 0 {
		builder.SetKeywords(rule.Keywords)
placeholder
	if len(rule.Platforms) > 0 {
		builder.SetPlatforms(rule.Platforms)
placeholder
	if rule.ResponseCode != nil {
		builder.SetResponseCode(*rule.ResponseCode)
placeholder
	if rule.CustomMessage != nil {
		builder.SetCustomMessage(*rule.CustomMessage)
placeholder
	if rule.Description != nil {
		builder.SetDescription(*rule.Description)
placeholder

	created, err := builder.Save(ctx)
	if err != nil {
		return nil, err
placeholder
	return r.toModel(created), nil
placeholder

// Update 更新规则
func (r *errorPassthroughRepository) Update(ctx context.Context, rule *model.ErrorPassthroughRule) (*model.ErrorPassthroughRule, error) {
	builder := r.client.ErrorPassthroughRule.UpdateOneID(rule.ID).
		SetName(rule.Name).
		SetEnabled(rule.Enabled).
		SetPriority(rule.Priority).
		SetMatchMode(rule.MatchMode).
		SetPassthroughCode(rule.PassthroughCode).
		SetPassthroughBody(rule.PassthroughBody).
		SetSkipMonitoring(rule.SkipMonitoring)

	// 处理可选字段
	if len(rule.ErrorCodes) > 0 {
		builder.SetErrorCodes(rule.ErrorCodes)
placeholder else {
		builder.ClearErrorCodes()
placeholder
	if len(rule.Keywords) > 0 {
		builder.SetKeywords(rule.Keywords)
placeholder else {
		builder.ClearKeywords()
placeholder
	if len(rule.Platforms) > 0 {
		builder.SetPlatforms(rule.Platforms)
placeholder else {
		builder.ClearPlatforms()
placeholder
	if rule.ResponseCode != nil {
		builder.SetResponseCode(*rule.ResponseCode)
placeholder else {
		builder.ClearResponseCode()
placeholder
	if rule.CustomMessage != nil {
		builder.SetCustomMessage(*rule.CustomMessage)
placeholder else {
		builder.ClearCustomMessage()
placeholder
	if rule.Description != nil {
		builder.SetDescription(*rule.Description)
placeholder else {
		builder.ClearDescription()
placeholder

	updated, err := builder.Save(ctx)
	if err != nil {
		return nil, err
placeholder
	return r.toModel(updated), nil
placeholder

// Delete 删除规则
func (r *errorPassthroughRepository) Delete(ctx context.Context, id int64) error {
	return r.client.ErrorPassthroughRule.DeleteOneID(id).Exec(ctx)
placeholder

// toModel 将 Ent 实体转换为服务模型
func (r *errorPassthroughRepository) toModel(e *ent.ErrorPassthroughRule) *model.ErrorPassthroughRule {
	rule := &model.ErrorPassthroughRule{
		ID:              int64(e.ID),
		Name:            e.Name,
		Enabled:         e.Enabled,
		Priority:        e.Priority,
		ErrorCodes:      e.ErrorCodes,
		Keywords:        e.Keywords,
		MatchMode:       e.MatchMode,
		Platforms:       e.Platforms,
		PassthroughCode: e.PassthroughCode,
		PassthroughBody: e.PassthroughBody,
		SkipMonitoring:  e.SkipMonitoring,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
placeholder

	if e.ResponseCode != nil {
		rule.ResponseCode = e.ResponseCode
placeholder
	if e.CustomMessage != nil {
		rule.CustomMessage = e.CustomMessage
placeholder
	if e.Description != nil {
		rule.Description = e.Description
placeholder

	// 确保切片不为 nil
	if rule.ErrorCodes == nil {
		rule.ErrorCodes = []int{placeholder
placeholder
	if rule.Keywords == nil {
		rule.Keywords = []string{placeholder
placeholder
	if rule.Platforms == nil {
		rule.Platforms = []string{placeholder
placeholder

	return rule
placeholder
