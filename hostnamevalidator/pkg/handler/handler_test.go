package handler

import (
	"reflect"
	"testing"
)

func TestEvalHostname(t *testing.T) {

	tests := []struct {
		hostname  string
		namespace string
		match     bool
	}{
		{hostname: "app.a.namespace-a.k8s.zach", namespace: "namespace-a", match: true},
		{hostname: "app.a.namespace-a.k8s.zach", namespace: "namespace-b", match: false},
		{hostname: "app.k8s.zach", namespace: "namespace-a", match: false},
		{hostname: "app.a.namespace-a.zach", namespace: "namespace-a", match: false},
		{hostname: "app.namespace-a.k8s.zach", namespace: "namespace-a", match: true},
	}

	for _, tc := range tests {
		valid, _ := evalHostname(tc.hostname, tc.namespace)
		if !reflect.DeepEqual(tc.match, valid) {
			t.Fatalf("hostname: %s, expected: %v, got: %v", tc.hostname, tc.match, valid)
		}
	}

}
