
# Baofeng Cloud 命令行工具

设置并保存AK/SK，会保存在$HOME/.bfcloud文件里
```shell
$bfcloud config
```


上传视频（默认Paas,public)

```shell
$bfcloud upload test.mp4 c:\test.mp4
```


上传视频（Saas, private)

```shell
$bfcloud -service 1 -private upload test.mp4 c:\test.mp4
```


查询

```shell
$bfcloud query test.mp4
```


删除

```shell
$bfcloud delete test.mp4
```

