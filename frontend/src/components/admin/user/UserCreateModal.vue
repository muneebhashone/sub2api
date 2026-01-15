<template>
  <BaseDialog
    :show="show"
    :title="t('admin.users.createUser')"
    width="normal"
    @close="$emit('close')"
  >
    <form id="create-user-form" @submit.prevent="submit" class="space-y-5">
      <div>
        <label class="input-label">{{ t('admin.users.email') placeholderplaceholder</label>
        <input v-model="form.email" type="email" required class="input" :placeholder="t('admin.users.enterEmail')" />
      </div>
      <div>
        <label class="input-label">{{ t('admin.users.password') placeholderplaceholder</label>
        <div class="flex gap-2">
          <div class="relative flex-1">
            <input v-model="form.password" type="text" required class="input pr-10" :placeholder="t('admin.users.enterPassword')" />
          </div>
          <button type="button" @click="generateRandomPassword" class="btn btn-secondary px-3">
            <Icon name="refresh" size="md" />
          </button>
        </div>
      </div>
      <div>
        <label class="input-label">{{ t('admin.users.username') placeholderplaceholder</label>
        <input v-model="form.username" type="text" class="input" :placeholder="t('admin.users.enterUsername')" />
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label class="input-label">{{ t('admin.users.columns.balance') placeholderplaceholder</label>
          <input v-model.number="form.balance" type="number" step="any" class="input" />
        </div>
        <div>
          <label class="input-label">{{ t('admin.users.columns.concurrency') placeholderplaceholder</label>
          <input v-model.number="form.concurrency" type="number" class="input" />
        </div>
      </div>
    </form>
    <template #footer>
      <div class="flex justify-end gap-3">
        <button @click="$emit('close')" type="button" class="btn btn-secondary">{{ t('common.cancel') placeholderplaceholder</button>
        <button type="submit" form="create-user-form" :disabled="loading" class="btn btn-primary">
          {{ loading ? t('admin.users.creating') : t('common.create') placeholderplaceholder
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { reactive, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'; import { adminAPI placeholder from '@/api/admin'
import { useForm placeholder from '@/composables/useForm'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'

const props = defineProps<{ show: boolean placeholder>()
const emit = defineEmits(['close', 'success']); const { t placeholder = useI18n()

const form = reactive({ email: '', password: '', username: '', notes: '', balance: 0, concurrency: 1 placeholder)

const { loading, submit placeholder = useForm({
  form,
  submitFn: async (data) => {
    await adminAPI.users.create(data)
    emit('success'); emit('close')
  placeholder,
  successMsg: t('admin.users.userCreated')
placeholder)

watch(() => props.show, (v) => { if(v) Object.assign(form, { email: '', password: '', username: '', notes: '', balance: 0, concurrency: 1 placeholder) placeholder)

const generateRandomPassword = () => {
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789!@#$%^&*'
  let p = ''; for (let i = 0; i < 16; i++) p += chars.charAt(Math.floor(Math.random() * chars.length))
  form.password = p
placeholder
</script>
