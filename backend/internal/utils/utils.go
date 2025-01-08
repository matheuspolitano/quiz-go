package utils

import "fmt"

func CombineIDs(id1, id2 string) string {
	return fmt.Sprintf("%s:%s", id1, id2)
}
