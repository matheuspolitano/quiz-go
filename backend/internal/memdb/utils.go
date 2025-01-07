package memdb

import "fmt"

func combineIDs(id1, id2 string) string {
	return fmt.Sprintf("%s:%s", id1, id2)
}
