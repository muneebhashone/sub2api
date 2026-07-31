//go:build unit

package service

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// newSMTPTestCert 生成 127.0.0.1/localhost 的自签证书及其信任池。
func newSMTPTestCert(t *testing.T) (tls.Certificate, *x509.CertPool) {
placeholder
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
placeholder
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "127.0.0.1"placeholder,
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuthplaceholder,
		BasicConstraintsValid: true,
		IsCA:                  true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")placeholder,
		DNSNames:              []string{"localhost"placeholder,
placeholder
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
placeholder
	leaf, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("parse certificate: %v", err)
placeholder
	pool := x509.NewCertPool()
	pool.AddCert(leaf)
	return tls.Certificate{Certificate: [][]byte{derplaceholder, PrivateKey: privplaceholder, pool
placeholder

// fakeSMTPServer 是覆盖三种连接形态的最小 SMTP 服务器：
// 隐式 TLS（465 语义）、明文+STARTTLS（587 语义）、纯明文。
type fakeSMTPServer struct {
	listener          net.Listener
	tlsConfig         *tls.Config
	advertiseStartTLS bool

	mu       sync.Mutex
	commands []string
	conns    atomic.Int64
	wg       sync.WaitGroup
placeholder

func startFakeSMTPServer(t *testing.T, implicitTLS, advertiseStartTLS bool) (*fakeSMTPServer, int) {
placeholder
	cert, pool := newSMTPTestCert(t)
	prevPool := smtpTestRootCAs
	smtpTestRootCAs = pool
	t.Cleanup(func() { smtpTestRootCAs = prevPool placeholder)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
placeholder
	srv := &fakeSMTPServer{
		listener:          listener,
		tlsConfig:         &tls.Config{Certificates: []tls.Certificate{certplaceholder, MinVersion: tls.VersionTLS12placeholder,
		advertiseStartTLS: advertiseStartTLS,
placeholder
	if implicitTLS {
		srv.listener = tls.NewListener(listener, srv.tlsConfig)
placeholder
	t.Cleanup(func() {
		_ = srv.listener.Close()
		srv.wg.Wait()
placeholder)

	srv.wg.Add(1)
	go func() {
		defer srv.wg.Done()
		for {
			conn, err := srv.listener.Accept()
			if err != nil {
				return
		placeholder
			srv.conns.Add(1)
			srv.wg.Add(1)
			go func() {
				defer srv.wg.Done()
				defer func() { _ = conn.Close() placeholder()
				_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
				srv.serve(conn, srv.advertiseStartTLS)
		placeholder()
	placeholder
placeholder()

	port := listener.Addr().(*net.TCPAddr).Port
	return srv, port
placeholder

func (srv *fakeSMTPServer) record(cmd string) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.commands = append(srv.commands, cmd)
placeholder

func (srv *fakeSMTPServer) sawCommand(prefix string) bool {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	for _, cmd := range srv.commands {
		if strings.HasPrefix(strings.ToUpper(cmd), prefix) {
			return true
	placeholder
placeholder
	return false
placeholder

func (srv *fakeSMTPServer) serve(conn net.Conn, allowStartTLS bool) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	writeLine := func(line string) bool {
		if _, err := writer.WriteString(line + "\r\n"); err != nil {
			return false
	placeholder
		return writer.Flush() == nil
placeholder
	if !writeLine("220 fake.test ESMTP ready") {
		return
placeholder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
	placeholder
		cmd := strings.TrimSpace(line)
		srv.record(cmd)
		upper := strings.ToUpper(cmd)
		switch {
		case strings.HasPrefix(upper, "EHLO"), strings.HasPrefix(upper, "HELO"):
			ok := writeLine("250-fake.test")
			if allowStartTLS {
				ok = ok && writeLine("250-STARTTLS")
		placeholder
			if !(ok && writeLine("250-AUTH PLAIN LOGIN") && writeLine("250 8BITMIME")) {
				return
		placeholder
		case upper == "STARTTLS" && allowStartTLS:
			if !writeLine("220 2.0.0 ready to start TLS") {
				return
		placeholder
			tlsConn := tls.Server(conn, srv.tlsConfig)
			if err := tlsConn.Handshake(); err != nil {
				return
		placeholder
			srv.serveUpgraded(tlsConn)
			return
		case strings.HasPrefix(upper, "AUTH"):
			if !writeLine("235 2.7.0 authentication successful") {
				return
		placeholder
		case strings.HasPrefix(upper, "MAIL"), strings.HasPrefix(upper, "RCPT"):
			if !writeLine("250 ok") {
				return
		placeholder
		case upper == "DATA":
			if !writeLine("354 go ahead") {
				return
		placeholder
			for {
				dataLine, err := reader.ReadString('\n')
				if err != nil {
					return
			placeholder
				if strings.TrimRight(dataLine, "\r\n") == "." {
					break
			placeholder
		placeholder
			if !writeLine("250 message accepted") {
				return
		placeholder
		case upper == "QUIT":
			_ = writeLine("221 bye")
			return
		default:
			if !writeLine("250 ok") {
				return
		placeholder
	placeholder
placeholder
placeholder

// serveUpgraded 复用命令循环处理 STARTTLS 升级后的会话（升级后不再提供 STARTTLS）。
func (srv *fakeSMTPServer) serveUpgraded(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	// net/smtp 在 StartTLS 成功后会重新发送 EHLO，直接进入命令循环即可。
	srv.serveCommands(reader, writer)
placeholder

func (srv *fakeSMTPServer) serveCommands(reader *bufio.Reader, writer *bufio.Writer) {
	writeLine := func(line string) bool {
		if _, err := writer.WriteString(line + "\r\n"); err != nil {
			return false
	placeholder
		return writer.Flush() == nil
placeholder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
	placeholder
		cmd := strings.TrimSpace(line)
		srv.record(cmd)
		upper := strings.ToUpper(cmd)
		switch {
		case strings.HasPrefix(upper, "EHLO"), strings.HasPrefix(upper, "HELO"):
			if !(writeLine("250-fake.test") && writeLine("250-AUTH PLAIN LOGIN") && writeLine("250 8BITMIME")) {
				return
		placeholder
		case strings.HasPrefix(upper, "AUTH"):
			if !writeLine("235 2.7.0 authentication successful") {
				return
		placeholder
		case strings.HasPrefix(upper, "MAIL"), strings.HasPrefix(upper, "RCPT"):
			if !writeLine("250 ok") {
				return
		placeholder
		case upper == "DATA":
			if !writeLine("354 go ahead") {
				return
		placeholder
			for {
				dataLine, err := reader.ReadString('\n')
				if err != nil {
					return
			placeholder
				if strings.TrimRight(dataLine, "\r\n") == "." {
					break
			placeholder
		placeholder
			if !writeLine("250 message accepted") {
				return
		placeholder
		case upper == "QUIT":
			_ = writeLine("221 bye")
			return
		default:
			if !writeLine("250 ok") {
				return
		placeholder
	placeholder
placeholder
placeholder

func smtpTestConfig(port int, useTLS bool) *SMTPConfig {
	return &SMTPConfig{
		Host:     "127.0.0.1",
		Port:     port,
		Username: "user",
		Password: "pass",
		From:     "noreply@example.com",
		FromName: "Test",
		UseTLS:   useTLS,
placeholder
placeholder

// 465 语义：UseTLS=true + 隐式 TLS 服务器，原有路径保持可用。
func TestSMTPConnectionImplicitTLS(t *testing.T) {
	srv, port := startFakeSMTPServer(t, true, false)
	svc := &EmailService{placeholder

	if err := svc.TestSMTPConnectionWithConfig(smtpTestConfig(port, true)); err != nil {
		t.Fatalf("expected implicit TLS connection to succeed, got: %v", err)
placeholder
	if !srv.sawCommand("EHLO") {
		t.Fatal("expected server to receive EHLO")
placeholder
placeholder

// 587 语义（#1470/#1488 核心场景）：UseTLS=true + 明文问候的 STARTTLS 服务器，
// 隐式 TLS 失败后必须自动降级为强制 STARTTLS 并成功。
func TestSMTPConnectionStartTLSFallbackWhenTLSEnabled(t *testing.T) {
	srv, port := startFakeSMTPServer(t, false, true)
	svc := &EmailService{placeholder

	if err := svc.TestSMTPConnectionWithConfig(smtpTestConfig(port, true)); err != nil {
		t.Fatalf("expected STARTTLS fallback to succeed, got: %v", err)
placeholder
	if !srv.sawCommand("STARTTLS") {
		t.Fatal("expected server to receive STARTTLS command")
placeholder
	if got := srv.conns.Load(); got < 2 {
		t.Fatalf("expected implicit TLS attempt before STARTTLS fallback (>=2 connections), got %d", got)
placeholder
placeholder

// UseTLS=true 但服务器不支持 STARTTLS：必须报错，且绝不能把凭据发到明文连接上。
func TestSMTPConnectionMandatoryStartTLSRefusesPlaintext(t *testing.T) {
	srv, port := startFakeSMTPServer(t, false, false)
	svc := &EmailService{placeholder

	err := svc.TestSMTPConnectionWithConfig(smtpTestConfig(port, true))
	if err == nil {
		t.Fatal("expected error when server does not support STARTTLS")
placeholder
	if !strings.Contains(err.Error(), "STARTTLS") {
		t.Fatalf("expected STARTTLS-related error, got: %v", err)
placeholder
	if srv.sawCommand("AUTH") {
		t.Fatal("credentials must not be sent over plaintext when TLS is required")
placeholder
placeholder

// UseTLS=false + 服务器支持 STARTTLS：测试连接与发送路径一致，机会式升级后认证成功。
// 这是 #1488 评论"测试连接不成功，发送测试邮件实际上能发"的回归用例。
func TestSMTPConnectionOpportunisticStartTLSWhenTLSDisabled(t *testing.T) {
	srv, port := startFakeSMTPServer(t, false, true)
	svc := &EmailService{placeholder

	if err := svc.TestSMTPConnectionWithConfig(smtpTestConfig(port, false)); err != nil {
		t.Fatalf("expected opportunistic STARTTLS test connection to succeed, got: %v", err)
placeholder
	if !srv.sawCommand("STARTTLS") {
		t.Fatal("expected test connection to upgrade via STARTTLS like the send path")
placeholder
placeholder

// UseTLS=false + 服务器不支持 STARTTLS：保持明文直连（既有行为不回归）。
func TestSMTPConnectionPlainWhenNoStartTLS(t *testing.T) {
	srv, port := startFakeSMTPServer(t, false, false)
	svc := &EmailService{placeholder

	if err := svc.TestSMTPConnectionWithConfig(smtpTestConfig(port, false)); err != nil {
		t.Fatalf("expected plain connection to succeed, got: %v", err)
placeholder
	if srv.sawCommand("STARTTLS") {
		t.Fatal("did not expect STARTTLS command when server does not advertise it")
placeholder
placeholder

// 发送路径全流程：UseTLS=true + STARTTLS 服务器（587 语义）完整走完 MAIL/RCPT/DATA。
func TestSendEmailWithConfigStartTLSFallback(t *testing.T) {
	srv, port := startFakeSMTPServer(t, false, true)
	svc := &EmailService{placeholder

	err := svc.SendEmailWithConfig(smtpTestConfig(port, true), "rcpt@example.com", "subject", "<p>body</p>")
	if err != nil {
		t.Fatalf("expected send via STARTTLS fallback to succeed, got: %v", err)
placeholder
	if !srv.sawCommand("STARTTLS") {
		t.Fatal("expected send path to upgrade via STARTTLS")
placeholder
	if !srv.sawCommand("DATA") {
		t.Fatal("expected send path to reach DATA")
placeholder
placeholder

// 发送路径全流程：UseTLS=true + 隐式 TLS 服务器（465 语义）保持既有行为。
func TestSendEmailWithConfigImplicitTLS(t *testing.T) {
	srv, port := startFakeSMTPServer(t, true, false)
	svc := &EmailService{placeholder

	err := svc.SendEmailWithConfig(smtpTestConfig(port, true), "rcpt@example.com", "subject", "<p>body</p>")
	if err != nil {
		t.Fatalf("expected send via implicit TLS to succeed, got: %v", err)
placeholder
	if !srv.sawCommand("DATA") {
		t.Fatal("expected send path to reach DATA")
placeholder
placeholder
