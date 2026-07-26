<template>
  <div class="card">
    <div class="flex items-start justify-between border-b border-gray-100 px-6 py-4 dark:border-dark-700">
      <div>
        <h2 class="text-lg font-medium text-gray-900 dark:text-white">
          {{ t('profile.passkey.title') placeholderplaceholder
        </h2>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          {{ t('profile.passkey.description') placeholderplaceholder
        </p>
      </div>
      <button
        v-if="enabled && supported && !showAddForm"
        type="button"
        class="btn btn-primary"
        :disabled="busy"
        @click="showAddForm = true"
      >
        {{ t('profile.passkey.add') placeholderplaceholder
      </button>
    </div>

    <div class="px-6 py-6">
      <div v-if="!enabled" class="text-sm text-gray-500 dark:text-gray-400">
        {{ t('profile.passkey.featureDisabled') placeholderplaceholder
      </div>
      <div v-else-if="!supported" class="text-sm text-amber-600 dark:text-amber-400">
        {{ t('profile.passkey.unsupported') placeholderplaceholder
      </div>
      <div v-else>
        <form
          v-if="showAddForm"
          class="mb-5 flex flex-col gap-3 rounded-lg border border-gray-200 p-4 dark:border-dark-700 sm:flex-row sm:items-end"
          @submit.prevent="addPasskey"
        >
          <div class="flex-1">
            <label for="passkey-name" class="input-label">{{ t('profile.passkey.name') placeholderplaceholder</label>
            <input
              id="passkey-name"
              v-model="newName"
              class="input"
              maxlength="100"
              :placeholder="t('profile.passkey.namePlaceholder')"
              autofocus
            />
          </div>
          <div class="flex gap-2">
            <button type="button" class="btn btn-secondary" :disabled="busy" @click="cancelAdd">
              {{ t('common.cancel') placeholderplaceholder
            </button>
            <button type="submit" class="btn btn-primary" :disabled="busy">
              {{ busy ? t('common.processing') : t('profile.passkey.continue') placeholderplaceholder
            </button>
          </div>
        </form>

        <div v-if="loading" class="flex justify-center py-6">
          <div class="h-8 w-8 animate-spin rounded-full border-b-2 border-primary-500"></div>
        </div>

        <div
          v-else-if="credentials.length === 0"
          class="rounded-lg border border-dashed border-gray-200 px-4 py-8 text-center text-sm text-gray-500 dark:border-dark-700 dark:text-gray-400"
        >
          {{ t('profile.passkey.empty') placeholderplaceholder
        </div>

        <div v-else class="divide-y divide-gray-100 dark:divide-dark-700">
          <div
            v-for="credential in credentials"
            :key="credential.id"
            class="flex items-center justify-between gap-4 py-4 first:pt-0 last:pb-0"
          >
            <div class="min-w-0">
              <div class="flex items-center gap-2">
                <Icon name="key" size="md" class="shrink-0 text-primary-500" />
                <p class="truncate font-medium text-gray-900 dark:text-white">
                  {{ credential.name placeholderplaceholder
                </p>
                <span
                  v-if="credential.backup"
                  class="rounded-full bg-green-50 px-2 py-0.5 text-xs text-green-700 dark:bg-green-900/30 dark:text-green-300"
                >
                  {{ t('profile.passkey.synced') placeholderplaceholder
                </span>
              </div>
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                {{ t('profile.passkey.createdAt', { date: formatDate(credential.created_at) placeholder) placeholderplaceholder
                <template v-if="credential.last_used_at">
                  · {{ t('profile.passkey.lastUsed', { date: formatDate(credential.last_used_at) placeholder) placeholderplaceholder
                </template>
              </p>
            </div>
            <div class="flex shrink-0 gap-2">
              <button
                type="button"
                class="btn btn-secondary btn-sm"
                :disabled="busy"
                @click="renamePasskey(credential)"
              >
                {{ t('common.edit') placeholderplaceholder
              </button>
              <button
                type="button"
                class="btn btn-ghost btn-sm text-red-600 hover:bg-red-50 dark:text-red-300 dark:hover:bg-red-950/30"
                :disabled="busy"
                @click="deletePasskey(credential)"
              >
                {{ t('common.delete') placeholderplaceholder
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { passkeyAPI, type PasskeyCredentialSummary placeholder from '@/api'
import { Icon placeholder from '@/components/icons'
import { useAppStore placeholder from '@/stores/app'

const props = defineProps<{ enabled: boolean placeholder>()

const { t placeholder = useI18n()
const appStore = useAppStore()
const supported = passkeyAPI.isSupported()
const loading = ref(false)
const busy = ref(false)
const showAddForm = ref(false)
const newName = ref('')
const credentials = ref<PasskeyCredentialSummary[]>([])

async function loadCredentials(): Promise<void> {
  if (!supported) return
  loading.value = true
  try {
    credentials.value = await passkeyAPI.list()
  placeholder catch (error) {
    const code = (error as { code?: string placeholder).code
    if (code !== 'PASSKEY_DISABLED') {
      appStore.showError(t('profile.passkey.loadFailed'))
    placeholder
  placeholder finally {
    loading.value = false
  placeholder
placeholder

async function addPasskey(): Promise<void> {
  busy.value = true
  try {
    await passkeyAPI.register(newName.value.trim())
    appStore.showSuccess(t('profile.passkey.added'))
    cancelAdd()
    await loadCredentials()
  placeholder catch (error) {
    if (!(error instanceof DOMException && error.name === 'NotAllowedError')) {
      appStore.showError(t('profile.passkey.addFailed'))
    placeholder
  placeholder finally {
    busy.value = false
  placeholder
placeholder

function cancelAdd(): void {
  showAddForm.value = false
  newName.value = ''
placeholder

async function renamePasskey(credential: PasskeyCredentialSummary): Promise<void> {
  const name = window.prompt(t('profile.passkey.renamePrompt'), credential.name)?.trim()
  if (!name || name === credential.name) return
  busy.value = true
  try {
    await passkeyAPI.rename(credential.id, name)
    credential.name = name
    appStore.showSuccess(t('profile.passkey.renamed'))
  placeholder catch {
    appStore.showError(t('profile.passkey.renameFailed'))
  placeholder finally {
    busy.value = false
  placeholder
placeholder

async function deletePasskey(credential: PasskeyCredentialSummary): Promise<void> {
  if (!window.confirm(t('profile.passkey.deleteConfirm', { name: credential.name placeholder))) return
  busy.value = true
  try {
    await passkeyAPI.remove(credential.id)
    credentials.value = credentials.value.filter((item) => item.id !== credential.id)
    appStore.showSuccess(t('profile.passkey.deleted'))
  placeholder catch {
    appStore.showError(t('profile.passkey.deleteFailed'))
  placeholder finally {
    busy.value = false
  placeholder
placeholder

function formatDate(value: string): string {
  return new Intl.DateTimeFormat(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  placeholder).format(new Date(value))
placeholder

watch(
  () => props.enabled,
  (enabled) => {
    if (enabled) void loadCredentials()
  placeholder,
  { immediate: true placeholder
)
</script>
