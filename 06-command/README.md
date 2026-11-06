# 06 · Command 命令模式

## 一句话意图

> 把"请求"封装成对象，从而使你可用不同的请求对客户进行参数化；对请求排队或记录请求日志，
> 以及支持可撤销的操作。

（GoF："Encapsulate a request as an object, thereby letting you parameterize clients with different
requests, queue or log requests, and support undoable operations."）

**关键那句是"封装成对象"**——一旦请求变成对象，它就能被存起来、传出去、排队、回放、撤销。
这不是"又一个多态的例子"，这是**把动词变成名词**。

## 它解决的问题（书里的故事）

遥控器有 7 个插槽，每个插槽对应一件家电的开/关。如果遥控器"知道"家电：

```java
if (slot == 0) kitchenLight.on();
else if (slot == 1) garageDoor.up();
else if (slot == 2) ceilingFan.high();
// 每加一件家电，改这个类 + 重新编译 + 重新烧录固件
```

问题不在 `if` 有多丑，而在**方向**：遥控器（调用者）必须知道每台家电（接收者），
而家电厂商和遥控器厂商根本不认识彼此。这和第 4 章的病是同一个：**依赖方向反了**。

命令模式的解法：遥控器只持有 7 个"能 Execute 的东西"。

```go
type Command interface {
	Execute() string
	Undo() string
}
```

于是遥控器代码里再也没有一件家电的名字。而"是谁 + 做什么 + 用什么参数"这三件事，
全被一个命令对象吃进去了：

```
LightOnCommand{ light: 客厅的灯 }        // 谁 + 做什么
StereoOnWithCDCommand{ stereo: 音响 }    // 谁 + 做什么 + 参数（CD 模式）
```

## 参与者 → Go 的对应物

| 书上（Java） | 本目录的 Go 写法 | 说明 |
|---|---|---|
| `Command` 接口（`execute`/`undo`） | `Command` 接口 | 两个方法都返回 string（书里是打印） |
| `LightOnCommand` 等具体命令 | 同名结构体，持有接收者指针 | `Execute` 调接收者，`Undo` 反着调 |
| `Light` / `GarageDoor` / `CeilingFan` / `Stereo`（接收者） | 同名结构体，**不知道命令存在** | 复用自己原有的 API |
| `RemoteControl`（调用者/Invoker） | `RemoteControl` 结构体 | 只有 `[]Command` + 一个撤销槽 |
| `NoCommand`（空对象） | `NoCommand struct{}` | Go 里空结构体零值即可用 |
| `MacroCommand`（派对模式） | `MacroCommand` | 命令 + 组合：Execute 顺序、Undo 逆序 |
| `CeilingFanHighCommand` 等 4 个类 | 1 个 `CeilingFanSpeedCommand` + 4 个构造函数 | 只差一个常量，Go 里没必要建 4 个类型 |

## 结构图

```
        调用者 Invoker                命令 Command                接收者 Receiver
   ┌────────────────────┐      ┌──────────────────────┐      ┌──────────────┐
   │   RemoteControl    │      │   «interface»        │      │    Light     │
   │  onCommands  []Cmd │─────▶│  Execute() / Undo()  │─────▶│  On() Off()  │
   │  offCommands []Cmd │      └──────────┬───────────┘      └──────────────┘
   │  undoCommand Command│                │ 实现                    ▲
   └────────────────────┘      ┌──────────┴──────────┐             │
                               │  LightOnCommand     │─────────────┘
                               │  GarageDoorUpCommand│──▶ GarageDoor
                               │  CeilingFanSpeed... │──▶ CeilingFan
                               │  MacroCommand       │──▶ 一组 Command
                               │  NoCommand          │──▶ 什么也不做
                               └─────────────────────┘
```

**三条边，两个"互相不认识"**：调用者不认识接收者，接收者也不认识命令。
唯一的公共语言是 `Command` 接口。

## 可运行演示

```bash
go test ./06-command        # 8 个场景测试
go run ./06-command/demo
```

```
------ 遥控器面板 ------
[0] LightOnCommand               | LightOffCommand
[3] CeilingFanSpeedCommand       | CeilingFanSpeedCommand
[4] StereoOnWithCDCommand        | StereoOffCommand
[5] （未配置）                    | （未配置）
[撤销] （未配置）

=== 1. 按按钮：遥控器不知道自己在开什么 ===
  [0] 厨房的灯打开了
  [2] 车库的车库门升起来了
  [4] 客厅的音响打开了，客厅的音响切换到 CD 模式

=== 2. 撤销上一次操作 ===
  撤销 → 客厅的音响关上了

=== 3. 吊扇：撤销要回到「上一个速度」，而不是「关闭」 ===
  先调到中速 → 客厅的吊扇调到中速
  [3] 客厅的吊扇调到高速
  撤销 → 客厅的吊扇调到中速        ← 关键：撤销 = 恢复执行前的状态，不是"反向操作"
  撤销后速度：2（2=中速，而不是 0=关闭）

=== 4. 派对模式：宏命令一键开、一键关 ===
  [5] 客厅的灯打开了；音响打开并切到 CD；吊扇调到高速
  [5] 客厅的吊扇调到关闭；音响关上；灯关上        ← 逆序撤销
```

## 三个技术要点（书里各自都有，但很容易看漏）

### 1. 撤销不是"反向操作"，而是"恢复旧状态"

很多人第一反应是 `LightOffCommand = LightOnCommand 的反向`，这在灯上碰巧成立，
在吊扇上立刻翻车：高速的"反向"是关闭吗？不是——如果执行前在中速，撤销应该回到**中速**。

所以命令对象必须携带**执行前的快照**：

```go
func (c *CeilingFanSpeedCommand) Execute() string {
	c.prev = c.fan.Speed()   // 记账
	return c.fan.SetSpeed(c.set)
}
func (c *CeilingFanSpeedCommand) Undo() string {
	return c.fan.SetSpeed(c.prev)  // 还账
}
```

这是命令模式"能记住历史"的物理基础：**状态存在命令对象里，而不是存在遥控器里。**

### 2. 空对象（NoCommand）让"未配置"变成正常状态

书里用 `NoCommand` 初始化所有插槽，于是按下未配置的按钮不会 NPE，`Undo` 也不会炸。
Go 里更省事：`NoCommand{}` 零值就能用，连构造函数都不用。

### 3. 宏命令是"命令 + 组合模式"

`MacroCommand` 自己也是 `Command`，肚子里装着一组 `Command`——同类型进同类型出，
所以宏里还能套宏。Execute 正序、Undo **逆序**（后开的先关），这是宏命令唯一的坑。

## 命令 vs 策略（最常被问的问题）

两者长得几乎一样：都是"一个接口 + 一堆实现 + 运行时替换"。区别在**意图**和**对象去哪了**：

| | 策略 Strategy（Ch.1） | 命令 Command（Ch.6） |
|---|---|---|
| **回答什么问题** | **怎么做**（how）：同一件事的多种算法 | **做什么/让谁做**（what/who）：一次具体请求 |
| **典型方法形状** | 有参数、有返回值（`Cost() float64`、`Fly() string`）—— 算出一个值 | 无参（`Execute()`）—— 产生副作用，重点不是返回值 |
| **是否持有接收者** | 不需要（`FlyWithWings` 自己就能飞） | **必须**（要知道动哪盏灯、哪个门） |
| **谁调用它** | 使用它的上下文对象（鸭子调用 `flyBehavior.Fly()`） | 一个不懂业务的**调用者**（遥控器、队列、线程池） |
| **放在哪里** | 通常作为一个**字段**被长期持有 | 通常放在**集合/队列/栈**里，被存起来、排期、回放 |
| **有没有历史** | 通常没有（替换行为，不记录） | 通常有（undo、日志、事务） |
| **能不能组合** | 一般不组合 | 可以宏命令组合（Composite） |
| **存在的目的** | 让**算法**可替换、可独立演化 | 让**请求**可存储、可传递、可撤销 |

一句话记法：

> **策略 = 把算法参数化；命令 = 把请求对象化。**
> 策略是"同一件事的不同做法"；命令是"不同的事，每件事都打包成一个对象"。

还有两个特别实用的判别信号（来自代码形状，比定义更好用）：

1. **看它放在哪**：`duck.flyBehavior = FlyWithRocket{}`（替换一个字段）→ 策略；
   `queue = append(queue, cmd)` / `undoStack.Push(cmd)`（进集合）→ 命令。
2. **看它有没有接收者**：命令里一定有个 `light *Light` 这样的字段，而且方法是转发；
   策略里通常没有"被动的对象"，算法自己是主角。

顺带说清一个常见误会：**命令模式里的"命令"不是业务逻辑的替代品**。
`LightOnCommand` 只做"调 `light.On()`"这一件事，它不含任何业务判断；
业务逻辑仍在接收者（灯）里。命令只是**请求的信封**。

## 书里没写、Go 里很自然的一种写法：命令用函数类型

只有一两个方法的"动作型"接口，Go 里常常直接用函数类型：

```go
type Command func() string
type CommandWithUndo struct{ Do func() string; Undo func() string }
```

标准库的 `http.HandlerFunc` 就是「函数类型实现接口」的著名例子。取舍：

| | 接口（本目录） | 函数类型 |
|---|---|---|
| 需要携带状态（prevSpeed、目标设备） | 结构体天然可以 | 靠闭包捕获 |
| 需要 undo / 多个方法 | ✅ 自然 | 得用结构体包两个函数 |
| 想像 `%T` 那样打印、想塞进 map 做名字 | ✅ 有类型名 | 闭包没有名字，调试信息弱 |
| 极简场景（一个动作、一次性） | 略显啰嗦 | ✅ 一行搞定 |

命令比策略更常需要 undo 和"自报名字"，所以本目录选了接口。**代码形状跟着需求走，不跟着课本走。**

## 什么时候别用

- **调用者和接收者本来就该耦合**（A 只是 B 的内部实现细节）：加一层命令只会让人多跳一次定义。
- **没有排队/撤销/日志/远程/延迟执行这些需求**：那你要的其实只是"一个方法调用"或"一个策略"。
- **命令数量爆炸、每个只有一行转发**：这时往往是"业务逻辑被削得太薄"，考虑把相关的几个命令合并，
  或者像本目录处理吊扇那样参数化。

## 命令模式的真实用途（Go 世界里）

| 场景 | 表现 |
|---|---|
| 撤销/重做（编辑器、绘图） | 命令进 undo/redo 两个栈 |
| 任务队列/工作池 | 命令对象送进 channel，worker 逐个 `Execute` |
| 事务/操作日志（WAL） | 命令序列化后落盘，崩溃后重放 |
| GUI 按钮/菜单/快捷键 | 同一个命令对象被多个控件共享 |
| 远程调用（RMI/gRPC 雏形） | 命令被序列化后送到另一台机器执行 |

共同点：**请求需要"活"得比这次函数调用更久**。这就是为什么它必须是个对象。

## 进一步阅读

- 《Head First 设计模式》第 2 版 第 6 章（遥控器 + 撤销 + 派对模式）
- Go 标准库里的"命令"味道：`http.Handler`（请求被封装成对象传下去）、`os/exec.Cmd`（待执行的命令对象）
  → https://pkg.go.dev/net/http#Handler 、 https://pkg.go.dev/os/exec#Cmd
- GoF 原文：把请求封装为对象 → 可参数化、排队、记录、撤销
