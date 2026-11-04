package strategy

// ============================================================================
// 变化的部分：算法族（family of algorithms）
//
// 在《Head First 设计模式》里，这些是 FlyBehavior / QuackBehavior 两个接口
// 以及它们各自的实现类。Java 里写成 interface + implements，Go 里就是
// interface + 方法集隐式满足（见 README 的 Go 注意事项）。
//
// 关键点：这两个接口里都"只有一个方法"，且方法不返回 error、不依赖别的状态。
// 这种极窄接口在 Go 里很容易实现、很容易替换、也很容易在测试里造假。
// ============================================================================

// FlyBehavior 是"怎么飞"这一族算法的抽象。
type FlyBehavior interface {
	// Fly 返回一句自我描述（书上用 fmt.Println 直接打印，
	// 这里改成返回值，方便测试断言，demo 里再统一打印）。
	Fly() string
}

// FlyWithWings 用翅膀飞。大多数真鸭子是这种。
type FlyWithWings struct{}

func (FlyWithWings) Fly() string { return "我在用翅膀飞" }

// FlyNoWay 不会飞。橡皮鸭、诱饵鸭用它。
type FlyNoWay struct{}

func (FlyNoWay) Fly() string { return "我根本不会飞" }

// FlyWithRocket 火箭助推——书上用来演示"运行时动态改变行为"。
type FlyWithRocket struct{}

func (FlyWithRocket) Fly() string { return "我在用火箭助推器飞！" }

// QuackBehavior 是"怎么叫"这一族算法的抽象。
type QuackBehavior interface {
	Quack() string
}

// Quack 正常嘎嘎叫。
type Quack struct{}

func (Quack) Quack() string { return "嘎嘎嘎" }

// Squeak 橡皮鸭的吱吱声。
type Squeak struct{}

func (Squeak) Quack() string { return "吱吱吱" }

// MuteQuack 诱饵鸭：不出声。
type MuteQuack struct{}

func (MuteQuack) Quack() string { return "<<沉默>>" }
