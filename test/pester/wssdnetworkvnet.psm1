

$ScriptPath = Split-Path $MyInvocation.MyCommand.Path
import-module "$PSScriptRoot\wssdglobal.psm1" -Force -Verbose:$false -DisableNameChecking
	$Global:sampleVirtualNetwork = "Default Switch"


function VirtualNetworkCreate($yamlFile) {
		Execute-WssdCommand -Arguments  "network vnet create --config $yamlFile"
}

function VirtualNetworkDelete($name) {
		Execute-WssdCommand -Arguments  "network vnet delete --name `"$name`""
}

function VirtualNetworkShow($name) {
		Execute-WssdCommand -Arguments  "network vnet show --name `"$name`""
}

function VirtualNetworkList() {
		Execute-WssdCommand -Arguments  "network vnet list"
}

function VirtualNetworkUpdate($name, $yamlFile) {
		Execute-WssdCommand -Arguments  "network vnet update --name `"$name`" --config $yamlFile"
}

function CreateSampleVirtualNetwork() {
	$yaml = @"
name: $Global:sampleVirtualNetwork
type: "ICS"
"@
		$yamlFile = "testVirtualNetwork.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		VirtualNetworkCreate $yamlFile
}

function DeleteSampleVirtualNetwork() {
	VirtualNetworkDelete $Global:sampleVirtualNetwork
}

Export-ModuleMember VirtualNetworkCreate
Export-ModuleMember VirtualNetworkDelete
Export-ModuleMember VirtualNetworkShow
Export-ModuleMember VirtualNetworkList
Export-ModuleMember VirtualNetworkUpdate
Export-ModuleMember CreateSampleVirtualNetwork
Export-ModuleMember DeleteSampleVirtualNetwork