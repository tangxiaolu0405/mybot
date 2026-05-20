package server

import (
	"mybot/internal/llm"
)

// trimHistoryToTokenBudget 从最早的用户/助手/tool 消息裁掉，使估算 token ≤ budget（保留前置 system 若有）。
func trimHistoryToTokenBudget(client *llm.Client, msgs []llm.Message, tools []llm.Tool, budget int) []llm.Message {
	if budget <= 0 || len(msgs) == 0 {
		return msgs
	}
	out := append([]llm.Message(nil), msgs...)
	for len(out) > 1 && client.EstimatedChatInputTokens(out, tools) > budget {
		drop := firstDroppableIndex(out)
		if drop < 0 {
			break
		}
		out = append(out[:drop], out[drop+1:]...)
	}
	return out
}

func firstDroppableIndex(msgs []llm.Message) int {
	for i, m := range msgs {
		switch m.Role {
		case "user", "assistant", "tool":
			return i
		}
	}
	return -1
}
