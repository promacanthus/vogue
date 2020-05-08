---
title: 49-xml
date: 2020-01-10T20:37:47.913857+08:00
draft: false
---

```go
// Go使用encoding.xml包提供对XML和类XML格式的内置支持。

package main

import (
	"encoding/xml"
	"fmt"
)

// Plant 结构体将被映射为XML
// 与JSON示例类似，字段标记包含编辑器和解码器的指令
// 使用XML包的一些特殊功能：
//  1. XMLName字段名称表示此结构的XML元素的名称
//  2. id,attr表示id字段是XML的属性而不是嵌套元素
type Plant struct {
	XMLName xml.Name `xml:"plant"`
	ID      int      `xml:"id,attr"`
	Name    string   `xml:"name"`
	Origin  []string `xml:"origin"`
}

func (p Plant) String() string {
	return fmt.Sprintf("Plant id = %v, name = %v, origin = %v", p.ID, p.Name, p.Origin)
}

func main() {
	coffee := &Plant{ID: 27, Name: "Coffee"}
	coffee.Origin = []string{"Ethiopia", "Brazil"}

	//使用MarshalIndent来生成更易读取的输出
	out, _ := xml.MarshalIndent(coffee, " ", " ")
	fmt.Println(string(out))

	// 要将通用XML标头添加到输出，需要显式添加它
	fmt.Println(xml.Header + string(out))

	// 使用Unmarshal将带有XML的字节流解析为数据结构
	// 如果XML格式错误或无法映射到Plant，将返回错误性描述
	var p Plant
	if err := xml.Unmarshal(out, &p); err != nil {
		panic(err)
	}
	fmt.Println(p)

	tomato := &Plant{ID: 81, Name: "Tomato"}
	tomato.Origin = []string{"Mexico", "California"}

	type Nesting struct {
		XMLName xml.Name `xml:"nesting"`
		Plants  []*Plant `xml:"parent>child>plant` // parent> child> plant field标记告诉编码器将所有植物嵌套在<parent> <child>下
	}

	nesting := &Nesting{}
	nesting.Plants = []*Plant{coffee, tomato}

	out, _ = xml.MarshalIndent(nesting, " ", " ")
	fmt.Println(string(out))
}

```