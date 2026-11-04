// Demo：跑一遍《Head First 设计模式》第 1 章的 SimUDuck 场景。
//
//	go run ./01-strategy/demo
package main

import (
	"fmt"

	strategy "github.com/spigcoder/design-patterns/01-strategy"
)

func main() {
	// 关键：这个切片的元素类型是 Duck 接口，不是具体鸭子类型。
	// 后面无论加多少种鸭子，这段循环都不用改 —— 这就是"针对接口编程"。
	ducks := []strategy.Duck{
		strategy.NewMallardDuck(),
		strategy.NewRedheadDuck(),
		strategy.NewRubberDuck(),
		strategy.NewDecoyDuck(),
	}

	fmt.Println("=== 1. 各种鸭子的当前表现 ===")
	for _, d := range ducks {
		show(d)
	}

	fmt.Println("\n=== 2. 运行时动态换行为（策略模式的核心收益）===")
	model := strategy.NewModelDuck()
	show(model)
	fmt.Println(">> 给模型鸭装上火箭助推器……")
	model.SetFlyBehavior(strategy.FlyWithRocket{})
	show(model)

	fmt.Println("\n=== 3. 换叫声也一样，鸭子本身一行代码都不用改 ===")
	fmt.Println(">> 给绿头鸭换成诱饵鸭的静音叫声……")
	greenHead := strategy.NewMallardDuck()
	greenHead.SetQuackBehavior(strategy.MuteQuack{})
	show(greenHead)
}

func show(d strategy.Duck) {
	fmt.Printf("%s | %s | %s | %s\n", d.Display(), d.PerformFly(), d.PerformQuack(), d.Swim())
}
