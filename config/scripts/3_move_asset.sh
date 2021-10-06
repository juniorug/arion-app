#!/bin/bash

peer chaincode invoke -n arion -c '{"Args":["MoveAssetItem", "1","3", "2", "2", "50", "5", "producing","1", "{}" ]}' -C myc
peer chaincode invoke -n arion -c '{"Args":["AddAssetItem", "3", "1" ]}' -C myc
peer chaincode invoke -n arion -c '{"Args":["MoveAssetItem", "1","4", "2", "2", "50", "5", "producing","1", "{}" ]}' -C myc
peer chaincode invoke -n arion -c '{"Args":["AddAssetItem", "4", "1" ]}' -C myc

peer chaincode invoke -n arion -c '{"Args":["TrackAssetItem", "1"]}' -C myc
peer chaincode invoke -n arion -c '{"Args":["TrackAssetItem", "3"]}' -C myc