package test

import "testing"

func IntegrationTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}
