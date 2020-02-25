

$ScriptPath = Split-Path $MyInvocation.MyCommand.Path
import-module "$PSScriptRoot\wssdglobal.psm1" -Force -Verbose:$false -DisableNameChecking


	$Global:sampleNetworkInterface = "sampleNetworkInterface"

function NetworkInterfaceCreate($yamlFile) {
		Execute-WssdCommand -Arguments  "network vnic create --config $yamlFile"
}

function NetworkInterfaceDelete($name) {
		Execute-WssdCommand -Arguments  "network vnic delete --name $name"
}

function NetworkInterfaceShow($name) {
		Execute-WssdCommand -Arguments  "network vnic show --name $name"
}

function NetworkInterfaceList() {
		Execute-WssdCommand -Arguments  "network vnic list"
}

function NetworkInterfaceUpdate($name, $yamlFile) {
		Execute-WssdCommand -Arguments  "network vnic update --name $name --config $yamlFile"
}

function CreateSampleNetworkInterface() {
		$yaml = @"
name: $Global:sampleNetworkInterface
virtualnetworkinterfaceproperties:
  ipconfigurations:
  - ipconfigurationproperties:
      subnetid: $Global:sampleVirtualNetwork
"@
		$yamlFile = "testNetworkInterface.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		NetworkInterfaceCreate $yamlFile
}

function DeleteSampleNetworkInterface() {
	NetworkInterfaceDelete $Global:sampleNetworkInterface
}


Export-ModuleMember NetworkInterfaceCreate
Export-ModuleMember NetworkInterfaceDelete
Export-ModuleMember NetworkInterfaceShow
Export-ModuleMember NetworkInterfaceList
Export-ModuleMember NetworkInterfaceUpdate
Export-ModuleMember CreateSampleNetworkInterface
Export-ModuleMember DeleteSampleNetworkInterface
