# 为指定路径显示页面
```
GET <path>[?sort=sortBy]
```
无论路径是否以“/”结尾，都应正常工作。

可用的排序key：
- `n` 按名称递增排序
- `N` 按名称递减排序
- `e` 按类型（后缀）递增排序
- `E` 按类型（后缀）递减排序
- `s` 按大小递增排序
- `S` 按大小递减排序
- `t` 按修改时间递增排序
- `T` 按修改时间递减排序
- `_` 不排序

目录顺序：
- `/<key>` 目录在文件之前
- `<key>/` 目录在文件之后
- `<key>` 目录与文件混合

举例：
```sh
curl 'http://localhost/ghfs/'
curl 'http://localhost/ghfs/?sort=/T'
```

# 获取指定路径JSON形式的数据
```
GET <path>[?sort=key]
Accept: application/json
```

举例：
```sh
curl -H 'Accept: application/json' 'http://localhost/ghfs/'
```

# 显示精简页面
```
GET <path>?simple[&sort=key]
```
类似于常规显示的页面，但隐藏路径列表、可排序的表头和上级目录链接。
这为“wget”之类的工具递归下载提供了方便。

举例：
```shell
wget --recursive -nc -nH -np 'http://localhost/dir/?download'
```

# 为文件添加下载参数
```
GET <path>?download[&sort=key]
```
为页面中的文件链接添加下载参数，使其可被下载，而不是显示其内容。

# 下载单个文件
通过输出`Content-Disposition`头，通知用户代理下载文件而不是显示其内容。
```
GET <path/to/file>?download[=filename]
```

举例：
```sh
curl 'http://localhost/ghfs/file?download'
```

# 以打包文件形式获取指定路径下的内容
仅在“archive”选项启用时有效。
```
GET <path>?tar[=filename]
GET <path>?tgz[=filename]
GET <path>?zip[=filename]
POST <path>?tar[=filename]
POST <path>?tgz[=filename]
POST <path>?zip[=filename]
```

举例：
```sh
curl 'http://localhost/tmp/?zip' > tmp.zip
```

要打包当前目录下的指定子项，用`name`参数指定：
```
GET <path>?tar&name=<path1>&name=<path2>&...name=<pathN>
GET <path>?tgz&name=<path1>&name=<path2>&...name=<pathN>
GET <path>?zip&name=<path1>&name=<path2>&...name=<pathN>
```

```
POST <path>?tar

name=<path1>&name=<path2>&...name=<pathN>
```

举例：
```sh
curl -X POST -d 'name=subdir1&name=subdir2/subdir21&name=file1&name=subdir3/file31' 'http://localhost/tmp/?zip' > tmp.zip
```

# 在指定路径下创建目录
仅在“mkdir”选项启用时有效。
```
POST <path>?mkdir
[Accept: application/json]

name=<dir1path>&name=<dir2path>&...name=<dirNpath>
```

举例：
```sh
curl -X POST -d 'name=dir1&name=dir2&name=foo/bar/baz' 'http://localhost/tmp/?mkdir'
```

# 上传文件到指定路径
仅在“upload”选项启用时有效。
```
POST <path>?upload
[Accept: application/json]
```
- 必须使用`POST`方法
- 必须使用`multipart/form-data`编码
- 每个文件内容占用一个段，表单字段名可以是`file`，`dirfile`或`innerdirfile`

举例：
```sh
curl -F 'file=@file1.txt' -F 'file=@file2.txt;filename=renamed.txt' 'http://localhost/tmp/?upload'
```

如果还启用了“mkdir”选项，可以将文件上传到相对于当前URL路径的特定路径，
使用表单字段`dirfile`代替`file`：
```sh
curl -F 'dirfile=@file1.txt;filename=subdir/childdir/filename.txt' 'http://localhost/tmp/?upload'
# 文件现在位于 http://localhost/tmp/subdir/childdir/filename.txt
```

另一表单字段`innerdirfile`与`dirfile`很相似，只是会去除第一级上传目录。
这对于上传一个目录中的内容很方便：
```sh
curl -F 'innerdirfile=@file1.txt;filename=subdir/childdir/filename.txt' 'http://localhost/tmp/?upload'
# 文件现在位于 http://localhost/tmp/childdir/filename.txt
```

# 在指定路径下删除文件或目录
仅在“delete”选项启用时有效。
目录将被递归删除。
```
POST <path>?delete
[Accept: application/json]

name=<dir1>&name=<dir2>&...name=<dirN>
```

举例：
```sh
curl -X POST -d 'name=dir1&name=dir2&name=dir3' 'http://localhost/tmp/?delete'
```

# 登录
发起登录认证，即使当前路径无需验证：
```
GET <path>?auth[=return_url]
```
