package repository

import (
	"errors"
	"strings"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/infrastructure/errors"
	"gorm.io/gorm"
)

func translatePersistenceError(err error, notFound, conflict *infraerrors.ApplicationError) error {
	if err == nil {
		return nil
placeholder

	if notFound != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound.WithCause(err)
placeholder

	if conflict != nil && isUniqueConstraintViolation(err) {
		return conflict.WithCause(err)
placeholder

	return err
placeholder

func isUniqueConstraintViolation(err error) bool {
	if err == nil {
		return false
placeholder

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
placeholder

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "duplicate entry")
placeholder
