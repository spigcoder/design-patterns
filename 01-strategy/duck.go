package strategy

// ============================================================================
// 不变的部分：Duck 骨架 + 客户端依赖的抽象
//
// 书上：Duck 是抽象类，持有两个行为接口字段，performFly()/performQuack()
// 把它们"委托"出去，display() 是留给子类的抽象方法。
//
// Go 没有继承也没有 abstract，所以拆成两块：
//   - Duck       接口：调用方（客户端）只认这个，这就是"面向接口"的那一侧；
//   - BaseDuck   结构体：放所有鸭子共有的状态与委托逻辑，供具体鸭子"嵌入"复用。
// 具体鸭子嵌入 *BaseDuck 后会自动获得 PerformFly 等方法（方法提升），
// 再自己实现 Display 来"覆盖"父类的抽象方法。
// ============================================================================

// Duck 是客户端唯一需要认识的类型：它只承诺"能被展示 / 会飞 / 会叫 / 能换行为"。
// 客户端代码拿到的如果是 Duck，就永远不知道也不关心背后是哪种鸭子。
type Duck interface {
	Display() string
	PerformFly() string
	PerformQuack() string
	Swim() string

	// 这两个 setter 是策略模式的灵魂：行为可以在运行时被替换。
	SetFlyBehavior(FlyBehavior)
	SetQuackBehavior(QuackBehavior)
}

// BaseDuck 持有两个行为"插槽"。
// 注意这里的字段是接口类型——这正是"针对接口编程"落到代码上的样子：
// BaseDuck 从不知道翅膀/火箭/橡皮的具体类型，只知道它们满足 FlyBehavior。
type BaseDuck struct {
	flyBehavior   FlyBehavior
	quackBehavior QuackBehavior
}

// NewBaseDuck 是行为的注入点（injection point）。
// 具体实现只在这里出现一次，之后全程只通过接口打交道。
func NewBaseDuck(fly FlyBehavior, quack QuackBehavior) *BaseDuck {
	return &BaseDuck{flyBehavior: fly, quackBehavior: quack}
}

// SetFlyBehavior 运行时换飞行方式。
func (d *BaseDuck) SetFlyBehavior(fly FlyBehavior) { d.flyBehavior = fly }

// SetQuackBehavior 运行时换叫声。
func (d *BaseDuck) SetQuackBehavior(quack QuackBehavior) { d.quackBehavior = quack }

// PerformFly 把"飞"这件事委托给 flyBehavior，而不是自己实现。
// 书上说：用组合（composition）代替继承（inheritance）——这里就是那一行。
func (d *BaseDuck) PerformFly() string { return d.flyBehavior.Fly() }

// PerformQuack 同理，委托给 quackBehavior。
func (d *BaseDuck) PerformQuack() string { return d.quackBehavior.Quack() }

// Swim 所有鸭子都会浮水，行为一致，不需要抽成策略。
// （判断标准：会不会变。不会变的东西不要抽象，见 README。）
func (d *BaseDuck) Swim() string { return "我在浮水（所有鸭子都会）" }
