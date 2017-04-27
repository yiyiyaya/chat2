package login

import(
    "fmt"
    "net/http"
    "strings"
    "github.com/julienschmidt/httprouter"
    "regexp"
    "log"
)
func login(username, password string){
    reg := `[\w\d\_]{6.10}`
    match, err := regexp.MatchString(reg, username)
    if err != "" {
    log.Println("username match err")    
    return
    }
    if match == false {
        fmt.Println("用户名是6-10位字母，数字，下划线")
        return
    }
    

}
