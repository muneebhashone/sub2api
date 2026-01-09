package repository

import (
	"context"
	"database/sql"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/proxy"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

type sqlQuerier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
placeholder

type proxyRepository struct {
	client *dbent.Client
	sql    sqlQuerier
placeholder

func NewProxyRepository(client *dbent.Client, sqlDB *sql.DB) service.ProxyRepository {
	return newProxyRepositoryWithSQL(client, sqlDB)
placeholder

func newProxyRepositoryWithSQL(client *dbent.Client, sqlq sqlQuerier) *proxyRepository {
	return &proxyRepository{client: client, sql: sqlqplaceholder
placeholder

func (r *proxyRepository) Create(ctx context.Context, proxyIn *service.Proxy) error {
	builder := r.client.Proxy.Create().
		SetName(proxyIn.Name).
		SetProtocol(proxyIn.Protocol).
		SetHost(proxyIn.Host).
		SetPort(proxyIn.Port).
		SetStatus(proxyIn.Status)
	if proxyIn.Username != "" {
		builder.SetUsername(proxyIn.Username)
placeholder
	if proxyIn.Password != "" {
		builder.SetPassword(proxyIn.Password)
placeholder

	created, err := builder.Save(ctx)
	if err == nil {
		applyProxyEntityToService(proxyIn, created)
placeholder
	return err
placeholder

func (r *proxyRepository) GetByID(ctx context.Context, id int64) (*service.Proxy, error) {
	m, err := r.client.Proxy.Get(ctx, id)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrProxyNotFound
	placeholder
		return nil, err
placeholder
	return proxyEntityToService(m), nil
placeholder

func (r *proxyRepository) Update(ctx context.Context, proxyIn *service.Proxy) error {
	builder := r.client.Proxy.UpdateOneID(proxyIn.ID).
		SetName(proxyIn.Name).
		SetProtocol(proxyIn.Protocol).
		SetHost(proxyIn.Host).
		SetPort(proxyIn.Port).
		SetStatus(proxyIn.Status)
	if proxyIn.Username != "" {
		builder.SetUsername(proxyIn.Username)
placeholder else {
		builder.ClearUsername()
placeholder
	if proxyIn.Password != "" {
		builder.SetPassword(proxyIn.Password)
placeholder else {
		builder.ClearPassword()
placeholder

	updated, err := builder.Save(ctx)
	if err == nil {
		applyProxyEntityToService(proxyIn, updated)
		return nil
placeholder
	if dbent.IsNotFound(err) {
		return service.ErrProxyNotFound
placeholder
	return err
placeholder

func (r *proxyRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.client.Proxy.Delete().Where(proxy.IDEQ(id)).Exec(ctx)
	return err
placeholder

func (r *proxyRepository) List(ctx context.Context, params pagination.PaginationParams) ([]service.Proxy, *pagination.PaginationResult, error) {
	return r.ListWithFilters(ctx, params, "", "", "")
placeholder

// ListWithFilters lists proxies with optional filtering by protocol, status, and search query
func (r *proxyRepository) ListWithFilters(ctx context.Context, params pagination.PaginationParams, protocol, status, search string) ([]service.Proxy, *pagination.PaginationResult, error) {
	q := r.client.Proxy.Query()
	if protocol != "" {
		q = q.Where(proxy.ProtocolEQ(protocol))
placeholder
	if status != "" {
		q = q.Where(proxy.StatusEQ(status))
placeholder
	if search != "" {
		q = q.Where(proxy.NameContainsFold(search))
placeholder

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	proxies, err := q.
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order(dbent.Desc(proxy.FieldID)).
		All(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	outProxies := make([]service.Proxy, 0, len(proxies))
	for i := range proxies {
		outProxies = append(outProxies, *proxyEntityToService(proxies[i]))
placeholder

	return outProxies, paginationResultFromTotal(int64(total), params), nil
placeholder

// ListWithFiltersAndAccountCount lists proxies with filters and includes account count per proxy
func (r *proxyRepository) ListWithFiltersAndAccountCount(ctx context.Context, params pagination.PaginationParams, protocol, status, search string) ([]service.ProxyWithAccountCount, *pagination.PaginationResult, error) {
	q := r.client.Proxy.Query()
	if protocol != "" {
		q = q.Where(proxy.ProtocolEQ(protocol))
placeholder
	if status != "" {
		q = q.Where(proxy.StatusEQ(status))
placeholder
	if search != "" {
		q = q.Where(proxy.NameContainsFold(search))
placeholder

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	proxies, err := q.
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order(dbent.Desc(proxy.FieldID)).
		All(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	// Get account counts
	counts, err := r.GetAccountCountsForProxies(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	// Build result with account counts
	result := make([]service.ProxyWithAccountCount, 0, len(proxies))
	for i := range proxies {
		proxyOut := proxyEntityToService(proxies[i])
		if proxyOut == nil {
			continue
	placeholder
		result = append(result, service.ProxyWithAccountCount{
			Proxy:        *proxyOut,
			AccountCount: counts[proxyOut.ID],
	placeholder)
placeholder

	return result, paginationResultFromTotal(int64(total), params), nil
placeholder

func (r *proxyRepository) ListActive(ctx context.Context) ([]service.Proxy, error) {
	proxies, err := r.client.Proxy.Query().
		Where(proxy.StatusEQ(service.StatusActive)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder
	outProxies := make([]service.Proxy, 0, len(proxies))
	for i := range proxies {
		outProxies = append(outProxies, *proxyEntityToService(proxies[i]))
placeholder
	return outProxies, nil
placeholder

// ExistsByHostPortAuth checks if a proxy with the same host, port, username, and password exists
func (r *proxyRepository) ExistsByHostPortAuth(ctx context.Context, host string, port int, username, password string) (bool, error) {
	q := r.client.Proxy.Query().
		Where(proxy.HostEQ(host), proxy.PortEQ(port))

	if username == "" {
		q = q.Where(proxy.Or(proxy.UsernameIsNil(), proxy.UsernameEQ("")))
placeholder else {
		q = q.Where(proxy.UsernameEQ(username))
placeholder
	if password == "" {
		q = q.Where(proxy.Or(proxy.PasswordIsNil(), proxy.PasswordEQ("")))
placeholder else {
		q = q.Where(proxy.PasswordEQ(password))
placeholder

	count, err := q.Count(ctx)
	return count > 0, err
placeholder

// CountAccountsByProxyID returns the number of accounts using a specific proxy
func (r *proxyRepository) CountAccountsByProxyID(ctx context.Context, proxyID int64) (int64, error) {
	var count int64
	if err := scanSingleRow(ctx, r.sql, "SELECT COUNT(*) FROM accounts WHERE proxy_id = $1", []any{proxyIDplaceholder, &count); err != nil {
		return 0, err
placeholder
	return count, nil
placeholder

// GetAccountCountsForProxies returns a map of proxy ID to account count for all proxies
func (r *proxyRepository) GetAccountCountsForProxies(ctx context.Context) (counts map[int64]int64, err error) {
	rows, err := r.sql.QueryContext(ctx, "SELECT proxy_id, COUNT(*) AS count FROM accounts WHERE proxy_id IS NOT NULL AND deleted_at IS NULL GROUP BY proxy_id")
	if err != nil {
		return nil, err
placeholder
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = closeErr
			counts = nil
	placeholder
placeholder()

	counts = make(map[int64]int64)
	for rows.Next() {
		var proxyID, count int64
		if err = rows.Scan(&proxyID, &count); err != nil {
			return nil, err
	placeholder
		counts[proxyID] = count
placeholder
	if err = rows.Err(); err != nil {
		return nil, err
placeholder
	return counts, nil
placeholder

// ListActiveWithAccountCount returns all active proxies with account count, sorted by creation time descending
func (r *proxyRepository) ListActiveWithAccountCount(ctx context.Context) ([]service.ProxyWithAccountCount, error) {
	proxies, err := r.client.Proxy.Query().
		Where(proxy.StatusEQ(service.StatusActive)).
		Order(dbent.Desc(proxy.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder

	// Get account counts
	counts, err := r.GetAccountCountsForProxies(ctx)
	if err != nil {
		return nil, err
placeholder

	// Build result with account counts
	result := make([]service.ProxyWithAccountCount, 0, len(proxies))
	for i := range proxies {
		proxyOut := proxyEntityToService(proxies[i])
		if proxyOut == nil {
			continue
	placeholder
		result = append(result, service.ProxyWithAccountCount{
			Proxy:        *proxyOut,
			AccountCount: counts[proxyOut.ID],
	placeholder)
placeholder

	return result, nil
placeholder

func proxyEntityToService(m *dbent.Proxy) *service.Proxy {
	if m == nil {
		return nil
placeholder
	out := &service.Proxy{
		ID:        m.ID,
		Name:      m.Name,
		Protocol:  m.Protocol,
		Host:      m.Host,
		Port:      m.Port,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
placeholder
	if m.Username != nil {
		out.Username = *m.Username
placeholder
	if m.Password != nil {
		out.Password = *m.Password
placeholder
	return out
placeholder

func applyProxyEntityToService(dst *service.Proxy, src *dbent.Proxy) {
	if dst == nil || src == nil {
		return
placeholder
	dst.ID = src.ID
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
placeholder
