package main

import (
	"bytes"
	"fmt"
	"github.com/Juniper/ipfix-translator/tools/rfc5102"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	var (
		buffer bytes.Buffer
		keys   []string
	)
	buffer.WriteString("//This file is auto generated, do not modify\n")
	buffer.WriteString("package msghandler")
	buffer.WriteString("\n\n")
	buffer.WriteString("type IpfixDataSet struct {\n")
	for _, v := range rfc5102.InfoModel {
		keys = append(keys, v.Name)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, v := range rfc5102.InfoModel {
			if k == v.Name {
				datatype := string(v.Type)
				if datatype == "" {
					/* Type is unknown */
					datatype = "interface{}"
				}
				buffer.WriteString("\t" + strings.Title(v.Name) + "\t" +
					datatype + "\t")
				buffer.WriteString("`json:" + "\"" + v.Name + ",omitempty\"`")
				buffer.WriteString("\n")
			}
		}
	}
	buffer.WriteString("}")
	file, err := os.Create("../msg-handler/rfc5102_datasets.go")
	if err != nil {
		log.Fatal("Create rfc5102_datasets.go error ", err)
	}
	defer file.Close()
	fmt.Fprintf(file, buffer.String())
}
