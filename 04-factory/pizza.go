package factory

import "strings"

// ============================================================================
// 产品族：披萨
//
// 书上是抽象类 Pizza：name/dough/sauce 等字段 + prepare() 抽象 + bake/cut/box 具体。
// Go 里拆成：
//   - Pizza 接口：客户（PizzaStore）唯一认识的东西；
//   - 具体结构体：自己持有原料字段，自己实现 Prepare。
//
// 关键设计：具体披萨**不知道**原料来自哪家工厂，它只知道"我要面团、酱、奶酪"。
// 这正是依赖倒置原则（DIP）落到代码上的样子。
// ============================================================================

// Pizza 是产品合同。注意它只有两个方法、也不关心原料是谁给的。
type Pizza interface {
	Name() string
	// Prepare 返回准备过程。书上这里是若干条 println，这里改成返回文本，方便测试断言。
	Prepare() string
}

// ---------------------------------------------------------------------------
// CheesePizza 奶酪披萨
// ---------------------------------------------------------------------------

type CheesePizza struct {
	factory PizzaIngredientFactory // 字段类型是接口 —— 不认识任何具体工厂

	dough  Dough
	sauce  Sauce
	cheese Cheese
}

func NewCheesePizza(f PizzaIngredientFactory) *CheesePizza {
	return &CheesePizza{factory: f}
}

func (p *CheesePizza) Name() string { return "奶酪披萨" }

func (p *CheesePizza) Prepare() string {
	p.dough = p.factory.CreateDough()
	p.sauce = p.factory.CreateSauce()
	p.cheese = p.factory.CreateCheese()

	return strings.Join([]string{
		"揉面团：" + string(p.dough),
		"加酱料：" + string(p.sauce),
		"撒奶酪：" + string(p.cheese),
	}, "\n")
}

// ---------------------------------------------------------------------------
// ClamPizza 蛤蜊披萨
// ---------------------------------------------------------------------------

type ClamPizza struct {
	factory PizzaIngredientFactory

	dough  Dough
	sauce  Sauce
	cheese Cheese
	clams  Clams
}

func NewClamPizza(f PizzaIngredientFactory) *ClamPizza {
	return &ClamPizza{factory: f}
}

func (p *ClamPizza) Name() string { return "蛤蜊披萨" }

func (p *ClamPizza) Prepare() string {
	p.dough = p.factory.CreateDough()
	p.sauce = p.factory.CreateSauce()
	p.cheese = p.factory.CreateCheese()
	p.clams = p.factory.CreateClams()

	return strings.Join([]string{
		"揉面团：" + string(p.dough),
		"加酱料：" + string(p.sauce),
		"撒奶酪：" + string(p.cheese),
		"加蛤蜊：" + string(p.clams),
	}, "\n")
}

// ---------------------------------------------------------------------------
// VeggiePizza 蔬菜披萨
// ---------------------------------------------------------------------------

// veggieToppings 是"不随地域变化的原料"，所以没有放进抽象工厂。
// 判断标准：它不会随"纽约/芝加哥"这种风格维度变化 —— 不属于那个产品族。
var veggieToppings = []string{"茄子", "菠菜", "黑橄榄"}

type VeggiePizza struct {
	factory PizzaIngredientFactory

	dough  Dough
	sauce  Sauce
	cheese Cheese
}

func NewVeggiePizza(f PizzaIngredientFactory) *VeggiePizza {
	return &VeggiePizza{factory: f}
}

func (p *VeggiePizza) Name() string { return "蔬菜披萨" }

func (p *VeggiePizza) Prepare() string {
	p.dough = p.factory.CreateDough()
	p.sauce = p.factory.CreateSauce()
	p.cheese = p.factory.CreateCheese()

	return strings.Join([]string{
		"揉面团：" + string(p.dough),
		"加酱料：" + string(p.sauce),
		"撒奶酪：" + string(p.cheese),
		"加蔬菜：" + strings.Join(veggieToppings, "、"),
	}, "\n")
}

// 编译期断言：三种披萨都满足产品合同。
var (
	_ Pizza = (*CheesePizza)(nil)
	_ Pizza = (*ClamPizza)(nil)
	_ Pizza = (*VeggiePizza)(nil)
)
