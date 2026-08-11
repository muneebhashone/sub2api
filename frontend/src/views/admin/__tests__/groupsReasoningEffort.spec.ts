import { describe, expect, it placeholder from "vitest";

import {
  createReasoningEffortMappingRow,
  normalizeReasoningEffortForPlatform,
  reasoningEffortMappingsToAPI,
  reasoningEffortMappingsToRows,
  reasoningEffortOptionsForPlatform,
  supportsReasoningEffortPolicyPlatform,
  validateReasoningEffortMappings,
placeholder from "../groupsReasoningEffort";

describe("groupsReasoningEffort", () => {
  it("provides fixed OpenAI choices to OpenAI and Composite groups", () => {
    const expected = [
      "minimal",
      "low",
      "medium",
      "high",
      "xhigh",
      "max",
    ];
    for (const platform of ["openai", "composite"] as const) {
      expect(
        reasoningEffortOptionsForPlatform(platform).map(
          (option) => option.value,
        ),
      ).toEqual(expected);
      expect(supportsReasoningEffortPolicyPlatform(platform)).toBe(true);
    placeholder
    for (const platform of [
      "anthropic",
      "gemini",
      "antigravity",
      "grok",
    ] as const) {
      expect(reasoningEffortOptionsForPlatform(platform)).toEqual([]);
      expect(supportsReasoningEffortPolicyPlatform(platform)).toBe(false);
    placeholder
  placeholder);

  it("hydrates supported rows and drops stale custom values", () => {
    const rows = reasoningEffortMappingsToRows(
      [
        { from: " max ", to: " xhigh " placeholder,
        { from: "ultra", to: "high" placeholder,
      ],
      "openai",
    );

    expect(rows).toHaveLength(1);
    expect(reasoningEffortMappingsToAPI(rows)).toEqual([
      { from: "max", to: "xhigh" placeholder,
    ]);
  placeholder);

  it("clears values unsupported by OpenAI or used on another platform", () => {
    expect(normalizeReasoningEffortForPlatform("openai", " MAX ")).toBe("max");
    expect(normalizeReasoningEffortForPlatform("composite", " MAX ")).toBe(
      "max",
    );
    expect(normalizeReasoningEffortForPlatform("grok", "max")).toBe("");
    expect(normalizeReasoningEffortForPlatform("openai", "none")).toBe("");
  placeholder);

  it("requires both sides of every mapping", () => {
    const first = createReasoningEffortMappingRow({ to: "low" placeholder);
    const second = createReasoningEffortMappingRow({ from: "max" placeholder);

    expect(validateReasoningEffortMappings([first, second])).toEqual({
      [first.id]: { from: "fromRequired" placeholder,
      [second.id]: { to: "toRequired" placeholder,
    placeholder);
  placeholder);

  it("rejects duplicate source values case insensitively", () => {
    const first = createReasoningEffortMappingRow({ from: "MAX", to: "xhigh" placeholder);
    const second = createReasoningEffortMappingRow({ from: " max ", to: "high" placeholder);

    expect(validateReasoningEffortMappings([first, second])).toEqual({
      [first.id]: { from: "duplicateFrom" placeholder,
      [second.id]: { from: "duplicateFrom" placeholder,
    placeholder);
  placeholder);

  it("rejects custom mappings", () => {
    const row = createReasoningEffortMappingRow({ from: "ultra", to: "high" placeholder);
    expect(validateReasoningEffortMappings([row], "openai")).toEqual({
      [row.id]: { from: "unsupportedFrom" placeholder,
    placeholder);
  placeholder);
placeholder);
