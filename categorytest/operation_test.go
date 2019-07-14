//
// @copyright: 2019 by Pauli Rikula <pauli.rikula@gmail.com>
// @license: MIT <http://www.opensource.org/licenses/mit-license.php>
//
package categorytest

import (
	"category"
	"fmt"
	"testing"
)

func TestOperationDiscard(t *testing.T) {
	connect := NewConnectionPrinter()
	a := NewConnectable("a")
	b := NewConnectable("b")

	opSet := category.NewOperationSet(connect)

	opSet.Add(category.NewFreezedOperation(connect, a, b))
	opSet.Add(category.NewFreezedOperation(connect, a, a))

	fmt.Printf("%+v", opSet.AsSortedArray())

	fmt.Printf("vs\n")

	opSet.Remove(category.NewFreezedOperation(connect, a, a))

	fmt.Printf("%+v", opSet.AsSortedArray())

}
