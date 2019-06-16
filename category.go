//
// @copyright: 2019 by Pauli Rikula <pauli.rikula@gmail.com>
//   @license: MIT <http://www.opensource.org/licenses/mit-license.php>
//

package category

// implementation details

type categoryImpl struct {
	Sources    ConnectableSet
	Sinks      ConnectableSet
	Operator   Operator
	Operations OperationSet
	isZero     bool
	isIdentity bool
	stringImpl func(c *categoryImpl) string
}

func (c *categoryImpl) GetSources() ConnectableSet {
	return c.Sources
}

func (c *categoryImpl) GetSinks() ConnectableSet {
	return c.Sinks
}
func (c *categoryImpl) GetOperator() Operator {
	return c.Operator
}
func (c *categoryImpl) GetOperations() OperationSet {
	return c.Operations
}

func (c *categoryImpl) Equals(another Category) bool {
	return another != nil && EqualOperators(
		c.Operator, another.GetOperator()) && c.Sources.Equals(
		another.GetSources()) && c.Sinks.Equals(
		another.GetSinks()) && c.Operations.Equals(another.GetOperations())

}
func (c *categoryImpl) IsZero() bool {
	return c.isZero
}

func (c *categoryImpl) IsIdentity() bool {
	return c.isIdentity
}

func (c *categoryImpl) Evaluate() error {
	for _, f := range c.Operations.AsArray() {
		err := f.Evaluate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *categoryImpl) EvaluateSorted() error {
	for _, f := range c.Operations.AsSortedArray() {
		err := f.Evaluate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *categoryImpl) String() string {
	return c.stringImpl(c)
}
