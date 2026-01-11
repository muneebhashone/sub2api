/**
 * Admin Promo Codes API endpoints
 */

import { apiClient placeholder from '../client'
import type {
  PromoCode,
  PromoCodeUsage,
  CreatePromoCodeRequest,
  UpdatePromoCodeRequest,
  BasePaginationResponse
placeholder from '@/types'

export async function list(
  page: number = 1,
  pageSize: number = 20,
  filters?: {
    status?: string
    search?: string
  placeholder
): Promise<BasePaginationResponse<PromoCode>> {
  const { data placeholder = await apiClient.get<BasePaginationResponse<PromoCode>>('/admin/promo-codes', {
    params: { page, page_size: pageSize, ...filters placeholder
  placeholder)
  return data
placeholder

export async function getById(id: number): Promise<PromoCode> {
  const { data placeholder = await apiClient.get<PromoCode>(`/admin/promo-codes/${idplaceholder`)
  return data
placeholder

export async function create(request: CreatePromoCodeRequest): Promise<PromoCode> {
  const { data placeholder = await apiClient.post<PromoCode>('/admin/promo-codes', request)
  return data
placeholder

export async function update(id: number, request: UpdatePromoCodeRequest): Promise<PromoCode> {
  const { data placeholder = await apiClient.put<PromoCode>(`/admin/promo-codes/${idplaceholder`, request)
  return data
placeholder

export async function deleteCode(id: number): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.delete<{ message: string placeholder>(`/admin/promo-codes/${idplaceholder`)
  return data
placeholder

export async function getUsages(
  id: number,
  page: number = 1,
  pageSize: number = 20
): Promise<BasePaginationResponse<PromoCodeUsage>> {
  const { data placeholder = await apiClient.get<BasePaginationResponse<PromoCodeUsage>>(
    `/admin/promo-codes/${idplaceholder/usages`,
    { params: { page, page_size: pageSize placeholder placeholder
  )
  return data
placeholder

const promoAPI = {
  list,
  getById,
  create,
  update,
  delete: deleteCode,
  getUsages
placeholder

export default promoAPI
