- 修改每台主机hostname 统一三台机器 u p

  - master: 192.168.1.5
  - node1: 192.168.1.8
  - node2: 192.168.1.7

- 执行sealos run 

  - ```shell
    sealos run registry.cn-shanghai.aliyuncs.com/labring/kubernetes:v1.29.7 registry.cn-shanghai.aliyuncs.com/labring/helm:v3.9.4 registry.cn-shanghai.aliyuncs.com/labring/cilium:v1.13.4 \
         --masters 192.168.1.5 \
         --nodes 192.168.1.8,192.168.1.7 -u zf -p root
    ```

- 查看各节点 /etc/host 文件 相应k8s地址解析已添加

- 执行 kubectl get nodes ||  kubectl get pod -A 

  - **踩坑1：kubectl权限问题**

  - 针对配置信息，`kubectl` 在 `$HOME/.kube` 目录中查找一个名为 `config` 的配置文件。 你可以通过设置 `KUBECONFIG` 环境变量或设置 [`--kubeconfig`](https://kubernetes.io/zh-cn/docs/concepts/configuration/organize-cluster-access-kubeconfig/)：[命令行工具 (kubectl) | Kubernetes](https://kubernetes.io/zh-cn/docs/reference/kubectl/)

    参数来指定其它 [kubeconfig](https://kubernetes.io/zh-cn/docs/concepts/configuration/organize-cluster-access-kubeconfig/) 文件
  
    ```
    kubectl get pod -A    E0908 11:27:07.185600    6191 memcache.go:265] couldn't get current server API group list: Get "http://localhost:8080/api?timeout=32s": dial tcp [::1]:8080: connect: connection refused
    E0908 11:27:07.186618    6191 memcache.go:265] couldn't get current server API group list: Get "http://localhost:8080/api?timeout=32s": dial tcp [::1]:8080: connect: connection refused
    E0908 11:27:07.187308    6191 memcache.go:265] couldn't get current server API group list: Get "http://localhost:8080/api?timeout=32s": dial tcp [::1]:8080: connect:
    ```

​		解决办法：cp /etc/kubernetes/kubelet.conf  ~/.kube/config 拷贝到用户目录下的config中，同时执行

​		export KUBECONFIG=/home/zf/.kube/config，按道理如果$HOME/.kube没有config文件 可以手动创建，然后拷贝实际的kubeconfig文件过来也可以，不用刻意通过 export修改 环境变量。

