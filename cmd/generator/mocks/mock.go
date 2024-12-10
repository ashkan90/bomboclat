package mocks

const mockTemplate = `
{{range .Signatures}}
type Mock{{.Name}}Client struct {
    mock.Mock
}

func (m *Mock{{.Name}}Client) {{.Name}}({{range .Args}}{{.Name}} {{.Type}},{{end}}) ({{range .ReturnType}}{{.Type}},{{end}}) {
    args := m.Called({{range .Args}}{{.Name}},{{end}})
    return args.Get(0).({{index .ReturnType 0}}), args.Error(1)
}
{{end}}
`
