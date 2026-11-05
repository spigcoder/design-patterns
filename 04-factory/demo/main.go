// Demo：比萨店——工厂方法 + 抽象工厂叠起来用。
//
//	go run ./04-factory/demo
package main

import (
	"fmt"
	"strings"

	factory "github.com/spigcoder/design-patterns/04-factory"
)

func main() {
	fmt.Println("=== 1. 两家店各点一份奶酪披萨：同一款产品，不同原料族 ===")
	for _, store := range []*factory.PizzaStore{
		factory.NewYorkPizzaStore(),
		factory.ChicagoPizzaStore(),
	} {
		printOrder(store, "cheese")
	}

	fmt.Println("\n=== 2. 芝加哥店的蛤蜊披萨：流程完全一样，只是钩子返回别的产品 ===")
	printOrder(factory.ChicagoPizzaStore(), "clam")

	fmt.Println("\n=== 3. 披萨只认工厂接口：同一块奶酪披萨，换一个原料工厂就换了一族原料 ===")
	for _, ing := range []factory.PizzaIngredientFactory{
		factory.NewYorkIngredientFactory{},
		factory.ChicagoIngredientFactory{},
	} {
		p := factory.NewCheesePizza(ing)
		fmt.Printf("  %T 的原料：\n%s\n", ing, indent(p.Prepare(), "    "))
	}

	fmt.Println("\n=== 4. 点一个菜单上没有的品类 ===")
	if _, err := factory.NewYorkPizzaStore().OrderPizza("榴莲"); err != nil {
		fmt.Println("  " + err.Error())
	}
}

func printOrder(store *factory.PizzaStore, kind string) {
	order, err := store.OrderPizza(kind)
	if err != nil {
		fmt.Println("  " + err.Error())
		return
	}
	fmt.Printf("  【%s】\n%s\n", order.Store, indent(strings.Join(order.Steps, "\n"), "    "))
}

func indent(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = prefix + lines[i]
	}
	return strings.Join(lines, "\n")
}
