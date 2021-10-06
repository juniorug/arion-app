#!/bin/bash

peer chaincode invoke -n arion -c '{"Args":["CreateAsset", "1", "Gravel", "[]", "[]", "[]", "{}" ]}' -C myc  | grep chaincodeCmd

peer chaincode invoke -n arion -c '{"Args":["CreateActor", "1", "rawMaterial","James JohnsonX9", "{}" ]}' -C myc | grep chaincodeCmd
peer chaincode invoke -n arion -c '{"Args":["CreateActor", "2", "manufacturer","Donald Jackson", "{}" ]}' -C myc | grep chaincodeCmd

peer chaincode invoke -n arion -c '{"Args":["CreateStep", "1", "exploration","1", "rawMaterial", "{}" ]}' -C myc | grep chaincodeCmd
peer chaincode invoke -n arion -c '{"Args":["CreateStep", "2", "production","2", "manufacturer", "{}" ]}' -C myc | grep chaincodeCmd

peer chaincode invoke -n arion -c '{"Args":["CreateAssetItem", "1", "1", "1", "2020-03-07T15:04:05", "", "", "order initiated","2", "{}" ]}' -C myc | grep chaincodeCmd
peer chaincode invoke -n arion -c '{"Args":["CreateAssetItem", "2", "1", "1", "2020-03-07T15:04:05", "", "", "order initiated","2", "{}" ]}' -C myc | grep chaincodeCmd

#peer chaincode invoke -n arion -c '{"Args":["AddActor", "1", "1" ]}' -C myc | grep chaincodeCmd
#peer chaincode invoke -n arion -c '{"Args":["AddActor", "2", "1" ]}' -C myc | grep chaincodeCmd
#peer chaincode invoke -n arion -c '{"Args":["AddStep", "1", "1" ]}' -C myc | grep chaincodeCmd
#peer chaincode invoke -n arion -c '{"Args":["AddStep", "2", "1" ]}' -C myc | grep chaincodeCmd
#peer chaincode invoke -n arion -c '{"Args":["AddAssetItem", "1", "1" ]}' -C myc | grep chaincodeCmd
#peer chaincode invoke -n arion -c '{"Args":["AddAssetItem", "2", "1" ]}' -C myc | grep chaincodeCmd



#peer chaincode invoke -n arion -c '{"Args":["MoveAssetItem", "1","3", "2", "2", "50", "5", "producing","1", "{}" ]}' -C myc
#peer chaincode invoke -n arion -c '{"Args":["AddAssetItem", "3", "1" ]}' -C myc
#peer chaincode invoke -n arion -c '{"Args":["MoveAssetItem", "1","4", "2", "2", "50", "5", "producing","1", "{}" ]}' -C myc
#peer chaincode invoke -n arion -c '{"Args":["AddAssetItem", "4", "1" ]}' -C myc

#peer chaincode invoke -n arion -c '{"Args":["TrackAssetItem", "1"]}' -C myc
#peer chaincode invoke -n arion -c '{"Args":["TrackAssetItem", "3"]}' -C myc
