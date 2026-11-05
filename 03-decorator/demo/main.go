// Demo：Starbuzz 咖啡——点单、加料、算钱。
//
//	go run ./03-decorator/demo
package main

import (
	"fmt"

	decorator "github.com/spigcoder/design-patterns/03-decorator"
)

func main() {
	fmt.Println("=== 1. 菜单上的基础饮料 ===")
	for _, b := range []decorator.Beverage{
		decorator.Espresso{},
		decorator.HouseBlend{},
		decorator.DarkRoast{},
		decorator.Decaf{},
	} {
		show(b)
	}

	fmt.Println("\n=== 2. 按顾客要求逐层加料 ===")
	orders := []decorator.Beverage{
		decorator.NewMocha(decorator.Espresso{}),
		decorator.NewWhip(decorator.NewMocha(decorator.NewMocha(decorator.DarkRoast{}))),
		decorator.NewWhip(decorator.NewMocha(decorator.NewSoy(decorator.HouseBlend{}))),
		decorator.NewMilk(decorator.NewMilk(decorator.Decaf{})),
	}

	total := 0.0
	for _, b := range orders {
		show(b)
		total += b.Cost()
	}
	fmt.Printf("—— 本单合计 %.2f 元 ——\n", total)

	fmt.Println("\n=== 3. 同一杯基础饮料，按顾客勾选在运行时拼装 ===")
	checkout := func(base decorator.Beverage, extras ...string) decorator.Beverage {
		for _, e := range extras {
			switch e {
			case "牛奶":
				base = decorator.NewMilk(base)
			case "摩卡":
				base = decorator.NewMocha(base)
			case "大豆":
				base = decorator.NewSoy(base)
			case "奶泡":
				base = decorator.NewWhip(base)
			}
		}
		return base
	}

	show(checkout(decorator.HouseBlend{}))
	show(checkout(decorator.HouseBlend{}, "摩卡"))
	show(checkout(decorator.HouseBlend{}, "摩卡", "摩卡", "奶泡"))
	show(checkout(decorator.HouseBlend{}, "大豆", "摩卡", "奶泡"))

	fmt.Println("\n注意：上面 4 张单子只用了 1 个饮料类型 + 4 个调料类型，")
	fmt.Println("没有为任何一种组合新建类型 —— 这就是装饰器替代'组合爆炸式继承'的地方。")
}

// show 只认识 Beverage 接口，不知道也不关心背后被包了几层。
func show(b decorator.Beverage) {
	fmt.Printf("  %s —— %.2f 元\n", b.Description(), b.Cost())
}
