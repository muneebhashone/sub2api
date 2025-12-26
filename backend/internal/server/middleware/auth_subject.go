package middleware

import "github.com/gin-gonic/gin"

// AuthSubject is the minimal authenticated identity stored in gin context.
// Decision: {UserID int64, Concurrency intplaceholder
type AuthSubject struct {
	UserID      int64
	Concurrency int
placeholder

func GetAuthSubjectFromContext(c *gin.Context) (AuthSubject, bool) {
	value, exists := c.Get(string(ContextKeyUser))
	if !exists {
		return AuthSubject{placeholder, false
placeholder
	subject, ok := value.(AuthSubject)
	return subject, ok
placeholder

func GetUserRoleFromContext(c *gin.Context) (string, bool) {
	value, exists := c.Get(string(ContextKeyUserRole))
	if !exists {
		return "", false
placeholder
	role, ok := value.(string)
	return role, ok
placeholder
