package shellrunner

import (
	"testing"
)

func TestRunCommandWithEnvUser(t *testing.T) {
	type args struct {
		CMDs     []string
		env      []string
		username string
		timeout  []int
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
		wantExit   int
	}{
		// TODO: Add test cases.
		{
			"test",
			args{
				CMDs: []string{"ls", "-l"},
			},
			"",
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, gotExit := RunCommandWithEnvUser(tt.args.CMDs, tt.args.env, tt.args.username, tt.args.timeout...)
			if gotOutput != tt.wantOutput {
				t.Errorf("RunCommandWithEnvUser() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
			if gotExit != tt.wantExit {
				t.Errorf("RunCommandWithEnvUser() gotExit = %v, want %v", gotExit, tt.wantExit)
			}
		})
	}
}
