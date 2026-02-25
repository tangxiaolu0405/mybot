#!/usr/bin/env python3
import os
import json
import re
from datetime import datetime
from pathlib import Path
from typing import Dict, List, Any


def parse_yaml_frontmatter(content: str) -> Dict[str, Any]:
    """解析YAML frontmatter"""
    frontmatter_match = re.match(r'^---\n(.*?)\n---', content, re.DOTALL)
    if not frontmatter_match:
        return {}
    
    yaml_content = frontmatter_match.group(1)
    result = {}
    
    for line in yaml_content.split('\n'):
        if ':' in line:
            key, value = line.split(':', 1)
            key = key.strip()
            value = value.strip()
            
            if key in ['tags', 'dependencies']:
                if value.startswith('[') and value.endswith(']'):
                    value = [item.strip().strip('"\'') for item in value[1:-1].split(',') if item.strip()]
                else:
                    value = []
            
            result[key] = value
    
    return result


def scan_skills_directory(skills_dir: str) -> List[Dict[str, Any]]:
    """扫描skills目录，读取所有技能"""
    skills = []
    skills_path = Path(skills_dir)
    
    if not skills_path.exists():
        print(f"Skills directory not found: {skills_dir}")
        return skills
    
    for skill_dir in skills_path.iterdir():
        if skill_dir.is_dir():
            skill_file = skill_dir / "SKILL.md"
            if skill_file.exists():
                with open(skill_file, 'r', encoding='utf-8') as f:
                    content = f.read()
                
                frontmatter = parse_yaml_frontmatter(content)
                
                if frontmatter.get('name'):
                    skill_info = {
                        'name': frontmatter.get('name', ''),
                        'path': str(skill_file.relative_to(Path(skills_dir).parent)),
                        'description': frontmatter.get('description', ''),
                        'version': frontmatter.get('version', '1.0.0'),
                        'author': frontmatter.get('author', ''),
                        'tags': frontmatter.get('tags', []),
                        'dependencies': frontmatter.get('dependencies', [])
                    }
                    skills.append(skill_info)
                    print(f"Found skill: {skill_info['name']}")
    
    return skills


def build_tags_index(skills: List[Dict[str, Any]]) -> Dict[str, List[str]]:
    """构建标签索引"""
    tags_index = {}
    
    for skill in skills:
        for tag in skill.get('tags', []):
            if tag not in tags_index:
                tags_index[tag] = []
            tags_index[tag].append(skill['name'])
    
    return tags_index


def generate_skills_index(skills_dir: str, output_file: str):
    """生成技能索引文件"""
    print(f"Scanning skills directory: {skills_dir}")
    
    skills = scan_skills_directory(skills_dir)
    
    if not skills:
        print("No skills found!")
        return
    
    tags_index = build_tags_index(skills)
    
    index_data = {
        'version': '1.0.0',
        'generated_at': datetime.now().isoformat(),
        'skills': skills,
        'tags_index': tags_index
    }
    
    output_path = Path(output_file)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    
    with open(output_path, 'w', encoding='utf-8') as f:
        json.dump(index_data, f, indent=2, ensure_ascii=False)
    
    print(f"\nGenerated skills index: {output_file}")
    print(f"Total skills: {len(skills)}")
    print(f"Total tags: {len(tags_index)}")


def main():
    script_dir = Path(__file__).parent
    project_root = script_dir.parent
    skills_dir = project_root / 'skills'
    output_file = project_root / 'skills' / 'skills-index.json'
    
    generate_skills_index(str(skills_dir), str(output_file))


if __name__ == '__main__':
    main()
