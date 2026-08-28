package fixtures

import "testing"

func TestDeterministicFixtureShape(t *testing.T) {
	items := []string{"声乐", "舞蹈", "戏剧", "器乐"}
	if len(items) != 4 || items[0] != "声乐" {
		t.Fatal("fixture mismatch")
	}
}
