package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func CheckSheet(f *excelize.File, checkSheetNames []string) error {
	var e error = nil

	sheetNames := f.GetSheetList()
	if len(sheetNames) >= len(checkSheetNames) {
		for _, checkSheet := range checkSheetNames {
			flag := true

			for _, sheetName := range sheetNames {
				if sheetName == checkSheet {
					flag = false
					break
				}
			}

			if flag {
				e = fmt.Errorf("缺少%ssheet", checkSheet)
				break
			}
		}
	} else {
		e = fmt.Errorf("sheet 数量太少，期望： %d; 实际： %d", len(checkSheetNames), len(sheetNames))
	}

	return e
}

func ParserFloat(s string, bitSize int) (float64, error) {
	if len(s) == 0 {
		return 0.0, nil
	}

	valueStr := strings.ReplaceAll(s, ",", "")
	valueStr = strings.ReplaceAll(valueStr, "￥", "")

	return strconv.ParseFloat(valueStr, 32)
}
