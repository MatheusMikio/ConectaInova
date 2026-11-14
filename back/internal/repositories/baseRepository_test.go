package repositories

import (
	"testing"

	"conectainova/internal/domain/schema"
)

func TestGetAllRejectsOverflowBeforeTouchingDatabase(t *testing.T) {
	repo := BaseRepository[schema.Startup]{}
	cases := []struct {
		page uint64
		size uint
	}{
		{1<<63 + 1, 2}, // Previously wrapped the offset to zero.
		{^uint64(0), 30},
		{1<<32 + 1, 12},
	}
	for _, tt := range cases {
		if tt.page > uint64(^uint(0)) {
			continue
		}
		// DB is intentionally nil. Any query would panic.
		result, err := repo.GetAll(uint(tt.page), tt.size)
		if err != nil || result == nil || len(result) != 0 {
			t.Fatalf("page=%d size=%d: expected empty array, got %v, %v", tt.page, tt.size, result, err)
		}
	}
}
