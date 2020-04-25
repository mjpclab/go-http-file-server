# 为指定路径显示页面
```
GET <path>[?sort=sortBy]
```
无论路径是否以“/”结尾，都应正常工作。

可用的排序key：
- `n` 按名称递增排序
- `N` 按名称递减排序
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
curl http://localhost/ghfs/
curl http://localhost/ghfs/?sort=/T
```

# 获取指定路径JSON形式的数据
```
GET <path>?json[&sort=key]
```

举例：
```sh
curl http://localhost/ghfs/?json
```

# 以打包文件形式获取指定路径下的内容
仅在“archive”选项启用时有效。
```
GET <path>?tar
GET <path>?tgz
GET <path>?zip
```

举例：
```sh
curl http://localhost/tmp/?zip > tmp.zip
```

# 上传文件到指定路径
仅在“upload”选项启用时有效。
```
POST <path>?upload[&json]
```
- 必须使用`POST`方法
- 必须使用`multipart/form-data`编码
- 每个文件内容占用一个段，字段名为`file`

举例：
```sh
curl -F 'file=@file1.txt' -F 'file=@file2.txt' http://localhost/tmp/?upload
```

# 在指定路径下创建目录
仅在“mkdir”选项启用时有效。
```
GET <path>?mkdir[&json]&name=<dir1>&name=<dir2>&...name=<dirN>
```
```
POST <path>?mkdir[&json]

name=<dir1>&name=<dir2>&...name=<dirN>
```

举例：
```sh
curl -X POST -d 'name=dir1&name=dir2&name=dir3' http://localhost/tmp/?mkdir
```

# 在指定路径下删除文件或目录
仅在“delete”选项启用时有效。
目录将被递归删除。
```
GET <path>?delete[&json]&name=<dir1>&name=<dir2>&...name=<dirN>
```
```
POST <path>?delete[&json]

name=<dir1>&name=<dir2>&...name=<dirN>
```

举例：
```sh
curl -X POST -d 'name=dir1&name=dir2&name=dir3' http://localhost/tmp/?delete
```
