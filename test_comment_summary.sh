#!/bin/bash

# AI评论总结功能测试脚本
# 使用方法: ./test_comment_summary.sh [测试用例名称]

echo "🚀 AI评论总结功能测试"
echo "========================================"

# 检查环境变量
check_env() {
    local missing_vars=()
    
    if [ -z "$LLM_API_KEY" ]; then
        missing_vars+=("LLM_API_KEY")
    fi
    
    if [ -z "$LLM_URL" ]; then
        missing_vars+=("LLM_URL")
    fi
    
    if [ -z "$LLM_MODEL" ]; then
        missing_vars+=("LLM_MODEL")
    fi
    
    if [ -z "$DB_HOST" ]; then
        missing_vars+=("DB_HOST")
    fi
    
    if [ -z "$DB_USER" ]; then
        missing_vars+=("DB_USER")
    fi
    
    if [ -z "$DB_DATABASE" ]; then
        missing_vars+=("DB_DATABASE")
    fi
    
    if [ ${#missing_vars[@]} -ne 0 ]; then
        echo "❌ 缺少必要的环境变量:"
        for var in "${missing_vars[@]}"; do
            echo "   - $var"
        done
        echo ""
        echo "请设置环境变量后再运行测试，例如:"
        echo "export LLM_API_KEY='your_api_key'"
        echo "export LLM_URL='your_llm_endpoint'"
        echo "export LLM_MODEL='your_model_name'"
        echo "export DB_HOST='your_db_host:port'"
        echo "export DB_USER='your_db_user'"
        echo "export DB_PASS='your_db_password'"
        echo "export DB_DATABASE='your_db_name'"
        exit 1
    fi
    
    echo "✅ 环境变量检查通过"
}

# 显示配置信息
show_config() {
    echo ""
    echo "📋 当前测试配置:"
    echo "   - 数据库: $DB_HOST/$DB_DATABASE"
    echo "   - LLM服务: $LLM_URL"
    echo "   - LLM模型: $LLM_MODEL"
    echo ""
    echo "⚠️  请确保在 server/service/comment_summary_test.go 中设置了正确的测试参数:"
    echo "   - TEST_BUSINESS_ID: 实际存在评论的文章ID"
    echo "   - TEST_TENANT_ID: 租户类型 (0=文章评论)"
    echo ""
}

# 运行测试
run_test() {
    local test_name=$1
    
    if [ -z "$test_name" ]; then
        echo "🧪 运行所有测试用例..."
        go test ./server/service -v -run TestCommentSummary
    else
        echo "🧪 运行测试用例: $test_name"
        go test ./server/service -v -run "$test_name"
    fi
    
    local exit_code=$?
    
    echo ""
    if [ $exit_code -eq 0 ]; then
        echo "✅ 测试完成！"
    else
        echo "❌ 测试失败，退出码: $exit_code"
        echo ""
        echo "💡 常见问题排查:"
        echo "   1. 检查数据库连接是否正常"
        echo "   2. 检查LLM服务是否可用"
        echo "   3. 检查TEST_BUSINESS_ID对应的文章是否有评论"
        echo "   4. 检查网络连接是否正常"
    fi
    
    return $exit_code
}

# 主函数
main() {
    check_env
    show_config
    
    echo "可用的测试用例:"
    echo "   - TestCommentSummaryWithRealData     (完整功能测试)"
    echo "   - TestCommentSummaryUpdate           (更新机制测试)"
    echo "   - TestCommentSummaryPerformance      (性能测试)"
    echo "   - TestCommentSummaryErrorHandling    (错误处理测试)"
    echo ""
    
    run_test "$1"
}

# 执行主函数
main "$@"