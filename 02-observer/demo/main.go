// Demo：跑一遍《Head First 设计模式》第 2 章的气象站（Weather-O-Rama）。
//
//	go run ./02-observer/demo
package main

import (
	"fmt"

	observer "github.com/spigcoder/design-patterns/02-observer"
)

func main() {
	station := observer.NewWeatherData()

	// 三块显示板自己把自己注册到气象站上 —— 注意 main 里没有一行"分发"代码。
	current := observer.NewCurrentConditionsDisplay(station)
	stats := observer.NewStatisticsDisplay(station)
	forecast := observer.NewForecastDisplay(station)

	// 只认识 DisplayElement 接口，不关心具体是哪块板子。
	displays := []observer.DisplayElement{current, stats, forecast}

	show := func(label string, m observer.Measurement) {
		fmt.Printf("\n=== %s ===\n", label)
		station.SetMeasurements(m)
		for _, d := range displays {
			fmt.Println("  " + d.Display())
		}
	}

	show("第一次测量", observer.Measurement{Temperature: 26.5, Humidity: 65, Pressure: 1013.2})
	show("第二次测量（气压下降）", observer.Measurement{Temperature: 28.0, Humidity: 70, Pressure: 1005.0})
	show("第三次测量（气压回升）", observer.Measurement{Temperature: 24.5, Humidity: 62, Pressure: 1010.5})

	fmt.Println("\n=== 运行时挂上一块新显示板（热指数），主题无需改动 ===")
	displays = append(displays, observer.NewHeatIndexDisplay(station))
	show("第四次测量（湿热，气压又降）", observer.Measurement{Temperature: 33.0, Humidity: 75, Pressure: 1002.0})

	fmt.Println("\n=== 运行时摘掉预报板 ===")
	fmt.Printf("  摘掉前订阅者数量：%d\n", station.ObserverCount())
	station.RemoveObserver(forecast)
	fmt.Printf("  摘掉后订阅者数量：%d\n", station.ObserverCount())
	station.SetMeasurements(observer.Measurement{Temperature: 31.0, Humidity: 68, Pressure: 1004.0})
	for _, d := range []observer.DisplayElement{current, stats} {
		fmt.Println("  " + d.Display())
	}
}
