package main

// (x): 사용자의 총 배터리 용량을 가져온다.
// (x): 사용자의 지금 배터리 충전 %를 가져온다.
// 배터리 사용 시간(시간)은 배터리 용량(mAh) / 기기 전력 소모량(mA)
// 기기 전력 소모량은 전압(V)과 전류(A)를 곱하여 구할 수 있습니다. 공식은 P(전력, W) = E(전압, V) x I(전류, A)
// Power: 현재 전력 소모량 (단위: W, 와트)
// Current: 현재 전류 (단위: A, 암페어)
// Voltage: 현재 전압 (단위: V, 볼트)
// 3000mAh 배터리를 사용하는 기기가 500mA를 소모한다면 3000mAh / 500mA = 6시간
// cap: 용량
// bat: 배터리
// Rem: remaining
// (0): 충전중일 때는, 완충 시간 뜨게 하기
// (0): 배터리 감소 때는 배터리 남은 시간 뜨게 하기

import (
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
	var cyan = color.New(color.FgCyan, color.Bold).PrintfFunc() // cyan colored printf

	batteries, err := battery.GetAll()
	if err != nil {
		fmt.Println("Could not get battery info!")
		return
	}

	var batteryData = batteries[0] // batteries: 구조체
	// var batState = batteryData.State
	var curBatCapMWh float64 = batteryData.Current // 현재 충전된  배터리 용량 (mWh)
	//var fullBatCapMWh = batteryData.Full           // 최대 가용 가능 배터리 용량 (mWh)
	var powerMW float64 = batteryData.ChargeRate // 현재 전력 소모량 (mW)
	var voltageV float64 = batteryData.Voltage
	var curBatCapMAh float64 = curBatCapMWh / voltageV // 현재 충전된  배터리 용량 (mAh)
	//var fullBatCapMAh float64 = fullBatCapMWh / voltageV // 최대 가용 가능 배터리 용량 (mWh)
	var powerMA float64 = powerMW / voltageV // 현재 전력 소모량 (mA)

	timeRem := curBatCapMAh / powerMA // 추정 배터리 잔여 시간
	hour := int(timeRem)
	minutes := int(60 * (float64(getFirstDecimalDigit(timeRem)) / 10))
	// 예상 남은 배터리 시간
	fmt.Println("Remaining Battery Time")
	cyan("%d : %d\n", hour, minutes)
}
