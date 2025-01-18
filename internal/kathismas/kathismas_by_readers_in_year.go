package kathismas

import (
	"maps"
	"time"
)

func getGenLoop(startLoop int, endLoop int) []int {
	newLoop := make([]int, 0)
	for i := startLoop; i < endLoop+1; i++ {
		newLoop = append(newLoop, i)
	}
	return newLoop
}

func cycleSlice(genLoop []int, loopFromTotalKathisma [20]int) map[int]int {
	loopFirst := make(map[int]int)
	kathismaIndex := 0
	for _, day := range genLoop {
		loopFirst[day] = loopFromTotalKathisma[kathismaIndex]
		kathismaIndex++
		if kathismaIndex >= len(loopFromTotalKathisma) {
			kathismaIndex = 0
		}
	}
	return loopFirst
}

func getCalendarMap(
	endNumberKathismaFirstLoop int,
	loopFromTotalKathisma [20]int,
	numberDaysInYear int,
	startZeroLoopSecond int,
	stepKathisma int,
) (map[int]int, map[int]int) {
	startNumberKathismaZeroLoopSecond := endNumberKathismaFirstLoop + stepKathisma
	zeroLoopSecond := make(map[int]int)
	count := 0
	for kathisma := startNumberKathismaZeroLoopSecond; kathisma < 21; kathisma++ {
		zeroLoopSecond[count+startZeroLoopSecond] = kathisma
		count++
	}
	startLoopSecond := startZeroLoopSecond + len(zeroLoopSecond)
	genLoopSecond := getGenLoop(startLoopSecond, numberDaysInYear+1)
	loopSecond := cycleSlice(genLoopSecond, loopFromTotalKathisma)
	return loopSecond, zeroLoopSecond
}

func GetListDate(
	startNoReading time.Time,
	endNoReading time.Time,
	startKathisma int,
	numberDaysInYear int,
	loopFromTotalKathisma [20]int,
) map[int]int {
	stepKathisma := 1
	zeroLoopFirst := make(map[int]int)
	count := 1
	for i := startKathisma; i < 21; i++ {
		zeroLoopFirst[count] = i
		count++
	}
	startLoopFirst := len(zeroLoopFirst) + stepKathisma
	endLoopFirst := startNoReading.YearDay() - 1
	startZeroLoopSecond := endNoReading.YearDay() + 1
	genLoopFirst := getGenLoop(startLoopFirst, endLoopFirst)
	loopFirst := cycleSlice(genLoopFirst, loopFromTotalKathisma)
	endNumberKathismaFirstLoop := loopFirst[endLoopFirst]
	loopSecond, zeroLoopSecond := getCalendarMap(
		endNumberKathismaFirstLoop,
		loopFromTotalKathisma,
		numberDaysInYear,
		startZeroLoopSecond,
		stepKathisma,
	)
	maps.Copy(zeroLoopFirst, loopFirst)
	maps.Copy(zeroLoopSecond, loopSecond)

	maps.Copy(zeroLoopFirst, zeroLoopSecond)
	return zeroLoopFirst
}
