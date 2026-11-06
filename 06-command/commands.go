package command

// ============================================================================
// 具体命令（Concrete Command）
//
// 每个命令只有三个要素：
//   ① 一个接收者指针（Light / GarageDoor / CeilingFan / Stereo）
//   ② Execute()：调用接收者的方法
//   ③ Undo()：把接收者恢复成执行前的样子
//
// 书上每种动作一个类（LightOnCommand、LightOffCommand...）。
// 这里基本照做，唯一的合并是吊扇的四个速度：
// 书里 CeilingFanHighCommand / MediumCommand / LowCommand / OffCommand 四个类
// 只差一个常量，Go 里当然可以用一个结构体 + 参数，再配四个构造函数。
// **注意：合并的是类型，不是实例** —— 每个速度仍然是独立的命令对象，
// 照样能分别塞进不同插槽、分别被撤销。
// ============================================================================

// ---------------------------------------------------------------------------
// 灯
// ---------------------------------------------------------------------------

type LightOnCommand struct{ light *Light }

func NewLightOnCommand(light *Light) *LightOnCommand { return &LightOnCommand{light: light} }

func (c *LightOnCommand) Execute() string { return c.light.On() }
func (c *LightOnCommand) Undo() string    { return c.light.Off() }

type LightOffCommand struct{ light *Light }

func NewLightOffCommand(light *Light) *LightOffCommand { return &LightOffCommand{light: light} }

func (c *LightOffCommand) Execute() string { return c.light.Off() }
func (c *LightOffCommand) Undo() string    { return c.light.On() }

// ---------------------------------------------------------------------------
// 车库门
// ---------------------------------------------------------------------------

type GarageDoorUpCommand struct{ door *GarageDoor }

func NewGarageDoorUpCommand(door *GarageDoor) *GarageDoorUpCommand {
	return &GarageDoorUpCommand{door: door}
}

func (c *GarageDoorUpCommand) Execute() string { return c.door.Up() }
func (c *GarageDoorUpCommand) Undo() string    { return c.door.Down() }

type GarageDoorDownCommand struct{ door *GarageDoor }

func NewGarageDoorDownCommand(door *GarageDoor) *GarageDoorDownCommand {
	return &GarageDoorDownCommand{door: door}
}

func (c *GarageDoorDownCommand) Execute() string { return c.door.Down() }
func (c *GarageDoorDownCommand) Undo() string    { return c.door.Up() }

// ---------------------------------------------------------------------------
// 吊扇：撤销必须"恢复上一次的速度"，所以命令自己得记住 prevSpeed
// ---------------------------------------------------------------------------

// CeilingFanSpeedCommand 把吊扇调到某个速度，并记住调之前的速度。
type CeilingFanSpeedCommand struct {
	fan  *CeilingFan
	set  int // 本次要调到的速度
	prev int // 执行前保存的旧速度 —— 这就是"命令对象携带状态"
}

func (c *CeilingFanSpeedCommand) Execute() string {
	c.prev = c.fan.Speed()
	return c.fan.SetSpeed(c.set)
}

func (c *CeilingFanSpeedCommand) Undo() string {
	return c.fan.SetSpeed(c.prev)
}

func NewCeilingFanHighCommand(fan *CeilingFan) *CeilingFanSpeedCommand {
	return &CeilingFanSpeedCommand{fan: fan, set: FanHigh}
}

func NewCeilingFanMediumCommand(fan *CeilingFan) *CeilingFanSpeedCommand {
	return &CeilingFanSpeedCommand{fan: fan, set: FanMedium}
}

func NewCeilingFanLowCommand(fan *CeilingFan) *CeilingFanSpeedCommand {
	return &CeilingFanSpeedCommand{fan: fan, set: FanLow}
}

func NewCeilingFanOffCommand(fan *CeilingFan) *CeilingFanSpeedCommand {
	return &CeilingFanSpeedCommand{fan: fan, set: FanOff}
}

// ---------------------------------------------------------------------------
// 音响
// ---------------------------------------------------------------------------

type StereoOnCommand struct{ stereo *Stereo }

func NewStereoOnCommand(stereo *Stereo) *StereoOnCommand { return &StereoOnCommand{stereo: stereo} }

func (c *StereoOnCommand) Execute() string { return c.stereo.On() }
func (c *StereoOnCommand) Undo() string    { return c.stereo.Off() }

// StereoOnWithCDCommand 演示"命令可以携带参数"：
// 同样是"打开音响"，但顺带把输入切到 CD —— 这个参数属于命令，不属于接收者。
type StereoOnWithCDCommand struct{ stereo *Stereo }

func NewStereoOnWithCDCommand(stereo *Stereo) *StereoOnWithCDCommand {
	return &StereoOnWithCDCommand{stereo: stereo}
}

func (c *StereoOnWithCDCommand) Execute() string {
	return c.stereo.On() + "，" + c.stereo.SetCD()
}

func (c *StereoOnWithCDCommand) Undo() string { return c.stereo.Off() }

type StereoOffCommand struct{ stereo *Stereo }

func NewStereoOffCommand(stereo *Stereo) *StereoOffCommand {
	return &StereoOffCommand{stereo: stereo}
}

func (c *StereoOffCommand) Execute() string { return c.stereo.Off() }
func (c *StereoOffCommand) Undo() string    { return c.stereo.On() }

// 编译期断言：全部具体命令都满足 Command。
var (
	_ Command = (*LightOnCommand)(nil)
	_ Command = (*LightOffCommand)(nil)
	_ Command = (*GarageDoorUpCommand)(nil)
	_ Command = (*GarageDoorDownCommand)(nil)
	_ Command = (*CeilingFanSpeedCommand)(nil)
	_ Command = (*StereoOnCommand)(nil)
	_ Command = (*StereoOnWithCDCommand)(nil)
	_ Command = (*StereoOffCommand)(nil)
)
