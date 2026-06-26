import { describe, it, expect placeholder from "vitest";
import {
  parseFingerprintSignalsToRows,
  serializeFingerprintRowsToJSON,
placeholder from "../codexFingerprintSignals";

describe("codex fingerprint signals 行编解码", () => {
  it("解析: 变体数组 → / 合并字符串", () => {
    const rows = parseFingerprintSignalsToRows(
      '[{"type":"header_exact","match":["session-id","session_id"],"required":trueplaceholder]',
    );
    expect(rows).toEqual([
      { type: "header_exact", match: "session-id / session_id", required: true placeholder,
    ]);
  placeholder);
  it("序列化: / 合并 → 变体数组, required 透传", () => {
    const json = serializeFingerprintRowsToJSON([
      { type: "header_prefix", match: "x-codex-", required: true placeholder,
      { type: "body_path", match: " a / b ", required: false placeholder,
    ]);
    expect(JSON.parse(json)).toEqual([
      { type: "header_prefix", match: ["x-codex-"], required: true placeholder,
      { type: "body_path", match: ["a", "b"], required: false placeholder,
    ]);
  placeholder);
  it("空/非法 → 空数组 / [] 串", () => {
    expect(parseFingerprintSignalsToRows("")).toEqual([]);
    expect(parseFingerprintSignalsToRows("nope")).toEqual([]);
    expect(serializeFingerprintRowsToJSON([])).toBe("[]");
  placeholder);
placeholder);
