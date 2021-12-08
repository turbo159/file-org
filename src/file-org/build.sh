echo "Started compile..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ../../bin/mac
cp tasks.json ../../bin/mac
echo "Mac arm64 build complete."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../bin/linux
cp tasks.json ../../bin/linux
echo "Linux amd64 build complete."
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o ../../bin/rasberry
cp tasks.json ../../bin/rasberry
echo "Rasberry Pi linux arm32 build complete."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ../..bin/windows
cp tasks.json ../../bin/windows
echo "Windows amd64 build complete."
echo "Done."
echo "."
ls -lhFR $PWD/../../bin|grep file-org