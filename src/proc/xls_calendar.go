package proc

import (
	"bytes"
	"fmt"
	"github.com/DjaPy/fot-twenty-readers-go/internal/kathismas"
	"github.com/xuri/excelize/v2"
	"strconv"
	"time"
)

const FONTTREBUCHET = "Trebuchet MS"

func addKathismaNumbersToXLS(xls *excelize.File, number int, sheetName string) {
	style, err := xls.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 3},
			{Type: "top", Color: "000000", Style: 3},
			{Type: "bottom", Color: "000000", Style: 3},
			{Type: "right", Color: "000000", Style: 3},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Family: FONTTREBUCHET, Bold: true, Size: 16},
	})
	if err != nil {
		panic(err)
	}
	xls.SetCellValue(sheetName, "A2", strconv.Itoa(number))
	xls.SetCellStyle(sheetName, "A2", "A2", style)
}

func addHeaderOfMonthToWs(xls *excelize.File, sheetName string) {
	cellAddressMonth := map[string]string{
		"B2": "ЯНВ", "C2": "ФЕВ", "D2": "МАРТ", "E2": "АПР",
		"F2": "МАЙ", "G2": "ИЮН", "H2": "ИЮЛ", "I2": "АВГ",
		"J2": "СЕН", "K2": "ОКТ", "L2": "НОЯ", "M2": "ДЕК",
	}
	style, err := xls.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 0},
			{Type: "top", Color: "000000", Style: 0},
			{Type: "bottom", Color: "000000", Style: 0},
			{Type: "right", Color: "000000", Style: 0},
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

func addColumnWithNumberDayToWs(xls *excelize.File, sheetName string) {
	style, _ := xls.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Font:      &excelize.Font{Family: FONTTREBUCHET, Size: 12},
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
		Font:      &excelize.Font{Family: FONTTREBUCHET, Size: 14, Color: "000000"},
	})
	cellStep := 1
	frameMonth := map[int]string{
		1: "B", 2: "C", 3: "D", 4: "E", 5: "F", 6: "G", 7: "H", 8: "I", 9: "J", 10: "K", 11: "L", 12: "M",
	}
	frameNumberDayA := getFrameNumberDay("A", 3, 33) // A = 1
	frameNumberDayN := getFrameNumberDay("N", 3, 33) // N = 1
	for num := range frameNumberDayN {
		xls.SetCellValue(sheetName, frameNumberDayN[num], strconv.Itoa(num-2))
		xls.SetCellValue(sheetName, frameNumberDayA[num], strconv.Itoa(num-2))
	}

	for month, days := range calendarTable {
		cellMonth := frameMonth[month]
		cellNameIndex := 2
		var keyDayStr string
		for _, day := range days {
			cellNameIndex += cellStep
			cellName := cellMonth + strconv.Itoa(cellNameIndex)
			targetDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			dayNow := targetDate.YearDay()

			if keyDay, ok := allKathisma[dayNow]; !ok {
				keyDayStr = ""
			} else {
				keyDayStr = strconv.Itoa(keyDay)
			}
			xls.SetCellValue(sheetName, cellName, keyDayStr)
			xls.SetCellStyle(sheetName, cellName, cellName, style)
		}
	}
}

func CreateXlSCalendar(startDate time.Time, startKathisma, year int) (*bytes.Buffer, error) {
	if year == 0 {
		year = startDate.Year()
	}
	calendarTable := kathismas.GetCalendarYear(startDate, year)
	calendarKathismas := kathismas.CreateCalendar(startDate, startKathisma, year)

	xls := excelize.NewFile()
	defer func() {
		if err := xls.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	for pair := calendarKathismas.Oldest(); pair != nil; pair = pair.Next() {
		sheetName := fmt.Sprintf("Чтец %d", pair.Key)
		xls.NewSheet(sheetName)
		addKathismaNumbersToXLS(xls, pair.Key, sheetName)
		addHeaderOfMonthToWs(xls, sheetName)
		addColumnWithNumberDayToWs(xls, sheetName)
		CreateCalendarForReaderToXLS(xls, calendarTable, pair.Value, year, sheetName)
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
