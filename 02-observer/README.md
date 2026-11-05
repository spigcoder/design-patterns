# 02 · Observer 观察者模式

## 一句话意图

> 定义对象之间的一对多依赖：当一个对象改变状态时，所有依赖它的对象都会自动收到通知并更新。

（GoF 原句意译："Defines a one-to-many dependency between objects so that when one object changes state,
all of its dependencies are notified and updated automatically."）

## 它解决的问题（书里的故事）

Weather-O-Rama 公司要造气象站：`WeatherData` 对象拿到新的温度/湿度/气压后，要更新三块显示板
（当前状况、统计、预报）。

天真的第一版是这段代码：

```java
public void measurementsChanged() {
    float temp = getTemperature();
    float humidity = getHumidity();
    float pressure = getPressure();
    currentConditionsDisplay.update(temp, humidity, pressure);
    statisticsDisplay.update(temp, humidity, pressure);
    forecastDisplay.update(temp, humidity, pressure);
}
```

书上一针见血地指出它的三个毛病（对应三条设计原则）：

1. **针对具体实现编程**：这里写死了三个具体显示板类型，加一块第四块板子就得改这个方法。
2. **每次有新的显示板都要改 `WeatherData`**：改的是"不变的部分"（气象站），而变化的其实是"谁想订阅"。
3. **无法在运行时增删订阅者**：显示板集合被编译期定死了。

## 参与者 → Go 的对应物

| 书上（Java） | 本目录的 Go 写法 | 说明 |
|---|---|---|
| `Subject` 接口（register/remove/notify） | `Subject` 接口 | 发布者合同 |
| `Observer` 接口（`update()`） | `Observer interface { Update(Measurement) }` | 订阅者合同，只有 1 个方法 |
| `DisplayElement` 接口（`display()`） | `DisplayElement interface { Display() string }` | 书里单列，强调"观察者≠显示" |
| `WeatherData implements Subject` | `WeatherData` 结构体 + `sync.RWMutex` | 自己实现 Subject，不用标准库的 `Observable` |
| `update(float,float,float)` 推模型 | `Update(Measurement)` | 见下文「为什么把三个 float 换成一个结构体」 |
| `measurementsChanged()` | `SetMeasurements()` | 更新状态 + 通知 |
| `CurrentConditionsDisplay` 等 | 同名的 4 个结构体 | 构造函数里自注册 |
| `new CurrentConditionsDisplay(weatherData)` | `NewCurrentConditionsDisplay(subject)` | 依赖 `Subject` 接口而不是 `*WeatherData` |
| 观察者列表（数组从头重建） | `[]Observer` + `append` | 动态集合，容量自动增长 |

## 结构图

```
                       ┌───────────────┐
                       │ «interface»   │  RegisterObserver / RemoveObserver / NotifyObservers
                       │    Subject    │
                       └───────┬───────┘
                               │ 实现
                        ┌──────┴────────┐
                        │  WeatherData  │  mu, observers []Observer, current Measurement
                        └──────┬────────┘
                               │ 遍历快照后逐个调用
                               ▼
                       ┌───────────────┐
                       │ «interface»   │  Update(Measurement)
                       │   Observer    │◄──────────┐
                       └───────┬───────┘           │
        ┌──────────────┬───────┴────────┬──────────┴───┐
 CurrentConditions  Statistics      Forecast      HeatIndex
   Display           Display         Display       Display
        │              │                │              │
        └──────────────┴────────────────┴──────────────┘
                     （都另实现 DisplayElement）
```

单向依赖是关键：**Subject 只知道 `Observer` 接口，观察者之间互相不知道**。

## 可运行演示

```bash
go test -race ./02-observer     # 7 个测试，含并发竞态检测
go run ./02-observer/demo
```

```
=== 第一次测量 ===
  当前状况：气温 26.5°C，湿度 65.0%
  统计（1 次）：平均 26.5°C，最高 26.5°C，最低 26.5°C
  预报：气压 1013.2 hPa，等下一次读数再预报

=== 第二次测量（气压下降） ===
  当前状况：气温 28.0°C，湿度 70.0%
  统计（2 次）：平均 27.2°C，最高 28.0°C，最低 26.5°C
  预报：气压 1013.2→1005.0 hPa，注意降温和降雨

=== 第三次测量（气压回升） ===
  当前状况：气温 24.5°C，湿度 62.0%
  统计（3 次）：平均 26.3°C，最高 28.0°C，最低 24.5°C
  预报：气压 1005.0→1010.5 hPa，天气在好转

=== 运行时挂上一块新显示板（热指数），主题无需改动 ===

=== 第四次测量（湿热，气压又降） ===
  当前状况：气温 33.0°C，湿度 75.0%
  统计（4 次）：平均 28.0°C，最高 33.0°C，最低 24.5°C
  预报：气压 1010.5→1002.0 hPa，注意降温和降雨
  热指数：体感 45.7°C

=== 运行时摘掉预报板 ===
  摘掉前订阅者数量：4
  摘掉后订阅者数量：3
  当前状况：气温 31.0°C，湿度 68.0%
  统计（5 次）：平均 28.6°C，最高 33.0°C，最低 24.5°C
```

注意"热指数"那段：`WeatherData` 的源码里**没有一处提到 HeatIndex**，它就是被挂上去的。

## 背后的设计原则

1. **为交互对象之间的松耦合设计而努力。**（本章新增原则）主题不知道观察者的具体类型，
   观察者也不知道彼此；双方只认同一个窄接口。
2. **针对接口编程** —— `observers []Observer`、构造函数参数 `s Subject`。
3. **封装变化** —— 会变的是"谁订阅、订阅几个"，这部分被抽成了列表；不变的是"数据变了要通知"。

一句话收益：**主题和观察者可以各自独立复用、独立演化。**

## Go 特有的注意点

### 1. 为什么把三个 float 换成一个 `Measurement` 结构体

书上的推模型是 `update(float temperature, float humidity, float pressure)`。这在 Java 里没问题，
但在 Go 里有两个现实问题：

- 加一个字段（比如风向）就要改**接口签名 + 所有实现**，而这三者往往在不同包；
- 三个同类型的相邻 `float64` 参数极易传错顺序，编译器完全帮不上忙。

改成 `Update(Measurement)` 后，加字段只影响需要它的那个观察者，其它人一行不改。
代价：观察者拿到全部数据，哪怕只用其中一个。（这也是"推模型 vs 拉模型"的经典取舍。）

### 2. 推模型 vs 拉模型：本实现同时支持两种

- **推（push）**：`Update(m Measurement)` 直接把数据递过去 —— 本目录的四个显示板都用这个。
- **拉（pull）**：`Update()` 不带参数，观察者自己调用 `subject.Measurement()` 去取。
  测试 `TestObserverCanCallSubjectInsideUpdate` 里的 `pullingObserver` 就是这么干的。

书里说拉模型"被认为更正确"（观察者不需要知道主题的细节），但实际上拉模型更灵活：
观察者可以只取自己需要的字段，主题也不用担心"要推什么"。**本实现把两种都留了口子**：
`Update` 带数据（推），`Measurement()` 是公开 getter（拉）。

### 3. `RemoveObserver` 用 `==` 比较接口值 —— 有坑

```go
for i, existing := range w.observers {
    if existing == o { ... }
}
```

Go 规范原文（[Comparison operators](https://go.dev/ref/spec#Comparison_operators)）：

> "A comparison of two interface values with identical dynamic types causes a **run-time panic**
>  if that type is not comparable. … Slice, map, and function types are not comparable."

也就是说：如果某个观察者的**动态类型不可比较**（结构体里含 slice / map / func 字段），
这一行会在运行时 panic，而不是编译报错。本目录所有观察者都用指针（`*XxxDisplay`）就没问题；
如果你见过 `type myObs struct{ events []string }` 这种值类型观察者，`go test -race` 能跑过，
但一调用 `RemoveObserver` 就炸。

真要在 Go 里稳妥地做这件事，有两条更常用的路：

```go
// 路子 A：注册时发一张"退订凭证"（像 context.CancelFunc / signal.Notify 的 stop）
func (w *WeatherData) Subscribe(fn func(Measurement)) (cancel func())

// 路子 B：用 ID / key 而不是对象身份
func (w *WeatherData) RegisterObserver(id string, o Observer)
func (w *WeatherData) RemoveObserver(id string)
```

### 4. 注册两次会通知两次

我写测试时就踩了这个坑：显示板的构造函数**自己会注册**，我又在外面
`station.RegisterObserver(display)` 注册了一次，结果同一块板子收到两次通知、`ObserverCount` 也多算一个。
书上的 Java 版也是"构造即注册"（`weatherData.registerObserver(this)`），所以这不是 Go 特有问题，
但 Go 里 `==` 只比指针、没有 Java 的 `equals`，更容易让你误以为"重复注册会被自动去重"——**不会**。
需要去重就自己在 `RegisterObserver` 里查一遍（同样要小心里面的 `==` 比较）。

### 5. 并发：书里没讲，Go 里必须讲

`WeatherData` 是可变状态，多个 goroutine 里注册/更新就会数据竞争。本实现做了两件事：

```go
w.mu.RLock()
snapshot := make([]Observer, len(w.observers))
copy(snapshot, w.observers)   // ① 锁内做快照
m := w.current
w.mu.RUnlock()

for _, o := range snapshot {  // ② 锁外做回调
    o.Update(m)
}
```

- **① 快照**：观察者在 `Update` 里注销自己/注册新观察者，都不会破坏正在进行的遍历
  （测试 `TestObserverCanUnregisterItselfDuringUpdate` 覆盖了这一点）。
- **② 锁外回调**：如果持锁调用 `Update`，观察者里再调用 `SetMeasurements` 或 `RemoveObserver`
  就会**自我死锁**（Go 的 `sync.Mutex` 不可重入）。这条是 Go 里最常见的观察者事故。

`go test -race ./02-observer` 通过，说明这层护栏是有效的。

### 6. 别把 `Observer` 接口一上来就放进"被观察者"包里

Go 官方 Code Review Comments 的建议是：**接口定义在使用它的包里**。这里 `Observer` 被
`WeatherData` 调用，所以放在同一个包是对的；如果以后气象站放进 `weather` 包、显示板放进 `displays` 包，
那么 `Observer` 应该定义在 `weather` 包（使用方），而不是 `displays` 包（实现方）。

## 什么时候别用

- **只有 1 个确定的通知对象**：直接调一个方法更清楚，别为"将来可能有第二个观察者"上模式。
- **需要背压 / 异步 / 多消费者队列**：Go 里通常直接用 channel（见下）。
- **通知是"命令式"的强耦合业务流**（A 必须等 B 处理完才继续）：这是调用链，不是观察者。
- **观察者链会形成环**（A 通知 B，B 又改 A）：会产生无限递归，此时该用事件队列或显式去重。

## 延伸：Go 里更常见的替代写法 —— channel

Go 程序员写"通知"时，第一反应常常不是接口，而是 channel：

```go
type WeatherData struct {
    updates chan Measurement   // 或者叫 subscribers chan Observer
}

func (w *WeatherData) Subscribe() <-chan Measurement { ... }
func (w *WeatherData) SetMeasurements(m Measurement) {
    select { case w.updates <- m: default: /* 丢弃或阻塞 */ }
}
```

对比一下，选型就有依据了：

| | 本文的接口版（同步回调） | channel 版 |
|---|---|---|
| 通知时机 | 立刻、同步，`NotifyObservers` 返回即全部处理完 | 异步，消费者自己何时取由调度决定 |
| 背压 | 天然有（同步即背压） | 需要缓冲/丢弃/阻塞策略 |
| 顺序 | 注册顺序稳定 | 多个消费者之间不保证 |
| 差错传播 | 观察者自己处理 | 需要额外 error channel |
| 适用 | 显示板、缓存失效、指标上报 | 流水线、工作队列、事件总线 |

两者不是替代关系：**channel 适合"数据流"，Observer 模式适合"状态变化广播"**。

## 进一步阅读

- 《Head First 设计模式》第 2 版 第 2 章（Weather-O-Rama）
- Go 规范 · Comparison operators（接口值比较何时 panic）: https://go.dev/ref/spec#Comparison_operators
- Go Code Review Comments · Interfaces: https://go.dev/wiki/CodeReviewComments#interfaces
- Go 官方博客 · Share Memory By Communicating（channel vs 回调的思路差异）: https://go.dev/blog/codelab-share
- NWS · The Heat Index Equation（本目录 `heatIndexCelsius` 用的 Rothfusz 回归公式）: https://www.wpc.ncep.noaa.gov/html/heatindex_equation.shtml
