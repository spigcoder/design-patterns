package command_test

import (
	"strings"
	"testing"

	command "github.com/spigcoder/design-patterns/06-command"
)

// 场景 1：遥控器不认识任何设备，只认识 Command 接口。
// 同一套按钮逻辑，装灯也能用、装音响也能用。
func TestRemoteWorksWithAnyDevice(t *testing.T) {
	remote := command.NewRemoteControl(3)
	light := command.NewLight("客厅")
	stereo := command.NewStereo("客厅")

	if err := remote.SetCommand(0,
		command.NewLightOnCommand(light), command.NewLightOffCommand(light)); err != nil {
		t.Fatal(err)
	}
	if err := remote.SetCommand(1,
		command.NewStereoOnWithCDCommand(stereo), command.NewStereoOffCommand(stereo)); err != nil {
		t.Fatal(err)
	}

	got, err := remote.OnButton(0)
	if err != nil || !strings.Contains(got, "客厅的灯打开了") {
		t.Fatalf("按 0 号开关键: %q, %v", got, err)
	}

	got, err = remote.OnButton(1)
	if err != nil || !strings.Contains(got, "CD 模式") {
		t.Fatalf("按 1 号开关键: %q, %v", got, err)
	}

	if !light.IsOn() || !stereo.IsOn() {
		t.Error("设备状态没有被命令改动")
	}
}

// 场景 2：撤销。命令自己记住了"执行前该怎样"，遥控器完全不知道细节。
func TestUndoRestoresPreviousState(t *testing.T) {
	remote := command.NewRemoteControl(1)
	light := command.NewLight("厨房")
	if err := remote.SetCommand(0,
		command.NewLightOnCommand(light), command.NewLightOffCommand(light)); err != nil {
		t.Fatal(err)
	}

	if _, err := remote.OnButton(0); err != nil {
		t.Fatal(err)
	}
	if !light.IsOn() {
		t.Fatal("灯没打开")
	}

	remote.UndoButton()
	if light.IsOn() {
		t.Error("撤销后灯应该回到关着的状态")
	}
}

// 场景 3：吊扇的撤销需要恢复"上一次的速度"——所以命令对象必须携带状态。
func TestCeilingFanUndoRestoresPreviousSpeed(t *testing.T) {
	remote := command.NewRemoteControl(2)
	fan := command.NewCeilingFan("卧室")

	low := command.NewCeilingFanLowCommand(fan)
	high := command.NewCeilingFanHighCommand(fan)
	off := command.NewCeilingFanOffCommand(fan)

	_ = remote.SetCommand(0, low, off)
	_ = remote.SetCommand(1, high, off)

	if _, err := remote.OnButton(0); err != nil { // 低速
		t.Fatal(err)
	}
	if fan.Speed() != command.FanLow {
		t.Fatalf("期望低速，实际 %d", fan.Speed())
	}

	if _, err := remote.OnButton(1); err != nil { // 高速（记住之前是低速）
		t.Fatal(err)
	}
	if fan.Speed() != command.FanHigh {
		t.Fatalf("期望高速，实际 %d", fan.Speed())
	}

	remote.UndoButton() // 撤销高速 → 应回到低速，而不是"关闭"
	if fan.Speed() != command.FanLow {
		t.Errorf("撤销高速后应回到低速，实际 %d", fan.Speed())
	}
}

// 场景 4：没配置的插槽按下去也不出事（NoCommand 空对象）。
func TestUnconfiguredSlotIsSafe(t *testing.T) {
	remote := command.NewRemoteControl(2)

	got, err := remote.OnButton(1)
	if err != nil {
		t.Fatalf("未配置插槽不该报错: %v", err)
	}
	if got != "" {
		t.Errorf("未配置插槽不该有输出，得到 %q", got)
	}
	if got := remote.UndoButton(); got != "" {
		t.Errorf("初始撤销不该有输出，得到 %q", got)
	}
}

// 场景 5：宏命令（派对模式）—— 一次执行一组，撤销时逆序。
func TestMacroCommand(t *testing.T) {
	light := command.NewLight("客厅")
	stereo := command.NewStereo("客厅")
	fan := command.NewCeilingFan("客厅")

	partyOn := command.NewMacroCommand(
		command.NewLightOnCommand(light),
		command.NewStereoOnWithCDCommand(stereo),
		command.NewCeilingFanHighCommand(fan),
	)

	out := partyOn.Execute()
	for _, want := range []string{"灯打开了", "CD 模式", "高速"} {
		if !strings.Contains(out, want) {
			t.Errorf("宏命令输出缺少 %q: %q", want, out)
		}
	}

	partyOn.Undo()
	if light.IsOn() || stereo.IsOn() || fan.Speed() != command.FanOff {
		t.Error("宏命令撤销后设备应全部复位")
	}
}

// 场景 6：命令是**对象**，所以可以被存起来、排队、稍后执行、记录日志。
// 这是它和"策略"最本质的区别（策略是替换某种行为，不需要被排队）。
func TestCommandsCanBeQueuedAndReplayed(t *testing.T) {
	light := command.NewLight("书房")

	// 先把三个请求装进队列（此时什么都没发生）
	queue := []command.Command{
		command.NewLightOnCommand(light),
		command.NewLightOffCommand(light),
		command.NewLightOnCommand(light),
	}
	if light.IsOn() {
		t.Fatal("还没执行，灯不该亮")
	}

	// 稍后统一执行
	var log []string
	for _, c := range queue {
		log = append(log, c.Execute())
	}

	if !light.IsOn() {
		t.Error("队列执行完，灯应该是亮的")
	}
	if len(log) != 3 {
		t.Errorf("应该记录了 3 条执行日志，实际 %d", len(log))
	}
}

// 场景 7：自定义一个只在本测试里存在的命令，遥控器照样能用 —— 不需要改遥控器。
func TestCustomCommandNeedsNoInvokerChange(t *testing.T) {
	sprinkler := &sprinkler{}
	remote := command.NewRemoteControl(1)
	if err := remote.SetCommand(0,
		command.NewLightOnCommand(command.NewLight("阳台")), // 占位
		command.NewLightOffCommand(command.NewLight("阳台"))); err != nil {
		t.Fatal(err)
	}
	// 直接塞自定义命令
	if err := remote.SetCommand(0, sprinkler, sprinkler); err != nil {
		t.Fatal(err)
	}

	got, err := remote.OnButton(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != "洒水器开始喷水" {
		t.Errorf("自定义命令没被调用: %q", got)
	}
}

type sprinkler struct{ on bool }

func (s *sprinkler) Execute() string { s.on = true; return "洒水器开始喷水" }
func (s *sprinkler) Undo() string    { s.on = false; return "洒水器停止喷水" }

// 场景 8：越界插槽返回 error，而不是 panic。
func TestInvalidSlotReturnsError(t *testing.T) {
	remote := command.NewRemoteControl(2)
	if err := remote.SetCommand(9, nil, nil); err == nil {
		t.Error("越界插槽应该报错")
	}
	if _, err := remote.OnButton(9); err == nil {
		t.Error("越界插槽应该报错")
	}
}
