package utils

import "fmt"

// JoinUintArray Joins the elements of the uint array into a string separated by the `delimiter`
func JoinUintArray(data []uint, delimiter string) string {
	result := ""
	for _, value := range data {
		if len(result) > 0 {
			result += delimiter
		}
		result += fmt.Sprintf("%d", value)
	}
	return result
}
