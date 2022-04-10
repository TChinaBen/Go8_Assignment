package check

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

func NewControls(t NodeType,in []byte)(*Controls,error){
	c := new(Controls)

	err := yaml.Unmarshal(in,c)

	if err != nil{
		fmt.Println("yaml文件转为结构体control失败")
		return nil,err
	}
	return c,nil
}


