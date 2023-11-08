package main

import (
	"context"
	"fmt"
	"log"

	"github.com/olivere/elastic/v7"
)

const (
	indexName = "students"
)

// Student represents a student with name and age.
type Student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	// 创建Elasticsearch客户端
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		log.Fatal(err)
	}

	// 创建索引（如果不存在）
	exists, err := client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		createIndex, err := client.CreateIndex(indexName).BodyString(mapping).Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		if !createIndex.Acknowledged {
			log.Fatal("Index creation not acknowledged")
		}
	}

	// 存储学生数据
	student := Student{Name: "John Doe", Age: 20}
	_, err = client.Index().Index(indexName).Id(student.Name).BodyJson(student).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Student data stored.")

	// 根据姓名查询年龄
	getResult, err := client.Get().Index(indexName).Id("John Doe").Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if getResult.Found {
		var foundStudent Student
		err := getResult.Source.Decode(&foundStudent)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Name: %s, Age: %d\n", foundStudent.Name, foundStudent.Age)
	}

	// 更新年龄
	updateAge := 25
	_, err = client.Update().Index(indexName).Id("John Doe").Doc(map[string]interface{}{"age": updateAge}).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Age updated to %d\n", updateAge)

	// 删除年龄
	_, err = client.Delete().Index(indexName).Id("John Doe").Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Student data deleted.")
}

const mapping = `
{
	"mappings": {
		"properties": {
			"name": {
				"type": "keyword"
			},
			"age": {
				"type": "integer"
			}
		}
	}
}
`
