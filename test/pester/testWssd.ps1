#  
param(
    [Parameter()]
    [Switch] $disableTls,
    [Switch] $setupOnly,
    [string] $marinerVhdPath,
    [string] $imageStore = "c:/wssdimagestore"
)

import-module "$PSScriptRoot\wssdcomputevm.psm1" -Force -Verbose:$false -DisableNameChecking
import-module "$PSScriptRoot\wssdcomputevmss.psm1" -Force -Verbose:$false -DisableNameChecking
import-module "$PSScriptRoot\wssdnetworkvnet.psm1" -Force -Verbose:$false -DisableNameChecking
import-module "$PSScriptRoot\wssdnetworkvnic.psm1" -Force -Verbose:$false -DisableNameChecking
import-module "$PSScriptRoot\wssdstoragevhd.psm1" -Force -Verbose:$false -DisableNameChecking
import-module "$PSScriptRoot\wssdstoragecontainer.psm1" -Force -Verbose:$false -DisableNameChecking
import-module "$PSScriptRoot\wssdsecurityvault.psm1" -Force -Verbose:$false -DisableNameChecking
import-module "$PSScriptRoot\wssdsecuritysecret.psm1" -Force -Verbose:$false -DisableNameChecking

if ($disableTls.IsPresent) {
    $Global:debugMode = $true
}

$Global:sampleContainerPath = $imageStore

if ($setupOnly.IsPresent) {
    # Just setup the required modules and return
    return
}

Describe 'Wssd Agent Pre-Requisite' {

	
    Context 'Checking for Agent' {
        It 'wssdagent.exe is running' {
            get-process -name 'wssdagent'  # | Should be $true
        }

        It 'wssdctl.exe is available' {
            get-command -name 'wssdctl.exe'  # | Should be $true
        }
    }
}

Describe 'Container BVT' {
    $script:testContainer = "testContainer1"

    It 'Should be able to create a storage container' {
        $yaml = @"
name: $script:testContainer
containerproperties:
  path: c:/containerpath	
"@
        $yamlFile = "testContainer.yaml"
        Set-Content -Path $yamlFile -Value $yaml 

        ContainerCreate $yamlFile  # | Should Not Throw
    }
    It 'Should be able to list all storage container' {
        ContainerList  # | Should Not Throw
    }
    <#
	# Uncomment once implemented
	It 'Should be able to show a storage container' {
		ContainerShow $script:testContainer  # | Should Not Throw
	}
	#>
    It 'Should not be able to delete a storage container while VHDs are using it' {
        CreateSampleVirtualHardDisk -container $script:testContainer
        { ContainerDelete $script:testContainer  } | Should Throw
        DeleteSampleVirtualHardDisk  -container $script:testContainer
    }

    It 'Should be able to delete a  storage container' {
        ContainerDelete $script:testContainer  # | Should Not Throw
    }
}

Describe 'VirtualNetwork BVT' {
    $script:testVirtualNetwork = "Default Switch"

    It 'Should be able to create a virtual network' {
        $yaml = @"
name: $script:testVirtualNetwork
type: "ICS"
"@
        $yamlFile = "testVirtualNetwork.yaml"
        Set-Content -Path $yamlFile -Value $yaml 

        VirtualNetworkCreate $yamlFile # | Should Not Throw
    }
    It 'Should be able to list all virtual network' {
        VirtualNetworkList  # | Should Not Throw
    }

    It 'Should be able to show a virtual network' {
        VirtualNetworkShow $script:testVirtualNetwork  # | Should Not Throw
    }

    It 'Should be able to delete a virtual network' {
        VirtualNetworkDelete $script:testVirtualNetwork  # | Should Not Throw
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
      subnetid: $Global:sampleVirtualNetwork
"@
		$yamlFile = "testNetworkInterface.yaml"
		Set-Content -Path $yamlFile -Value $yaml 

		NetworkInterfaceCreate $yamlFile  # | Should Not Throw
	}
	#>
    It 'Should be able to create a network interface without specifying an IPAddress' {
        $script:testNetworkInterface1 = "testNetworkInterface1"
        $yaml = @"
name: $script:testNetworkInterface1
virtualnetworkinterfaceproperties:
  ipconfigurations:
  - ipconfigurationproperties:
      subnetid: $Global:sampleVirtualNetwork
"@
        $yamlFile = "testNetworkInterface.yaml"
        Set-Content -Path $yamlFile -Value $yaml 

        NetworkInterfaceCreate $yamlFile  # | Should Not Throw
    }

    It 'Should be able to list all network interface' {
        NetworkInterfaceList  # | Should Not Throw
    }

    It 'Should be able to show a network interface' {
        NetworkInterfaceShow $script:testNetworkInterface1  # | Should Not Throw
    }

    It 'Should be able to delete a network interface' {
        NetworkInterfaceDelete $script:testNetworkInterface1  # | Should Not Throw
    }
}

Describe 'VirtualHardDisk BVT' {
    BeforeAll {
        CreateSampleContainer
        CreateSampleVirtualHardDisk -container $global:sampleContainer
    }
    AfterAll {
        DeleteSampleVirtualHardDisk -container $global:sampleContainer
        DeleteSampleContainer
    }

    $script:testOSVirtualHardDisk = "testOSVirtualHardDisk"
    $script:testCloneVirtualHardDisk = "testcloneVirtualHardDisk"
    $script:testDataVirtualHardDisk = "testDataVirtualHardDisk"

    It 'Should be able to upload a virtual hard disk' {
        $yaml = @"
name: $script:testOSVirtualHardDisk
virtualharddiskproperties:
  source: $Global:testVirtualHardDiskSource	
"@
        $yamlFile = "testUploadVirtualHardDisk.yaml"
        Set-Content -Path $yamlFile -Value $yaml 

        VirtualHardDiskCreate -yamlFile $yamlFile -container $global:sampleContainer  # | Should Not Throw
    }
    It 'Should be able to clone a virtual hard disk' {
        $yaml = @"
name: $script:testCloneVirtualHardDisk
virtualharddiskproperties:
  source: $Global:sampleVirtualHardDisk
"@
        $yamlFile = "testCloneVirtualHardDisk.yaml"
        Set-Content -Path $yamlFile -Value $yaml 

        VirtualHardDiskCreate -yamlFile $yamlFile -container $global:sampleContainer  # | Should Not Throw
    }

    It 'Should be able to create a data virtual hard disk' {
        $yaml = @"
name: $script:testDataVirtualHardDisk
virtualharddiskproperties:
  disksizebytes: 10737418240
  dynamic: true
  blocksizebytes: 33554432
  logicalsectorbytes: 4096
  physicalsectorbytes: 4096
  virtualharddisktype: DATADISK_VIRTUALHARDDISK
"@
        $yamlFile = "testVirtualHardDiskDataDisk.yaml"
        Set-Content -Path $yamlFile -Value $yaml 

        VirtualHardDiskCreate -yamlFile $yamlFile -container $global:sampleContainer 
    }
    It 'Should be able to list all virtual hard disk' {
        VirtualHardDiskList -container $global:sampleContainer  # | Should Not Throw
    }

    It 'Should be able to resize the virtual hard disk data disk' {
        ResizeVirtualHardDiskDataDisk -name $script:testDataVirtualHardDisk -sizeBytes 21474836480 ‬ -container $global:sampleContainer  # | Should Not Throw
    }
    <#
	# Uncomment once implemented
	It 'Should be able to show a virtual hard disk' {
		VirtualHardDiskShow $script:testVirtualHardDisk  # | Should Not Throw
	}
	#>
    It 'Should be able to delete a virtual hard disk' {
        VirtualHardDiskDelete -name $script:testOSVirtualHardDisk  -container $global:sampleContainer # | Should Not Throw
        VirtualHardDiskDelete -name $script:testCloneVirtualHardDisk  -container $global:sampleContainer # | Should Not Throw
        VirtualHardDiskDelete -name $script:testDataVirtualHardDisk  -container $global:sampleContainer # | Should Not Throw
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

        KeyVaultCreate $yamlFile  # | Should Not Throw
    }
    It 'Should be able to list all keyvault' {
        KeyVaultList  # | Should Not Throw
    }

    It 'Should be able to show a keyvault' {
        KeyVaultShow $script:testKeyVault  # | Should Not Throw
    }

    It 'Should be able to delete a keyvault' {
        KeyVaultDelete $script:testKeyVault  # | Should Not Throw
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
value: test
secretproperties:
  vaultname: $Global:sampleKeyVault
"@
        $yamlFile = "testSecret.yaml"
        Set-Content -Path $yamlFile -Value $yaml 

        SecretSet -yamlFile $yamlFile -vaultName $Global:sampleKeyVault  # | Should Not Throw
    }
    It 'Should be able to list all secret' {
        SecretList -vaultName $Global:sampleKeyVault  # | Should Not Throw
    }

    It 'Should be able to show a secret' {
        SecretShow -name  $script:testSecret -vaultName $Global:sampleKeyVault  # | Should Not Throw
    }

    It 'Should be able to delete a secret' {
        SecretDelete -name $script:testSecret -vaultName $Global:sampleKeyVault  # | Should Not Throw
    }
}


Describe 'VirtualMachine BVT' {
    $script:testVirtualMachine = "testVirtualMachine1"

    BeforeAll {
        CreateSampleContainer
        CreateSampleVirtualNetwork
        CreateSampleNetworkInterface
        CreateSampleVirtualHardDisk -container $global:sampleContainer
        CreateSampleVirtualHardDiskDataDisk -container $global:sampleContainer
    }

    AfterAll {
        DeleteSampleNetworkInterface
        DeleteSampleVirtualNetwork
        DeleteSampleVirtualHardDiskDataDisk -container $global:sampleContainer
        DeleteSampleVirtualHardDisk -container $global:sampleContainer
        DeleteSampleContainer
    }

    It 'Should be able to list all virtual machine' {
        VirtualMachineList  # | Should Not Throw
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

        VirtualMachineCreate $yamlFile  # | Should Not Throw
    }
    It 'Should be able to list all virtual machine' {
        VirtualMachineList  # | Should Not Throw
    }

    It 'Should be able to show a virtual machine' {
        VirtualMachineShow $script:testVirtualMachine  # | Should Not Throw
    }

    It 'Should be able to attach a data disk to a virtual machine' {
        AttachVirtualHardDiskDataDisk -name $Global:sampleVirtualHardDiskDataDisk -container $global:sampleContainer -vmName   $script:testVirtualMachine
    }

    It 'Should be able to resize the virtual hard disk data disk while attached' {
        ResizeVirtualHardDiskDataDisk -name $Global:sampleVirtualHardDiskDataDisk -sizeBytes 21474836480 -container $global:sampleContainer  # | Should Not Throw
    }
    It 'Should be able to detach a data disk to a virtual machine' {
        DetachVirtualHardDiskDataDisk -name $Global:sampleVirtualHardDiskDataDisk -container $global:sampleContainer -vmName   $script:testVirtualMachine
    }

    It 'Should be able to delete a virtual machine' {
        VirtualMachineDelete $script:testVirtualMachine  # | Should Not Throw
    }
}

Describe 'VirtualMachineScaleSet BVT' {
    BeforeAll {
        CreateSampleContainer
        CreateSampleVirtualNetwork
        CreateSampleVirtualHardDisk -container $global:sampleContainer
    }

    AfterAll {
        DeleteSampleVirtualHardDisk -container $global:sampleContainer
        DeleteSampleVirtualNetwork
        DeleteSampleContainer
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
            ipconfigurations:
            - ipconfigurationproperties:
                subnetid: $Global:sampleVirtualNetwork
"@
        $yamlFile = "testVirtualMachineScaleSet.yaml"
        Set-Content -Path $yamlFile -Value $yaml 

        VirtualMachineScaleSetCreate $yamlFile  # | Should Not Throw
    }
    It 'Should be able to list all virtual machine scale set' {
        VirtualMachineScaleSetList  # | Should Not Throw
    }

    It 'Should be able to show a virtual machine scale set' {
        VirtualMachineScaleSetShow $script:testVirtualMachineScaleSet  # | Should Not Throw
    }

    It 'Should be able to delete a virtual machine scale set' {
        VirtualMachineScaleSetDelete $script:testVirtualMachineScaleSet  # | Should Not Throw
    }
}
