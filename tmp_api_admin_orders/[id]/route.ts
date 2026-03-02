import { NextRequest, NextResponse placeholder from 'next/server';
import { prisma placeholder from '@/lib/db';
import { verifyAdminToken, unauthorizedResponse placeholder from '@/lib/admin-auth';

export async function GET(
  request: NextRequest,
  { params placeholder: { params: Promise<{ id: string placeholder> placeholder,
) {
  if (!verifyAdminToken(request)) return unauthorizedResponse();

  const { id placeholder = await params;

  const order = await prisma.order.findUnique({
    where: { id placeholder,
    include: {
      auditLogs: {
        orderBy: { createdAt: 'desc' placeholder,
      placeholder,
    placeholder,
  placeholder);

  if (!order) {
    return NextResponse.json({ error: '订单不存在' placeholder, { status: 404 placeholder);
  placeholder

  return NextResponse.json({
    ...order,
    amount: Number(order.amount),
    refundAmount: order.refundAmount ? Number(order.refundAmount) : null,
  placeholder);
placeholder
