package cmd

var (
	masterFile = "master.yaml"
	nodeFile="node.yaml"
	etcdFile = "etcd.yaml"
	cfgDir = "./configs/"
	envVarsPrefix = "KUBE_BENCH"
	configFileError error
)
var TypeMap = map[string][]string{
	"ca": {"cafile","defaultcafile"},
	"kubeconfig": {"kubeconfig","defaultkubeconfig"},
	"service": {"svc","defaultsvc"},
	"config": {"confs","defaultconf"},
}