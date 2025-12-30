/**
 * Pinia Stores Export
 * Central export point for all application stores
 */

export { useAuthStore placeholder from './auth'
export { useAppStore placeholder from './app'
export { useSubscriptionStore placeholder from './subscriptions'
export { useOnboardingStore placeholder from './onboarding'

// Re-export types for convenience
export type { User, LoginRequest, RegisterRequest, AuthResponse placeholder from '@/types'
export type { Toast, ToastType, AppState placeholder from '@/types'
