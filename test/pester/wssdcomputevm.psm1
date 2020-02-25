

$ScriptPath = Split-Path $MyInvocation.MyCommand.Path
import-module "$PSScriptRoot\wssdglobal.psm1" -Force -Verbose:$false -DisableNameChecking

$Global:sampleVirtualMachine = "sampleVirtualMachine"

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

function CreateSampleVirtualMachine($virtualHardDisk, $networkInterface) {
	$yaml = @"
name: $Global:sampleVirtualMachine
virtualmachineproperties:
  storageprofile:
    osdisk:
      name: null
      ostype: "Linux"
      vhdname: $virtualHardDisk
    datadisks: []
  osprofile:
    computername: "lumaster"
    adminusername: "localadmin"
    adminpassword: ""
    customdata: ""
    windowsconfiguration: null
    linuxconfiguration:
      ssh:
        publickeys:
        - keydata: ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDKxmVSZOphCI2RWMJf5qvtMwmLBo0OlG1knLt4Yk26JqOTtGWdqmJM7QcQevBp6wnBKhzEIheq/kUJ8lRMoGplZ4wPsTu/BO2IgoAi0/xIX9NalRCD1TpLPOmaa7nqGi/7+BbTznbqtDDqKST80juLT+bbz5g3UIxsSu+R2Rpm782AzDkQ61K3YFuRiK4c58+ANZv790NTltQ3Y9iO0ivJ1dbiNXj1qVEEuXAuP80d4MgQHt+rwNdpex/2p5NHRpC/GYuSwrjQBgBX2hgOT2kvAq19x55D0bcvZ99+M9Ar9TBCfVfme7GGFceD1qrhJdXQapqhO9FJG9qk75Iti2BX
      disablepasswordauthentication: true
  networkprofile:
    networkinterfaces:
    - virtualnetworkinterfacereference: $networkInterface
"@
	$yamlFile = "testVirtualMachine.yaml"
	Set-Content -Path $yamlFile -Value $yaml 

	VirtualMachineCreate $yamlFile  # | Should Not Throw
}

function DeleteSampleVirtualMachine() {
	VirtualMachineDelete -name $Global:sampleVirtualMachine
}
Export-ModuleMember VirtualMachineCreate
Export-ModuleMember VirtualMachineDelete
Export-ModuleMember VirtualMachineShow
Export-ModuleMember VirtualMachineList
Export-ModuleMember VirtualMachineUpdate
Export-ModuleMember CreateSampleVirtualMachine
Export-ModuleMember DeleteSampleVirtualMachine