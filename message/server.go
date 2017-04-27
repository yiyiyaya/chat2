package message



import (
    "net"
    "fmt"
    "bufio"
   // "os"
	"log"
	//"github.com/asiainfoLDP/datafoundry_proxy/messages"
	//"os"
	//"hash/adler32"
)
    
func ServerConnection(port int){
   Server := Server{
        CreateConn:  make(chan net.Conn, 8),
        CloseConn:   make(chan net.Conn, 8),
        Visitors:    make(map[net.Conn]*Visitor, 30),
        Messages:    make(chan string, 30),
    }
    fmt.Println("will start server......")
   // 
    go Server.Run()
    adress := fmt.Sprintf(":%d",port)
    listener, err := net.Listen("tcp", adress)
    if err != nil {
        log.Println("Listen net err")
        return
    }
    for {
        //建立一个客户端的连接
        conn, err := listener.Accept()
        if err != nil {
            log.Println("get client conn err")
        } else {
            Server.CreateConn <- conn
            log.Println(conn.RemoteAddr)
        }
    }
    
   
}
type Server struct{
    CreateConn     chan net.Conn
    CloseConn      chan net.Conn
    Visitors       map[net.Conn]*Visitor
    Messages       chan string
    NextId         int
}
 func (server *Server)Run(){
        for{
            //log.Println("hjbngjnjnjjjj")
            select {
                case conn := <- server.CreateConn :
                    visitor := &Visitor{
                            Conn:        conn,
                            Name:        fmt.Sprintf("Visitor#%d", server.NextId),
                            Server:      server,

                            Message:      make(chan string, 30),
                            ClientReader: bufio.NewReader(conn),
                            CloseInfo:    make(chan int, 3),
                            Closed:       make(chan struct{}),

                    }  
                    server.Visitors[conn] = visitor
                    server.NextId++ 
                    go visitor.Run() 
                    log.Println("connection success")
                case conn :=  <-server.CloseConn :
                    delete(server.Visitors, conn)
                    log.Println("server close conn")
                case  msg := <- server.Messages :
                            for _, visitor := range server.Visitors {
                                select{
                                    case  visitor.Message <- msg :
                                    log.Println("send message")
                                default:
                                    visitor.CloseInfo <- 2
                                    log.Println("Server read and write failed ")
                                }
                            }                  
            }
        }
    }
    type Visitor struct{
    Conn        net.Conn
    Name        string
    Server      *Server

    Message      chan string
    ClientReader *bufio.Reader
    CloseInfo    chan int
    Closed       chan struct{}
    }
    func (visitor Visitor)Run(){
        //read from client
        go func(){
            for{
             select{
                 case <- visitor.Closed :
                    return
                default :
                log.Println("read")
                    msg, err :=  visitor.ClientReader.ReadString('\n')
                    if err != nil {
                        log.Println("visitor read client err ")
                        visitor.CloseInfo <- 0
                        return
                    }
                    visitor.Server.Messages <- msg
                    //log.Println(msg)
                    


             }
         }
        }()
        go func (){
            for{
                select{
                  case <- visitor.Closed :
                    return
                  case msg := <- visitor.Message :
                 // log.Println("visitor", visitor)
                 // log.Println("visitor.Conn", visitor.Conn)
                    _, err := visitor.Conn.Write([]byte(msg))   
                    if err != nil {
                        log.Println("visitor writer err")
                        visitor.CloseInfo <- 1
                        return
                    }
                }
            }
        }()
        defer func(){
           log.Println("%s exit", visitor.Name)
           visitor.Conn.Close()
            visitor.Server.CloseConn <- visitor.Conn
        }()
         <- visitor.CloseInfo
         close(visitor.Closed)
    }
 