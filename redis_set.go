package main

import (
"github.com/garyburd/redigo/redis"
"fmt"
)


func main()  {
    conn,err := redis.Dial("tcp","192.168.1.86:6379")
    if err != nil {
        fmt.Println("connect redis error :",err)
        return
    }
    defer conn.Close()
    _, err = conn.Do("SET", "name", "ctd")
    if err != nil {
        fmt.Println("redis set error:", err)
    }
}
