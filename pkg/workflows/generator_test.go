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
			resp, err := generator.GenerateWorkflow(tt.agent, tt.trigger)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateWorkflow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if resp == nil {
					t.Error("GenerateWorkflow() returned nil response")
					return
				}
				if !strings.Contains(resp.Content, tt.want) {
					t.Errorf("GenerateWorkflow() content = %v, want substring %v", resp.Content, tt.want)
				}
			}
		})
	}
}

func TestListTemplates(t *testing.T) {
	g := NewGenerator()
	list := g.ListTemplates()
	if len(list) == 0 {
		t.Fatalf("expected templates to be non-empty")
	}
	for _, tmpl := range list {
		if tmpl.ID == "" || tmpl.Name == "" || tmpl.Content == "" {
			t.Fatalf("template missing required fields: %+v", tmpl)
		}
		if len(tmpl.DefaultFilename) == 0 {
			t.Fatalf("template %s missing default filename", tmpl.ID)
		}
	}
}
