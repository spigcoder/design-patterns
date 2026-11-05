# 04 · Factory 工厂方法 / Abstract Factory 抽象工厂

> 这一章真正值钱的不是"怎么造对象"，而是**"谁被允许知道具体类"**。代码很短，思想很长。

## 一句话意图

- **工厂方法（Factory Method）**：定义创建对象的接口，但让**实现这个接口的类**决定实例化哪一个类。
  工厂方法让类的实例化推迟到子类。
- **抽象工厂（Abstract Factory）**：提供一个接口，用于创建**相关或依赖对象的家族**，而不需要明确指定具体类。

两者都管"对象怎么出生"，区别只有一个：**工厂方法一次造一个产品，抽象工厂一次造一族产品。**

## 它解决的问题（书里的故事）

比萨店要加盟：纽约店用薄饼面团 + Marinara 酱 + Reggiano 奶酪 + 新鲜蛤蜊；
芝加哥店用厚饼 + 李子番茄酱 + Mozzarella + 冷冻蛤蜊。**流程完全一样**（准备→烤→切→装盒），
**只有原料和品类不同**。

如果不用工厂，`orderPizza()` 里会写满 `new NYStyleCheesePizza()`、`new ChicagoStyleClamPizza()`——
每开一家新店就要回来改一遍这个方法。问题不在"要 new"，而在**new 的位置**：
决定"用哪个具体类"的代码，混进了"根本不关心具体类"的流程里。

书用两步解决：

1. **工厂方法**：把"做哪种披萨"抽成一个钩子 `createPizza(type)`，`orderPizza()` 只调钩子。
   具体店（`NYPizzaStore`）实现钩子 → 流程不动，新增店只需新增一个实现。
2. **抽象工厂**：把"这一族原料从哪来"抽成 `PizzaIngredientFactory`。
   披萨类只写 `factory.CreateDough()`，永远不知道拿到的是薄饼还是厚饼。

## 三种"工厂"的对照表（本章最该背下来的东西）

| | 简单工厂（**不是 GoF 模式**） | 工厂方法 | 抽象工厂 |
|---|---|---|---|
| 谁决定具体类 | 一个函数里的 `switch` | **子类 / 注入的钩子** | **一个工厂对象** |
| 复用机制 | 无 | 继承 + 覆写（=模板方法） | 组合 + 委托 |
| 一次产出 | 一种产品 | 一种产品 | **一族**产品 |
| 扩展方式 | **改 switch**（违反 OCP） | **加一个实现**（不改旧代码） | 加一个新工厂 + 新产品类 |
| 变化的维度 | 1 个（品类） | 1 个（品类） | **2 个**（品类 × 风格族） |
| Go 里的样子 | `func newPizza(kind string)` | `createPizza func(string) (Pizza, error)` 字段 | `PizzaIngredientFactory` 接口 + 每家一个实现 |
| 软肋 | 新增品类要改老代码 | 每种产品都要一个实现者 | **加一个产品方法要改所有工厂** |

一句话：
**简单工厂把所有变化留在函数里；工厂方法把变化推给新的实现者；抽象工厂把"族"打包成对象。**

本书的做法是把两者**叠起来**：

```
PizzaStore（店 = 风格维度）── 工厂方法决定"做哪种披萨"
   └── PizzaIngredientFactory（原料 = 族维度）── 抽象工厂决定"用什么原料"
```

## 参与者 → Go 的对应物

| 书上（Java） | 本目录的 Go 写法 | 说明 |
|---|---|---|
| `Pizza` 抽象类 | `Pizza` 接口 | 产品合同 |
| `CheesePizza`/`ClamPizza`/`VeggiePizza` | 三个结构体，字段是 `Dough`/`Sauce`/... | 持有原料，但**不认识工厂的具体类型** |
| `PizzaIngredientFactory` 接口 | 同名接口 | 抽象工厂，4 个 `CreateXxx` |
| `NYPizzaIngredientFactory` | `NewYorkIngredientFactory struct{}` | 无状态 → 空结构体 |
| `Dough`/`Sauce`/`Cheese`/`Clams` 各自一个接口 + 每族一个类 | 带名字的字符串类型 + 常量 | **不要为空行为建空结构体**，见下文 |
| `PizzaStore.orderPizza()`（模板方法） | `(*PizzaStore).OrderPizza` | 流程固定，只调钩子 |
| `abstract Pizza createPizza(String)` | `CreatePizzaFunc` 函数字段 | 工厂方法的 Go 表达 |
| `NYPizzaStore extends PizzaStore` | `NewYorkPizzaStore()` 构造函数 | 不需要新类型，填钩子即可 |
| `SimplePizzaFactory` | `NewSimplePizza` | 书自己说"不是模式"，保留作对照 |

## 结构图

```
        ┌────────────────────────────┐
        │        PizzaStore          │  name + createPizza func（钩子）
        │  OrderPizza(kind)          │  ← 模板方法：prepare/bake/cut/box
        └──────────┬─────────────────┘
                   │ 调用钩子（不认识具体产品）
                   ▼
        ┌────────────────────────────┐
        │     «interface» Pizza      │  Name() / Prepare()
        └──────────┬─────────────────┘
     ┌─────────────┼──────────────┐
 CheesePizza   ClamPizza     VeggiePizza
     │             │              │
     └─────────────┴──────────────┘  都只持有 ↓ 这个接口
                   ▼
        ┌────────────────────────────────┐
        │ «interface» PizzaIngredient-    │  CreateDough / CreateSauce /
        │   Factory                       │  CreateCheese / CreateClams
        └───────┬────────────────┬────────┘
   NewYorkIngredientFactory  ChicagoIngredientFactory
        （薄饼族）                （厚饼族）
```

两个箭头方向值得盯一眼：**店 → 产品**、**产品 → 原料工厂**，两条都指向接口，没有一条指向具体类。

## 可运行演示

```bash
go test ./04-factory       # 7 个场景测试（含"新增东京店不改旧代码"）
go run ./04-factory/demo
```

```
=== 1. 两家店各点一份奶酪披萨：同一款产品，不同原料族 ===
  【纽约披萨店】
    制作：奶酪披萨
    揉面团：薄饼面团
    加酱料：Marinara 酱
    撒奶酪：Reggiano 奶酪
    烘烤 25 分钟
    切成 8 块
    装入纽约披萨店的盒子
  【芝加哥披萨店】
    制作：奶酪披萨
    揉面团：厚饼面团
    ...

=== 3. 披萨只认工厂接口：同一块奶酪披萨，换一个原料工厂就换了一族原料 ===
  factory.NewYorkIngredientFactory 的原料：
    揉面团：薄饼面团
    ...
  factory.ChicagoIngredientFactory 的原料：
    揉面团：厚饼面团
    ...

=== 4. 点一个菜单上没有的品类 ===
  纽约披萨店 点单失败: 没有 "榴莲" 这种披萨
```

第 3 段是抽象工厂的"神迹时刻"：**同一个 `CheesePizza` 类型**，什么都不改，换一个工厂参数，
就换了一整族原料。它证明了披萨类里确实没有任何"纽约/芝加哥"的知识。

测试 `TestAddingANewCityNeedsNoChange` 是工厂方法的对应证明：东京店和东京原料工厂
**全部写在测试包里**，`factory` 包一个字没动，就能开出一家新店。

## 背后的设计原则：依赖倒置（DIP）

> **依赖抽象，不要依赖具体类。**（Dependency Inversion Principle）

书里给了三条"准则"，但明确说**不可能 100% 做到**：

1. 变量不应该持有具体类的引用 → 用接口/抽象类型声明；
2. 类不应该从具体类派生 → 别为了复用去继承一个实现；
3. 方法不应该覆盖基类里已实现的方法 → 覆写会把父类的假设带进来。

**为什么工厂是 DIP 的化身？** 因为任何程序都必须 `new`，DIP 做不到"零具体类依赖"。它能做到的是
**把"知道具体类"的代码压缩到极少的、固定的几个地方**——那几个地方就是工厂。
换句话说：

> 工厂 = 应用里被官方允许依赖具体类的角落。它的存在，是为了让**其余所有代码**不必依赖具体类。

这句话解释了本章那个"好像学到了又好像没学到"的感觉：**你学到的不是一种造对象的技巧，
而是"把不变的部分与决定具体类的部分分开"的工序。**

## Go 特有的注意点

### 1. 工厂方法有两种 Go 写法，本目录用的是 A

```go
// 写法 A：函数字段（本目录）
type PizzaStore struct{ createPizza func(string) (Pizza, error) }

// 写法 B：接口 + 嵌入（更接近书里的形状）
type Creator interface{ CreatePizza(kind string) (Pizza, error) }
type PizzaStore struct{ Creator }
```

| | A 函数字段 | B 接口嵌入 |
|---|---|---|
| 代码量 | 3 行搞定 | 要额外定义一个接口 + 每家店一个类型 |
| 能不能内联闭包 | ✅ | ❌ |
| 能不能被 mock / 复用命名实现 | 一般 | ✅ 更好 |
| 阅读体验 | 需要看构造函数才知道钩子是什么 | 类型名即文档 |

Go 标准库里两种都有：`http.HandlerFunc`（函数类型当接口用，类似 A）和 `http.Handler`（接口，B）。
**没有绝对正确的那个，看实现会不会被复用。**

### 2. 原料不要为"对称"建空结构体

书里是 `ThinCrustDough implements Dough` —— 因为 Java 每个原料是一个类。
Go 里原料没有行为、只需要"能区分能打印"，所以：

```go
type Dough string
const ThinCrustDough Dough = "薄饼面团"
```

保留**具名类型**依然有类型安全（`CreateDough()` 返回 `Dough`，把 `Sauce` 塞进去编译不过），
但省掉了 8 个只用来当标签的类。这是本章最实用的 Go 经验：**结构体该有行为，标签才用常量。**

### 3. 构造函数要老实返回 error

Java 版靠异常；Go 里 `createPizza` 直接返回 `(Pizza, error)`，`OrderPizza` 用 `%w` 包装后往上抛。
别写"找不到就返回 nil 然后祈祷"。

### 4. 包名别叫 `factory`

本目录为了和 Ch.1–3 对齐才叫 `factory`。真实项目里包名应该表达**领域**（`pizza`），
而不是表达模式（Go 官方建议避免 `util`/`common`/`types` 这类泛名，`factory` 也有同样的味道）。

### 5. 抽象工厂的软肋：加产品方法要改所有工厂

想让工厂也提供"蔬菜"（`CreateVeggies() []string`），就得改 `PizzaIngredientFactory` 接口，
于是纽约、芝加哥、以及你测试里那个东京工厂**全都要改**。这正是接口演化的老问题（Ch.1 聊过）。
所以抽象工厂的产品**族成员要一次性想清楚**，别边写边加。

### 6. 延伸：注册表式工厂（Go 里非常常见）

当"品类"会不断增长且来自不同包时，Go 项目通常不用 switch，而是用注册表：

```go
var registry = map[string]func() Pizza{}

func Register(name string, make func() Pizza) { registry[name] = make }

func New(kind string) (Pizza, error) {
    make, ok := registry[kind]
    if !ok { return nil, fmt.Errorf("未知品类 %q", kind) }
    return make(), nil
}
```

标准库就有现成的例子：`database/sql` 的驱动注册。驱动包在自己的 `init()` 里 `sql.Register("mysql", ...)`，
上层只写 `sql.Open("mysql", dsn)` —— **上层不认识任何具体驱动，新增数据库不用改上层代码**。
这就是"工厂方法 + 注册表"在 Go 里最成功的应用之一。

## 什么时候用哪个（决策清单）

```
你有几个正交的变化维度？
├─ 0 个（就这么一种，永远不会变） → 直接 new，别上工厂
├─ 1 个（品类会不断增加）
│   ├─ 只有一处 new、改动不频繁、团队不共享 → 简单工厂足够（别过度设计）
│   └─ 多处 new、要按场景（城市/租户/环境）分派 → 工厂方法（钩子注入）
└─ 2 个（品类 × 风格族）→ 抽象工厂 + 工厂方法叠起来（本章）
```

**反证**：如果某个"工厂"只有一个实现、也永远不会有第二个，那它大概率应该被删掉，
直接 `new` 就好——模式不是资历，是成本。

## 进一步阅读

- 《Head First 设计模式》第 2 版 第 4 章（比萨店 + DIP）
- Go 标准库的注册表工厂：`database/sql` 的 `Register` / `Open`
  → https://pkg.go.dev/database/sql#Register
- Go Code Review Comments（返回具体类型、接口放使用方、包名）
  → https://go.dev/wiki/CodeReviewComments
- Martin Fowler · Inversion of Control（DIP 与 IoC 的关系，容易混）
  → https://martinfowler.com/bliki/InversionOfControl.html
