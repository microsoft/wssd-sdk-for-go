$ScriptPath = Split-Path $MyInvocation.MyCommand.Path
$Global:wssdcloutctl = (Get-Command 'wssdctl.exe').Source
$Global:debugMode = $false

function Execute-Command(
    $Command,
    $Arguments
) {

    $result = (& $Command $Arguments.Split(" ") 2>&1)

    $out = $result | ?{$_.gettype().Name -ne "ErrorRecord"}  # On a non-zero exit code, this may contain the error
    $outString = ($out | Out-String).ToLowerInvariant()

    if ($LASTEXITCODE) {
       $err = $result | ?{$_.gettype().Name -eq "ErrorRecord"}
       throw "$Command $Arguments failed to execute [$err]"
    }
    return $out
}

function Execute-WssdCommand(
    $Arguments
) {
    $newArgs = $Arguments
    if ($Global:debugMode) {
        $newArgs = "$Arguments --debug"
    }
    Execute-Command -Command $Global:wssdcloutctl -Arguments $newArgs

}
