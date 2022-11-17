package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var y = flag.String("y", "", "")
var m = flag.String("m", "", "")
var WeekDayMap = map[string]string{
	"Monday":    "周一",
	"Tuesday":   "周二",
	"Wednesday": "周三",
	"Thursday":  "周四",
	"Friday":    "周五",
	"Saturday":  "周六",
	"Sunday":    "周日",
}

type AutoGenerated []struct {
	Date            string `json:"date"`
	Name            string `json:"name"`
	IsHoliday       string `json:"isHoliday"`
	HolidayCategory string `json:"holidayCategory"`
	Description     string `json:"description"`
}

func main() {

	flag.Parse()
	if string(*y) == "" {
		log.Fatal("please input your year: -y=111")
		return
	}

	if string(*m) == "" {
		log.Fatal("please input your month: -m=6")
		return
	}

	// 讀取excel, 準備建立範例
	excel := getExcel("example.xlsx")

	excel, err := copySheet(excel, "空白範本", "正確空白範本")
	if err != nil {
		return
	}

	// 抓取線上政府例假日
	excel = getJson(excel)

	// 先個別讀出資料
	holiday := readHoliday(excel)
	leaveList := readLeaveList(excel)
	employee := readEmployee(excel)

	// 正確版的excel
	newExcel := sortOutExcelTemp(excel, holiday)

	// 開始進行每一個人的處理
	// 1: 部門,	2: 職稱, 3: 姓名
	for z, rowE := range employee {
		if z == 0 {
			continue
		}

		_, err := copySheet(newExcel, "正確空白範本", rowE[3])
		if err != nil {
			return
		}

		title := fmt.Sprintf("%s  年  %s  月份   出勤紀錄", *y, *m)
		err = newExcel.SetCellValue(rowE[3], "A1", title)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 設定部門
		err = newExcel.SetCellValue(rowE[3], "B2", rowE[1])
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 設定職稱
		err = newExcel.SetCellValue(rowE[3], "H2", rowE[2])
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 設定名字
		err = newExcel.SetCellValue(rowE[3], "B3", rowE[3])
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 各別請假處理
		leaveCheck(newExcel, rowE[3], leaveList)

		// 填入隨機上班時間
		randTime(newExcel, rowE[3])

	}
	fileName := fmt.Sprintf("gaia-%s-%s.xlsx", *y, *m)
	err = newExcel.SaveAs(fileName)
	if err != nil {
		fmt.Println("save error", err)
		return
	}
}
func getExcel(fileName string) *excelize.File {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	return f
}
func readEmployee(excel *excelize.File) [][]string {
	rows, err := excel.GetRows("員工清單")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return rows
}
func readHoliday(excel *excelize.File) [][]string {
	rows, err := excel.GetRows("當月假日")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return rows
}
func sortOutExcelTemp(excel *excelize.File, holiday [][]string) *excelize.File {
	rows, err := excel.GetRows("空白範本")
	if err != nil {
		fmt.Println(err)
	}

	// 0: 幾號, 1: 節日
	for i, row := range rows {
		if i < 6 || i > 37 {
			continue
		}

		// 先處理國定假日
		for x, rowH := range holiday {
			if x < 1 {
				continue
			}
			// 假日
			if row[0] == rowH[0] {
				axis := fmt.Sprintf("B%d", i+1)
				err = excel.SetCellValue("正確空白範本", axis, rowH[1])
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
		}

		// 在處理週末例假
		yy, _ := strconv.Atoi(*y)
		mm, _ := strconv.Atoi(*m)
		d := getYearMonthToDay(yy, mm)
		for f := 1; f <= d; f = f + 1 {
			axis := fmt.Sprintf("B%d", f+6)
			date := fmt.Sprintf("%d-%02d-%02d", yy+1911, mm, f)
			t, err := time.Parse("2006-01-02", date)
			if err != nil {
				panic(err)
			}
			week := WeekDayMap[t.Weekday().String()]

			if week == "周六" || week == "周日" {
				// 發現有國定假日
				txt, _ := excel.GetCellValue("正確空白範本", axis)
				if txt != "周六" && txt != "周日" && txt != "" {
					if week == "周六" {
						changeAxis := fmt.Sprintf("B%d", f+6-1)
						changeTxt := fmt.Sprintf("%s-遇假日補休", txt)
						_ = excel.SetCellValue("正確空白範本", changeAxis, changeTxt)
						_ = excel.MergeCell("正確空白範本", changeAxis, fmt.Sprintf("C%d", f+6-1))
					}

					if week == "周日" {
						changeAxis := fmt.Sprintf("B%d", f+6+1)
						changeTxt := fmt.Sprintf("%s-遇假日補休", txt)
						_ = excel.SetCellValue("正確空白範本", changeAxis, changeTxt)
						_ = excel.MergeCell("正確空白範本", changeAxis, fmt.Sprintf("C%d", f+6+1))
					}
				} else {
					_ = excel.SetCellValue("正確空白範本", axis, week)
				}
			}
		}
	}

	return excel
}
func copySheet(excel *excelize.File, oldSheetName, newSheetName string) (*excelize.File, error) {
	fromIndex := excel.GetSheetIndex(oldSheetName)
	index := excel.NewSheet(newSheetName)
	err := excel.CopySheet(fromIndex, index)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return excel, nil
}
func readLeaveList(excel *excelize.File) [][]string {
	rows, err := excel.GetRows("請假清單")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return rows
}
func getYearMonthToDay(year int, month int) int {
	day31 := map[int]struct{}{
		1:  struct{}{},
		3:  struct{}{},
		5:  struct{}{},
		7:  struct{}{},
		8:  struct{}{},
		10: struct{}{},
		12: struct{}{},
	}
	if _, ok := day31[month]; ok {
		return 31
	}

	day30 := map[int]struct{}{
		4:  struct{}{},
		6:  struct{}{},
		9:  struct{}{},
		11: struct{}{},
	}
	if _, ok := day30[month]; ok {
		return 30
	}

	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		return 29
	}

	return 28
}
func randTime(newExcel *excelize.File, name string) {
	person, err := newExcel.GetRows(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	for u, rowP := range person {
		if u < 6 {
			continue
		}

		if u == 36 {
			if *m == "2" || *m == "4" || *m == "6" || *m == "9" || *m == "11" {
				continue
			}
		}

		// fmt.Println(u, rowP, *m)
		time.Sleep(100 * time.Millisecond)

		sH := randomInt(9, 10)
		sM := randomInt(0, 30)
		eH := sH + 9
		eM := randomInt(sM, sM+29)

		if rowP[1] == "" && rowP[5] == "" {
			// fmt.Println(rowP[1], rowP[5], u+1)
			axis := fmt.Sprintf("%s%d", "B", u+1)
			_ = newExcel.SetCellValue(name, axis, sH)

			axis = fmt.Sprintf("%s%d", "C", u+1)
			_ = newExcel.SetCellValue(name, axis, sM)

			axis = fmt.Sprintf("%s%d", "D", u+1)
			_ = newExcel.SetCellValue(name, axis, eH)

			axis = fmt.Sprintf("%s%d", "E", u+1)
			_ = newExcel.SetCellValue(name, axis, eM)
		}
	}
}
func leaveCheck(newExcel *excelize.File, name string, leaveList [][]string) {
	// 1: 部門,	2: 職稱, 3: 姓名, 4: 假別, 5: 請假日期, 6: 起時間, 7: 迄時間, 8: 時數
	for w, rowL := range leaveList {
		if w == 0 {
			continue
		}

		if name == rowL[3] {
			// fromIndex := excel.GetSheetIndex(rowE[3])
			// 寫入請假假別
			tmpDay := strings.Split(rowL[5], "-")
			day, _ := strconv.Atoi(tmpDay[1])
			month, _ := strconv.Atoi(tmpDay[0])
			mString := strconv.Itoa(month)
			if mString != *m {
				log.Fatalln("輸入請假日期有非本月份, 請確認後再執行一次")
			}

			axis := fmt.Sprintf("%s%d", "F", day+6)
			_ = newExcel.SetCellValue(name, axis, rowL[4])

			// 寫入時間
			if rowL[4] != "特休" {
				hh, _ := strconv.Atoi(rowL[8])

				sTime := strings.Split(rowL[6], ":")
				axis1 := fmt.Sprintf("%s%d", "B", day+6)
				axis2 := fmt.Sprintf("%s%d", "C", day+6)
				//eTime := strings.Split(rowL[7], ":")
				axis3 := fmt.Sprintf("%s%d", "D", day+6)
				axis4 := fmt.Sprintf("%s%d", "E", day+6)

				if hh == 8 {
					//_ = newExcel.SetCellValue(name, axis1, sTime[0])
					//_ = newExcel.SetCellValue(name, axis2, sTime[1])
					//_ = newExcel.SetCellValue(name, axis3, eTime[0])
					//_ = newExcel.SetCellValue(name, axis4, eTime[1])
				} else {
					// 處理半天請假上下午
					qq, _ := strconv.Atoi(sTime[0])
					a := randomInt(0, 30)

					// 請上午
					if qq <= 13 {
						_ = newExcel.SetCellValue(name, axis1, "14")
						_ = newExcel.SetCellValue(name, axis2, a)
						_ = newExcel.SetCellValue(name, axis3, "18")
						_ = newExcel.SetCellValue(name, axis4, randomInt(a, 59))
					} else { // 請下午

						_ = newExcel.SetCellValue(name, axis1, "9")
						_ = newExcel.SetCellValue(name, axis2, a)
						_ = newExcel.SetCellValue(name, axis3, "13")
						_ = newExcel.SetCellValue(name, axis4, randomInt(a, 59))
					}
				}
			}

			// 寫入請假時數
			axis = fmt.Sprintf("%s%d", "G", day+6)
			err := newExcel.SetCellValue(name, axis, rowL[8])
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}
func getJson(excel *excelize.File) *excelize.File {
	resp, err := http.Get("https://data.ntpc.gov.tw/api/datasets/308DCD75-6434-45BC-A95F-584DA4FED251/json?page=4&size=270")
	if err != nil {
		log.Fatalln("抓取行事曆失敗:", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("抓取行事曆失敗:", err)
	}

	var data AutoGenerated
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalln("行事曆解析失敗:", err)
	}

	excel.NewSheet("當月假日")
	_ = excel.SetCellValue("當月假日", "A1", "日期")
	_ = excel.SetCellValue("當月假日", "B1", "節日")

	x := 1
	for _, row := range data {
		month := strings.Split(row.Date, "/")
		q, _ := strconv.Atoi(*y)
		q = q + 1911

		// fmt.Println(row)
		// fmt.Println(q, month[0], month[1], row.Name)
		// fmt.Println((month[1] == *m && month[0] == strconv.Itoa(q)), (row.Name != "" && row.IsHoliday == "是" && row.HolidayCategory == "星期六、星期日"))

		if (month[1] == *m && month[0] == strconv.Itoa(q)) && (row.Name != "" && row.IsHoliday == "是" && (row.HolidayCategory == "星期六、星期日" || row.HolidayCategory == "放假之紀念日及節日")) {
			//fmt.Println(month[1])
			// fmt.Println(q, month[0], month[1], row.Name)
			x = x + 1

			axis := fmt.Sprintf("A%d", x)
			dateS := fmt.Sprintf("%s號", month[2])
			_ = excel.SetCellValue("當月假日", axis, dateS)

			axis = fmt.Sprintf("B%d", x)
			_ = excel.SetCellValue("當月假日", axis, row.Name)
		}
	}

	return excel
}
func randomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}
