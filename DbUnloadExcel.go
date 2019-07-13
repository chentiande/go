package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
    "io/ioutil"

    _ "github.com/mattn/go-oci8"
    "github.com/360EntSecGroup-Skylar/excelize"
)

func query() {
    os.Setenv("NLS_LANG", "Simplified Chinese_China.AL32UTF8")
    log.SetFlags(log.Lshortfile | log.LstdFlags)
    f := excelize.NewFile()
    //db, err := sql.Open("oci8", "system/system@192.168.1.86:1521/ORCLCDB")
     str:=os.Args[1]
     str=strings.Replace(str, "\\", "", -1)
    db, err := sql.Open("oci8", str)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
     b, err := ioutil.ReadFile(os.Args[2])
    if err != nil {
        fmt.Println("ioutil ReadFile error: ", err)
        return
    }
    fmt.Printf("sql=%v",string(b))
    rows, err := db.Query(strings.Replace(string(b),";","",-1))
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
     if (j-2)%10000==0 {
       fmt.Println("导出"+strconv.Itoa(j-2)+"行")
       
      }
//        fmt.Printf("%s\t%s\n", result[0],result[1])
    }
      fmt.Println("总数据量为"+strconv.Itoa(j-2)+"行")
    rows.Close()
       fmt.Println("开始生成EXCEL，请耐心等待.....")
    err = f.SaveAs(os.Args[3])
    if err != nil {
        fmt.Println(err)
    }else{
       fmt.Println("导出数据成功")    
         }
}

func main() {
     if len(os.Args)!=4{
     fmt.Println("参数设置不正确,需要包含三个参数:连接串,包含sql的文件名,保存的excel文件名")
     fmt.Println("./DbUnloadExcel user/passwd@ip:port/service_name objects.sql object.xlsx")
     fmt.Println("如有其它问题请联系chentiande@boco.com.cn")
     }else {
  
     fmt.Println("开始导出数据......")
    query()
    } 
}
