#  

$ScriptPath = Split-Path $MyInvocation.MyCommand.Path
import-module "$PSScriptRoot\wssdcomputevm.psm1" -Force -Verbose:$false -DisableNameChecking
import-module "$PSScriptRoot\wssdcomputevmss.psm1" -Force -Verbose:$false -DisableNameChecking
import-module "$PSScriptRoot\wssdnetworkvnet.psm1" -Force -Verbose:$false -DisableNameChecking
import-module "$PSScriptRoot\wssdnetworkvnic.psm1" -Force -Verbose:$false -DisableNameChecking
import-module "$PSScriptRoot\wssdstoragevhd.psm1" -Force -Verbose:$false -DisableNameChecking
import-module "$PSScriptRoot\wssdsecurityvault.psm1" -Force -Verbose:$false -DisableNameChecking
import-module "$PSScriptRoot\wssdsecuritysecret.psm1" -Force -Verbose:$false -DisableNameChecking

Describe 'Wssd Agent Pre-Requisite' {
	Context 'Checking for Agent' {
		It 'wssdagent.exe is running' {
			get-process -name 'wssdagent'  | Should be $true
		}

		It 'wssdctl.exe is available' {
			get-command -name 'wssdctl.exe'  | Should be $true
		}
	}
}


Describe 'VirtualNetwork BVT' {
	$script:testVirtualNetwork = "testVirtualNetwork1"

	It 'Should be able to create a virtual network' {
		$yaml = @"
name: $script:testVirtualNetwork
type: "ICS"
"@
		$yamlFile = "testVirtualNetwork.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		VirtualNetworkCreate $yamlFile
	}
	It 'Should be able to list all virtual network' {
		VirtualNetworkList
	}

	It 'Should be able to show a virtual network' {
		VirtualNetworkShow $script:testVirtualNetwork
	}

	It 'Should be able to delete a virtual network' {
		VirtualNetworkDelete $script:testVirtualNetwork
	}
}

Describe 'NetworkInterface BVT' {
	BeforeAll {
		CreateSampleVirtualNetwork
	}

	AfterAll {
		DeleteSampleVirtualNetwork
	}


	<#
	It 'Should be able to create a network interface with an IPAddress' {
	$script:testNetworkInterface = "testNetworkInterface1"
		$yaml = @"
name: $script:testNetworkInterface
virtualnetworkinterfaceproperties:
  virtualnetworkname: $Global:sampleVirtualNetwork
  ipconfigurations:
  - name: test
    ipconfigurationproperties:
      ipaddress: "192.168.1.188"
      prefixlength: "24"		
"@
		$yamlFile = "testNetworkInterface.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		NetworkInterfaceCreate $yamlFile
	}
	#>
	It 'Should be able to create a network interface without specifying an IPAddress' {
	$script:testNetworkInterface1 = "testNetworkInterface1"
		$yaml = @"
name: $script:testNetworkInterface1
virtualnetworkinterfaceproperties:
  virtualnetworkname: $Global:sampleVirtualNetwork
"@
		$yamlFile = "testNetworkInterface.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		NetworkInterfaceCreate $yamlFile
	}

	It 'Should be able to list all network interface' {
		NetworkInterfaceList
	}

	It 'Should be able to show a network interface' {
		NetworkInterfaceShow $script:testNetworkInterface
	}

	It 'Should be able to delete a network interface' {
		NetworkInterfaceDelete $script:testNetworkInterface
	}
}

Describe 'VirtualHardDisk BVT' {
	BeforeAll {
		$script:testVirtualHardDiskSource = "./test1.vhdx"
		New-VHD $script:testVirtualHardDiskSource -SizeBytes 4MB
	}
	AfterAll {
		del $script:testVirtualHardDiskSource
	}

	$script:testVirtualHardDisk = "testVirtualHardDisk1"

	It 'Should be able to create a virtual hard disk' {
		$yaml = @"
name: $script:testVirtualHardDisk
virtualharddiskproperties:
  source: $script:testVirtualHardDiskSource	
"@
		$yamlFile = "testVirtualHardDisk.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		VirtualHardDiskCreate $yamlFile
	}
	It 'Should be able to list all virtual hard disk' {
		VirtualHardDiskList
	}

	It 'Should be able to show a virtual hard disk' {
		VirtualHardDiskShow $script:testVirtualHardDisk
	}

	It 'Should be able to delete a virtual hard disk' {
		VirtualHardDiskDelete $script:testVirtualHardDisk
	}
}

Describe 'KeyVault BVT' {
	$script:testKeyVault = "testKeyVault1"

	It 'Should be able to create a keyvault' {
		$yaml = @"
name: $script:testKeyVault		
"@
		$yamlFile = "testKeyVault.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		KeyVaultCreate $yamlFile
	}
	It 'Should be able to list all keyvault' {
		KeyVaultList
	}

	It 'Should be able to show a keyvault' {
		KeyVaultShow $script:testKeyVault
	}

	It 'Should be able to delete a keyvault' {
		KeyVaultDelete $script:testKeyVault
	}
}

Describe 'Secret BVT' {
	$script:testSecret = "testSecret1"
	BeforeAll {
		CreateSampleKeyVault
	}
	AfterAll {
		DeleteSampleKeyVault
	}

	It 'Should be able to set a secret' {
		$yaml = @"
name: $script:testSecret
secretproperties:
  vaultname: $Global:sampleKeyVault
  value: test	
"@
		$yamlFile = "testSecret.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		SecretSet -yamlFile $yamlFile -vaultName $Global:sampleKeyVault
	}
	It 'Should be able to list all secret' {
		SecretList -vaultName $Global:sampleKeyVault
	}

	It 'Should be able to show a secret' {
		SecretShow -name  $script:testSecret -vaultName $Global:sampleKeyVault
	}

	It 'Should be able to delete a secret' {
		SecretDelete -name $script:testSecret -vaultName $Global:sampleKeyVault
	}
}


Describe 'VirtualMachine BVT' {
	$script:testVirtualMachine = "testVirtualMachine1"

	BeforeAll {
		CreateSampleVirtualNetwork
		CreateSampleNetworkInterface
		CreateSampleVirtualHardDisk
	}

	AfterAll {
		DeleteSampleNetworkInterface
		DeleteSampleVirtualNetwork
		DeleteSampleVirtualHardDisk
	}

	It 'Should be able to create a virtual machine' {
		$yaml = @"
name: $script:testVirtualMachine
virtualmachineproperties:
  storageprofile:
    osdisk:
      name: null
      ostype: "Linux"
      vhdname: $Global:sampleVirtualHardDisk
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
    - virtualnetworkinterfacereference: $Global:sampleNetworkInterface
"@
		$yamlFile = "testVirtualMachine.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		VirtualMachineCreate $yamlFile
		Start-Sleep 30
	}
	It 'Should be able to list all virtual machine' {
		VirtualMachineList
	}

	It 'Should be able to show a virtual machine' {
		VirtualMachineShow $script:testVirtualMachine
	}

	It 'Should be able to delete a virtual machine' {
		VirtualMachineDelete $script:testVirtualMachine
	}
}

Describe 'VirtualMachineScaleSet BVT' {
	BeforeAll {
		CreateSampleVirtualNetwork
		CreateSampleNetworkInterface
		CreateSampleVirtualHardDisk
	}

	AfterAll {
		DeleteSampleNetworkInterface
		DeleteSampleVirtualNetwork
		DeleteSampleVirtualHardDisk
	}
	$script:testVirtualMachineScaleSet = "testVirtualMachineScaleSet1"

	It 'Should be able to create a virtual machine scale set' {
		$yaml = @"
name: $script:testVirtualMachineScaleSet
sku:
  name: "test"
  capacity: 1
virtualmachinescalesetproperties:
  virtualmachineprofile:
    name: "ubuntuvm"
    virtualmachinescalesetvmprofileproperties:
      storageprofile:
        osdisk:
          name: null
          ostype: "Linux"
          vhdname: $Global:sampleVirtualHardDisk
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
        networkinterfaceconfigurations:
        - virtualmachinescalesetnetworkconfigurationproperties:
            virtualnetworkname: $Global:sampleVirtualNetwork	
"@
		$yamlFile = "testVirtualMachineScaleSet.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		VirtualMachineScaleSetCreate $yamlFile
		Start-Sleep 30
	}
	It 'Should be able to list all virtual machine scale set' {
		VirtualMachineScaleSetList
	}

	It 'Should be able to show a virtual machine scale set' {
		VirtualMachineScaleSetShow $script:testVirtualMachineScaleSet
	}

	It 'Should be able to delete a virtual machine scale set' {
		VirtualMachineScaleSetDelete $script:testVirtualMachineScaleSet
	}
}
