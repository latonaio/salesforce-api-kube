package str

// ToFirstUppercase converts the first letter to uppercase.
func ToFirstUppercase(s string) string {
	r := []rune(s)
	if 'a' <= r[0] && r[0] <= 'z' {
		r[0] += 'A' - 'a'
	}
	return string(r)
}