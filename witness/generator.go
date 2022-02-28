package witness

import (
	"time"

	"github.com/vivijj/zigo/internal/smt"
	"github.com/vivijj/zigo/pkg/account"
	"github.com/vivijj/zigo/pkg/block"
)

// Generator will generate witness and store it in the db.
type Generator struct {
	CircuitAccTree smt.SparseMerkleTree[account.CircuitAccount]
}

func (g Generator) PrepareWitnessAndSave(blk block.Block) {
	start := time.Now()
}

func buildBlockWitness() {

}
