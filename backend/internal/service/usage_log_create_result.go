package service

import "errors"

type usageLogCreateDisposition int

const (
	usageLogCreateDispositionUnknown usageLogCreateDisposition = iota
	usageLogCreateDispositionNotPersisted
	usageLogCreateDispositionDropped
)

type UsageLogCreateError struct {
	err         error
	disposition usageLogCreateDisposition
placeholder

func (e *UsageLogCreateError) Error() string {
	if e == nil || e.err == nil {
		return "usage log create error"
placeholder
	return e.err.Error()
placeholder

func (e *UsageLogCreateError) Unwrap() error {
	if e == nil {
		return nil
placeholder
	return e.err
placeholder

func MarkUsageLogCreateNotPersisted(err error) error {
	if err == nil {
		return nil
placeholder
	return &UsageLogCreateError{
		err:         err,
		disposition: usageLogCreateDispositionNotPersisted,
placeholder
placeholder

func MarkUsageLogCreateDropped(err error) error {
	if err == nil {
		return nil
placeholder
	return &UsageLogCreateError{
		err:         err,
		disposition: usageLogCreateDispositionDropped,
placeholder
placeholder

func IsUsageLogCreateNotPersisted(err error) bool {
	if err == nil {
		return false
placeholder
	var target *UsageLogCreateError
	if !errors.As(err, &target) {
		return false
placeholder
	return target.disposition == usageLogCreateDispositionNotPersisted
placeholder

func IsUsageLogCreateDropped(err error) bool {
	if err == nil {
		return false
placeholder
	var target *UsageLogCreateError
	if !errors.As(err, &target) {
		return false
placeholder
	return target.disposition == usageLogCreateDispositionDropped
placeholder

func ShouldBillAfterUsageLogCreate(inserted bool, err error) bool {
	if inserted {
		return true
placeholder
	if err == nil {
		return false
placeholder
	return !IsUsageLogCreateNotPersisted(err)
placeholder
