# Simple File Download
**sfd**是一个GO语言开发的，简单易用的下载网络文件（图片，HTML，视频，音频）小工具

## 功能
- **单文件下载**： 可以直接下载一个网络文件
- **多文件下载**： 可以将网络文件放在一个文件里，同时下载多个文件

## 使用方法

```bash
# sfd -h
Usage of sfd:
  -f string
        Individual remote files to download
  -l string
        List of remote files to download
  -n int
        Number of parallel downloads (default 8)
  -o string
        Output path for remote file download（default current path）
```