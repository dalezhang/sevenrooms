#!/bin/bash
echo "Start Run Go Fmt..."
cd $PROJECT_DIR

unformatted=$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"  -not -path "./routers/*" -not -path "./tests/*"))
[ -z "$unformatted" ] && exit 0

echo >&2 "\033[0;31m Go files must be formatted with gofmt. Please run:"
  for fn in $unformatted; do
      echo >&2 "\033[0;31m  gofmt -w $PWD/$fn"
  done

exit 1