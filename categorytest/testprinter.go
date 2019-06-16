//
// @copyright: 2019 by Pauli Rikula <pauli.rikula@gmail.com>
//   @license: MIT <http://www.opensource.org/licenses/mit-license.php>
//

package categorytest

import (
	"category"
	"fmt"
)

// Connectable implements the category.Connectable -interface for testing purposes
type Connectable struct {
	Id string
}

// GetId returns unique identifier value for each connectable instance
func (c *Connectable) GetId() string {
	return c.Id
}

// NewConnectable creates a new category.Connectable instance for testing purposes
func NewConnectable(id string) category.Connectable {
	return &Connectable{Id: id}
}

// PrintConnection simulates the connecting operation for testing purposes
func PrintConnection(a category.Connectable, b category.Connectable) error {
	if a == nil || b == nil {
		return fmt.Errorf("connecting nil")
	}
	fmt.Printf("%s -> %s\n", a.GetId(), b.GetId())
	return nil
}

// ConnectionPrinter implements the category.Operator -interface for testing purposes
type ConnectionPrinter struct {
}

// Evaluate prints a.GetId() -> b.GetId()
func (c *ConnectionPrinter) Evaluate(a category.Connectable, b category.Connectable) error {
	return PrintConnection(a, b)
}

// GetId returns and identifier for this test printer
func (c *ConnectionPrinter) GetId() string {
	return "testprint"
}

// NewConnectionPrinter returns a new category.Operator instance for testing purposes
func NewConnectionPrinter() category.Operator {
	return &ConnectionPrinter{}
}
