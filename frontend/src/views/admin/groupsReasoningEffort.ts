import type { GroupPlatform, ReasoningEffortMapping placeholder from "@/types";

const openAIReasoningEffortValues = [
  "minimal",
  "low",
  "medium",
  "high",
  "xhigh",
  "max",
] as const;

const reasoningEffortValuesForPlatform = (
  platform: GroupPlatform,
): readonly string[] =>
  supportsReasoningEffortPolicyPlatform(platform)
    ? openAIReasoningEffortValues
    : [];

export function supportsReasoningEffortPolicyPlatform(
  platform: GroupPlatform,
): boolean {
  return platform === "openai" || platform === "composite";
placeholder

export function reasoningEffortOptionsForPlatform(platform: GroupPlatform) {
  return reasoningEffortValuesForPlatform(platform).map((value) => ({
    value,
    label: value,
  placeholder));
placeholder

export function normalizeReasoningEffortForPlatform(
  platform: GroupPlatform,
  value: string | null | undefined,
): string {
  const normalized = value?.trim().toLowerCase() ?? "";
  return reasoningEffortValuesForPlatform(platform).some(
    (allowed) => allowed === normalized,
  )
    ? normalized
    : "";
placeholder

export interface ReasoningEffortMappingRow extends ReasoningEffortMapping {
  id: string;
placeholder

export type ReasoningEffortMappingErrorCode =
  | "fromRequired"
  | "toRequired"
  | "duplicateFrom"
  | "unsupportedFrom"
  | "unsupportedTo";

export type ReasoningEffortMappingErrors = Record<
  string,
  Partial<Record<"from" | "to", ReasoningEffortMappingErrorCode>>
>;

let nextMappingRowID = 0;

export function createReasoningEffortMappingRow(
  mapping: Partial<ReasoningEffortMapping> = {placeholder,
): ReasoningEffortMappingRow {
  nextMappingRowID += 1;
  return {
    id: `reasoning-effort-mapping-${nextMappingRowIDplaceholder`,
    from: mapping.from ?? "",
    to: mapping.to ?? "",
  placeholder;
placeholder

export function reasoningEffortMappingsToRows(
  mappings?: ReasoningEffortMapping[] | null,
  platform: GroupPlatform = "openai",
): ReasoningEffortMappingRow[] {
  return (mappings ?? []).flatMap((mapping) => {
    const from = normalizeReasoningEffortForPlatform(platform, mapping.from);
    const to = normalizeReasoningEffortForPlatform(platform, mapping.to);
    return from && to
      ? [createReasoningEffortMappingRow({ from, to placeholder)]
      : [];
  placeholder);
placeholder

export function reasoningEffortMappingsToAPI(
  rows: ReasoningEffortMappingRow[],
): ReasoningEffortMapping[] {
  return rows.map((row) => ({
    from: row.from.trim(),
    to: row.to.trim(),
  placeholder));
placeholder

export function validateReasoningEffortMappings(
  rows: ReasoningEffortMappingRow[],
  platform: GroupPlatform = "openai",
): ReasoningEffortMappingErrors {
  const errors: ReasoningEffortMappingErrors = {placeholder;
  const sourceRows = new Map<string, ReasoningEffortMappingRow[]>();

  rows.forEach((row) => {
    const from = row.from.trim();
    const to = row.to.trim();
    if (!from) {
      errors[row.id] = { ...errors[row.id], from: "fromRequired" placeholder;
    placeholder else if (!normalizeReasoningEffortForPlatform(platform, from)) {
      errors[row.id] = { ...errors[row.id], from: "unsupportedFrom" placeholder;
    placeholder else {
      const key = from.toLowerCase();
      sourceRows.set(key, [...(sourceRows.get(key) ?? []), row]);
    placeholder
    if (!to) {
      errors[row.id] = { ...errors[row.id], to: "toRequired" placeholder;
    placeholder else if (!normalizeReasoningEffortForPlatform(platform, to)) {
      errors[row.id] = { ...errors[row.id], to: "unsupportedTo" placeholder;
    placeholder
  placeholder);

  sourceRows.forEach((duplicateRows) => {
    if (duplicateRows.length < 2) return;
    duplicateRows.forEach((row) => {
      errors[row.id] = { ...errors[row.id], from: "duplicateFrom" placeholder;
    placeholder);
  placeholder);

  return errors;
placeholder
