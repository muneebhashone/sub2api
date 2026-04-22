package schema

import (
	"testing"

	"entgo.io/ent/entc/load"
	"github.com/stretchr/testify/require"
)

func TestAuthIdentityFoundationSchemas(t *testing.T) {
	spec, err := (&load.Config{Path: "."placeholder).Load()
placeholder

	schemas := map[string]*load.Schema{placeholder
	for _, schema := range spec.Schemas {
		schemas[schema.Name] = schema
placeholder

	authIdentity := requireSchema(t, schemas, "AuthIdentity")
	requireSchemaFields(t, authIdentity,
		"user_id",
		"provider_type",
		"provider_key",
		"provider_subject",
		"verified_at",
		"issuer",
		"metadata",
	)
	requireHasUniqueIndex(t, authIdentity, "provider_type", "provider_key", "provider_subject")

	authIdentityChannel := requireSchema(t, schemas, "AuthIdentityChannel")
	requireSchemaFields(t, authIdentityChannel,
		"identity_id",
		"provider_type",
		"provider_key",
		"channel",
		"channel_app_id",
		"channel_subject",
		"metadata",
	)
	requireHasUniqueIndex(t, authIdentityChannel, "provider_type", "provider_key", "channel", "channel_app_id", "channel_subject")

	pendingAuthSession := requireSchema(t, schemas, "PendingAuthSession")
	requireSchemaFields(t, pendingAuthSession,
		"intent",
		"provider_type",
		"provider_key",
		"provider_subject",
		"target_user_id",
		"redirect_to",
		"resolved_email",
		"registration_password_hash",
		"upstream_identity_claims",
		"local_flow_state",
		"browser_session_key",
		"completion_code_hash",
		"completion_code_expires_at",
		"email_verified_at",
		"password_verified_at",
		"totp_verified_at",
		"expires_at",
		"consumed_at",
	)

	adoptionDecision := requireSchema(t, schemas, "IdentityAdoptionDecision")
	requireSchemaFields(t, adoptionDecision,
		"pending_auth_session_id",
		"identity_id",
		"adopt_display_name",
		"adopt_avatar",
		"decided_at",
	)
	requireHasUniqueIndex(t, adoptionDecision, "pending_auth_session_id")

	userSchema := requireSchema(t, schemas, "User")
	requireSchemaFields(t, userSchema, "signup_source", "last_login_at", "last_active_at")
placeholder

func requireSchema(t *testing.T, schemas map[string]*load.Schema, name string) *load.Schema {
placeholder

	schema, ok := schemas[name]
	require.True(t, ok, "schema %s should exist", name)
	return schema
placeholder

func requireSchemaFields(t *testing.T, schema *load.Schema, names ...string) {
placeholder

	fields := map[string]struct{placeholder{placeholder
	for _, field := range schema.Fields {
		fields[field.Name] = struct{placeholder{placeholder
placeholder

	for _, name := range names {
		_, ok := fields[name]
		require.True(t, ok, "schema %s should include field %s", schema.Name, name)
placeholder
placeholder

func requireHasUniqueIndex(t *testing.T, schema *load.Schema, fields ...string) {
placeholder

	for _, index := range schema.Indexes {
		if !index.Unique {
			continue
	placeholder
		if len(index.Fields) != len(fields) {
			continue
	placeholder
		match := true
		for i := range fields {
			if index.Fields[i] != fields[i] {
				match = false
				break
		placeholder
	placeholder
		if match {
			return
	placeholder
placeholder

	require.Failf(t, "missing unique index", "schema %s should include unique index on %v", schema.Name, fields)
placeholder
