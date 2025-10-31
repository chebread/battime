package main

import (
	"flag"
	"fmt"
	"math"
	"os/exec"   // 명령어 실행
    "runtime"   // OS 확인
    "strings"   // 문자열 처리
	"strconv"

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
	// 인수가 없을 때 실행되는 블록
// 인수가 없을 때 실행되는 블록
if flag.NFlag() == 0 {
	// runtime.GOOS로 현재 OS를 확인
	switch runtime.GOOS {
	case "linux", "darwin", "windows":
		// 지원하는 OS일 경우, OS별 명령어를 실행
		var cmd *exec.Cmd

		if runtime.GOOS == "linux" {
			cmd = exec.Command("sh", "-c", "upower -i $(upower -e | grep 'BAT') | grep 'time to empty'")
		} else if runtime.GOOS == "darwin" { // macOS
			cmd = exec.Command("sh", "-c", "pmset -g batt | grep 'remaining'")
		} else { // windows
			cmd = exec.Command("WMIC", "Path", "Win32_Battery", "Get", "EstimatedRunTime")
		}

		// 명령어 실행
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Could not get battery time from OS. The battery might be charging or fully charged.")
			return
		}

		// OS별 결과값을 파싱하여 시간만 추출
		remainingTime, err := parseRemainingTime(string(output), runtime.GOOS)
		if err != nil {
			fmt.Printf("Could not parse battery output: %v\n", err)
			return
		}
		
		// 최종 결과 출력
		boldCyan("%s\n", remainingTime)

	default:
		// 지원하지 않는 OS일 경우, 기존의 수동 계산 방식(Fallback)으로 실행
		fmt.Println("Unsupported OS for direct query, using manual calculation...")
		
		var batteryData = batteries[0]
		var batState = batteryData.State
		var curBatCapMWh float64 = batteryData.Current
		var fullBatCapMWh float64 = batteryData.Full
		var powerMW float64 = batteryData.ChargeRate

		if powerMW <= 0 {
			fmt.Println("Cannot calculate time: Power consumption is zero or negative.")
			return
		}

		switch batState.String() {
		case battery.Charging.String():
			remainingCapMWh := fullBatCapMWh - curBatCapMWh
			timeRem := remainingCapMWh / powerMW
			hour := int(timeRem)
			minutes := int(60 * (float64(getFirstDecimalDigit(timeRem)) / 10))
			fmt.Println("Battery Charge Time")
			boldRed("%d:%02d\n", hour, minutes)
		default:
			timeRem := curBatCapMWh / powerMW
			hour := int(timeRem)
			minutes := int(60 * (float64(getFirstDecimalDigit(timeRem)) / 10))
			fmt.Println("Remaining Battery Time")
			boldCyan("%d:%02d\n", hour, minutes)
		}
	}
}
}

// parseRemainingTime 함수는 OS별 명령어 출력(문자열)을 받아 시간 값만 추출합니다.
func parseRemainingTime(output string, os string) (string, error) {
	switch os {
	case "darwin": // macOS
		// 출력 예: "... discharging; 1:30 remaining ..."
		// 'remaining' 바로 앞의 단어를 찾습니다.
		fields := strings.Fields(output)
		for i, field := range fields {
			if strings.Contains(field, "remaining") && i > 0 {
				return fields[i-1], nil
			}
		}
		return "", fmt.Errorf("could not find 'remaining' keyword in macOS output")

	case "linux":
		// 출력 예: "  time to empty:      1.5 hours"
		// ':'를 기준으로 문장을 나누고 뒷부분을 가져옵니다.
		parts := strings.SplitN(output, ":", 2)
		if len(parts) == 2 {
			return strings.TrimSpace(parts[1]), nil
		}
		return "", fmt.Errorf("could not parse linux output with ':'")

	case "windows":
		// 출력 예:
		// EstimatedRunTime
		// 90
		// 숫자(분)를 "H:MM" 형식으로 변환합니다.
		lines := strings.Split(strings.TrimSpace(output), "\n")
		if len(lines) >= 2 {
			minutesStr := strings.TrimSpace(lines[len(lines)-1]) // 마지막 줄이 숫자 값
			minutes, err := strconv.Atoi(minutesStr)
			if err != nil {
				return "", fmt.Errorf("failed to convert windows minutes to number")
			}
			h := minutes / 60
			m := minutes % 60
			return fmt.Sprintf("%d:%02d", h, m), nil
		}
		return "", fmt.Errorf("could not parse windows output")
	}
	return "", fmt.Errorf("unsupported OS for parsing")
}
