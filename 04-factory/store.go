package factory

import "fmt"

// ============================================================================
// 工厂方法：把"创建哪种披萨"这件事交给一个**钩子**
//
// 书上结构：
//   abstract class PizzaStore {
//       Pizza orderPizza(String type) {          // 模板方法：流程固定
//           Pizza p = createPizza(type);         // ← 抽象方法，子类决定
//           p.prepare(); p.bake(); p.cut(); p.box();
//       }
//       abstract Pizza createPizza(String type); // ← 工厂方法
//   }
//   class NYPizzaStore extends PizzaStore { Pizza createPizza(...) {...} }
//
// Go 没有继承，所以"子类覆写抽象方法"用**函数字段**来表达：
//   PizzaStore.createPizza 就是一个钩子，具体店在构造函数里把它填上。
// 这比"每个店建一个结构体 + 一个方法"更短，而且同样满足"新增店不改旧代码"。
//
// （另一种写法是把钩子做成接口并嵌入结构体，见 README「两种 Go 写法」。）
// ============================================================================

// CreatePizzaFunc 是工厂方法的签名：给一个品类，返回一个产品。
type CreatePizzaFunc func(kind string) (Pizza, error)

// PizzaStore 既是"客户端"（使用产品），也是"工厂方法的持有者"（决定用哪个实现）。
//
// 注意 OrderPizza 里从头到尾**没有出现过任何具体披萨类型** —— 这就是 DIP。
type PizzaStore struct {
	name        string
	createPizza CreatePizzaFunc
}

// NewPizzaStore 允许在外部组装一家新店，而不需要修改本包。
// 这是"工厂方法"能对扩展开放的关键：钩子是注入的，不是写死的 switch。
func NewPizzaStore(name string, create CreatePizzaFunc) *PizzaStore {
	return &PizzaStore{name: name, createPizza: create}
}

// Order 是一张订单：谁做的、做了什么、全过程步骤。
type Order struct {
	Store string
	Pizza Pizza
	Steps []string
}

// OrderPizza 是模板方法：定好流程，只在第一行调用钩子。
// 「哪里变」和「哪里不变」被清清楚楚分开：变的是 createPizza，不变的是这套流程。
func (s *PizzaStore) OrderPizza(kind string) (Order, error) {
	p, err := s.createPizza(kind) // ← 唯一的扩展点
	if err != nil {
		return Order{}, fmt.Errorf("%s 点单失败: %w", s.name, err)
	}

	return Order{
		Store: s.name,
		Pizza: p,
		Steps: []string{
			"制作：" + p.Name(),
			p.Prepare(),
			"烘烤 25 分钟",
			"切成 8 块",
			"装入" + s.name + "的盒子",
		},
	}, nil
}

// ---------------------------------------------------------------------------
// 两家具体的店
//
// 每家店做两件事：选一个原料工厂（抽象工厂），并决定每个品类 new 哪种披萨。
// 这两件事正是"工厂方法"和"抽象工厂"叠在一起的地方。
// ---------------------------------------------------------------------------

// NewYorkPizzaStore 纽约店：薄饼风格。
func NewYorkPizzaStore() *PizzaStore {
	ingredients := NewYorkIngredientFactory{} // 选族：抽象工厂

	return NewPizzaStore("纽约披萨店", func(kind string) (Pizza, error) {
		switch kind { // 选品类：工厂方法
		case "cheese":
			return NewCheesePizza(ingredients), nil
		case "clam":
			return NewClamPizza(ingredients), nil
		case "veggie":
			return NewVeggiePizza(ingredients), nil
		}
		return nil, fmt.Errorf("没有 %q 这种披萨", kind)
	})
}

// ChicagoPizzaStore 芝加哥店：厚饼风格。
func ChicagoPizzaStore() *PizzaStore {
	ingredients := ChicagoIngredientFactory{}

	return NewPizzaStore("芝加哥披萨店", func(kind string) (Pizza, error) {
		switch kind {
		case "cheese":
			return NewCheesePizza(ingredients), nil
		case "clam":
			return NewClamPizza(ingredients), nil
		case "veggie":
			return NewVeggiePizza(ingredients), nil
		}
		return nil, fmt.Errorf("没有 %q 这种披萨", kind)
	})
}

// ---------------------------------------------------------------------------
// 书里第 4 章开头的"简单工厂"
// ---------------------------------------------------------------------------

// NewSimplePizza 是"简单工厂"——书里明确说这**不是**一个 GoF 模式，
// 只是程序员们习以为常的习惯。它把 new 收拢到一个地方，但：
//
//   - 新增品类必须改这个 switch（对修改开放了）；
//   - 所有店的差异只能靠参数传进来，没法让"纽约店"和"芝加哥店"各自演化。
//
// 它和工厂方法的区别不是"好不好"，而是"谁来承担变化"：
// 简单工厂把变化留在函数里，工厂方法把变化推给新的实现者。
func NewSimplePizza(f PizzaIngredientFactory, kind string) (Pizza, error) {
	switch kind {
	case "cheese":
		return NewCheesePizza(f), nil
	case "clam":
		return NewClamPizza(f), nil
	case "veggie":
		return NewVeggiePizza(f), nil
	}
	return nil, fmt.Errorf("没有 %q 这种披萨", kind)
}
