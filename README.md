# 设计模式 · Go 实现（对照《Head First 设计模式》第 2 版）

每学完一章就新增一个目录，用 Go 重写书里的例子，并记录"这个模式在 Go 里会长成什么样"。

## 目录

| 章 | 模式 | 目录 | 状态 |
|---|---|---|---|
| 1 | Strategy 策略 | [`01-strategy`](./01-strategy) | ✅ |
| 2 | Observer 观察者 | [`02-observer`](./02-observer) | ✅ |
| 3 | Decorator 装饰者 | `03-decorator` | ⬜ |
| 4 | Factory Method / Abstract Factory 工厂 | `04-factory` | ⬜ |
| 5 | Singleton 单件 | `05-singleton` | ⬜ |
| 6 | Command 命令 | `06-command` | ⬜ |
| 7 | Adapter / Facade 适配器与外观 | `07-adapter-facade` | ⬜ |
| 8 | Template Method 模板方法 | `08-template-method` | ⬜ |
| 9 | Iterator / Composite 迭代器与组合 | `09-iterator-composite` | ⬜ |
| 10 | State 状态 | `10-state` | ⬜ |
| 11 | Proxy 代理 | `11-proxy` | ⬜ |
| 12 | 复合模式（鸭子模拟器） | `12-compound` | ⬜ |
| 13 | 真实世界中的模式（MVC） | `13-mvc` | ⬜ |
| 14 | 附录：剩余模式（Bridge/Builder/CoR/Flyweight/Interpreter/Mediator/Memento/Prototype/Visitor） | `14-leftover` | ⬜ |

## 目录约定

- 一个模式一个目录，前缀是章号，方便顺序阅读；目录名用英文小写加连字符。
- Go 的包名不能用数字开头，所以目录 `01-strategy` 里的包名叫 `strategy`。
- 每个模式目录包含：
  - `*.go` —— 模式的实现，包名 = 模式名
  - `*_test.go` —— 用测试表达"这个模式解决了什么问题"（可替换性、运行时可换、扩展时不改旧代码）
  - `demo/main.go` —— 可运行的场景演示，**一个模式一个 main 包**，避免包名冲突
  - `README.md` —— 模式卡片，固定结构见下
- 只有一个 `go.mod`（模块 `github.com/spigcoder/design-patterns`），所有模式同一个模块。
- 导入本模块内的包时，导入路径形如 `github.com/spigcoder/design-patterns/01-strategy`。

## 常用命令

```bash
go test ./...                 # 跑所有模式的测试
go test -race ./...           # 带竞态检测（观察者那章开始必须跑）
go run ./01-strategy/demo     # 跑某个模式的演示
go vet ./... && gofmt -l .    # 体检
```

## 模式卡片模板（每个模式目录的 README）

1. **一句话意图** —— 这个模式干了什么（尽量背下 GoF 原句的意译）
2. **它解决的问题** —— 书里的故事 + "如果不这么写会怎样"的具体痛点
3. **参与者 → Go 里的对应物** —— 表格：书上 Java 概念 ↔ Go 写法
4. **结构图** —— 文本 UML
5. **可运行演示** —— 命令 + 真实输出
6. **背后的设计原则** —— 记下这章原则，累积成自己的清单
7. **Go 特有的注意点** —— 语言差异导致的写法差别、常见误用
8. **什么时候别用** —— 反例比正例更能防滥用
9. **进一步阅读**

## 设计原则累积清单

- 封装变化：找出应用中可能需要变化之处，把它们独立出来，不要和那些不需要变化的代码混在一起。（Ch.1）
- 针对接口编程，而不是针对实现编程。（Ch.1）
- 多用组合，少用继承。（Ch.1）
- 为交互对象之间的松耦合设计而努力。（Ch.2）
