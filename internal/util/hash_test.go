package util

import "testing"

func TestGetMD5HashReturnsHexEncodedMD5(t *testing.T) {
	got := GetMD5Hash("FeedCraft")
	const want = "01ec1a92ca7590baf10433c270707cd5"
	if got != want {
		t.Fatalf("GetMD5Hash() = %q, want %q", got, want)
	}
}
