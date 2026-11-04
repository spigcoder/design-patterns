package strategy

// ============================================================================
// 具体鸭子：只在构造函数里"选"行为，其余全部复用 BaseDuck
//
// 每只鸭子只多了两样东西：一个构造它的函数 + 一句 Display()。
// 书上用继承 + 方法覆盖做到的事，这里用"嵌入 + 同名方法遮蔽"做到。
// ============================================================================

// MallardDuck 绿头鸭：正常飞、正常叫。
type MallardDuck struct{ *BaseDuck }

func NewMallardDuck() *MallardDuck {
	return &MallardDuck{NewBaseDuck(FlyWithWings{}, Quack{})}
}

func (d *MallardDuck) Display() string { return "我是一只绿头鸭" }

// RedheadDuck 红头鸭：飞得像样，叫得也正常。
type RedheadDuck struct{ *BaseDuck }

func NewRedheadDuck() *RedheadDuck {
	return &RedheadDuck{NewBaseDuck(FlyWithWings{}, Quack{})}
}

func (d *RedheadDuck) Display() string { return "我是一只红头鸭" }

// RubberDuck 橡皮鸭：不会飞，只会吱吱叫。
// 书里正是它暴露了"把 fly() 放进父类"的坑。
type RubberDuck struct{ *BaseDuck }

func NewRubberDuck() *RubberDuck {
	return &RubberDuck{NewBaseDuck(FlyNoWay{}, Squeak{})}
}

func (d *RubberDuck) Display() string { return "我是一只橡皮鸭" }

// DecoyDuck 诱饵鸭：既不会飞，也不出声。
type DecoyDuck struct{ *BaseDuck }

func NewDecoyDuck() *DecoyDuck {
	return &DecoyDuck{NewBaseDuck(FlyNoWay{}, MuteQuack{})}
}

func (d *DecoyDuck) Display() string { return "我是一只木头诱饵鸭" }

// ModelDuck 模型鸭：一开始不会飞，等装上火箭才飞得起来（书上第 1 章结尾的演示）。
type ModelDuck struct{ *BaseDuck }

func NewModelDuck() *ModelDuck {
	return &ModelDuck{NewBaseDuck(FlyNoWay{}, Quack{})}
}

func (d *ModelDuck) Display() string { return "我是一只塑料模型鸭" }

// 编译期断言：这几种具体类型确实都满足 Duck 接口。
// Go 的接口是隐式满足的，这几行不会报错就说明"合同已履行"，
// 比等到运行时才发现漏实现友好得多。
var (
	_ Duck = (*MallardDuck)(nil)
	_ Duck = (*RedheadDuck)(nil)
	_ Duck = (*RubberDuck)(nil)
	_ Duck = (*DecoyDuck)(nil)
	_ Duck = (*ModelDuck)(nil)
)
