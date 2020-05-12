

$ScriptPath = Split-Path $MyInvocation.MyCommand.Path
import-module "$PSScriptRoot\wssdglobal.psm1" -Force -Verbose:$false -DisableNameChecking

	$Global:testVirtualHardDiskSource = "$Script:ScriptPath\test.vhdx"
	$Global:sampleVirtualHardDisk = "sampleVirtualHardDisk"
	$Global:sampleVirtualHardDiskSource = "$Script:ScriptPath\test.vhdx"
	$Global:sampleVirtualHardDiskDataDisk = "sampleVirtualHardDiskDataDisk"

function VirtualHardDiskCreate($containerName, $yamlFile) {
		Execute-WssdCommand -Arguments  "storage vhd create --config $yamlFile --container $containerName"
}

function VirtualHardDiskDelete($name, $containerName) {
		Execute-WssdCommand -Arguments  "storage vhd delete --name $name --container $containerName"
}

function VirtualHardDiskShow($name, $containerName) {
		Execute-WssdCommand -Arguments  "storage vhd show --name $name --container $containerName"
}

function VirtualHardDiskList($containerName) {
		Execute-WssdCommand -Arguments  "storage vhd list --container $containerName"
}

function VirtualHardDiskUpdate($name, $containerName, $yamlFile) {
		Execute-WssdCommand -Arguments  "storage vhd update --name $name --config $yamlFile --container $containerName"
}

function CreateSampleVirtualHardDisk($containerName) {
	$pwd = (pwd).Path
$yaml = @"
name: $Global:sampleVirtualHardDisk
virtualharddiskproperties:
  source: $Global:sampleVirtualHardDiskSource	
"@
		$yamlFile = "testVirtualHardDisk.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		VirtualHardDiskCreate -yamlFile $yamlFile -container $containerName
}

function CreateSampleVirtualHardDiskDataDisk($containerName) {

	$yaml = @"
name: $Global:sampleVirtualHardDiskDataDisk
virtualharddiskproperties:
  disksizebytes: 10737418240
  dynamic: true
  blocksizebytes: 33554432
  logicalsectorbytes: 4096
  physicalsectorbytes: 4096
  virtualmachinename: ""
  virtualharddisktype: DATADISK_VIRTUALHARDDISK	
"@
	$yamlFile = "testVirtualHardDiskDataDisk.yaml"
	Set-Content -Path $yamlFile -Value $yaml 

	VirtualHardDiskCreate -yamlFile $yamlFile -container $containerName
}

function DeleteSampleVirtualHardDisk($containerName) {
	VirtualHardDiskDelete $Global:sampleVirtualHardDisk -container $containerName
}

function DeleteSampleVirtualHardDiskDataDisk($containerName) {
	VirtualHardDiskDelete $Global:sampleVirtualHardDiskDataDisk -container $containerName
}

function AttachVirtualHardDiskDataDisk($name, $vmName, $containerName) {
	Execute-WssdCommand -Arguments  "storage vhd attach --name $name --vm-name $vmName --container $containerName" 
	IsVirtualHardDiskAttached -name $name  -container $containerName -vmName   $vmName
}

function ResizeVirtualHardDiskDataDisk($name, $sizeBytes, $containerName) {
	Execute-WssdCommand -Arguments  "storage vhd resize --size-bytes $sizeBytes --name $name  --container $containerName" 
}

function DetachVirtualHardDiskDataDisk($name, $vmName, $containerName) {
	Execute-WssdCommand -Arguments  "storage vhd detach --name $name --container $containerName" 
	IsVirtualHardDiskDetached -name $name  -container $containerName
}

function IsVirtualHardDiskAttached($name, $vmName,  $container) {
	$out = VirtualHardDiskShow -name $name -container $container
	if (($out -Match "virtualmachinename: $vmName").Count -gt 0) {
		return
	}
	throw "VirtualHardDisk $name is not attached to the VirtualMachine $vmName"
}

function IsVirtualHardDiskDetached($name,  $container) {
	$out = VirtualHardDiskShow -name $name  -container $container
	if ( ($out -Match "virtualmachinename: """"").Count -gt 0) {
		return
	}
	throw "VirtualHardDisk $name is not detached from the VirtualMachine -  $out"
}



Export-ModuleMember VirtualHardDiskCreate
Export-ModuleMember VirtualHardDiskDelete
Export-ModuleMember VirtualHardDiskShow
Export-ModuleMember VirtualHardDiskList
Export-ModuleMember VirtualHardDiskUpdate
Export-ModuleMember CreateSampleVirtualHardDisk
Export-ModuleMember CreateSampleVirtualHardDiskDataDisk
Export-ModuleMember DeleteSampleVirtualHardDiskDataDisk
Export-ModuleMember DeleteSampleVirtualHardDisk
Export-ModuleMember AttachVirtualHardDiskDataDisk
Export-ModuleMember DetachVirtualHardDiskDataDisk
Export-ModuleMember ResizeVirtualHardDiskDataDisk
Export-ModuleMember CreateVMMSVhd
Export-ModuleMember CleanupVMMSVhd
Export-ModuleMember IsVirtualHardDiskAttached
Export-ModuleMember IsVirtualHardDiskDetached
