# Copyright (c) Microsoft Corporation.
# Licensed under the MIT license.

Module="network"
echo "Generating $Module protoc"
protoc -I $Module/common $Module/common/common.proto --go_out=plugins=grpc:$Module
protoc -I $Module/virtualnetwork -I $Module/common $Module/virtualnetwork/virtualnetwork.proto --go_out=plugins=grpc:$Module
protoc -I $Module/loadbalancer -I $Module/common   $Module/loadbalancer/loadbalancer.proto --go_out=plugins=grpc:$Module
protoc -I $Module/virtualnetworkinterface -I $Module/common $Module/virtualnetworkinterface/virtualnetworkinterface.proto --go_out=plugins=grpc:$Module

# Generate compute agent protoc
Module="compute"
echo "Generating $Module protoc"
protoc -I $Module/common $Module/common/common.proto --go_out=plugins=grpc:$Module
protoc -I $Module/virtualmachine -I $Module/common $Module/virtualmachine/virtualmachine.proto --go_out=plugins=grpc:$Module
protoc -I $Module/virtualmachinescaleset -I $Module/virtualmachine -I $Module/common $Module/virtualmachinescaleset/virtualmachinescaleset.proto --go_out=plugins=grpc:$Module

Module="storage"
echo "Generating $Module protoc"
protoc -I $Module/common $Module/common/common.proto --go_out=plugins=grpc:$Module
protoc -I $Module/virtualharddisk -I $Module/common $Module/virtualharddisk/virtualharddisk.proto  --go_out=plugins=grpc:$Module

Module="security"
echo "Generating $Module protoc"
protoc -I $Module/common $Module/common/common.proto --go_out=plugins=grpc:$Module
protoc -I $Module/keyvault/secret -I $Module/common $Module/keyvault/secret/secret.proto  --go_out=plugins=grpc:$Module
protoc -I $Module/keyvault -I $Module/common -I $Module/keyvault/secret $Module/keyvault/keyvault.proto  --go_out=plugins=grpc:$Module
