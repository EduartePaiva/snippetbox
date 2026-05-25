package forms

type errors map[string][]string

func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e errors) Get(field string) string {
	message := e[field]
	if len(message) == 0 {
		return ""
	}
	return message[0]
}
