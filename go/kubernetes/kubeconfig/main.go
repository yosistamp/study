package main

import (
	"fmt"
	"log"
	"os"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	types "k8s.io/client-go/tools/clientcmd/api"
)

func main() {
	fileName := "./kubeconfig"
	config, err := readConfig(fileName)
	if err != nil {
		log.Fatalln(err)
	}

	writeConfig(*config, "kube1", "https://xxxx:6443/", "xxxxxxxxxxxxxxxxxxxxxx")
}

func readConfig(fileName string) (*clientcmdapi.Config, error) {
	if exists(fileName) {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		loadingRules.ExplicitPath = fileName

		kubeconfig, err := loadingRules.Load()
		if err != nil {
			return kubeconfig, err
		}

		return kubeconfig, nil
	}

	clusters := make(map[string]*clientcmdapi.Cluster)
	contexts := make(map[string]*clientcmdapi.Context)
	auths := make(map[string]*clientcmdapi.AuthInfo)
	kubeconfig := &types.Config{
		Clusters:       clusters,
		Contexts:       contexts,
		AuthInfos:      auths,
		CurrentContext: "",
	}
	return kubeconfig, nil
}

// kubeconfigの書き込み
// 同じキーがあればマップを更新、なければ新規項目として追加する
func writeConfig(config clientcmdapi.Config, id string, endpoint string, cert string) {
	clusterName := fmt.Sprintf("%s-cluster", id)
	authInfoName := fmt.Sprintf("%s-user", id)
	contextName := fmt.Sprintf("%s@%s", authInfoName, id)

	cluster := &types.Cluster{
		Server:                   endpoint,
		CertificateAuthorityData: []byte(cert), // Base64エンコードされるので、エンコード済みの場合はデコードしてセット
	}
	config.Clusters[clusterName] = cluster

	context := &types.Context{
		Cluster:  clusterName,
		AuthInfo: authInfoName,
	}
	config.Contexts[contextName] = context

	authInfo := &types.AuthInfo{
		ClientCertificateData: []byte(cert), // Base64エンコードされるので、エンコード済みの場合はデコードしてセット
		ClientKeyData:         []byte(cert), // Base64エンコードされるので、エンコード済みの場合はデコードしてセット
	}
	config.AuthInfos[authInfoName] = authInfo

	if err := clientcmd.WriteToFile(config, "./kubeconfig"); err != nil {
		log.Fatalln(err)
	}

}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
