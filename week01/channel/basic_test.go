package channel

import "testing"

func TestSendAndReceive(t *testing.T) {
	got := SendAndReceive(42)
	if got != 42 {
		t.Errorf("SendAndReceive(42) = %d, want 42", got)
	}
}

func TestPipeline(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{5, 11},  // 5*2 + 1 = 11
		{0, 1},   // 0*2 + 1 = 1
		{-3, -5}, // -3*2 + 1 = -5
	}
	for _, tt := range tests {
		got := Pipeline(tt.input)
		if got != tt.want {
			t.Errorf("Pipeline(%d) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
