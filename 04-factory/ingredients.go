package factory

// ============================================================================
// 抽象工厂：一族相关产品的"生产接口"
//
// 书上是：Dough / Sauce / Cheese / Clams 四个接口 + 四个实现类（ThinCrustDough、
// MarinaraSauce、ReggianoCheese、FreshClams ...），一共 8 个原料类 + 2 个工厂类。
//
// Go 的取舍：原料只需要"能被区分、能被打印"，不需要行为，
// 所以用**带名字的字符串类型 + 常量**，而不是为每种原料建一个类。
// 这是本章最值得记的一条 Go 经验：不要为了对称去建空结构体。
//
// 用带名字的类型（而不是裸 string）仍然保留了类型安全：
// CreateDough 返回 Dough 而不是 string，把 Sauce 传进 CreateDough 会编译失败。
// ============================================================================

type (
	// Dough 面团。
	Dough string
	// Sauce 酱料。
	Sauce string
	// Cheese 奶酪。
	Cheese string
	// Clams 蛤蜊。
	Clams string
)

const (
	// 纽约族
	ThinCrustDough Dough  = "薄饼面团"
	MarinaraSauce  Sauce  = "Marinara 酱"
	ReggianoCheese Cheese = "Reggiano 奶酪"
	FreshClams     Clams  = "新鲜蛤蜊"

	// 芝加哥族
	ThickCrustDough  Dough  = "厚饼面团"
	PlumTomatoSauce  Sauce  = "李子番茄酱"
	MozzarellaCheese Cheese = "Mozzarella 奶酪"
	FrozenClams      Clams  = "冷冻蛤蜊"
)

// PizzaIngredientFactory 就是抽象工厂。
//
// 它和"工厂方法"的区别看这个接口的形状就能看出来：
//   - 工厂方法：一个方法，产出**一种**产品，具体实现由子类/钩子决定；
//   - 抽象工厂：一组方法，产出**一族**产品，具体实现由一个工厂对象决定。
//
// 它保证的是一致性：只要拿的是纽约工厂，四个方法返回的原料一定都来自纽约族，
// 披萨类完全不用操心"别把芝加哥的厚饼配纽约的 Marinara 酱"。
type PizzaIngredientFactory interface {
	CreateDough() Dough
	CreateSauce() Sauce
	CreateCheese() Cheese
	CreateClams() Clams
}

// NewYorkIngredientFactory 纽约原料工厂：薄饼 + Marinara 酱 + Reggiano 奶酪 + 新鲜蛤蜊。
type NewYorkIngredientFactory struct{}

func (NewYorkIngredientFactory) CreateDough() Dough   { return ThinCrustDough }
func (NewYorkIngredientFactory) CreateSauce() Sauce   { return MarinaraSauce }
func (NewYorkIngredientFactory) CreateCheese() Cheese { return ReggianoCheese }
func (NewYorkIngredientFactory) CreateClams() Clams   { return FreshClams }

// ChicagoIngredientFactory 芝加哥原料工厂：厚饼 + 李子番茄酱 + Mozzarella 奶酪 + 冷冻蛤蜊。
type ChicagoIngredientFactory struct{}

func (ChicagoIngredientFactory) CreateDough() Dough   { return ThickCrustDough }
func (ChicagoIngredientFactory) CreateSauce() Sauce   { return PlumTomatoSauce }
func (ChicagoIngredientFactory) CreateCheese() Cheese { return MozzarellaCheese }
func (ChicagoIngredientFactory) CreateClams() Clams   { return FrozenClams }

var (
	_ PizzaIngredientFactory = NewYorkIngredientFactory{}
	_ PizzaIngredientFactory = ChicagoIngredientFactory{}
)
