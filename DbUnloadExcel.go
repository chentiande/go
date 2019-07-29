package main

//从oracle数据库导出文件，可导出xlsx和csv两种格式
//作者：chentiande@boco.com.cn
//

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/axgle/mahonia"
	_ "github.com/mattn/go-oci8"
)

//增加GBK到utf8函数转换，将数据库取出的数据转成uft8然后保存到excel
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func querytocsv() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	f, err := os.Create(os.Args[3]) //创建文件
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)         //创建一个新的写入文件流
	str := os.Args[1]
	//解决密码中有特殊字符进行转义后去掉转义斜杠
	str = strings.Replace(str, "\\", "", -1)
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
	//sql语句中末尾如果有分号，自动去掉
	strsql := strings.Replace(string(b), ";", "", -1)
	fmt.Printf("sql=%v", strsql)
	//如果配置了第四个参数，将进行字符转换，将sql语句中的中文从utf8转为GBK，然后提交数据库
	if os.Args[4] == "Y" || os.Args[4] == "y" {
		enc := mahonia.NewEncoder("gbk")
		strsql = enc.ConvertString(strsql)
	}
	rows, err := db.Query(strsql)
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

	w.Write(cols) //写入表头数据
	w.Flush()
	j := 2
	for rows.Next() {
		err = rows.Scan(dest...)
		for i, raw := range rawResult {
			if raw == nil {
				result[i] = ""
			} else {
				//第四个参数确定是否数据库需要GBK到utf的转码
				if os.Args[4] == "Y" || os.Args[4] == "y" {
					result[i] = ConvertToString(string(raw), "gbk", "utf-8")
				} else {
					result[i] = string(raw)
				}
			}
		}

		w.Write(result) //写入数据

		j++
		if (j-2)%10000 == 0 {
			fmt.Println("导出" + strconv.Itoa(j-2) + "行")
			w.Flush()

		}
		//        fmt.Printf("%s\t%s\n", result[0],result[1])
	}
	w.Flush()
	fmt.Println("总数据量为" + strconv.Itoa(j-2) + "行")
	rows.Close()
	fmt.Println("导出数据成功")
}

func querytoxlsx() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	f := excelize.NewFile()
	//db, err := sql.Open("oci8", "system/system@192.168.1.86:1521/ORCLCDB")
	str := os.Args[1]
	//解决密码中有特殊字符进行转义后去掉转义斜杠
	str = strings.Replace(str, "\\", "", -1)
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
	//sql语句中末尾如果有分号，自动去掉
	strsql := strings.Replace(string(b), ";", "", -1)
	fmt.Printf("sql=%v", strsql)
	//如果配置了第四个参数，将进行字符转换，将sql语句中的中文从utf8转为GBK，然后提交数据库
	if os.Args[4] == "Y" || os.Args[4] == "y" {
		enc := mahonia.NewEncoder("gbk")
		strsql = enc.ConvertString(strsql)
	}
	rows, err := db.Query(strsql)
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
	for k := range result {
		f.SetCellValue("Sheet1", string(65+k)+strconv.Itoa(1), cols[k])
	}
	j := 2
	for rows.Next() {
		err = rows.Scan(dest...)
		for i, raw := range rawResult {
			if raw == nil {
				result[i] = ""
			} else {
				//第四个参数确定是否数据库需要GBK到utf的转码
				if os.Args[4] == "Y" || os.Args[4] == "y" {
					result[i] = ConvertToString(string(raw), "gbk", "utf-8")
				} else {
					result[i] = string(raw)
				}
			}
		}
		for m := range result {
			f.SetCellValue("Sheet1", string(65+m)+strconv.Itoa(j), result[m])
		}
		j++
		if (j-2)%10000 == 0 {
			fmt.Println("导出" + strconv.Itoa(j-2) + "行")

		}
		//        fmt.Printf("%s\t%s\n", result[0],result[1])
	}
	fmt.Println("总数据量为" + strconv.Itoa(j-2) + "行")
	rows.Close()
	fmt.Println("开始生成EXCEL，请耐心等待.....")
	err = f.SaveAs(os.Args[3])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("导出数据成功")
	}
}

func main() {
	if len(os.Args) != 5 {
		fmt.Println("参数设置不正确,需要包含四个参数:连接串,包含sql的文件名,保存的excel文件名,是否GBK编码")
		fmt.Println("./DbUnloadExcel user/passwd@ip:port/service_name objects.sql object.xlsx Y|N")
		fmt.Println("如有其它问题请联系chentiande@boco.com.cn")
	} else {

		fmt.Println("开始导出数据......")
		if os.Args[3][len(os.Args[3])-3:] == "csv" {
			querytocsv()
		} else {
			querytoxlsx()
		}
	}
}

