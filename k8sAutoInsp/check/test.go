package check

import (
	"strconv"
	"strings"
)


func testItemOp(test testItem,actualValue string)bool{
	var res bool
	op := test.Compare.Op
	op = strings.TrimSpace(op)
	if len(op)>0{
		switch op {
		case "bitmask":
			requested,_ := strconv.ParseInt(test.Compare.Value,8,64)
			max,_ := strconv.ParseInt(actualValue,8,64)
			res = (requested & max) == max
		case "eq":
			val := strings.ToLower(actualValue)
			res = val==Replace_n(test.Compare.Value)
		case "noteq":
			res = !(actualValue==Replace_n(test.Value))
		case "nothave":
			res = !strings.Contains(actualValue,test.Compare.Value)
		case "has":
			res = strings.Contains(actualValue,Replace_n(test.Compare.Value))
		case "gt","gte","lt","lte":
			a,b,_:= toNumeric(actualValue,Replace_n(test.Compare.Value))
			switch op {
			case "gt":
				res = a>b
			case "gte":
				res = a>=b
			case "lt":
				res = a<b
			case "lte":
				res = a<=b
			}
		case "valid_elements":
			s := splitAndRemoveLastSeparator(actualValue,",")
			target := splitAndRemoveLastSeparator(test.Compare.Value,",")
			res = allElementsValid(s,target)

		}
	}

	return res
}

func toNumeric(a,b string)(c,d int,err error){
	c,_ = strconv.Atoi(strings.TrimSpace(a))
	d,_ = strconv.Atoi(strings.TrimSpace(b))
	return c,d,nil
}

func splitAndRemoveLastSeparator(s,sep string) []string{
	cleanS := strings.TrimRight(strings.TrimSpace(s),sep)
	if len(cleanS) == 0{
		return []string{}
	}
	ts := strings.Split(cleanS,sep)
	for i := range ts{
		ts[i] = strings.TrimSpace(ts[i])
	}
	return ts
}
func allElementsValid(s,t []string)bool{
	sourceEmpty := len(s) == 0
	targetEmpty := len(t) == 0
	if (sourceEmpty || targetEmpty) && !(sourceEmpty && targetEmpty){
		return false
	}
	for _,sv := range s{
		found := false
		for _,tv := range t{
			if sv==tv{
				found=true
				break
			}
		}
		if !found{
			return false
		}
	}
	return true
}

func (test tests)execute(auditOutput string)bool{
	    var result bool
		lens := len(test.TestItems)
		res := make([]bool,lens)
		for i,item := range test.TestItems{
			res[i] = item.execute(auditOutput)
		}

		if len(res) == 1{
			return res[0]
		}else{
			switch test.BinOp {
			case "and":{
				result = true
				for i := range res{
					result = result && res[i]
				}
			}
			case "or": {
				result = false
				for i := range res{
					result = result || res[i]
				}
			}

			}

		}

	return result
}
func (item testItem) execute(auditOutput string)bool {
	if len(item.Set) > 0 && item.Set == "false" {
		if strings.Contains(auditOutput, Replace_n(item.Flag)) {
			return false
		} else {
			return true
		}
	} else {
		// 拥有多行输出
		if strings.Count(auditOutput, "/n") >= 2 {
			return checkMultiOutput(auditOutput,item)
		}else{
			return checkNotMultiOutput(auditOutput,item)
		}
	}
}

func checkMultiOutput(auditOutput string,item testItem)bool{
	// 声明输出切片
	var sta = true
	var flag bool
	var output []string
	output = strings.Split(auditOutput,"\n")
	for i:=0;i<len(output)-1;i++{
		flag = evaluate(item,output[i])
		if !flag{
			sta = false
		}
	}

	return sta
}

func checkNotMultiOutput(auditOutput string,item testItem)bool{
	// 声明输出切片
	var flag bool
	flag = evaluate(item,auditOutput)
	return flag

}
