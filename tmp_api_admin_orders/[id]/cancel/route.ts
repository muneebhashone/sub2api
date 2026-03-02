import { NextRequest, NextResponse placeholder from 'next/server';
import { verifyAdminToken, unauthorizedResponse placeholder from '@/lib/admin-auth';
import { adminCancelOrder, OrderError placeholder from '@/lib/order/service';

export async function POST(
  request: NextRequest,
  { params placeholder: { params: Promise<{ id: string placeholder> placeholder,
) {
  if (!verifyAdminToken(request)) return unauthorizedResponse();

  try {
    const { id placeholder = await params;
    await adminCancelOrder(id);
    return NextResponse.json({ success: true placeholder);
  placeholder catch (error) {
    if (error instanceof OrderError) {
      return NextResponse.json(
        { error: error.message, code: error.code placeholder,
        { status: error.statusCode placeholder,
      );
    placeholder
    console.error('Admin cancel order error:', error);
    return NextResponse.json({ error: '取消订单失败' placeholder, { status: 500 placeholder);
  placeholder
placeholder
