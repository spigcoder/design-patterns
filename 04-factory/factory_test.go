package factory_test

import (
	"strings"
	"testing"

	factory "github.com/spigcoder/design-patterns/04-factory"
)

// 场景 1：同样的披萨，换一个原料工厂 → 原料整族变化，披萨类一行不改。
// 这就是抽象工厂的核心价值：保证"族"的一致性。
func TestSamePizzaDifferentIngredientFamily(t *testing.T) {
	ny := factory.NewCheesePizza(factory.NewYorkIngredientFactory{})
	chicago := factory.NewCheesePizza(factory.ChicagoIngredientFactory{})

	if got := ny.Prepare(); !strings.Contains(got, "薄饼面团") || !strings.Contains(got, "Marinara 酱") {
		t.Errorf("纽约奶酪披萨原料不对:\n%s", got)
	}
	if got := chicago.Prepare(); !strings.Contains(got, "厚饼面团") || !strings.Contains(got, "李子番茄酱") {
		t.Errorf("芝加哥奶酪披萨原料不对:\n%s", got)
	}
}

// 场景 2：两家店点同一款披萨，出来的东西不同 —— 差异全部来自"店 + 工厂"的组合。
func TestStoresProduceDifferentStyles(t *testing.T) {
	nyOrder, err := factory.NewYorkPizzaStore().OrderPizza("cheese")
	if err != nil {
		t.Fatalf("纽约店点单失败: %v", err)
	}
	chicagoOrder, err := factory.ChicagoPizzaStore().OrderPizza("cheese")
	if err != nil {
		t.Fatalf("芝加哥店点单失败: %v", err)
	}

	if nyOrder.Store != "纽约披萨店" || chicagoOrder.Store != "芝加哥披萨店" {
		t.Errorf("店名不对: %q / %q", nyOrder.Store, chicagoOrder.Store)
	}
	if nyOrder.Pizza.Name() != chicagoOrder.Pizza.Name() {
		t.Errorf("同一款披萨，名字应该相同: %q vs %q", nyOrder.Pizza.Name(), chicagoOrder.Pizza.Name())
	}

	nySteps := strings.Join(nyOrder.Steps, "\n")
	chicagoSteps := strings.Join(chicagoOrder.Steps, "\n")
	if !strings.Contains(nySteps, "薄饼面团") {
		t.Errorf("纽约店没用薄饼:\n%s", nySteps)
	}
	if !strings.Contains(chicagoSteps, "厚饼面团") {
		t.Errorf("芝加哥店没用厚饼:\n%s", chicagoSteps)
	}
}

// 场景 3：模板方法把那套固定流程跑全了（烤、切、装盒）。
func TestOrderRunsTheFixedSteps(t *testing.T) {
	order, err := factory.ChicagoPizzaStore().OrderPizza("clam")
	if err != nil {
		t.Fatalf("点单失败: %v", err)
	}

	all := strings.Join(order.Steps, "\n")
	for _, want := range []string{"制作：蛤蜊披萨", "加蛤蜊：冷冻蛤蜊", "烘烤 25 分钟", "切成 8 块", "装入芝加哥披萨店的盒子"} {
		if !strings.Contains(all, want) {
			t.Errorf("流程缺少 %q:\n%s", want, all)
		}
	}
}

// 场景 4：不认识的品类 → 返回错误（Go 用 error 代替 Java 的异常）。
func TestUnknownKindReturnsError(t *testing.T) {
	if _, err := factory.NewYorkPizzaStore().OrderPizza("durian"); err == nil {
		t.Fatal("点了不存在的披萨，居然没报错")
	}
}

// 场景 5：新增一座城市的店 —— 不用改已有代码。
// 下面的东京店和东京原料工厂都定义在测试包里，
// factory 包里一个字都没动，这就是"对扩展开放、对修改关闭"。
func TestAddingANewCityNeedsNoChange(t *testing.T) {
	store := factory.NewPizzaStore("东京披萨店", func(kind string) (factory.Pizza, error) {
		ing := tokyoIngredients{}
		switch kind {
		case "cheese":
			return factory.NewCheesePizza(ing), nil
		case "clam":
			return factory.NewClamPizza(ing), nil
		}
		return nil, errUnknown
	})

	order, err := store.OrderPizza("clam")
	if err != nil {
		t.Fatalf("东京店点单失败: %v", err)
	}

	all := strings.Join(order.Steps, "\n")
	if !strings.Contains(all, "米饼面团") || !strings.Contains(all, "芥末酱") {
		t.Errorf("东京店原料不对:\n%s", all)
	}
}

type tokyoIngredients struct{}

func (tokyoIngredients) CreateDough() factory.Dough   { return "米饼面团" }
func (tokyoIngredients) CreateSauce() factory.Sauce   { return "芥末酱" }
func (tokyoIngredients) CreateCheese() factory.Cheese { return "白奶酪" }
func (tokyoIngredients) CreateClams() factory.Clams   { return "北海道蛤蜊" }

// 场景 6：依赖倒置的实证 —— 披萨只认工厂接口，连"这是哪家的工厂"都不知道。
// 测试里给一个只在测试里存在的工厂，披萨照样能造出来。
func TestPizzaKnowsOnlyTheInterface(t *testing.T) {
	p := factory.NewCheesePizza(tokyoIngredients{})

	if got := p.Prepare(); !strings.Contains(got, "白奶酪") {
		t.Errorf("披萨没能用上注入的工厂:\n%s", got)
	}
}

// 场景 7：简单工厂也能干活，但它的扩展点在函数内部（对比场景 5）。
func TestSimpleFactoryStillWorks(t *testing.T) {
	p, err := factory.NewSimplePizza(factory.NewYorkIngredientFactory{}, "veggie")
	if err != nil {
		t.Fatalf("简单工厂失败: %v", err)
	}
	if got := p.Prepare(); !strings.Contains(got, "黑橄榄") {
		t.Errorf("蔬菜披萨内容不对:\n%s", got)
	}
}

var errUnknown = &unknownKindError{}

type unknownKindError struct{}

func (*unknownKindError) Error() string { return "没有这种披萨" }
