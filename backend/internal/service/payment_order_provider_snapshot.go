package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/payment"
)

type paymentOrderProviderSnapshot struct {
	SchemaVersion      int
	ProviderInstanceID string
	ProviderKey        string
	PaymentMode        string
placeholder

func psOrderProviderSnapshot(order *dbent.PaymentOrder) *paymentOrderProviderSnapshot {
	if order == nil || len(order.ProviderSnapshot) == 0 {
		return nil
placeholder

	snapshot := &paymentOrderProviderSnapshot{
		SchemaVersion:      psSnapshotIntValue(order.ProviderSnapshot["schema_version"]),
		ProviderInstanceID: psSnapshotStringValue(order.ProviderSnapshot["provider_instance_id"]),
		ProviderKey:        psSnapshotStringValue(order.ProviderSnapshot["provider_key"]),
		PaymentMode:        psSnapshotStringValue(order.ProviderSnapshot["payment_mode"]),
placeholder
	if snapshot.SchemaVersion == 0 && snapshot.ProviderInstanceID == "" && snapshot.ProviderKey == "" && snapshot.PaymentMode == "" {
		return nil
placeholder
	return snapshot
placeholder

func psSnapshotStringValue(value any) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	default:
		return ""
placeholder
placeholder

func psSnapshotIntValue(value any) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int32:
		return int(typed)
	case int64:
		return int(typed)
	case float32:
		return int(typed)
	case float64:
		return int(typed)
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(typed))
		if err == nil {
			return n
	placeholder
placeholder
	return 0
placeholder

func (s *PaymentService) resolveSnapshotOrderProviderInstance(ctx context.Context, order *dbent.PaymentOrder, snapshot *paymentOrderProviderSnapshot) (*dbent.PaymentProviderInstance, error) {
	if s == nil || s.entClient == nil || order == nil || snapshot == nil {
		return nil, nil
placeholder

	snapshotInstanceID := strings.TrimSpace(snapshot.ProviderInstanceID)
	columnInstanceID := strings.TrimSpace(psStringValue(order.ProviderInstanceID))
	if snapshotInstanceID == "" {
		snapshotInstanceID = columnInstanceID
placeholder
	if snapshotInstanceID == "" {
		return nil, fmt.Errorf("order %d provider snapshot is missing provider_instance_id", order.ID)
placeholder
	if columnInstanceID != "" && snapshot.ProviderInstanceID != "" && !strings.EqualFold(columnInstanceID, snapshot.ProviderInstanceID) {
		return nil, fmt.Errorf("order %d provider snapshot instance mismatch: snapshot=%s order=%s", order.ID, snapshot.ProviderInstanceID, columnInstanceID)
placeholder

	instID, err := strconv.ParseInt(snapshotInstanceID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("order %d provider snapshot instance id is invalid: %s", order.ID, snapshotInstanceID)
placeholder

	inst, err := s.entClient.PaymentProviderInstance.Get(ctx, instID)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, fmt.Errorf("order %d provider snapshot instance %s is missing", order.ID, snapshotInstanceID)
	placeholder
		return nil, err
placeholder

	if snapshot.ProviderKey != "" && !strings.EqualFold(strings.TrimSpace(inst.ProviderKey), snapshot.ProviderKey) {
		return nil, fmt.Errorf("order %d provider snapshot key mismatch: snapshot=%s instance=%s", order.ID, snapshot.ProviderKey, inst.ProviderKey)
placeholder

	return inst, nil
placeholder

func expectedNotificationProviderKeyForOrder(registry *payment.Registry, order *dbent.PaymentOrder, instanceProviderKey string) string {
	if order == nil {
		return strings.TrimSpace(instanceProviderKey)
placeholder

	orderProviderKey := psStringValue(order.ProviderKey)
	if snapshot := psOrderProviderSnapshot(order); snapshot != nil && snapshot.ProviderKey != "" {
		orderProviderKey = snapshot.ProviderKey
placeholder

	return expectedNotificationProviderKey(registry, order.PaymentType, orderProviderKey, instanceProviderKey)
placeholder
