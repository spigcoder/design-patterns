// Demo：可编程遥控器（书上的 Remote Control）+ 撤销 + 派对模式。
//
//	go run ./06-command/demo
package main

import (
	"fmt"

	command "github.com/spigcoder/design-patterns/06-command"
)

func main() {
	remote := command.NewRemoteControl(7)

	// 设备（接收者）
	kitchenLight := command.NewLight("厨房")
	livingLight := command.NewLight("客厅")
	garage := command.NewGarageDoor("车库")
	fan := command.NewCeilingFan("客厅")
	stereo := command.NewStereo("客厅")

	// 把命令装进插槽（这一步之后，遥控器就不再需要认识这些设备了）
	must(remote.SetCommand(0, command.NewLightOnCommand(kitchenLight), command.NewLightOffCommand(kitchenLight)))
	must(remote.SetCommand(1, command.NewLightOnCommand(livingLight), command.NewLightOffCommand(livingLight)))
	must(remote.SetCommand(2, command.NewGarageDoorUpCommand(garage), command.NewGarageDoorDownCommand(garage)))
	must(remote.SetCommand(3, command.NewCeilingFanHighCommand(fan), command.NewCeilingFanOffCommand(fan)))
	must(remote.SetCommand(4, command.NewStereoOnWithCDCommand(stereo), command.NewStereoOffCommand(stereo)))

	fmt.Print(remote)

	fmt.Println("\n=== 1. 按按钮：遥控器不知道自己在开什么 ===")
	press(remote.OnButton, 0)
	press(remote.OnButton, 2)
	press(remote.OnButton, 4)

	fmt.Println("\n=== 2. 撤销上一次操作 ===")
	fmt.Println("  撤销 → " + remote.UndoButton())
	fmt.Printf("  音响现在：on=%v\n", stereo.IsOn())

	fmt.Println("\n=== 3. 吊扇：撤销要回到「上一个速度」，而不是「关闭」 ===")
	fmt.Println("  先调到中速 → " + command.NewCeilingFanMediumCommand(fan).Execute())
	press(remote.OnButton, 3) // 高速
	fmt.Printf("  当前速度：%d（3=高速）\n", fan.Speed())
	fmt.Println("  撤销 → " + remote.UndoButton())
	fmt.Printf("  撤销后速度：%d（2=中速，而不是 0=关闭）\n", fan.Speed())

	fmt.Println("\n=== 4. 派对模式：宏命令一键开、一键关 ===")
	partyOn := command.NewMacroCommand(
		command.NewLightOnCommand(livingLight),
		command.NewStereoOnWithCDCommand(stereo),
		command.NewCeilingFanHighCommand(fan),
	)
	partyOff := command.NewMacroCommand(
		command.NewCeilingFanOffCommand(fan),
		command.NewStereoOffCommand(stereo),
		command.NewLightOffCommand(livingLight),
	)
	must(remote.SetCommand(5, partyOn, partyOff))

	press(remote.OnButton, 5)
	press(remote.OffButton, 5)
	fmt.Printf("  派对结束后：灯=%v 音响=%v 吊扇=%d\n", livingLight.IsOn(), stereo.IsOn(), fan.Speed())
}

func press(button func(int) (string, error), slot int) {
	out, err := button(slot)
	if err != nil {
		fmt.Println("  错误：" + err.Error())
		return
	}
	fmt.Printf("  [%d] %s\n", slot, out)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
