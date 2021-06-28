package notification

type pageInfo struct {
	name   string
	url    string
	status int
}

func (p *pageInfo) Name() string {
	return p.name
}

func (p *pageInfo) URL() string {
	return p.url
}

func (p *pageInfo) Status() int {
	return p.status
}

func NewPageInfo(name string, url string, status int) PageInfo {
	return &pageInfo{name: name, url: url, status: status}
}
