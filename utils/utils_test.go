package utils

import (
	"fmt"
	"strings"
	"testing"
)

var (
	ip0   = "10.0.0.0"
	ip1   = "10.0.0.1"
	ip2   = "10.0.0.2"
	ip3   = "10.0.0.3"
	ip4   = "10.0.0.4"
	list1 = []string{ip0, ip1}
	list2 = []string{ip1, ip2}
	list3 = []string{ip0}
	list4 = []string{ip4, ip3}
	list5 = []string{ip3, ip4, ip0}
	list6 = []string{ip0, ip3, ip0}
)

// TestGetIPs tests the GetIPs function
func TestGetEnv(t *testing.T) {
	result := GetEnv("DEBUG_LEVEL", "TEST")
	if result != "TEST" {
		t.Errorf("result was incorrect, got: %s, want: %s.", result, "TEST")
	}
}

// TestDecryptKey tests the DecryptKey function
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
				t.Errorf("topic name mismatch wanted %s got: %s", tt.want, result)
			}
		})
	}
}

// TestContains tests the Contains function
func TestContains(t *testing.T) {
	var tests = []struct {
		name  string
		ip    string
		list1 []string
		want  bool
	}{
		{"IP found", ip0, list1, true},
		{"IP not found", ip2, list1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Contains(tt.list1, tt.ip)
			if result != tt.want {
				t.Errorf("test %s, result %+v", tt.name, tt.want)
			}
		})
	}
}

// TestDifferenceInSlices tests the DifferenceInSlices function
func TestDifferenceInSlices(t *testing.T) {
	var tests = []struct {
		name      string
		list1     []string
		list2     []string
		wantList1 []string
		wantList2 []string
		wantList3 []string
	}{
		{"oneResultEachList", list1, list2, []string{ip2}, []string{ip0}, []string{ip1}},
		{"oneResultFirstList", list1, list3, []string{}, []string{ip1}, []string{ip0}},
		{"oneResultSecondList", list3, list1, []string{ip1}, []string{}, []string{ip0}},
		{"newOneResultEachList", list1, list4, []string{ip3, ip4}, []string{ip0, ip1}, []string{}},
		{"unbalanced", list1, list5, []string{ip3, ip4}, []string{ip1}, []string{ip0}},
		{"doubleMatch", list1, list1, []string{}, []string{}, []string{ip0, ip1}},
		{"doubleIP", list1, list6, []string{ip3}, []string{ip1}, []string{ip0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wantList1 []string
			var wantList2 []string
			var wantList3 []string
			wantList1, wantList2, wantList3 = DifferenceInSlices(tt.list1, tt.list2)
			if strings.Compare(strings.Join(tt.wantList1, " "), strings.Join(wantList1, " ")) != 0 {
				t.Errorf("test %s, result %+v wanted %+v", tt.name, wantList1, tt.wantList1)
			}
			if strings.Compare(strings.Join(wantList2, " "), strings.Join(tt.wantList2, " ")) != 0 {
				t.Errorf("test %s, result %+v, wanted %+v", tt.name, wantList2, tt.wantList2)
			}
			if strings.Compare(strings.Join(wantList3, " "), strings.Join(tt.wantList3, " ")) != 0 {
				t.Errorf("test %s, result %+v, wanted %+v", tt.name, wantList3, tt.wantList3)
			}
		})
	}
}
