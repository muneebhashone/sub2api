<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-wrap items-center gap-3">
          <div class="relative w-full md:w-80">
            <Icon name="search" size="md" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            <input v-model="filters.search" type="text" class="input pl-10" :placeholder="t('admin.affiliates.records.searchPlaceholder')" @input="debounceLoad" />
          </div>
          <input v-model="filters.start_at" type="date" class="input w-full sm:w-44" :title="t('admin.affiliates.records.startAt')" @change="reloadFromFirstPage" />
          <input v-model="filters.end_at" type="date" class="input w-full sm:w-44" :title="t('admin.affiliates.records.endAt')" @change="reloadFromFirstPage" />
          <button class="btn btn-secondary px-2 md:px-3" :disabled="loading" :title="t('common.refresh')" @click="loadRecords">
            <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
          </button>
        </div>
      </template>

      <template #table>
        <DataTable
          :columns="columns"
          :data="records"
          :loading="loading"
          :server-side-sort="true"
          default-sort-key="created_at"
          default-sort-order="desc"
          :sort-storage-key="sortStorageKey"
          @sort="handleSort"
        >
          <template #cell-inviter="{ row placeholder">
            <UserCell
              :id="row.inviter_id"
              :email="row.inviter_email"
              :username="row.inviter_username"
              :clickable="props.type !== 'transfers'"
              @open="openUserOverview"
            />
          </template>
          <template #cell-invitee="{ row placeholder">
            <UserCell
              :id="row.invitee_id"
              :email="row.invitee_email"
              :username="row.invitee_username"
              :clickable="props.type !== 'transfers'"
              @open="openUserOverview"
            />
          </template>
          <template #cell-user="{ row placeholder">
            <UserCell
              :id="row.user_id"
              :email="row.user_email"
              :username="row.username"
              :clickable="true"
              @open="openUserOverview"
            />
          </template>
          <template #cell-aff_code="{ row placeholder">
            <span class="font-mono text-sm text-gray-700 dark:text-gray-300">{{ row.aff_code || '-' placeholderplaceholder</span>
          </template>
          <template #cell-order="{ row placeholder">
            <div class="space-y-0.5">
              <div class="font-mono text-sm text-gray-900 dark:text-white">#{{ row.order_id placeholderplaceholder</div>
              <div class="max-w-56 truncate text-sm text-gray-500 dark:text-dark-400">{{ row.out_trade_no placeholderplaceholder</div>
            </div>
          </template>
          <template #cell-payment_type="{ row placeholder">
            {{ t('payment.methods.' + row.payment_type, row.payment_type || '-') placeholderplaceholder
          </template>
          <template #cell-order_status="{ row placeholder">
            <OrderStatusBadge :status="row.order_status" />
          </template>
          <template #cell-total_rebate="{ row placeholder">
            <AmountText :value="row.total_rebate" />
          </template>
          <template #cell-order_amount="{ row placeholder">
            <AmountText :value="row.order_amount" />
          </template>
          <template #cell-pay_amount="{ row placeholder">
            <span class="text-sm text-gray-900 dark:text-white">¥{{ formatAmount(row.pay_amount) placeholderplaceholder</span>
          </template>
          <template #cell-rebate_amount="{ row placeholder">
            <AmountText :value="row.rebate_amount" strong />
          </template>
          <template #cell-amount="{ row placeholder">
            <AmountText :value="row.amount" strong />
          </template>
          <template #cell-balance_after="{ row placeholder">
            <NullableAmountText :value="row.balance_after" />
          </template>
          <template #cell-available_quota_after="{ row placeholder">
            <NullableAmountText :value="row.available_quota_after" />
          </template>
          <template #cell-frozen_quota_after="{ row placeholder">
            <NullableAmountText :value="row.frozen_quota_after" />
          </template>
          <template #cell-history_quota_after="{ row placeholder">
            <NullableAmountText :value="row.history_quota_after" />
          </template>
          <template #cell-created_at="{ row placeholder">
            <span class="text-sm text-gray-700 dark:text-gray-300">{{ formatDateTime(row.created_at) placeholderplaceholder</span>
          </template>
        </DataTable>
      </template>

      <template #pagination>
        <Pagination
          v-if="pagination.total > 0"
          :page="pagination.page"
          :total="pagination.total"
          :page-size="pagination.page_size"
          @update:page="handlePageChange"
          @update:pageSize="handlePageSizeChange"
        />
      </template>
    </TablePageLayout>

    <BaseDialog
      :show="overviewDialog"
      :title="t('admin.affiliates.overview.title')"
      width="normal"
      @close="overviewDialog = false"
    >
      <div v-if="overviewLoading" class="flex justify-center py-8">
        <div class="h-6 w-6 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
      </div>
      <div v-else-if="selectedOverview" class="space-y-4">
        <div class="rounded-lg border border-gray-100 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800">
          <div class="font-mono text-sm text-gray-900 dark:text-white">#{{ selectedOverview.user_id placeholderplaceholder</div>
          <div class="mt-1 text-sm font-medium text-gray-900 dark:text-white">{{ selectedOverview.email || '-' placeholderplaceholder</div>
          <div class="mt-0.5 text-sm text-gray-500 dark:text-dark-400">{{ selectedOverview.username || '-' placeholderplaceholder</div>
        </div>
        <div class="grid gap-3 sm:grid-cols-2">
          <OverviewStat :label="t('admin.affiliates.overview.affCode')" :value="selectedOverview.aff_code || '-'" mono />
          <OverviewStat :label="t('admin.affiliates.overview.rebateRate')" :value="formatPercent(selectedOverview.rebate_rate_percent)" />
          <OverviewStat :label="t('admin.affiliates.overview.invitedCount')" :value="String(selectedOverview.invited_count)" />
          <OverviewStat :label="t('admin.affiliates.overview.rebatedInviteeCount')" :value="String(selectedOverview.rebated_invitee_count)" />
          <OverviewStat :label="t('admin.affiliates.overview.availableQuota')" :value="'$' + formatAmount(selectedOverview.available_quota)" />
          <OverviewStat :label="t('admin.affiliates.overview.historyQuota')" :value="'$' + formatAmount(selectedOverview.history_quota)" />
        </div>
      </div>
    </BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, defineComponent, h, onMounted, reactive, ref, type PropType placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'
import OrderStatusBadge from '@/components/payment/OrderStatusBadge.vue'
import type { Column placeholder from '@/components/common/types'
import { useAppStore placeholder from '@/stores/app'
import { affiliatesAPI, type AffiliateInviteRecord, type AffiliateRebateRecord, type AffiliateTransferRecord, type AffiliateUserOverview, type ListAffiliateRecordsParams placeholder from '@/api/admin/affiliates'
import type { PaginatedResponse placeholder from '@/types'
import { extractI18nErrorMessage placeholder from '@/utils/apiError'
import { formatDateTime as formatDisplayDateTime placeholder from '@/utils/format'

type RecordType = 'invites' | 'rebates' | 'transfers'
type AffiliateRecord = AffiliateInviteRecord | AffiliateRebateRecord | AffiliateTransferRecord

const props = defineProps<{
  type: RecordType
placeholder>()

const { t placeholder = useI18n()
const appStore = useAppStore()
const loading = ref(false)
const records = ref<AffiliateRecord[]>([])
const filters = reactive({ search: '', start_at: '', end_at: '' placeholder)
const pagination = reactive({ page: 1, page_size: 20, total: 0 placeholder)
const overviewDialog = ref(false)
const overviewLoading = ref(false)
const selectedOverview = ref<AffiliateUserOverview | null>(null)
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const columns = computed<Column[]>(() => {
  if (props.type === 'invites') {
    return [
      { key: 'inviter', label: t('admin.affiliates.records.inviter'), sortable: true placeholder,
      { key: 'invitee', label: t('admin.affiliates.records.invitee'), sortable: true placeholder,
      { key: 'aff_code', label: t('admin.affiliates.records.affCode'), sortable: true placeholder,
      { key: 'total_rebate', label: t('admin.affiliates.records.totalRebate'), sortable: true placeholder,
      { key: 'created_at', label: t('admin.affiliates.records.invitedAt'), sortable: true placeholder,
    ]
  placeholder
  if (props.type === 'rebates') {
    return [
      { key: 'order', label: t('admin.affiliates.records.order'), sortable: true placeholder,
      { key: 'inviter', label: t('admin.affiliates.records.inviter'), sortable: true placeholder,
      { key: 'invitee', label: t('admin.affiliates.records.invitee'), sortable: true placeholder,
      { key: 'order_amount', label: t('admin.affiliates.records.orderAmount'), sortable: true placeholder,
      { key: 'pay_amount', label: t('admin.affiliates.records.payAmount'), sortable: true placeholder,
      { key: 'rebate_amount', label: t('admin.affiliates.records.rebateAmount') placeholder,
      { key: 'payment_type', label: t('admin.affiliates.records.paymentType'), sortable: true placeholder,
      { key: 'order_status', label: t('admin.affiliates.records.orderStatus'), sortable: true placeholder,
      { key: 'created_at', label: t('admin.affiliates.records.rebatedAt'), sortable: true placeholder,
    ]
  placeholder
  return [
    { key: 'user', label: t('admin.affiliates.records.user'), sortable: true placeholder,
    { key: 'amount', label: t('admin.affiliates.records.transferAmount'), sortable: true placeholder,
    { key: 'balance_after', label: t('admin.affiliates.records.balanceAfter'), sortable: true placeholder,
    { key: 'available_quota_after', label: t('admin.affiliates.records.availableQuotaAfter'), sortable: true placeholder,
    { key: 'frozen_quota_after', label: t('admin.affiliates.records.frozenQuotaAfter'), sortable: true placeholder,
    { key: 'history_quota_after', label: t('admin.affiliates.records.historyQuotaAfter'), sortable: true placeholder,
    { key: 'created_at', label: t('admin.affiliates.records.transferredAt'), sortable: true placeholder,
  ]
placeholder)

const sortStorageKey = computed(() => `admin-affiliate-${props.typeplaceholder-table-sort`)

function loadInitialSortState(): { sort_by: string; sort_order: 'asc' | 'desc' placeholder {
  const fallback = { sort_by: 'created_at', sort_order: 'desc' as 'asc' | 'desc' placeholder
  try {
    const raw = localStorage.getItem(sortStorageKey.value)
    if (!raw) return fallback
    const parsed = JSON.parse(raw) as { key?: string; order?: string placeholder
    const key = typeof parsed.key === 'string' ? parsed.key : ''
    if (!columns.value.some((column) => column.key === key && column.sortable)) return fallback
    return {
      sort_by: key,
      sort_order: parsed.order === 'asc' ? 'asc' : 'desc',
    placeholder
  placeholder catch {
    return fallback
  placeholder
placeholder

const sortState = reactive(loadInitialSortState())

function userTimezone(): string {
  try {
    return Intl.DateTimeFormat().resolvedOptions().timeZone
  placeholder catch {
    return 'UTC'
  placeholder
placeholder

function buildParams(): ListAffiliateRecordsParams {
  return {
    page: pagination.page,
    page_size: pagination.page_size,
    search: filters.search.trim() || undefined,
    start_at: filters.start_at || undefined,
    end_at: filters.end_at || undefined,
    sort_by: sortState.sort_by,
    sort_order: sortState.sort_order,
    timezone: userTimezone(),
  placeholder
placeholder

async function fetchRecords(params: ListAffiliateRecordsParams): Promise<PaginatedResponse<AffiliateRecord>> {
  if (props.type === 'invites') {
    return affiliatesAPI.listInviteRecords(params)
  placeholder
  if (props.type === 'rebates') {
    return affiliatesAPI.listRebateRecords(params)
  placeholder
  return affiliatesAPI.listTransferRecords(params)
placeholder

async function loadRecords() {
  loading.value = true
  try {
    const res = await fetchRecords(buildParams())
    records.value = res.items || []
    pagination.total = res.total || 0
  placeholder catch (error) {
    appStore.showError(extractI18nErrorMessage(error, t, 'admin.affiliates.errors', t('common.error')))
  placeholder finally {
    loading.value = false
  placeholder
placeholder

function debounceLoad() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => reloadFromFirstPage(), 300)
placeholder

function reloadFromFirstPage() {
  pagination.page = 1
  void loadRecords()
placeholder

function handlePageChange(page: number) {
  pagination.page = page
  void loadRecords()
placeholder

function handlePageSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  void loadRecords()
placeholder

function handleSort(key: string, order: 'asc' | 'desc') {
  sortState.sort_by = key
  sortState.sort_order = order
  pagination.page = 1
  void loadRecords()
placeholder

function formatAmount(value: number | null | undefined): string {
  return Number(value || 0).toFixed(2)
placeholder

function formatPercent(value: number | null | undefined): string {
  const rounded = Math.round(Number(value || 0) * 100) / 100
  return `${Number.isInteger(rounded) ? rounded.toString() : rounded.toString()placeholder%`
placeholder

function formatDateTime(value: string | null | undefined): string {
  return value ? formatDisplayDateTime(value) : '-'
placeholder

async function openUserOverview(userId: number) {
  if (!userId) return
  overviewDialog.value = true
  overviewLoading.value = true
  selectedOverview.value = null
  try {
    selectedOverview.value = await affiliatesAPI.getUserOverview(userId)
  placeholder catch (error) {
    overviewDialog.value = false
    appStore.showError(extractI18nErrorMessage(error, t, 'admin.affiliates.errors', t('common.error')))
  placeholder finally {
    overviewLoading.value = false
  placeholder
placeholder

const UserCell = defineComponent({
  props: {
    id: { type: Number, required: true placeholder,
    email: { type: String, default: '' placeholder,
    username: { type: String, default: '' placeholder,
    clickable: { type: Boolean, default: false placeholder,
  placeholder,
  emits: ['open'],
  setup(cellProps, { emit placeholder) {
    return () => h('div', { class: 'space-y-0.5' placeholder, [
      h('div', { class: 'font-mono text-sm text-gray-900 dark:text-white' placeholder, `#${cellProps.idplaceholder`),
      h(cellProps.clickable ? 'button' : 'div', {
        class: cellProps.clickable
          ? 'max-w-56 truncate text-left text-sm font-medium text-primary-600 hover:text-primary-700 hover:underline dark:text-primary-400 dark:hover:text-primary-300'
          : 'max-w-56 truncate text-sm text-gray-700 dark:text-gray-300',
        type: cellProps.clickable ? 'button' : undefined,
        onClick: cellProps.clickable ? () => emit('open', cellProps.id) : undefined,
      placeholder, cellProps.email || '-'),
      h('div', { class: 'max-w-56 truncate text-sm text-gray-500 dark:text-dark-400' placeholder, cellProps.username || '-'),
    ])
  placeholder,
placeholder)

const AmountText = defineComponent({
  props: {
    value: { type: Number, default: 0 placeholder,
    strong: { type: Boolean, default: false placeholder,
  placeholder,
  setup(amountProps) {
    return () => h('span', {
      class: amountProps.strong
        ? 'text-sm font-semibold text-emerald-600 dark:text-emerald-400'
        : 'text-sm text-gray-900 dark:text-white',
    placeholder, `$${formatAmount(amountProps.value)placeholder`)
  placeholder,
placeholder)

const NullableAmountText = defineComponent({
  props: {
    value: { type: Number as PropType<number | null | undefined>, default: null placeholder,
  placeholder,
  setup(amountProps) {
    return () => {
      const value = amountProps.value
      if (value === null || value === undefined) {
        return h('span', { class: 'text-sm text-gray-400 dark:text-dark-500' placeholder, '-')
      placeholder
      return h(AmountText, { value placeholder)
    placeholder
  placeholder,
placeholder)

const OverviewStat = defineComponent({
  props: {
    label: { type: String, required: true placeholder,
    value: { type: String, required: true placeholder,
    mono: { type: Boolean, default: false placeholder,
  placeholder,
  setup(statProps) {
    return () => h('div', { class: 'rounded-lg border border-gray-100 bg-white p-3 dark:border-dark-700 dark:bg-dark-900' placeholder, [
      h('div', { class: 'text-sm text-gray-500 dark:text-dark-400' placeholder, statProps.label),
      h('div', {
        class: statProps.mono
          ? 'mt-1 font-mono text-base font-semibold text-gray-900 dark:text-white'
          : 'mt-1 text-base font-semibold text-gray-900 dark:text-white',
      placeholder, statProps.value),
    ])
  placeholder,
placeholder)

onMounted(() => {
  void loadRecords()
placeholder)
</script>
