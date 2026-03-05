/**
 * Admin Scheduled Tests API endpoints
 * Handles scheduled test plan management for account connectivity monitoring
 */

import { apiClient placeholder from '../client'
import type {
  ScheduledTestPlan,
  ScheduledTestResult,
  CreateScheduledTestPlanRequest,
  UpdateScheduledTestPlanRequest
placeholder from '@/types'

/**
 * List all scheduled test plans for an account
 * @param accountId - Account ID
 * @returns List of scheduled test plans
 */
export async function listByAccount(accountId: number): Promise<ScheduledTestPlan[]> {
  const { data placeholder = await apiClient.get<ScheduledTestPlan[]>(
    `/admin/accounts/${accountIdplaceholder/scheduled-test-plans`
  )
  return data
placeholder

/**
 * Create a new scheduled test plan
 * @param req - Plan creation request
 * @returns Created plan
 */
export async function create(req: CreateScheduledTestPlanRequest): Promise<ScheduledTestPlan> {
  const { data placeholder = await apiClient.post<ScheduledTestPlan>(
    '/admin/scheduled-test-plans',
    req
  )
  return data
placeholder

/**
 * Update an existing scheduled test plan
 * @param id - Plan ID
 * @param req - Fields to update
 * @returns Updated plan
 */
export async function update(id: number, req: UpdateScheduledTestPlanRequest): Promise<ScheduledTestPlan> {
  const { data placeholder = await apiClient.put<ScheduledTestPlan>(
    `/admin/scheduled-test-plans/${idplaceholder`,
    req
  )
  return data
placeholder

/**
 * Delete a scheduled test plan
 * @param id - Plan ID
 */
export async function deletePlan(id: number): Promise<void> {
  await apiClient.delete(`/admin/scheduled-test-plans/${idplaceholder`)
placeholder

/**
 * List test results for a plan
 * @param planId - Plan ID
 * @param limit - Optional max number of results to return
 * @returns List of test results
 */
export async function listResults(planId: number, limit?: number): Promise<ScheduledTestResult[]> {
  const { data placeholder = await apiClient.get<ScheduledTestResult[]>(
    `/admin/scheduled-test-plans/${planIdplaceholder/results`,
    {
      params: limit ? { limit placeholder : undefined
    placeholder
  )
  return data
placeholder

export const scheduledTestsAPI = {
  listByAccount,
  create,
  update,
  delete: deletePlan,
  listResults
placeholder

export default scheduledTestsAPI
