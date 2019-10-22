$Global:wssdcloutctl = (Get-Command "wssdctl.exe").Source

function Execute-Command(
    $Command,
    $Arguments
) {

    & $command $Arguments.Split(" ")
}

function Execute-WssdCommand(
    $Arguments
) {

    Execute-Command -Command $Global:wssdcloutctl -Arguments $Arguments
}
