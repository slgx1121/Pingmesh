# Pingmesh

##安装fping4.2
Ubuntu要先安装build-essential
apt-get install build-essential
下载安装包
wget https://fping.org/dist/fping-4.2.tar.gz
解压
tar zxvf fping-4.2.tar.gz
配置
cd fping-4.2/
./configure --prefix=/usr/local/fping
make
sudo make install
查看版本
/usr/local/fping/sbin/fping -V
设置环境变量
> vim /etc/profile
#在最后面添加
export PATH=$PATH:/usr/local/fping/sbin
> source /etc/profile
> fping -v

##V2
server通过rpc把pinglist数据传给client，client对服务器进行ping并产生结果
<img width="1081" alt="image" src="https://user-images.githubusercontent.com/121349317/210778343-c4d64fc4-be68-47b6-a7f6-6362ca77f638.png">
