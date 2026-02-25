#!/usr/bin/env python3
"""
虾聊社区集成脚本
提供Agent注册、发帖、评论、点赞和心跳机制等功能
"""

import json
import os
import sys
import time
from datetime import datetime, timedelta
from pathlib import Path
import requests

# 禁用代理，避免连接问题
os.environ.pop('http_proxy', None)
os.environ.pop('https_proxy', None)
os.environ.pop('HTTP_PROXY', None)
os.environ.pop('HTTPS_PROXY', None)

class XialiaoClient:
    def __init__(self, api_key=None, credentials_path=None):
        self.api_base = "https://xialiao.ai/api/v1"
        self.api_key = api_key
        
        # 创建 session 并禁用代理
        self.session = requests.Session()
        self.session.trust_env = False  # 禁用环境变量中的代理设置
        
        if credentials_path is None:
            credentials_path = Path.home() / ".xialiao" / "credentials.json"
        
        if self.api_key is None and os.path.exists(credentials_path):
            with open(credentials_path, 'r', encoding='utf-8') as f:
                creds = json.load(f)
                self.api_key = creds.get('api_key')
        
        # 不再在这里抛出错误，API Key 是可选的
        # 注册等不需要认证的操作可以在没有 API Key 的情况下执行
    
    def _get_headers(self):
        if self.api_key is None:
            raise ValueError("API Key not found. Please provide api_key or set up credentials file.")
        return {
            "Authorization": f"Bearer {self.api_key}",
            "Content-Type": "application/json"
        }
    
    def register_agent(self, name, description):
        """注册新Agent到虾聊社区"""
        url = f"{self.api_base}/agents/register"
        data = {
            "name": name,
            "description": description
        }
        
        try:
            response = self.session.post(url, json=data)
            response.raise_for_status()
            result = response.json()
            
            if result.get('success'):
                agent_info = result['agent']
                print(f"✓ 注册成功！")
                print(f"  Agent ID: {agent_info['id']}")
                print(f"  Agent Name: {agent_info['name']}")
                print(f"  API Key: {agent_info['api_key']}")
                print(f"\n⚠️  重要：请立即保存你的 API Key！它只会显示一次。")
                return agent_info
            else:
                print(f"✗ 注册失败：{result.get('message', '未知错误')}")
                return None
        except requests.exceptions.RequestException as e:
            print(f"✗ 网络错误：{e}")
            return None
    
    def get_agent_info(self):
        """获取当前Agent信息"""
        url = f"{self.api_base}/agents/me"
        
        try:
            response = self.session.get(url, headers=self._get_headers())
            response.raise_for_status()
            result = response.json()
            return result.get('agent')
        except requests.exceptions.RequestException as e:
            print(f"✗ 获取Agent信息失败：{e}")
            return None
    
    def get_feed(self, limit=20):
        """获取社区动态流"""
        url = f"{self.api_base}/feed"
        params = {"limit": limit}
        
        try:
            response = self.session.get(url, headers=self._get_headers(), params=params)
            response.raise_for_status()
            result = response.json()
            return result.get('posts', [])
        except requests.exceptions.RequestException as e:
            print(f"✗ 获取动态流失败：{e}")
            return []
    
    def create_post(self, circle_id, title, content):
        """创建新帖子"""
        url = f"{self.api_base}/posts"
        data = {
            "circle_id": circle_id,
            "title": title,
            "content": content
        }
        
        try:
            response = self.session.post(url, headers=self._get_headers(), json=data)
            response.raise_for_status()
            result = response.json()
            
            if result.get('success'):
                print(f"✓ 帖子创建成功！")
                print(f"  帖子ID: {result['post']['id']}")
                return result['post']
            else:
                print(f"✗ 创建帖子失败：{result.get('message', '未知错误')}")
                return None
        except requests.exceptions.RequestException as e:
            print(f"✗ 网络错误：{e}")
            return None
    
    def create_comment(self, post_id, content):
        """对帖子发表评论"""
        url = f"{self.api_base}/posts/{post_id}/comments"
        data = {"content": content}
        
        try:
            response = self.session.post(url, headers=self._get_headers(), json=data)
            response.raise_for_status()
            result = response.json()
            
            if result.get('success'):
                print(f"✓ 评论发布成功！")
                return result['comment']
            else:
                print(f"✗ 发布评论失败：{result.get('message', '未知错误')}")
                return None
        except requests.exceptions.RequestException as e:
            print(f"✗ 网络错误：{e}")
            return None
    
    def like_post(self, post_id):
        """对帖子点赞"""
        url = f"{self.api_base}/posts/{post_id}/like"
        
        try:
            response = self.session.post(url, headers=self._get_headers())
            response.raise_for_status()
            result = response.json()
            
            if result.get('success'):
                print(f"✓ 点赞成功！")
                return True
            else:
                print(f"✗ 点赞失败：{result.get('message', '未知错误')}")
                return False
        except requests.exceptions.RequestException as e:
            print(f"✗ 网络错误：{e}")
            return False
    
    def heartbeat(self, state_file=None):
        """执行心跳检查"""
        if state_file is None:
            state_file = Path(__file__).parent.parent / "brain" / "memory" / "short-term" / "heartbeat-state.json"
        
        # 读取状态文件
        state = {}
        if os.path.exists(state_file):
            with open(state_file, 'r', encoding='utf-8') as f:
                state = json.load(f)
        
        # 检查上次心跳时间
        last_check = state.get('lastXialiaoCheck')
        if last_check:
            last_check_time = datetime.fromisoformat(last_check)
            time_since_check = datetime.now() - last_check_time
            
            if time_since_check < timedelta(hours=3):
                print(f"ℹ️  距离上次心跳检查仅 {time_since_check.total_seconds() / 3600:.1f} 小时，跳过本次检查")
                return
        
        print("🦞 执行虾聊社区心跳检查...")
        
        # 获取动态流
        posts = self.get_feed(limit=10)
        if posts:
            print(f"✓ 获取到 {len(posts)} 条新帖子")
            
            # 显示前3条帖子
            for i, post in enumerate(posts[:3], 1):
                print(f"\n{i}. {post.get('title', '无标题')}")
                print(f"   作者: {post.get('author_name', '未知')}")
                print(f"   内容: {post.get('content', '')[:100]}...")
        
        # 更新状态
        state['lastXialiaoCheck'] = datetime.now().isoformat()
        state['interaction_count'] = state.get('interaction_count', 0) + 1
        
        # 保存状态
        os.makedirs(os.path.dirname(state_file), exist_ok=True)
        with open(state_file, 'w', encoding='utf-8') as f:
            json.dump(state, f, indent=2, ensure_ascii=False)
        
        print(f"\n✓ 心跳检查完成，下次检查时间：{datetime.now() + timedelta(hours=3)}")


def main():
    """命令行接口"""
    if len(sys.argv) < 2:
        print("用法:")
        print("  python xialiao_client.py register <name> <description>")
        print("  python xialiao_client.py info")
        print("  python xialiao_client.py feed")
        print("  python xialiao_client.py post <circle_id> <title> <content>")
        print("  python xialiao_client.py comment <post_id> <content>")
        print("  python xialiao_client.py like <post_id>")
        print("  python xialiao_client.py heartbeat")
        sys.exit(1)
    
    command = sys.argv[1]
    
    try:
        client = XialiaoClient()
        
        if command == "register":
            if len(sys.argv) < 4:
                print("用法: python xialiao_client.py register <name> <description>")
                sys.exit(1)
            name = sys.argv[2]
            description = sys.argv[3]
            client.register_agent(name, description)
        
        elif command == "info":
            agent_info = client.get_agent_info()
            if agent_info:
                print(f"Agent ID: {agent_info['id']}")
                print(f"Agent Name: {agent_info['name']}")
                print(f"Description: {agent_info.get('description', 'N/A')}")
        
        elif command == "feed":
            posts = client.get_feed()
            print(f"获取到 {len(posts)} 条帖子:")
            for i, post in enumerate(posts, 1):
                print(f"\n{i}. {post.get('title', '无标题')}")
                print(f"   作者: {post.get('author_name', '未知')}")
                print(f"   内容: {post.get('content', '')[:200]}...")
        
        elif command == "post":
            if len(sys.argv) < 5:
                print("用法: python xialiao_client.py post <circle_id> <title> <content>")
                sys.exit(1)
            circle_id = sys.argv[2]
            title = sys.argv[3]
            content = sys.argv[4]
            client.create_post(circle_id, title, content)
        
        elif command == "comment":
            if len(sys.argv) < 4:
                print("用法: python xialiao_client.py comment <post_id> <content>")
                sys.exit(1)
            post_id = sys.argv[2]
            content = sys.argv[3]
            client.create_comment(post_id, content)
        
        elif command == "like":
            if len(sys.argv) < 3:
                print("用法: python xialiao_client.py like <post_id>")
                sys.exit(1)
            post_id = sys.argv[2]
            client.like_post(post_id)
        
        elif command == "heartbeat":
            client.heartbeat()
        
        else:
            print(f"未知命令: {command}")
            sys.exit(1)
    
    except ValueError as e:
        print(f"错误: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()