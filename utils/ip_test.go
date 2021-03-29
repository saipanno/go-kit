package utils

import (
	"reflect"
	"testing"
)

func TestParseIPRangeToCIDRs(t *testing.T) {
	type args struct {
		start string
		end   string
	}
	tests := []struct {
		name      string
		args      args
		wantCIDRs []string
	}{
		// TODO: Add test cases.
		{
			"test1",
			args{
				start: "192.168.1.1",
				end:   "192.168.1.21",
			},
			[]string{
				"192.168.1.1/32",
				"192.168.1.2/31",
				"192.168.1.4/30",
				"192.168.1.8/29",
				"192.168.1.16/30",
				"192.168.1.20/31",
			},
		},
		{
			"test2",
			args{
				start: "192.168.1.1",
				end:   "192.168.1.1",
			},
			[]string{
				"192.168.1.1/32",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCIDRs := ParseIPRangeToCIDRs(tt.args.start, tt.args.end); !reflect.DeepEqual(gotCIDRs, tt.wantCIDRs) {
				t.Errorf("ParseIPRangeToCIDRs() = %v, want %v", gotCIDRs, tt.wantCIDRs)
			}
		})
	}
}

func TestGetIPAddrs(t *testing.T) {
	tests := []struct {
		name        string
		wantPrivate []string
		wantPublic  []string
	}{
		// TODO: Add test cases.
		{
			"test",
			[]string{},
			[]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrivate, gotPublic := GetIPAddrs()
			if !reflect.DeepEqual(gotPrivate, tt.wantPrivate) {
				t.Errorf("GetIPAddrs() gotPrivate = %v, want %v", gotPrivate, tt.wantPrivate)
			}
			if !reflect.DeepEqual(gotPublic, tt.wantPublic) {
				t.Errorf("GetIPAddrs() gotPublic = %v, want %v", gotPublic, tt.wantPublic)
			}
		})
	}
}
