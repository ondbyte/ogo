/* a simple inventory management system in Go.
Design a set of Go types to represent different items in the inventory,
including Product, Equipment, and Material. Each item should have common attributes like ID,
Name, and Quantity, as well as specific attributes and methods based on their types.
Implement methods to update quantities, add new items, and list all items in the inventory.
*/

package cb

import "fmt"

type IBase interface {
	UpdateQty(qty uint)
}

type Base struct {
	Id   string
	Name string
	Qty  uint
}

func (b *Base) UpdateQty(qty uint) {
	b.Qty = qty
}

type Product struct {
	Base
	ProductCategory string
}

type Equipment struct {
	Base
	EquipmentCategory string
}

type Material struct {
	Base
	MaterialCategory string
}

type Inventory struct {
	items []IBase
}

func Example() {
	product := &Product{
		Base: Base{
			Id:   "xyz",
			Name: "xyz",
			Qty:  0,
		},
	}
	b := []IBase{}
	b = append(b, product)
	product.UpdateQty(2)

	fmt.Println(product.Qty)
}
