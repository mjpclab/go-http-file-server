VERSION=$(git describe --tags --candidates=999 2> /dev/null || git rev-parse --short HEAD 2> /dev/null)

# remove prefixing "v"
VERSION=${VERSION#v}

# remove branch identifier like "-go1.9"
VERSION=$(sed -e 's/-go[0-9]*\.[0-9]*//' <<< "$VERSION")

# remove additional commit count "X.Y.Z-n-gCOMMIT" -> "X.Y.Z-gCOMMIT"
VERSION=${VERSION/-*-/-}
