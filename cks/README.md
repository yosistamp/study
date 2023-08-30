# CKS

## CIS Benchmark

### インタラクティブモードの実行

```
./Assessor -i
```

### 出力指定

```
./Assessor-CLI.sh -i -rd /var/www/html/ -nts -rp index
```

あとは、レポートを表示して、指示された修正を行う。

## kube-bench

kubernetes関連のセキュリティ診断

## AppArmor

### ロード

以下のコマンドでロードすればよい。ファイルはどこにあってもいい。

```
apparmor_parser -q　<ファイル>
```

結果は aa-status で確認できる

### Podへの適用

アノテーションで適用する。提要場所はPodのmetadataになることに注意。（Deploymentのmetadataではない）

ローカルのプロファイル指定はlocalhost/<プロファイル名>とする。

```
metadata:
    annotations:
      container.apparmor.security.beta.kubernetes.io/nginx: localhost/restricted-nginx
```

## Falco

ランタイムの脆弱性検出

### ログの確認

```
journalctl -fu falco
```

### ファイル出力

/etc/falco/falco.yaml

```
file_output:
  enabled: true
  keep_alive: false
```

### /etc/falco/falco_rules.local.yamlの修正

デフォルトを上書きせずに.localに追加する

```
- rule: Write below binary dir
  desc: an attempt to write to any file below a set of binary directories
  condition: >
    bin_dir and evt.dir = < and open_write
    and not package_mgmt_procs
    and not exe_running_docker_save
    and not python_running_get_pip
    and not python_running_ms_oms
    and not user_known_write_below_binary_dir_activities
  output: >
    File below a known binary directory opened for writing (user_id=%user.uid file_updated=%fd.name command=%proc.cmdline)
  priority: CRITICAL
  tags: [filesystem, mitre_persistence]
```

### 再起動

```
service falco restart
or
systemctl restart falco
```

### ルールの変更

問題の条件に合致するルールを探して、Falcoのドキュメントを見ながらoutputのフォーマットを変更する

```
- rule: Terminal shell in container
  desc: A shell was used as the entrypoint/exec point into a container with an attached terminal.
  condition: >
    spawned_process and container
    and shell_procs and proc.tty != 0
    and container_entrypoint
    and not user_expected_terminal_shell_in_container_conditions
  output: >
    %evt.time.s,%user.uid,%container.id,%container.image.repository
  priority: ALERT
  tags: [container, shell, mitre_execution]
```

## RuntimeClass

ランタイムクラスの確認

```
k get runtimeclass
```

Podへの適用

```
apiVersion: v1
kind: Pod
metadata:
  name: mypod
spec:
  runtimeClassName: myclass
```

## 秘密情報の漏洩検出

Podにマウントされている環境変数やファイルなどから検出する

もしくはマウントされているServiceAccountのトークンを使ってシークレットにAPIアクセスする

### Pod内からアクセスする手順

```
curl https://kubernetes.default/api/v1/namespaces/restricted/secrets -H "Authorization: Bearer $(cat /run/secrets/kubernetes.io/serviceaccount/token)" -k
```

## certificate

### apiserver

/etc/kubernetes/pki/apiserver.crt
/etc/kubernetes/manifests

### CRI log

crictl ps -a

## kubectl

- kubectl proxy
  - APIサーバーへのアクセスをプロキシする
- kubectl port-forwrad
  - デプロイされているPodへのアクセスをプロキシする

## upgrade 

https://v1-26.docs.kubernetes.io/docs/tasks/administer-cluster/kubeadm/kubeadm-upgrade/

kubeadm upgrade

## systemd

/lib/systemd/system配下

## ポートを調べる

netstat -natp

## firewall

ufw allow/deny from <CIDR> to <CIDR>|any port <PORT> proto <tcp/udp>

## systemcall trace

strace

## security context

securityContext:
  runAsUser: 1000
  runAsGroup: 3000
  fsGroup: 2000
  fsGroupChangePolicy: "OnRootMismatch"

## seccomp

### デフォルトパス

```
/var/lib/kubelet/seccomp/profiles/
```

### Podに反映

```
spec:
  securityContext:
    seccompProfile:
      type: Localhost
      localhostProfile: profiles/audit.json
```

## securityContext

Podレベル（specへの指定）とコンテナレベル（containersへの指定）ができる。
capabilitiesはコンテナ指定のみ。

## admission controller 

kube-apiserverのパラメータ確認

```
ps -ef | grep kube-apiserver | grep admission-plugins
```

/etc/kubernetes/manifests/kube-apiserver.yaml

```
    - --enable-admission-plugins=NodeRestriction,NamespaceAutoProvision
    - --disable-admission-plugins=DefaultStorageClass
```

## OPA

サーバーモード

```
opa run -s
```

ポリシー適用

```
curl -X PUT --data-binary @sample.rego http://localhost:8181/v1/policies/samplepolicy
```

## secret

account-token
opaque

```
      envFrom:
      - secretRef:
          name: mysecret
```


## registry whitelist


## kubesec

Kubernetesマニフェストを静的解析するツール

kubesec scan <json|yaml>

結果のスコアを見てマイナス部分を修正する

## trivy

trivy image xxx 

dockerコマンドが使えないときはcrictlを確認する


## immutable

ファイルシステムへの変更を防ぐにはsecurityContextでreadOnlyRootFilesystemを指定する。

```
containers:
  securityContext:
    readOnlyRootFilesystem: true
```

ユーザーとグループの指定

```
containers:
  securityContext:
    runAsUser: <uid>
    runAsGroup: <gid>
```


## audit log

kube-apiserverが管理する

ResponseComplete

RequestResponse


### kube-apiserver.yamlの修正

ポリシーファイルとログファイルの定義

```
 - --audit-policy-file=/etc/kubernetes/prod-audit.yaml
 - --audit-log-path=/var/log/prod-secrets.log
 - --audit-log-maxage=30
```

### ボリューム指定

```
  - name: audit
    hostPath:
      path: /etc/kubernetes/prod-audit.yaml
      type: File

  - name: audit-log
    hostPath:
      path: /var/log/prod-secrets.log
      type: FileOrCreate
```

### ボリュームマウント指定

yamlへはreadOnlyでパスを指定。ログファイルへはwriteできる権限でマウント

```
    volumeMounts:
    - mountPath: /etc/kubernetes/prod-audit.yaml
      name: audit
      readOnly: true
    - mountPath: /var/log/prod-secrets.log
      name: audit-log
      readOnly: false
```

## NetworkPolicy

```
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: allow-app1-app2
  namespace: apps-xyz
spec:
  podSelector:
    matchLabels:
      tier: backend
      role: db
  ingress:
  - from:
    - podSelector:
        matchLabels:
          name: app1
          tier: frontend
    - podSelector:
        matchLabels:
          name: app2
          tier: frontend
```

## OPA

## admission

admisstion controllerのkubeconfigfileはpkiフォルダ配下のものを使う

```
apiVersion: apiserver.config.k8s.io/v1
kind: AdmissionConfiguration
plugins:
- name: ImagePolicyWebhook
  configuration:
    imagePolicy:
      kubeConfigFile: /etc/kubernetes/pki/admission_kube_config.yaml
      allowTTL: 50
      denyTTL: 50
      retryBackoff: 500
      defaultAllow: false
```

有効化するにはenable-admission-pluginsへの追加とadmission-control-config-fileへのパス設定
```
   - --enable-admission-plugins=NodeRestriction,ImagePolicyWebhook
   - --admission-control-config-file=/etc/kubernetes/pki/admission_configuration.yaml
```

## 不変

readOnlyRootFilesystem: true指定を明示的に指定しているものが不変

```
Pod solaris is immutable as it have readOnlyRootFilesystem: true so it should not be deleted.

Pod sonata is running with privileged: true and triton doesn't define readOnlyRootFilesystem: true so both break the concept of immutability and should be deleted.
```

## default serviceaccount

どんな権限になっている？


## etcdctl

kubectl exec -n kube-system -it <etcd pod> -- sh -c "ETCDCTL=3 etcdctl --endpoints IP:PORT --cert=<CERT> --key=<KEY> --cacert=<CACERT> get / --prefix --keys-only=true

## システムコールを取得する

strace -p <プロセス>

## Service

Internal Access のみ → ClusterIP

## kubeletの設定

/var/lib/kubelet/config.yaml

## APIServerのNodePortからClusterIPへの切り替え

- APIServerの--kubernetes-service-node-portを0にする
- k delete svc kubernetes でAPIServerのServiceを一度削除する（自動で再作成されたときにClusterIPで作られる）

## crictlでコンテナIDからPod特定

```
crictl ps -id <コンテナID>
crictl pods -id <POD ID>
```

