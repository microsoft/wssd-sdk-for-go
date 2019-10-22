
$ScriptPath = Split-Path $MyInvocation.MyCommand.Path
import-module "$PSScriptRoot\wssdglobal.psm1" -Force -Verbose:$false -DisableNameChecking


function VirtualMachineScaleSetCreate($yamlFile) {
		Execute-WssdCommand -Arguments  "compute vmss create --config $yamlFile"
}

function VirtualMachineScaleSetDelete($name) {
		Execute-WssdCommand -Arguments  "compute vmss delete --name $name"
}

function VirtualMachineScaleSetShow($name) {
		Execute-WssdCommand -Arguments  "compute vmss show --name $name"
}

function VirtualMachineScaleSetList() {
		Execute-WssdCommand -Arguments  "compute vmss list"
}

function VirtualMachineScaleSetUpdate($name, $yamlFile) {
		Execute-WssdCommand -Arguments  "compute vmss update --name $name --config $yamlFile"
}

Export-ModuleMember VirtualMachineScaleSetCreate
Export-ModuleMember VirtualMachineScaleSetDelete
Export-ModuleMember VirtualMachineScaleSetShow
Export-ModuleMember VirtualMachineScaleSetList
Export-ModuleMember VirtualMachineScaleSetUpdate