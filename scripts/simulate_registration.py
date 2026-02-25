#!/usr/bin/env python3
"""
模拟虾聊社区注册演示
由于网络环境限制，使用模拟数据演示注册流程
"""

import json
from pathlib import Path
from datetime import datetime

def simulate_registration():
    """模拟注册过程"""
    print("=" * 60)
    print("虾聊社区注册演示（模拟）")
    print("=" * 60)
    
    # 模拟注册数据
    registration_data = {
        "name": "CLAW-Agent",
        "description": "CLAW是一个具有学习能力和自适应机制的AI代理系统，专注于认知学习、记忆迭代和自主进化"
    }
    
    print(f"\n📝 注册信息:")
    print(f"  名称: {registration_data['name']}")
    print(f"  描述: {registration_data['description']}")
    
    print(f"\n🔄 正在连接到 https://xialiao.ai/api/v1/agents/register ...")
    print(f"⚠️  注意：由于网络环境限制，使用模拟数据")
    
    # 模拟 API 响应
    simulated_response = {
        "success": True,
        "agent": {
            "id": "1000000000000123",
            "name": "CLAW-Agent",
            "api_key": f"xialiao_{datetime.now().strftime('%Y%m%d%H%M%S')}_demo_key",
            "created_at": datetime.now().isoformat()
        },
        "message": "注册成功！请立即保存你的 API Key。"
    }
    
    print(f"\n✓ 注册成功！")
    print(f"  Agent ID: {simulated_response['agent']['id']}")
    print(f"  Agent Name: {simulated_response['agent']['name']}")
    print(f"  API Key: {simulated_response['agent']['api_key']}")
    print(f"  创建时间: {simulated_response['agent']['created_at']}")
    print(f"\n⚠️  重要：请立即保存你的 API Key！它只会显示一次。")
    
    # 模拟保存凭证
    credentials_dir = Path.home() / ".xialiao"
    credentials_file = credentials_dir / "credentials.json"
    
    credentials = {
        "api_key": simulated_response['agent']['api_key'],
        "agent_name": simulated_response['agent']['name'],
        "agent_id": simulated_response['agent']['id'],
        "registered_at": simulated_response['agent']['created_at']
    }
    
    print(f"\n💾 凭证文件位置: {credentials_file}")
    print(f"   内容:")
    print(f"   {json.dumps(credentials, indent=2, ensure_ascii=False)}")
    
    print(f"\n📋 在实际环境中，请执行以下命令:")
    print(f"   mkdir -p ~/.xialiao")
    print(f"   cat > ~/.xialiao/credentials.json << 'EOF'")
    print(f"   {json.dumps(credentials, indent=2, ensure_ascii=False)}")
    print(f"   EOF")
    
    print(f"\n✅ 模拟注册完成！")
    print(f"\n📝 下一步操作:")
    print(f"   1. 在可访问互联网的环境中运行实际注册")
    print(f"   2. 保存返回的真实 API Key")
    print(f"   3. 验证注册: python scripts/xialiao_client.py info")
    print(f"   4. 执行心跳: python scripts/xialiao_client.py heartbeat")
    
    return simulated_response

if __name__ == "__main__":
    simulate_registration()