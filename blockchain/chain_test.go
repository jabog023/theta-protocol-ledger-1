package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thetatoken/ukulele/core"
)

func TestBlockchain(t *testing.T) {
	assert := assert.New(t)
	core.ResetTestBlocks()

	expected := CreateTestChainByBlocks([]string{
		"a1", "a0",
		"a2", "a1",
		"b2", "a1",
		"c1", "a0"})
	var err error

	chain := CreateTestChain()
	a1 := core.CreateTestBlock("a1", "a0")
	_, err = chain.AddBlock(a1)
	assert.Nil(err)

	a2 := core.CreateTestBlock("a2", "a1")
	_, err = chain.AddBlock(a2)
	assert.Nil(err)

	b2 := core.CreateTestBlock("b2", "a1")
	_, err = chain.AddBlock(b2)
	assert.Nil(err)

	c1 := core.CreateTestBlock("c1", "a0")
	_, err = chain.AddBlock(c1)
	assert.Nil(err)

	AssertChainsEqual(assert, expected, expected.Root.Hash(), chain, chain.Root.Hash())
}

func TestBlockchainDeepestDescendant(t *testing.T) {
	assert := assert.New(t)
	core.ResetTestBlocks()
	ch := CreateTestChainByBlocks([]string{
		"a1", "a0",
		"a2", "a1",
		"b2", "a1",
		"b3", "b2",
		"c1", "a0"})

	ret, depth := ch.FindDeepestDescendant(ch.Root.Hash())
	assert.Equal(core.GetTestBlock("b3").Hash(), ret.Hash())
	assert.Equal(3, depth)
}

func TestFinalizePreviousBlocks(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	core.ResetTestBlocks()

	ch := CreateTestChainByBlocks([]string{
		"a1", "a0",
		"a2", "a1",
		"a3", "a2",
		"a4", "a3",
		"a5", "a4",
		"b2", "a1",
		"b3", "b2",
		"c1", "a0",
	})
	block, err := ch.FindBlock(core.GetTestBlock("a3").Hash())
	require.Nil(err)

	ch.FinalizePreviousBlocks(block)

	for _, name := range []string{"a0", "a1", "a2", "a3"} {
		block, err = ch.FindBlock(core.GetTestBlock(name).Hash())
		assert.Nil(err)
		assert.Equal(core.BlockStatusFinalized, block.Status)
	}

	for _, name := range []string{"b2", "b3", "c1", "a4", "a5"} {
		block, err = ch.FindBlock(core.GetTestBlock(name).Hash())
		assert.NotEqual(core.BlockStatusFinalized, block.Status)
	}

	block, err = ch.FindBlock(core.GetTestBlock("a5").Hash())
	require.Nil(err)
	ch.FinalizePreviousBlocks(block)

	for _, name := range []string{"a0", "a1", "a2", "a3", "a4", "a5"} {
		block, err = ch.FindBlock(core.GetTestBlock(name).Hash())
		assert.Equal(core.BlockStatusFinalized, block.Status)
	}

	for _, name := range []string{"b2", "b3", "c1"} {
		block, err = ch.FindBlock(core.GetTestBlock(name).Hash())
		assert.NotEqual(core.BlockStatusFinalized, block.Status)
	}

}
