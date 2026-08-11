import { afterEach, beforeEach, describe, expect, it, vi placeholder from 'vitest'
import { flushPromises, mount placeholder from '@vue/test-utils'

import BackupView from '../BackupView.vue'

const {
  getS3Config,
  getImageStorageConfig,
  getSchedule,
  listBackups,
  getDownloadURL,
placeholder = vi.hoisted(() => ({
  getS3Config: vi.fn(),
  getImageStorageConfig: vi.fn(),
  getSchedule: vi.fn(),
  listBackups: vi.fn(),
  getDownloadURL: vi.fn(),
placeholder))

vi.mock('@/api', () => ({
  adminAPI: {
    backup: {
      getS3Config,
      updateS3Config: vi.fn(),
      testS3Connection: vi.fn(),
      getImageStorageConfig,
      updateImageStorageConfig: vi.fn(),
      testImageStorageConnection: vi.fn(),
      getSchedule,
      updateSchedule: vi.fn(),
      createBackup: vi.fn(),
      listBackups,
      getBackup: vi.fn(),
      deleteBackup: vi.fn(),
      getDownloadURL,
      restoreBackup: vi.fn(),
    placeholder,
  placeholder,
placeholder))

vi.mock('@/stores', () => ({
  useAppStore: () => ({
    showError: vi.fn(),
    showSuccess: vi.fn(),
    showWarning: vi.fn(),
  placeholder),
placeholder))

vi.mock('@/composables/useStepUp', () => ({
  useStepUp: () => ({ run: (fn: () => unknown) => fn() placeholder),
  isStepUpBlocked: () => false,
  isStepUpCancelled: () => false,
  stepUpBlockReason: () => '',
placeholder))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string, params?: Record<string, unknown>) =>
      params?.index === undefined ? key : `${keyplaceholder:${params.indexplaceholder`,
  placeholder),
placeholder))

const baseRecord = (id: string, parts?: unknown[]) => ({
  id,
  status: 'completed',
  backup_type: 'postgres',
  file_name: `${idplaceholder.sql.gz`,
  s3_key: `backups/${idplaceholder.sql.gz`,
  parts,
  size_bytes: 10,
  triggered_by: 'manual',
  started_at: '2026-08-09T00:00:00Z',
placeholder)

function mountBackupView() {
  return mount(BackupView, {
    global: {
      stubs: {
        TotpStepUpDialog: true,
      placeholder,
    placeholder,
  placeholder)
placeholder

describe('admin BackupView 分卷备份', () => {
  beforeEach(() => {
    getS3Config.mockResolvedValue({placeholder)
    getImageStorageConfig.mockResolvedValue({ config: {placeholder, secret_configured: false placeholder)
    getSchedule.mockResolvedValue({ enabled: false, cron_expr: '', retain_days: 14, retain_count: 10 placeholder)
    getDownloadURL.mockReset()
    vi.spyOn(HTMLAnchorElement.prototype, 'click').mockImplementation(() => {placeholder)
  placeholder)

  afterEach(() => {
    vi.restoreAllMocks()
    document.body.innerHTML = ''
  placeholder)

  it('显示分卷数并在下载时列出每个分卷链接', async () => {
    listBackups.mockResolvedValue({
      items: [baseRecord('split', [{ index: 1 placeholder, { index: 2 placeholder, { index: 3 placeholder])],
    placeholder)
    getDownloadURL.mockResolvedValue({
      parts: [
        { index: 1, size_bytes: 5, url: 'https://example.test/part-1' placeholder,
        { index: 2, size_bytes: 6, url: 'https://example.test/part-2' placeholder,
        { index: 3, size_bytes: 7, url: 'https://example.test/part-3' placeholder,
      ],
    placeholder)

    const wrapper = mountBackupView()
    await flushPromises()

    expect(wrapper.text()).toContain('3')
    const downloadButton = wrapper.findAll('button').find(button =>
      button.text().includes('admin.backup.actions.download'),
    )
    expect(downloadButton).toBeDefined()
    await downloadButton!.trigger('click')
    await flushPromises()

    expect(document.body.textContent).toContain('admin.backup.actions.partLabel:1')
    expect(document.body.textContent).toContain('admin.backup.actions.partLabel:3')
    expect(document.body.querySelector('a[href="https://example.test/part-2"]')).not.toBeNull()
  placeholder)

  it('旧单文件记录仍使用单个下载地址', async () => {
    listBackups.mockResolvedValue({ items: [baseRecord('legacy')] placeholder)
    getDownloadURL.mockResolvedValue({ url: 'https://example.test/legacy.sql.gz' placeholder)

    const wrapper = mountBackupView()
    await flushPromises()
    const downloadButton = wrapper.findAll('button').find(button =>
      button.text().includes('admin.backup.actions.download'),
    )
    await downloadButton!.trigger('click')
    await flushPromises()

    expect(getDownloadURL).toHaveBeenCalledWith('legacy')
    expect(document.body.textContent).not.toContain('admin.backup.actions.downloadParts')
  placeholder)

  it('运行中的备份不显示删除入口', async () => {
    listBackups.mockResolvedValue({
      items: [{ ...baseRecord('running'), status: 'running', progress: 'uploading' placeholder],
    placeholder)

    const wrapper = mountBackupView()
    await flushPromises()

    expect(wrapper.find('tbody tr td:nth-child(5)').text()).toBe('-')
    expect(wrapper.findAll('button').some(button => button.text() === 'common.delete')).toBe(false)
  placeholder)
placeholder)
