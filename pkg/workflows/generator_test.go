package workflows

import (
	"strings"
	"testing"
)

func TestGenerateWorkflow(t *testing.T) {
	generator := NewGenerator()

	tests := []struct {
		name    string
		agent   string
		trigger string
		wantErr bool
		want    string
	}{
		{
			name:    "Claude Comment",
			agent:   "Claude",
			trigger: "Comment",
			wantErr: false,
			want:    "name: Claude Agent",
		},
		{
			name:    "Jules Label",
			agent:   "Jules",
			trigger: "Label",
			wantErr: false,
			want:    "name: Jules Agent",
		},
		{
			name:    "Invalid Combination",
			agent:   "Unknown",
			trigger: "Trigger",
			wantErr: true,
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, _, err := generator.GenerateWorkflow(tt.agent, tt.trigger)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !strings.Contains(got, tt.want) {
				t.Errorf("GenerateWorkflow() = %v, want substring %v", got, tt.want)
			}
		})
	}
}
