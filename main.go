package main

import (
	"flag"
	"fmt"
	"math"

	"github.com/distatus/battery"
	"github.com/fatih/color"
)

func getFirstDecimalDigit(n float64) int {
	truncated := math.Trunc(n*10) / 10
	firstDecimal := math.Mod(truncated*10, 10) // 3.9 -> 9
	return int(firstDecimal)
}

func main() {
	// colored printf
	var boldCyan = color.New(color.FgCyan, color.Bold).PrintfFunc() // cyan colored printf
	var boldRed = color.New(color.FgRed, color.Bold).PrintfFunc()
	var underlineBoldWhite = color.New(color.FgWhite, color.Bold, color.Underline).PrintfFunc()

	// flags
	var iFlag = flag.Bool("i", false, "battery informations")
	var infoFlag = flag.Bool("info", false, "battery informations")
	flag.Parse() // flag init

	// battery
	batteries, err := battery.GetAll()
	if err != nil {
		fmt.Println("error: Could not get battery info!")
		return
	}

	// 데스크탑 감지
	if len(batteries) == 0 {
		fmt.Println("This system does not support batteries.")
		return
	}

	var batteryData = batteries[0]                 // batteries: 구조체
	var batState = batteryData.State               // Charging, Discharging
	var curBatCapMWh float64 = batteryData.Current // 현재 충전된  배터리 용량 (mWh)
	var fullBatCapMWh float64 = batteryData.Full   // 최대 가용 가능 배터리 용량 (mWh)
	var designBatCapMWh float64 = batteryData.Design
	var powerMW float64 = batteryData.ChargeRate // 현재 전력 소모량 (mW)
	var voltageV float64 = batteryData.Voltage
	var designVoltageV float64 = batteryData.DesignVoltage

	// -i flag
	if *iFlag || *infoFlag {
		underlineBoldWhite("Battery Information:\n")
		fmt.Printf("state: %s\n", batState.String())
		fmt.Printf("current capacity: %f mWh\n", curBatCapMWh)
		fmt.Printf("last full capacity: %f mWh\n", fullBatCapMWh)
		fmt.Printf("design capacity: %f mWh\n", designBatCapMWh)
		fmt.Printf("charge rate: %f mW\n", powerMW)
		fmt.Printf("voltage: %f V\n", voltageV)
		fmt.Printf("design voltage: %f V\n", designVoltageV)
	}

	// 인수가 없다면
	if flag.NFlag() == 0 {
		switch batState.String() {
		case battery.Charging.String():
			// 예상 남은 충전 시간
			remainingCapMWh := fullBatCapMWh - curBatCapMWh
			timeRem := remainingCapMWh / powerMW
			hour := int(timeRem)
			minutes := int(60 * (float64(getFirstDecimalDigit(timeRem)) / 10))

			fmt.Println("Battery Charge Time")
			boldRed("%d : %d\n", hour, minutes)
		default:
			// 예상 남은 배터리 시간
			timeRem := curBatCapMWh / powerMW // 추정 배터리 잔여 시간
			hour := int(timeRem)
			minutes := int(60 * (float64(getFirstDecimalDigit(timeRem)) / 10))

			fmt.Println("Remaining Battery Time")
			boldCyan("%d : %d\n", hour, minutes)
		}
	}
}
