import http from 'k6/http';
import { check, sleep placeholder from 'k6';
import { Rate, Trend placeholder from 'k6/metrics';

const baseURL = (__ENV.BASE_URL || 'http://127.0.0.1:5231').replace(/\/$/, '');
const wsAPIKey = (__ENV.WS_API_KEY || '').trim();
const wsHotspotAPIKey = (__ENV.WS_HOTSPOT_API_KEY || wsAPIKey).trim();
const model = __ENV.MODEL || 'gpt-5.3-codex';
const duration = __ENV.DURATION || '5m';
const timeout = __ENV.TIMEOUT || '180s';

const shortRPS = Number(__ENV.SHORT_RPS || 12);
const longRPS = Number(__ENV.LONG_RPS || 4);
const errorRPS = Number(__ENV.ERROR_RPS || 2);
const hotspotRPS = Number(__ENV.HOTSPOT_RPS || 10);
const preAllocatedVUs = Number(__ENV.PRE_ALLOCATED_VUS || 50);
const maxVUs = Number(__ENV.MAX_VUS || 400);

const reqDurationMs = new Trend('openai_ws_v2_perf_req_duration_ms', true);
const ttftMs = new Trend('openai_ws_v2_perf_ttft_ms', true);
const non2xxRate = new Rate('openai_ws_v2_perf_non2xx_rate');
const doneRate = new Rate('openai_ws_v2_perf_done_rate');
const expectedErrorRate = new Rate('openai_ws_v2_perf_expected_error_rate');

export const options = {
  scenarios: {
    short_request: {
      executor: 'constant-arrival-rate',
      exec: 'runShortRequest',
      rate: shortRPS,
      timeUnit: '1s',
      duration,
      preAllocatedVUs,
      maxVUs,
      tags: { scenario: 'short_request' placeholder,
    placeholder,
    long_request: {
      executor: 'constant-arrival-rate',
      exec: 'runLongRequest',
      rate: longRPS,
      timeUnit: '1s',
      duration,
      preAllocatedVUs: Math.max(20, Math.ceil(longRPS * 6)),
      maxVUs: Math.max(100, Math.ceil(longRPS * 20)),
      tags: { scenario: 'long_request' placeholder,
    placeholder,
    error_injection: {
      executor: 'constant-arrival-rate',
      exec: 'runErrorInjection',
      rate: errorRPS,
      timeUnit: '1s',
      duration,
      preAllocatedVUs: Math.max(8, Math.ceil(errorRPS * 4)),
      maxVUs: Math.max(40, Math.ceil(errorRPS * 12)),
      tags: { scenario: 'error_injection' placeholder,
    placeholder,
    hotspot_account: {
      executor: 'constant-arrival-rate',
      exec: 'runHotspotAccount',
      rate: hotspotRPS,
      timeUnit: '1s',
      duration,
      preAllocatedVUs: Math.max(16, Math.ceil(hotspotRPS * 3)),
      maxVUs: Math.max(80, Math.ceil(hotspotRPS * 10)),
      tags: { scenario: 'hotspot_account' placeholder,
    placeholder,
  placeholder,
  thresholds: {
    openai_ws_v2_perf_non2xx_rate: ['rate<0.05'],
    openai_ws_v2_perf_req_duration_ms: ['p(95)<5000', 'p(99)<9000'],
    openai_ws_v2_perf_ttft_ms: ['p(99)<2000'],
    openai_ws_v2_perf_done_rate: ['rate>0.95'],
  placeholder,
placeholder;

function buildHeaders(apiKey, opts = {placeholder) {
  const headers = {
    'Content-Type': 'application/json',
    'User-Agent': 'codex_cli_rs/0.104.0',
    'OpenAI-Beta': 'responses_websockets=2026-02-06,responses=experimental',
  placeholder;
  if (apiKey) {
    headers.Authorization = `Bearer ${apiKeyplaceholder`;
  placeholder
  if (opts.sessionID) {
    headers.session_id = opts.sessionID;
  placeholder
  if (opts.conversationID) {
    headers.conversation_id = opts.conversationID;
  placeholder
  return headers;
placeholder

function shortBody() {
  return JSON.stringify({
    model,
    stream: false,
    input: [
      {
        role: 'user',
        content: [{ type: 'input_text', text: '请回复一个词：pong' placeholder],
      placeholder,
    ],
    max_output_tokens: 64,
  placeholder);
placeholder

function longBody() {
  const tools = [];
  for (let i = 0; i < 28; i += 1) {
    tools.push({
      type: 'function',
      name: `perf_tool_${iplaceholder`,
      description: 'load test tool schema',
      parameters: {
        type: 'object',
        properties: {
          query: { type: 'string' placeholder,
          limit: { type: 'number' placeholder,
          with_cache: { type: 'boolean' placeholder,
        placeholder,
        required: ['query'],
      placeholder,
    placeholder);
  placeholder

  const input = [];
  for (let i = 0; i < 20; i += 1) {
    input.push({
      role: 'user',
      content: [{ type: 'input_text', text: `长请求压测消息 ${iplaceholder: 请输出简要摘要。` placeholder],
    placeholder);
  placeholder

  return JSON.stringify({
    model,
    stream: false,
    input,
    tools,
    parallel_tool_calls: true,
    max_output_tokens: 256,
    reasoning: { effort: 'medium' placeholder,
    instructions: '你是压测助手，简洁回复。',
  placeholder);
placeholder

function errorInjectionBody() {
  return JSON.stringify({
    model,
    stream: false,
    previous_response_id: `resp_not_found_${__VUplaceholder_${__ITERplaceholder`,
    input: [
      {
        role: 'user',
        content: [{ type: 'input_text', text: '触发错误注入路径。' placeholder],
      placeholder,
    ],
  placeholder);
placeholder

function postResponses(apiKey, body, tags, opts = {placeholder) {
  const res = http.post(`${baseURLplaceholder/v1/responses`, body, {
    headers: buildHeaders(apiKey, opts),
    timeout,
    tags,
  placeholder);
  reqDurationMs.add(res.timings.duration, tags);
  ttftMs.add(res.timings.waiting, tags);
  non2xxRate.add(res.status < 200 || res.status >= 300, tags);
  return res;
placeholder

function hasDone(res) {
  return !!res && !!res.body && res.body.indexOf('[DONE]') >= 0;
placeholder

export function runShortRequest() {
  const tags = { scenario: 'short_request' placeholder;
  const res = postResponses(wsAPIKey, shortBody(), tags);
  check(res, { 'short status is 2xx': (r) => r.status >= 200 && r.status < 300 placeholder);
  doneRate.add(hasDone(res) || (res.status >= 200 && res.status < 300), tags);
placeholder

export function runLongRequest() {
  const tags = { scenario: 'long_request' placeholder;
  const res = postResponses(wsAPIKey, longBody(), tags);
  check(res, { 'long status is 2xx': (r) => r.status >= 200 && r.status < 300 placeholder);
  doneRate.add(hasDone(res) || (res.status >= 200 && res.status < 300), tags);
placeholder

export function runErrorInjection() {
  const tags = { scenario: 'error_injection' placeholder;
  const res = postResponses(wsAPIKey, errorInjectionBody(), tags);
  // 错误注入场景允许 4xx/5xx，重点观测 fallback 和错误路径抖动。
  expectedErrorRate.add(res.status >= 400, tags);
  doneRate.add(hasDone(res), tags);
placeholder

export function runHotspotAccount() {
  const tags = { scenario: 'hotspot_account' placeholder;
  const opts = {
    sessionID: 'perf-hotspot-session-fixed',
    conversationID: 'perf-hotspot-conversation-fixed',
  placeholder;
  const res = postResponses(wsHotspotAPIKey, shortBody(), tags, opts);
  check(res, { 'hotspot status is 2xx': (r) => r.status >= 200 && r.status < 300 placeholder);
  doneRate.add(hasDone(res) || (res.status >= 200 && res.status < 300), tags);
  sleep(0.01);
placeholder

export function handleSummary(data) {
  return {
    stdout: `\nOpenAI WSv2 性能套件压测完成\n${JSON.stringify(data.metrics, null, 2)placeholder\n`,
    'docs/perf/openai-ws-v2-perf-suite-summary.json': JSON.stringify(data, null, 2),
  placeholder;
placeholder
