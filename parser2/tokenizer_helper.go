package parser2

type tokenData struct {
	t  token
	td []byte
}

func tokenize(data string) ([]tokenData, bool) {
	result := []tokenData{}

	res := tokenizeRaw([]byte(data), func(t token, td []byte) {
		result = append(result, tokenData{t, td})
	})

	if !res {
		return nil, false
	}

	return result, true
}
