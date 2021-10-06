#!/bin/bash

peer chaincode invoke -n arion -c '{"Args":["AddActor", "1", "1" ]}' -C myc
peer chaincode invoke -n arion -c '{"Args":["AddActor", "2", "1" ]}' -C myc
peer chaincode invoke -n arion -c '{"Args":["AddStep", "1", "1" ]}' -C myc
peer chaincode invoke -n arion -c '{"Args":["AddStep", "2", "1" ]}' -C myc
peer chaincode invoke -n arion -c '{"Args":["AddAssetItem", "1", "1" ]}' -C myc
peer chaincode invoke -n arion -c '{"Args":["AddAssetItem", "2", "1" ]}' -C myc