**重启虚拟机后 node节点无法访问api-server**

在master节点执行“kubectl get node，发现控制平面 not ready

执行： journalctl -f -u kubelet.service 查看kubelet日志， cni plugin not initialized （正常情况下就是镜像拉去或者启动失败 原因有很多种 我遇到的只是磁盘没有空间 拉取镜像）

具体查看pod情况：kubectl get po -A

 cilium-jnqbk  拉取镜像失败

执行：kubectl describe pod cilium-jnqbk -n kube-system 查看具体 pod 日志

发现：no spaces left

![image-20240910224146332](C:\Users\ASUS\AppData\Roaming\Typora\typora-user-images\image-20240910224146332.png)

最终问题虚机默认分配的存储空间太小，准备扩容

扩容后 重启kubelet



**如何对vmware虚拟机中的Linux系统进行扩容并将扩大的空间应用在linux中**
具体链接参考：https://cloud.tencent.com/developer/article/2395337

首先在VMware中对虚拟机进行扩容操作，如图，虚拟机必须关机才可以进行“扩展”，我的原先为8G，要扩展到13G（此时截屏为扩展后） 



![img](https://developer.qcloudimg.com/http-save/yehe-5814261/0b470e256b0104ccf8686fb257ab2612.png)

 扩展需要一段时间，扩展成功后启动虚拟机 fdisk -l 命令查看分区情况，可以看到框1中/dev/sda已经拥有了扩大的空间，但下面的框中并没有展示出扩大的空间，是因为还没有分区，还不能使用。 

![img](https://developer.qcloudimg.com/http-save/yehe-5814261/08ecec50b1849a9d90284a22b690b8ae.png)

 接下来使用Linux的fdisk分区工具给磁盘/dev/sda分区,命令如下

代码语言：javascript

复制

```javascript
fdisk /dev/sda
```

可以根据提示输入m查看帮助信息，在这里输入n(增加分区)，回车后输入p(创建主分区)，回车后partition number输入3(因为上面已经有两个分区sda1和sda2)，回车会提示输入分区的start值(通过fdisk -l 可以看出sda2的end值为27262975)，我们可以根据提示指定start值为16777216，end值为默认即可(即当前最大值)，回车后输入W进行保存，分区划分完毕。 

![img](https://developer.qcloudimg.com/http-save/yehe-5814261/493f65f94803ac0d666e5891c1a0030b.png)

![img](https://developer.qcloudimg.com/http-save/yehe-5814261/fc840b7f043fe6be23229869cc44774a.png)

 可以看到/dev/sda3的Id号为83，我们要将其改成8e(LVM卷文件系统的Id)，具体方法输入fdisk /dev/sda,选择t（改变一个分区的系统ID）回车，然后选择分区3回车，然后输入L回车。然后输入8e回车，然后输入w，保存修改的分区信息。最后输入fdisk -l ,查看ID是否修改成功。修改成功后必须重新启动linux系统才能进行后面的操作。 

![img](https://developer.qcloudimg.com/http-save/yehe-5814261/373ed50fe1b2e0d5f03fec4ed5853299.png)

 如下图，修改成功 

![img](https://developer.qcloudimg.com/http-save/yehe-5814261/576abad16c0f799c7dcb36fc2395af5f.png)

 系统重启后，格式化新的分区为ext4格式。命令如下

代码语言：javascript

复制

```javascript
mkfs.ext4 /dev/sda3
```

![img](https://developer.qcloudimg.com/http-save/yehe-5814261/94a25e3a67d4190a8fe7f9ec9b735f85.png)

 格式化后，创建PV. 用`pvdisplay`查看当前的物理卷。然后用pvcreate指令用于将物理硬盘分区初始化为物理卷，以便被LVM使用。要创建物理卷必须首先对硬盘进行分区，并且将硬盘分区的类型设置为“8e”后，才能使用pvcreat指令将分区初始化为物理卷。`pvcreate /dev/sda3`创建完后，我们可以再用`pvdisplay`查看到新创建的物理卷。 

![img](https://developer.qcloudimg.com/http-save/yehe-5814261/84a0657f43b8278fe828e0ee19973613.png)

![img](https://developer.qcloudimg.com/http-save/yehe-5814261/45e37a171934c68914828853dcab199a.png)

 扩展VG：当前需要查看扩充的lvm组名，可以通过`vgdisplay`查看，在此例中我们的组名为 centos,并可以看到里面的空间只有20多G。然后用vgextend指令用于动态的扩展卷组，它通过向卷组中添加物理卷来增加卷组的容量。 `vgextend centos /dev/sda3` 添加成功后，我们可以用`vgdisplay`再次查看，容量已经添加进去。 

![img](https://developer.qcloudimg.com/http-save/yehe-5814261/606622f00233d0aad9d77523c5f8379a.png)

![img](https://developer.qcloudimg.com/http-save/yehe-5814261/ec203faeef6a6cda06c1475efdda1d23.png)

 用`lvextend -L+4.98G /dev/centos/root /dev/sda3` 命令扩展空间到root下，扩容的空间要略小于VG的free空间，因此这里只输入了4.98G.然后通过`df -h`查看，root空间并没有改变，因为还差一步

代码语言：javascript

复制

```javascript
lvextend -L+4.98G /dev/centos/root /dev/sda3
```

![img](https://developer.qcloudimg.com/http-save/yehe-5814261/b6dfdc72d1c9e2cb55ba2b09b920ed08.png)

![img](https://developer.qcloudimg.com/http-save/yehe-5814261/a705819550e6aefdf875aa5b9a377bd9.png)

 使用使用resize2fs命令，用于扩大或者缩小未挂载的ext2,ext3或者是ext4文件系统。具体命令为：resize2fs -p /dev/mapper/VolGroup-lv_root 290G 。然后再用df -h 查看

代码语言：javascript

复制

```javascript
resize2fs -p /dev/mapper/centos-root
```

此时报错， 

![img](https://developer.qcloudimg.com/http-save/yehe-5814261/74088415ca35ae45a47d78270aeb507d.png)

 不是resize2fs我们可以用xfs

代码语言：javascript

复制

```javascript
  xfs_growfs /dev/mapper/centos-root
```

![img](https://developer.qcloudimg.com/http-save/yehe-5814261/84fa5531e3ccd916880652aa2f594d9f.png)

 成功后用下面的命令查看

代码语言：javascript

复制

```javascript
df -h
```

![img](https://developer.qcloudimg.com/http-save/yehe-5814261/e658bbb4a0068e2634425108e380ec83.png)