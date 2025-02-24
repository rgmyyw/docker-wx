package lib

func ALLGather(Data []string) string {
	var M string
	for i := 0; i < len(Data); i++ {
		M += Data[i]
	}
	return M
}
