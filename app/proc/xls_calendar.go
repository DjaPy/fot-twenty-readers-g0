package proc

import (
	"bytes"
	"fmt"
	"github.com/xuri/excelize/v2"
	"maps"
	"strconv"
	"time"
)

func getCalendarForTable(startCalendarDate time.Time, year int) map[int][]int {
	tableYear := make(map[int][]int)
	var currentDayList []int
	currentMonth := 1

	for {
		if startCalendarDate.Year() == year+1 {
			break
		}
		if startCalendarDate.Day() == 1 {
			if len(currentDayList) > 0 {
				tableYear[currentMonth] = currentDayList
			}
			currentDayList = []int{}
			currentMonth = int(startCalendarDate.Month())
		}
		currentDayList = append(currentDayList, startCalendarDate.Day())
		startCalendarDate = startCalendarDate.AddDate(0, 0, 1)
	}

	if len(currentDayList) > 0 {
		tableYear[currentMonth] = currentDayList
	}

	return tableYear
}

func GetEasterDate(year int) time.Time {
	a := year % 4
	b := year % 7
	c := year % 19
	d := (19*c + 15) % 30
	e := (2*a + 4*b - d + 34) % 7
	month := (d + e + 114) / 31
	day := ((d + e + 114) % 31) + 1

	// Перевод даты из Юлианского в Григорианский календарь
	easter := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	offset := year/100 - year/400 - 2
	easter = easter.AddDate(0, 0, offset)
	return easter
}

func GetBoundaryDays(easterDay time.Time) (time.Time, time.Time) {
	startNoReading := easterDay.AddDate(0, 0, -3)
	endNoReading := easterDay.AddDate(0, 0, 6)
	return startNoReading, endNoReading
}

func GetNumberDaysInYear(year int) int {
	startYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	endYear := startYear.AddDate(1, 0, 0)
	return int(endYear.Sub(startYear).Hours() / 24)
}

func addKathismaNumbersToXLS(xls *excelize.File, number int, sheetName string) {
	style, err := xls.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 8},
			{Type: "top", Color: "000000", Style: 8},
			{Type: "bottom", Color: "000000", Style: 8},
			{Type: "right", Color: "000000", Style: 8},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Family: "Trebuchet MS", Bold: true, Size: 16},
	})
	if err != nil {
		panic(err)
	}
	xls.SetCellValue(sheetName, "A2", strconv.Itoa(number))
	xls.SetCellStyle(sheetName, "A2", "A2", style)
}

func AddHeaderOfMonthToWs(xls *excelize.File, sheetName string) {
	cellAddressMonth := map[string]string{
		"B2": "ЯНВ", "C2": "ФЕВ", "D2": "МАРТ", "E2": "АПР",
		"F2": "МАЙ", "G2": "ИЮН", "H2": "ИЮЛ", "I2": "АВГ",
		"J2": "СЕН", "K2": "ОКТ", "L2": "НОЯ", "M2": "ДЕК",
	}
	style, err := xls.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 8},
			{Type: "top", Color: "000000", Style: 8},
			{Type: "bottom", Color: "000000", Style: 8},
			{Type: "right", Color: "000000", Style: 8},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Family: "Calibri", Color: "FF8080", Size: 16},
	})
	if err != nil {
		panic(err)
	}
	for k, v := range cellAddressMonth {
		xls.SetCellValue(sheetName, k, v)
		xls.SetCellStyle(sheetName, k, k, style)
	}
}

func AddColumnWithNumberDayToWs(xls *excelize.File, sheetName string) {
	style, _ := xls.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Family: "Trebuchet MS", Size: 12},
	})

	for number := 1; number <= 31; number++ {
		numberCell := number + 2
		cellNameLeft := fmt.Sprintf("A%d", numberCell)
		cellNameRight := fmt.Sprintf("N%d", numberCell)
		xls.SetCellValue(sheetName, cellNameLeft, strconv.Itoa(number))
		xls.SetCellValue(sheetName, cellNameRight, strconv.Itoa(number))
		xls.SetCellStyle(sheetName, cellNameLeft, cellNameLeft, style)
		xls.SetCellStyle(sheetName, cellNameRight, cellNameRight, style)
	}
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

func getGenLoopFirst(startLoop int, endLoop int) []int {
	newLoop := make([]int, 0)
	for i := startLoop; i < endLoop+1; i++ {
		newLoop = append(newLoop, i)
	}
	return newLoop
}

func getCalendarHash(
	endNumberKathismaFirstLoop int,
	loopFromTotalKathisma [20]int,
	numberDaysInYear int,
	startZeroLoopSecond int,
	stepKathisma int,
) (map[int]int, map[int]int) {
	startNumberKathismaZeroLoopSecond := endNumberKathismaFirstLoop + stepKathisma
	zeroLoopSecond := make(map[int]int)
	count := 0
	for kathisma := startNumberKathismaZeroLoopSecond; kathisma <= 21; kathisma++ {
		zeroLoopSecond[count+startNumberKathismaZeroLoopSecond] = kathisma
		count++
	}
	startLoopSecond := startZeroLoopSecond + len(zeroLoopSecond)
	genLoopSecond := getGenLoopFirst(startLoopSecond, numberDaysInYear+1)
	loopSecond := cycleSlice(genLoopSecond, loopFromTotalKathisma)
	return loopSecond, zeroLoopSecond
}

func getListDate(
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
	genLoopFirst := getGenLoopFirst(startLoopFirst, endLoopFirst)
	loopFirst := cycleSlice(genLoopFirst, loopFromTotalKathisma)
	endNumberKathismaFirstLoop := loopFirst[endLoopFirst]
	loopSecond, zeroLoopSecond := getCalendarHash(
		endNumberKathismaFirstLoop,
		loopFromTotalKathisma,
		numberDaysInYear,
		startZeroLoopSecond,
		startKathisma,
	)
	maps.Copy(zeroLoopFirst, loopFirst)
	maps.Copy(zeroLoopSecond, loopSecond)
	maps.Copy(zeroLoopFirst, zeroLoopSecond)
	fmt.Println(zeroLoopFirst[endLoopFirst+1], zeroLoopFirst[endLoopFirst])

	return zeroLoopFirst
}

func getFrameNumberDay(symbol string, start int, end int) map[int]string {
	frameNumberDay := make(map[int]string)
	for num := start; num <= end; num++ {
		frameNumberDay[num] = symbol + strconv.Itoa(num)
	}
	return frameNumberDay
}

func CreateCalendarForReaderToXLS(
	xls *excelize.File,
	calendarTable map[int][]int,
	allKathisma map[int]int,
	year int,
	sheetName string,
) {
	style, _ := xls.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Family: "Trebuchet MS", Size: 14, Color: "000000"},
	})
	cellStep := 1
	frameMonth := map[int]string{
		1: "B", 2: "C", 3: "D", 4: "E", 5: "F", 6: "G", 7: "H", 8: "I", 9: "J", 10: "K", 11: "L", 12: "M",
	}
	frameNumberDayA := getFrameNumberDay("A", 3, 33) // A = 1
	frameNumberDayN := getFrameNumberDay("N", 3, 33) // N = 1
	for num, _ := range frameNumberDayN {
		xls.SetCellValue(sheetName, frameNumberDayN[num], strconv.Itoa(num))
		xls.SetCellValue(sheetName, frameNumberDayA[num], strconv.Itoa(num))
	}

	for month, days := range calendarTable {
		cellMonth := frameMonth[month]
		cellNameIndex := 2
		for _, day := range days {
			cellNameIndex += cellStep
			cellName := cellMonth + strconv.Itoa(cellNameIndex)
			targetDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			dayNow := targetDate.Day()
			xls.SetCellValue(sheetName, cellName, strconv.Itoa(allKathisma[dayNow]))
			xls.SetCellStyle(sheetName, cellName, cellName, style)
		}
	}
}

func CreateXlSCalendar(startDate time.Time, startKathisma, year int) (*bytes.Buffer, error) {
	if year == 0 {
		year = startDate.Year()
	}
	calendarTable := getCalendarForTable(startDate, year)
	easterDay := GetEasterDate(year)
	startNoReading, endNoReading := GetBoundaryDays(easterDay)
	numberDaysInYear := GetNumberDaysInYear(year)
	totalKathismas := [20]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	xls := excelize.NewFile()
	defer func() {
		if err := xls.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	for _, numberKathisma := range totalKathismas {
		sheetName := fmt.Sprintf("Чтец %d", numberKathisma)
		xls.NewSheet(sheetName)
		addKathismaNumbersToXLS(xls, numberKathisma, sheetName)
		AddHeaderOfMonthToWs(xls, sheetName)
		AddColumnWithNumberDayToWs(xls, sheetName)
		allKathismas := getListDate(startNoReading, endNoReading, numberKathisma, numberDaysInYear, totalKathismas)
		CreateCalendarForReaderToXLS(xls, calendarTable, allKathismas, year, sheetName)
		if startKathisma > 19 {
			startKathisma = 0
		}
		startKathisma += 1
	}
	p := getPathForFile()
	err := xls.SaveAs(p.outFile)
	if err != nil {
		return nil, err
	}
	return xls.WriteToBuffer()
}
