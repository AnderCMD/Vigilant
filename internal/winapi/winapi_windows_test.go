package winapi

import (
	"testing"
	"unsafe"
)

// SendInput rejects calls whose cbSize differs from sizeof(INPUT).
func TestInputSize(t *testing.T) {
	want := uintptr(40) // 64-bit Windows
	if unsafe.Sizeof(uintptr(0)) == 4 {
		want = 28
	}
	if got := unsafe.Sizeof(input{}); got != want {
		t.Fatalf("sizeof(input) = %d, want %d", got, want)
	}
}

func TestSingleInstance(t *testing.T) {
	ok, err := AcquireSingleInstance()
	if err != nil || !ok {
		t.Fatalf("first acquire: ok=%v err=%v", ok, err)
	}
}
