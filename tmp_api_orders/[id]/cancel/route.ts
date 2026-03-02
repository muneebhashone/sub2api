import { NextRequest, NextResponse placeholder from 'next/server';
import { z placeholder from 'zod';
import { cancelOrder, OrderError placeholder from '@/lib/order/service';

const cancelSchema = z.object({
  user_id: z.number().int().positive(),
placeholder);

export async function POST(
  request: NextRequest,
  { params placeholder: { params: Promise<{ id: string placeholder> placeholder,
) {
  try {
    const { id placeholder = await params;
    const body = await request.json();
    const parsed = cancelSchema.safeParse(body);

    if (!parsed.success) {
      return NextResponse.json(
        { error: '参数错误', details: parsed.error.flatten().fieldErrors placeholder,
        { status: 400 placeholder,
      );
    placeholder

    await cancelOrder(id, parsed.data.user_id);
    return NextResponse.json({ success: true placeholder);
  placeholder catch (error) {
    if (error instanceof OrderError) {
      return NextResponse.json(
        { error: error.message, code: error.code placeholder,
        { status: error.statusCode placeholder,
      );
    placeholder
    console.error('Cancel order error:', error);
    return NextResponse.json({ error: '取消订单失败' placeholder, { status: 500 placeholder);
  placeholder
placeholder
