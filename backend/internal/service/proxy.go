package service

import (
	"fmt"
	"time"
)

type Proxy struct {
	ID        int64
	Name      string
	Protocol  string
	Host      string
	Port      int
	Username  string
	Password  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
placeholder

func (p *Proxy) IsActive() bool {
	return p.Status == StatusActive
placeholder

func (p *Proxy) URL() string {
	if p.Username != "" && p.Password != "" {
		return fmt.Sprintf("%s://%s:%s@%s:%d", p.Protocol, p.Username, p.Password, p.Host, p.Port)
placeholder
	return fmt.Sprintf("%s://%s:%d", p.Protocol, p.Host, p.Port)
placeholder

type ProxyWithAccountCount struct {
	Proxy
	AccountCount   int64
	LatencyMs      *int64
	LatencyStatus  string
	LatencyMessage string
placeholder

type ProxyAccountSummary struct {
	ID       int64
	Name     string
	Platform string
	Type     string
	Notes    *string
placeholder
