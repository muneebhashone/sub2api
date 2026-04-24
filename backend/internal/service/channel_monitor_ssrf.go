package service

import (
	"context"
	"net"
	"strings"
)

// SSRF 防护 helper：
//   - validateEndpoint 在 admin 提交时阻止 http/loopback/私网/云元数据 URL
//   - safeDialContext 在 socket 层再次校验真实 IP，防止 DNS rebinding
//
// 已知 cloud metadata hostname 拒绝列表（小写比较）。
var monitorBlockedHostnames = map[string]struct{placeholder{
	"localhost":                  {placeholder,
	"localhost.localdomain":      {placeholder,
	"metadata":                   {placeholder,
	"metadata.google.internal":   {placeholder,
	"metadata.goog":              {placeholder,
	"instance-data":              {placeholder,
	"instance-data.ec2.internal": {placeholder,
placeholder

// CIDR 列表：包含所有需要拒绝的 IPv4/IPv6 段。
// 解析时只 panic 一次（启动时确认），生产路径只做 Contains。
var monitorBlockedCIDRs = mustParseCIDRs([]string{
	"127.0.0.0/8",    // IPv4 loopback
	"10.0.0.0/8",     // RFC1918
	"172.16.0.0/12",  // RFC1918
	"192.168.0.0/16", // RFC1918
	"169.254.0.0/16", // link-local（含云元数据 169.254.169.254）
	"100.64.0.0/10",  // CGNAT
	"0.0.0.0/8",      // "this network"
	"::1/128",        // IPv6 loopback
	"fc00::/7",       // IPv6 ULA
	"fe80::/10",      // IPv6 link-local
	"::/128",         // IPv6 unspecified
placeholder)

// monitorDialer 共享 Dialer，与 net/http 默认值对齐。
var monitorDialer = &net.Dialer{
	Timeout:   monitorDialTimeout,
	KeepAlive: monitorDialKeepAlive,
placeholder

// mustParseCIDRs 在包初始化时解析 CIDR 字符串，失败 panic。
func mustParseCIDRs(cidrs []string) []*net.IPNet {
	out := make([]*net.IPNet, 0, len(cidrs))
	for _, c := range cidrs {
		_, n, err := net.ParseCIDR(c)
		if err != nil {
			panic("channel_monitor_ssrf: invalid CIDR " + c + ": " + err.Error())
	placeholder
		out = append(out, n)
placeholder
	return out
placeholder

// isBlockedHostname 判断 hostname 是否命中黑名单。
func isBlockedHostname(hostname string) bool {
	if hostname == "" {
		return true
placeholder
	_, blocked := monitorBlockedHostnames[strings.ToLower(hostname)]
	return blocked
placeholder

// isPrivateIP 判断 IP 是否落在禁止段（loopback/RFC1918/link-local/ULA 等）。
func isPrivateIP(ip net.IP) bool {
	if ip == nil {
		return true
placeholder
	if ip.IsUnspecified() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsInterfaceLocalMulticast() {
		return true
placeholder
	for _, n := range monitorBlockedCIDRs {
		if n.Contains(ip) {
			return true
	placeholder
placeholder
	return false
placeholder

// isPrivateOrLoopbackHost 解析 hostname 的所有 A/AAAA 记录，
// 任一 IP 落在私网/loopback 段即认为不安全。
//
// hostname 是 IP 字面量时也走同一路径。
func isPrivateOrLoopbackHost(ctx context.Context, hostname string) (bool, error) {
	if isBlockedHostname(hostname) {
		return true, nil
placeholder
	// IP 字面量直接判断。
	if ip := net.ParseIP(hostname); ip != nil {
		return isPrivateIP(ip), nil
placeholder
	resolver := net.DefaultResolver
	addrs, err := resolver.LookupIPAddr(ctx, hostname)
	if err != nil {
		return false, err
placeholder
	if len(addrs) == 0 {
		return true, nil
placeholder
	for _, a := range addrs {
		if isPrivateIP(a.IP) {
			return true, nil
	placeholder
placeholder
	return false, nil
placeholder

// safeDialContext 在真实 dial 前再次校验目标 IP，防止 DNS rebinding。
// 解析 hostname 后逐个 IP 尝试连接，命中私网即拒绝（即便 validateEndpoint 时返回的是公网 IP）。
func safeDialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
placeholder
	// 字面量 IP 走快速路径。
	if ip := net.ParseIP(host); ip != nil {
		if isPrivateIP(ip) {
			return nil, &net.AddrError{Err: "blocked by SSRF policy", Addr: addressplaceholder
	placeholder
		return monitorDialer.DialContext(ctx, network, address)
placeholder
	if isBlockedHostname(host) {
		return nil, &net.AddrError{Err: "blocked by SSRF policy", Addr: addressplaceholder
placeholder
	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
placeholder
	if len(addrs) == 0 {
		return nil, &net.AddrError{Err: "no addresses for host", Addr: hostplaceholder
placeholder
	var lastErr error
	for _, a := range addrs {
		if isPrivateIP(a.IP) {
			lastErr = &net.AddrError{Err: "blocked by SSRF policy", Addr: a.IP.String()placeholder
			continue
	placeholder
		conn, err := monitorDialer.DialContext(ctx, network, net.JoinHostPort(a.IP.String(), port))
		if err == nil {
			return conn, nil
	placeholder
		lastErr = err
placeholder
	if lastErr == nil {
		lastErr = &net.AddrError{Err: "no usable addresses", Addr: hostplaceholder
placeholder
	return nil, lastErr
placeholder
