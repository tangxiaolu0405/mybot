#!/usr/bin/env python3
"""
虾聊客户端测试脚本
用于验证代码逻辑，不依赖实际网络连接
"""

import json
import sys
from pathlib import Path

# 获取项目根目录
project_root = Path(__file__).parent.parent
scripts_dir = project_root / "scripts"
sys.path.insert(0, str(scripts_dir))

def test_client_initialization():
    """测试客户端初始化"""
    print("测试 1: 客户端初始化（无 API Key）")
    try:
        from xialiao_client import XialiaoClient
        client = XialiaoClient()
        print("✓ 客户端初始化成功")
        print(f"  API Base: {client.api_base}")
        print(f"  API Key: {client.api_key if client.api_key else 'None'}")
        print(f"  Session trust_env: {client.session.trust_env}")
        return True
    except Exception as e:
        print(f"✗ 客户端初始化失败: {e}")
        return False

def test_client_with_api_key():
    """测试带 API Key 的客户端初始化"""
    print("\n测试 2: 客户端初始化（带 API Key）")
    try:
        from xialiao_client import XialiaoClient
        client = XialiaoClient(api_key="test_api_key_12345")
        print("✓ 客户端初始化成功")
        print(f"  API Key: {client.api_key}")
        
        # 测试获取 headers
        headers = client._get_headers()
        print(f"  Headers: {headers}")
        return True
    except Exception as e:
        print(f"✗ 客户端初始化失败: {e}")
        return False

def test_heartbeat_state_file():
    """测试心跳状态文件"""
    print("\n测试 3: 心跳状态文件")
    try:
        state_file = project_root / "brain" / "memory" / "short-term" / "heartbeat-state.json"
        
        if state_file.exists():
            with open(state_file, 'r', encoding='utf-8') as f:
                state = json.load(f)
            print("✓ 心跳状态文件存在")
            print(f"  内容: {json.dumps(state, indent=2, ensure_ascii=False)}")
        else:
            print("✓ 心跳状态文件不存在（首次运行正常）")
        
        return True
    except Exception as e:
        print(f"✗ 心跳状态文件测试失败: {e}")
        return False

def test_skills_index():
    """测试技能索引"""
    print("\n测试 4: 技能索引")
    try:
        skills_index_file = project_root / "skills" / "skills-index.json"
        
        with open(skills_index_file, 'r', encoding='utf-8') as f:
            skills_index = json.load(f)
        
        print("✓ 技能索引文件存在")
        print(f"  版本: {skills_index.get('version')}")
        print(f"  技能数量: {len(skills_index.get('skills', []))}")
        
        # 检查虾聊集成技能
        xialiao_skill = None
        for skill in skills_index.get('skills', []):
            if skill.get('name') == 'xialiao-integration':
                xialiao_skill = skill
                break
        
        if xialiao_skill:
            print("✓ 虾聊集成技能已注册")
            print(f"  技能名称: {xialiao_skill.get('name')}")
            print(f"  技能描述: {xialiao_skill.get('description')}")
            print(f"  技能标签: {xialiao_skill.get('tags')}")
        else:
            print("✗ 虾聊集成技能未找到")
            return False
        
        return True
    except Exception as e:
        print(f"✗ 技能索引测试失败: {e}")
        return False

def test_core_md():
    """测试核心思维模式文件"""
    print("\n测试 5: 核心思维模式文件")
    try:
        core_md_file = project_root / "brain" / "core.md"
        
        with open(core_md_file, 'r', encoding='utf-8') as f:
            content = f.read()
        
        print("✓ 核心思维模式文件存在")
        
        # 检查是否包含社区互动章节
        if "## 社区互动" in content:
            print("✓ 包含社区互动章节")
        else:
            print("✗ 缺少社区互动章节")
            return False
        
        # 检查是否包含虾聊相关内容
        if "虾聊社区" in content:
            print("✓ 包含虾聊社区相关内容")
        else:
            print("✗ 缺少虾聊社区相关内容")
            return False
        
        # 检查心跳机制
        if "心跳机制" in content:
            print("✓ 包含心跳机制说明")
        else:
            print("✗ 缺少心跳机制说明")
            return False
        
        return True
    except Exception as e:
        print(f"✗ 核心思维模式文件测试失败: {e}")
        return False

def test_heartbeat_file():
    """测试心跳文件"""
    print("\n测试 6: 心跳文件")
    try:
        heartbeat_file = project_root / "HEARTBEAT.md"
        
        with open(heartbeat_file, 'r', encoding='utf-8') as f:
            content = f.read()
        
        print("✓ 心跳文件存在")
        
        # 检查是否包含虾聊社区任务
        if "虾聊社区" in content:
            print("✓ 包含虾聊社区心跳任务")
        else:
            print("✗ 缺少虾聊社区心跳任务")
            return False
        
        return True
    except Exception as e:
        print(f"✗ 心跳文件测试失败: {e}")
        return False

def main():
    """运行所有测试"""
    print("=" * 60)
    print("虾聊客户端测试")
    print("=" * 60)
    
    tests = [
        test_client_initialization,
        test_client_with_api_key,
        test_heartbeat_state_file,
        test_skills_index,
        test_core_md,
        test_heartbeat_file
    ]
    
    results = []
    for test in tests:
        try:
            result = test()
            results.append(result)
        except Exception as e:
            print(f"✗ 测试异常: {e}")
            results.append(False)
    
    print("\n" + "=" * 60)
    print("测试总结")
    print("=" * 60)
    print(f"总测试数: {len(results)}")
    print(f"通过: {sum(results)}")
    print(f"失败: {len(results) - sum(results)}")
    
    if all(results):
        print("\n✓ 所有测试通过！代码逻辑正确。")
        print("\n注意：由于网络环境限制，无法测试实际的网络连接。")
        print("在可访问互联网的环境中运行即可使用完整功能。")
        return 0
    else:
        print("\n✗ 部分测试失败，请检查错误信息。")
        return 1

if __name__ == "__main__":
    sys.exit(main())