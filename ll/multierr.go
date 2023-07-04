package ll

import "bytes"

type multiErr struct {
	errs []error
}

func (m *multiErr) add(err error) {
	if err == nil {
		return
	}
	m.errs = append(m.errs, err)
}

func (m *multiErr) Error() string {
	var buff bytes.Buffer
	buff.WriteString("Multiple errors: ")
	lastI := len(m.errs) - 1
	for i, err := range m.errs {
		buff.WriteString(err.Error())
		if i != lastI {
			buff.WriteString(", ")
		}
	}
	return buff.String()
}
