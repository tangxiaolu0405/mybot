package evolve

import "testing"

func TestAppendCrystallizeTriggers_repeatedBrowser(t *testing.T) {
	s := &Snapshot{ShortTermBytes: 1024}
	excerpt := "browser_snapshot a\nbrowser_snapshot b\nbrowser_snapshot c\nbrowser_navigate x"
	appendCrystallizeTriggers(s, excerpt)
	found := false
	for _, tr := range s.Triggers {
		if tr == "repeated_browser_tools" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("triggers=%v", s.Triggers)
	}
}

func TestAppendCrystallizeTriggers_taskKeyword(t *testing.T) {
	s := &Snapshot{ShortTermBytes: 1024}
	appendCrystallizeTriggers(s, "今日涨停板连板统计")
	found := false
	for _, tr := range s.Triggers {
		if tr == "task_keyword:涨停" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("triggers=%v", s.Triggers)
	}
}

func TestShouldInvokeCrystallize_highToken(t *testing.T) {
	s := &Snapshot{Triggers: []string{"high_token_session"}}
	if !shouldInvokeCrystallize(s) {
		t.Fatal("expected crystallize for high_token_session")
	}
}

func TestShouldInvokeCrystallize_noTriggers(t *testing.T) {
	s := &Snapshot{Triggers: []string{"short_term>=4096B"}}
	if shouldInvokeCrystallize(s) {
		t.Fatal("periodic consolidate trigger should not invoke crystallize alone")
	}
}
