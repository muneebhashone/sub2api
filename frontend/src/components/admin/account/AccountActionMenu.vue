<template>
  <Teleport to="body">
    <div v-if="show && position" class="action-menu-content fixed z-[9999] w-52 overflow-hidden rounded-xl bg-white shadow-lg ring-1 ring-black/5 dark:bg-dark-800" :style="{ top: position.top + 'px', left: position.left + 'px' placeholder">
      <div class="py-1">
        <template v-if="account">
          <button @click="$emit('test', account); $emit('close')" class="flex w-full items-center gap-2 px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-dark-700">
            <Icon name="play" size="sm" class="text-green-500" :stroke-width="2" />
            {{ t('admin.accounts.testConnection') placeholderplaceholder
          </button>
          <button @click="$emit('stats', account); $emit('close')" class="flex w-full items-center gap-2 px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-dark-700">
            <Icon name="chart" size="sm" class="text-indigo-500" />
            {{ t('admin.accounts.viewStats') placeholderplaceholder
          </button>
          <template v-if="account.type === 'oauth' || account.type === 'setup-token'">
            <button @click="$emit('reauth', account); $emit('close')" class="flex w-full items-center gap-2 px-4 py-2 text-sm text-blue-600 hover:bg-gray-100 dark:hover:bg-dark-700">
              <Icon name="link" size="sm" />
              {{ t('admin.accounts.reAuthorize') placeholderplaceholder
            </button>
            <button @click="$emit('refresh-token', account); $emit('close')" class="flex w-full items-center gap-2 px-4 py-2 text-sm text-purple-600 hover:bg-gray-100 dark:hover:bg-dark-700">
              <Icon name="refresh" size="sm" />
              {{ t('admin.accounts.refreshToken') placeholderplaceholder
            </button>
          </template>
          <div v-if="account.status === 'error' || isRateLimited || isOverloaded" class="my-1 border-t border-gray-100 dark:border-dark-700"></div>
          <button v-if="account.status === 'error'" @click="$emit('reset-status', account); $emit('close')" class="flex w-full items-center gap-2 px-4 py-2 text-sm text-yellow-600 hover:bg-gray-100 dark:hover:bg-dark-700">
            <Icon name="sync" size="sm" />
            {{ t('admin.accounts.resetStatus') placeholderplaceholder
          </button>
          <button v-if="isRateLimited || isOverloaded" @click="$emit('clear-rate-limit', account); $emit('close')" class="flex w-full items-center gap-2 px-4 py-2 text-sm text-amber-600 hover:bg-gray-100 dark:hover:bg-dark-700">
            <Icon name="clock" size="sm" />
            {{ t('admin.accounts.clearRateLimit') placeholderplaceholder
          </button>
        </template>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { Icon placeholder from '@/components/icons'
import type { Account placeholder from '@/types'

const props = defineProps<{ show: boolean; account: Account | null; position: { top: number; left: number placeholder | null placeholder>()
defineEmits(['close', 'test', 'stats', 'reauth', 'refresh-token', 'reset-status', 'clear-rate-limit'])
const { t placeholder = useI18n()
const isRateLimited = computed(() => props.account?.rate_limit_reset_at && new Date(props.account.rate_limit_reset_at) > new Date())
const isOverloaded = computed(() => props.account?.overload_until && new Date(props.account.overload_until) > new Date())
</script>
