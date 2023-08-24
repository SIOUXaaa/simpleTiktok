# simpleTiktok

**本项目使用字节的Hertz框架，gorm**

数据库：

项目中配置好了mysql，可以使用docker起mysql，或者使用本地的mysql，手动新建数据库douyin，DSN可以在pkg/constants/constant.go中修改

推荐使用docker，后面可能会使用minio对象存储来处理视频流

## docker部署mysql

在docker-compose.yml文件中可以修改docker的端口映射，现在是docker的3306端口映射到本地的13306端口，如果端口冲突的话可以将13306修改为其他值，再将DSN中的端口一起修改
```shell
cd simpleTiktok
docker-compose up -d
```

---
**minio改为部署到服务器上了，服务器地址在constants里**
## docker部署minio

minio的后台网页：localhost:9001，API端口：localhost:9000

后台网页可以查看bucket情况，启动minio之后，进入后台网页，点击Access Keys，Create access key，然后把simpleTiktok/biz/dal/minio/init.go的Init()函数的accessKeyID和secretAccessKey改为自己新建的

```
# 使用即可同时启动mysql与minio
docker-compose up -d 
```
---

## 上传视频
视频和封面保存格式如下

| play_url    | cover_url       |
| ----------- | --------------- |
| video/1.mp4 | snapshot/1.jpeg |


## 启动
```shell
go mod tidy
go build 
./simpleTiktok
```

**gorm文档：** https://gorm.cn/zh_CN/docs/index.html

**Hertz文档：** https://www.cloudwego.io/zh/docs/hertz/overview/

**参考项目：** https://github.com/cloudwego/hertz-examples/tree/main/bizdemo/tiktok_demo
