

$ScriptPath = Split-Path $MyInvocation.MyCommand.Path
import-module "$PSScriptRoot\wssdglobal.psm1" -Force -Verbose:$false -DisableNameChecking


function VirtualHardDiskCreate($yamlFile) {
		Execute-WssdCommand -Arguments  "storage vhd create --config $yamlFile"
}

function VirtualHardDiskDelete($name) {
		Execute-WssdCommand -Arguments  "storage vhd delete --name $name"
}

function VirtualHardDiskShow($name) {
		Execute-WssdCommand -Arguments  "storage vhd show --name $name"
}

function VirtualHardDiskList() {
		Execute-WssdCommand -Arguments  "storage vhd list"
}

function VirtualHardDiskUpdate($name, $yamlFile) {
		Execute-WssdCommand -Arguments  "storage vhd update --name $name --config $yamlFile"
}

function CreateSampleVirtualHardDisk() {
	$Global:sampleVirtualHardDisk = "sampleVirtualHardDisk"
	$pwd = (pwd).Path
	$Global:sampleVirtualHardDiskSource = "$Script:ScriptPath\test.vhdx"
$yaml = @"
name: $Global:sampleVirtualHardDisk
virtualharddiskproperties:
  source: $Global:sampleVirtualHardDiskSource	
"@
		$yamlFile = "testVirtualHardDisk.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		VirtualHardDiskCreate $yamlFile
}

function CreateSampleVirtualHardDiskDataDisk() {
	$script:sampleVirtualHardDisk = "sampleVirtualHardDisk"

	It 'Should be able to create a virtual hard disk of type data disk' {
		$yaml = @"
name: $Global:sampleVirtualHardDisk
virtualharddiskproperties:
  source: ""
  path: "c:\\cluster\\volume1\\testdatadisk.vhdx"
  disksizegb: 10737418240
  dynamic: true
  blocksizebytes: 33554432
  logicalsectorbytes: 4096
  physicalsectorbytes: 4096
  controllernumber: 0
  controllerlocation: 0
  disknumber: 0
  vmname: ""
  vmid: ""
  scsipath: "0.0.0.0"
  virtualharddisktype: DATADISK_VIRTUALHARDDISK	
"@
		$yamlFile = "testVirtualHardDiskDataDisk.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		VirtualHardDiskCreate $yamlFile
	}
}

function DeleteSampleVirtualHardDisk() {
	VirtualHardDiskDelete $Global:sampleVirtualHardDisk
}

function CreateVMMSVhd() {
	$Global:testVirtualHardDiskSource = "$Script:ScriptPath\test.vhdx"
}

function CleanupVMMSVhd() {
}

Export-ModuleMember VirtualHardDiskCreate
Export-ModuleMember VirtualHardDiskDelete
Export-ModuleMember VirtualHardDiskShow
Export-ModuleMember VirtualHardDiskList
Export-ModuleMember VirtualHardDiskUpdate
Export-ModuleMember CreateSampleVirtualHardDisk
Export-ModuleMember CreateSampleVirtualHardDiskDataDisk
Export-ModuleMember DeleteSampleVirtualHardDisk
Export-ModuleMember CreateVMMSVhd
Export-ModuleMember CleanupVMMSVhd
