/** WebSearch emulation mode values (must match backend WebSearchMode* constants in account.go) */
export const WEB_SEARCH_MODE_DEFAULT = 'default' as const
export const WEB_SEARCH_MODE_ENABLED = 'enabled' as const
export const WEB_SEARCH_MODE_DISABLED = 'disabled' as const
export type WebSearchMode = typeof WEB_SEARCH_MODE_DEFAULT | typeof WEB_SEARCH_MODE_ENABLED | typeof WEB_SEARCH_MODE_DISABLED

/** Quota notification threshold type values (must match thresholdType* constants in balance_notify_service.go) */
export const QUOTA_THRESHOLD_TYPE_FIXED = 'fixed' as const
export const QUOTA_THRESHOLD_TYPE_PERCENTAGE = 'percentage' as const
export type QuotaThresholdType = typeof QUOTA_THRESHOLD_TYPE_FIXED | typeof QUOTA_THRESHOLD_TYPE_PERCENTAGE

/** Quota reset mode values */
export const QUOTA_RESET_MODE_ROLLING = 'rolling' as const
export const QUOTA_RESET_MODE_FIXED = 'fixed' as const
export type QuotaResetMode = typeof QUOTA_RESET_MODE_ROLLING | typeof QUOTA_RESET_MODE_FIXED

/** Vertex AI location options for Service Account accounts */
export const VERTEX_LOCATION_OPTIONS = [
  {
    label: 'Common',
    options: [
      { value: 'us-central1', label: 'us-central1 (Iowa)' placeholder,
      { value: 'global', label: 'global' placeholder,
      { value: 'us', label: 'us' placeholder,
      { value: 'eu', label: 'eu' placeholder
    ]
  placeholder,
  {
    label: 'United States',
    options: [
      { value: 'us-east1', label: 'us-east1 (South Carolina)' placeholder,
      { value: 'us-east4', label: 'us-east4 (Northern Virginia)' placeholder,
      { value: 'us-east5', label: 'us-east5 (Columbus)' placeholder,
      { value: 'us-south1', label: 'us-south1 (Dallas)' placeholder,
      { value: 'us-west1', label: 'us-west1 (Oregon)' placeholder,
      { value: 'us-west4', label: 'us-west4 (Las Vegas)' placeholder
    ]
  placeholder,
  {
    label: 'Europe',
    options: [
      { value: 'europe-west1', label: 'europe-west1 (Belgium)' placeholder,
      { value: 'europe-west2', label: 'europe-west2 (London)' placeholder,
      { value: 'europe-west3', label: 'europe-west3 (Frankfurt)' placeholder,
      { value: 'europe-west4', label: 'europe-west4 (Netherlands)' placeholder,
      { value: 'europe-west6', label: 'europe-west6 (Zurich)' placeholder,
      { value: 'europe-west8', label: 'europe-west8 (Milan)' placeholder,
      { value: 'europe-west9', label: 'europe-west9 (Paris)' placeholder
    ]
  placeholder,
  {
    label: 'Asia Pacific',
    options: [
      { value: 'asia-east1', label: 'asia-east1 (Taiwan)' placeholder,
      { value: 'asia-east2', label: 'asia-east2 (Hong Kong)' placeholder,
      { value: 'asia-northeast1', label: 'asia-northeast1 (Tokyo)' placeholder,
      { value: 'asia-northeast3', label: 'asia-northeast3 (Seoul)' placeholder,
      { value: 'asia-south1', label: 'asia-south1 (Mumbai)' placeholder,
      { value: 'asia-southeast1', label: 'asia-southeast1 (Singapore)' placeholder,
      { value: 'australia-southeast1', label: 'australia-southeast1 (Sydney)' placeholder
    ]
  placeholder
] as const
