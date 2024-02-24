# Go HTTP File Server
基于命令行的HTTP文件共享服务器。

![Go HTTP File Server pages](doc/ghfs.gif)

## 软件特色
- 比Apache/Nginx更友好的目录列表
- 适配移动设备显示
- 单一可执行文件
- 可以将当前浏览目录内容打包下载
- 可以开启某个目录的上传权限
- 可以指定自定义主题来渲染页面
- 支持目录别名（将另一个目录挂载到某个URL路径）

## 编译
至少需要Go 1.18版本。
```sh
go build main.go
```
会在当前目录生成"main"可执行文件。

## 举例
在8080端口启动服务器，根目录为当前工作目录：
```sh
ghfs -l 8080
``` 

在8080端口启动服务器，根目录为 /usr/share/doc：
```sh
ghfs -l 8080 -r /usr/share/doc
```

在默认端口启动服务器，根目录为/tmp，并允许上传文件到/tmp/data：
```sh
ghfs -r /tmp -u /data
```

共享/etc下的文件，同时把/usr/share/doc挂载到URL路径/doc下：
```sh
ghfs -r /etc -a :/doc:/usr/share/doc
```

在8080端口启动服务器，使用HTTPS协议：
```sh
ghfs -k /path/to/certificate/key -c /path/to/certificate/file -l 8080
```

不显示`.`开头的unix隐藏目录和文件。提示：用引号括起通配符以避免shell展开：
```sh
ghfs -H '.*'
```

在命令行显示访问日志：
```sh
ghfs -L -
```

http基本验证：
- 对URL /files 启用验证
- 用户名：user1，密码：pass1
- 用户名：user2，密码：pass2
```sh
ghfs --auth /files --user user1:pass1 --user-sha1 user2:8be52126a6fde450a7162a3651d589bb51e9579d
```

启动2台虚拟主机：
- 服务器1
    - 在80端口提供http服务
    - 在443端口提供https服务
        - 证书文件：/cert/server1.pem
        - 私钥文件：/cert/server1.key
    - 主机名：server1.example.com
    - 根目录：/var/www/server1
- 服务器2
    - 在80端口提供http服务
    - 在443端口提供https服务
        - 证书文件：/cert/server2.pem
        - 私钥文件：/cert/server2.key
    - 主机名：server2.example.com
    - 根目录：/var/www/server2
```sh
ghfs --listen-plain 80 --listen-tls 443 -c /cert/server1.pem -k /cert/server1.key --hostname server1.example.com -r /var/www/server1 ,, --listen-plain 80 --listen-tls 443 -c /cert/server2.pem -k /cert/server2.key --hostname server2.example.com -r /var/www/server2
```

## 使用方法
```
ghfs [选项]

-l|--listen <IP|端口|:端口|IP:端口|socket> ...
    指定服务器要侦听的IP和端口，例如“:80”或“127.0.0.1:80”。
    如果指定了--cert和--key，端口接受TLS连接。
    如果未指定端口，则在纯HTTP模式下使用80端口，TLS模式下使用443端口。
    如果值中包含“/”，则将其当作unix socket路径。
    标志“-l”或“--listen”可以省略。
--listen-plain <IP|端口|:端口|IP:端口|socket> ...
    与--listen类似，但强制使用非TLS模式。
--listen-tls <IP|端口|:端口|IP:端口|socket> ...
    与--listen类似，但强制使用TLS模式。若未指定证书和私钥，则启动失败。

--hostname <主机名> ...
    指定与当前虚拟主机关联的主机名。
    如果值以“.”开头，则将其当作后缀，匹配该域下的所有子域名，例如“.example.com”。
    如果值以“.”结尾，则将其当作前缀，匹配所有域名后缀。

-r|--root <目录>
    服务器的根目录。
    默认为当前目录。

-R|--empty-root
    使用空的虚拟目录作为根目录。
    在仅需挂载别名的情况下较实用。

-a|--alias <分隔符><URL路径><分隔符><文件系统路径> ...
    设置路径别名。
    将某个文件系统路径挂载到URL路径下。
    例如：“:/doc:/usr/share/doc”。

--prefix <path> ...
    在指定的URL子路径下提供服务。
    如果服务器在反向代理之后，且收到的请求并未去除代理路径前缀，可能较有用。

-/|--auto-dir-slash [<状态码>=301]
    如果在请求目录列表页时URL没有以“/”结尾，重定向到带有该结尾的URL。
    如果在请求文件时URL以“/”结尾，重定向到不带有该结尾的URL。

--default-sort <排序规则>
    指定文件和目录的默认排序规则。
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

-I|--dir-index <文件> ...
    指定目录默认页面文件。

--global-restrict-access [<允许的主机> ...]
    限制第三方主机对所有URL路径的访问，它是通过检测请求头中的`Referer`或`Origin`实现的。
    如果该请求头为空，仍然能够访问目录列表。
    如果未指定允许的第三方主机，文件内容仅可被当前主机访问。注意这样无法限制把域名指向你的
    主机且能匹配当前虚拟主机的人，除非明确指定允许的主机。
    “主机”可以是主机名，即使用默认端口，也可以是“主机:端口”的形式。
--restrict-access <分隔符><URL路径>[<分隔符><允许的主机>...]
    与--global-restrict-access类似，但仅限于指定的URL路径（及子路径）。
    例如"#/url/path#example1.com#example2.com"。
--restrict-access-dir <分隔符><文件系统路径>[<分隔符><允许的主机>...]
    与--global-restrict-access类似，但仅限于指定的文件系统路径（及子路径）。
    例如"#/fs/path#example1.com#example2.com"。

--global-header <名称>:<值> ...
    添加自定义HTTP响应头。
--header <分隔符><URL路径><分隔符><名称><分隔符><值> ...
    为指定的URL路径（及子路径）添加自定义HTTP响应头。
--header-dir <分隔符><文件系统路径><分隔符><名称><分隔符><值> ...
    与--header类似，但指定的是文件系统路径，而不是URL路径。

-U|--global-upload
    对所有URL路径开启上传权限。
    请谨慎使用。
-u|--upload <URL路径> ...
    设置允许上传的URL路径（及子路径）。
    请谨慎使用。
--upload-dir <文件系统路径> ...
    与--upload类似，但指定的是文件系统路径，而不是URL路径。

    上传选项注意事项：
        如果名称已存在且是常规文件，
        若已启用删除（例如--delete选项），则尝试先删除文件，
        否则尝试添加或递增数字后缀。
        目录上传模式中，仅当启用创建目录时，才会上传子目录。

--global-mkdir
    对所有URL路径开启创建子目录权限。
--mkdir <URL路径> ...
    设置允许创建子目录的URL路径（及子路径）。
--mkdir-dir <文件系统路径> ...
    与--mkdir类似，但指定的是文件系统路径，而不是URL路径。

    创建子目录选项注意事项：
        为避免歧义，被别名遮蔽的目录名不能被创建。

--global-delete
    对所有URL路径开启删除子项权限。
--delete <URL路径> ...
    设置允许删除子项的URL路径（及子路径）。
--delete-dir <文件系统路径> ...
    与--delete类似，但指定的是文件系统路径，而不是URL路径。

    删除选项注意事项：
        为避免歧义，URL路径下挂载的别名不能被删除。
        别名下的非别名文件/目录仍然可以被删除。
        为避免歧义，被别名遮蔽的正常文件/目录不能被删除。

-A|--global-archive
    对所有URL路径开启打包下载当前目录内容的功能。
    页面顶部会出现下载链接。
    请确保符号链接没有循环引用。
--archive <URL路径> ...
    对指定URL路径（及子路径）开启打包下载当前目录内容的功能。
--archive-dir <文件系统路径> ...
    与--archive类似，但指定的是文件系统路径，而不是URL路径。

--global-cors
    接受所有URL路径的CORS跨域请求。
--cors <URL路径> ...
    接受指定URL路径（及子路径）的CORS跨域请求。
--cors-dir <文件系统路径> ...
    接受指定文件系统路径（及子路径）的CORS跨域请求。

--user [<用户名>]:[<密码>] ...
    为当前虚拟主机指定用于http基本验证的用户，允许空的用户名和/或密码。
--user-base64 [<用户名>]:[<base64密码>] ...
--user-md5 [<用户名>]:<md5密码> ...
--user-sha1 [<用户名>]:<sha1密码> ...
--user-sha256 [<用户名>]:<sha256密码> ...
--user-sha512 [<用户名>]:<sha512密码> ...
    指定http基本验证的用户，对密码使用特定的编码。

--global-auth
    对所有URL路径启用http基本验证(Basic Auth)。
--auth <URL路径> ...
--auth-user <分隔符><URL路径>[<分隔符><允许的用户名>...] ...
    对指定URL路径（及子路径）启用http基本验证。
--auth-dir <文件系统路径> ...
--auth-dir-user <分隔符><文件系统路径>[<分隔符><允许的用户名>...] ...
    对指定文件系统路径（及子路径）启用http基本验证。

-c|--cert <证书文件> ...
    指定TLS证书文件。

-k|--key <私钥文件> ...
    指定TLS私钥文件。

--theme <主题文件>
    指定用于渲染页面和静态资源的自定义主题zip压缩文件，代替内建主题。
    主题的内容在运行时一直缓存在内存中。
--theme-dir <主题目录>
    指定主题文件所在的目录。
    每次请求时主题内容都会重新计算。
    这为开发主题提供了便利。

    主题选项注意事项：
        --theme和--theme-dir是互斥的。
        --theme-dir更为优先。
        页面模板文件名固定为“index.html”。
        使用“?asset=<asset-path>”格式来引用主题中的静态资源。

--hsts [<有效时长>]
    启用HSTS(HTTP Strict Transport Security)。
    仅当当前虚拟主机的纯HTTP和TLS模式都监听在标准端口上时才有效。
--to-https [<目标端口>]
    将纯HTTP请求重定向到HTTPS端口。
    目标端口必须存在于当前虚拟主机--listen-tls中。
    如果省略目标端口，则使用--listen-tls中的第一项。

-S|--show <通配符> ...
-SD|--show-dir <通配符> ...
-SF|--show-file <通配符> ...
    如果指定该选项，只有匹配通配符的目录或文件（除了被hide选项隐藏的）才会显示出来。

-H|--hide <通配符> ...
-HD|--hide-dir <通配符> ...
-HF|--hide-file <通配符> ...
    如果指定该选项，匹配通配符的目录或文件不会显示出来。

-L|--access-log <文件>
    访问日志。
    使用“-”指定为标准输出。
    设为空来禁用。

-E|--error-log <文件>
    错误日志。
    使用“-”指定为标准错误输出。
    设为空来禁用。
    默认为“-”。

--config <文件>
    指定外部配置文件。

    其内容为任何其他选项，
    与在命令行指定的形式相同，
    用空白符分割。

    外部配置的优先级低于命令行选项。
    如果在命令行指定了某个选项，则其外部配置被忽略。

    使用“-”指定为标准输入。

,,
    要指定多台虚拟主机的选项，用此符号分割每台主机的选项。
    可以为每台虚拟主机分别指定以上选项。

    如果多个虚拟主机共享相同的IP和端口，
    使用--hostname指定主机名，根据请求头中的主机名来区分虚拟主机。
    如果请求的主机名不匹配任何虚拟主机，
    服务器尝试使用第一个没有指定主机名的虚拟主机，
    如果失败则使用第一个虚拟主机。
```

## 环境变量

### GHFS_PID_FILE
指定进程ID文件路径。进程ID会在应用启动时被写入文件。

### GHFS_QUIET
为避免在控制台输出额外信息，例如可访问的URL等，可将值设为“1”。

### GHFS_CPU_PROFILE_FILE
生成Go的CPU pprof profile到指定的文件路径。

## 默认主题的快捷键
- `←`, `→`：使焦点在路径项之间移动
- `Ctrl`/`Opt` + `←`：把焦点移动到第一个路径项
- `Ctrl`/`Opt` + `→`：把焦点移动到最后一个路径项
- `↑`, `↓`：使焦点在文件项之间移动
- `Ctrl`/`Opt` + `↑`：把焦点移动到第一个文件项
- `Ctrl`/`Opt` + `↓`：把焦点移动到最后一个文件项
- 重复输入相同字符将查找下一个前缀匹配该字符的文件。+ `Shift`执行反向查找。
- 输入不重复的字符，将在短时间内被记忆为字符串，用于查找下一个前缀匹配的文件。+ `Shift`执行反向查找。
- 当已启用上传时，粘贴（`Ctrl`/`Cmd` + `v`）图像或文本会将此内容上传为文件。
