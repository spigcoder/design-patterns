package command

import "strings"

// ============================================================================
// 命令模式的唯一合同
//
// 书上是：
//   public interface Command {
//       public void execute();
//       public void undo();      // 第 6 章后半段加的
//   }
//
// Go 里一模一样，只是我让两个方法返回 string：
// 书上是 println 打印，这里返回文本便于测试断言；副作用（改设备状态）照样发生。
// ============================================================================

// Command 把"一次请求"包装成对象。
//
// 注意它的形状：**没有参数、没有返回值（除了我们的可测试性返回）**。
// 一个命令对象自己知道：要动谁（接收者）、怎么动、怎么撤回。
// 调用者只知道"我手上有个东西能 Execute"。
type Command interface {
	Execute() string
	Undo() string
}

// ---------------------------------------------------------------------------
// 空对象（Null Object）
//
// 书里用 NoCommand 解决"某个插槽还没配置命令"的情况，避免每次都判断 nil。
// Go 里它甚至更短：空结构体零值即可用，不需要构造函数。
// 这是"空对象模式"最典型的应用，也是命令模式的常见搭档。
// ---------------------------------------------------------------------------

// NoCommand 是一个什么都不做的命令。
type NoCommand struct{}

func (NoCommand) Execute() string { return "" }
func (NoCommand) Undo() string    { return "" }

// ---------------------------------------------------------------------------
// 宏命令（Macro Command）= 命令 + 组合模式
//
// 一个命令里面装着一堆命令：Execute 依次执行，Undo 逆序撤销。
// 书里"派对模式"就是这个：一键把灯、音响、吊扇全打开，再一键全关掉。
// ---------------------------------------------------------------------------

// MacroCommand 一次执行一组命令。
type MacroCommand struct {
	commands []Command
}

func NewMacroCommand(commands ...Command) *MacroCommand {
	return &MacroCommand{commands: commands}
}

func (m *MacroCommand) Execute() string {
	results := make([]string, 0, len(m.commands))
	for _, c := range m.commands {
		results = append(results, c.Execute())
	}
	return strings.Join(results, "；")
}

// Undo 必须**逆序**：后开的先关，否则撤销顺序本身就会变成新的 bug。
func (m *MacroCommand) Undo() string {
	results := make([]string, 0, len(m.commands))
	for i := len(m.commands) - 1; i >= 0; i-- {
		results = append(results, m.commands[i].Undo())
	}
	return strings.Join(results, "；")
}

var (
	_ Command = NoCommand{}
	_ Command = (*MacroCommand)(nil)
)
