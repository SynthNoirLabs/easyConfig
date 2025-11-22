package main

import (
	"testing"
)

func TestApp_Greet(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		app  *App
		args args
		want string
	}{
		{
			name: "Greeting with name",
			app:  NewApp(),
			args: args{name: "World"},
			want: "Hello World, It's show time!",
		},
		{
			name: "Greeting with empty name",
			app:  NewApp(),
			args: args{name: ""},
			want: "Hello , It's show time!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.app.Greet(tt.args.name); got != tt.want {
				t.Errorf("App.Greet() = %v, want %v", got, tt.want)
			}
		})
	}
}
