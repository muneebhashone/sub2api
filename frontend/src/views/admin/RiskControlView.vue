<template>
  <AppLayout>
    <div class="space-y-6">
      <div v-if="loading" class="flex items-center justify-center py-16">
        <div class="h-8 w-8 animate-spin rounded-full border-b-2 border-primary-600"></div>
      </div>

      <template v-else>
        <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">{{ t('admin.riskControl.title') placeholderplaceholder</h1>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.description') placeholderplaceholder</p>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <button type="button" class="btn btn-secondary inline-flex items-center gap-2" :disabled="statusLoading" @click="loadStatus(false)">
              <Icon name="refresh" size="sm" :class="statusLoading ? 'animate-spin' : ''" />
              {{ t('admin.riskControl.refreshStatus') placeholderplaceholder
            </button>
            <button type="button" class="btn btn-primary inline-flex items-center gap-2" @click="openSettings">
              <Icon name="cog" size="sm" />
              {{ t('admin.riskControl.openSettings') placeholderplaceholder
            </button>
          </div>
        </div>

        <div class="grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-4">
          <div
            v-for="item in overviewItems"
            :key="item.key"
            class="rounded-lg border border-gray-100 bg-white px-4 py-3 shadow-sm dark:border-dark-700 dark:bg-dark-800"
          >
            <div class="flex min-w-0 items-center gap-3">
              <div class="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-lg" :class="item.iconClass">
                <Icon :name="item.icon" size="sm" />
              </div>
              <div class="min-w-0 flex-1">
                <div class="flex min-w-0 items-center justify-between gap-2">
                  <p class="truncate text-xs font-medium text-gray-500 dark:text-gray-400">{{ item.label placeholderplaceholder</p>
                  <span
                    v-if="item.badge"
                    class="inline-flex flex-shrink-0 items-center rounded-full px-2 py-0.5 text-xs font-medium"
                    :class="item.badgeClass"
                  >
                    {{ item.badge placeholderplaceholder
                  </span>
                </div>
                <div class="mt-1 flex min-w-0 items-baseline gap-2">
                  <p class="truncate text-xl font-semibold leading-7 text-gray-900 dark:text-white">{{ item.value placeholderplaceholder</p>
                  <p v-if="item.meta" class="truncate text-xs text-gray-500 dark:text-gray-400">{{ item.meta placeholderplaceholder</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="card">
          <div class="flex flex-col gap-4 border-b border-gray-100 px-6 py-4 dark:border-dark-700 lg:flex-row lg:items-center lg:justify-between">
            <div>
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('admin.riskControl.workerStatus') placeholderplaceholder</h2>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.workerStatusHint') placeholderplaceholder</p>
            </div>
            <div class="flex flex-wrap items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
              <span>{{ t('admin.riskControl.autoRefresh') placeholderplaceholder</span>
              <span v-if="status?.last_cleanup_at">
                {{ t('admin.riskControl.lastCleanup', { time: formatDateTime(status.last_cleanup_at) placeholder) placeholderplaceholder
              </span>
            </div>
          </div>

          <div class="grid grid-cols-1 gap-6 p-6 xl:grid-cols-[minmax(0,360px)_1fr]">
            <div class="space-y-4">
              <div class="rounded-lg border border-gray-100 p-4 dark:border-dark-700">
                <div class="flex items-center justify-between gap-3">
                  <div>
                    <p class="text-sm font-medium text-gray-900 dark:text-white">{{ t('admin.riskControl.queueUsage') placeholderplaceholder</p>
                    <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                      {{ formatNumber(status?.queue_length ?? 0) placeholderplaceholder / {{ formatNumber(status?.queue_size ?? configForm.queue_size) placeholderplaceholder
                    </p>
                  </div>
                  <span class="text-sm font-semibold text-gray-900 dark:text-white">{{ queueUsagePercent placeholderplaceholder</span>
                </div>
                <div class="mt-4 h-2 overflow-hidden rounded-full bg-gray-100 dark:bg-dark-700">
                  <div class="h-full rounded-full bg-primary-500 transition-all duration-300" :style="queueUsageStyle"></div>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-3">
                <div class="rounded-lg bg-gray-50 p-4 dark:bg-dark-700/50">
                  <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.activeWorkers') placeholderplaceholder</p>
                  <p class="mt-2 text-2xl font-semibold text-gray-900 dark:text-white">{{ status?.active_workers ?? 0 placeholderplaceholder</p>
                </div>
                <div class="rounded-lg bg-emerald-50 p-4 dark:bg-emerald-900/10">
                  <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.idleWorkers') placeholderplaceholder</p>
                  <p class="mt-2 text-2xl font-semibold text-emerald-700 dark:text-emerald-300">{{ status?.idle_workers ?? configForm.worker_count placeholderplaceholder</p>
                </div>
                <div class="rounded-lg bg-gray-50 p-4 dark:bg-dark-700/50">
                  <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.processed') placeholderplaceholder</p>
                  <p class="mt-2 text-2xl font-semibold text-gray-900 dark:text-white">{{ formatNumber(status?.processed ?? 0) placeholderplaceholder</p>
                </div>
                <div class="rounded-lg bg-gray-50 p-4 dark:bg-dark-700/50">
                  <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.droppedErrors') placeholderplaceholder</p>
                  <p class="mt-2 text-2xl font-semibold text-gray-900 dark:text-white">{{ formatNumber((status?.dropped ?? 0) + (status?.errors ?? 0)) placeholderplaceholder</p>
                </div>
              </div>
            </div>

            <div>
              <div class="mb-3 flex items-center justify-between gap-3">
                <div>
                  <p class="text-sm font-medium text-gray-900 dark:text-white">{{ t('admin.riskControl.workerPool') placeholderplaceholder</p>
                  <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                    {{ t('admin.riskControl.workerPoolMeta', { active: status?.active_workers ?? 0, idle: status?.idle_workers ?? configForm.worker_count, total: status?.worker_count ?? configForm.worker_count placeholder) placeholderplaceholder
                  </p>
                </div>
                <span class="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-1 text-xs font-medium text-gray-600 dark:bg-dark-700 dark:text-gray-300">
                  {{ modeLabel(status?.mode ?? configForm.mode) placeholderplaceholder
                </span>
              </div>
              <div class="grid grid-cols-2 gap-2 sm:grid-cols-4 md:grid-cols-6 xl:grid-cols-8 2xl:grid-cols-10">
                <div
                  v-for="worker in workerSlots"
                  :key="worker.id"
                  class="flex h-12 items-center justify-between rounded-lg border px-3 transition-colors"
                  :class="workerSlotClass(worker.state)"
                  :title="worker.label"
                >
                  <span class="text-sm font-semibold">#{{ worker.id placeholderplaceholder</span>
                  <span class="h-2.5 w-2.5 rounded-full" :class="workerDotClass(worker.state)"></span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="card">
          <div class="flex flex-col gap-4 border-b border-gray-100 px-6 py-4 dark:border-dark-700">
            <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
              <div>
                <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('admin.riskControl.records') placeholderplaceholder</h2>
                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.recordsHint') placeholderplaceholder</p>
              </div>
              <button type="button" class="btn btn-secondary inline-flex items-center gap-2" :disabled="logsLoading" @click="loadLogs">
                <Icon name="refresh" size="sm" :class="logsLoading ? 'animate-spin' : ''" />
                {{ t('admin.riskControl.refresh') placeholderplaceholder
              </button>
            </div>

            <div class="grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-6">
              <Select v-model="filters.result" :options="resultOptions" @change="reloadLogsFromFirstPage" />
              <Select v-model="filters.group_id" :options="groupFilterOptions" @change="reloadLogsFromFirstPage" />
              <Select v-model="filters.endpoint" :options="endpointOptions" @change="reloadLogsFromFirstPage" />
              <input v-model.trim="filters.search" type="search" class="input" :placeholder="t('admin.riskControl.filters.search')" @keyup.enter="reloadLogsFromFirstPage" />
              <input v-model="filters.from" type="datetime-local" class="input" :title="t('admin.riskControl.filters.from')" @change="reloadLogsFromFirstPage" />
              <input v-model="filters.to" type="datetime-local" class="input" :title="t('admin.riskControl.filters.to')" @change="reloadLogsFromFirstPage" />
            </div>
          </div>

          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200 dark:divide-dark-700">
              <thead class="bg-gray-50 dark:bg-dark-800">
                <tr>
                  <th class="px-5 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.table.time') placeholderplaceholder</th>
                  <th class="px-5 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.table.group') placeholderplaceholder</th>
                  <th class="px-5 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.table.user') placeholderplaceholder</th>
                  <th class="px-5 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.table.apiKey') placeholderplaceholder</th>
                  <th class="px-5 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.table.endpoint') placeholderplaceholder</th>
                  <th class="px-5 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.table.result') placeholderplaceholder</th>
                  <th class="px-5 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.table.highest') placeholderplaceholder</th>
                  <th class="px-5 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.table.actionMeta') placeholderplaceholder</th>
                  <th class="px-5 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.table.latency') placeholderplaceholder</th>
                  <th class="px-5 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.table.input') placeholderplaceholder</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100 bg-white dark:divide-dark-800 dark:bg-dark-800">
                <tr v-if="logsLoading">
                  <td colspan="10" class="px-5 py-12 text-center text-sm text-gray-500 dark:text-gray-400">{{ t('common.loading') placeholderplaceholder</td>
                </tr>
                <tr v-else-if="logs.length === 0">
                  <td colspan="10" class="px-5 py-12 text-center text-sm text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.emptyLogs') placeholderplaceholder</td>
                </tr>
                <template v-else>
                  <tr v-for="row in logs" :key="row.id" class="hover:bg-gray-50 dark:hover:bg-dark-700/60">
                    <td class="whitespace-nowrap px-5 py-4 text-sm text-gray-700 dark:text-gray-300">{{ formatDateTime(row.created_at) placeholderplaceholder</td>
                    <td class="whitespace-nowrap px-5 py-4 text-sm text-gray-700 dark:text-gray-300">{{ row.group_name || '-' placeholderplaceholder</td>
                    <td class="whitespace-nowrap px-5 py-4 text-sm text-gray-700 dark:text-gray-300">
                      <div>{{ row.user_email || '-' placeholderplaceholder</div>
                      <div v-if="row.user_id" class="text-xs text-gray-400">UID {{ row.user_id placeholderplaceholder</div>
                    </td>
                    <td class="whitespace-nowrap px-5 py-4 text-sm text-gray-700 dark:text-gray-300">{{ row.api_key_name || '-' placeholderplaceholder</td>
                    <td class="whitespace-nowrap px-5 py-4 text-sm text-gray-700 dark:text-gray-300">
                      <div>{{ row.endpoint || '-' placeholderplaceholder</div>
                      <div class="text-xs text-gray-400">{{ row.provider || '-' placeholderplaceholder / {{ row.model || '-' placeholderplaceholder</div>
                    </td>
                    <td class="whitespace-nowrap px-5 py-4">
                      <span class="inline-flex rounded-md px-2 py-1 text-xs font-medium" :class="resultBadgeClass(row)">
                        {{ resultLabel(row) placeholderplaceholder
                      </span>
                    </td>
                    <td class="whitespace-nowrap px-5 py-4 text-sm text-gray-700 dark:text-gray-300">
                      <div>{{ row.highest_category || '-' placeholderplaceholder</div>
                      <div class="text-xs text-gray-400">{{ percent(row.highest_score) placeholderplaceholder</div>
                    </td>
                    <td class="whitespace-nowrap px-5 py-4 text-sm text-gray-700 dark:text-gray-300">
                      <div>{{ violationCountText(row) placeholderplaceholder</div>
                      <div class="text-xs text-gray-400">
                        {{ row.email_sent ? t('admin.riskControl.emailSent') : t('admin.riskControl.emailNotSent') placeholderplaceholder
                        <span v-if="row.auto_banned"> / {{ t('admin.riskControl.autoBanned') placeholderplaceholder</span>
                      </div>
                      <button
                        v-if="canUnbanRow(row)"
                        type="button"
                        class="mt-2 inline-flex items-center gap-1 rounded-md border border-emerald-200 bg-emerald-50 px-2 py-1 text-xs font-medium text-emerald-700 transition-colors hover:bg-emerald-100 disabled:cursor-not-allowed disabled:opacity-60 dark:border-emerald-900/60 dark:bg-emerald-900/20 dark:text-emerald-300 dark:hover:bg-emerald-900/30"
                        :disabled="unbanningUserID === row.user_id"
                        @click="unbanUser(row)"
                      >
                        <Icon name="checkCircle" size="xs" :class="unbanningUserID === row.user_id ? 'animate-spin' : ''" />
                        {{ unbanningUserID === row.user_id ? t('common.processing') : t('admin.riskControl.unbanUser') placeholderplaceholder
                      </button>
                    </td>
                    <td class="whitespace-nowrap px-5 py-4 text-sm text-gray-700 dark:text-gray-300">
                      <div>{{ latencyText(row.upstream_latency_ms) placeholderplaceholder</div>
                      <div v-if="row.queue_delay_ms !== null && row.queue_delay_ms !== undefined" class="text-xs text-gray-400">
                        {{ t('admin.riskControl.queueDelay', { ms: row.queue_delay_ms placeholder) placeholderplaceholder
                      </div>
                    </td>
                    <td class="w-[320px] max-w-sm px-5 py-4 text-sm text-gray-700 dark:text-gray-300">
                      <button
                        type="button"
                        class="group flex w-full min-w-0 items-center gap-2 rounded-lg px-2 py-1.5 text-left transition-colors hover:bg-gray-100 dark:hover:bg-dark-700"
                        :title="inputSummaryText(row)"
                        @click="openInputDetail(row)"
                      >
                        <span class="min-w-0 flex-1 truncate">{{ inputSummaryText(row) placeholderplaceholder</span>
                        <Icon name="eye" size="xs" class="flex-shrink-0 text-gray-300 transition-colors group-hover:text-primary-500 dark:text-gray-500" />
                      </button>
                    </td>
                  </tr>
                </template>
              </tbody>
            </table>
          </div>

          <Pagination
            v-if="pagination.total > 0"
            :page="pagination.page"
            :total="pagination.total"
            :page-size="pagination.page_size"
            @update:page="onPageChange"
            @update:pageSize="onPageSizeChange"
          />
        </div>
      </template>

      <BaseDialog :show="settingsOpen" :title="t('admin.riskControl.settingsTitle')" width="extra-wide" @close="settingsOpen = false">
        <div class="space-y-6">
          <div class="flex gap-2 overflow-x-auto border-b border-gray-100 pb-3 dark:border-dark-700">
            <button
              v-for="tab in settingsTabs"
              :key="tab.id"
              type="button"
              class="inline-flex whitespace-nowrap rounded-lg px-3 py-2 text-sm font-medium transition-colors"
              :class="activeSettingsTab === tab.id ? 'bg-primary-50 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300' : 'text-gray-500 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-dark-700 dark:hover:text-white'"
              @click="activeSettingsTab = tab.id"
            >
              {{ tab.label placeholderplaceholder
            </button>
          </div>

          <div v-if="activeSettingsTab === 'basic'" class="space-y-5">
            <div class="grid grid-cols-1 gap-5 lg:grid-cols-2">
              <div class="flex items-center justify-between rounded-lg border border-gray-100 p-4 dark:border-dark-700">
                <div>
                  <p class="text-sm font-medium text-gray-900 dark:text-white">{{ t('admin.riskControl.enabled') placeholderplaceholder</p>
                  <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.enabledHint') placeholderplaceholder</p>
                </div>
                <Toggle v-model="configForm.enabled" />
              </div>
              <div>
                <label class="input-label">{{ t('admin.riskControl.mode') placeholderplaceholder</label>
                <Select v-model="configForm.mode" :options="modeOptions" />
                <p class="mt-2 text-xs leading-5 text-gray-500 dark:text-gray-400">{{ modeDescription(configForm.mode) placeholderplaceholder</p>
              </div>
              <div>
                <label class="input-label">{{ t('admin.riskControl.baseUrl') placeholderplaceholder</label>
                <input v-model.trim="configForm.base_url" type="url" class="input" placeholder="https://api.openai.com" />
              </div>
              <div>
                <label class="input-label">{{ t('admin.riskControl.model') placeholderplaceholder</label>
                <input v-model.trim="configForm.model" type="text" class="input" placeholder="omni-moderation-latest" />
              </div>
              <div>
                <label class="input-label">{{ t('admin.riskControl.timeoutMs') placeholderplaceholder</label>
                <input v-model.number="configForm.timeout_ms" type="number" min="500" max="30000" class="input" />
              </div>
              <div>
                <label class="input-label">{{ t('admin.riskControl.retryCount') placeholderplaceholder</label>
                <input v-model.number="configForm.retry_count" type="number" min="0" max="5" class="input" />
              </div>
              <div>
                <label class="input-label">{{ t('admin.riskControl.sampleRate') placeholderplaceholder</label>
                <div class="relative">
                  <input v-model.number="configForm.sample_rate" type="number" min="0" max="100" step="1" class="input pr-8" />
                  <span class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-gray-400">%</span>
                </div>
              </div>
            </div>

            <div class="overflow-hidden rounded-xl border border-gray-100 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-800">
              <div class="flex flex-col gap-4 border-b border-gray-100 bg-gray-50 px-4 py-4 dark:border-dark-700 dark:bg-dark-800/60 lg:flex-row lg:items-center lg:justify-between">
                <div class="flex items-start gap-3">
                  <span class="mt-0.5 flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-lg bg-primary-50 text-primary-600 dark:bg-primary-900/30 dark:text-primary-300">
                    <Icon name="key" size="md" />
                  </span>
                  <div>
                    <label class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('admin.riskControl.apiKeys') placeholderplaceholder</label>
                    <p class="mt-1 max-w-3xl text-xs leading-5 text-gray-500 dark:text-gray-400">
                      {{ t('admin.riskControl.apiKeysHint', { count: configForm.api_key_count placeholder) placeholderplaceholder
                    </p>
                  </div>
                </div>
                <div class="flex flex-wrap items-center gap-2">
                  <button
                    type="button"
                    class="btn btn-secondary inline-flex items-center gap-2"
                    :disabled="apiKeyTesting || inputApiKeyCount === 0 || configForm.clear_api_key"
                    @click="testApiKeys(true)"
                  >
                    <Icon name="beaker" size="sm" :class="apiKeyTesting ? 'animate-pulse' : ''" />
                    {{ apiKeyTesting ? t('admin.riskControl.testingApiKeys') : t('admin.riskControl.testInputApiKeys') placeholderplaceholder
                  </button>
                  <button
                    type="button"
                    class="btn btn-secondary inline-flex items-center gap-2"
                    :disabled="apiKeyTesting || !configForm.api_key_configured || configForm.clear_api_key"
	                    @click="testApiKeys(false)"
	                  >
	                    <Icon name="shield" size="sm" />
	                    {{ storedApiKeyTestButtonText placeholderplaceholder
	                  </button>
                  <button
                    v-if="configForm.api_key_configured"
                    type="button"
                    class="btn btn-secondary inline-flex items-center gap-2"
                    @click="toggleClearApiKey"
                  >
                    <Icon :name="configForm.clear_api_key ? 'x' : 'trash'" size="sm" />
                    {{ configForm.clear_api_key ? t('admin.riskControl.keepApiKey') : t('admin.riskControl.clearApiKey') placeholderplaceholder
                  </button>
                </div>
              </div>

              <div class="grid grid-cols-1 gap-4 p-4 xl:grid-cols-[minmax(0,1fr)_minmax(360px,440px)]">
                <div class="space-y-3">
                  <textarea
                    v-model="configForm.api_keys_text"
                    class="input min-h-44 resize-y font-mono text-sm"
                    :placeholder="apiKeyPlaceholder"
                    autocomplete="new-password"
                    :disabled="configForm.clear_api_key"
                  ></textarea>
                  <div class="flex flex-wrap items-center gap-2 text-xs text-gray-500 dark:text-gray-400">
                    <span class="inline-flex rounded-md bg-gray-100 px-2 py-1 dark:bg-dark-700">
                      {{ t('admin.riskControl.inputApiKeyCount', { count: inputApiKeyCount placeholder) placeholderplaceholder
                    </span>
                    <span v-if="configForm.api_key_configured" class="inline-flex rounded-md bg-gray-100 px-2 py-1 dark:bg-dark-700">
                      {{ t('admin.riskControl.storedApiKeyCount', { count: configForm.api_key_count placeholder) placeholderplaceholder
                    </span>
                    <span v-if="configForm.clear_api_key" class="inline-flex rounded-md bg-red-50 px-2 py-1 text-red-700 dark:bg-red-900/20 dark:text-red-300">
                      {{ t('admin.riskControl.apiKeyWillClear') placeholderplaceholder
                    </span>
                  </div>

                  <div class="rounded-lg border border-gray-100 bg-gray-50 p-3 dark:border-dark-700 dark:bg-dark-900/30" @paste="handleModerationImagePaste">
                    <div class="mb-3 flex items-center justify-between gap-3">
                      <div>
                        <p class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('admin.riskControl.auditTestInput') placeholderplaceholder</p>
                        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.auditTestInputHint') placeholderplaceholder</p>
                      </div>
                      <button
                        v-if="moderationTestPrompt || moderationTestImages.length > 0 || moderationTestResult"
                        type="button"
                        class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-gray-500 hover:bg-white hover:text-gray-900 dark:text-gray-400 dark:hover:bg-dark-800 dark:hover:text-white"
                        @click="clearModerationTestInput"
                      >
                        <Icon name="x" size="xs" />
                        {{ t('admin.riskControl.clearAuditTest') placeholderplaceholder
                      </button>
                    </div>
                    <textarea
                      v-model="moderationTestPrompt"
                      class="input min-h-24 resize-y text-sm"
                      :placeholder="t('admin.riskControl.auditTestPromptPlaceholder')"
                    ></textarea>
                    <div
                      class="mt-3 rounded-lg border border-dashed border-gray-200 bg-white p-3 dark:border-dark-700 dark:bg-dark-800"
                      @dragover.prevent
                      @drop.prevent="handleModerationImageDrop"
                    >
                      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                        <div class="flex items-start gap-2">
                          <Icon name="upload" size="md" class="mt-0.5 text-gray-400" />
                          <div>
                            <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ t('admin.riskControl.auditTestImages') placeholderplaceholder</p>
                            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.auditTestImagesHint') placeholderplaceholder</p>
                          </div>
                        </div>
                        <label class="btn btn-secondary inline-flex cursor-pointer items-center gap-2">
                          <Icon name="plus" size="sm" />
                          {{ t('admin.riskControl.addAuditTestImage') placeholderplaceholder
                          <input type="file" accept="image/*" multiple class="sr-only" @change="handleModerationImageUpload" />
                        </label>
                      </div>
                      <div v-if="moderationTestImages.length > 0" class="mt-3 grid grid-cols-2 gap-2 sm:grid-cols-4">
                        <div
                          v-for="(image, index) in moderationTestImages"
                          :key="image.slice(0, 64) + index"
                          class="group relative aspect-square overflow-hidden rounded-lg border border-gray-100 bg-gray-100 dark:border-dark-700 dark:bg-dark-700"
                        >
                          <img :src="image" alt="" class="h-full w-full object-cover" />
                          <button
                            type="button"
                            class="absolute right-1.5 top-1.5 flex h-7 w-7 items-center justify-center rounded-full bg-black/60 text-white opacity-0 transition-opacity group-hover:opacity-100"
                            @click="removeModerationTestImage(index)"
                          >
                            <Icon name="x" size="xs" :stroke-width="2" />
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="rounded-lg border border-gray-100 bg-gray-50 p-3 dark:border-dark-700 dark:bg-dark-900/30">
                  <div class="mb-3 flex items-center justify-between gap-3">
                    <div>
                      <p class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('admin.riskControl.apiKeyHealth') placeholderplaceholder</p>
                      <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.apiKeyFreezeRule') placeholderplaceholder</p>
                    </div>
                    <span class="inline-flex rounded-md bg-white px-2 py-1 text-xs font-medium text-gray-600 shadow-sm dark:bg-dark-800 dark:text-gray-300">
                      {{ t('admin.riskControl.apiKeyRows', { count: apiKeyRows.length placeholder) placeholderplaceholder
                    </span>
                  </div>

                  <div v-if="apiKeyRows.length === 0" class="flex min-h-32 flex-col items-center justify-center rounded-lg border border-dashed border-gray-200 bg-white px-4 py-6 text-center dark:border-dark-700 dark:bg-dark-800">
                    <Icon name="infoCircle" size="lg" class="text-gray-300 dark:text-dark-500" />
                    <p class="mt-2 text-sm font-medium text-gray-700 dark:text-gray-200">{{ t('admin.riskControl.apiKeyHealthEmpty') placeholderplaceholder</p>
                    <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.apiKeyHealthEmptyHint') placeholderplaceholder</p>
                  </div>
                  <div v-else class="max-h-72 space-y-2 overflow-y-auto pr-1">
                    <div
                      v-for="(row, index) in apiKeyRows"
                      :key="apiKeyRowKey(row, index)"
                      class="rounded-lg border border-gray-100 bg-white p-3 shadow-sm dark:border-dark-700 dark:bg-dark-800"
                    >
                      <div class="flex items-start justify-between gap-3">
                        <div class="min-w-0">
                          <div class="flex min-w-0 flex-wrap items-center gap-2">
                            <span class="truncate font-mono text-sm font-semibold text-gray-900 dark:text-white">{{ row.masked || '-' placeholderplaceholder</span>
                            <span
                              class="inline-flex rounded-md px-1.5 py-0.5 text-[11px] font-medium"
                              :class="row.configured ? 'bg-primary-50 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300' : 'bg-purple-50 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300'"
                            >
                              {{ row.configured ? t('admin.riskControl.apiKeyConfigured') : t('admin.riskControl.apiKeyTemporary') placeholderplaceholder
                            </span>
                          </div>
                          <p class="mt-1 text-xs leading-5 text-gray-500 dark:text-gray-400">{{ apiKeyStatusMeta(row) placeholderplaceholder</p>
                        </div>
                        <span class="inline-flex flex-shrink-0 items-center gap-1.5 rounded-full px-2 py-1 text-xs font-medium" :class="apiKeyStatusBadgeClass(row.status)">
                          <span class="h-1.5 w-1.5 rounded-full" :class="apiKeyStatusDotClass(row.status)"></span>
                          {{ apiKeyStatusLabel(row.status) placeholderplaceholder
                        </span>
                      </div>
                      <p v-if="row.last_error" class="mt-2 rounded-md bg-amber-50 px-2 py-1.5 text-xs leading-5 text-amber-700 dark:bg-amber-900/20 dark:text-amber-300">
                        {{ row.last_error placeholderplaceholder
                      </p>
                    </div>
                  </div>

                  <div v-if="moderationTestResult" class="mt-4 rounded-lg border border-gray-100 bg-white p-3 dark:border-dark-700 dark:bg-dark-800">
                    <div class="flex items-start justify-between gap-3">
                      <div>
                        <p class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('admin.riskControl.auditTestResult') placeholderplaceholder</p>
                        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                          {{ t('admin.riskControl.auditTestHighest', { category: moderationTestResult.highest_category || '-', score: percent(moderationTestResult.highest_score) placeholder) placeholderplaceholder
                        </p>
                      </div>
                      <span class="inline-flex rounded-full px-2 py-1 text-xs font-medium" :class="moderationTestResult.flagged ? 'bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-300' : 'bg-emerald-50 text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300'">
                        {{ moderationTestResult.flagged ? t('admin.riskControl.auditTestFlagged') : t('admin.riskControl.auditTestPassed') placeholderplaceholder
                      </span>
                    </div>
                    <div class="mt-3">
                      <div class="mb-2 flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
                        <span>{{ t('admin.riskControl.auditTestComposite') placeholderplaceholder</span>
                        <span class="font-semibold text-gray-900 dark:text-white">{{ percent(moderationTestResult.composite_score) placeholderplaceholder</span>
                      </div>
                      <div class="h-2 overflow-hidden rounded-full bg-gray-100 dark:bg-dark-700">
                        <div class="h-full rounded-full" :class="moderationTestResult.flagged ? 'bg-red-500' : 'bg-emerald-500'" :style="{ width: percentWidth(moderationTestResult.composite_score) placeholder"></div>
                      </div>
                    </div>
                    <div class="mt-3 max-h-52 space-y-2 overflow-y-auto pr-1">
                      <div v-for="score in moderationScoreRows" :key="score.category">
                        <div class="mb-1 flex items-center justify-between gap-3 text-xs">
                          <span class="truncate text-gray-600 dark:text-gray-300">{{ score.category placeholderplaceholder</span>
                          <span class="font-mono text-gray-500 dark:text-gray-400">{{ percent(score.score) placeholderplaceholder / {{ percent(score.threshold) placeholderplaceholder</span>
                        </div>
                        <div class="h-1.5 overflow-hidden rounded-full bg-gray-100 dark:bg-dark-700">
                          <div class="h-full rounded-full" :class="score.hit ? 'bg-red-500' : 'bg-primary-500'" :style="{ width: percentWidth(score.score) placeholder"></div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-else-if="activeSettingsTab === 'scope'" class="space-y-5">
            <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
              <div>
                <h3 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('admin.riskControl.groupScope') placeholderplaceholder</h3>
                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.groupScopeHint') placeholderplaceholder</p>
              </div>
              <div class="inline-flex rounded-lg bg-gray-100 p-1 dark:bg-dark-700">
                <button
                  type="button"
                  class="rounded-md px-3 py-1.5 text-sm font-medium transition-colors"
                  :class="configForm.all_groups ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-800 dark:text-white' : 'text-gray-500 dark:text-gray-400'"
                  @click="configForm.all_groups = true"
                >
                  {{ t('admin.riskControl.allGroups') placeholderplaceholder
                </button>
                <button
                  type="button"
                  class="rounded-md px-3 py-1.5 text-sm font-medium transition-colors"
                  :class="!configForm.all_groups ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-800 dark:text-white' : 'text-gray-500 dark:text-gray-400'"
                  @click="configForm.all_groups = false"
                >
                  {{ t('admin.riskControl.selectedGroups') placeholderplaceholder
                </button>
              </div>
            </div>

            <div v-if="!configForm.all_groups" class="space-y-4">
              <div class="relative">
                <Icon name="search" size="sm" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                <input v-model.trim="groupSearch" type="search" class="input pl-9" :placeholder="t('admin.riskControl.searchGroups')" />
              </div>
              <div class="grid max-h-[420px] grid-cols-1 gap-3 overflow-y-auto pr-1 md:grid-cols-2 xl:grid-cols-3">
                <button
                  v-for="group in filteredGroups"
                  :key="group.id"
                  type="button"
                  class="flex min-h-20 items-center justify-between rounded-lg border p-4 text-left transition-colors"
                  :class="isGroupSelected(group.id) ? 'border-primary-300 bg-primary-50 dark:border-primary-700 dark:bg-primary-900/20' : 'border-gray-100 hover:bg-gray-50 dark:border-dark-700 dark:hover:bg-dark-700/60'"
                  @click="toggleGroup(group.id)"
                >
                  <span class="min-w-0">
                    <span class="block truncate text-sm font-semibold text-gray-900 dark:text-white">{{ group.name placeholderplaceholder</span>
                    <span class="mt-1 inline-flex rounded-md bg-gray-100 px-2 py-0.5 text-xs text-gray-500 dark:bg-dark-700 dark:text-gray-400">{{ group.platform placeholderplaceholder</span>
                  </span>
                  <span
                    class="flex h-5 w-5 flex-shrink-0 items-center justify-center rounded-full border"
                    :class="isGroupSelected(group.id) ? 'border-primary-500 bg-primary-500 text-white' : 'border-gray-300 text-transparent dark:border-dark-500'"
                  >
                    <Icon name="check" size="xs" :stroke-width="2" />
                  </span>
                </button>
                <p v-if="filteredGroups.length === 0" class="text-sm text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.noGroups') placeholderplaceholder</p>
              </div>
            </div>
          </div>

          <div v-else-if="activeSettingsTab === 'runtime'" class="grid grid-cols-1 gap-5 lg:grid-cols-2">
            <div>
              <label class="input-label">{{ t('admin.riskControl.workerCount') placeholderplaceholder</label>
              <input v-model.number="configForm.worker_count" type="number" min="1" max="32" class="input" />
            </div>
            <div>
              <label class="input-label">{{ t('admin.riskControl.queueSize') placeholderplaceholder</label>
              <input v-model.number="configForm.queue_size" type="number" min="100" max="100000" class="input" />
            </div>
            <div class="flex items-center justify-between rounded-lg border border-gray-100 p-4 dark:border-dark-700 lg:col-span-2">
              <div>
                <p class="text-sm font-medium text-gray-900 dark:text-white">{{ t('admin.riskControl.recordNonHits') placeholderplaceholder</p>
                <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.recordNonHitsHint') placeholderplaceholder</p>
              </div>
              <Toggle v-model="configForm.record_non_hits" />
            </div>
            <div class="space-y-4 rounded-lg border border-gray-100 p-4 dark:border-dark-700 lg:col-span-2">
              <div class="flex items-center justify-between gap-4">
                <div>
                  <p class="text-sm font-medium text-gray-900 dark:text-white">{{ t('admin.riskControl.preHashCheck') placeholderplaceholder</p>
                  <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.preHashCheckHint') placeholderplaceholder</p>
                </div>
                <Toggle v-model="configForm.pre_hash_check_enabled" />
              </div>
              <div class="rounded-lg bg-gray-50 p-3 dark:bg-dark-900/30">
                <div class="flex flex-col gap-3 xl:flex-row xl:items-center xl:justify-between">
                  <div>
                    <p class="text-sm font-medium text-gray-900 dark:text-white">
                      {{ t('admin.riskControl.flaggedHashCount', { count: formatNumber(status?.flagged_hash_count ?? 0) placeholder) placeholderplaceholder
                    </p>
                    <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.flaggedHashHint') placeholderplaceholder</p>
                  </div>
                  <button
                    type="button"
                    class="btn btn-secondary inline-flex items-center justify-center gap-2 text-red-600 hover:text-red-700 dark:text-red-300"
                    :disabled="hashActionLoading || (status?.flagged_hash_count ?? 0) === 0"
                    @click="clearFlaggedHashes"
                  >
                    <Icon name="trash" size="sm" :class="hashActionLoading ? 'animate-pulse' : ''" />
                    {{ t('admin.riskControl.clearFlaggedHashes') placeholderplaceholder
                  </button>
                </div>
                <div class="mt-3 flex flex-col gap-2 sm:flex-row">
                  <input
                    v-model.trim="flaggedHashInput"
                    type="text"
                    class="input font-mono text-sm"
                    :placeholder="t('admin.riskControl.flaggedHashPlaceholder')"
                  />
                  <button
                    type="button"
                    class="btn btn-secondary inline-flex items-center justify-center gap-2"
                    :disabled="hashActionLoading || !isFlaggedHashInputValid"
                    @click="deleteFlaggedHash"
                  >
                    <Icon name="trash" size="sm" />
                    {{ t('admin.riskControl.deleteFlaggedHash') placeholderplaceholder
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div v-else-if="activeSettingsTab === 'response'" class="space-y-5">
            <div class="grid grid-cols-1 gap-5 lg:grid-cols-2">
              <div>
                <label class="input-label">{{ t('admin.riskControl.blockStatus') placeholderplaceholder</label>
                <input v-model.number="configForm.block_status" type="number" min="400" max="599" class="input" />
              </div>
              <div>
                <label class="input-label">{{ t('admin.riskControl.blockMessage') placeholderplaceholder</label>
                <input v-model.trim="configForm.block_message" type="text" class="input" />
              </div>
              <div class="flex items-center justify-between rounded-lg border border-gray-100 p-4 dark:border-dark-700">
                <div>
                  <p class="text-sm font-medium text-gray-900 dark:text-white">{{ t('admin.riskControl.emailOnHit') placeholderplaceholder</p>
                  <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.emailOnHitHint') placeholderplaceholder</p>
                </div>
                <Toggle v-model="configForm.email_on_hit" />
              </div>
              <div class="flex items-center justify-between rounded-lg border border-gray-100 p-4 dark:border-dark-700">
                <div>
                  <p class="text-sm font-medium text-gray-900 dark:text-white">{{ t('admin.riskControl.autoBan') placeholderplaceholder</p>
                  <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.autoBanHint') placeholderplaceholder</p>
                </div>
                <Toggle v-model="configForm.auto_ban_enabled" />
              </div>
              <div>
                <label class="input-label">{{ t('admin.riskControl.banThreshold') placeholderplaceholder</label>
                <input v-model.number="configForm.ban_threshold" type="number" min="1" max="1000" class="input" />
              </div>
              <div>
                <label class="input-label">{{ t('admin.riskControl.violationWindowHours') placeholderplaceholder</label>
                <input v-model.number="configForm.violation_window_hours" type="number" min="1" max="8760" class="input" />
              </div>
            </div>
          </div>

          <div v-else class="grid grid-cols-1 gap-5 lg:grid-cols-2">
            <div>
              <label class="input-label">{{ t('admin.riskControl.hitRetentionDays') placeholderplaceholder</label>
              <input v-model.number="configForm.hit_retention_days" type="number" min="1" max="3650" class="input" />
            </div>
            <div>
              <label class="input-label">{{ t('admin.riskControl.nonHitRetentionDays') placeholderplaceholder</label>
              <input v-model.number="configForm.non_hit_retention_days" type="number" min="1" max="3" class="input" />
            </div>
            <div class="rounded-lg border border-gray-100 p-4 text-sm text-gray-500 dark:border-dark-700 dark:text-gray-400 lg:col-span-2">
              <div class="flex flex-wrap items-center gap-3">
                <Icon name="database" size="md" class="text-gray-400" />
                <span>{{ t('admin.riskControl.cleanupStats', { hit: status?.last_cleanup_deleted_hit ?? 0, nonHit: status?.last_cleanup_deleted_non_hit ?? 0 placeholder) placeholderplaceholder</span>
              </div>
            </div>
          </div>
        </div>

        <template #footer>
          <div class="flex justify-end gap-2">
            <button type="button" class="btn btn-secondary" @click="settingsOpen = false">{{ t('common.cancel') placeholderplaceholder</button>
            <button type="button" class="btn btn-primary inline-flex items-center gap-2" :disabled="saving" @click="saveConfig">
              <Icon v-if="saving" name="refresh" size="sm" class="animate-spin" />
              <Icon v-else name="check" size="sm" />
              {{ saving ? t('common.saving') : t('admin.riskControl.saveConfig') placeholderplaceholder
            </button>
          </div>
        </template>
      </BaseDialog>

      <BaseDialog
        :show="inputDetailRow !== null"
        :title="t('admin.riskControl.inputDetailTitle')"
        width="wide"
        @close="closeInputDetail"
      >
        <div v-if="inputDetailRow" class="space-y-5">
          <div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
            <div class="rounded-lg border border-gray-100 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800/70">
              <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.table.time') placeholderplaceholder</p>
              <p class="mt-1 truncate text-sm font-semibold text-gray-900 dark:text-white">{{ formatDateTime(inputDetailRow.created_at) placeholderplaceholder</p>
            </div>
            <div class="rounded-lg border border-gray-100 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800/70">
              <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.table.user') placeholderplaceholder</p>
              <p class="mt-1 truncate text-sm font-semibold text-gray-900 dark:text-white">{{ inputDetailRow.user_email || '-' placeholderplaceholder</p>
            </div>
            <div class="rounded-lg border border-gray-100 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800/70">
              <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.table.result') placeholderplaceholder</p>
              <span class="mt-1 inline-flex rounded-md px-2 py-1 text-xs font-medium" :class="resultBadgeClass(inputDetailRow)">
                {{ resultLabel(inputDetailRow) placeholderplaceholder
              </span>
            </div>
            <div class="rounded-lg border border-gray-100 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800/70">
              <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('admin.riskControl.table.highest') placeholderplaceholder</p>
              <p class="mt-1 truncate text-sm font-semibold text-gray-900 dark:text-white">
                {{ inputDetailRow.highest_category || '-' placeholderplaceholder / {{ percent(inputDetailRow.highest_score) placeholderplaceholder
              </p>
            </div>
          </div>

          <div class="rounded-xl border border-gray-100 bg-white p-4 shadow-sm dark:border-dark-700 dark:bg-dark-800">
            <div class="flex flex-wrap items-center justify-between gap-3">
              <div>
                <p class="text-sm font-semibold text-gray-900 dark:text-white">{{ t('admin.riskControl.inputDetailContent') placeholderplaceholder</p>
                <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                  {{ inputDetailRow.endpoint || '-' placeholderplaceholder · {{ inputDetailRow.provider || '-' placeholderplaceholder / {{ inputDetailRow.model || '-' placeholderplaceholder
                </p>
              </div>
              <span v-if="inputDetailRow.group_name" class="inline-flex rounded-md bg-sky-50 px-2.5 py-1 text-xs font-medium text-sky-700 dark:bg-sky-900/20 dark:text-sky-300">
                {{ inputDetailRow.group_name placeholderplaceholder
              </span>
            </div>
            <pre class="mt-4 max-h-[420px] overflow-auto whitespace-pre-wrap break-words rounded-lg bg-gray-950 p-4 text-sm leading-6 text-gray-100 shadow-inner dark:bg-black/50">{{ inputDetailText placeholderplaceholder</pre>
          </div>
        </div>

        <template #footer>
          <div class="flex justify-end">
            <button type="button" class="btn btn-secondary" @click="closeInputDetail">{{ t('common.close') placeholderplaceholder</button>
          </div>
        </template>
      </BaseDialog>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, ref placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'
import Select from '@/components/common/Select.vue'
import Toggle from '@/components/common/Toggle.vue'
import Pagination from '@/components/common/Pagination.vue'
import { adminAPI placeholder from '@/api/admin'
import type {
  ContentModerationAPIKeyStatus,
  ContentModerationConfig,
  ContentModerationLog,
  ContentModerationRuntimeStatus,
  ContentModerationTestAuditResult,
  ModerationMode,
  UpdateContentModerationConfig,
placeholder from '@/api/admin/riskControl'
import type { AdminGroup, SelectOption placeholder from '@/types'
import { useAppStore placeholder from '@/stores/app'
import { extractApiErrorMessage placeholder from '@/utils/apiError'
import { formatDateTime as formatDateTimeValue placeholder from '@/utils/format'

type SettingsTab = 'basic' | 'scope' | 'runtime' | 'response' | 'retention'
type WorkerSlotState = 'active' | 'idle' | 'disabled'
type OverviewIcon = 'shield' | 'key' | 'users' | 'document'
type OverviewItem = {
  key: string
  label: string
  value: string
  meta: string
  icon: OverviewIcon
  iconClass: string
  badge?: string
  badgeClass?: string
placeholder
type ModerationScoreRow = {
  category: string
  score: number
  threshold: number
  hit: boolean
placeholder

const maxModerationTestImages = 4
const maxModerationTestImageSize = 8 * 1024 * 1024

const { t placeholder = useI18n()
const appStore = useAppStore()

const loading = ref(true)
const saving = ref(false)
const logsLoading = ref(false)
const statusLoading = ref(false)
const apiKeyTesting = ref(false)
const hashActionLoading = ref(false)
const unbanningUserID = ref<number | null>(null)
const settingsOpen = ref(false)
const activeSettingsTab = ref<SettingsTab>('basic')
const groupSearch = ref('')
const flaggedHashInput = ref('')
const groups = ref<AdminGroup[]>([])
const logs = ref<ContentModerationLog[]>([])
const status = ref<ContentModerationRuntimeStatus | null>(null)
const testedApiKeyStatuses = ref<ContentModerationAPIKeyStatus[]>([])
const moderationTestPrompt = ref('')
const moderationTestImages = ref<string[]>([])
const moderationTestResult = ref<ContentModerationTestAuditResult | null>(null)
const inputDetailRow = ref<ContentModerationLog | null>(null)
let statusTimer: number | null = null

const configForm = reactive({
  enabled: false,
  mode: 'pre_block' as ModerationMode,
  base_url: 'https://api.openai.com',
  model: 'omni-moderation-latest',
  api_keys_text: '',
  api_key_configured: false,
  api_key_masked: '',
  api_key_count: 0,
  api_key_masks: [] as string[],
  api_key_statuses: [] as ContentModerationAPIKeyStatus[],
  clear_api_key: false,
  timeout_ms: 3000,
  retry_count: 2,
  sample_rate: 100,
  all_groups: true,
  group_ids: [] as number[],
  record_non_hits: false,
  worker_count: 4,
  queue_size: 32768,
  block_status: 403,
  block_message: '内容审计命中风险规则，请调整输入后重试',
  email_on_hit: true,
  auto_ban_enabled: true,
  ban_threshold: 10,
  violation_window_hours: 720,
  hit_retention_days: 180,
  non_hit_retention_days: 3,
  pre_hash_check_enabled: false,
placeholder)

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0,
  pages: 1,
placeholder)

const filters = reactive({
  result: '',
  group_id: 0,
  endpoint: '',
  search: '',
  from: '',
  to: '',
placeholder)

const settingsTabs = computed<Array<{ id: SettingsTab; label: string placeholder>>(() => [
  { id: 'basic', label: t('admin.riskControl.tabs.basic') placeholder,
  { id: 'scope', label: t('admin.riskControl.tabs.scope') placeholder,
  { id: 'runtime', label: t('admin.riskControl.tabs.runtime') placeholder,
  { id: 'response', label: t('admin.riskControl.tabs.response') placeholder,
  { id: 'retention', label: t('admin.riskControl.tabs.retention') placeholder,
])

const modeOptions = computed<SelectOption[]>(() => [
  { value: 'pre_block', label: t('admin.riskControl.modePreBlock') placeholder,
  { value: 'observe', label: t('admin.riskControl.modeObserve') placeholder,
  { value: 'off', label: t('admin.riskControl.modeOff') placeholder,
])

const resultOptions = computed<SelectOption[]>(() => [
  { value: '', label: t('admin.riskControl.result.all') placeholder,
  { value: 'hit', label: t('admin.riskControl.result.hit') placeholder,
  { value: 'blocked', label: t('admin.riskControl.result.blocked') placeholder,
  { value: 'pass', label: t('admin.riskControl.result.pass') placeholder,
  { value: 'error', label: t('admin.riskControl.result.error') placeholder,
])

const endpointOptions = computed<SelectOption[]>(() => [
  { value: '', label: t('admin.riskControl.filters.allEndpoints') placeholder,
  { value: '/v1/messages', label: '/v1/messages' placeholder,
  { value: '/v1/responses', label: '/v1/responses' placeholder,
  { value: '/v1/chat/completions', label: '/v1/chat/completions' placeholder,
  { value: '/v1beta/models', label: '/v1beta/models' placeholder,
  { value: '/v1/images/generations', label: '/v1/images/generations' placeholder,
  { value: '/v1/images/edits', label: '/v1/images/edits' placeholder,
])

const groupFilterOptions = computed<SelectOption[]>(() => [
  { value: 0, label: t('admin.riskControl.filters.allGroups') placeholder,
  ...groups.value.map((group) => ({
    value: group.id,
    label: `${group.nameplaceholder (${group.platformplaceholder)`,
  placeholder)),
])

const selectedGroupCount = computed(() => String(configForm.group_ids.length))

const filteredGroups = computed(() => {
  const keyword = groupSearch.value.trim().toLowerCase()
  if (!keyword) return groups.value
  return groups.value.filter((group) => {
    return group.name.toLowerCase().includes(keyword) || String(group.platform).toLowerCase().includes(keyword)
  placeholder)
placeholder)

const apiKeyPlaceholder = computed(() => {
  if (configForm.clear_api_key) return t('admin.riskControl.apiKeyWillClear')
  if (configForm.api_key_configured) return t('admin.riskControl.apiKeysPlaceholderKeep')
  return t('admin.riskControl.apiKeysPlaceholder')
placeholder)

const inputApiKeyCount = computed(() => parseApiKeys(configForm.api_keys_text).length)

const hasModerationAuditInput = computed(() => {
  return moderationTestPrompt.value.trim() !== '' || moderationTestImages.value.length > 0
placeholder)

const isFlaggedHashInputValid = computed(() => /^[a-fA-F0-9]{64placeholder$/.test(flaggedHashInput.value.trim()))

const storedApiKeyTestButtonText = computed(() => {
  if (apiKeyTesting.value) return t('admin.riskControl.testingApiKeys')
  if (hasModerationAuditInput.value) return t('admin.riskControl.testContentWithStoredApiKey')
  return t('admin.riskControl.testStoredApiKeys')
placeholder)

const savedApiKeyRows = computed<ContentModerationAPIKeyStatus[]>(() => {
  const rows = status.value?.api_key_statuses?.length
    ? status.value.api_key_statuses
    : configForm.api_key_statuses
  return Array.isArray(rows) ? rows : []
placeholder)

const apiKeyRows = computed<ContentModerationAPIKeyStatus[]>(() => [
  ...savedApiKeyRows.value,
  ...testedApiKeyStatuses.value,
])

const apiKeyHealthBadges = computed<Array<{ status: ContentModerationAPIKeyStatus['status']; count: number placeholder>>(() => {
  const counts: Record<ContentModerationAPIKeyStatus['status'], number> = {
    ok: 0,
    error: 0,
    frozen: 0,
    unknown: 0,
  placeholder
  for (const row of savedApiKeyRows.value) {
    counts[row.status] = (counts[row.status] ?? 0) + 1
  placeholder
  if (savedApiKeyRows.value.length === 0 && configForm.api_key_count > 0) {
    counts.unknown = configForm.api_key_count
  placeholder
  return (['ok', 'frozen', 'error', 'unknown'] as Array<ContentModerationAPIKeyStatus['status']>)
    .map((item) => ({ status: item, count: counts[item] placeholder))
    .filter((item) => item.count > 0)
placeholder)

const apiKeyHealthSummary = computed(() => {
  if (!configForm.api_key_configured) return ''
  if (apiKeyHealthBadges.value.length === 0) return t('admin.riskControl.apiKeyStatusUnknown')
  return apiKeyHealthBadges.value
    .map((badge) => `${apiKeyStatusLabel(badge.status)placeholder ${badge.countplaceholder`)
    .join(' · ')
placeholder)

const overviewItems = computed<OverviewItem[]>(() => [
  {
    key: 'status',
    label: t('admin.riskControl.overview.status'),
    value: configForm.enabled ? t('admin.riskControl.overview.enabled') : t('admin.riskControl.overview.disabled'),
    meta: modeLabel(configForm.mode),
    icon: 'shield',
    iconClass: configForm.enabled
      ? 'bg-emerald-50 text-emerald-600 dark:bg-emerald-900/20 dark:text-emerald-300'
      : 'bg-gray-100 text-gray-500 dark:bg-dark-700 dark:text-gray-400',
    badge: runtimeBadgeText.value,
    badgeClass: runtimeBadgeClass.value,
  placeholder,
  {
    key: 'api-key',
    label: t('admin.riskControl.overview.apiKey'),
    value: configForm.api_key_configured ? t('admin.riskControl.apiKeyCount', { count: configForm.api_key_count placeholder) : t('admin.riskControl.notConfigured'),
    meta: configForm.api_key_configured ? apiKeyHealthSummary.value || configForm.model || '-' : configForm.model || '-',
    icon: 'key',
    iconClass: 'bg-sky-50 text-sky-600 dark:bg-sky-900/20 dark:text-sky-300',
  placeholder,
  {
    key: 'scope',
    label: t('admin.riskControl.overview.groupScope'),
    value: configForm.all_groups ? t('admin.riskControl.allGroups') : selectedGroupCount.value,
    meta: configForm.all_groups ? t('admin.riskControl.allGroupsHint') : t('admin.riskControl.selectedGroupsHint'),
    icon: 'users',
    iconClass: 'bg-violet-50 text-violet-600 dark:bg-violet-900/20 dark:text-violet-300',
  placeholder,
  {
    key: 'logs',
    label: t('admin.riskControl.overview.logs'),
    value: formatNumber(pagination.total),
    meta: t('admin.riskControl.overview.currentFilter'),
    icon: 'document',
    iconClass: 'bg-amber-50 text-amber-600 dark:bg-amber-900/20 dark:text-amber-300',
  placeholder,
])

const moderationScoreRows = computed<ModerationScoreRow[]>(() => {
  const result = moderationTestResult.value
  if (!result) return []
  return Object.entries(result.category_scores || {placeholder)
    .map(([category, score]) => {
      const threshold = result.thresholds?.[category] ?? 1
      return {
        category,
        score,
        threshold,
        hit: score >= threshold,
      placeholder
    placeholder)
    .sort((a, b) => b.score - a.score)
placeholder)

const inputDetailText = computed(() => {
  if (!inputDetailRow.value) return '-'
  return inputDetailRow.value.input_excerpt || inputDetailRow.value.error || '-'
placeholder)

const queueUsagePercent = computed(() => `${Math.min(100, Math.max(0, status.value?.queue_usage_percent ?? 0)).toFixed(1)placeholder%`)

const queueUsageStyle = computed(() => ({
  width: queueUsagePercent.value,
placeholder))

const workerSlots = computed(() => {
  const total = Math.max(0, status.value?.worker_count ?? configForm.worker_count)
  const active = Math.max(0, status.value?.active_workers ?? 0)
  const enabled = Boolean(status.value?.risk_control_enabled && status.value?.enabled && status.value?.mode !== 'off')
  return Array.from({ length: total placeholder, (_, index) => ({
    id: index + 1,
    state: (!enabled ? 'disabled' : index < active ? 'active' : 'idle') as WorkerSlotState,
    label: !enabled
      ? t('admin.riskControl.workerDisabled')
      : index < active
        ? t('admin.riskControl.workerActive')
        : t('admin.riskControl.workerIdle'),
  placeholder))
placeholder)

const runtimeBadgeText = computed(() => {
  if (!status.value?.risk_control_enabled) return t('admin.riskControl.riskSwitchOff')
  if (!configForm.enabled || configForm.mode === 'off') return t('admin.riskControl.overview.disabled')
  return t('admin.riskControl.overview.enabled')
placeholder)

const runtimeBadgeClass = computed(() => {
  if (!status.value?.risk_control_enabled || !configForm.enabled || configForm.mode === 'off') {
    return 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-gray-300'
  placeholder
  return 'bg-emerald-50 text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300'
placeholder)

function applyConfig(config: ContentModerationConfig) {
  configForm.enabled = config.enabled
  configForm.mode = config.mode
  configForm.base_url = config.base_url || 'https://api.openai.com'
  configForm.model = config.model || 'omni-moderation-latest'
  configForm.api_keys_text = ''
  configForm.api_key_configured = config.api_key_configured
  configForm.api_key_masked = config.api_key_masked || ''
  configForm.api_key_count = config.api_key_count || 0
  configForm.api_key_masks = Array.isArray(config.api_key_masks) ? [...config.api_key_masks] : []
  configForm.api_key_statuses = Array.isArray(config.api_key_statuses) ? [...config.api_key_statuses] : []
  configForm.clear_api_key = false
  testedApiKeyStatuses.value = []
  configForm.timeout_ms = config.timeout_ms || 3000
  configForm.retry_count = config.retry_count ?? 2
  configForm.sample_rate = config.sample_rate ?? 100
  configForm.all_groups = config.all_groups
  configForm.group_ids = Array.isArray(config.group_ids) ? [...config.group_ids] : []
  configForm.record_non_hits = config.record_non_hits
  configForm.worker_count = config.worker_count || 4
  configForm.queue_size = config.queue_size || 32768
  configForm.block_status = config.block_status || 403
  configForm.block_message = config.block_message || '内容审计命中风险规则，请调整输入后重试'
  configForm.email_on_hit = config.email_on_hit ?? true
  configForm.auto_ban_enabled = config.auto_ban_enabled ?? true
  configForm.ban_threshold = config.ban_threshold || 10
  configForm.violation_window_hours = config.violation_window_hours || 720
  configForm.hit_retention_days = config.hit_retention_days || 180
  configForm.non_hit_retention_days = Math.min(Math.max(config.non_hit_retention_days || 3, 1), 3)
  configForm.pre_hash_check_enabled = config.pre_hash_check_enabled ?? false
placeholder

async function loadAll() {
  loading.value = true
  try {
    const [config, groupItems, runtimeStatus] = await Promise.all([
      adminAPI.riskControl.getConfig(),
      adminAPI.groups.getAll(),
      adminAPI.riskControl.getStatus(),
    ])
    applyConfig(config)
    groups.value = groupItems
    status.value = runtimeStatus
    if (Array.isArray(runtimeStatus.api_key_statuses)) {
      configForm.api_key_statuses = [...runtimeStatus.api_key_statuses]
    placeholder
    await loadLogs()
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('admin.riskControl.loadFailed')))
  placeholder finally {
    loading.value = false
  placeholder
placeholder

async function loadStatus(silent = true) {
  statusLoading.value = true
  try {
    const runtimeStatus = await adminAPI.riskControl.getStatus()
    status.value = runtimeStatus
    if (Array.isArray(runtimeStatus.api_key_statuses)) {
      configForm.api_key_statuses = [...runtimeStatus.api_key_statuses]
    placeholder
  placeholder catch (err: unknown) {
    if (!silent) {
      appStore.showError(extractApiErrorMessage(err, t('admin.riskControl.statusFailed')))
    placeholder
  placeholder finally {
    statusLoading.value = false
  placeholder
placeholder

async function saveConfig() {
  saving.value = true
  try {
    const payload: UpdateContentModerationConfig = {
      enabled: configForm.enabled,
      mode: configForm.mode,
      base_url: configForm.base_url,
      model: configForm.model,
      timeout_ms: Number(configForm.timeout_ms) || 3000,
      retry_count: Number(configForm.retry_count) || 0,
      sample_rate: Number(configForm.sample_rate) || 0,
      all_groups: configForm.all_groups,
      group_ids: configForm.all_groups ? [] : [...configForm.group_ids],
      record_non_hits: configForm.record_non_hits,
      clear_api_key: configForm.clear_api_key,
      worker_count: Number(configForm.worker_count) || 4,
      queue_size: Number(configForm.queue_size) || 32768,
      block_status: Number(configForm.block_status) || 403,
      block_message: configForm.block_message || '内容审计命中风险规则，请调整输入后重试',
      email_on_hit: configForm.email_on_hit,
      auto_ban_enabled: configForm.auto_ban_enabled,
      ban_threshold: Number(configForm.ban_threshold) || 10,
      violation_window_hours: Number(configForm.violation_window_hours) || 720,
      hit_retention_days: Number(configForm.hit_retention_days) || 180,
      non_hit_retention_days: Math.min(Math.max(Number(configForm.non_hit_retention_days) || 3, 1), 3),
      pre_hash_check_enabled: configForm.pre_hash_check_enabled,
    placeholder
    const keys = parseApiKeys(configForm.api_keys_text)
    if (keys.length > 0) {
      payload.api_keys = keys
      payload.clear_api_key = false
    placeholder

    const updated = await adminAPI.riskControl.updateConfig(payload)
    applyConfig(updated)
    settingsOpen.value = false
    appStore.showSuccess(t('admin.riskControl.saved'))
    await Promise.all([loadStatus(true), loadLogs()])
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('admin.riskControl.saveFailed')))
  placeholder finally {
    saving.value = false
  placeholder
placeholder

async function loadLogs() {
  logsLoading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
      result: filters.result || undefined,
      group_id: filters.group_id || undefined,
      endpoint: filters.endpoint || undefined,
      search: filters.search || undefined,
      from: normalizeDateTimeLocal(filters.from),
      to: normalizeDateTimeLocal(filters.to),
    placeholder
    const result = await adminAPI.riskControl.listLogs(params)
    logs.value = result.items
    pagination.total = result.total
    pagination.page = result.page
    pagination.page_size = result.page_size
    pagination.pages = result.pages
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('admin.riskControl.logsFailed')))
  placeholder finally {
    logsLoading.value = false
  placeholder
placeholder

function canUnbanRow(row: ContentModerationLog): boolean {
  return Boolean(row.auto_banned && row.user_id && row.user_status === 'disabled')
placeholder

function inputSummaryText(row: ContentModerationLog): string {
  return row.input_excerpt || row.error || '-'
placeholder

function openInputDetail(row: ContentModerationLog) {
  inputDetailRow.value = row
placeholder

function closeInputDetail() {
  inputDetailRow.value = null
placeholder

async function unbanUser(row: ContentModerationLog) {
  if (!row.user_id || unbanningUserID.value !== null) return
  unbanningUserID.value = row.user_id
  try {
    const result = await adminAPI.riskControl.unbanUser(row.user_id)
    logs.value = logs.value.map((item) => {
      if (item.user_id !== row.user_id) return item
      return { ...item, user_status: result.status placeholder
    placeholder)
    appStore.showSuccess(t('admin.riskControl.unbanSuccess'))
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('admin.riskControl.unbanFailed')))
  placeholder finally {
    unbanningUserID.value = null
  placeholder
placeholder

async function deleteFlaggedHash() {
  if (!isFlaggedHashInputValid.value || hashActionLoading.value) return
  hashActionLoading.value = true
  try {
    const result = await adminAPI.riskControl.deleteFlaggedHash(flaggedHashInput.value)
    flaggedHashInput.value = ''
    await loadStatus(true)
    appStore.showSuccess(result.deleted ? t('admin.riskControl.flaggedHashDeleted') : t('admin.riskControl.flaggedHashNotFound'))
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('admin.riskControl.flaggedHashDeleteFailed')))
  placeholder finally {
    hashActionLoading.value = false
  placeholder
placeholder

async function clearFlaggedHashes() {
  if (hashActionLoading.value) return
  const confirmed = window.confirm(t('admin.riskControl.clearFlaggedHashesConfirm'))
  if (!confirmed) return
  hashActionLoading.value = true
  try {
    const result = await adminAPI.riskControl.clearFlaggedHashes()
    await loadStatus(true)
    appStore.showSuccess(t('admin.riskControl.flaggedHashesCleared', { count: result.deleted placeholder))
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('admin.riskControl.flaggedHashesClearFailed')))
  placeholder finally {
    hashActionLoading.value = false
  placeholder
placeholder

function openSettings() {
  activeSettingsTab.value = 'basic'
  settingsOpen.value = true
placeholder

function reloadLogsFromFirstPage() {
  pagination.page = 1
  void loadLogs()
placeholder

function onPageChange(page: number) {
  pagination.page = page
  void loadLogs()
placeholder

function onPageSizeChange(pageSize: number) {
  pagination.page = 1
  pagination.page_size = pageSize
  void loadLogs()
placeholder

function toggleClearApiKey() {
  configForm.clear_api_key = !configForm.clear_api_key
  if (configForm.clear_api_key) {
    configForm.api_keys_text = ''
    testedApiKeyStatuses.value = []
  placeholder
placeholder

async function testApiKeys(useInputKeys: boolean) {
  const keys = useInputKeys ? parseApiKeys(configForm.api_keys_text) : []
  if (useInputKeys && keys.length === 0) {
    appStore.showError(t('admin.riskControl.apiKeyTestNoInput'))
    return
  placeholder
  apiKeyTesting.value = true
  try {
    const result = await adminAPI.riskControl.testAPIKeys({
      api_keys: keys,
      base_url: configForm.base_url,
      model: configForm.model,
      timeout_ms: Number(configForm.timeout_ms) || 3000,
      prompt: moderationTestPrompt.value,
      images: moderationTestImages.value,
    placeholder)
    moderationTestResult.value = result.audit_result ?? null
    if (useInputKeys) {
      testedApiKeyStatuses.value = result.items.map((item) => ({ ...item, configured: false placeholder))
    placeholder else {
      mergeConfiguredAPIKeyStatuses(result.items)
      testedApiKeyStatuses.value = []
      await loadStatus(true)
    placeholder
    appStore.showSuccess(t('admin.riskControl.apiKeyTestDone', { count: result.items.length placeholder))
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('admin.riskControl.apiKeyTestFailed')))
  placeholder finally {
    apiKeyTesting.value = false
  placeholder
placeholder

function mergeConfiguredAPIKeyStatuses(items: ContentModerationAPIKeyStatus[]) {
  if (!hasModerationAuditInput.value || configForm.api_key_statuses.length === 0) {
    configForm.api_key_statuses = items
    return
  placeholder
  const updates = new Map(items.map((item) => [item.key_hash, item]))
  configForm.api_key_statuses = configForm.api_key_statuses.map((item) => updates.get(item.key_hash) ?? item)
placeholder

function clearModerationTestInput() {
  moderationTestPrompt.value = ''
  moderationTestImages.value = []
  moderationTestResult.value = null
placeholder

function removeModerationTestImage(index: number) {
  moderationTestImages.value.splice(index, 1)
placeholder

async function handleModerationImageUpload(event: Event) {
  const input = event.target as HTMLInputElement
  await addModerationTestFiles(input.files)
  input.value = ''
placeholder

async function handleModerationImageDrop(event: DragEvent) {
  await addModerationTestFiles(event.dataTransfer?.files ?? null)
placeholder

async function handleModerationImagePaste(event: ClipboardEvent) {
  const files = Array.from(event.clipboardData?.files ?? []).filter((file) => file.type.startsWith('image/'))
  if (files.length === 0) return
  event.preventDefault()
  await addModerationTestFiles(files)
placeholder

async function addModerationTestFiles(files: FileList | File[] | null) {
  if (!files) return
  const items = Array.from(files).filter((file) => file.type.startsWith('image/'))
  for (const file of items) {
    if (moderationTestImages.value.length >= maxModerationTestImages) {
      appStore.showError(t('admin.riskControl.auditTestImageLimit', { count: maxModerationTestImages placeholder))
      return
    placeholder
    if (file.size > maxModerationTestImageSize) {
      appStore.showError(t('admin.riskControl.auditTestImageTooLarge'))
      continue
    placeholder
    try {
      moderationTestImages.value.push(await fileToDataURL(file))
    placeholder catch {
      appStore.showError(t('admin.riskControl.auditTestImageReadFailed'))
    placeholder
  placeholder
placeholder

function fileToDataURL(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(String(reader.result || ''))
    reader.onerror = () => reject(reader.error)
    reader.readAsDataURL(file)
  placeholder)
placeholder

function toggleGroup(groupID: number) {
  const index = configForm.group_ids.indexOf(groupID)
  if (index >= 0) {
    configForm.group_ids.splice(index, 1)
  placeholder else {
    configForm.group_ids.push(groupID)
  placeholder
placeholder

function isGroupSelected(groupID: number): boolean {
  return configForm.group_ids.includes(groupID)
placeholder

function modeLabel(mode: ModerationMode): string {
  const found = modeOptions.value.find((option) => option.value === mode)
  return found?.label ?? mode
placeholder

function modeDescription(mode: ModerationMode): string {
  const descriptions: Record<ModerationMode, string> = {
    pre_block: t('admin.riskControl.modePreBlockDesc'),
    observe: t('admin.riskControl.modeObserveDesc'),
    off: t('admin.riskControl.modeOffDesc'),
  placeholder
  return descriptions[mode] ?? ''
placeholder

function resultLabel(row: ContentModerationLog): string {
  if (row.action === 'block') return t('admin.riskControl.action.block')
  if (row.action === 'error' || row.error) return t('admin.riskControl.action.error')
  if (row.flagged) return t('admin.riskControl.result.hit')
  return t('admin.riskControl.result.pass')
placeholder

function resultBadgeClass(row: ContentModerationLog): string {
  if (row.action === 'block') return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300'
  if (row.action === 'error' || row.error) return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300'
  if (row.flagged) return 'bg-pink-100 text-pink-700 dark:bg-pink-900/30 dark:text-pink-300'
  return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300'
placeholder

function workerSlotClass(state: WorkerSlotState): string {
  if (state === 'active') {
    return 'border-sky-200 bg-sky-50 text-sky-700 dark:border-sky-900/60 dark:bg-sky-900/20 dark:text-sky-300'
  placeholder
  if (state === 'idle') {
    return 'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-900/60 dark:bg-emerald-900/20 dark:text-emerald-300'
  placeholder
  return 'border-gray-100 bg-white text-gray-400 dark:border-dark-700 dark:bg-dark-800 dark:text-gray-500'
placeholder

function workerDotClass(state: WorkerSlotState): string {
  if (state === 'active') return 'bg-sky-500'
  if (state === 'idle') return 'bg-emerald-500'
  return 'bg-gray-300 dark:bg-dark-500'
placeholder

function percent(value: number): string {
  if (!Number.isFinite(value)) return '-'
  return `${(value * 100).toFixed(1)placeholder%`
placeholder

function percentWidth(value: number): string {
  if (!Number.isFinite(value)) return '0%'
  return `${Math.min(100, Math.max(0, value * 100)).toFixed(1)placeholder%`
placeholder

function latencyText(value: number | null): string {
  if (value === null || value === undefined) return '-'
  return `${valueplaceholder ms`
placeholder

function apiKeyRowKey(row: ContentModerationAPIKeyStatus, index: number): string {
  return `${row.configured ? 'saved' : 'test'placeholder-${row.key_hash || indexplaceholder`
placeholder

function apiKeyStatusLabel(statusValue: ContentModerationAPIKeyStatus['status']): string {
  const labels: Record<ContentModerationAPIKeyStatus['status'], string> = {
    ok: t('admin.riskControl.apiKeyStatusOk'),
    error: t('admin.riskControl.apiKeyStatusError'),
    frozen: t('admin.riskControl.apiKeyStatusFrozen'),
    unknown: t('admin.riskControl.apiKeyStatusUnknown'),
  placeholder
  return labels[statusValue] ?? labels.unknown
placeholder

function apiKeyStatusBadgeClass(statusValue: ContentModerationAPIKeyStatus['status']): string {
  const classes: Record<ContentModerationAPIKeyStatus['status'], string> = {
    ok: 'bg-emerald-50 text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300',
    error: 'bg-amber-50 text-amber-700 dark:bg-amber-900/20 dark:text-amber-300',
    frozen: 'bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-300',
    unknown: 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-gray-300',
  placeholder
  return classes[statusValue] ?? classes.unknown
placeholder

function apiKeyStatusDotClass(statusValue: ContentModerationAPIKeyStatus['status']): string {
  const classes: Record<ContentModerationAPIKeyStatus['status'], string> = {
    ok: 'bg-emerald-500',
    error: 'bg-amber-500',
    frozen: 'bg-red-500',
    unknown: 'bg-gray-400',
  placeholder
  return classes[statusValue] ?? classes.unknown
placeholder

function apiKeyStatusMeta(row: ContentModerationAPIKeyStatus): string {
  const parts: string[] = []
  parts.push(t('admin.riskControl.apiKeyFailureCount', { count: row.failure_count || 0 placeholder))
  if (row.last_latency_ms > 0) {
    parts.push(t('admin.riskControl.apiKeyLatency', { ms: row.last_latency_ms placeholder))
  placeholder
  if (row.last_http_status > 0) {
    parts.push(t('admin.riskControl.apiKeyHTTPStatus', { status: row.last_http_status placeholder))
  placeholder
  if (row.frozen_until) {
    parts.push(t('admin.riskControl.apiKeyFrozenUntil', { time: formatDateTime(row.frozen_until) placeholder))
  placeholder else if (row.last_checked_at) {
    parts.push(t('admin.riskControl.apiKeyLastChecked', { time: formatDateTime(row.last_checked_at) placeholder))
  placeholder else {
    parts.push(t('admin.riskControl.apiKeyNotTested'))
  placeholder
  return parts.join(' / ')
placeholder

function parseApiKeys(value: string): string[] {
  return value
    .split(/\r?\n/)
    .map((item) => item.trim())
    .filter((item, index, arr) => item && arr.indexOf(item) === index)
placeholder

function violationCountText(row: ContentModerationLog): string {
  if (!row.flagged) return '-'
  return t('admin.riskControl.violationCount', { count: row.violation_count || 1 placeholder)
placeholder

function normalizeDateTimeLocal(value: string): string | undefined {
  if (!value) return undefined
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return undefined
  return date.toISOString()
placeholder

function formatDateTime(value: string): string {
  return formatDateTimeValue(value) || '-'
placeholder

function formatNumber(value: number): string {
  return new Intl.NumberFormat().format(value)
placeholder

onMounted(() => {
  void loadAll()
  statusTimer = window.setInterval(() => {
    void loadStatus(true)
  placeholder, 15000)
placeholder)

onUnmounted(() => {
  if (statusTimer !== null) {
    window.clearInterval(statusTimer)
    statusTimer = null
  placeholder
placeholder)
</script>
