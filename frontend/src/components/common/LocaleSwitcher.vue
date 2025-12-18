<template>
  <div class="relative" ref="dropdownRef">
    <button
      @click="toggleDropdown"
      class="flex items-center gap-1.5 px-2 py-1.5 rounded-lg text-sm font-medium text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-dark-700 transition-colors"
      :title="currentLocale?.name"
    >
      <span class="text-base">{{ currentLocale?.flag placeholderplaceholder</span>
      <span class="hidden sm:inline">{{ currentLocale?.code.toUpperCase() placeholderplaceholder</span>
      <svg
        class="w-3.5 h-3.5 text-gray-400 transition-transform duration-200"
        :class="{ 'rotate-180': isOpen placeholder"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        stroke-width="2"
      >
        <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
      </svg>
    </button>

    <transition name="dropdown">
      <div
        v-if="isOpen"
        class="absolute right-0 mt-1 w-32 rounded-lg bg-white dark:bg-dark-800 shadow-lg border border-gray-200 dark:border-dark-700 overflow-hidden z-50"
      >
        <button
          v-for="locale in availableLocales"
          :key="locale.code"
          @click="selectLocale(locale.code)"
          class="w-full flex items-center gap-2 px-3 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-dark-700 transition-colors"
          :class="{ 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400': locale.code === currentLocaleCode placeholder"
        >
          <span class="text-base">{{ locale.flag placeholderplaceholder</span>
          <span>{{ locale.name placeholderplaceholder</span>
          <svg
            v-if="locale.code === currentLocaleCode"
            class="w-4 h-4 ml-auto text-primary-500"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="2"
          >
            <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
          </svg>
        </button>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { setLocale, availableLocales placeholder from '@/i18n'

const { locale placeholder = useI18n()

const isOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)

const currentLocaleCode = computed(() => locale.value)
const currentLocale = computed(() => availableLocales.find(l => l.code === locale.value))

function toggleDropdown() {
  isOpen.value = !isOpen.value
placeholder

function selectLocale(code: string) {
  setLocale(code)
  isOpen.value = false
placeholder

function handleClickOutside(event: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    isOpen.value = false
  placeholder
placeholder

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
placeholder)

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
placeholder)
</script>

<style scoped>
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.15s ease;
placeholder

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: scale(0.95) translateY(-4px);
placeholder
</style>
