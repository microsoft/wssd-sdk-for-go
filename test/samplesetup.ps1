
[Environment]::SetEnvironmentVariable("WSSD_DEBUG_MODE", "on", [EnvironmentVariableTarget]::Machine)

$pname = "wssdagent"
$process = Get-Process $pname -ErrorAction SilentlyContinue
if ($process) { $process | Stop-Process }

$pname = "wssdcloudagent"
$process = Get-Process $pname -ErrorAction SilentlyContinue
if ($process) { $process | Stop-Process }

scp madhanm@NetAppLinux:/home/madhanm/repo/gopath/src/github.com/microsoft/wssdagent/bin/wssdagent.exe .
scp madhanm@NetAppLinux:/home/madhanm/repo/gopath/src/github.com/microsoft/wssdcloudagent/bin/wssdcloudagent.exe . 
scp madhanm@NetAppLinux:/home/madhanm/repo/gopath/src/github.com/microsoft/wssd-sdk-for-go/bin/wssdctl.exe c:/windows/system32
scp madhanm@NetAppLinux:/home/madhanm/repo/gopath/src/github.com/microsoft/wssdcloud-sdk-for-go/bin/wssdcloudctl.exe c:/windows/system32
#scp madhanm@NetAppLinux:~/repo/gopath/src/github.com/microsoft/wssdagent/certs/wssd*.pem .

start ./wssdagent.exe
start ./wssdcloudagent.exe


scp madhanm@NetAppLinux:~/repo/gopath/src/github.com/microsoft/wssdcloud-sdk-for-go/test/pester/* cloud/test
scp madhanm@NetAppLinux:~/repo/gopath/src/github.com/microsoft/wssd-sdk-for-go/test/pester/* node/test
scp madhanm@NetAppLinux:~/repo/gopath/src/github.com/microsoft/wssdagent/pkg/psdriver/scripts/vhd.ps1 .