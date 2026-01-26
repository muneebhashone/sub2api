/**
 * API Client for Sub2API Backend
 * Central export point for all API modules
 */

// Re-export the HTTP client
export { apiClient placeholder from './client'

// Auth API
export { authAPI, isTotp2FARequired, type LoginResponse placeholder from './auth'

// User APIs
export { keysAPI placeholder from './keys'
export { usageAPI placeholder from './usage'
export { userAPI placeholder from './user'
export { redeemAPI, type RedeemHistoryItem placeholder from './redeem'
export { userGroupsAPI placeholder from './groups'
export { totpAPI placeholder from './totp'

// Admin APIs
export { adminAPI placeholder from './admin'

// Default export
export { default placeholder from './client'
