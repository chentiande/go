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
    name, err := redis.String(conn.Do("GET", "name"))
    if err != nil {
        fmt.Println("redis get error:", err)
    } else {
        fmt.Printf("Got name: %s \n", name)
    }
}
