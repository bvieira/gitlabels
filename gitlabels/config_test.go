package gitlabels

import (
	"reflect"
	"testing"
)

const ownerConfigExample = `owner: user
project-regex: prefix-.*
labels:
  'issue: bug': 
    color: 801515
    description: this is a description`

const orgConfigExample = `org: orgname
project-regex: prefix-.*
labels:
  'issue: bug': 
    color: 801515
    description: this is a description`

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    Config
		wantErr bool
	}{
		{"valid owner config", []byte(ownerConfigExample), Config{Owner: "user", ProjectRegex: "prefix-.*", Labels: map[string]LabelConfig{"issue: bug": LabelConfig{Color: "801515", Description: "this is a description"}}}, false},
		{"valid org config", []byte(orgConfigExample), Config{ORG: "orgname", ProjectRegex: "prefix-.*", Labels: map[string]LabelConfig{"issue: bug": LabelConfig{Color: "801515", Description: "this is a description"}}}, false},
		{"invalid config", []byte("	"), Config{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConfig(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestConfig_getUser(t *testing.T) {
	type fields struct {
		Owner string
		ORG   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"owner", fields{Owner: "owner"}, "owner"},
		{"org", fields{ORG: "org"}, "org"},
		{"org precedence", fields{Owner: "owner", ORG: "org"}, "org"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Owner: tt.fields.Owner,
				ORG:   tt.fields.ORG,
			}
			if got := c.getUser(); got != tt.want {
				t.Errorf("Config.getUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
