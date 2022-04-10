package check

import (
	"bufio"
	"fmt"
	"os"
)

func WriteFile(content string){
	var filename = "result.txt"
	file,err := os.OpenFile(filename,os.O_CREATE|os.O_APPEND|os.O_RDWR,0666)
	if err != nil{
		fmt.Println("Error:",err.Error())
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.Write([]byte(content+"\n"))
	defer writer.Flush()
}
