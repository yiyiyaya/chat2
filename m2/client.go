package m2



import (

	"net"
	"github.com/asiainfoLDP/datahub_commons/log"
	//"bufio"
	"fmt"
    //"github.com/asiainfoLDP/datahub_commons/log"  
	"bufio"
	"os"  
	"strings"
)

func ClienttosServer(){
   /* info, _ := net.InterfaceAddrs()
    for _, addr := range info{
        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback(){
            if ipnet.IP.To4() != nil{
                ip :=ipnet.IP.String()
            }
        }
    }*/
    conn, err := net.Dial("tcp", "localhost:9999")
    if err != nil {
        log.DefaultlLogger().Error("client conn failed err")
        return
    }
    fmt.Println("client conn success")
    defer conn.Close()
   
    go ReadFromServer(conn)
     WriteToServer(conn)
     
   

    
    
}
/*func Write(conn net.Conn){
    defer conn.Close()
    for{
        //data := make([]byte,1024)
        data, err := bufio.NewReader(conn).ReadString('\n')
        //c, err := conn.Read(data)
        if err != nil {
            fmt.Println("read server write err")
            return
        }
        fmt.Println(data)
    }
}
func Read(conn net.Conn){
    defer conn.Close()
    for {

    }
}*/
func WriteToServer(conn net.Conn){
    
     reader := bufio.NewReader(os.Stdin)
        
    //content := string(data)
    //fmt.Println(content)
        for {
           // var talkContent string
            data, _, err := reader.ReadLine()
            if err != nil {
                log.DefaultlLogger().Error("ReadFromConsole err")
                return
            }
            data =append(data, "\n"...)
            _, err = conn.Write(data)
            //log.DefaultlLogger().Info("cvvvvvvbhjjjjjjjhjbvfvukbjfgbjfbvjhbfv")
            if err != nil {
                log.DefaultlLogger().Error("write to server err")
                return
            }
        }
        

}
func ReadFromServer(conn net.Conn){
    for{
        data, err := bufio.NewReader(conn).ReadString('\n')
        //c, err := conn.Read(data)
        if err != nil {
            fmt.Println("read server write err")
            return
        }
        ip := strings.Split(conn.LocalAddr().String(), ":")[0]
        fmt.Printf("%s saied:%s", ip, data)

    }
   
}


