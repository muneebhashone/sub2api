import { NextRequest, NextResponse placeholder from 'next/server';
import { z placeholder from 'zod';
import { createOrder, OrderError placeholder from '@/lib/order/service';
import { getEnv placeholder from '@/lib/config';

const createOrderSchema = z.object({
  user_id: z.number().int().positive(),
  amount: z.number().positive(),
  payment_type: z.enum(['alipay', 'wxpay']),
placeholder);

export async function POST(request: NextRequest) {
  try {
    const env = getEnv();
    const body = await request.json();
    const parsed = createOrderSchema.safeParse(body);

    if (!parsed.success) {
      return NextResponse.json(
        { error: '参数错误', details: parsed.error.flatten().fieldErrors placeholder,
        { status: 400 placeholder,
      );
    placeholder

    const { user_id, amount, payment_type placeholder = parsed.data;

    // Validate amount range
    if (amount < env.MIN_RECHARGE_AMOUNT || amount > env.MAX_RECHARGE_AMOUNT) {
      return NextResponse.json(
        { error: `充值金额需在 ${env.MIN_RECHARGE_AMOUNTplaceholder - ${env.MAX_RECHARGE_AMOUNTplaceholder 之间` placeholder,
        { status: 400 placeholder,
      );
    placeholder

    // Validate payment type is enabled
    if (!env.ENABLED_PAYMENT_TYPES.includes(payment_type)) {
      return NextResponse.json(
        { error: `不支持的支付方式: ${payment_typeplaceholder` placeholder,
        { status: 400 placeholder,
      );
    placeholder

    const clientIp = request.headers.get('x-forwarded-for')?.split(',')[0]?.trim()
      || request.headers.get('x-real-ip')
      || '127.0.0.1';

    const result = await createOrder({
      userId: user_id,
      amount,
      paymentType: payment_type,
      clientIp,
    placeholder);

    return NextResponse.json(result);
  placeholder catch (error) {
    if (error instanceof OrderError) {
      return NextResponse.json(
        { error: error.message, code: error.code placeholder,
        { status: error.statusCode placeholder,
      );
    placeholder
    console.error('Create order error:', error);
    return NextResponse.json(
      { error: '创建订单失败，请稍后重试' placeholder,
      { status: 500 placeholder,
    );
  placeholder
placeholder
