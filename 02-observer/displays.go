package observer

import (
	"fmt"
	"math"
)

// ============================================================================
// 具体观察者：书上的四个显示板
//
// 每个显示板只做两件事：
//   1. 在构造函数里把自己注册到主题（subject）上；
//   2. 实现 Update 保存数据、实现 Display 输出。
//
// 关键：它们谁都不认识别的显示板，也不认识 WeatherData 的具体类型 ——
// 只认识 Subject 接口。这就是"主题和观察者之间松耦合"落到代码上的样子。
// ============================================================================

// CurrentConditionsDisplay 当前状况显示板。
type CurrentConditionsDisplay struct {
	subject Subject
	last    Measurement
}

func NewCurrentConditionsDisplay(s Subject) *CurrentConditionsDisplay {
	d := &CurrentConditionsDisplay{subject: s}
	s.RegisterObserver(d)
	return d
}

func (d *CurrentConditionsDisplay) Update(m Measurement) { d.last = m }

func (d *CurrentConditionsDisplay) Display() string {
	return fmt.Sprintf("当前状况：气温 %.1f°C，湿度 %.1f%%", d.last.Temperature, d.last.Humidity)
}

var (
	_ Observer       = (*CurrentConditionsDisplay)(nil)
	_ DisplayElement = (*CurrentConditionsDisplay)(nil)
)

// StatisticsDisplay 统计显示板：跟踪最低/最高/平均气温。
type StatisticsDisplay struct {
	subject Subject

	min, max, sum float64
	count         int
}

func NewStatisticsDisplay(s Subject) *StatisticsDisplay {
	d := &StatisticsDisplay{subject: s, min: math.MaxFloat64, max: -math.MaxFloat64}
	s.RegisterObserver(d)
	return d
}

func (d *StatisticsDisplay) Update(m Measurement) {
	d.min = math.Min(d.min, m.Temperature)
	d.max = math.Max(d.max, m.Temperature)
	d.sum += m.Temperature
	d.count++
}

func (d *StatisticsDisplay) Display() string {
	if d.count == 0 {
		return "统计：还没有数据"
	}
	return fmt.Sprintf("统计（%d 次）：平均 %.1f°C，最高 %.1f°C，最低 %.1f°C",
		d.count, d.sum/float64(d.count), d.max, d.min)
}

var (
	_ Observer       = (*StatisticsDisplay)(nil)
	_ DisplayElement = (*StatisticsDisplay)(nil)
)

// ForecastDisplay 天气预报显示板：靠“气压和上次相比是升还是降”来推断。
// 它必须自己保存上一次的状态 —— 观察者是可以有状态的。
// 注意：主题并不知道预报需要两次读数，这个“历史”只存在观察者自己肚子里。
type ForecastDisplay struct {
	subject Subject
	last    float64 // 上一次气压
	current float64 // 本次气压
	count   int     // 收到几次读数
}

func NewForecastDisplay(s Subject) *ForecastDisplay {
	d := &ForecastDisplay{subject: s}
	s.RegisterObserver(d)
	return d
}

func (d *ForecastDisplay) Update(m Measurement) {
	d.last = d.current
	d.current = m.Pressure
	d.count++
}

func (d *ForecastDisplay) Display() string {
	switch {
	case d.count == 0:
		return "预报：还没有数据"
	case d.count == 1:
		return fmt.Sprintf("预报：气压 %.1f hPa，等下一次读数再预报", d.current)
	case d.current > d.last:
		return fmt.Sprintf("预报：气压 %.1f→%.1f hPa，天气在好转", d.last, d.current)
	case d.current == d.last:
		return fmt.Sprintf("预报：气压持平 %.1f hPa，天气没变化", d.current)
	default:
		return fmt.Sprintf("预报：气压 %.1f→%.1f hPa，注意降温和降雨", d.last, d.current)
	}
}

var (
	_ Observer       = (*ForecastDisplay)(nil)
	_ DisplayElement = (*ForecastDisplay)(nil)
)

// HeatIndexDisplay 热指数显示板（书上第 2 章的练习题）。
//
// 用的是 NWS（美国国家气象局）的 Rothfusz 回归公式，输入必须是华氏度：
//
//	HI = -42.379 + 2.04901523*T + 10.14333127*RH - 0.22475541*T*RH
//	     - 0.00683783*T² - 0.05481717*RH² + 0.00122874*T²*RH
//	     + 0.00085282*T*RH² - 0.00000199*T²*RH²
//
// 出处见 README「进一步阅读」。这里只演示"再挂一块显示板"不需要动主题。
type HeatIndexDisplay struct {
	subject Subject
	index   float64 // 摄氏度
}

func NewHeatIndexDisplay(s Subject) *HeatIndexDisplay {
	d := &HeatIndexDisplay{subject: s}
	s.RegisterObserver(d)
	return d
}

func (d *HeatIndexDisplay) Update(m Measurement) {
	d.index = heatIndexCelsius(m.Temperature, m.Humidity)
}

func (d *HeatIndexDisplay) Display() string {
	return fmt.Sprintf("热指数：体感 %.1f°C", d.index)
}

// heatIndexCelsius 把摄氏温度 + 相对湿度换算成体感温度（摄氏）。
func heatIndexCelsius(celsius, humidity float64) float64 {
	t := celsius*9/5 + 32 // 公式要求华氏度
	rh := humidity

	f := -42.379 +
		2.04901523*t +
		10.14333127*rh -
		0.22475541*t*rh -
		0.00683783*t*t -
		0.05481717*rh*rh +
		0.00122874*t*t*rh +
		0.00085282*t*rh*rh -
		0.00000199*t*t*rh*rh

	return (f - 32) * 5 / 9
}

var (
	_ Observer       = (*HeatIndexDisplay)(nil)
	_ DisplayElement = (*HeatIndexDisplay)(nil)
)
