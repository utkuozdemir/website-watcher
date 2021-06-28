package notification

type Handler interface {
	Handle(pageInfo PageInfo)
}

type PageInfo interface {
	Name() string
	URL() string
	Status() int
}
