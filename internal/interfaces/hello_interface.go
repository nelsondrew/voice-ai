package interfaces

type HelloInterface interface {
	GetMessage() string
	CreateMessage(name string) string
}
