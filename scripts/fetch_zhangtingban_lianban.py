#!/usr/bin/env python3
"""Dev reference for evolve crystallize → brain skills/zhangtingban-lianban/script.py

Writes CSV under output cwd (not brain dir). Params: optional JSON argv[1] with {"date": "YYYY-MM-DD"}.
"""
import json
import os
import sys
from datetime import date


def main() -> int:
    params = {}
    if len(sys.argv) > 1:
        try:
            params = json.loads(sys.argv[1])
        except json.JSONDecodeError:
            pass
    day = params.get("date") or date.today().isoformat()
    out_dir = os.path.join("zhangtingban_analysis", "lianban")
    os.makedirs(out_dir, exist_ok=True)
    out_path = os.path.join(out_dir, "连板票统计.csv")
    # Placeholder: replace with real fetch; evolve should paste verified logic from short-term.
    with open(out_path, "w", encoding="utf-8") as f:
        f.write("date,code,name,boards\n")
        f.write(f"{day},000000,placeholder,0\n")
    print(f"wrote {out_path} (cwd={os.getcwd()})")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
