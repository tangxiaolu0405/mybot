# CLAW Skills Indexing System

## Overview

The CLAW project now has an automated skills indexing system that makes it easy to discover and manage available skills.

## Files

### Core Files

- **skills/skills-index.json** - Generated index file containing all skill information
- **scripts/generate_skills_index.py** - Script to scan skills directory and generate index
- **scripts/list_skills.py** - Script to list and search skills using the index

## Usage

### Generate Skills Index

```bash
python scripts/generate_skills_index.py
```

This scans the `skills/` directory and creates/updates `skills/skills-index.json` with:
- All skill metadata (name, description, version, author, tags, dependencies)
- Tags index for quick lookup
- Generation timestamp

### List All Skills

```bash
python scripts/list_skills.py
```

Shows a summary of all available skills.

### List Skills with Details

```bash
python scripts/list_skills.py -d
```

Shows detailed information for each skill.

### Show Tags Index

```bash
python scripts/list_skills.py -t
```

Shows all tags and which skills have each tag.

### Filter by Tag

```bash
python scripts/list_skills.py --tag <tag-name>
```

Example:
```bash
python scripts/list_skills.py --tag skill-management
```

### Search by Keyword

```bash
python scripts/list_skills.py --search <keyword>
```

Searches skill names and descriptions.

## Index Structure

The `skills-index.json` file contains:

```json
{
  "version": "1.0.0",
  "generated_at": "ISO timestamp",
  "skills": [
    {
      "name": "skill-name",
      "path": "relative/path/to/SKILL.md",
      "description": "Skill description",
      "version": "1.0.0",
      "author": "Author name",
      "tags": ["tag1", "tag2"],
      "dependencies": []
    }
  ],
  "tags_index": {
    "tag1": ["skill-name"],
    "tag2": ["skill-name"]
  }
}
```

## Integration with AI

When the AI needs to know available skills:

1. Check if `skills/skills-index.json` exists
2. If not, run `python scripts/generate_skills_index.py`
3. Load and parse the JSON index
4. Use the index to display or search skills

## Benefits

- **Fast lookup**: No need to parse markdown files each time
- **Tag-based filtering**: Quick filtering by skill tags
- **Searchable**: Easy keyword search across skills
- **Automated**: Simple script to regenerate index
- **Structured**: JSON format easy to parse programmatically

## Maintenance

Run `python scripts/generate_skills_index.py` whenever:
- A new skill is added
- An existing skill is modified
- Skill metadata changes

The script will automatically update the index with the latest information.
