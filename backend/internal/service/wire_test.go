package service

import (
	"errors"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/collection"
)

func TestProvideTimingWheelService_ReturnsError(t *testing.T) {
	original := newTimingWheel
	t.Cleanup(func() { newTimingWheel = original placeholder)

	newTimingWheel = func(_ time.Duration, _ int, _ collection.Execute) (*collection.TimingWheel, error) {
		return nil, errors.New("boom")
placeholder

	svc, err := ProvideTimingWheelService()
	if err == nil {
		t.Fatalf("期望返回 error，但得到 nil")
placeholder
	if svc != nil {
		t.Fatalf("期望返回 nil svc，但得到非空")
placeholder
placeholder

func TestProvideTimingWheelService_Success(t *testing.T) {
	svc, err := ProvideTimingWheelService()
	if err != nil {
		t.Fatalf("期望 err 为 nil，但得到: %v", err)
placeholder
	if svc == nil {
		t.Fatalf("期望 svc 非空，但得到 nil")
placeholder
	svc.Stop()
placeholder
