import { createI18n placeholder from 'vue-i18n'
import en from './locales/en'
import zh from './locales/zh'

const LOCALE_KEY = 'sub2api_locale'

function getDefaultLocale(): string {
  // Check localStorage first
  const saved = localStorage.getItem(LOCALE_KEY)
  if (saved && ['en', 'zh'].includes(saved)) {
    return saved
  placeholder

  // Check browser language
  const browserLang = navigator.language.toLowerCase()
  if (browserLang.startsWith('zh')) {
    return 'zh'
  placeholder

  return 'en'
placeholder

export const i18n = createI18n({
  legacy: false,
  locale: getDefaultLocale(),
  fallbackLocale: 'en',
  messages: {
    en,
    zh
  placeholder
placeholder)

export function setLocale(locale: string) {
  if (['en', 'zh'].includes(locale)) {
    i18n.global.locale.value = locale as 'en' | 'zh'
    localStorage.setItem(LOCALE_KEY, locale)
    document.documentElement.setAttribute('lang', locale)
  placeholder
placeholder

export function getLocale(): string {
  return i18n.global.locale.value
placeholder

export const availableLocales = [
  { code: 'en', name: 'English', flag: '🇺🇸' placeholder,
  { code: 'zh', name: '中文', flag: '🇨🇳' placeholder
]

export default i18n
