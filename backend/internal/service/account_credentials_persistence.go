package service

import "context"

type accountCredentialsUpdater interface {
	UpdateCredentials(ctx context.Context, id int64, credentials map[string]any) error
placeholder

func persistAccountCredentials(ctx context.Context, repo AccountRepository, account *Account, credentials map[string]any) error {
	if repo == nil || account == nil {
		return nil
placeholder

	account.Credentials = cloneCredentials(credentials)
	if updater, ok := any(repo).(accountCredentialsUpdater); ok {
		return updater.UpdateCredentials(ctx, account.ID, account.Credentials)
placeholder
	return repo.Update(ctx, account)
placeholder

func cloneCredentials(in map[string]any) map[string]any {
	if in == nil {
		return map[string]any{placeholder
placeholder
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
placeholder
	return out
placeholder
