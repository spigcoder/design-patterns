package decorator

// ============================================================================
// 装饰器模式的唯一前提：装饰者与被装饰者 **是同一个类型**。
//
// 书上是抽象类 Beverage，两个方法：getDescription()（部分实现）和 cost()（抽象）。
// Go 里没有抽象类，所以：
//   - Beverage 变成接口（合同：能被描述、能算钱）；
//   - "部分实现"这件事没有了 —— 每个具体饮料自己给描述、自己算价。
//     书上 getDescription 提供默认值、子类覆盖，这种"默认值"在 Go 里用嵌入结构体做，
//     但这里只有两三个字段，直接各写各的更清楚（不要为了对称而抽象）。
// ============================================================================

// Beverage 是"一杯能点的饮料"，也是"一个能包住饮料的调料"。
// 正是这种'同类型'让装饰者可以层层嵌套、无限组合。
type Beverage interface {
	Description() string
	Cost() float64 // 单位：元
}

// 以下是书上的四种基础饮料（Starbuzz 的菜单）。
// 都是空结构体：它们没有状态，只有身份。零值即合法值，不需要构造函数。
type (
	// Espresso 浓缩咖啡。
	Espresso struct{}
	// HouseBlend 综合咖啡。
	HouseBlend struct{}
	// DarkRoast 深焙咖啡。
	DarkRoast struct{}
	// Decaf 低因咖啡。
	Decaf struct{}
)

func (Espresso) Description() string   { return "浓缩咖啡" }
func (HouseBlend) Description() string { return "综合咖啡" }
func (DarkRoast) Description() string  { return "深焙咖啡" }
func (Decaf) Description() string      { return "低因咖啡" }

func (Espresso) Cost() float64   { return 1.99 }
func (HouseBlend) Cost() float64 { return 0.89 }
func (DarkRoast) Cost() float64  { return 0.99 }
func (Decaf) Cost() float64      { return 1.05 }

// 编译期断言：四种基础饮料都满足 Beverage。
var (
	_ Beverage = Espresso{}
	_ Beverage = HouseBlend{}
	_ Beverage = DarkRoast{}
	_ Beverage = Decaf{}
)
