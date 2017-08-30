package context

type Logger interface {
	Info(comment interface{})
	Warning(comment interface{}, err error)
	Fatal(comment interface{}, err error) string
}
