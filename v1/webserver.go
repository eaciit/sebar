package sebar

type WebServer struct {
	SebarServer

	_webUrl string
}

func (w *WebServer) SetWebUrl(urlText string) *WebServer {
	w._webUrl = urlText
	return w
}

func (w *WebServer) WebUrl() string {
	return w._webUrl
}

func (w *WebServer) Start() error {
	e := w.SebarServer.Start()
	if e != nil {
		return e
	}
	return nil
}
