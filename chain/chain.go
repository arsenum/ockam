package chain

import (
	"fmt"

	"github.com/ockam-network/ockam"
	"github.com/ockam-network/ockam/claim"
)

// Chain represents a local instace of the ockam blockchain that is maintained by
// a network of ockam nodes.
//
// It implements the ockam.Chain interface and communicates with the ockam network
// using a trusted local (typically) ockam.Node instance.
type Chain struct {
	id          string
	trustedNode ockam.Node
}

// Option is used to provide optional arguments to the New function which creates a Chain instace.
type Option func(*Chain)

// New returns a new Chain instace
func New(options ...Option) (*Chain, error) {
	// create the struct
	c := &Chain{}

	// apply options to the new struct
	for _, option := range options {
		option(c)
	}

	// return the new Chain struct
	return c, nil
}

// ID is used to optionally set the id of a new Chain struct that is created
func ID(id string) Option {
	return func(c *Chain) {
		c.id = id
	}
}

// TrustedNode is used to optionally set the trustedNode of a new Chain
func TrustedNode(n ockam.Node) Option {
	return func(c *Chain) {
		c.trustedNode = n
	}
}

// ID returns the identifier of the chain
func (c *Chain) ID() string {
	return c.id
}

// Sync causes the Chain's trusted node to synchronize its state with its network
// Which includes fetching the latest block.
func (c *Chain) Sync() error {
	return c.trustedNode.Sync()
}

// LatestBlock returns the latest block in the chain
func (c *Chain) LatestBlock() ockam.Block {
	return c.trustedNode.LatestBlock()
}

// Register is
func (c *Chain) Register(e ockam.Entity) (ockam.Claim, error) {
	cl, err := claim.New(
		claim.Data{"id": e.ID().String()},
		claim.Issuer(e),
		claim.Subject(e),
	)
	if err != nil {
		return nil, err
	}

	err = c.Submit(cl)
	if err != nil {
		return nil, err
	}

	return cl, nil
}

// Submit is
func (c *Chain) Submit(cl ockam.Claim) error {
	bin, err := cl.MarshalBinary()
	if err != nil {
		return err
	}

	tx, err := c.trustedNode.Submit(bin)
	if err != nil {
		return err
	}

	fmt.Println(tx)

	return err
}