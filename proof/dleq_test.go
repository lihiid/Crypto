package proof

import (
	"fmt"
	"testing"

	"github.com/dedis/crypto/abstract"
	"github.com/dedis/crypto/edwards"
	"github.com/dedis/crypto/random"
	"github.com/stretchr/testify/require"
)

func TestDLEQProof(t *testing.T) {
	suite := edwards.NewAES128SHA256Ed25519(false)
	n := 10
	for i := 0; i < n; i++ {
		// Create some random secrets and base points
		x := suite.Scalar().Pick(random.Stream)
		g, _ := suite.Point().Pick([]byte(fmt.Sprintf("G%d", i)), random.Stream)
		h, _ := suite.Point().Pick([]byte(fmt.Sprintf("H%d", i)), random.Stream)
		proof, xG, xH, err := NewDLEQProof(suite, g, h, x)
		require.Equal(t, err, nil)
		require.Nil(t, proof.Verify(suite, g, h, xG, xH))
	}
}

func TestDLEQProofBatch(t *testing.T) {
	suite := edwards.NewAES128SHA256Ed25519(false)
	n := 10
	x := make([]abstract.Scalar, n)
	g := make([]abstract.Point, n)
	h := make([]abstract.Point, n)
	for i := range x {
		x[i] = suite.Scalar().Pick(random.Stream)
		g[i], _ = suite.Point().Pick([]byte(fmt.Sprintf("G%d", i)), random.Stream)
		h[i], _ = suite.Point().Pick([]byte(fmt.Sprintf("H%d", i)), random.Stream)
	}
	proofs, xG, xH, err := NewDLEQProofBatch(suite, g, h, x)
	require.Equal(t, err, nil)
	for i := range proofs {
		require.Nil(t, proofs[i].Verify(suite, g[i], h[i], xG[i], xH[i]))
	}
}

func TestDLEQLengths(t *testing.T) {
	suite := edwards.NewAES128SHA256Ed25519(false)
	n := 10
	x := make([]abstract.Scalar, n)
	g := make([]abstract.Point, n)
	h := make([]abstract.Point, n)
	for i := range x {
		x[i] = suite.Scalar().Pick(random.Stream)
		g[i], _ = suite.Point().Pick([]byte(fmt.Sprintf("G%d", i)), random.Stream)
		h[i], _ = suite.Point().Pick([]byte(fmt.Sprintf("H%d", i)), random.Stream)
	}
	// Remove an element to make the test fail
	x = append(x[:5], x[6:]...)
	_, _, _, err := NewDLEQProofBatch(suite, g, h, x)
	require.Equal(t, err, errorDifferentLengths)
}
