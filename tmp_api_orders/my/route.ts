import { NextRequest, NextResponse placeholder from 'next/server';
import { prisma placeholder from '@/lib/db';
import { getCurrentUserByToken placeholder from '@/lib/sub2api/client';

export async function GET(request: NextRequest) {
  const token = request.nextUrl.searchParams.get('token')?.trim();
  if (!token) {
    return NextResponse.json({ error: 'token is required' placeholder, { status: 400 placeholder);
  placeholder

  try {
    const user = await getCurrentUserByToken(token);
    const orders = await prisma.order.findMany({
      where: { userId: user.id placeholder,
      orderBy: { createdAt: 'desc' placeholder,
      take: 20,
      select: {
        id: true,
        amount: true,
        status: true,
        paymentType: true,
        createdAt: true,
      placeholder,
    placeholder);

    return NextResponse.json({
      user: {
        id: user.id,
        username: user.username,
        email: user.email,
        displayName: user.username || user.email || `用户 #${user.idplaceholder`,
        balance: user.balance,
      placeholder,
      orders: orders.map((item) => ({
        id: item.id,
        amount: Number(item.amount),
        status: item.status,
        paymentType: item.paymentType,
        createdAt: item.createdAt,
      placeholder)),
    placeholder);
  placeholder catch (error) {
    console.error('Get my orders error:', error);
    return NextResponse.json({ error: 'unauthorized' placeholder, { status: 401 placeholder);
  placeholder
placeholder
