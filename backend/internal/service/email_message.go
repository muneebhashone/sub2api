package service

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net/mail"
	"strings"
	"time"
)

type smtpMessage struct {
	envelopeFrom string
	envelopeTo   string
	data         []byte
placeholder

func buildSMTPMessage(config *SMTPConfig, to, subject, body string) (smtpMessage, error) {
	if config == nil {
		return smtpMessage{placeholder, errors.New("missing SMTP configuration")
placeholder

	fromAddress, err := parseSMTPAddress(config.From, "from")
	if err != nil {
		return smtpMessage{placeholder, err
placeholder
	recipientAddress, err := parseSMTPAddress(to, "recipient")
	if err != nil {
		return smtpMessage{placeholder, err
placeholder
	messageID, err := generateEmailMessageID(fromAddress.Address, config.Host)
	if err != nil {
		return smtpMessage{placeholder, fmt.Errorf("generate message ID: %w", err)
placeholder

	fromName := sanitizeEmailHeader(config.FromName)
	if strings.TrimSpace(fromName) == "" {
		fromName = fromAddress.Name
placeholder
	fromHeader := (&mail.Address{
		Name:    fromName,
		Address: fromAddress.Address,
placeholder).String()
	toHeader := (&mail.Address{
		Name:    recipientAddress.Name,
		Address: recipientAddress.Address,
placeholder).String()
	subjectHeader := mime.QEncoding.Encode("UTF-8", sanitizeEmailHeader(subject))

	var message bytes.Buffer
	fmt.Fprintf(&message, "From: %s\r\n", fromHeader)
	fmt.Fprintf(&message, "To: %s\r\n", toHeader)
	fmt.Fprintf(&message, "Date: %s\r\n", time.Now().UTC().Format(time.RFC1123Z))
	fmt.Fprintf(&message, "Message-ID: %s\r\n", messageID)
	fmt.Fprintf(&message, "Subject: %s\r\n", subjectHeader)
	fmt.Fprint(&message, "MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"Content-Transfer-Encoding: quoted-printable\r\n\r\n")

	bodyWriter := quotedprintable.NewWriter(&message)
	if _, err := bodyWriter.Write([]byte(body)); err != nil {
		return smtpMessage{placeholder, fmt.Errorf("encode email body: %w", err)
placeholder
	if err := bodyWriter.Close(); err != nil {
		return smtpMessage{placeholder, fmt.Errorf("close email body encoder: %w", err)
placeholder

	return smtpMessage{
		envelopeFrom: fromAddress.Address,
		envelopeTo:   recipientAddress.Address,
		data:         message.Bytes(),
placeholder, nil
placeholder

func parseSMTPAddress(value, field string) (*mail.Address, error) {
	if strings.ContainsAny(value, "\r\n") {
		return nil, fmt.Errorf("invalid SMTP %s address: contains a line break", field)
placeholder

	cleaned := strings.TrimSpace(value)
	address, err := mail.ParseAddress(cleaned)
	if err != nil || strings.TrimSpace(address.Address) == "" {
		if err == nil {
			err = fmt.Errorf("address is empty")
	placeholder
		return nil, fmt.Errorf("invalid SMTP %s address: %w", field, err)
placeholder
	return address, nil
placeholder

func generateEmailMessageID(fromAddress, smtpHost string) (string, error) {
	randomID := make([]byte, 16)
	if _, err := rand.Read(randomID); err != nil {
		return "", err
placeholder

	domain := strings.TrimSpace(sanitizeEmailHeader(smtpHost))
	if at := strings.LastIndexByte(fromAddress, '@'); at >= 0 && at < len(fromAddress)-1 {
		domain = fromAddress[at+1:]
placeholder
	domain = strings.Trim(domain, "[]<>")
	if domain == "" {
		domain = "localhost"
placeholder

	return fmt.Sprintf("<%s@%s>", hex.EncodeToString(randomID), domain), nil
placeholder
