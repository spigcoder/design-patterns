package observer

import "sync"

// ============================================================================
// WeatherData：Subject 的具体实现（书里那个追踪气象数据的对象）
//
// 书上用 java.util.Observable 的继承版 → 后来重写成自己实现 Subject 接口。
// 这里从一开始就是"自己实现"，因为这正是书里最终推荐的形态。
// ============================================================================

// WeatherData 是被观察的主题。
//
// 三个字段类型说明一切：
//   - observers 是 []Observer —— 具体怎么通知，到运行时才知道；
//   - current 是值类型 Measurement —— 传的是快照，观察者改不到我的状态；
//   - mu 是并发安全的护栏（书里没讲，但 Go 里必须讲，见 README）。
type WeatherData struct {
	mu        sync.RWMutex
	observers []Observer
	current   Measurement
}

// NewWeatherData 构造一个还没有观察者的气象站。
func NewWeatherData() *WeatherData { return &WeatherData{} }

var _ Subject = (*WeatherData)(nil)

// RegisterObserver 加入订阅列表。
func (w *WeatherData) RegisterObserver(o Observer) {
	if o == nil {
		return // 防御：nil 观察者会在 NotifyObservers 里 panic
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.observers = append(w.observers, o)
}

// RemoveObserver 退订。
//
// 注意 `existing == o` 这一行：Go 规范说
// "A comparison of two interface values with identical dynamic types causes a
//
//	run-time panic if that type is not comparable."
//
// 也就是说如果某个观察者的动态类型里含 slice/map/func 字段，这里会比较到 panic。
// 观察者一律用指针（*XxxDisplay）就没这个问题。详见 README。
func (w *WeatherData) RemoveObserver(o Observer) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for i, existing := range w.observers {
		if existing == o {
			w.observers = append(w.observers[:i], w.observers[i+1:]...)
			return
		}
	}
}

// NotifyObservers 依次通知所有观察者。
//
// 两个细节都是书里没写、但实战里必须做的：
//  1. 先复制一份快照，再在锁外回调 —— 否则观察者在 Update 里调用
//     SetMeasurements / RegisterObserver 会自己把自己锁死（可重入死锁）；
//  2. 回调期间不持锁 —— 否则观察者想读 Measurement() 也会阻塞。
func (w *WeatherData) NotifyObservers() {
	w.mu.RLock()
	snapshot := make([]Observer, len(w.observers))
	copy(snapshot, w.observers)
	m := w.current
	w.mu.RUnlock()

	for _, o := range snapshot {
		o.Update(m)
	}
}

// SetMeasurements 是书里的 measurementsChanged() 的入口：更新数据 + 通知所有人。
func (w *WeatherData) SetMeasurements(m Measurement) {
	w.mu.Lock()
	w.current = m
	w.mu.Unlock()

	w.NotifyObservers()
}

// Measurement 返回当前值（拉模型用的那个 getter）。
func (w *WeatherData) Measurement() Measurement {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.current
}

// ObserverCount 方便测试断言注册/注销是否生效。
func (w *WeatherData) ObserverCount() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.observers)
}
