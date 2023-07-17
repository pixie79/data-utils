package utils

import (
	"fmt"
	"testing"
)

func TestGetEnv(t *testing.T) {
	result := GetEnv("DEBUG_LEVEL", "TEST")
	if result != "TEST" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result, "TEST")
	}
}

func TestDecryptKey(t *testing.T) {
	var tests = []struct {
		name string
		key  string
		want string
	}{
		{"test1", "AAAAAhACRmNsaWVudC1pbmRpdmlkdWFsLWRldmljZS1ldmVudC12MDAx", "client-individual-device-event-v001"},
		{"test2", "AAAAAhACRGNsaWVudC1pbmRpdmlkdWFsLXBheWVlLWV2ZW50LXYwMDE=", "client-individual-payee-event-v001"},
		{"test3", "AAAAAhACTGNsaWVudC1zdXNwZW5zZS10cmFuc2FjdGlvbi1ldmVudC12MDAy", "client-suspense-transaction-event-v002"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := B64DecodeMsg(tt.key)
			if err != nil {
				t.Errorf(fmt.Sprintf("%+v", err))
			}
			if string(result) != tt.want {
				t.Errorf("Topic name mismatch wanted %s got: %s", tt.want, result)
			}
		})
	}
}
