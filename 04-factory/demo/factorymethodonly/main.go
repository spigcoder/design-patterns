// Demo：如果只有"工厂方法"、没有"抽象工厂"，代码会长成什么样？
//
//	go run ./04-factory/demo/factorymethodonly
//
// 原料族（纽约 / 芝加哥）有 4 个成员：面团、酱、奶酪、蛤蜊。
// 没有"抽象工厂"这个可替换的族对象，你就只有两条路能表达"纽约风格"：
//
//	路 A：每种 "品类 × 风格" 建一个类型  → 3 品类 × 2 风格 = 6 个类型（本文件）
//	路 B：保持 3 个类型，但在每个类型内部 switch 城市 → 每个类型里都长出一个 switch
//
// 本文件走的是路 A：原料字符串**硬编码在类型内部**。
// 看看它有多重，再看 04-factory 里用抽象工厂的版本有多轻。
package main

import "fmt"

// ---------------------------------------------------------------------------
// 路 A：6 个类型，每个都把自己的原料抄了一遍
// ---------------------------------------------------------------------------

type nyCheesePizza struct{}

func (nyCheesePizza) Name() string { return "纽约奶酪披萨" }
func (nyCheesePizza) Prepare() string {
	return "揉面团：薄饼面团\n加酱料：Marinara 酱\n撒奶酪：Reggiano 奶酪"
}

type nyClamPizza struct{}

func (nyClamPizza) Name() string { return "纽约蛤蜊披萨" }
func (nyClamPizza) Prepare() string {
	return "揉面团：薄饼面团\n加酱料：Marinara 酱\n撒奶酪：Reggiano 奶酪\n加蛤蜊：新鲜蛤蜊"
}

type nyVeggiePizza struct{}

func (nyVeggiePizza) Name() string { return "纽约蔬菜披萨" }
func (nyVeggiePizza) Prepare() string {
	return "揉面团：薄饼面团\n加酱料：Marinara 酱\n撒奶酪：Reggiano 奶酪\n加蔬菜：茄子、菠菜、黑橄榄"
}

type chicagoCheesePizza struct{}

func (chicagoCheesePizza) Name() string { return "芝加哥奶酪披萨" }
func (chicagoCheesePizza) Prepare() string {
	return "揉面团：厚饼面团\n加酱料：李子番茄酱\n撒奶酪：Mozzarella 奶酪"
}

type chicagoClamPizza struct{}

func (chicagoClamPizza) Name() string { return "芝加哥蛤蜊披萨" }
func (chicagoClamPizza) Prepare() string {
	return "揉面团：厚饼面团\n加酱料：李子番茄酱\n撒奶酪：Mozzarella 奶酪\n加蛤蜊：冷冻蛤蜊"
}

type chicagoVeggiePizza struct{}

func (chicagoVeggiePizza) Name() string { return "芝加哥蔬菜披萨" }
func (chicagoVeggiePizza) Prepare() string {
	return "揉面团：厚饼面团\n加酱料：李子番茄酱\n撒奶酪：Mozzarella 奶酪\n加蔬菜：茄子、菠菜、黑橄榄"
}

// ---------------------------------------------------------------------------

func main() {
	all := []interface {
		Name() string
		Prepare() string
	}{
		nyCheesePizza{}, nyClamPizza{}, nyVeggiePizza{},
		chicagoCheesePizza{}, chicagoClamPizza{}, chicagoVeggiePizza{},
	}

	fmt.Println("=== 只用工厂方法：每个「品类 × 风格」一个类型，共 6 个 ===")
	for _, p := range all {
		fmt.Printf("■ %s\n%s\n", p.Name(), p.Prepare())
	}

	fmt.Println("=== 数一数代价 ===")
	fmt.Println("  3 品类 × 2 风格 = 6 个类型，每个类型里都抄了一份原料清单")
	fmt.Println("  加一种风格（东京）  → 再抄 3 个类型、3 段原料字符串")
	fmt.Println("  加一种品类（香肠）  → 每种风格各抄 1 个，2 个新类型")
	fmt.Println("  改原料（Marinara 酱换配方）→ 纽约的 3 个类型都要改")
	fmt.Println()
	fmt.Println("  换成抽象工厂（见 04-factory 包）后：")
	fmt.Println("  3 个品类 × 1 份原料族 = 3 个披萨类型（原料从工厂参数来，不再抄）")
	fmt.Println("  加一种风格 → 新增 1 个原料工厂类型，披萨类型一个都不动")
	fmt.Println("  改原料     → 只改那 1 个原料工厂，披萨类型一个都不动")
}
