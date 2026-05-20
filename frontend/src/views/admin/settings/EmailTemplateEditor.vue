<template>
  <div class="card">
    <div
      class="flex flex-col gap-3 border-b border-gray-100 px-6 py-4 dark:border-dark-700 lg:flex-row lg:items-start lg:justify-between"
    >
      <div>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ t("admin.settings.emailTemplates.title") placeholderplaceholder
        </h2>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          {{ t("admin.settings.emailTemplates.description") placeholderplaceholder
        </p>
      </div>
      <div class="flex flex-wrap gap-2">
        <button
          type="button"
          class="btn btn-secondary btn-sm"
          :disabled="loadingTemplate || previewing || !canPreview"
          @click="refreshPreview"
        >
          {{ previewing ? t("admin.settings.emailTemplates.previewing") : t("admin.settings.emailTemplates.preview") placeholderplaceholder
        </button>
        <button
          type="button"
          class="btn btn-secondary btn-sm"
          :disabled="loadingTemplate || restoring || !selectedEvent || !selectedLocale"
          @click="restoreOfficial"
        >
          {{ restoring ? t("admin.settings.emailTemplates.restoring") : t("admin.settings.emailTemplates.restoreOfficial") placeholderplaceholder
        </button>
        <button
          type="button"
          class="btn btn-primary btn-sm"
          :disabled="loadingTemplate || saving || !canSave"
          @click="saveTemplate"
        >
          {{ saving ? t("admin.settings.emailTemplates.saving") : t("admin.settings.emailTemplates.save") placeholderplaceholder
        </button>
      </div>
    </div>

    <div class="space-y-6 p-6">
      <div
        v-if="loadingList"
        class="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400"
      >
        <span
          class="h-4 w-4 animate-spin rounded-full border-b-2 border-primary-600"
        ></span>
        {{ t("common.loading") placeholderplaceholder
      </div>

      <template v-else>
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="input-label" for="email-template-event">
              {{ t("admin.settings.emailTemplates.event") placeholderplaceholder
            </label>
            <select
              id="email-template-event"
              v-model="selectedEvent"
              class="input"
              :disabled="loadingTemplate || eventOptions.length === 0"
            >
              <option
                v-for="option in eventOptions"
                :key="option.value"
                :value="option.value"
              >
                {{ option.label || option.value placeholderplaceholder
              </option>
            </select>
            <p v-if="selectedEventDescription" class="input-hint">
              {{ selectedEventDescription placeholderplaceholder
            </p>
          </div>
          <div>
            <label class="input-label" for="email-template-locale">
              {{ t("admin.settings.emailTemplates.locale") placeholderplaceholder
            </label>
            <select
              id="email-template-locale"
              v-model="selectedLocale"
              class="input"
              :disabled="loadingTemplate || localeOptions.length === 0"
            >
              <option
                v-for="localeOption in localeOptions"
                :key="localeOption"
                :value="localeOption"
              >
                {{ formatLocale(localeOption) placeholderplaceholder
              </option>
            </select>
          </div>
        </div>

        <div
          v-if="!eventOptions.length || !localeOptions.length"
          class="rounded-lg border border-amber-200 bg-amber-50 p-4 text-sm text-amber-700 dark:border-amber-800 dark:bg-amber-900/20 dark:text-amber-300"
        >
          {{ t("admin.settings.emailTemplates.empty") placeholderplaceholder
        </div>

        <div v-else class="grid grid-cols-1 gap-6 xl:grid-cols-2">
          <div class="space-y-4">
            <div>
              <label class="input-label" for="email-template-subject">
                {{ t("admin.settings.emailTemplates.subject") placeholderplaceholder
              </label>
              <input
                id="email-template-subject"
                v-model="subject"
                type="text"
                class="input"
                :disabled="loadingTemplate"
                :placeholder="t('admin.settings.emailTemplates.subjectPlaceholder')"
              />
            </div>

            <div>
              <label class="input-label" for="email-template-html">
                {{ t("admin.settings.emailTemplates.html") placeholderplaceholder
              </label>
              <textarea
                id="email-template-html"
                v-model="html"
                rows="18"
                class="input min-h-[28rem] resize-y font-mono text-sm leading-6"
                :disabled="loadingTemplate"
                :placeholder="t('admin.settings.emailTemplates.htmlPlaceholder')"
              ></textarea>
            </div>

            <div
              class="rounded-lg border border-gray-200 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800/60"
            >
              <div class="text-sm font-medium text-gray-900 dark:text-white">
                {{ t("admin.settings.emailTemplates.placeholders") placeholderplaceholder
              </div>
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                {{ t("admin.settings.emailTemplates.placeholdersHelp") placeholderplaceholder
              </p>
              <div class="mt-3 flex flex-wrap gap-2">
                <button
                  v-for="placeholder in placeholderList"
                  :key="placeholder"
                  type="button"
                  class="rounded-full border border-gray-200 bg-white px-3 py-1 font-mono text-xs text-gray-700 transition-colors hover:border-primary-300 hover:text-primary-600 dark:border-dark-600 dark:bg-dark-700 dark:text-gray-200 dark:hover:border-primary-500 dark:hover:text-primary-300"
                  @click="copyPlaceholder(placeholder)"
                >
                  {{ placeholder placeholderplaceholder
                </button>
              </div>
            </div>
          </div>

          <div class="space-y-4">
            <div
              class="rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800"
            >
              <div
                class="flex items-center justify-between border-b border-gray-100 px-4 py-3 dark:border-dark-700"
              >
                <div>
                  <div class="text-sm font-medium text-gray-900 dark:text-white">
                    {{ t("admin.settings.emailTemplates.livePreview") placeholderplaceholder
                  </div>
                  <div class="mt-0.5 text-xs text-gray-500 dark:text-gray-400">
                    {{ previewSubject || t("admin.settings.emailTemplates.noPreview") placeholderplaceholder
                  </div>
                </div>
                <span
                  v-if="isCustomTemplate"
                  class="rounded-full bg-primary-50 px-2.5 py-1 text-xs font-medium text-primary-700 dark:bg-primary-900/30 dark:text-primary-300"
                >
                  {{ t("admin.settings.emailTemplates.customized") placeholderplaceholder
                </span>
              </div>
              <div class="bg-gray-100 p-3 dark:bg-dark-900">
                <iframe
                  class="h-[36rem] w-full rounded-md border border-gray-200 bg-white dark:border-dark-700"
                  sandbox=""
                  :srcdoc="previewHtml"
                  :title="t('admin.settings.emailTemplates.livePreview')"
                ></iframe>
              </div>
            </div>

            <p class="text-xs text-gray-500 dark:text-gray-400">
              {{ t("admin.settings.emailTemplates.previewSecurityHint") placeholderplaceholder
            </p>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch placeholder from "vue";
import { useI18n placeholder from "vue-i18n";
import { adminAPI placeholder from "@/api";
import type {
  EmailTemplateEventOption,
  EmailTemplateOption,
placeholder from "@/api/admin/settings";
import { useAppStore placeholder from "@/stores";
import { extractApiErrorMessage placeholder from "@/utils/apiError";

const { t placeholder = useI18n();
const appStore = useAppStore();

const fallbackPlaceholders = [
  "{{site_nameplaceholderplaceholder",
  "{{recipient_nameplaceholderplaceholder",
  "{{recipient_emailplaceholderplaceholder",
  "{{subscription_groupplaceholderplaceholder",
  "{{expiry_timeplaceholderplaceholder",
  "{{days_remainingplaceholderplaceholder",
  "{{current_balanceplaceholderplaceholder",
  "{{thresholdplaceholderplaceholder",
  "{{recharge_urlplaceholderplaceholder",
  "{{recharge_amountplaceholderplaceholder",
  "{{order_idplaceholderplaceholder",
  "{{unsubscribe_urlplaceholderplaceholder",
];

const loadingList = ref(true);
const loadingTemplate = ref(false);
const saving = ref(false);
const previewing = ref(false);
const restoring = ref(false);
const eventOptions = ref<EmailTemplateOption[]>([]);
const localeOptions = ref<string[]>([]);
const selectedEvent = ref("");
const selectedLocale = ref("");
const subject = ref("");
const html = ref("");
const isCustomTemplate = ref(false);
const placeholders = ref<string[]>([]);
const previewSubject = ref("");
const previewHtml = ref("");
const initializingSelection = ref(false);

function normalizeEventOption(option: EmailTemplateEventOption): EmailTemplateOption {
  if (typeof option === "string") {
    return { value: option placeholder;
  placeholder
  return option;
placeholder

const selectedEventDescription = computed(() => {
  return (
    eventOptions.value.find((option) => option.value === selectedEvent.value)
      ?.description || ""
  );
placeholder);

const placeholderList = computed(() => {
  const combined = [...placeholders.value, ...fallbackPlaceholders];
  return Array.from(
    new Set(
      combined
        .map((item) => formatPlaceholder(item))
        .filter((item) => item.length > 0),
    ),
  );
placeholder);

function formatPlaceholder(placeholder: string): string {
  const trimmed = placeholder.trim();
  if (!trimmed) return "";
  if (trimmed.startsWith("{{") && trimmed.endsWith("placeholderplaceholder")) return trimmed;
  return `{{${trimmedplaceholderplaceholderplaceholder`;
placeholder

const canSave = computed(
  () =>
    Boolean(selectedEvent.value && selectedLocale.value) &&
    subject.value.trim().length > 0 &&
    html.value.trim().length > 0,
);

const canPreview = computed(
  () => Boolean(selectedEvent.value && selectedLocale.value) && html.value.trim().length > 0,
);

function formatLocale(locale: string): string {
  const lower = locale.toLowerCase();
  if (lower === "zh" || lower.startsWith("zh-")) {
    return t("admin.settings.emailTemplates.localeZh");
  placeholder
  if (lower === "en" || lower.startsWith("en-")) {
    return t("admin.settings.emailTemplates.localeEn");
  placeholder
  return locale;
placeholder

function applyTemplate(template: {
  subject: string;
  html: string;
  is_custom?: boolean;
  placeholders?: string[];
placeholder) {
  subject.value = template.subject;
  html.value = template.html;
  isCustomTemplate.value = template.is_custom === true;
  if (template.placeholders?.length) {
    placeholders.value = template.placeholders;
  placeholder
placeholder

async function loadTemplate() {
  if (!selectedEvent.value || !selectedLocale.value) return;
  loadingTemplate.value = true;
  try {
    const template = await adminAPI.settings.getEmailTemplate(
      selectedEvent.value,
      selectedLocale.value,
    );
    applyTemplate(template);
    await refreshPreview();
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t("common.error")));
  placeholder finally {
    loadingTemplate.value = false;
  placeholder
placeholder

async function loadTemplateList() {
  loadingList.value = true;
  try {
    const response = await adminAPI.settings.getEmailTemplates();
    eventOptions.value = response.events.map(normalizeEventOption);
    localeOptions.value = response.locales;
    placeholders.value = response.placeholders || [];
    initializingSelection.value = true;
    selectedEvent.value = eventOptions.value[0]?.value || "";
    selectedLocale.value = response.locales[0] || "";
    await loadTemplate();
    initializingSelection.value = false;
  placeholder catch (err: unknown) {
    initializingSelection.value = false;
    appStore.showError(extractApiErrorMessage(err, t("common.error")));
  placeholder finally {
    loadingList.value = false;
  placeholder
placeholder

async function saveTemplate() {
  if (!canSave.value) {
    appStore.showError(t("admin.settings.emailTemplates.validationRequired"));
    return;
  placeholder
  saving.value = true;
  try {
    const template = await adminAPI.settings.updateEmailTemplate(
      selectedEvent.value,
      selectedLocale.value,
      {
        subject: subject.value,
        html: html.value,
      placeholder,
    );
    applyTemplate(template);
    await refreshPreview();
    appStore.showSuccess(t("admin.settings.emailTemplates.saveSuccess"));
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t("common.error")));
  placeholder finally {
    saving.value = false;
  placeholder
placeholder

async function refreshPreview() {
  if (!canPreview.value) {
    previewSubject.value = "";
    previewHtml.value = "";
    return;
  placeholder
  previewing.value = true;
  try {
    const preview = await adminAPI.settings.previewEmailTemplate({
      event: selectedEvent.value,
      locale: selectedLocale.value,
      subject: subject.value,
      html: html.value,
    placeholder);
    previewSubject.value = preview.subject;
    previewHtml.value = preview.html;
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t("common.error")));
  placeholder finally {
    previewing.value = false;
  placeholder
placeholder

async function restoreOfficial() {
  if (!selectedEvent.value || !selectedLocale.value) return;
  if (!window.confirm(t("admin.settings.emailTemplates.restoreConfirm"))) return;

  restoring.value = true;
  try {
    const template = await adminAPI.settings.restoreOfficialEmailTemplate(
      selectedEvent.value,
      selectedLocale.value,
    );
    applyTemplate(template);
    await refreshPreview();
    appStore.showSuccess(t("admin.settings.emailTemplates.restoreSuccess"));
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t("common.error")));
  placeholder finally {
    restoring.value = false;
  placeholder
placeholder

async function copyPlaceholder(placeholder: string) {
  try {
    await navigator.clipboard.writeText(placeholder);
    appStore.showSuccess(t("admin.settings.emailTemplates.placeholderCopied"));
  placeholder catch {
    appStore.showError(t("common.error"));
  placeholder
placeholder

watch([selectedEvent, selectedLocale], ([eventValue, localeValue], [oldEvent, oldLocale]) => {
  if (initializingSelection.value) return;
  if (!eventValue || !localeValue) return;
  if (eventValue === oldEvent && localeValue === oldLocale) return;
  void loadTemplate();
placeholder);

onMounted(() => {
  void loadTemplateList();
placeholder);
</script>
