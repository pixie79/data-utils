package aws

import (
	"github.com/pixie79/data-utils"
	"reflect"
	"testing"
)

func TestFetchCredentials(t *testing.T) {
	type args struct {
		credentialsKey string
	}
	var tests []struct {
		name string
		args args
		want data_utils.CredentialsType
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FetchCredentials(tt.args.credentialsKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSecretManagerValue(t *testing.T) {
	type args struct {
		passwordKey string
	}
	var tests []struct {
		name string
		args args
		want string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSecretManagerValue(tt.args.passwordKey); got != tt.want {
				t.Errorf("GetSecretManagerValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSsmParam(t *testing.T) {
	type args struct {
		parameterPath string
	}
	var tests []struct {
		name string
		args args
		want string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSsmParam(tt.args.parameterPath); got != tt.want {
				t.Errorf("GetSsmParam() = %v, want %v", got, tt.want)
			}
		})
	}
}
