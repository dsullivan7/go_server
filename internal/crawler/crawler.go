package crawler

type Crawler interface {
	Login(url string, username string, password string) string
}
