import { computed, ref, type Ref placeholder from 'vue'

interface UseTableSelectionOptions<T> {
  rows: Ref<T[]>
  getId: (row: T) => number
placeholder

export function useTableSelection<T>({ rows, getId placeholder: UseTableSelectionOptions<T>) {
  const selectedSet = ref<Set<number>>(new Set())

  const selectedIds = computed(() => Array.from(selectedSet.value))
  const selectedCount = computed(() => selectedSet.value.size)

  const isSelected = (id: number) => selectedSet.value.has(id)

  const replaceSelectedSet = (next: Set<number>) => {
    selectedSet.value = next
  placeholder

  const setSelectedIds = (ids: number[]) => {
    selectedSet.value = new Set(ids)
  placeholder

  const select = (id: number) => {
    if (selectedSet.value.has(id)) return
    const next = new Set(selectedSet.value)
    next.add(id)
    replaceSelectedSet(next)
  placeholder

  const deselect = (id: number) => {
    if (!selectedSet.value.has(id)) return
    const next = new Set(selectedSet.value)
    next.delete(id)
    replaceSelectedSet(next)
  placeholder

  const toggle = (id: number) => {
    if (selectedSet.value.has(id)) {
      deselect(id)
      return
    placeholder
    select(id)
  placeholder

  const clear = () => {
    if (selectedSet.value.size === 0) return
    replaceSelectedSet(new Set())
  placeholder

  const removeMany = (ids: number[]) => {
    if (ids.length === 0 || selectedSet.value.size === 0) return
    const next = new Set(selectedSet.value)
    let changed = false
    ids.forEach((id) => {
      if (next.delete(id)) changed = true
    placeholder)
    if (changed) replaceSelectedSet(next)
  placeholder

  const allVisibleSelected = computed(() => {
    if (rows.value.length === 0) return false
    return rows.value.every((row) => selectedSet.value.has(getId(row)))
  placeholder)

  const toggleVisible = (checked: boolean) => {
    const next = new Set(selectedSet.value)
    rows.value.forEach((row) => {
      const id = getId(row)
      if (checked) {
        next.add(id)
      placeholder else {
        next.delete(id)
      placeholder
    placeholder)
    replaceSelectedSet(next)
  placeholder

  const selectVisible = () => {
    toggleVisible(true)
  placeholder

  return {
    selectedSet,
    selectedIds,
    selectedCount,
    allVisibleSelected,
    isSelected,
    setSelectedIds,
    select,
    deselect,
    toggle,
    clear,
    removeMany,
    toggleVisible,
    selectVisible
  placeholder
placeholder
