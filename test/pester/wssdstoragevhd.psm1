

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
	$Global:sampleVirtualHardDiskSource = "$pwd/test.vhdx"
$yaml = @"
name: $Global:sampleVirtualHardDisk
virtualharddiskproperties:
  source: $Global:sampleVirtualHardDiskSource	
"@
		$yamlFile = "testVirtualHardDisk.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		VirtualHardDiskCreate $yamlFile
}

function DeleteSampleVirtualHardDisk() {
	VirtualHardDiskDelete $Global:sampleVirtualHardDisk
	remove-item $Global:sampleVirtualHardDiskSource
}

function CreateVMMSVhd() {
	$pwd = (pwd).Path
	$Global:testVirtualHardDiskSource = "$pwd/test.vhdx"
}

function CleanupVMMSVhd() {
		del $Global:testVirtualHardDiskSource
}

Export-ModuleMember VirtualHardDiskCreate
Export-ModuleMember VirtualHardDiskDelete
Export-ModuleMember VirtualHardDiskShow
Export-ModuleMember VirtualHardDiskList
Export-ModuleMember VirtualHardDiskUpdate
Export-ModuleMember CreateSampleVirtualHardDisk
Export-ModuleMember DeleteSampleVirtualHardDisk
Export-ModuleMember CreateVMMSVhd
Export-ModuleMember CleanupVMMSVhd
