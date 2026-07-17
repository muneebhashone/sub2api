package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

const ingressRejectUpsertChunkSize = 500

func (r *opsRepository) BatchUpsertIngressRejects(ctx context.Context, items []*service.OpsIngressRejectAggregate) error {
	if r == nil || r.db == nil || len(items) == 0 {
		return nil
placeholder
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
placeholder
	defer func() { _ = tx.Rollback() placeholder()

	for start := 0; start < len(items); start += ingressRejectUpsertChunkSize {
		end := start + ingressRejectUpsertChunkSize
		if end > len(items) {
			end = len(items)
	placeholder
		valid := make([]*service.OpsIngressRejectAggregate, 0, end-start)
		for _, item := range items[start:end] {
			if item != nil && item.RequestCount > 0 {
				valid = append(valid, item)
		placeholder
	placeholder
		if len(valid) == 0 {
			continue
	placeholder

		var query strings.Builder
		query.WriteString(`INSERT INTO ops_ingress_reject_aggregates
  (bucket_start, reject_reason, route_family, protocol, client_ip, user_id, api_key_id, request_count, first_seen, last_seen)
VALUES `)
		args := make([]any, 0, len(valid)*10)
		for i, item := range valid {
			if i > 0 {
				query.WriteByte(',')
		placeholder
			base := len(args)
			fmt.Fprintf(&query, "($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
				base+1, base+2, base+3, base+4, base+5, base+6, base+7, base+8, base+9, base+10)
			var userID, apiKeyID int64
			if item.UserID != nil {
				userID = *item.UserID
		placeholder
			if item.APIKeyID != nil {
				apiKeyID = *item.APIKeyID
		placeholder
			args = append(args, item.BucketStart.UTC(), item.RejectReason, item.RouteFamily, item.Protocol,
				item.ClientIP, userID, apiKeyID, item.RequestCount, item.FirstSeen.UTC(), item.LastSeen.UTC())
	placeholder
		query.WriteString(`
ON CONFLICT (bucket_start, reject_reason, route_family, protocol, client_ip, user_id, api_key_id)
DO UPDATE SET request_count = ops_ingress_reject_aggregates.request_count + EXCLUDED.request_count,
              first_seen = LEAST(ops_ingress_reject_aggregates.first_seen, EXCLUDED.first_seen),
              last_seen = GREATEST(ops_ingress_reject_aggregates.last_seen, EXCLUDED.last_seen),
              updated_at = NOW()`)
		if _, err := tx.ExecContext(ctx, query.String(), args...); err != nil {
			return err
	placeholder
placeholder
	return tx.Commit()
placeholder

func (r *opsRepository) ListIngressRejects(ctx context.Context, filter *service.OpsIngressRejectFilter) (*service.OpsIngressRejectList, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("nil ops repository")
placeholder
	if filter == nil {
		filter = &service.OpsIngressRejectFilter{placeholder
placeholder
	page, pageSize := filter.Page, filter.PageSize
	if page <= 0 {
		page = 1
placeholder
	if pageSize <= 0 {
		pageSize = 50
placeholder
	if pageSize > 200 {
		pageSize = 200
placeholder

	clauses := []string{"1=1"placeholder
	args := make([]any, 0)
	add := func(expr string, value any) {
		args = append(args, value)
		clauses = append(clauses, fmt.Sprintf(expr, len(args)))
placeholder
	if filter.StartTime != nil {
		add("bucket_start >= $%d", filter.StartTime.UTC())
placeholder
	if filter.EndTime != nil {
		add("bucket_start < $%d", filter.EndTime.UTC())
placeholder
	if value := strings.TrimSpace(filter.RejectReason); value != "" {
		add("reject_reason = $%d", value)
placeholder
	if value := strings.TrimSpace(filter.RouteFamily); value != "" {
		add("route_family = $%d", value)
placeholder
	if value := strings.TrimSpace(filter.Protocol); value != "" {
		add("protocol = $%d", value)
placeholder
	if value := strings.TrimSpace(filter.ClientIP); value != "" {
		add("client_ip = $%d::inet", value)
placeholder
	if filter.UserID != nil {
		add("user_id = $%d", *filter.UserID)
placeholder
	if filter.APIKeyID != nil {
		add("api_key_id = $%d", *filter.APIKeyID)
placeholder
	where := "WHERE " + strings.Join(clauses, " AND ")

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM ops_ingress_reject_aggregates "+where, args...).Scan(&total); err != nil {
		return nil, err
placeholder
	args = append(args, pageSize, (page-1)*pageSize)
	query := fmt.Sprintf(`SELECT id,bucket_start,reject_reason,route_family,protocol,host(client_ip),user_id,api_key_id,request_count,first_seen,last_seen
FROM ops_ingress_reject_aggregates %s ORDER BY bucket_start DESC,id DESC LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args))
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
placeholder
	defer rows.Close()

	result := &service.OpsIngressRejectList{
		Items: make([]*service.OpsIngressRejectAggregate, 0, pageSize), Total: total, Page: page, PageSize: pageSize,
placeholder
	for rows.Next() {
		item := &service.OpsIngressRejectAggregate{placeholder
		var userID, apiKeyID int64
		if err := rows.Scan(&item.ID, &item.BucketStart, &item.RejectReason, &item.RouteFamily, &item.Protocol,
			&item.ClientIP, &userID, &apiKeyID, &item.RequestCount, &item.FirstSeen, &item.LastSeen); err != nil {
			return nil, err
	placeholder
		if userID > 0 {
			item.UserID = &userID
	placeholder
		if apiKeyID > 0 {
			item.APIKeyID = &apiKeyID
	placeholder
		result.Items = append(result.Items, item)
placeholder
	if err := rows.Err(); err != nil {
		return nil, err
placeholder
	return result, nil
placeholder
