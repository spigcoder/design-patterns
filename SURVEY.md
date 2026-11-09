# 剩余模式速览（Ch.7 – Ch.14）

> 不打算逐个写代码的模式集中在这里：每个模式 3~5 行，够用来判断"这个场景该不该上它"。
> 已实现的模式见各目录的 README；本文件只做地图。

## 全景：按"它想让什么变"分类

| 家族 | 模式 | 它们共同的一句话 |
|---|---|---|
| **流程与步骤** | Template Method(8)、Iterator(9)、State(10) | 把"步骤/遍历/状态迁移"从业务里抽出来 |
| **结构与包装** | Adapter(7)、Facade(7)、Decorator(3)、Proxy(11)、Bridge(14) | 都是"包一层"，区别只在这层包裹**改变了什么** |
| **组合与共享** | Composite(9)、Flyweight(14) | 把多个对象当整体看 / 把相同的部分共享掉 |
| **请求与协调** | Command(6)、Chain of Responsibility(14)、Mediator(14)、Observer(2) | 请求怎么被传递、排队、分发、广播 |
| **创建** | Factory Method / Abstract Factory(4)、Singleton(5)、Builder(14)、Prototype(14) | 对象怎么出生 |
| **解析与操作** | Visitor(14)、Interpreter(14) | 对一棵结构稳定的树做不断新增的操作 |

---

## 流程与步骤

### Template Method 模板方法（Ch.8）

**骨架在父类，洞留给子类。** 父类写好 `prepare → cook → serve`，中间某些步骤 `abstract` 由子类填。

- Go 里的形状：`结构体 + 注入的钩子`。
  ```go
  type Beverage struct {
      brew func() string   // ← 洞
      addCondiments func() string
  }
  func (b Beverage) Prepare() string {   // ← 骨架
      return "boil water; " + b.brew() + "; pour; " + b.addCondiments()
  }
  ```
- **我们已经用过它了**：`04-factory/store.go` 里的 `OrderPizza`（流程固定、只调 `createPizza` 钩子）就是模板方法。
- 与 Strategy 的区别：模板方法用**继承/嵌入**把"部分步骤"参数化，流程本身在父类；策略用**组合**把"整块算法"换掉，流程在外面的调用者手里。

### Iterator 迭代器（Ch.9，书里和 Composite 同章）

**把"怎么遍历"从"遍历什么"里拆出去。** 书里 Diner 用数组、PancakeHouse 用链表，用同一个迭代器接口统一。

- Go 里你**每天都在用**：`for range` 就是迭代器；Go 1.23+ 还有 `iter.Seq[T]`（range-over-func）。
- 记住一句：**"自己写迭代器"在 Go 里通常是坏味道**，返回 `[]T` 或 `iter.Seq[T]` 更符合口味。真正会自己写的时候，是数据源大得不能一次装进内存（数据库游标、文件扫描）。

### State 状态模式（Ch.10）

**把"if 状态 == X 就 blah"变成一组状态对象，每个状态自己决定"下一步是哪个状态"。**
书中例子：口香糖机（待币 → 摇柄 → 出货 → 售罄）。

```go
type State interface {
    InsertCoin(m *Machine) string
    TurnCrank(m *Machine) string
}
type NoCoinState struct{}          // 每种状态一个类型
func (NoCoinState) InsertCoin(m *Machine) string {
    m.state = HasCoinState{}       // ← 状态自己完成状态迁移
    return "收到硬币"
}
```

**它是策略模式的双胞胎，区别只有一个**：状态迁移由状态对象自己触发（内部驱动、自发），
策略由客户端注入与替换（外部驱动）。见本文件末尾的对照表。

- 常见误用：状态少（2 个）且没有迁移时，一个 `enum + switch` 就够了，别上状态模式。
- 加分用法：状态对象可以**内存级唯一**（空结构体 / `sync.Once`），只有"上下文"在变。

### Proxy 代理模式（Ch.11）

**同接口进、同接口出，但目的是"控制访问"而不是"增加功能"。**
书中例子：口香糖机的远程监视器（RMI 代理）、虚拟代理（图片先用占位图，真图加载完再换）。

Go 里几乎不需要"自己写代理"：`net/rpc`、`grpc` 的 stub 就是代理；懒加载用 `sync.Once`；
访问控制在 `http.Handler` 中间件里做。**能叫出它名字，比会写它更重要**（面试/读源码时有用）。

---

## 结构与包装：五个"包一层"，只差意图

| 模式 | 对外接口 | 里面包的是 | 目的 |
|---|---|---|---|
| **Decorator**(3) | **同一个** | 一个同类对象 | **加**责任，可叠加 |
| **Proxy**(11) | **同一个** | 一个同类对象 | **控**访问（延迟/权限/远程/计数） |
| **Adapter**(7) | **变成另一个** | 一个**接口不兼容**的对象 | **转**接口 |
| **Facade**(7) | **简化版** | **一组**子系统对象 | **收**入口（一键式操作） |
| **Bridge**(14) | 抽象侧**持有**实现侧 | 一条独立变化的实现线 | **分**两个维度，各自演化 |

- **Adapter vs Facade**：Adapter 让旧的能用（接口转接，1 对 1）；Facade 让复杂的变简单（1 对多，简化调用）。
- **Bridge vs Strategy**：Bridge 是"结构上的两个维度（形状 × 颜色）各自有继承线"，策略只是"一个维度可替换"。Go 里两者长得像，因为 Go 只有接口；**目的不同**：Bridge 是为了避免 `M×N` 个子类。
- Go 里的 Adapter：`io.Reader` 适配任意数据源（`strings.Reader`）、`http.HandlerFunc`（函数适配成接口）。

---

## 组合与共享

### Composite 组合模式（Ch.9）

**把"单个"和"一堆"当成同一种类型**，于是树可以无限嵌套。

- 书里：菜单 + 子菜单，`Menu` 和 `MenuItem` 都实现同一个接口。
- Go 里我们已经见过了：`06-command` 的 `MacroCommand`（宏命令本身也是 Command）。
- 与 Decorator 的区别：Composite 的**孩子是复数（一组）**，Decorator 的孩子**通常只有一个**。

### Flyweight 享元（Ch.14）

**把大量对象的"共同部分"提取共享，只让"差异部分"按实例存。**
书里：树（Tree）共享 TreeType，只为每棵树存坐标。

- Go 里的样子：`map[Key]*Shared` 缓存 + 实例只存 key；字符串常量池、`intern` 都是它。
- 现代 Go 里往往被"结构体拆成指针 + 值"或"用 ID 引用"替代，很少专门实现。

---

## 请求与协调

| 模式 | 一句话 | Go 里最像它的东西 |
|---|---|---|
| **Observer**(2) | 一变全通知（一对多广播） | `chan` 事件总线、`sync.Cond`、回调 |
| **Command**(6) | 请求变成对象（可排队/撤销） | `http.Handler`、`os/exec.Cmd`、任务队列 |
| **Chain of Responsibility**(14) | 请求沿链传递，谁处理谁停 | **中间件链** `logging(auth(handler))` |
| **Mediator**(14) | 把 N×N 的互相通信改成 N×1 的中介 | 事件循环、调度器、聊天室 hub |

- **CoR vs Decorator**：CoR 的链成员可以**终止**请求（谁处理谁负责），Decorator 一定**继续往下传**。
- **Mediator vs Observer**：Observer 是被观察者直接广播；Mediator 让对象之间**谁也不认识谁**，都只认识中介。

---

## 创建

| 模式 | 一句话 | Go 里的现实写法 |
|---|---|---|
| Factory Method / Abstract Factory(4) | 把 new 收拢、把族打包 | 钩子函数 / `Provider` 接口（已实现，见 `04-factory`） |
| Singleton(5) | 全局唯一实例 | `sync.Once` + package 变量；更多时候用 DI 传实例 |
| **Builder**(14) | 分步构建复杂对象 | **Functional Options**（`func(*Config)` 变参），比 Builder 类型更 Go |
| **Prototype**(14) | 克隆自己 | Go 值语义天然是拷贝；深拷贝写 `Clone()` 方法 |

```go
// Go 里 Builder 的替代品：functional options
type Option func(*Server)
func WithPort(p int) Option { return func(s *Server) { s.port = p } }
func NewServer(opts ...Option) *Server { /* 逐个应用 */ }
```

---

## 解析与操作

### Visitor 访问者（Ch.14）

**在不修改元素类的前提下，为它增加新操作**：元素 `Accept(v)`，访问者 `VisitXxx(e)`，靠**双重分派**决定调用谁。

- Go 里的现实：**大多数情况用不上**。Go 的 `switch v := node.(type)` 已经能按类型分派，
  Visitor 带来的收益（把操作聚在一个类型里）在 Go 里往往抵不过它的啰嗦。
- 什么时候真的用：元素类型**稳定**而操作**不断新增**，且你想让每个操作自成一体（如 AST 的多个 pass）。

### Interpreter 解释器（Ch.14）

**把一门小语言的语法表示成类树，对树求值。** 书里是社会关系表达式、正则表达式。

- Go 里：真要做 DSL，通常是**自己写词法+语法解析 → AST → 一个 eval 函数**，
  很少为每种语法规则建一个类。**知道它存在即可。**

---

## Ch.12 复合模式 / Ch.13 真实世界（MVC）

- **Ch.12（鸭子模拟器）**：把 Observer + Decorator + Factory + Adapter 揉进一个程序，
  想说的道理只有一句：**真实系统里模式是叠着用的**。我们仓库已经演示了这一点——`04-factory` 里
  工厂方法和抽象工厂就是叠起来的。
- **Ch.13（MVC）**：MVC 本身就是**复合模式**——Model 用 Observer 通知 View，
  View 用 Composite 组织界面，Controller 用 Strategy 处理输入。
  在 Go 里对应的是 `net/http`：`Handler`（控制）→ 业务（模型）→ `html/template`（视图）。

---

## 一张对照表：最容易混的三对

| 容易混 | 区别一句话 |
|---|---|
| **State vs Strategy** | 结构一样。State 的**下一个状态由状态自己决定**（内部驱动）；Strategy 由**外部注入/选择**（外部驱动）。State 通常记得"当前是谁"，Strategy 只是被持有。 |
| **Proxy vs Decorator** | 都同接口包同接口。Decorator **加功能且可叠加**；Proxy **控访问**，通常一层、且常由框架生成。 |
| **Adapter vs Facade** | Adapter 改成**别人的接口**（1→1）；Facade 提供**简化入口**（多→1）。 |

## 在 Go 里"知道就好、别急着实现"的几个

`Visitor`、`Interpreter`、`Prototype`、`Bridge`、`Flyweight`、`Singleton`。
它们的设计动机大多来自 Java/C++ 的语法限制（没有隐式接口、没有值语义、没有函数值、没有 `sync.Once`），
在 Go 里要么有更轻的替代，要么根本不成为问题。**看它们是为了读懂别人的代码，不是为了写出来。**
