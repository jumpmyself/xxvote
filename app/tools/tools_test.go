package tools

import (
	"fmt"
	"testing"
)

func TestGetUUID(t *testing.T) {
	GetUUID()
}

func TestGetUid(t *testing.T) {
	id := GetUid()
	fmt.Printf("id: %d\n", id)

	id1 := GetUid()
	fmt.Printf("id1: %d\n", id1)

	id2 := GetUid()
	fmt.Printf("id2: %d\n", id2)

	id3 := GetUid()
	fmt.Printf("id3: %d\n", id3)
}
