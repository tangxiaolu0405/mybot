package brain

import "testing"

func TestParseCapabilitiesYAML(t *testing.T) {
	c := ParseCapabilitiesYAML([]byte("skills:\n  - babysit\nmcp:\n  - browser\n"))
	if len(c.Skills) != 1 || c.Skills[0] != "babysit" {
		t.Fatalf("skills=%v", c.Skills)
	}
	if !c.AllowsMCPServer("browser") {
		t.Fatal("expected browser allowed")
	}
}
