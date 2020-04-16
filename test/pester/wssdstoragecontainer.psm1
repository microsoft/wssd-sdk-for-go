

$ScriptPath = Split-Path $MyInvocation.MyCommand.Path
import-module "$PSScriptRoot\wssdglobal.psm1" -Force -Verbose:$false -DisableNameChecking

	$Global:sampleContainer = "sampleContainer"
	$Global:sampleContainerPath = "c:/wssdimagestore"

function ContainerCreate($yamlFile) {
		Execute-WssdCommand -Arguments  "storage container create --config $yamlFile"
}

function ContainerDelete($name) {
		Execute-WssdCommand -Arguments  "storage container delete --name $name"
}

function ContainerShow($name) {
		Execute-WssdCommand -Arguments  "storage container show --name $name"
}

function ContainerList() {
		Execute-WssdCommand -Arguments  "storage container list"
}

function ContainerUpdate($name, $yamlFile) {
		Execute-WssdCommand -Arguments  "storage container update --name $name --config $yamlFile"
}

function CreateSampleContainer() {
	$pwd = (pwd).Path
$yaml = @"
name: $Global:sampleContainer
containerproperties:
  path: $Global:sampleContainerPath	
"@
		$yamlFile = "testContainer.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		ContainerCreate $yamlFile
}

function DeleteSampleContainer() {
	ContainerDelete $Global:sampleContainer
}

Export-ModuleMember ContainerCreate
Export-ModuleMember ContainerDelete
Export-ModuleMember ContainerShow
Export-ModuleMember ContainerList
Export-ModuleMember ContainerUpdate
Export-ModuleMember CreateSampleContainer
Export-ModuleMember DeleteSampleContainer