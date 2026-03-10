#!/bin/bash

# Gemini 粘性会话压力测试脚本
# 测试目标：验证不同会话分配不同账号，同一会话保持同一账号

BASE_URL="http://host.clicodeplus.com:8080"
API_KEY="sk-32ad0a3197e528c840ea84f0dc6b2056dd3fead03526b5c605a60709bd408f7e"
MODEL="gemini-2.5-flash"

# 创建临时目录存放结果
RESULT_DIR="/tmp/gemini_stress_test_$(date +%s)"
mkdir -p "$RESULT_DIR"

echo "=========================================="
echo "Gemini 粘性会话压力测试"
echo "结果目录: $RESULT_DIR"
echo "=========================================="

# 函数：发送请求并记录
send_request() {
    local session_id=$1
    local round=$2
    local system_prompt=$3
    local contents=$4
    local output_file="$RESULT_DIR/session_${session_idplaceholder_round_${roundplaceholder.json"

    local request_body=$(cat <<EOF
{
    "systemInstruction": {
        "parts": [{"text": "$system_prompt"placeholder]
    placeholder,
    "contents": $contents
placeholder
EOF
)

    curl -s -X POST "${BASE_URLplaceholder/v1beta/models/${MODELplaceholder:generateContent" \
        -H "Content-Type: application/json" \
        -H "x-goog-api-key: ${API_KEYplaceholder" \
        -d "$request_body" > "$output_file" 2>&1

    echo "[Session $session_id Round $round] 完成"
placeholder

# 会话1：数学计算器（累加序列）
run_session_1() {
    local sys_prompt="你是一个数学计算器，只返回计算结果数字，不要任何解释"

    # Round 1: 1+1=?
    send_request 1 1 "$sys_prompt" '[{"role":"user","parts":[{"text":"1+1=?"placeholder]placeholder]'

    # Round 2: 继续 2+2=?（累加历史）
    send_request 1 2 "$sys_prompt" '[{"role":"user","parts":[{"text":"1+1=?"placeholder]placeholder,{"role":"model","parts":[{"text":"2"placeholder]placeholder,{"role":"user","parts":[{"text":"2+2=?"placeholder]placeholder]'

    # Round 3: 继续 3+3=?
    send_request 1 3 "$sys_prompt" '[{"role":"user","parts":[{"text":"1+1=?"placeholder]placeholder,{"role":"model","parts":[{"text":"2"placeholder]placeholder,{"role":"user","parts":[{"text":"2+2=?"placeholder]placeholder,{"role":"model","parts":[{"text":"4"placeholder]placeholder,{"role":"user","parts":[{"text":"3+3=?"placeholder]placeholder]'

    # Round 4: 批量计算 10+10, 20+20, 30+30
    send_request 1 4 "$sys_prompt" '[{"role":"user","parts":[{"text":"1+1=?"placeholder]placeholder,{"role":"model","parts":[{"text":"2"placeholder]placeholder,{"role":"user","parts":[{"text":"2+2=?"placeholder]placeholder,{"role":"model","parts":[{"text":"4"placeholder]placeholder,{"role":"user","parts":[{"text":"3+3=?"placeholder]placeholder,{"role":"model","parts":[{"text":"6"placeholder]placeholder,{"role":"user","parts":[{"text":"计算: 10+10=? 20+20=? 30+30=?"placeholder]placeholder]'
placeholder

# 会话2：英文翻译器（不同系统提示词 = 不同会话）
run_session_2() {
    local sys_prompt="你是一个英文翻译器，将中文翻译成英文，只返回翻译结果"

    send_request 2 1 "$sys_prompt" '[{"role":"user","parts":[{"text":"你好"placeholder]placeholder]'
    send_request 2 2 "$sys_prompt" '[{"role":"user","parts":[{"text":"你好"placeholder]placeholder,{"role":"model","parts":[{"text":"Hello"placeholder]placeholder,{"role":"user","parts":[{"text":"世界"placeholder]placeholder]'
    send_request 2 3 "$sys_prompt" '[{"role":"user","parts":[{"text":"你好"placeholder]placeholder,{"role":"model","parts":[{"text":"Hello"placeholder]placeholder,{"role":"user","parts":[{"text":"世界"placeholder]placeholder,{"role":"model","parts":[{"text":"World"placeholder]placeholder,{"role":"user","parts":[{"text":"早上好"placeholder]placeholder]'
placeholder

# 会话3：日文翻译器
run_session_3() {
    local sys_prompt="你是一个日文翻译器，将中文翻译成日文，只返回翻译结果"

    send_request 3 1 "$sys_prompt" '[{"role":"user","parts":[{"text":"你好"placeholder]placeholder]'
    send_request 3 2 "$sys_prompt" '[{"role":"user","parts":[{"text":"你好"placeholder]placeholder,{"role":"model","parts":[{"text":"こんにちは"placeholder]placeholder,{"role":"user","parts":[{"text":"谢谢"placeholder]placeholder]'
    send_request 3 3 "$sys_prompt" '[{"role":"user","parts":[{"text":"你好"placeholder]placeholder,{"role":"model","parts":[{"text":"こんにちは"placeholder]placeholder,{"role":"user","parts":[{"text":"谢谢"placeholder]placeholder,{"role":"model","parts":[{"text":"ありがとう"placeholder]placeholder,{"role":"user","parts":[{"text":"再见"placeholder]placeholder]'
placeholder

# 会话4：乘法计算器（另一个数学会话，但系统提示词不同）
run_session_4() {
    local sys_prompt="你是一个乘法专用计算器，只计算乘法，返回数字结果"

    send_request 4 1 "$sys_prompt" '[{"role":"user","parts":[{"text":"2*3=?"placeholder]placeholder]'
    send_request 4 2 "$sys_prompt" '[{"role":"user","parts":[{"text":"2*3=?"placeholder]placeholder,{"role":"model","parts":[{"text":"6"placeholder]placeholder,{"role":"user","parts":[{"text":"4*5=?"placeholder]placeholder]'
    send_request 4 3 "$sys_prompt" '[{"role":"user","parts":[{"text":"2*3=?"placeholder]placeholder,{"role":"model","parts":[{"text":"6"placeholder]placeholder,{"role":"user","parts":[{"text":"4*5=?"placeholder]placeholder,{"role":"model","parts":[{"text":"20"placeholder]placeholder,{"role":"user","parts":[{"text":"计算: 10*10=? 20*20=?"placeholder]placeholder]'
placeholder

# 会话5：诗人（完全不同的角色）
run_session_5() {
    local sys_prompt="你是一位诗人，用简短的诗句回应每个话题，每次只写一句诗"

    send_request 5 1 "$sys_prompt" '[{"role":"user","parts":[{"text":"春天"placeholder]placeholder]'
    send_request 5 2 "$sys_prompt" '[{"role":"user","parts":[{"text":"春天"placeholder]placeholder,{"role":"model","parts":[{"text":"春风拂面花满枝"placeholder]placeholder,{"role":"user","parts":[{"text":"夏天"placeholder]placeholder]'
    send_request 5 3 "$sys_prompt" '[{"role":"user","parts":[{"text":"春天"placeholder]placeholder,{"role":"model","parts":[{"text":"春风拂面花满枝"placeholder]placeholder,{"role":"user","parts":[{"text":"夏天"placeholder]placeholder,{"role":"model","parts":[{"text":"蝉鸣蛙声伴荷香"placeholder]placeholder,{"role":"user","parts":[{"text":"秋天"placeholder]placeholder]'
placeholder

echo ""
echo "开始并发测试 5 个独立会话..."
echo ""

# 并发运行所有会话
run_session_1 &
run_session_2 &
run_session_3 &
run_session_4 &
run_session_5 &

# 等待所有后台任务完成
wait

echo ""
echo "=========================================="
echo "所有请求完成，结果保存在: $RESULT_DIR"
echo "=========================================="

# 显示结果摘要
echo ""
echo "响应摘要:"
for f in "$RESULT_DIR"/*.json; do
    filename=$(basename "$f")
    response=$(cat "$f" | head -c 200)
    echo "[$filename]: ${responseplaceholder..."
done

echo ""
echo "请检查服务器日志确认账号分配情况"
