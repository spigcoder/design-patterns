package decorator_test

import (
	"fmt"
	"math"
	"strings"
	"testing"

	decorator "github.com/spigcoder/design-patterns/03-decorator"
)

// approx 比较金额。用 float64 表示钱是跟着书走（书上用 double），
// 真实项目该用整数分或定点小数——但那属于工程细节，本目录不展开。
func approx(t *testing.T, got, want float64, label string) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s: got %.2f, want %.2f", label, got, want)
	}
}

// 场景 1：纯饮料，没有任何装饰。
func TestPlainBeverage(t *testing.T) {
	var b decorator.Beverage = decorator.Espresso{}

	if got := b.Description(); got != "浓缩咖啡" {
		t.Errorf("Description() = %q", got)
	}
	approx(t, b.Cost(), 1.99, "Espresso.Cost()")
}

// 场景 2：一层层包上去，价格和描述都累加。
func TestDecoratorsAccumulate(t *testing.T) {
	// 双份摩卡 + 奶泡的深焙咖啡：0.99 + 0.20 + 0.20 + 0.10 = 1.49
	b := decorator.NewWhip(decorator.NewMocha(decorator.NewMocha(decorator.DarkRoast{})))

	approx(t, b.Cost(), 1.49, "双份摩卡 + 奶泡")
	if got := b.Description(); got != "深焙咖啡，摩卡，摩卡，奶泡" {
		t.Errorf("Description() = %q", got)
	}
}

// 场景 3：书上的 beverage3——综合咖啡 + 大豆 + 摩卡 + 奶泡 = 1.34。
func TestBookExample(t *testing.T) {
	b := decorator.NewWhip(decorator.NewMocha(decorator.NewSoy(decorator.HouseBlend{})))

	approx(t, b.Cost(), 1.34, "0.89 + 0.15 + 0.20 + 0.10")
}

// 场景 4：包装顺序不影响总价，但影响描述（价格可交换，文本不可交换）。
func TestOrderOfWrappingDoesNotChangeCost(t *testing.T) {
	milkThenMocha := decorator.NewMocha(decorator.NewMilk(decorator.Espresso{}))
	mochaThenMilk := decorator.NewMilk(decorator.NewMocha(decorator.Espresso{}))

	approx(t, milkThenMocha.Cost(), mochaThenMilk.Cost(), "两种包装顺序")
	approx(t, milkThenMocha.Cost(), 2.29, "1.99 + 0.10 + 0.20")

	if milkThenMocha.Description() == mochaThenMilk.Description() {
		t.Errorf("描述按包装顺序拼接，两者应当不同")
	}
}

// 场景 5：装饰者本身也是 Beverage。
// 任何只认 Beverage 的函数都能接收装饰过的饮料，不需要知道它被包了几层。
func TestDecoratedIsStillABeverage(t *testing.T) {
	b := decorator.NewWhip(decorator.NewSoy(decorator.Decaf{}))

	if got := describe(b); !strings.Contains(got, "奶泡") {
		t.Errorf("装饰过的饮料没有正确传给只认 Beverage 的函数: %q", got)
	}
}

// describe 是一个"只认识 Beverage"的普通函数 —— 它不知道也不需要知道装饰链。
func describe(b decorator.Beverage) string {
	return fmt.Sprintf("%s（%.2f 元）", b.Description(), b.Cost())
}

// 场景 6：4 种调料 → 16 种组合，全部用同一批类型拼出来。
// 如果按"每种组合建一个类"的老办法，这里需要 16 个类型；
// 用装饰器只需要 4 个调料类型 + 循环。
func TestAllCombinationsNeedNoNewTypes(t *testing.T) {
	type condiment struct {
		name  string
		price float64
		wrap  func(decorator.Beverage) decorator.Beverage
	}
	all := []condiment{
		{"牛奶", 0.10, func(b decorator.Beverage) decorator.Beverage { return decorator.NewMilk(b) }},
		{"摩卡", 0.20, func(b decorator.Beverage) decorator.Beverage { return decorator.NewMocha(b) }},
		{"大豆", 0.15, func(b decorator.Beverage) decorator.Beverage { return decorator.NewSoy(b) }},
		{"奶泡", 0.10, func(b decorator.Beverage) decorator.Beverage { return decorator.NewWhip(b) }},
	}

	const basePrice = 0.89 // HouseBlend
	for mask := 0; mask < 1<<len(all); mask++ {
		var b decorator.Beverage = decorator.HouseBlend{}
		want := basePrice
		var names []string

		for i, c := range all {
			if mask&(1<<i) == 0 {
				continue
			}
			b = c.wrap(b)
			want += c.price
			names = append(names, c.name)
		}

		approx(t, b.Cost(), want, "组合 "+strings.Join(names, "+"))
	}
}

// 场景 7：同一个基础饮料可以按顾客选择在运行时被包装 —— 这正是装饰器相对继承的收益。
func TestBuildAtRuntime(t *testing.T) {
	order := func(base decorator.Beverage, extras ...string) decorator.Beverage {
		for _, e := range extras {
			switch e {
			case "牛奶":
				base = decorator.NewMilk(base)
			case "摩卡":
				base = decorator.NewMocha(base)
			case "奶泡":
				base = decorator.NewWhip(base)
			}
		}
		return base
	}

	b := order(decorator.Espresso{}, "摩卡", "摩卡", "奶泡")
	approx(t, b.Cost(), 2.49, "1.99 + 0.20*2 + 0.10")

	if got := b.Description(); got != "浓缩咖啡，摩卡，摩卡，奶泡" {
		t.Errorf("Description() = %q", got)
	}
}
