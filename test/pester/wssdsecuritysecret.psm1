
$script:wssdcloutctl 
$ScriptPath = Split-Path $MyInvocation.MyCommand.Path
import-module "$PSScriptRoot\wssdglobal.psm1" -Force -Verbose:$false -DisableNameChecking

function SecretCreate($yamlFile) {
		Execute-WssdCommand -Arguments  "security keyvault secret create --config $yamlFile"
}

function SecretSet($yamlFile) {
		Execute-WssdCommand -Arguments  "security keyvault secret set --config $yamlFile"
}


function SecretDelete($name, $vaultName) {
		Execute-WssdCommand -Arguments  "security keyvault secret delete --name $name --vault-name $vaultName"
}

function SecretShow($name, $vaultName) {
		Execute-WssdCommand -Arguments  "security keyvault secret show --name $name --vault-name $vaultName"
}

function SecretList($vaultName) {
		Execute-WssdCommand -Arguments  "security keyvault secret list --vault-name $vaultName"
}

function SecretUpdate($name, $yamlFile, $vaultName) {
		Execute-WssdCommand -Arguments  "security keyvault secret update --name $name --config $yamlFile --vault-name $vaultName"
}

Export-ModuleMember SecretCreate
Export-ModuleMember SecretSet
Export-ModuleMember SecretDelete
Export-ModuleMember SecretShow
Export-ModuleMember SecretList
Export-ModuleMember SecretUpdate