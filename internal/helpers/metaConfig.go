package helpers

import (
	"fmt"
	"math"
	"strconv"
)

func MetaConfig(totalData, page int, host, query string) []interface{} {
	var result []interface{}
	// url := fmt.Sprintf("%s%s", host, query)
	totalPage := math.Ceil(float64(totalData) / 5)
	isLastPage := page >= int(totalPage)
	var next any
	var prev any
	result = append(result, page)
	result = append(result, totalPage)
	result = append(result, totalData)
	if isLastPage && page == 1 {
		next = nil
		prev = nil
		result = append(result, next)
		result = append(result, prev)
		return result
	}
	if isLastPage {
		next = nil
		prev = fmt.Sprintf("%s%s", host, query[:len(query)-1]+strconv.Itoa(page-1))
		result = append(result, next)
		result = append(result, prev)
		return result
	}
	if page == 1 {
		next = fmt.Sprintf("%s%s", host, query[:len(query)-1]+strconv.Itoa(page+1))
		prev = nil
		result = append(result, next)
		result = append(result, prev)
		return result
	}

	next = fmt.Sprintf("%s%s", host, query[:len(query)-1]+strconv.Itoa(page+1))
	prev = fmt.Sprintf("%s%s", host, query[:len(query)-1]+strconv.Itoa(page-1))
	result = append(result, next)
	result = append(result, prev)
	return result
}