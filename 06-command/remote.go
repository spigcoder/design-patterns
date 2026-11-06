package command

import (
	"fmt"
	"strings"
)

// ============================================================================
// 调用者（Invoker）：遥控器
//
// 这是命令模式最值得玩味的地方：遥控器**完全不知道**灯、车库门、吊扇的存在。
// 它只有两排 Command 接口 + 一个撤销槽位。
// 加一台新家电 = 写一个新命令 + 塞进插槽，遥控器代码一行不改。
//
// 对比第 3 章装饰器、第 4 章工厂：那两章解决"对象怎么组"和"对象怎么生"，
// 这里解决的是"请求怎么被传递"。
// ============================================================================

// RemoteControl 是调用者：它只认识 Command 接口。
type RemoteControl struct {
	onCommands  []Command
	offCommands []Command
	undoCommand Command // 书里的"单步撤销"：只记住最后一个执行的命令
}

// NewRemoteControl 造一个有 slots 个插槽的遥控器，每个插槽默认是 NoCommand。
// 这样"某个按钮还没配置"时按下去也不会 panic —— 空对象的价值。
func NewRemoteControl(slots int) *RemoteControl {
	on := make([]Command, slots)
	off := make([]Command, slots)
	for i := range on {
		on[i] = NoCommand{}
		off[i] = NoCommand{}
	}
	return &RemoteControl{onCommands: on, offCommands: off, undoCommand: NoCommand{}}
}

// SetCommand 把一个"开/关"命令对装到指定插槽。
func (r *RemoteControl) SetCommand(slot int, onCommand, offCommand Command) error {
	if slot < 0 || slot >= len(r.onCommands) {
		return fmt.Errorf("插槽 %d 不存在（共 %d 个）", slot, len(r.onCommands))
	}
	r.onCommands[slot] = onCommand
	r.offCommands[slot] = offCommand
	return nil
}

// OnButton 按下"开"按钮。
func (r *RemoteControl) OnButton(slot int) (string, error) {
	if slot < 0 || slot >= len(r.onCommands) {
		return "", fmt.Errorf("插槽 %d 不存在（共 %d 个）", slot, len(r.onCommands))
	}
	result := r.onCommands[slot].Execute()
	r.undoCommand = r.onCommands[slot] // 记住它，供撤销用
	return result, nil
}

// OffButton 按下"关"按钮。
func (r *RemoteControl) OffButton(slot int) (string, error) {
	if slot < 0 || slot >= len(r.offCommands) {
		return "", fmt.Errorf("插槽 %d 不存在（共 %d 个）", slot, len(r.offCommands))
	}
	result := r.offCommands[slot].Execute()
	r.undoCommand = r.offCommands[slot]
	return result, nil
}

// UndoButton 撤销**上一次**按下的按钮。
// 注意：命令里保存的是"执行前的状态"，所以撤销不需要遥控器知道任何细节。
func (r *RemoteControl) UndoButton() string {
	return r.undoCommand.Undo()
}

// String 把遥控器面板打印出来（插槽号 + 命令类型名）。
func (r *RemoteControl) String() string {
	var b strings.Builder
	b.WriteString("------ 遥控器面板 ------\n")
	for i := range r.onCommands {
		fmt.Fprintf(&b, "[%d] %-28s | %s\n", i, commandName(r.onCommands[i]), commandName(r.offCommands[i]))
	}
	fmt.Fprintf(&b, "[撤销] %s\n", commandName(r.undoCommand))
	return b.String()
}

// commandName 只用于打印面板。
// 注意：类型判断只放在"显示"里，绝不放进业务逻辑 —— 那是第 3 章讲过的反模式。
func commandName(c Command) string {
	if _, ok := c.(NoCommand); ok {
		return "（未配置）"
	}
	name := fmt.Sprintf("%T", c)
	return strings.TrimPrefix(strings.TrimPrefix(name, "*"), "command.")
}
