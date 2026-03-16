<template>
  <div class="card p-6">
    <!-- Toolbar: left filters (multi-line) + right actions -->
    <div class="flex flex-wrap items-end justify-between gap-4">
      <!-- Left: filters (allowed to wrap to multiple rows) -->
      <div class="flex flex-1 flex-wrap items-end gap-4">
        <!-- User Search -->
        <div ref="userSearchRef" class="usage-filter-dropdown relative w-full sm:w-auto sm:min-w-[240px]">
          <label class="input-label">{{ t('admin.usage.userFilter') placeholderplaceholder</label>
          <input
            v-model="userKeyword"
            type="text"
            class="input pr-8"
            :placeholder="t('admin.usage.searchUserPlaceholder')"
            @input="debounceUserSearch"
            @focus="showUserDropdown = true"
          />
          <button
            v-if="filters.user_id"
            type="button"
            @click="clearUser"
            class="absolute right-2 top-9 text-gray-400"
            aria-label="Clear user filter"
          >
            ✕
          </button>
          <div
            v-if="showUserDropdown && (userResults.length > 0 || userKeyword)"
            class="absolute z-50 mt-1 max-h-60 w-full overflow-auto rounded-lg border bg-white shadow-lg dark:bg-gray-800"
          >
            <button
              v-for="u in userResults"
              :key="u.id"
              type="button"
              @click="selectUser(u)"
              class="w-full px-4 py-2 text-left hover:bg-gray-100 dark:hover:bg-gray-700"
            >
              <span>{{ u.email placeholderplaceholder</span>
              <span class="ml-2 text-xs text-gray-400">#{{ u.id placeholderplaceholder</span>
            </button>
          </div>
        </div>

        <!-- API Key Search -->
        <div ref="apiKeySearchRef" class="usage-filter-dropdown relative w-full sm:w-auto sm:min-w-[240px]">
          <label class="input-label">{{ t('usage.apiKeyFilter') placeholderplaceholder</label>
          <input
            v-model="apiKeyKeyword"
            type="text"
            class="input pr-8"
            :placeholder="t('admin.usage.searchApiKeyPlaceholder')"
            @input="debounceApiKeySearch"
            @focus="onApiKeyFocus"
          />
          <button
            v-if="filters.api_key_id"
            type="button"
            @click="onClearApiKey"
            class="absolute right-2 top-9 text-gray-400"
            aria-label="Clear API key filter"
          >
            ✕
          </button>
          <div
            v-if="showApiKeyDropdown && apiKeyResults.length > 0"
            class="absolute z-50 mt-1 max-h-60 w-full overflow-auto rounded-lg border bg-white shadow-lg dark:bg-gray-800"
          >
            <button
              v-for="k in apiKeyResults"
              :key="k.id"
              type="button"
              @click="selectApiKey(k)"
              class="w-full px-4 py-2 text-left hover:bg-gray-100 dark:hover:bg-gray-700"
            >
              <span class="truncate">{{ k.name || `#${k.idplaceholder` placeholderplaceholder</span>
              <span class="ml-2 text-xs text-gray-400">#{{ k.id placeholderplaceholder</span>
            </button>
          </div>
        </div>

        <!-- Model Filter -->
        <div class="w-full sm:w-auto sm:min-w-[220px]">
          <label class="input-label">{{ t('usage.model') placeholderplaceholder</label>
          <Select v-model="filters.model" :options="modelOptions" searchable @change="emitChange" />
        </div>

        <!-- Account Filter -->
        <div ref="accountSearchRef" class="usage-filter-dropdown relative w-full sm:w-auto sm:min-w-[220px]">
          <label class="input-label">{{ t('admin.usage.account') placeholderplaceholder</label>
          <input
            v-model="accountKeyword"
            type="text"
            class="input pr-8"
            :placeholder="t('admin.usage.searchAccountPlaceholder')"
            @input="debounceAccountSearch"
            @focus="showAccountDropdown = true"
          />
          <button
            v-if="filters.account_id"
            type="button"
            @click="clearAccount"
            class="absolute right-2 top-9 text-gray-400"
            aria-label="Clear account filter"
          >
            ✕
          </button>
          <div
            v-if="showAccountDropdown && (accountResults.length > 0 || accountKeyword)"
            class="absolute z-50 mt-1 max-h-60 w-full overflow-auto rounded-lg border bg-white shadow-lg dark:bg-gray-800"
          >
            <button
              v-for="a in accountResults"
              :key="a.id"
              type="button"
              @click="selectAccount(a)"
              class="w-full px-4 py-2 text-left hover:bg-gray-100 dark:hover:bg-gray-700"
            >
              <span class="truncate">{{ a.name placeholderplaceholder</span>
              <span class="ml-2 text-xs text-gray-400">#{{ a.id placeholderplaceholder</span>
            </button>
          </div>
        </div>

        <!-- Request Type Filter -->
        <div class="w-full sm:w-auto sm:min-w-[180px]">
          <label class="input-label">{{ t('usage.type') placeholderplaceholder</label>
          <Select v-model="filters.request_type" :options="requestTypeOptions" @change="emitChange" />
        </div>

        <!-- Billing Type Filter -->
        <div class="w-full sm:w-auto sm:min-w-[200px]">
          <label class="input-label">{{ t('admin.usage.billingType') placeholderplaceholder</label>
          <Select v-model="filters.billing_type" :options="billingTypeOptions" @change="emitChange" />
        </div>

        <!-- Group Filter -->
        <div class="w-full sm:w-auto sm:min-w-[200px]">
          <label class="input-label">{{ t('admin.usage.group') placeholderplaceholder</label>
          <Select v-model="filters.group_id" :options="groupOptions" searchable @change="emitChange" />
        </div>

      </div>

      <!-- Right: actions -->
      <div v-if="showActions" class="flex w-full flex-wrap items-center justify-end gap-3 sm:w-auto">
        <button type="button" @click="$emit('refresh')" class="btn btn-secondary">
          {{ t('common.refresh') placeholderplaceholder
        </button>
        <button type="button" @click="$emit('reset')" class="btn btn-secondary">
          {{ t('common.reset') placeholderplaceholder
        </button>
        <slot name="after-reset" />
        <button type="button" @click="$emit('cleanup')" class="btn btn-danger">
          {{ t('admin.usage.cleanup.button') placeholderplaceholder
        </button>
        <button type="button" @click="$emit('export')" :disabled="exporting" class="btn btn-primary">
          {{ t('usage.exportExcel') placeholderplaceholder
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, toRef, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { adminAPI placeholder from '@/api/admin'
import Select, { type SelectOption placeholder from '@/components/common/Select.vue'
import type { SimpleApiKey, SimpleUser placeholder from '@/api/admin/usage'

type ModelValue = Record<string, any>

interface Props {
  modelValue: ModelValue
  exporting: boolean
  startDate: string
  endDate: string
  showActions?: boolean
placeholder

const props = withDefaults(defineProps<Props>(), {
  showActions: true
placeholder)
const emit = defineEmits([
  'update:modelValue',
  'change',
  'refresh',
  'reset',
  'export',
  'cleanup'
])

const { t placeholder = useI18n()
const filters = toRef(props, 'modelValue')

const userSearchRef = ref<HTMLElement | null>(null)
const apiKeySearchRef = ref<HTMLElement | null>(null)
const accountSearchRef = ref<HTMLElement | null>(null)

const userKeyword = ref('')
const userResults = ref<SimpleUser[]>([])
const showUserDropdown = ref(false)
let userSearchTimeout: ReturnType<typeof setTimeout> | null = null

const apiKeyKeyword = ref('')
const apiKeyResults = ref<SimpleApiKey[]>([])
const showApiKeyDropdown = ref(false)
let apiKeySearchTimeout: ReturnType<typeof setTimeout> | null = null

interface SimpleAccount {
  id: number
  name: string
placeholder
const accountKeyword = ref('')
const accountResults = ref<SimpleAccount[]>([])
const showAccountDropdown = ref(false)
let accountSearchTimeout: ReturnType<typeof setTimeout> | null = null

const modelOptions = ref<SelectOption[]>([{ value: null, label: t('admin.usage.allModels') placeholder])
const groupOptions = ref<SelectOption[]>([{ value: null, label: t('admin.usage.allGroups') placeholder])

const requestTypeOptions = ref<SelectOption[]>([
  { value: null, label: t('admin.usage.allTypes') placeholder,
  { value: 'ws_v2', label: t('usage.ws') placeholder,
  { value: 'stream', label: t('usage.stream') placeholder,
  { value: 'sync', label: t('usage.sync') placeholder
])

const billingTypeOptions = ref<SelectOption[]>([
  { value: null, label: t('admin.usage.allBillingTypes') placeholder,
  { value: 0, label: t('admin.usage.billingTypeBalance') placeholder,
  { value: 1, label: t('admin.usage.billingTypeSubscription') placeholder
])

const emitChange = () => emit('change')

const debounceUserSearch = () => {
  if (userSearchTimeout) clearTimeout(userSearchTimeout)
  userSearchTimeout = setTimeout(async () => {
    if (!userKeyword.value) {
      userResults.value = []
      return
    placeholder
    try {
      userResults.value = await adminAPI.usage.searchUsers(userKeyword.value)
    placeholder catch {
      userResults.value = []
    placeholder
  placeholder, 300)
placeholder

const debounceApiKeySearch = () => {
  if (apiKeySearchTimeout) clearTimeout(apiKeySearchTimeout)
  apiKeySearchTimeout = setTimeout(async () => {
    try {
      apiKeyResults.value = await adminAPI.usage.searchApiKeys(
        filters.value.user_id,
        apiKeyKeyword.value || ''
      )
    placeholder catch {
      apiKeyResults.value = []
    placeholder
  placeholder, 300)
placeholder

const selectUser = async (u: SimpleUser) => {
  userKeyword.value = u.email
  showUserDropdown.value = false
  filters.value.user_id = u.id
  clearApiKey()

  // Auto-load API keys for this user
  try {
    apiKeyResults.value = await adminAPI.usage.searchApiKeys(u.id, '')
  placeholder catch {
    apiKeyResults.value = []
  placeholder

  emitChange()
placeholder

const clearUser = () => {
  userKeyword.value = ''
  userResults.value = []
  showUserDropdown.value = false
  filters.value.user_id = undefined
  clearApiKey()
  emitChange()
placeholder

const selectApiKey = (k: SimpleApiKey) => {
  apiKeyKeyword.value = k.name || String(k.id)
  showApiKeyDropdown.value = false
  filters.value.api_key_id = k.id
  emitChange()
placeholder

const clearApiKey = () => {
  apiKeyKeyword.value = ''
  apiKeyResults.value = []
  showApiKeyDropdown.value = false
  filters.value.api_key_id = undefined
placeholder

const onClearApiKey = () => {
  clearApiKey()
  emitChange()
placeholder

const debounceAccountSearch = () => {
  if (accountSearchTimeout) clearTimeout(accountSearchTimeout)
  accountSearchTimeout = setTimeout(async () => {
    if (!accountKeyword.value) {
      accountResults.value = []
      return
    placeholder
    try {
      const res = await adminAPI.accounts.list(1, 20, { search: accountKeyword.value placeholder)
      accountResults.value = res.items.map((a) => ({ id: a.id, name: a.name placeholder))
    placeholder catch {
      accountResults.value = []
    placeholder
  placeholder, 300)
placeholder

const selectAccount = (a: SimpleAccount) => {
  accountKeyword.value = a.name
  showAccountDropdown.value = false
  filters.value.account_id = a.id
  emitChange()
placeholder

const clearAccount = () => {
  accountKeyword.value = ''
  accountResults.value = []
  showAccountDropdown.value = false
  filters.value.account_id = undefined
  emitChange()
placeholder

const onApiKeyFocus = () => {
  showApiKeyDropdown.value = true
  // Trigger search if no results yet
  if (apiKeyResults.value.length === 0) {
    debounceApiKeySearch()
  placeholder
placeholder

const onDocumentClick = (e: MouseEvent) => {
  const target = e.target as Node | null
  if (!target) return

  const clickedInsideUser = userSearchRef.value?.contains(target) ?? false
  const clickedInsideApiKey = apiKeySearchRef.value?.contains(target) ?? false
  const clickedInsideAccount = accountSearchRef.value?.contains(target) ?? false

  if (!clickedInsideUser) showUserDropdown.value = false
  if (!clickedInsideApiKey) showApiKeyDropdown.value = false
  if (!clickedInsideAccount) showAccountDropdown.value = false
placeholder

watch(
  () => props.startDate,
  (value) => {
    filters.value.start_date = value
  placeholder,
  { immediate: true placeholder
)

watch(
  () => props.endDate,
  (value) => {
    filters.value.end_date = value
  placeholder,
  { immediate: true placeholder
)

watch(
  () => filters.value.user_id,
  (userId) => {
    if (!userId) {
      userKeyword.value = ''
      userResults.value = []
    placeholder
  placeholder
)

watch(
  () => filters.value.api_key_id,
  (apiKeyId) => {
    if (!apiKeyId) {
      apiKeyKeyword.value = ''
      apiKeyResults.value = []
    placeholder
  placeholder
)

watch(
  () => filters.value.account_id,
  (accountId) => {
    if (!accountId) {
      accountKeyword.value = ''
      accountResults.value = []
    placeholder
  placeholder
)

onMounted(async () => {
  document.addEventListener('click', onDocumentClick)

  try {
    const [gs, ms] = await Promise.all([
      adminAPI.groups.list(1, 1000),
      adminAPI.dashboard.getModelStats({ start_date: props.startDate, end_date: props.endDate placeholder)
    ])

    groupOptions.value.push(...gs.items.map((g: any) => ({ value: g.id, label: g.name placeholder)))

    const uniqueModels = new Set<string>()
    ms.models?.forEach((s: any) => {
      if (s.model) {
        uniqueModels.add(s.model)
      placeholder
    placeholder)
    modelOptions.value.push(
      ...Array.from(uniqueModels)
        .sort()
        .map((m) => ({ value: m, label: m placeholder))
    )
  placeholder catch {
    // Ignore filter option loading errors (page still usable)
  placeholder
placeholder)

onUnmounted(() => {
  document.removeEventListener('click', onDocumentClick)
placeholder)
</script>
