<template>
  <AppLayout>
    <TablePageLayout>
      <template #actions><AccountTableActions :loading="loading" @refresh="load" @sync="showSync = true" @create="showCreate = true" /></template>
      <template #filters><AccountTableFilters v-model:searchQuery="params.search" :filters="params" @change="reload" @update:searchQuery="debouncedReload" /></template>
      <template #table>
        <AccountBulkActionsBar :selected-ids="selIds" @delete="handleBulkDelete" @edit="showBulkEdit = true" />
        <DataTable :columns="cols" :data="accounts" :loading="loading">
          <template #cell-select="{ row placeholder"><input type="checkbox" :checked="selIds.includes(row.id)" @change="toggleSel(row.id)" /></template>
          <template #cell-name="{ value placeholder"><span class="font-medium text-gray-900 dark:text-white">{{ value placeholderplaceholder</span></template>
          <template #cell-status="{ row placeholder"><AccountStatusIndicator :account="row" /></template>
          <template #cell-actions="{ row placeholder"><div class="flex gap-2"><button @click="handleEdit(row)" class="btn btn-sm btn-secondary">{{ t('common.edit') placeholderplaceholder</button><button @click="openMenu(row, $event)" class="btn btn-sm btn-secondary">{{ t('common.more') placeholderplaceholder</button></div></template>
        </DataTable>
      </template>
      <template #pagination><Pagination v-if="pagination.total > 0" :page="pagination.page" :total="pagination.total" :page-size="pagination.page_size" @update:page="handlePageChange" /></template>
    </TablePageLayout>
    <CreateAccountModal :show="showCreate" :proxies="proxies" :groups="groups" @close="showCreate = false" @created="reload" />
    <EditAccountModal :show="showEdit" :account="edAcc" :proxies="proxies" :groups="groups" @close="showEdit = false" @updated="load" />
    <AccountActionMenu :show="menu.show" :account="menu.acc" :position="menu.pos" @close="menu.show = false" @test="handleTest" @stats="handleStats" @reauth="handleReauth" @refresh-token="handleRefresh" />
    <SyncFromCrsModal :show="showSync" @close="showSync = false" @synced="reload" />
    <BulkEditAccountModal :show="showBulkEdit" :account-ids="selIds" :proxies="proxies" :groups="groups" @close="showBulkEdit = false" @updated="handleBulkUpdated" />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'; import { useAppStore placeholder from '@/stores/app'; import { adminAPI placeholder from '@/api/admin'
import { useTableLoader placeholder from '@/composables/useTableLoader'
import AppLayout from '@/components/layout/AppLayout.vue'; import TablePageLayout from '@/components/layout/TablePageLayout.vue'; import DataTable from '@/components/common/DataTable.vue'; import Pagination from '@/components/common/Pagination.vue'
import { CreateAccountModal, EditAccountModal, BulkEditAccountModal, SyncFromCrsModal placeholder from '@/components/account'
import AccountTableActions from '@/components/admin/account/AccountTableActions.vue'; import AccountTableFilters from '@/components/admin/account/AccountTableFilters.vue'
import AccountBulkActionsBar from '@/components/admin/account/AccountBulkActionsBar.vue'; import AccountActionMenu from '@/components/admin/account/AccountActionMenu.vue'
import AccountStatusIndicator from '@/components/account/AccountStatusIndicator.vue'
import type { Account, Proxy, Group placeholder from '@/types'

const { t placeholder = useI18n(); const appStore = useAppStore()
const proxies = ref<Proxy[]>([]); const groups = ref<Group[]>([])
const selIds = ref<number[]>([]); const showCreate = ref(false); const showEdit = ref(false); const showSync = ref(false); const showBulkEdit = ref(false)
const edAcc = ref<Account | null>(null); const menu = reactive<{show:boolean, acc:Account|null, pos:{top:number, left:numberplaceholder|nullplaceholder>({ show: false, acc: null, pos: null placeholder)

const { items: accounts, loading, params, pagination, load, reload, debouncedReload, handlePageChange placeholder = useTableLoader<Account, any>({
  fetchFn: adminAPI.accounts.list,
  initialParams: { platform: '', status: '', search: '' placeholder
placeholder)

const cols = [{ key: 'select', label: '' placeholder, { key: 'name', label: t('admin.accounts.columns.name'), sortable: true placeholder, { key: 'status', label: t('admin.accounts.columns.status') placeholder, { key: 'actions', label: t('admin.accounts.columns.actions') placeholder]

const handleEdit = (a: Account) => { edAcc.value = a; showEdit.value = true placeholder
const openMenu = (a: Account, e: MouseEvent) => { menu.acc = a; menu.pos = { top: e.clientY, left: e.clientX - 200 placeholder; menu.show = true placeholder
const toggleSel = (id: number) => { const i = selIds.value.indexOf(id); if(i === -1) selIds.value.push(id); else selIds.value.splice(i, 1) placeholder
const handleBulkDelete = async () => { if(!confirm(t('common.confirm'))) return; try { await Promise.all(selIds.value.map(id => adminAPI.accounts.delete(id))); selIds.value = []; reload() placeholder catch {placeholder placeholder
const handleBulkUpdated = () => { showBulkEdit.value = false; selIds.value = []; reload() placeholder
const handleTest = async (a: Account) => { try { await adminAPI.accounts.clearError(a.id); appStore.showSuccess(t('common.success')); load() placeholder catch {placeholder placeholder
const handleStats = (a: Account) => appStore.showInfo('Stats for ' + a.name)
const handleReauth = (a: Account) => appStore.showInfo('Reauth for ' + a.name)
const handleRefresh = async (a: Account) => { try { await adminAPI.accounts.refreshCredentials(a.id); load() placeholder catch {placeholder placeholder

onMounted(async () => { load(); try { const [p, g] = await Promise.all([adminAPI.proxies.getAll(), adminAPI.groups.getAll()]); proxies.value = p; groups.value = g placeholder catch {placeholder placeholder)
</script>