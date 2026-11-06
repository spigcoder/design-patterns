package command

import "fmt"

// ============================================================================
// 接收者（Receiver）：真正干活的设备
//
// 书上这些类完全不知道命令模式的存在 —— 它们只有自己的业务方法（On/Off/Up...）。
// 这一点很重要：命令模式**不侵入**接收者，接收者复用它自己的 API。
// ============================================================================

// ---------------------------------------------------------------------------
// Light 灯
// ---------------------------------------------------------------------------

type Light struct {
	location string
	on       bool
}

func NewLight(location string) *Light { return &Light{location: location} }

func (l *Light) On() string {
	l.on = true
	return l.location + "的灯打开了"
}

func (l *Light) Off() string {
	l.on = false
	return l.location + "的灯关上了"
}

func (l *Light) IsOn() bool { return l.on }

// ---------------------------------------------------------------------------
// GarageDoor 车库门
// ---------------------------------------------------------------------------

type GarageDoor struct {
	location string
	open     bool
}

func NewGarageDoor(location string) *GarageDoor { return &GarageDoor{location: location} }

func (g *GarageDoor) Up() string {
	g.open = true
	return g.location + "的车库门升起来了"
}

func (g *GarageDoor) Down() string {
	g.open = false
	return g.location + "的车库门降下去了"
}

func (g *GarageDoor) IsOpen() bool { return g.open }

// ---------------------------------------------------------------------------
// CeilingFan 吊扇（书上用它演示"撤销需要记住之前的状态"）
// ---------------------------------------------------------------------------

const (
	FanOff = iota
	FanLow
	FanMedium
	FanHigh
)

type CeilingFan struct {
	location string
	speed    int
}

func NewCeilingFan(location string) *CeilingFan { return &CeilingFan{location: location} }

func (f *CeilingFan) Speed() int { return f.speed }

func (f *CeilingFan) SetSpeed(speed int) string {
	f.speed = speed
	return fmt.Sprintf("%s的吊扇调到%s", f.location, fanSpeedName(speed))
}

func fanSpeedName(speed int) string {
	switch speed {
	case FanLow:
		return "低速"
	case FanMedium:
		return "中速"
	case FanHigh:
		return "高速"
	default:
		return "关闭"
	}
}

// ---------------------------------------------------------------------------
// Stereo 音响
// ---------------------------------------------------------------------------

type Stereo struct {
	location string
	on       bool
	mode     string
}

func NewStereo(location string) *Stereo { return &Stereo{location: location} }

func (s *Stereo) On() string {
	s.on = true
	return s.location + "的音响打开了"
}

func (s *Stereo) Off() string {
	s.on = false
	s.mode = ""
	return s.location + "的音响关上了"
}

func (s *Stereo) SetCD() string {
	s.mode = "CD"
	return s.location + "的音响切换到 CD 模式"
}

func (s *Stereo) IsOn() bool { return s.on }

func (s *Stereo) Mode() string { return s.mode }
