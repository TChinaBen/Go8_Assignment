package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"k8sAutoInsp/check"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func getFilePath(bv string,filename string)(path string,err error){
	path = filepath.Join(cfgDir,bv)
	filepath := filepath.Join(path,string(filename))

	_, err1 := os.Stat(filepath)
	if err1 != nil{
		fmt.Println("没有找到这个文件")
		return
	}
	return path,nil
}

// 获取可执行的命令
func getBinaries(v *viper.Viper,nodetype check.NodeType)(map[string]string){
	binmap := make(map[string]string)
	// component 包含 apiServer scheduler controller-manager etcd
	for _,component := range v.GetStringSlice("components"){
		s := v.Sub(component)
		if s == nil{
			continue
		}
		bins := s.GetStringSlice("bins")
		if len(bins) > 0 {
			bin:= findExecutable(bins)
			binmap[component] = bin
		}
	}
	return binmap
}

func findExecutable(bins []string)string{
	var res string
	for _,bin := range bins{
		bin = strings.Trim(bin,"'\"")
		proc := strings.Fields(bin)[0]
		out := psFunc(proc)
		reFirstWord := regexp.MustCompile(`^(\S*\/)*`+bin)
		lines := strings.Split(out,"\n")
		for _,l := range lines{
			if reFirstWord.Match([]byte(l)){
				res = bin
				break
			}
		}

	}

 return res
}

// 该函数等同于命令ps -ef | grep proc
func psFunc(proc string) string{
	cmd := exec.Command("/bin/ps","-C",proc,"-o","cmd","--no-headers")
	out,_ := cmd.Output()
	return string(out)
}
