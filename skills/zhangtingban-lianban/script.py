#!/usr/bin/env python3
"""Workspace brain copy target for evolve; dev copy of scripts/fetch_zhangtingban_lianban.py"""
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
    with open(out_path, "w", encoding="utf-8") as f:
        f.write("date,code,name,boards\n")
        f.write(f"{day},000000,placeholder,0\n")
    print(f"wrote {out_path} (cwd={os.getcwd()})")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
