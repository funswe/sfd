# Simple File Download
**sfd**是一个GO语言开发的，简单易用的下载网络文件（图片，HTML，视频，音频）小工具

## 功能
- **单文件下载**： 可以直接下载一个网络文件
- **多文件下载**： 可以将网络文件放在一个文件里，同时下载多个文件

## 安装
```bash
# go install github.com/funswe/sfd@v1.0.0
```
执行完之后会在GOPATH路径下的bin目录里看到sfd可执行程序

## 使用方法

```bash
# sfd -h
Usage of sfd:
  -f string
        Individual remote files to download
  -l string
        List of remote files to download
  -n int
        Number of parallel downloads (default cpu nums)
  -o string
        Output path for remote file download（default current path）
```

## 下载单个文件
```bash
# sfd -f https://cos-anonymous-cdn.bw-yx.com/test.pdf -o ./pdf
total download files: 1
test.pdf [=============>----------------------------------------------------------------]  19 %
```

## 下载多个文件
```bash
# cat file-list.txt
https://cos-anonymous-cdn.bw-yx.com/test.pdf
https://cos-anonymous-cdn.bw-yx.com/videos/painting/4-36f152a0113d658280c28b3e634d3c82/ab1bd410bef711ec9bf71b9cd156e33a.mp4
https://cos-anonymous-cdn.bw-yx.com/videos/painting/4-36f152a0113d658280c28b3e634d3c82/25f25c90bef811ecb66d0b0a269545fe.mp4
```
```bash
# sfd -l file-list.txt -o ./pdf
total download files: 3
test.pdf [==================>-----------------------------------------------------------]  24 %
25f25c90bef811ecb66d0b0a269545fe.mp4 [=========================================================>--------------------]  74 %
ab1bd410bef711ec9bf71b9cd156e33a.mp4 [>-----------------------------------------------------------------------------]   2 %
```
