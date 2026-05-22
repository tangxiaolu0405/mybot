# zhangtingban-lianban

## 适用

- 东方财富涨停板 / 连板统计（已知 API 或已验证脚本）
- 产出写到产出区：`zhangtingban_analysis/<date>/`

## 不适用

- 其它站点、需登录或新页面结构 → 使用 **browser_*** 探索，勿强行 `run_skill`

## 用法

```json
{ "skill": "zhangtingban-lianban", "params": { "date": "2026-05-21" } }
```

## 输出

- `zhangtingban_analysis/lianban/连板票统计.csv`（相对产出区 cwd）
