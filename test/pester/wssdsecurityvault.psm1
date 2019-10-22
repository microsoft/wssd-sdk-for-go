

$ScriptPath = Split-Path $MyInvocation.MyCommand.Path
import-module "$PSScriptRoot\wssdglobal.psm1" -Force -Verbose:$false -DisableNameChecking



function KeyVaultCreate($yamlFile) {
		Execute-WssdCommand -Arguments  "security keyvault create --config $yamlFile"
}

function KeyVaultDelete($name) {
		Execute-WssdCommand -Arguments  "security keyvault delete --name $name"
}

function KeyVaultShow($name) {
		Execute-WssdCommand -Arguments  "security keyvault show --name $name"
}

function KeyVaultList() {
		Execute-WssdCommand -Arguments  "security keyvault list"
}

function KeyVaultUpdate($name, $yamlFile) {
		Execute-WssdCommand -Arguments  "security keyvault update --name $name --config $yamlFile"
}

function CreateSampleKeyVault() {
	$Global:sampleKeyVault = "sampleKeyVault1"
	$yaml = @"
name: $Global:sampleKeyVault 		
"@
		$yamlFile = "testKeyVault.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		KeyVaultCreate $yamlFile
}

function DeleteSampleKeyVault() {
	KeyVaultDelete  $Global:sampleKeyVault
}

Export-ModuleMember KeyVaultCreate
Export-ModuleMember KeyVaultDelete
Export-ModuleMember KeyVaultShow
Export-ModuleMember KeyVaultList
Export-ModuleMember KeyVaultUpdate
Export-ModuleMember CreateSampleKeyVault
Export-ModuleMember DeleteSampleKeyVault