#!/bin/bash
echo "Start Run Go Vet..."
cd $PROJECT_DIR

go tool vet -all -unreachable=false   $(find . -type f -name '*.go' -not -path "./vendor/*"  -not -path "./routers/*" -not -path "./tests/*")
retVal=$?
exit $retVal
#!/bin/bash
echo "Start Run Go Vet..."
cd $PROJECT_DIR

go tool vet -all -unreachable=false   $(find . -type f -name '*.go' -not -path "./vendor/*"  -not -path "./routers/*" -not -path "./tests/*")
retVal=$?
exit $retVal
