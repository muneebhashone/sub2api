//go:build unit

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatePlanRequired_AllValid(t *testing.T) {
	err := validatePlanRequired("Pro", 1, 9.99, 30, "days", nil)
placeholder
placeholder

func TestValidatePlanRequired_EmptyName(t *testing.T) {
	err := validatePlanRequired("", 1, 9.99, 30, "days", nil)
placeholder
	require.Contains(t, err.Error(), "plan name")
placeholder

func TestValidatePlanRequired_WhitespaceName(t *testing.T) {
	err := validatePlanRequired("   ", 1, 9.99, 30, "days", nil)
placeholder
	require.Contains(t, err.Error(), "plan name")
placeholder

func TestValidatePlanRequired_ZeroGroupID(t *testing.T) {
	err := validatePlanRequired("Pro", 0, 9.99, 30, "days", nil)
placeholder
	require.Contains(t, err.Error(), "group")
placeholder

func TestValidatePlanRequired_NegativeGroupID(t *testing.T) {
	err := validatePlanRequired("Pro", -1, 9.99, 30, "days", nil)
placeholder
	require.Contains(t, err.Error(), "group")
placeholder

func TestValidatePlanRequired_ZeroPrice(t *testing.T) {
	err := validatePlanRequired("Pro", 1, 0, 30, "days", nil)
placeholder
	require.Contains(t, err.Error(), "price")
placeholder

func TestValidatePlanRequired_NegativePrice(t *testing.T) {
	err := validatePlanRequired("Pro", 1, -5, 30, "days", nil)
placeholder
	require.Contains(t, err.Error(), "price")
placeholder

func TestValidatePlanRequired_ZeroValidityDays(t *testing.T) {
	err := validatePlanRequired("Pro", 1, 9.99, 0, "days", nil)
placeholder
	require.Contains(t, err.Error(), "validity days")
placeholder

func TestValidatePlanRequired_NegativeValidityDays(t *testing.T) {
	err := validatePlanRequired("Pro", 1, 9.99, -7, "days", nil)
placeholder
	require.Contains(t, err.Error(), "validity days")
placeholder

func TestValidatePlanRequired_EmptyValidityUnit(t *testing.T) {
	err := validatePlanRequired("Pro", 1, 9.99, 30, "", nil)
placeholder
	require.Contains(t, err.Error(), "validity unit")
placeholder

func TestValidatePlanRequired_WhitespaceValidityUnit(t *testing.T) {
	err := validatePlanRequired("Pro", 1, 9.99, 30, "   ", nil)
placeholder
	require.Contains(t, err.Error(), "validity unit")
placeholder

func TestValidatePlanRequired_NameValidatedFirst(t *testing.T) {
	err := validatePlanRequired("", 0, 0, 0, "", nil)
placeholder
	require.Contains(t, err.Error(), "plan name")
placeholder

func TestValidatePlanRequired_TrimmedValidName(t *testing.T) {
	err := validatePlanRequired("  Pro  ", 1, 9.99, 30, "days", nil)
placeholder
placeholder

func TestValidatePlanRequired_NegativeOriginalPrice(t *testing.T) {
	neg := -10.0
	err := validatePlanRequired("Pro", 1, 9.99, 30, "days", &neg)
placeholder
	require.Contains(t, err.Error(), "original price")
placeholder

func TestValidatePlanRequired_ZeroOriginalPrice(t *testing.T) {
	zero := 0.0
	err := validatePlanRequired("Pro", 1, 9.99, 30, "days", &zero)
placeholder
placeholder

func TestValidatePlanRequired_ValidOriginalPrice(t *testing.T) {
	op := 19.99
	err := validatePlanRequired("Pro", 1, 9.99, 30, "days", &op)
placeholder
placeholder

// --- validatePlanPatch tests ---

func TestValidatePlanPatch_NegativeOriginalPrice(t *testing.T) {
	neg := -5.0
	err := validatePlanPatch(UpdatePlanRequest{OriginalPrice: &negplaceholder)
placeholder
	require.Contains(t, err.Error(), "original price")
placeholder

func TestValidatePlanPatch_ZeroOriginalPrice(t *testing.T) {
	zero := 0.0
	err := validatePlanPatch(UpdatePlanRequest{OriginalPrice: &zeroplaceholder)
placeholder
placeholder

func TestValidatePlanPatch_ValidOriginalPrice(t *testing.T) {
	op := 29.99
	err := validatePlanPatch(UpdatePlanRequest{OriginalPrice: &opplaceholder)
placeholder
placeholder

func TestValidatePlanPatch_NilOriginalPrice(t *testing.T) {
	err := validatePlanPatch(UpdatePlanRequest{OriginalPrice: nilplaceholder)
placeholder
placeholder
