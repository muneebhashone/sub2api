package schema

import (
	"fmt"

	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

var pendingAuthIntents = map[string]struct{placeholder{
	"login":                        {placeholder,
	"bind_current_user":            {placeholder,
	"adopt_existing_user_by_email": {placeholder,
placeholder

func validatePendingAuthIntent(value string) error {
	if _, ok := pendingAuthIntents[value]; ok {
		return nil
placeholder
	return fmt.Errorf("invalid pending auth intent %q", value)
placeholder

// PendingAuthSession stores a short-lived post-auth decision session.
type PendingAuthSession struct {
	ent.Schema
placeholder

func (PendingAuthSession) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "pending_auth_sessions"placeholder,
placeholder
placeholder

func (PendingAuthSession) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{placeholder,
placeholder
placeholder

func (PendingAuthSession) Fields() []ent.Field {
	return []ent.Field{
		field.String("session_token").
			MaxLen(255).
			NotEmpty(),
		field.String("intent").
			MaxLen(40).
			NotEmpty().
			Validate(validatePendingAuthIntent),
		field.String("provider_type").
			MaxLen(20).
			NotEmpty().
			Validate(validateAuthProviderType),
		field.String("provider_key").
			NotEmpty().
			SchemaType(map[string]string{dialect.Postgres: "text"placeholder),
		field.String("provider_subject").
			NotEmpty().
			SchemaType(map[string]string{dialect.Postgres: "text"placeholder),
		field.Int64("target_user_id").
			Optional().
			Nillable(),
		field.String("redirect_to").
			Default("").
			SchemaType(map[string]string{dialect.Postgres: "text"placeholder),
		field.String("resolved_email").
			Default("").
			SchemaType(map[string]string{dialect.Postgres: "text"placeholder),
		field.String("registration_password_hash").
			Default("").
			SchemaType(map[string]string{dialect.Postgres: "text"placeholder),
		field.JSON("upstream_identity_claims", map[string]any{placeholder).
			Default(func() map[string]any { return map[string]any{placeholder placeholder).
			SchemaType(map[string]string{dialect.Postgres: "jsonb"placeholder),
		field.JSON("local_flow_state", map[string]any{placeholder).
			Default(func() map[string]any { return map[string]any{placeholder placeholder).
			SchemaType(map[string]string{dialect.Postgres: "jsonb"placeholder),
		field.String("browser_session_key").
			Default("").
			SchemaType(map[string]string{dialect.Postgres: "text"placeholder),
		field.String("completion_code_hash").
			Default("").
			SchemaType(map[string]string{dialect.Postgres: "text"placeholder),
		field.Time("completion_code_expires_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"placeholder),
		field.Time("email_verified_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"placeholder),
		field.Time("password_verified_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"placeholder),
		field.Time("totp_verified_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"placeholder),
		field.Time("expires_at").
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"placeholder),
		field.Time("consumed_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.Postgres: "timestamptz"placeholder),
placeholder
placeholder

func (PendingAuthSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("target_user", User.Type).
			Ref("pending_auth_sessions").
			Field("target_user_id").
			Unique(),
		edge.To("adoption_decision", IdentityAdoptionDecision.Type).
			Unique(),
placeholder
placeholder

func (PendingAuthSession) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("session_token").Unique(),
		index.Fields("target_user_id"),
		index.Fields("expires_at"),
		index.Fields("provider_type", "provider_key", "provider_subject"),
		index.Fields("completion_code_hash"),
placeholder
placeholder
