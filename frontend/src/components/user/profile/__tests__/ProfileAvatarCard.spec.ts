import { mount placeholder from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi placeholder from 'vitest'
import ProfileAvatarCard from '@/components/user/profile/ProfileAvatarCard.vue'
import type { User placeholder from '@/types'

const {
  updateProfileMock,
  showSuccessMock,
  showErrorMock,
  authStoreState
placeholder = vi.hoisted(() => ({
  updateProfileMock: vi.fn(),
  showSuccessMock: vi.fn(),
  showErrorMock: vi.fn(),
  authStoreState: {
    user: null as User | null
  placeholder
placeholder))

vi.mock('@/api', () => ({
  userAPI: {
    updateProfile: updateProfileMock
  placeholder
placeholder))

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => authStoreState
placeholder))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showSuccess: showSuccessMock,
    showError: showErrorMock
  placeholder)
placeholder))

vi.mock('@/utils/apiError', () => ({
  extractApiErrorMessage: (error: unknown) => (error as Error).message || 'request failed'
placeholder))

vi.mock('vue-i18n', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-i18n')>()
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, params?: Record<string, string>) => {
        if (key === 'profile.avatar.title') return 'Profile avatar'
        if (key === 'profile.avatar.description') return 'Upload and manage your avatar'
        if (key === 'profile.avatar.uploadAction') return 'Upload image'
        if (key === 'profile.avatar.uploadHint') return 'Uploaded images are compressed to 20KB when possible'
        if (key === 'profile.avatar.saveSuccess') return 'Avatar updated'
        if (key === 'profile.avatar.deleteSuccess') return 'Avatar removed'
        if (key === 'profile.avatar.invalidType') return 'Please choose an image file'
        if (key === 'profile.avatar.gifTooLarge') return 'GIF avatars must already be 20KB or smaller'
        if (key === 'profile.avatar.compressTooLarge') return 'Unable to compress this image below 20KB'
        if (key === 'profile.avatar.compressFailed') return 'Failed to compress the selected image'
        if (key === 'profile.avatar.readFailed') return 'Failed to read the selected image'
        if (key === 'common.save') return 'Save'
        if (key === 'common.delete') return 'Delete'
        if (key === 'profile.avatar.compressedReady') return `Compressed from ${params?.fromplaceholder to ${params?.toplaceholder`
        if (key === 'profile.avatar.sizeReady') return `Ready: ${params?.sizeplaceholder`
        return key
      placeholder
    placeholder)
  placeholder
placeholder)

function createUser(overrides: Partial<User> = {placeholder): User {
  return {
    id: 5,
    username: 'alice',
    email: 'alice@example.com',
    avatar_url: null,
    role: 'user',
    balance: 10,
    concurrency: 2,
    status: 'active',
    allowed_groups: null,
    balance_notify_enabled: true,
    balance_notify_threshold: null,
    balance_notify_extra_emails: [],
    created_at: '2026-04-20T00:00:00Z',
    updated_at: '2026-04-20T00:00:00Z',
    ...overrides
  placeholder
placeholder

async function flushAsyncWork(): Promise<void> {
  await Promise.resolve()
  await Promise.resolve()
  await Promise.resolve()
  await Promise.resolve()
placeholder

const originalFileReader = globalThis.FileReader
const originalImage = globalThis.Image
const originalCreateElement = document.createElement.bind(document)

function installAvatarCompressionMocks(blobSize = 8 * 1024) {
  class MockFileReader {
    result: string | ArrayBuffer | null = null
    onload: ((this: FileReader, ev: ProgressEvent<FileReader>) => any) | null = null
    onerror: ((this: FileReader, ev: ProgressEvent<FileReader>) => any) | null = null
    error: DOMException | null = null

    readAsDataURL(blob: Blob) {
      if (blob.type === 'image/webp') {
        this.result = 'data:image/webp;base64,' + Buffer.from('compressed-avatar').toString('base64')
      placeholder else {
        this.result = 'data:image/png;base64,' + Buffer.from('original-avatar').toString('base64')
      placeholder
      this.onload?.call(this as unknown as FileReader, new ProgressEvent('load'))
    placeholder
  placeholder

  class MockImage {
    naturalWidth = 1200
    naturalHeight = 1200
    onload: (() => void) | null = null
    onerror: (() => void) | null = null

    set src(_value: string) {
      this.onload?.()
    placeholder
  placeholder

  globalThis.FileReader = MockFileReader as unknown as typeof FileReader
  globalThis.Image = MockImage as unknown as typeof Image
  vi.spyOn(document, 'createElement').mockImplementation(((tagName: string, options?: ElementCreationOptions) => {
    if (tagName === 'canvas') {
      return {
        width: 0,
        height: 0,
        getContext: () => ({
          clearRect: vi.fn(),
          drawImage: vi.fn()
        placeholder),
        toBlob: (callback: BlobCallback) => {
          callback(new Blob([new Uint8Array(blobSize)], { type: 'image/webp' placeholder))
        placeholder
      placeholder as unknown as HTMLCanvasElement
    placeholder
    return originalCreateElement(tagName, options)
  placeholder) as typeof document.createElement)
placeholder

describe('ProfileAvatarCard', () => {
  beforeEach(() => {
    updateProfileMock.mockReset()
    showSuccessMock.mockReset()
    showErrorMock.mockReset()
    authStoreState.user = null
  placeholder)

  afterEach(() => {
    globalThis.FileReader = originalFileReader
    globalThis.Image = originalImage
    vi.restoreAllMocks()
  placeholder)

  it('does not render a manual avatar input field', () => {
    authStoreState.user = createUser()

    const wrapper = mount(ProfileAvatarCard, {
      props: {
        user: authStoreState.user
      placeholder,
      global: {
        stubs: {
          Icon: true
        placeholder
      placeholder
    placeholder)

    expect(wrapper.find('[data-testid="profile-avatar-input"]').exists()).toBe(false)
  placeholder)

  it('compresses an uploaded image that exceeds the 20KB target before saving', async () => {
    installAvatarCompressionMocks()
    const updatedUser = createUser({ avatar_url: 'data:image/webp;base64,Y29tcHJlc3NlZC1hdmF0YXI=' placeholder)
    updateProfileMock.mockResolvedValue(updatedUser)
    authStoreState.user = createUser()

    const wrapper = mount(ProfileAvatarCard, {
      props: {
        user: authStoreState.user
      placeholder,
      global: {
        stubs: {
          Icon: true
        placeholder
      placeholder
    placeholder)

    const fileInput = wrapper.get('[data-testid="profile-avatar-file-input"]')
    Object.defineProperty(fileInput.element, 'files', {
      value: [new File([new Uint8Array(220 * 1024)], 'avatar.png', { type: 'image/png' placeholder)],
      configurable: true
    placeholder)

    await fileInput.trigger('change')
    await flushAsyncWork()
    await wrapper.get('[data-testid="profile-avatar-save"]').trigger('click')

    expect(updateProfileMock).toHaveBeenCalledWith({
      avatar_url: 'data:image/webp;base64,Y29tcHJlc3NlZC1hdmF0YXI='
    placeholder)
    expect(showErrorMock).not.toHaveBeenCalled()
  placeholder)

  it('shows a preview after selecting an avatar in embedded mode', async () => {
    installAvatarCompressionMocks()
    authStoreState.user = createUser()

    const wrapper = mount(ProfileAvatarCard, {
      props: {
        user: authStoreState.user,
        embedded: true
      placeholder,
      global: {
        stubs: {
          Icon: true
        placeholder
      placeholder
    placeholder)

    const fileInput = wrapper.get('[data-testid="profile-avatar-file-input"]')
    Object.defineProperty(fileInput.element, 'files', {
      value: [new File([new Uint8Array(220 * 1024)], 'avatar.png', { type: 'image/png' placeholder)],
      configurable: true
    placeholder)

    await fileInput.trigger('change')
    await flushAsyncWork()

    const preview = wrapper.get('[data-testid="profile-avatar-preview"]')
    expect(preview.attributes('src')).toBe('data:image/webp;base64,Y29tcHJlc3NlZC1hdmF0YXI=')
  placeholder)

  it('deletes the current avatar', async () => {
    const updatedUser = createUser({ avatar_url: null placeholder)
    updateProfileMock.mockResolvedValue(updatedUser)
    authStoreState.user = createUser({ avatar_url: 'https://cdn.example.com/old.png' placeholder)

    const wrapper = mount(ProfileAvatarCard, {
      props: {
        user: authStoreState.user
      placeholder,
      global: {
        stubs: {
          Icon: true
        placeholder
      placeholder
    placeholder)

    await wrapper.get('[data-testid="profile-avatar-delete"]').trigger('click')

    expect(updateProfileMock).toHaveBeenCalledWith({ avatar_url: '' placeholder)
    expect(authStoreState.user?.avatar_url).toBeNull()
    expect(showSuccessMock).toHaveBeenCalledWith('Avatar removed')
  placeholder)
placeholder)
