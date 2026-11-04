package observer_test

import (
	"strings"
	"testing"

	observer "github.com/spigcoder/design-patterns/02-observer"
)

// 场景 1：主题一变，所有注册过的观察者都收到通知。
func TestAllObserversAreNotified(t *testing.T) {
	station := observer.NewWeatherData()
	cur := observer.NewCurrentConditionsDisplay(station)
	stats := observer.NewStatisticsDisplay(station)
	forecast := observer.NewForecastDisplay(station)

	station.SetMeasurements(observer.Measurement{Temperature: 26.0, Humidity: 65, Pressure: 1013.2})
	station.SetMeasurements(observer.Measurement{Temperature: 28.0, Humidity: 70, Pressure: 1005.0})

	if got := cur.Display(); !strings.Contains(got, "28.0°C") {
		t.Errorf("当前状况没更新: %q", got)
	}
	if got := stats.Display(); !strings.Contains(got, "平均 27.0°C") {
		t.Errorf("统计不对: %q", got)
	}
	if got := forecast.Display(); !strings.Contains(got, "1005.0 hPa") {
		t.Errorf("预报没更新: %q", got)
	}
}

// 场景 2：退订之后不再收到通知。
func TestRemoveObserverStopsNotifications(t *testing.T) {
	station := observer.NewWeatherData()
	cur := observer.NewCurrentConditionsDisplay(station)
	station.SetMeasurements(observer.Measurement{Temperature: 20, Humidity: 50, Pressure: 1010})

	station.RemoveObserver(cur)
	if n := station.ObserverCount(); n != 0 {
		t.Fatalf("注销后还剩 %d 个观察者", n)
	}

	station.SetMeasurements(observer.Measurement{Temperature: 35, Humidity: 90, Pressure: 990})
	if got := cur.Display(); !strings.Contains(got, "20.0°C") {
		t.Errorf("已退订却收到了新数据: %q", got)
	}
}

// 场景 3：统计显示板自己的状态是对的（观察者可以持有状态）。
func TestStatisticsTrackMinMaxAvg(t *testing.T) {
	station := observer.NewWeatherData()
	stats := observer.NewStatisticsDisplay(station)

	for _, temp := range []float64{10, 20, 30} {
		station.SetMeasurements(observer.Measurement{Temperature: temp, Humidity: 50, Pressure: 1010})
	}

	got := stats.Display()
	for _, want := range []string{"平均 20.0°C", "最高 30.0°C", "最低 10.0°C", "3 次"} {
		if !strings.Contains(got, want) {
			t.Errorf("统计缺少 %q: %q", want, got)
		}
	}
}

// 场景 4：观察者在自己的 Update 里调用主题 -> 不能死锁。
// 这验证了 NotifyObservers 是"锁内快照、锁外回调"。
func TestObserverCanCallSubjectInsideUpdate(t *testing.T) {
	station := observer.NewWeatherData()
	probe := &pullingObserver{station: station}
	station.RegisterObserver(probe)

	station.SetMeasurements(observer.Measurement{Temperature: 18, Humidity: 60, Pressure: 1008})

	if probe.pulled.Temperature != 18 {
		t.Errorf("观察者在 Update 里拉取到了错误数据: %+v", probe.pulled)
	}
}

// pullingObserver 演示"拉模型"：Update 的参数不用，自己去主题要数据。
type pullingObserver struct {
	station observer.Subject
	pulled  observer.Measurement
}

func (o *pullingObserver) Update(observer.Measurement) {
	if s, ok := o.station.(interface {
		Measurement() observer.Measurement
	}); ok {
		o.pulled = s.Measurement()
	}
}

// 场景 5：观察者在自己被通知时注销自己（"用完就走"）—— 也不能出问题。
func TestObserverCanUnregisterItselfDuringUpdate(t *testing.T) {
	station := observer.NewWeatherData()
	once := &selfRemovingObserver{station: station}
	station.RegisterObserver(once)
	observer.NewCurrentConditionsDisplay(station) // 构造函数里会自己注册一次

	station.SetMeasurements(observer.Measurement{Temperature: 22, Humidity: 55, Pressure: 1012})
	station.SetMeasurements(observer.Measurement{Temperature: 23, Humidity: 55, Pressure: 1011})

	if once.calls != 1 {
		t.Errorf("自注销的观察者被调用了 %d 次，应为 1 次", once.calls)
	}
	if n := station.ObserverCount(); n != 1 {
		t.Errorf("注销后应剩 1 个观察者，实际 %d", n)
	}
}

type selfRemovingObserver struct {
	station observer.Subject
	calls   int
}

func (o *selfRemovingObserver) Update(observer.Measurement) {
	o.calls++
	o.station.RemoveObserver(o)
}

// 场景 6：加一块新显示板，不需要改 WeatherData 一行代码。
// 这个观察者定义在测试包里（外部包），照样能订阅 —— 接口是隐式满足的。
func TestNewObserverNeedsNoSubjectChange(t *testing.T) {
	station := observer.NewWeatherData()
	spy := &spyObserver{}
	station.RegisterObserver(spy)

	station.SetMeasurements(observer.Measurement{Temperature: 30, Humidity: 80, Pressure: 999})

	if spy.count != 1 {
		t.Fatalf("新观察者收到 %d 次通知", spy.count)
	}
	if spy.last.Temperature != 30 {
		t.Errorf("新观察者拿到 %+v", spy.last)
	}
}

type spyObserver struct {
	count int
	last  observer.Measurement
}

func (o *spyObserver) Update(m observer.Measurement) {
	o.count++
	o.last = m
}

// 场景 7：热指数是"派生数据"，说明观察者可以自己做计算，主题不必知道。
func TestHeatIndexIsDerivedFromMeasurement(t *testing.T) {
	station := observer.NewWeatherData()
	hi := observer.NewHeatIndexDisplay(station)

	station.SetMeasurements(observer.Measurement{Temperature: 32, Humidity: 70, Pressure: 1005})

	got := hi.Display()
	if !strings.Contains(got, "体感") {
		t.Errorf("热指数输出异常: %q", got)
	}
	if strings.Contains(got, "体感 32.0°C") {
		t.Errorf("热指数不应该等于气温本身（湿热天气体感更高）: %q", got)
	}
}
