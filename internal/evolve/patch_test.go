package evolve

import (
	"testing"

	"mybot/internal/brain"
)

func TestNormalizeWorkspaceRel_skills(t *testing.T) {
	cases := []struct {
		path string
		ok   bool
	}{
		{"skills/foo/SKILL.md", true},
		{"skills/foo/manifest.yaml", true},
		{"skills/foo/script.py", true},
		{"skills/foo/evil.exe", false},
		{"modes/_default/capabilities.yaml", false},
	}
	for _, c := range cases {
		_, err := normalizeWorkspaceRel(c.path)
		if c.ok && err != nil {
			t.Fatalf("%s: %v", c.path, err)
		}
		if !c.ok && err == nil {
			t.Fatalf("expected error for %s", c.path)
		}
	}
}

func TestRejectCapabilitiesPatch(t *testing.T) {
	if err := brain.RejectCapabilitiesPatch("modes/_default/capabilities.yaml", "overwrite", "skills: []\nmcp: []\n"); err == nil {
		t.Fatal("expected reject clear mcp")
	}
}
