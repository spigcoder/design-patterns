# 03 · Decorator 装饰者模式

## 一句话意图

> 动态地将责任附加到对象上。若要扩展功能，装饰者提供了比继承更有弹性的替代方案。

（GoF："Attach additional responsibilities to an object dynamically. Decorators provide a flexible
alternative to subclassing for extending functionality."）

## 它解决的问题（书里的故事）

Starbuzz 咖啡要按"基础饮料 + 各种调料"算钱：4 种咖啡、4 种调料（牛奶/摩卡/大豆/奶泡），
而且**调料可以重复加、顺序随意**。

用继承硬做，会发生组合爆炸：

```
HouseBlendWithMilkAndMocha
HouseBlendWithMochaAndWhip
HouseBlendWithDoubleMochaAndSoy...
```

算一下：m 种调料可以选或不选 → 2^m 种组合；再乘上 n 种基础饮料 → **n × 2^m 个类**。
本目录的参数（n=4, m=4）就是 **64 个类型**。而且改动一处（比如摩卡涨价）要改一大批类。

真正的病根是：**"加调料"不是一种新的饮料，而是给已有饮料*加上一层责任*。**
继承表达的是"是一种（is-a）"，而这里的关系是"包着一个（wraps-a）"。所以书上给的结论是：

> 类应该对扩展开放，对修改关闭。（**开闭原则 OCP**）

加调料时，`HouseBlend`、`Espresso` 这些已有类型**一行都不用改**——扩展靠新增类型完成。

## 参与者 → Go 的对应物

| 书上（Java） | 本目录的 Go 写法 | 说明 |
|---|---|---|
| `Beverage` 抽象类（`getDescription()` 具体 + `cost()` 抽象） | `Beverage` 接口 | Go 无抽象类；"默认实现"这里用不上，各写各的更短 |
| `Espresso` / `HouseBlend` / `DarkRoast` / `Decaf extends Beverage` | 四个空结构体 | 无状态 → 空结构体，零值即合法值 |
| `CondimentDecorator extends Beverage`（抽象，持有 `Beverage beverage`） | **不需要**这个中间抽象类 | Go 里每个装饰者直接 `struct{ Beverage }`，一行就得到"包着一个" |
| `Mocha extends CondimentDecorator` | `type Mocha struct{ Beverage }` | 嵌入**接口**，不是嵌入具体类型（关键！） |
| `getDescription()` 里 `beverage.getDescription() + ", Mocha"` | `m.Beverage.Description() + "，摩卡"` | 方法遮蔽（shadow）= Go 版 override |
| `cost()` 返回 `beverage.cost() + .20` | `m.Beverage.Cost() + 0.20` | 转发 + 加价 |
| 客户端 `new Mocha(new Espresso())` | `NewMocha(Espresso{})` | 构造器返回具体类型（Go 惯例） |
| `Java I/O` 的 `new BufferedInputStream(new FileInputStream(...))` | `bufio.NewReader(gzip.NewReader(f))`、`http.Handler` 中间件 | 同一个模式，名字叫 wrapper / middleware |

## 结构图

```
              ┌──────────────────────┐
              │   «interface»        │   Description() string
              │     Beverage         │   Cost() float64
              └───┬──────────────┬───┘
        implements│              │implements（装饰者也实现它！）
   ┌──────────┬───┴────┬────────┐│
 Espresso HouseBlend DarkRoast Decaf
                              │
              ┌───────────────┴────────────────────────┐
              │ Milk / Mocha / Soy / Whip              │
              │ type Milk struct{ Beverage }           │  ← 内嵌一个 Beverage
              └────────────────┬───────────────────────┘
                               │ 委托 + 加价
                               ▼
                       被装饰者（可以是任意 Beverage，
                       包括另一个装饰者 → 无限嵌套）
```

## 可运行演示

本目录有两个 demo：

```bash
go test ./03-decorator             # 7 个场景测试
go run ./03-decorator/demo         # Starbuzz 点单
go run ./03-decorator/demo/stdlib  # 同一个模式在 Go 标准库里的样子
```

### Starbuzz

```
=== 1. 菜单上的基础饮料 ===
  浓缩咖啡 —— 1.99 元
  综合咖啡 —— 0.89 元
  深焙咖啡 —— 0.99 元
  低因咖啡 —— 1.05 元

=== 2. 按顾客要求逐层加料 ===
  浓缩咖啡，摩卡 —— 2.19 元
  深焙咖啡，摩卡，摩卡，奶泡 —— 1.49 元
  综合咖啡，大豆，摩卡，奶泡 —— 1.34 元
  低因咖啡，蒸牛奶，蒸牛奶 —— 1.25 元
—— 本单合计 6.27 元 ——

=== 3. 同一杯基础饮料，按顾客勾选在运行时拼装 ===
  综合咖啡 —— 0.89 元
  综合咖啡，摩卡 —— 1.09 元
  综合咖啡，摩卡，摩卡，奶泡 —— 1.39 元
  综合咖啡，大豆，摩卡，奶泡 —— 1.34 元
```

第 2 段的"深焙咖啡，摩卡，摩卡，奶泡 = 1.49"就是书上的 beverage2，一个字都没改。

### 标准库版（书里 Java I/O 那节的 Go 对应物）

```
=== 1. io.Reader 装饰器：同一个接口进，同一个接口出 ===
  原始            : *strings.Reader
  +capsReader     : main.capsReader
  +bufio.Reader   : *bufio.Reader
  +io.LimitReader : *io.LimitedReader
  读出来          : "DECORATOR"

=== 2. http.Handler 中间件：中间件就是装饰器 ===
  组装顺序：logging(auth(hello))
       [logging] GET -> /order
  没带令牌: 401 unauthorized （auth 短路，没有进到 hello）
  ...
  令牌正确: 200 hello from /order （穿过 logging → auth → hello）
```

这两段就是书里"装饰者模式在真实世界的用法"在 Go 的答案：**`io.Reader` 装饰器链**和
**`http.Handler` 中间件**。后者你以后写 Web 服务每天都会写，只是没人叫它 Decorator。

## 背后的设计原则

1. **类应该对扩展开放，对修改关闭**（OCP，本章新增）。加调料 = 加类型，不改旧类型。
2. **多用组合，少用继承**。装饰者就是"用组合（包一层）代替继承（复制 64 个子类）"。
3. 顺带把前两章的"封装变化""针对接口编程"用上了：`struct{ Beverage }` 里那个字段类型正是接口。

## Go 特有的注意点

### 1. 嵌入**接口**，而不是嵌入具体类型——这是本模式在 Go 里最关键的一行

```go
type Milk struct{ Beverage }   // ✅ 能包住任何 Beverage：Espresso、Decaf，甚至另一个 Milk
type Milk struct{ Espresso }   // ❌ 只能包 Espresso，装饰器当场失去意义
```

嵌入接口后，`Milk` 自动获得被装饰者的 `Description()` / `Cost()`（方法提升），
我们再用同名方法把它们**遮蔽**成"多加一份钱、多一句描述"。这就是 Go 版的 override。

### 2. 不需要 `CondimentDecorator` 这个中间抽象类

书上有三层：`Beverage`（抽象）→ `CondimentDecorator`（抽象，持有 beverage）→ `Mocha`。
Go 里 `struct{ Beverage }` 一行就同时表达了"我是 Beverage"和"我包着一个 Beverage"，
所以直接 `type Mocha struct{ Beverage }`。**抽象层数是 Java 的语法约束带来的，不是模式的本质。**

### 3. 装饰者用值类型，构造器返回具体类型

```go
func NewMocha(b Beverage) Mocha   // 而不是返回 Beverage
```

Go 官方建议实现方返回具体类型（新方法可以随便加，不动调用方）。虽然这里返回 `Beverage`
也能链式调用，但返回 `Mocha` 更符合 Go 的习惯，且在需要时仍可当 `Beverage` 用。
调料本身无状态、不可变，所以用值类型（`Mocha{}`），不需要指针。

### 4. 代价：装饰之后，具体类型就丢了

```go
b := NewMocha(Espresso{})
// b 的静态类型是 Beverage（或被装饰者），运行时是 Mocha
```

书上也提醒过：**装饰者会让"类型"变得不可靠**。所以不要写这种代码：

```go
if _, ok := b.(Mocha); ok { ... }   // ❌ 反模式：靠类型断言判断"加了什么调料"
```

要判断"加了什么"，应该让行为本身携带信息（比如 `Description()` 里已经有"摩卡"），
或者在接口上暴露语义方法（`HasMocha() bool`），而不是靠类型反查。
同理，装饰链很深时调试要一层层剥（Java I/O 的"注释地狱"就是前车之鉴）。

### 5. 顺序：价格可交换，语义不一定

- 算钱：`Milk(Mocha(x))` 和 `Mocha(Milk(x))` 同价（加法可交换）；
- 描述：拼接顺序不同，文本就不同；
- 中间件：**顺序有真实语义**。`logging(auth(hello))` 意味着未授权时也会打日志；
  `auth(logging(hello))` 则不会。装饰器不保证顺序无关，写链的时候要当成"洋葱从外往里"思考。

## 什么时候别用

- **只是参数不同**（加不加糖、大杯中杯）：加个字段/参数就够了，别上装饰器。
- **需要依赖具体类型做分支**：装饰器会毁掉类型信息，这种场景用策略 + 枚举更稳。
- **装饰链会固定到三五层以上且很少变化**：直接写一个组合结构体更易读；
  装饰器的价值在于"运行时任意组合"。

## 与相邻模式的区别（防混淆）

| 模式 | 接口关系 | 干什么 |
|---|---|---|
| **Decorator** | 同接口进、同接口出 | **加**责任（层层包装，可叠加） |
| **Strategy**（Ch.1） | 同接口，但作为**字段**被替换 | **换**算法（不包装，不叠加） |
| **Proxy**（Ch.11） | 同接口进、同接口出 | **控制**访问（不改变功能，只是加关卡） |
| **Adapter**（Ch.7） | **不同**接口进、不同接口出 | **转换**接口（让不兼容的东西能对接） |

Decorator / Proxy 长得几乎一样，区别在**意图**：装饰者是为了"加功能且允许叠加"，
代理是为了"替身/延迟加载/访问控制"，通常只有一层。

## 进一步阅读

- 《Head First 设计模式》第 2 版 第 3 章（Starbuzz + Java I/O 装饰者）
- Go 标准库里的同款：`io.LimitReader` / `bufio.NewReader` / `http.StripPrefix` / `http.TimeoutHandler`
  → https://pkg.go.dev/io#LimitReader 、 https://pkg.go.dev/net/http#Handler
- Go 官方博客 · Working with Errors in Go（用包装器逐层附加信息，是装饰者思想的另一面）:
  https://go.dev/blog/go1.13-errors
