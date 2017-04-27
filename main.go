package main

import(
    "github.com/gorilla/websocket"
    "net/http"
    "strings"
    "log" 
    "github.com/asiainfoLDP/chat/message"
    "html/template"
   // "github.com/asiainfoLDP/chat/m2"
    "bytes"
    "fmt"
    "net"
    "time"
    "io"
)
const PrefixWS = "/ws/"
var wsUpgrader = websocket.Upgrader{
    ReadBufferSize:      512,
    WriteBufferSize:     512,
    CheckOrigin:         func(r *http.Request)  bool {return true},

}
func websocketHandler (w http.ResponseWriter, r *http.Request)  {
    log.Println("starbfffffffffffffffffffffn")
    if r.Method != "GET" {
        http.Error(w, "Method don't allow", 405)
        return
    }
    log.Println(":11111111111")
    if (!strings.HasPrefix(r.URL.Path, PrefixWS)) || len(r.URL.Path) < len(PrefixWS){
        http.Error(w, "bad uri", 400)
        return 
    }
        log.Println(":22222222222222222")

    var wsconn, err = wsUpgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "Method not allowed", 405)
        return
    }
        log.Println("33333333333333333333")
        server.CreateConn <- &ChatConn{Conn: wsconn}

}
func sendPageData(w http.ResponseWriter, pageDataBytes []byte, contextType string) error {

	w.Header().Set("Content-Type", contextType)

	var numWrites int
	var err error

	numBytes := len(pageDataBytes)
	for numBytes > 0 {
		numWrites, err = w.Write(pageDataBytes)
		if err != nil {
			return err
		}
		numBytes -= numWrites
	}

	return nil
}
func httpHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
    var httptemplate *template.Template = nil
    var httpContentCache    []byte = nil
    //var err error
    if httptemplate == nil {
        httpTemplate, err :=template.ParseFiles("template/websocketclient.html")
        if err != nil { 
			sendPageData(w, []byte("Parse template error."), "text/plain; charset=utf-8")
			return
		}
        if httpContentCache == nil {
		var buf bytes.Buffer
		err = httpTemplate.Execute(&buf, nil)
		if err != nil {
			sendPageData(w, []byte("Render page error."), "text/plain; charset=utf-8")
			return
		}

		httpContentCache = buf.Bytes()
	}

	sendPageData(w, httpContentCache, "text/html; charset=utf-8")
}
    



}




func createWebsocketServer(port int){
    log.Printf("Websocket listening at : %d ...\n", port)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    http.HandleFunc("/ws/", websocketHandler)
    http.HandleFunc("/", httpHandler)
    adress := fmt.Sprintf(":%d", port)
    err := http.ListenAndServe(adress, nil)
    if err != nil {
        log.Fatal("Websocket server failt to start: ", err)
    }
}
var server *message.Server
func main(){
     server = &message.Server{
        CreateConn:  make(chan net.Conn, 8),
        CloseConn:   make(chan net.Conn, 8),
        Visitors:    make(map[net.Conn]*message.Visitor, 30),
        Messages:    make(chan string, 30),
    }
    go server.Run()
     createWebsocketServer(9999)
   // go m2.ClienttosServer()
   // message.ServerConnection(9999)
   
     
} 
func (cc *ChatConn) Close() error {
	return cc.Conn.Close()
}

func (cc *ChatConn) LocalAddr() net.Addr {
	return cc.Conn.LocalAddr()
}

func (cc *ChatConn) RemoteAddr() net.Addr {
	return cc.Conn.RemoteAddr()
}

func (cc *ChatConn) SetDeadline(t time.Time) error {
	var err = cc.Conn.SetReadDeadline(t)
	if err == nil {
		err = cc.Conn.SetWriteDeadline(t)
	}

	return err
}

func (cc *ChatConn) SetReadDeadline(t time.Time) error {
	return cc.Conn.SetReadDeadline(t)
}

func (cc *ChatConn) SetWriteDeadline(t time.Time) error {
	return cc.Conn.SetWriteDeadline(t)
}
func (cc *ChatConn) Read(b []byte) (int, error) {
	// read from buffer
	//bs := []byte{1,   0, 1, 12, 2, 22,   2,   0, 2, 12, 1, 22}
	//nn := copy(b, bs)
	//return nn, nil

	index := 0
	n, err := cc.InputBuffer.Read(b)
	index += n
	if err != nil && err != io.EOF {
		return index, err
	}

	if index > 0 {
		return index, nil
	}

	for {
		// try to read more message data and cache it
		messageType, p, err := cc.Conn.ReadMessage()
		if err != nil && err != io.EOF {
			return index, err
		}
        log.Println(messageType, string(p))
		if messageType != websocket.TextMessage { // only accept BinaryMessage messages
			continue
	    }
        p = append(p, '\n')
        log.Println("1111111111111111")
		// n2 must be len(p) if err2 is nil
		n2, err2 := cc.InputBuffer.Write(p) // cache it
		if err2 != nil {
			return index, err2
		}
        log.Println("1111111111111111", n2)
		if n2 > 0 {
			break
		}
        log.Println("3333333333333333333")
	}
log.Println("444444444444444444")
	// read from buffer again

	n, err = cc.InputBuffer.Read(b[index:])
	index += n
	if err != nil && err != io.EOF {
		return index, err
	}

	return index, nil
}
func (cc *ChatConn) Write(b []byte) (int, error) {
	return len(b), cc.Conn.WriteMessage(websocket.TextMessage, b) // todo: maybe not ok
}
type ChatConn struct { // implement chat.ReadWriteCloser
	Conn *websocket.Conn

	InputBuffer  bytes.Buffer
	OutputBuffer bytes.Buffer
}

