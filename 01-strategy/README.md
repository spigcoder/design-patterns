# 01 · Strategy 策略模式

## 一句话意图

> 定义算法族，分别封装起来，让它们之间可以互相替换；策略让算法的变化独立于使用它的客户。

在书里，"算法族"就是一群**会飞的方式**（`FlyBehavior`）和一群**会叫的方式**（`QuackBehavior`）；
"客户"就是鸭子。鸭子不实现飞和叫，它只是**持有一个会飞的东西**、**持有一个会叫的东西**，把活儿转包出去。

## 它解决的问题（书里的故事）

SimUDuck 一开始用继承：`Duck` 有 `fly()`、`quack()`，子类各显神通。第一周的"改进"——把 `fly()`
提到父类——立刻出事：

1. `RubberDuck`（橡皮鸭）会飞了；
2. `DecoyDuck`（诱饵鸭）会飞还会嘎嘎叫；
3. 用覆盖（override）修补之后，橡皮鸭得写 `fly() { // do nothing }`，每加一种鸭子就多一份重复代码；
4. 真正的灾难是**需求变更**：书上说想让鸭子装火箭助推器。用继承，你得去改**每一个会飞的子类**。

根本原因一句话：**"飞"和"叫"是变化的部分，却被写死进了不变的部分（鸭子本身）。**

## 参与者 → Go 的对应物

| 书上（Java） | 本目录的 Go 写法 | 说明 |
|---|---|---|
| `FlyBehavior` 接口 | `FlyBehavior interface { Fly() string }` | Go 接口隐式满足，不需要 `implements` |
| `FlyWithWings` / `FlyNoWay` / `FlyWithRocket` | 同名空结构体 | 无状态算法用空结构体最省事（零值可用、无分配） |
| `QuackBehavior` 接口 | `QuackBehavior interface { Quack() string }` | 一个模式里有多个平行策略族 |
| `Duck` 抽象类（持有策略 + `display()` 抽象） | `Duck` 接口 + `BaseDuck` 结构体 | Go 无继承 → **接口 + 嵌入** 两件套 |
| `MallardDuck extends Duck` | `type MallardDuck struct{ *BaseDuck }` | 嵌入获得方法提升，等价于"继承来的方法" |
| 覆盖 `display()` | 自己实现 `Display()` | 外层方法遮蔽（shadow）晋升来的同名方法 |
| `setFlyBehavior()` | `(*BaseDuck).SetFlyBehavior` | 晋升给所有具体鸭子 → 运行时换策略 |
| `performFly() { flyBehavior.fly(); }` | `func (d *BaseDuck) PerformFly() string { return d.flyBehavior.Fly() }` | 这就是"组合 + 委托"那一行 |
| 客户端 `Duck d = new MallardDuck()` | `var d strategy.Duck = strategy.NewMallardDuck()` | 变量声明为**接口**，具体类型只出现在构造处 |

## 结构图

```
                       ┌──────────────┐
                       │  «interface» │      ← 客户端只依赖这两个抽象
                       │     Duck     │
                       └──────┬───────┘
                              │ 实现（嵌入 + 遮蔽 Display）
      ┌───────────┬───────────┼───────────┬────────────┐
 MallardDuck  RedheadDuck  RubberDuck  DecoyDuck   ModelDuck
      │           │            │           │            │
      └───────────┴────────────┴───────────┴────────────┘
                              │ 各自嵌入
                       ┌──────┴───────┐
                       │   BaseDuck   │   flyBehavior   : «interface» FlyBehavior
                       │              │   quackBehavior : «interface» QuackBehavior
                       └──────────────┘
                              │ 委托调用
        FlyBehavior ◄─────────┴─────────► QuackBehavior
        ├ FlyWithWings                     ├ Quack
        ├ FlyNoWay                         ├ Squeak
        └ FlyWithRocket                    └ MuteQuack
```

## 可运行演示

```bash
go test ./01-strategy              # 4 个测试：委托正确 / 运行时可换 / 新增策略不改旧代码 / 客户端只认接口
go run ./01-strategy/demo
```

```
=== 1. 各种鸭子的当前表现 ===
我是一只绿头鸭 | 我在用翅膀飞 | 嘎嘎嘎 | 我在浮水（所有鸭子都会）
我是一只红头鸭 | 我在用翅膀飞 | 嘎嘎嘎 | 我在浮水（所有鸭子都会）
我是一只橡皮鸭 | 我根本不会飞 | 吱吱吱 | 我在浮水（所有鸭子都会）
我是一只木头诱饵鸭 | 我根本不会飞 | <<沉默>> | 我在浮水（所有鸭子都会）

=== 2. 运行时动态换行为（策略模式的核心收益）===
我是一只塑料模型鸭 | 我根本不会飞 | 嘎嘎嘎 | 我在浮水（所有鸭子都会）
>> 给模型鸭装上火箭助推器……
我是一只塑料模型鸭 | 我在用火箭助推器飞！ | 嘎嘎嘎 | 我在浮水（所有鸭子都会）

=== 3. 换叫声也一样，鸭子本身一行代码都不用改 ===
>> 给绿头鸭换成诱饵鸭的静音叫声……
我是一只绿头鸭 | 我在用翅膀飞 | <<沉默>> | 我在浮水（所有鸭子都会）
```

第 3 段最能说明问题：`MallardDuck` 的源码里**没有**任何"静音"字样，却能静音。行为是插进去的，不是继承来的。

## 背后的设计原则

1. **封装变化** —— 找出会变的部分，独立出来。会变的是 fly/quack，不会变的是 swim/display。
2. **针对接口编程，而不是针对实现编程** —— `BaseDuck` 的字段类型是接口；客户端变量的类型是 `Duck`。
3. **多用组合，少用继承** —— 鸭子"有一个"飞行行为（has-a），而不是"是一种会飞的鸭子"（is-a）。
   Erich Gamma 本人说过：Strategy 是他讲"组合优于继承"时用的原型例子，因为组合才能做到**运行时插拔**。

## Go 特有的注意点

1. **没有抽象类，所以"抽象类"要拆成接口 + 可嵌入结构体。** 书里 `Duck` 既是类型又是代码复用单元；
   Go 里 `Duck`（接口，定义合同）和 `BaseDuck`（结构体，提供复用）分开，这是常规做法。
2. **接口是隐式满足的**，所以"策略族"扩展起来极便宜：测试文件里那个 `flyWithJetpack` 甚至定义在
   另一个包里，什么都不用注册就能塞给鸭子。这也是 Go 里 `sort.Interface`、`io.Writer` 那套风格。
3. **单方法策略可以用函数类型代替接口。** 如果只有一个方法，`type FlyBehavior func() string` 往往更轻。
   这里保留接口，是为了和书上结构一一对应、对照学习；实战里两种都常见。
4. **接口应该定义在使用方，而不是实现方；实现方返回具体类型。**
   Go 官方 Code Review Comments 明确写了这一点（见下方链接）。所以别在"行为包"里预先堆一堆接口。
5. **不要为了模式而模式。** `Swim()` 就没有抽成策略——因为"所有鸭子都这样浮水"，它不变。
   给不变的东西加一层抽象，只会让读者多跳一次定义。

## 什么时候别用

- 变化点根本不存在（只有一个实现、且永远不会有第二个）：多一层接口只是多一跳间接。
- 策略之间需要互相知道对方、需要共享大量状态时：策略会退化成"散落一地的 if-else"，此时可能该用 State。
- 行为差异只是参数不同（飞行速度、叫声音量）：一个带参数的方法比一族类更清楚。

## 进一步阅读

- 《Head First 设计模式》第 2 版 第 1 章
- Erich Gamma 谈这两条原则（Artima 访谈，一手材料）: https://www.artima.com/articles/design-principles-from-design-patterns
- Go Code Review Comments · Interfaces: https://go.dev/wiki/CodeReviewComments#interfaces
- Effective Go · Interfaces: https://go.dev/doc/effective_go#interfaces
