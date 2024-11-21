package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func emptyStorage(t *testing.T) *Storage {
	st, err := NewStorage(":memory:")
	if err != nil {
		t.Fatalf("error creating storage: %v", err)
	}
	return st
}
func TestStorageListCategories(t *testing.T) {
	st := emptyStorage(t)

	cats, err := st.ListCategories(1)
	require.NoError(t, err)
	require.Len(t, cats, 0)

	st.AddLink(2, "category", "link", "headline")

	cats, err = st.ListCategories(1)
	require.NoError(t, err)
	require.Len(t, cats, 0)

	cats, err = st.ListCategories(2)
	require.NoError(t, err)
	require.Len(t, cats, 1)
	require.EqualValues(t, cats[0], CategoryCounts{"category", 1})
}
