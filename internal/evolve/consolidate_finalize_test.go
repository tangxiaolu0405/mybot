package evolve

import "testing"

func TestShouldFinalizeShortTerm(t *testing.T) {
	snap := &Snapshot{ShortTermBytes: 100}
	dec := &Decision{Action: "consolidate"}
	if shouldFinalizeShortTerm(dec, []string{"modes/_default/persona.md"}, snap, false) {
		t.Fatal("expected false without enough bytes")
	}
	snap.ShortTermBytes = shortTermActivityBytes
	if !shouldFinalizeShortTerm(dec, []string{"modes/_default/persona.md"}, snap, false) {
		t.Fatal("expected true for consolidate + persona")
	}
	if shouldFinalizeShortTerm(&Decision{Action: "idle"}, nil, snap, false) {
		t.Fatal("expected false for idle")
	}
	if !shouldFinalizeShortTerm(&Decision{Action: "consolidate"}, []string{"memory/long/note.md"}, snap, true) {
		t.Fatal("expected session compress finalize")
	}
}
