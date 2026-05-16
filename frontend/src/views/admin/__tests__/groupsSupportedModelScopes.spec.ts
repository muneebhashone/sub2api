import { describe, expect, it placeholder from "vitest";
import { normalizeSupportedModelScopesForPlatform placeholder from "../groupsSupportedModelScopes";

describe("normalizeSupportedModelScopesForPlatform", () => {
  it("preserves model scopes for Antigravity groups", () => {
    expect(
      normalizeSupportedModelScopesForPlatform("antigravity", [
        "claude",
        "gemini_text",
      ]),
    ).toEqual(["claude", "gemini_text"]);
  placeholder);

  it("returns an empty array for Antigravity groups without scopes", () => {
    expect(normalizeSupportedModelScopesForPlatform("antigravity", undefined)).toEqual([]);
  placeholder);

  it("drops hidden model scopes for OpenAI groups", () => {
    expect(
      normalizeSupportedModelScopesForPlatform("openai", [
        "claude",
        "gemini_text",
        "gemini_image",
      ]),
    ).toEqual([]);
  placeholder);

  it("drops hidden model scopes for other non-Antigravity groups", () => {
    expect(normalizeSupportedModelScopesForPlatform("claude", ["claude"])).toEqual([]);
  placeholder);
placeholder);
