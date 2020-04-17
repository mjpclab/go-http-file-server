# Render page of specified path
```
GET <path>
```
Should work no matter tailing “/” is present or not in path.

Example:
```sh
curl http://localhost/ghfs/
```

# Get JSON data of specified path
```
GET <path>?json
```

Example:
```sh
curl http://localhost/ghfs/?json
```

# Get contents of specified path as archive file
Only work when "archive" is enabled.
```
GET <path>?tar
GET <path>?tgz
GET <path>?zip
```

Example:
```sh
curl http://localhost/tmp/?zip > tmp.zip
```

# Upload files to specific path
Only work when "upload" is enabled.
```
POST <path>?upload[&json]
```
- Must use `POST` method
- Must use `multipart/form-data` encoding type
- Each file content use one part, field name is `file`

Example:
```sh
curl -F 'file=@file1.txt' -F 'file=@file2.txt' http://localhost/tmp/?upload
```

# Create directories in specific path
Only work when "mkdir" is enabled.
```
GET <path>?mkdir[&json]&name=<dir1>&name=<dir2>&...name=<dirN>
```
```
POST <path>?mkdir[&json]

name=<dir1>&name=<dir2>&...name=<dirN>
```

Example:
```sh
curl -X POST -d 'name=dir1&name=dir2&name=dir3' http://localhost/tmp/?mkdir
```

# Delete files or directories in specific path
Only work when "delete" is enabled.
Directories will be deleted recursively.
```
GET <path>?delete[&json]&name=<dir1>&name=<dir2>&...name=<dirN>
```
```
POST <path>?delete[&json]

name=<dir1>&name=<dir2>&...name=<dirN>
```

Example:
```sh
curl -X POST -d 'name=dir1&name=dir2&name=dir3' http://localhost/tmp/?delete
```
