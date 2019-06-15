//
// @copyright: 2019 by Pauli Rikula <pauli.rikula@gmail.com>
//   @license: MIT <http://www.opensource.org/licenses/mit-license.php>
//
package categorytest

import (
	"category"
	"fmt"
)

type Connectable struct {
	Id string
}

func (c *Connectable) GetId() string {
	return c.Id
}

func NewConnectable(id string) category.Connectable {
	return &Connectable{Id: id}
}

func PrintConnection(a category.Connectable, b category.Connectable) error {
	if a == nil || b == nil {
		return fmt.Errorf("connecting nil")
	}
	fmt.Printf("%s -> %s\n", a.GetId(), b.GetId())
	return nil
}

type ConnectionPrinter struct {
}

func (c *ConnectionPrinter) Evaluate(a category.Connectable, b category.Connectable) error {
	return PrintConnection(a, b)
}

func (c *ConnectionPrinter) GetId() string {
	return "print"
}

func NewConnectionPrinter() category.Operator {
	return &ConnectionPrinter{}
}
