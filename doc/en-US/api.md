# Render page of specified path
```
GET <path>[?sort=<sortBy>]
```
Should work no matter tailing “/” is present or not in path.

Available sort key:
- `n` sort by name ascending
- `N` sort by name descending
- `e` sort by type(suffix) ascending
- `E` sort by type(suffix) descending
- `s` sort by size ascending
- `S` sort by size descending
- `t` sort by modify time ascending
- `T` sort by modify time descending
- `_` no sort

Directory sort:
- `/<key>` directories before files
- `<key>/` directories after files
- `<key>` directories mixed with files

Example:
```sh
curl 'http://localhost/ghfs/'
curl 'http://localhost/ghfs/?sort=/T'
```

# Get JSON data of specified path
```
GET <path>[?sort=key]
Accept: application/json
```

Example:
```sh
curl -H 'Accept: application/json' 'http://localhost/ghfs/'
```

# Render page in simple mode
```
GET <path>?simple[&sort=key]
```
Similar to regular page rendering, but
hide path list, sortable list header, and parent directory link.
It's convenient for tools like "wget" to download files recursively.

Example:
```shell
wget --recursive -nc -nH -np 'http://localhost/dir/?download'
```

# Add download param to files
```
GET <path>?download[&sort=key]
```
Add download parameter to file links in directory list page, which
makes them downloadable instead of displaying content.

# Download a file
Notify user agent download a file rather than displaying its content,
by outputting `Content-Disposition` header.
```
GET <path/to/file>?download[=filename]
```

Example:
```sh
curl 'http://localhost/ghfs/file?download'
```

# Get contents of specified path as archive file
Only work when "archive" is enabled.
```
GET <path>?tar[=filename]
GET <path>?tgz[=filename]
GET <path>?zip[=filename]
POST <path>?tar[=filename]
POST <path>?tgz[=filename]
POST <path>?zip[=filename]
```

Example:
```sh
curl 'http://localhost/tmp/?zip' > tmp.zip
```

To archive specific sub items under current directory, pass `name` params:
```
GET <path>?tar&name=<path1>&name=<path2>&...name=<pathN>
GET <path>?tgz&name=<path1>&name=<path2>&...name=<pathN>
GET <path>?zip&name=<path1>&name=<path2>&...name=<pathN>
```

```
POST <path>?tar

name=<path1>&name=<path2>&...name=<pathN>
```

Example:
```sh
curl -X POST -d 'name=subdir1&name=subdir2/subdir21&name=file1&name=subdir3/file31' 'http://localhost/tmp/?zip' > tmp.zip
```

# Create directories in specific path
Only work when "mkdir" is enabled.
```
POST <path>?mkdir
[Accept: application/json]

name=<dir1path>&name=<dir2path>&...name=<dirNpath>
```

Example:
```sh
curl -X POST -d 'name=dir1&name=dir2&name=foo/bar/baz' 'http://localhost/tmp/?mkdir'
```

# Upload files to specific path
Only work when "upload" is enabled.
```
POST <path>?upload
[Accept: application/json]
```
- Must use `POST` method
- Must use `multipart/form-data` encoding type
- Each file content use one part, form field name can be `file`, `dirfile` or `innerdirfile`

Example:
```sh
curl -F 'file=@file1.txt' -F 'file=@file2.txt;filename=renamed.txt' 'http://localhost/tmp/?upload'
```

If "mkdir" is also enabled, it is possible to upload file to a specific path relative to current URL path,
using form name `dirfile` instead of `file`:
```sh
curl -F 'dirfile=@file1.txt;filename=subdir/childdir/filename.txt' 'http://localhost/tmp/?upload'
# file is now available at http://localhost/tmp/subdir/childdir/filename.txt
```

Another form name `innerdirfile` is similar to `dirfile`, but strip first level of upload directory.
It is convenient to upload contents of a directory:
```sh
curl -F 'innerdirfile=@file1.txt;filename=subdir/childdir/filename.txt' 'http://localhost/tmp/?upload'
# file is now available at http://localhost/tmp/childdir/filename.txt
```

# Delete files or directories in specific path
Only work when "delete" is enabled.
Directories will be deleted recursively.
```
POST <path>?delete
[Accept: application/json]

name=<dir1>&name=<dir2>&...name=<dirN>
```

Example:
```sh
curl -X POST -d 'name=dir1&name=dir2&name=dir3' 'http://localhost/tmp/?delete'
```

# Login
Perform a login authentication even not required by current path:
```
GET <path>?auth[=return_url]
```
