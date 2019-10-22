

$ScriptPath = Split-Path $MyInvocation.MyCommand.Path
import-module "$PSScriptRoot\wssdglobal.psm1" -Force -Verbose:$false -DisableNameChecking



function VirtualMachineCreate($yamlFile) {
		Execute-WssdCommand -Arguments  "compute vm create --config $yamlFile"
}

function VirtualMachineDelete($name) {
		Execute-WssdCommand -Arguments  "compute vm delete --name $name"
}

function VirtualMachineShow($name) {
		Execute-WssdCommand -Arguments  "compute vm show --name $name"
}

function VirtualMachineList() {
		Execute-WssdCommand -Arguments  "compute vm list"
}

function VirtualMachineUpdate($name, $yamlFile) {
		Execute-WssdCommand -Arguments  "compute vm update --name $name --config $yamlFile"
}

Export-ModuleMember VirtualMachineCreate
Export-ModuleMember VirtualMachineDelete
Export-ModuleMember VirtualMachineShow
Export-ModuleMember VirtualMachineList
Export-ModuleMember VirtualMachineUpdate