#!/usr/bin/env python3
import json
import sys
from pathlib import Path
from typing import List, Dict, Any


def load_skills_index(index_file: str) -> Dict[str, Any]:
    """加载技能索引文件"""
    index_path = Path(index_file)
    
    if not index_path.exists():
        print(f"Error: Skills index file not found: {index_file}")
        print("Please run: python scripts/generate_skills_index.py")
        return {}
    
    with open(index_path, 'r', encoding='utf-8') as f:
        return json.load(f)


def display_skills_summary(skills: List[Dict[str, Any]]):
    """显示技能摘要"""
    print(f"\n{'='*60}")
    print(f"Available Skills ({len(skills)} total)")
    print(f"{'='*60}\n")
    
    for i, skill in enumerate(skills, 1):
        print(f"{i}. **{skill['name']}** (v{skill['version']})")
        print(f"   Description: {skill['description']}")
        print(f"   Author: {skill['author']}")
        print(f"   Tags: {', '.join(skill['tags'])}")
        print()


def display_skills_detailed(skills: List[Dict[str, Any]]):
    """显示技能详细信息"""
    print(f"\n{'='*60}")
    print(f"Skills Detailed Information ({len(skills)} total)")
    print(f"{'='*60}\n")
    
    for skill in skills:
        print(f"Skill: {skill['name']}")
        print(f"Version: {skill['version']}")
        print(f"Author: {skill['author']}")
        print(f"Description: {skill['description']}")
        print(f"Tags: {skill['tags']}")
        print(f"Dependencies: {skill['dependencies']}")
        print(f"Path: {skill['path']}")
        print(f"{'-'*60}\n")


def filter_by_tag(skills: List[Dict[str, Any]], tag: str, tags_index: Dict[str, List[str]]) -> List[Dict[str, Any]]:
    """按标签过滤技能"""
    if tag not in tags_index:
        print(f"No skills found with tag: {tag}")
        return []
    
    skill_names = tags_index[tag]
    filtered_skills = [s for s in skills if s['name'] in skill_names]
    return filtered_skills


def search_by_keyword(skills: List[Dict[str, Any]], keyword: str) -> List[Dict[str, Any]]:
    """按关键词搜索技能"""
    keyword_lower = keyword.lower()
    filtered_skills = []
    
    for skill in skills:
        if (keyword_lower in skill['name'].lower() or
            keyword_lower in skill['description'].lower()):
            filtered_skills.append(skill)
    
    return filtered_skills


def display_tags_index(tags_index: Dict[str, List[str]]):
    """显示标签索引"""
    print(f"\n{'='*60}")
    print(f"Tags Index ({len(tags_index)} tags)")
    print(f"{'='*60}\n")
    
    for tag, skill_names in sorted(tags_index.items()):
        print(f"{tag}: {', '.join(skill_names)}")
    print()


def main():
    script_dir = Path(__file__).parent
    project_root = script_dir.parent
    index_file = project_root / 'skills' / 'skills-index.json'
    
    index_data = load_skills_index(str(index_file))
    
    if not index_data:
        sys.exit(1)
    
    skills = index_data.get('skills', [])
    tags_index = index_data.get('tags_index', {})
    
    print(f"Index generated at: {index_data.get('generated_at', 'Unknown')}")
    
    args = sys.argv[1:]
    
    if not args:
        display_skills_summary(skills)
    elif args[0] == '--detailed' or args[0] == '-d':
        display_skills_detailed(skills)
    elif args[0] == '--tags' or args[0] == '-t':
        display_tags_index(tags_index)
    elif args[0] == '--tag' and len(args) > 1:
        tag = args[1]
        filtered = filter_by_tag(skills, tag, tags_index)
        if filtered:
            display_skills_summary(filtered)
    elif args[0] == '--search' and len(args) > 1:
        keyword = args[1]
        filtered = search_by_keyword(skills, keyword)
        if filtered:
            print(f"\nFound {len(filtered)} skill(s) matching '{keyword}':")
            display_skills_summary(filtered)
        else:
            print(f"\nNo skills found matching '{keyword}'")
    else:
        print("Usage:")
        print("  python scripts/list_skills.py              - Show skills summary")
        print("  python scripts/list_skills.py -d           - Show detailed information")
        print("  python scripts/list_skills.py -t           - Show tags index")
        print("  python scripts/list_skills.py --tag <tag>   - Filter by tag")
        print("  python scripts/list_skills.py --search <keyword> - Search by keyword")


if __name__ == '__main__':
    main()
