package observer

import "fmt"

// ============================================================================
// 观察者模式的两个合同（书上是 Subject / Observer / DisplayElement 三个接口）
//
// 一句话意图：定义对象间的一对多依赖，当一个对象改变状态时，所有依赖它的对象
// 都会自动收到通知并更新。
//
// 书里的三个角色 → 本文件的三个类型：
//   Subject       发布者："我有变化会通知你"          → 注册/注销/通知
//   Observer      订阅者："有变化就调我的 Update"      → 只要一个方法
//   DisplayElement 展示："我能把自己显示出来"          → 只要一个方法
//
// Go 里这三个接口都很窄（1~3 个方法），这是有意的：
// 接口越窄，能实现它的东西越多，替换成本越低。
// ============================================================================

// Measurement 是"会变化的数据"的载体。
//
// 书上的推模型（push）签名是 update(float temp, float humidity, float pressure)。
// 这里改成传一个结构体，原因见 README「为什么我把三个 float 换成一个结构体」：
// 加一个字段（比如风向）时，不用改 Observer 接口、也不用改所有实现。
type Measurement struct {
	Temperature float64 // 气温，摄氏度
	Humidity    float64 // 相对湿度，百分比 0~100
	Pressure    float64 // 气压，hPa
}

func (m Measurement) String() string {
	return fmt.Sprintf("%.1f°C / 湿度 %.1f%% / %.1f hPa", m.Temperature, m.Humidity, m.Pressure)
}

// Observer 是订阅者合同。只有一个方法 —— 任何类型都能实现它。
type Observer interface {
	// Update 由主题在数据变化时调用。参数是本次的完整测量值（推模型）。
	Update(Measurement)
}

// DisplayElement 是展示合同。书上把它单列，是为了强调"观察者不必等于显示元素"。
type DisplayElement interface {
	Display() string
}

// Subject 是发布者合同。客户端只依赖这个抽象，不关心背后是气象站还是股票行情。
type Subject interface {
	RegisterObserver(o Observer)
	RemoveObserver(o Observer)
	NotifyObservers()
}
