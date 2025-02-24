package Algorithm

func MakeKeyHash(Key int) int {
	var a int
	var b int
	a_result := 0
	b_result := 0

	for i := 0; i < 16; i++ {
		a = 1 << (2 * i)
		b = 1 << (2 * i + 1)

		a &= Key
		b &= Key

		a = a << (2 * i)
		b = b << (2 * i + 1)

		a_result |= a
		b_result |= b
	}
	return a_result | b_result
}
