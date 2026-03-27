package repository

import (
	"context"

	"github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/tlsfingerprintprofile"
	"github.com/Wei-Shaw/sub2api/internal/model"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type tlsFingerprintProfileRepository struct {
	client *ent.Client
placeholder

// NewTLSFingerprintProfileRepository 创建 TLS 指纹模板仓库
func NewTLSFingerprintProfileRepository(client *ent.Client) service.TLSFingerprintProfileRepository {
	return &tlsFingerprintProfileRepository{client: clientplaceholder
placeholder

// List 获取所有模板
func (r *tlsFingerprintProfileRepository) List(ctx context.Context) ([]*model.TLSFingerprintProfile, error) {
	profiles, err := r.client.TLSFingerprintProfile.Query().
		Order(ent.Asc(tlsfingerprintprofile.FieldName)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder

	result := make([]*model.TLSFingerprintProfile, len(profiles))
	for i, p := range profiles {
		result[i] = r.toModel(p)
placeholder
	return result, nil
placeholder

// GetByID 根据 ID 获取模板
func (r *tlsFingerprintProfileRepository) GetByID(ctx context.Context, id int64) (*model.TLSFingerprintProfile, error) {
	p, err := r.client.TLSFingerprintProfile.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
	placeholder
		return nil, err
placeholder
	return r.toModel(p), nil
placeholder

// Create 创建模板
func (r *tlsFingerprintProfileRepository) Create(ctx context.Context, p *model.TLSFingerprintProfile) (*model.TLSFingerprintProfile, error) {
	builder := r.client.TLSFingerprintProfile.Create().
		SetName(p.Name).
		SetEnableGrease(p.EnableGREASE)

	if p.Description != nil {
		builder.SetDescription(*p.Description)
placeholder
	if len(p.CipherSuites) > 0 {
		builder.SetCipherSuites(p.CipherSuites)
placeholder
	if len(p.Curves) > 0 {
		builder.SetCurves(p.Curves)
placeholder
	if len(p.PointFormats) > 0 {
		builder.SetPointFormats(p.PointFormats)
placeholder
	if len(p.SignatureAlgorithms) > 0 {
		builder.SetSignatureAlgorithms(p.SignatureAlgorithms)
placeholder
	if len(p.ALPNProtocols) > 0 {
		builder.SetAlpnProtocols(p.ALPNProtocols)
placeholder
	if len(p.SupportedVersions) > 0 {
		builder.SetSupportedVersions(p.SupportedVersions)
placeholder
	if len(p.KeyShareGroups) > 0 {
		builder.SetKeyShareGroups(p.KeyShareGroups)
placeholder
	if len(p.PSKModes) > 0 {
		builder.SetPskModes(p.PSKModes)
placeholder
	if len(p.Extensions) > 0 {
		builder.SetExtensions(p.Extensions)
placeholder

	created, err := builder.Save(ctx)
	if err != nil {
		return nil, err
placeholder
	return r.toModel(created), nil
placeholder

// Update 更新模板
func (r *tlsFingerprintProfileRepository) Update(ctx context.Context, p *model.TLSFingerprintProfile) (*model.TLSFingerprintProfile, error) {
	builder := r.client.TLSFingerprintProfile.UpdateOneID(p.ID).
		SetName(p.Name).
		SetEnableGrease(p.EnableGREASE)

	if p.Description != nil {
		builder.SetDescription(*p.Description)
placeholder else {
		builder.ClearDescription()
placeholder

	if len(p.CipherSuites) > 0 {
		builder.SetCipherSuites(p.CipherSuites)
placeholder else {
		builder.ClearCipherSuites()
placeholder
	if len(p.Curves) > 0 {
		builder.SetCurves(p.Curves)
placeholder else {
		builder.ClearCurves()
placeholder
	if len(p.PointFormats) > 0 {
		builder.SetPointFormats(p.PointFormats)
placeholder else {
		builder.ClearPointFormats()
placeholder
	if len(p.SignatureAlgorithms) > 0 {
		builder.SetSignatureAlgorithms(p.SignatureAlgorithms)
placeholder else {
		builder.ClearSignatureAlgorithms()
placeholder
	if len(p.ALPNProtocols) > 0 {
		builder.SetAlpnProtocols(p.ALPNProtocols)
placeholder else {
		builder.ClearAlpnProtocols()
placeholder
	if len(p.SupportedVersions) > 0 {
		builder.SetSupportedVersions(p.SupportedVersions)
placeholder else {
		builder.ClearSupportedVersions()
placeholder
	if len(p.KeyShareGroups) > 0 {
		builder.SetKeyShareGroups(p.KeyShareGroups)
placeholder else {
		builder.ClearKeyShareGroups()
placeholder
	if len(p.PSKModes) > 0 {
		builder.SetPskModes(p.PSKModes)
placeholder else {
		builder.ClearPskModes()
placeholder
	if len(p.Extensions) > 0 {
		builder.SetExtensions(p.Extensions)
placeholder else {
		builder.ClearExtensions()
placeholder

	updated, err := builder.Save(ctx)
	if err != nil {
		return nil, err
placeholder
	return r.toModel(updated), nil
placeholder

// Delete 删除模板
func (r *tlsFingerprintProfileRepository) Delete(ctx context.Context, id int64) error {
	return r.client.TLSFingerprintProfile.DeleteOneID(id).Exec(ctx)
placeholder

// toModel 将 Ent 实体转换为服务模型
func (r *tlsFingerprintProfileRepository) toModel(e *ent.TLSFingerprintProfile) *model.TLSFingerprintProfile {
	p := &model.TLSFingerprintProfile{
		ID:                  e.ID,
		Name:                e.Name,
		Description:         e.Description,
		EnableGREASE:        e.EnableGrease,
		CipherSuites:        e.CipherSuites,
		Curves:              e.Curves,
		PointFormats:        e.PointFormats,
		SignatureAlgorithms: e.SignatureAlgorithms,
		ALPNProtocols:       e.AlpnProtocols,
		SupportedVersions:   e.SupportedVersions,
		KeyShareGroups:      e.KeyShareGroups,
		PSKModes:            e.PskModes,
		Extensions:          e.Extensions,
		CreatedAt:           e.CreatedAt,
		UpdatedAt:           e.UpdatedAt,
placeholder

	// 确保切片不为 nil
	if p.CipherSuites == nil {
		p.CipherSuites = []uint16{placeholder
placeholder
	if p.Curves == nil {
		p.Curves = []uint16{placeholder
placeholder
	if p.PointFormats == nil {
		p.PointFormats = []uint16{placeholder
placeholder
	if p.SignatureAlgorithms == nil {
		p.SignatureAlgorithms = []uint16{placeholder
placeholder
	if p.ALPNProtocols == nil {
		p.ALPNProtocols = []string{placeholder
placeholder
	if p.SupportedVersions == nil {
		p.SupportedVersions = []uint16{placeholder
placeholder
	if p.KeyShareGroups == nil {
		p.KeyShareGroups = []uint16{placeholder
placeholder
	if p.PSKModes == nil {
		p.PSKModes = []uint16{placeholder
placeholder
	if p.Extensions == nil {
		p.Extensions = []uint16{placeholder
placeholder

	return p
placeholder
