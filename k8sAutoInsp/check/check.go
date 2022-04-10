package check

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func ElementTest(runner Runner,c *Controls){
	// 对于每一条命令Audit，这里都有输出结果
	for _,group := range c.Groups{
		for _,check := range group.Checks{
			runner.Run(check)
		}}
	c.RunChecks()

}

type Runner interface {
	Run(c *Check) State
}
type defaultRunner struct{}
func NewRunner() Runner {
	return &defaultRunner{}
}
func (r *defaultRunner) Run(c *Check) State{
	return c.run()
}
func (c *Check) run() State{
	if c.Tests == nil || len(c.Tests.TestItems) == 0{
		c.State = WARN
		content :=c.ID+":"+"there is no testItems"
		WriteFile(content)
		return c.State
	}

   err := c.runAuditCommands()
   if err == nil{

   }
   return c.State
}



func (c *Check)runAuditCommands()error{
	output,err := runAudit(c.Audit)
	c.AuditOutput = output
	return err
}

func runAudit(audit string)(output string,err error){
	var out bytes.Buffer
	audit = strings.TrimSpace(audit)
	audit = strings.Replace(audit,"\\u003e",">",-1)
	if len(audit) == 0{
		return output,err
	}
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = strings.NewReader(audit)
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	output = out.String()
    return output,err
}


//  在字符串s 里面找到匹配到t.flag的值，并将匹配到的结果输出到value
func (t testItem)findValue(s string)(match bool,value string,err error){
	if s == "" || t.Flag == ""{
		return
	}

	// match 里面是否含有flag这个值
	match = strings.Contains(s,t.Flag)
	if match{
		pttn := `(` + t.Flag +`)(=|: *)*([^\s]*) *`
		flagRe := regexp.MustCompile(pttn)
		vals := flagRe.FindStringSubmatch(s)
		if len(vals) > 0{
			if vals[3] != ""{
				value = vals[3]
			}else{
				if strings.HasPrefix(t.Flag,"--"){
					value = "true"
				}else{
					value = vals[1]
				}
			}
		}else{
			err = fmt.Errorf("invalid flag in testItem definition: %s",s)
		}
	}

	return match,value,err
}

func (controls *Controls)RunChecks(){
	for _,group := range controls.Groups{
		for _,check := range group.Checks{
			RunTestItems(*check)
		}
	}


}

func RunTestItems(check Check){
	var auditOutput string = check.AuditOutput
	var test = check.Tests
	res := test.execute(auditOutput)
	var tes testOutput
	tes.Audit = check.Audit

	if res{
		tes.ID = check.ID
		tes.State = PASS
	}else{
		var str string = check.Remediation
		Replace_n(str)
		tes.Remediation=str
		tes.ID = check.ID
		tes.State = FAIL
	}
	resjson,_ := json.Marshal(tes)
	WriteFile(string(resjson))
}




// 对于每一个testItem
func evaluate(item testItem,output string)bool{
	var flag = true
	var op = item.Compare.Op
	op = strings.TrimSpace(op)
	if len(op)>0{
		match,value,_ := item.findValue(output)
		value = Replace_n(value)
		if match{
			var res = testItemOp(item,value)
			if !res{
				flag = false
			}
		}else{
			flag = false
		}
	}else{
		flag = strings.Contains(output,item.Flag)
	}

	return flag
}




func Replace_n(remediation string)string{
	str := strings.Replace(remediation,"\n","",-1)
	return str
}








