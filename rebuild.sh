./kall.bat
./resetdata.sh
go build cmd/git-server/git-server.go
rm assets/bin/git/git-server.exe
mv git-server.exe assets/bin/git/