package craft

import (
	"testing"

	"FeedCraft/internal/dao"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveCraftAtoms_UnknownCraftNameReturnsError(t *testing.T) {
	db := newCraftRuntimeTestDB(t)

	_, err := ResolveCraftAtoms(db, "does-not-exist")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a valid craft name")
}

func TestResolveCraftAtoms_CommaSeparatedTopLevel(t *testing.T) {
	db := newCraftRuntimeTestDB(t)

	atoms, err := ResolveCraftAtoms(db, "proxy, limit ,guid-fix")
	require.NoError(t, err)
	require.Len(t, atoms, 3)
	assert.Equal(t, "proxy", atoms[0].TemplateName)
	assert.Equal(t, "limit", atoms[1].TemplateName)
	assert.Equal(t, "guid-fix", atoms[2].TemplateName)
}

func TestResolveCraftAtoms_CommaSeparatedSkipsEmptyParts(t *testing.T) {
	db := newCraftRuntimeTestDB(t)

	atoms, err := ResolveCraftAtoms(db, "proxy,,limit,")
	require.NoError(t, err)
	require.Len(t, atoms, 2)
	assert.Equal(t, "proxy", atoms[0].TemplateName)
	assert.Equal(t, "limit", atoms[1].TemplateName)
}

func TestResolveCraftAtoms_InvalidTemplateNameOnCustomAtom(t *testing.T) {
	db := newCraftRuntimeTestDB(t)
	require.NoError(t, dao.CreateCraftAtom(db, &dao.CraftAtom{
		Name:         "broken-atom",
		TemplateName: "no-such-template",
	}))

	_, err := ResolveCraftAtoms(db, "broken-atom")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid tmpl name")
}

func TestResolveCraftAtoms_NestedFlowExceedsMaxDepth(t *testing.T) {
	db := newCraftRuntimeTestDB(t)
	// A self-referencing flow guarantees recursion until the max depth guard trips.
	require.NoError(t, dao.CreateCraftFlow(db, &dao.CraftFlow{
		Name: "loop-flow",
		CraftFlowConfig: []dao.CraftFlowItem{
			{CraftName: "loop-flow"},
		},
	}))

	_, err := ResolveCraftAtoms(db, "loop-flow")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "max call depth hit")
}

func TestResolveCraftAtoms_NestedFlowWithinDepthResolves(t *testing.T) {
	db := newCraftRuntimeTestDB(t)
	require.NoError(t, dao.CreateCraftFlow(db, &dao.CraftFlow{
		Name: "inner-flow",
		CraftFlowConfig: []dao.CraftFlowItem{
			{CraftName: "limit"},
			{CraftName: "guid-fix"},
		},
	}))
	require.NoError(t, dao.CreateCraftFlow(db, &dao.CraftFlow{
		Name: "outer-flow",
		CraftFlowConfig: []dao.CraftFlowItem{
			{CraftName: "proxy"},
			{CraftName: "inner-flow"},
		},
	}))

	atoms, err := ResolveCraftAtoms(db, "outer-flow")
	require.NoError(t, err)
	require.Len(t, atoms, 3)
	assert.Equal(t, "proxy", atoms[0].TemplateName)
	assert.Equal(t, "limit", atoms[1].TemplateName)
	assert.Equal(t, "guid-fix", atoms[2].TemplateName)
}
