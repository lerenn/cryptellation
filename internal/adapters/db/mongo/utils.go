package mongo

func ToInterfaceList[T any](l []T) []interface{} {
	res := make([]interface{}, len(l))
	for i, v := range l {
		res[i] = v
	}
	return res
}
