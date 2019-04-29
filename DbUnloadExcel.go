package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "strconv"
    "io/ioutil"

    _ "github.com/mattn/go-oci8"
    "github.com/360EntSecGroup-Skylar/excelize"
)

func query() {
    os.Setenv("NLS_LANG", "AMERICAN_AMERICA.AL32UTF8")
    log.SetFlags(log.Lshortfile | log.LstdFlags)
    f := excelize.NewFile()
    //db, err := sql.Open("oci8", "system/system@192.168.1.86:1521/ORCLCDB")
    db, err := sql.Open("oci8", os.Args[1])
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
     b, err := ioutil.ReadFile(os.Args[2])
    if err != nil {
        fmt.Println("ioutil ReadFile error: ", err)
        return
    }
    rows, err := db.Query(string(b))
    if err != nil {
        log.Fatal(err)
    }
    cols, _ := rows.Columns()
    rawResult := make([][]byte, len(cols))
    result := make([]string, len(cols))
    dest := make([]interface{}, len(cols))
    for i := range rawResult {
        dest[i] = &rawResult[i]
    }
    for k := range result{
      f.SetCellValue("Sheet1", string(65+k)+strconv.Itoa(1), cols[k])
    }
    j:=2
    for rows.Next() {
        err = rows.Scan(dest...)
        for i, raw := range rawResult {
            if raw == nil {
                result[i] = ""
            } else {
                result[i] = string(raw)
            }
        }
    for m := range result{
      f.SetCellValue("Sheet1", string(65+m)+strconv.Itoa(j), result[m])
    }
        j++
     if j%10000==0 {
       fmt.Println("导出"+strconv.Itoa(j)+"行")
       
      }
//        fmt.Printf("%s\t%s\n", result[0],result[1])
    }
      fmt.Println("总数据量为"+strconv.Itoa(j)+"行")
    rows.Close()
       fmt.Println("开始生成EXCEL，请耐心等待.....")
    err = f.SaveAs(os.Args[3])
    if err != nil {
        fmt.Println(err)
    }
}

func main() {
     fmt.Println("开始导出数据......")
    query()
      fmt.Println("导出数据成功")
     
}
