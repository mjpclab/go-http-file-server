VERSION=$(git describe --abbrev=0 --tags 2> /dev/null || git rev-parse --abbrev-ref HEAD 2> /dev/null)
VERSION=${VERSION#v}
VERSION=${VERSION%-go*}
