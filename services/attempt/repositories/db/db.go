package db

type Database interface {
	Application() Application
	Endpoint() Endpoint
}
