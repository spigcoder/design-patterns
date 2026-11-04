package strategy_test

import (
	"testing"

	strategy "github.com/spigcoder/design-patterns/01-strategy"
)

// 场景 1：每种鸭子把"飞"和"叫"都委托给正确的行为对象。
func TestDucksDelegateToBehaviors(t *testing.T) {
	cases := []struct {
		name     string
		duck     strategy.Duck
		wantFly  string
		wantQuak string
	}{
		{"绿头鸭", strategy.NewMallardDuck(), "我在用翅膀飞", "嘎嘎嘎"},
		{"红头鸭", strategy.NewRedheadDuck(), "我在用翅膀飞", "嘎嘎嘎"},
		{"橡皮鸭", strategy.NewRubberDuck(), "我根本不会飞", "吱吱吱"},
		{"诱饵鸭", strategy.NewDecoyDuck(), "我根本不会飞", "<<沉默>>"},
		{"模型鸭", strategy.NewModelDuck(), "我根本不会飞", "嘎嘎嘎"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.duck.PerformFly(); got != c.wantFly {
				t.Errorf("PerformFly() = %q, want %q", got, c.wantFly)
			}
			if got := c.duck.PerformQuack(); got != c.wantQuak {
				t.Errorf("PerformQuack() = %q, want %q", got, c.wantQuak)
			}
		})
	}
}

// 场景 2：运行时替换行为，不需要新建子类、不需要改鸭子代码。
func TestBehaviorIsSwappableAtRuntime(t *testing.T) {
	duck := strategy.NewModelDuck()
	if got := duck.PerformFly(); got != "我根本不会飞" {
		t.Fatalf("初始 PerformFly() = %q", got)
	}

	duck.SetFlyBehavior(strategy.FlyWithRocket{})
	if got := duck.PerformFly(); got != "我在用火箭助推器飞！" {
		t.Fatalf("换行为后 PerformFly() = %q", got)
	}

	duck.SetQuackBehavior(strategy.MuteQuack{})
	if got := duck.PerformQuack(); got != "<<沉默>>" {
		t.Fatalf("换叫声后 PerformQuack() = %q", got)
	}
}

// 场景 3：加一族新算法，不用动任何已有代码。
// 这里在测试文件里"外部"实现一个 FlyBehavior 就够了 ——
// 接口是隐式满足的，连 import 一个注册表都不用。
type flyWithJetpack struct{}

func (flyWithJetpack) Fly() string { return "我在用喷气背包飞" }

func TestNewBehaviorNeedsNoChangesToExistingCode(t *testing.T) {
	duck := strategy.NewMallardDuck()
	duck.SetFlyBehavior(flyWithJetpack{})

	if got := duck.PerformFly(); got != "我在用喷气背包飞" {
		t.Fatalf("PerformFly() = %q", got)
	}
}

// 场景 4：客户端只依赖 Duck 接口，就能统一处理所有鸭子。
func TestClientOnlyNeedsTheInterface(t *testing.T) {
	all := []strategy.Duck{
		strategy.NewMallardDuck(),
		strategy.NewRedheadDuck(),
		strategy.NewRubberDuck(),
		strategy.NewDecoyDuck(),
		strategy.NewModelDuck(),
	}
	for _, d := range all {
		if d.Display() == "" {
			t.Errorf("%T 的 Display() 为空", d)
		}
		if d.PerformFly() == "" || d.PerformQuack() == "" {
			t.Errorf("%T 的行为委托返回了空字符串", d)
		}
	}
}
