import { NextRequest, NextResponse placeholder from 'next/server';
import { verifyAdminToken, unauthorizedResponse placeholder from '@/lib/admin-auth';
import { retryRecharge, OrderError placeholder from '@/lib/order/service';

export async function POST(
  request: NextRequest,
  { params placeholder: { params: Promise<{ id: string placeholder> placeholder,
) {
  if (!verifyAdminToken(request)) return unauthorizedResponse();

  try {
    const { id placeholder = await params;
    await retryRecharge(id);
    return NextResponse.json({ success: true placeholder);
  placeholder catch (error) {
    if (error instanceof OrderError) {
      return NextResponse.json(
        { error: error.message, code: error.code placeholder,
        { status: error.statusCode placeholder,
      );
    placeholder
    console.error('Retry recharge error:', error);
    return NextResponse.json({ error: '重试充值失败' placeholder, { status: 500 placeholder);
  placeholder
placeholder
