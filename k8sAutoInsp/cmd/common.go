package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"k8sAutoInsp/check"
	"os"
	"path/filepath"
	"strings"
)

func loadConfig(nodeType check.NodeType,bv string) string{
	var file string
	switch nodeType {
	case check.MASTER:
		file = masterFile
	case check.NODE:
		file = nodeFile
	case check.ETCD:
		file =etcdFile
	}
	path,err := getFilePath(bv,file)
	if err!=nil{
		fmt.Println("查找文件失败")
	}
	mergeConfig(path)
	return filepath.Join(path,file)

}

func mergeConfig(path string){
	// 设置配置文件信息
	viper.SetConfigType("yaml")
	viper.SetConfigFile(path + "/config.yaml")
	// 读取配置文件
	//err := viper.ReadInConfig()
	err := viper.MergeInConfig()
	if err != nil{
		fmt.Println("合并文件失败")
		return
	}
}

func runChecks(nodetype check.NodeType,testYamlFile string){

	in,err := ioutil.ReadFile(testYamlFile)
	if err != nil{
		fmt.Println("读取文件失败")
	}
	typeConf := viper.Sub(string(nodetype))
	// binmap这个map里面存放的binmap[apiserver]=kube-apiserver
	binmap := getBinaries(typeConf,nodetype)
	confmap := getFiles(typeConf,"config")
	svcmap := getFiles(typeConf,"service")
	kubeconfmap := getFiles(typeConf,"kubeconfig")
	cafilemap := getFiles(typeConf,"ca")
	// 现在就需要把audit命令中的变量$apiserver改成系统中可用的变量
	s := string(in)
	s,_ = makeSubstitutions(s,"bin",binmap)
	s,_ = makeSubstitutions(s,"conf",confmap)
	s,_ = makeSubstitutions(s,"svc",svcmap)
	s,_ = makeSubstitutions(s,"kubeconfig",kubeconfmap)
	s,_ = makeSubstitutions(s,"cafile",cafilemap)
	// 现在controls 里面存放了yaml文件转化的结构体
	controls,err := check.NewControls(nodetype,[]byte(s))
	if err != nil{
		fmt.Println("yaml文件转化为结构体controls失败")
	}
	runner := check.NewRunner()
	check.ElementTest(runner,controls)

}
// s  "bin" binmap
func makeSubstitutions(s string,ext string,m map[string]string)(string,[]string){
	substitutions := make([]string,0)
	for k,v := range m{
		subset := "$"+ k + ext
		if v == ""{
			continue
		}
		beforeS := s
		s = multiWordReplace(s,subset,v)
		if beforeS != s{
			substitutions = append(substitutions,v)
		}
	}
	return s,substitutions
}

// 只要sub不为空，就将subname 替换为sub,此处n=-1,全部替换
func multiWordReplace(s string,subname string,sub string)string{
	f := strings.Fields(sub)
	if len(f) > 1{
		sub = "'" + sub + "'"
	}
	return strings.Replace(s,subname,sub,-1)
}

func getFiles(v *viper.Viper,fileType string) map[string]string{
	filemap := make(map[string]string)
	mainOpt := TypeMap[fileType][0]
	defaultOpt := TypeMap[fileType][1]
	for _,component := range v.GetStringSlice("components"){
		s := v.Sub(component)
		if s==nil{
			continue
		}
		file := findConfigFile(s.GetStringSlice(mainOpt))
		if file == ""{
			if viper.IsSet(defaultOpt){
				file = s.GetString(defaultOpt)
			}else{
				file = component
			}
		}else{

		}
		filemap[component] = file
	}
	return filemap
}

func findConfigFile(candidates []string)string{
	for _,c := range candidates{
		_,err := os.Stat(c)
		if err == nil {
			return c
		}
		if !os.IsNotExist(err){

		}
	}
	return ""
}
func isMaster() bool{
	return isThisNodeRunning(check.MASTER)
}
func isThisNodeRunning(nodeType check.NodeType) bool{
	nodeTypeConf := viper.Sub(string(nodeType))
	if nodeTypeConf == nil{
		return false
	}
	components := getBinaries(nodeTypeConf,nodeType)
	if len(components) == 0{
		return false
	}
	return true
}

func ValidTargets(bv string,targets []string,v *viper.Viper)(bool,error){
	benchmarkVersionToTargetsMap,err := loadTargetMapping(v)
	if err != nil{
		return false,err
	}
	providedTargets,found := benchmarkVersionToTargetsMap[bv]
	if !found{
		return false,fmt.Errorf("No targets configured for %s",bv)
	}
	for _,pt := range targets{
		f := false
		for _,t := range providedTargets{
			if pt == strings.ToLower(t){
				f = true
				break
			}
		}
		if !f{
			return false,nil
		}
	}
	return true,nil
}
func loadTargetMapping(v *viper.Viper)(map[string][]string,error){
	benchmarkVersionToTargetsMap := v.GetStringMapStringSlice("target_mapping")
	if len(benchmarkVersionToTargetsMap) == 0{
		return nil,fmt.Errorf("")
	}
	return benchmarkVersionToTargetsMap,nil
}