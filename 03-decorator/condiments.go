package decorator

// ============================================================================
// 四个调料（装饰者）
//
// 每个调料只有三行长这样：
//
//	type Xxx struct{ Beverage }                                  // ① 嵌入"被装饰者"
//	func NewXxx(b Beverage) Xxx { return Xxx{Beverage: b} }      // ② 构造：包住它
//	func (x Xxx) Cost() float64 { return x.Beverage.Cost() + 价 } // ③ 加价，并转发
//
// 为什么嵌入的是 **接口** 而不是具体类型？这是本模式在 Go 里最关键的一行：
//
//	type Milk struct{ Beverage }   // ✅ 能包住任何 Beverage：Espresso、Decaf、甚至别的 Milk
//	type Milk struct{ Espresso }   // ❌ 只能包 Espresso，装饰器立刻失去意义
//
// 嵌入接口后，Milk 自动"继承"了被装饰者的两个方法；
// 我们再定义同名方法把它们遮蔽（shadow）成"多加一份钱、多一句描述"。
// 这就是 Go 版的"override"，也是 Go 里表达装饰器最紧凑的写法。
// ============================================================================

// Milk 蒸牛奶，+0.10。
type Milk struct{ Beverage }

func NewMilk(b Beverage) Milk { return Milk{Beverage: b} }

func (m Milk) Description() string { return m.Beverage.Description() + "，蒸牛奶" }
func (m Milk) Cost() float64       { return m.Beverage.Cost() + 0.10 }

// Mocha 摩卡，+0.20。可以加多次（双份摩卡），因为每加一次就多包一层。
type Mocha struct{ Beverage }

func NewMocha(b Beverage) Mocha { return Mocha{Beverage: b} }

func (m Mocha) Description() string { return m.Beverage.Description() + "，摩卡" }
func (m Mocha) Cost() float64       { return m.Beverage.Cost() + 0.20 }

// Soy 大豆，+0.15。
type Soy struct{ Beverage }

func NewSoy(b Beverage) Soy { return Soy{Beverage: b} }

func (s Soy) Description() string { return s.Beverage.Description() + "，大豆" }
func (s Soy) Cost() float64       { return s.Beverage.Cost() + 0.15 }

// Whip 奶泡，+0.10。
type Whip struct{ Beverage }

func NewWhip(b Beverage) Whip { return Whip{Beverage: b} }

func (w Whip) Description() string { return w.Beverage.Description() + "，奶泡" }
func (w Whip) Cost() float64       { return w.Beverage.Cost() + 0.10 }

// 编译期断言：装饰者自己也是 Beverage —— 所以可以继续被装饰。
// 这四行是"装饰者与被装饰者同类型"这一前提的可执行证明。
var (
	_ Beverage = Milk{}
	_ Beverage = Mocha{}
	_ Beverage = Soy{}
	_ Beverage = Whip{}
)
