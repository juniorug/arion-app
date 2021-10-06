/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// AssetTransferSmartContract provides functions for managing an asset
type AssetTransferSmartContract struct {
	contractapi.Contract
}

// Actor describes basic details of an actor in the SCM
type Actor struct {
	ActorID          string            `json:"actorID"`
	ActorType        string            `json:"actorType"`
	ActorName        string            `json:"actorName"`
	Deleted          bool              `json:"deleted"`
	AditionalInfoMap map[string]string `json:"aditionalInfoMap"`
	//assetID        string            `json:"assetID"`
}

// Step describes basic details of an step in the SCM
type Step struct {
	StepID           string            `json:"stepID"`
	StepName         string            `json:"stepName"`
	StepOrder        uint              `json:"stepOrder"`
	ActorType        string            `json:"actorType"`
	Deleted          bool              `json:"deleted"`
	AditionalInfoMap map[string]string `json:"aditionalInfoMap"`
}

// AssetItem describes ths single item in the SCM
type AssetItem struct {
	AssetItemID      string            `json:"assetItemID"`
	OwnerID          string            `json:"ownerID"`
	StepID           string            `json:"stepID"`
	ParentID         string            `json:"parentID"`
	ProcessDate      string            `json:"processDate"`  //(date which currenct actor acquired the item)
	DeliveryDate     string            `json:"deliveryDate"` //(date which currenct actor received the item)
	OrderPrice       string            `json:"orderPrice"`
	ShippingPrice    string            `json:"shippingPrice"`
	Status           string            `json:"status"`
	Quantity         string            `json:"quantity"`
	Deleted          bool              `json:"deleted"`
	AditionalInfoMap map[string]string `json:"aditionalInfoMap"`
}

// Asset is the collection of assetItems in the SCM
type Asset struct {
	AssetID          string            `json:"assetID"`
	AssetName        string            `json:"assetName"`
	AssetItems       []AssetItem       `json:"assetItems"`
	Actors           []Actor           `json:"actors"`
	Steps            []Step            `json:"steps"`
	Deleted          bool              `json:"deleted"`
	AditionalInfoMap map[string]string `json:"aditionalInfoMap"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *AssetItem
}

// InitLedger adds a base set of items to the ledger
func (s *AssetTransferSmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	aditionalInfoMap := make(map[string]string)

	actors := []Actor{
		{ActorID: "1", ActorType: "rawMaterial", ActorName: "James Johnson", AditionalInfoMap: aditionalInfoMap},
		{ActorID: "2", ActorType: "manufacturer", ActorName: "Helena Smith", AditionalInfoMap: aditionalInfoMap},
	}

	for i, actor := range actors {
		actorAsBytes, _ := json.Marshal(actor)
		err := ctx.GetStub().PutState("ACTOR_"+strconv.Itoa(i), actorAsBytes)
		if err != nil {
			return fmt.Errorf("Failed to put the actor[%d] to world state. %s", i, err.Error())
		}
	}

	steps := []Step{
		{StepID: "1", StepName: "exploration", StepOrder: 1, ActorType: "rawMaterial", AditionalInfoMap: aditionalInfoMap},
		{StepID: "2", StepName: "Production", StepOrder: 2, ActorType: "manufacturer", AditionalInfoMap: aditionalInfoMap},
	}

	for i, step := range steps {
		stepAsBytes, _ := json.Marshal(step)
		err := ctx.GetStub().PutState("STEP_"+strconv.Itoa(i), stepAsBytes)
		if err != nil {
			return fmt.Errorf("Failed to put the step[%d] to world state. %s", i, err.Error())
		}
	}

	assetItems := []AssetItem{
		{
			AssetItemID:      "1",
			OwnerID:          "1",
			StepID:           "1",
			ParentID:         "0",
			ProcessDate:      "2020-03-07T15:04:05",
			DeliveryDate:     "",
			OrderPrice:       "",
			ShippingPrice:    "",
			Status:           "order initiated",
			Quantity:         "",
			Deleted:          false,
			AditionalInfoMap: aditionalInfoMap,
		},
		{
			AssetItemID:      "2",
			OwnerID:          "1",
			StepID:           "1",
			ParentID:         "0",
			ProcessDate:      "2020-03-07T15:04:05",
			DeliveryDate:     "",
			OrderPrice:       "",
			ShippingPrice:    "",
			Status:           "order initiated",
			Quantity:         "",
			Deleted:          false,
			AditionalInfoMap: aditionalInfoMap,
		},
	}

	for i, assetItem := range assetItems {
		assetItemAsBytes, _ := json.Marshal(assetItem)
		err := ctx.GetStub().PutState("ASSET_ITEM_"+strconv.Itoa(i), assetItemAsBytes)
		if err != nil {
			return fmt.Errorf("Failed to put the assetItem[%d] to world state. %s", i, err.Error())
		}
	}

	assets := []Asset{
		{
			AssetID:          "1",
			AssetName:        "Gravel",
			AssetItems:       assetItems,
			Actors:           actors,
			Steps:            steps,
			AditionalInfoMap: aditionalInfoMap,
		},
	}

	for i, asset := range assets {
		assetAsBytes, _ := json.Marshal(asset)
		err := ctx.GetStub().PutState("ASSET_"+strconv.Itoa(i), assetAsBytes)
		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateActor adds a new Actor to the world state with given details
func (s *AssetTransferSmartContract) CreateActor(ctx contractapi.TransactionContextInterface, actorID string, actorType string, actorName string, aditionalInfoMap map[string]string) error {
	log.Print("[CreateActor] called with actorID: ", actorID)
	actorJSON, err := ctx.GetStub().GetState("ACTOR_" + actorID)

	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	if actorJSON != nil {
		return fmt.Errorf("The actor %s already exists", actorID)
	}

	actor := Actor{
		ActorID:          actorID,
		ActorType:        actorType,
		ActorName:        actorName,
		Deleted:          false,
		AditionalInfoMap: aditionalInfoMap,
	}

	actorAsBytes, err := json.Marshal(actor)
	if err != nil {
		return err
	}

	log.Print("[CreateActor] beinging saved to Stub.state and returning its response")
	return ctx.GetStub().PutState("ACTOR_"+actorID, actorAsBytes)
}

// CreateStep adds a new Step to the world state with given details
func (s *AssetTransferSmartContract) CreateStep(ctx contractapi.TransactionContextInterface, stepID string, stepName string, stepOrder uint, actorType string, aditionalInfoMap map[string]string) error {
	log.Print("[CreateStep] called with stepID: ", stepID)
	stepJSON, err := ctx.GetStub().GetState("STEP_" + stepID)

	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	if stepJSON != nil {
		return fmt.Errorf("The step %s already exists", stepID)
	}

	step := Step{
		StepID:           stepID,
		StepName:         stepName,
		StepOrder:        stepOrder,
		ActorType:        actorType,
		Deleted:          false,
		AditionalInfoMap: aditionalInfoMap,
	}

	stepAsBytes, err := json.Marshal(step)
	if err != nil {
		return err
	}

	log.Print("[CreateStep] beinging saved to Stub.state and returning its response")
	return ctx.GetStub().PutState("STEP_"+stepID, stepAsBytes)
}

// CreateAssetItem adds a new AssetItem to the world state with given details
func (s *AssetTransferSmartContract) CreateAssetItem(ctx contractapi.TransactionContextInterface, assetItemID string, ownerID string, stepId string, deliveryDate string, orderPrice string, shippingPrice string, status string, quantity string, aditionalInfoMap map[string]string) error {
	log.Print("[CreateAssetItem] called with assetItemID: ", assetItemID)
	assetItemJSON, err := ctx.GetStub().GetState("ASSET_ITEM_" + assetItemID)

	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	if assetItemJSON != nil {
		return fmt.Errorf("The assetItem %s already exists", assetItemID)
	}

	assetItem := AssetItem{
		AssetItemID:      assetItemID,
		OwnerID:          ownerID,
		StepID:           "1",
		ParentID:         "0",
		ProcessDate:      time.Now().Format("2006-01-02 15:04:05"),
		DeliveryDate:     deliveryDate,
		OrderPrice:       orderPrice,
		ShippingPrice:    shippingPrice,
		Status:           status,
		Quantity:         quantity,
		Deleted:          false,
		AditionalInfoMap: aditionalInfoMap,
	}

	assetItemAsBytes, err := json.Marshal(assetItem)
	if err != nil {
		return err
	}

	log.Print("[CreateAssetItem] beinging saved to Stub.state and returning its response")
	return ctx.GetStub().PutState("ASSET_ITEM_"+assetItemID, assetItemAsBytes)
}

// CreateAsset adds a new Asset with empty arrays to the world state with given details
func (s *AssetTransferSmartContract) CreateEmptyAsset(ctx contractapi.TransactionContextInterface, assetID string, assetName string, aditionalInfoMap map[string]string) error {
	assetJSON, err := ctx.GetStub().GetState("ASSET_" + assetID)

	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	if assetJSON != nil {
		return fmt.Errorf("The asset %s already exists", assetID)
	}

	asset := Asset{
		AssetID:          assetID,
		AssetName:        assetName,
		AssetItems:       []AssetItem{},
		Actors:           []Actor{},
		Steps:            []Step{},
		Deleted:          false,
		AditionalInfoMap: aditionalInfoMap,
	}

	assetAsBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("ASSET_"+assetID, assetAsBytes)
}

// CreateAsset adds a new Asset to the world state with given details
func (s *AssetTransferSmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, assetID string, assetName string, assetItems []AssetItem, actors []Actor, steps []Step, aditionalInfoMap map[string]string) error {
	assetJSON, err := ctx.GetStub().GetState("ASSET_" + assetID)

	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	if assetJSON != nil {
		return fmt.Errorf("The asset %s already exists", assetID)
	}

	asset := Asset{
		AssetID:          assetID,
		AssetName:        assetName,
		AssetItems:       assetItems,
		Actors:           actors,
		Steps:            steps,
		Deleted:          false,
		AditionalInfoMap: aditionalInfoMap,
	}

	assetAsBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("ASSET_"+assetID, assetAsBytes)
}

// QueryActor returns the actor stored in the world state with given id
func (s *AssetTransferSmartContract) QueryActor(ctx contractapi.TransactionContextInterface, actorID string) (*Actor, error) {
	actorAsBytes, err := ctx.GetStub().GetState("ACTOR_" + actorID)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if actorAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", actorID)
	}

	actor := new(Actor)
	err = json.Unmarshal(actorAsBytes, actor)
	if err != nil {
		return nil, err
	}

	return actor, nil
}

// QueryStep returns the step stored in the world state with given id
func (s *AssetTransferSmartContract) QueryStep(ctx contractapi.TransactionContextInterface, stepID string) (*Step, error) {
	stepAsBytes, err := ctx.GetStub().GetState("STEP_" + stepID)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if stepAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", stepID)
	}

	step := new(Step)
	err = json.Unmarshal(stepAsBytes, step)
	if err != nil {
		return nil, err
	}

	return step, nil
}

// QueryAssetItem returns the AssetItem stored in the world state with given id
func (s *AssetTransferSmartContract) QueryAssetItem(ctx contractapi.TransactionContextInterface, assetItemID string) (*AssetItem, error) {
	assetItemAsBytes, err := ctx.GetStub().GetState("ASSET_ITEM_" + assetItemID)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if assetItemAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", assetItemID)
	}

	assetItem := new(AssetItem)
	err = json.Unmarshal(assetItemAsBytes, assetItem)
	if err != nil {
		return nil, err
	}

	return assetItem, nil
}

// QueryAsset returns the Asset stored in the world state with given id
func (s *AssetTransferSmartContract) QueryAsset(ctx contractapi.TransactionContextInterface, assetID string) (*Asset, error) {
	log.Print("[QueryAsset] called with assetId: ", assetID)
	fullAssetID := "ASSET_" + assetID
	log.Print("[QueryAsset] fullAssetID: ", fullAssetID)
	assetAsBytes, err := ctx.GetStub().GetState(fullAssetID)
	log.Print("[QueryAsset] assetAsBytes called: ", assetAsBytes)
	log.Print("[QueryAsset] will check if err != nil")
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	log.Print("[QueryAsset] after check if err != nil. Now will check if assetAsBytes == nil ")
	if assetAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", assetID)
	}

	log.Print("[QueryAsset] after all validations will marshal asset to send it back")
	asset := new(Asset)
	err = json.Unmarshal(assetAsBytes, asset)
	log.Print("[QueryAsset] after asset marshal. will verify if err is != nil")
	if err != nil {
		return nil, err
	}
	log.Print("[QueryAsset] after all validations with no errors. Will send the asset back")
	return asset, nil
}

//QueryAllActors returns the actors' list stored in the world state
func (s *AssetTransferSmartContract) QueryAllActors(ctx contractapi.TransactionContextInterface) ([]*Actor, error) {

	actorsIterator, err := ctx.GetStub().GetStateByRange("ACTOR_0", "ACTOR_9999999999999999999")
	if err != nil {
		return nil, err
	}

	defer actorsIterator.Close()

	var actors []*Actor
	for actorsIterator.HasNext() {
		actorResponse, err := actorsIterator.Next()
		if err != nil {
			return nil, err
		}

		var actor *Actor
		err = json.Unmarshal(actorResponse.Value, &actor)
		if err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}

	return actors, nil
}

//QueryAllSteps returns the steps' list stored in the world state
func (s *AssetTransferSmartContract) QueryAllSteps(ctx contractapi.TransactionContextInterface) ([]*Step, error) {

	stepsIterator, err := ctx.GetStub().GetStateByRange("STEP_0", "STEP_9999999999999999999")
	if err != nil {
		return nil, err
	}

	defer stepsIterator.Close()

	var steps []*Step
	for stepsIterator.HasNext() {
		stepResponse, err := stepsIterator.Next()
		if err != nil {
			return nil, err
		}

		var step *Step
		err = json.Unmarshal(stepResponse.Value, &step)
		if err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}

	return steps, nil
}

//QueryAllAssetItems returns the asset items' list stored in the world state
func (s *AssetTransferSmartContract) QueryAllAssetItems(ctx contractapi.TransactionContextInterface) ([]*AssetItem, error) {

	assetItemsIterator, err := ctx.GetStub().GetStateByRange("ASSET_ITEM_0", "ASSET_ITEM_9999999999999999999")
	if err != nil {
		return nil, err
	}

	defer assetItemsIterator.Close()

	var assetItems []*AssetItem
	for assetItemsIterator.HasNext() {
		assetItemResponse, err := assetItemsIterator.Next()
		if err != nil {
			return nil, err
		}

		var assetItem *AssetItem
		err = json.Unmarshal(assetItemResponse.Value, &assetItem)
		if err != nil {
			return nil, err
		}
		assetItems = append(assetItems, assetItem)
	}

	return assetItems, nil
}

//QueryAllAssets returns the assets' list stored in the world state
func (s *AssetTransferSmartContract) QueryAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {

	assetsIterator, err := ctx.GetStub().GetStateByRange("ASSET_0", "ASSET_9999999999999999999")
	if err != nil {
		return nil, err
	}

	defer assetsIterator.Close()

	var assets []*Asset
	for assetsIterator.HasNext() {
		assetResponse, err := assetsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset *Asset
		err = json.Unmarshal(assetResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

// UpdateActor updates an existing Actor to the world state with given details
func (s *AssetTransferSmartContract) UpdateActor(ctx contractapi.TransactionContextInterface, actorID string, actorType string, actorName string, aditionalInfoMap map[string]string) error {
	actorJSON, err := ctx.GetStub().GetState("ACTOR_" + actorID)

	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	if actorJSON == nil {
		return fmt.Errorf("The actor %s does not exists", actorID)
	}

	actor := Actor{
		ActorID:          actorID,
		ActorType:        actorType,
		ActorName:        actorName,
		Deleted:          false,
		AditionalInfoMap: aditionalInfoMap,
	}

	actorAsBytes, err := json.Marshal(actor)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("ACTOR_"+actorID, actorAsBytes)
}

// UpdateStep updates an existing Step to the world state with given details
func (s *AssetTransferSmartContract) UpdateStep(ctx contractapi.TransactionContextInterface, stepID string, stepName string, stepOrder uint, actorType string, aditionalInfoMap map[string]string) error {
	stepJSON, err := ctx.GetStub().GetState("STEP_" + stepID)

	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	if stepJSON == nil {
		return fmt.Errorf("The step %s does not exists", stepID)
	}

	step := Step{
		StepID:           stepID,
		StepName:         stepName,
		StepOrder:        stepOrder,
		ActorType:        actorType,
		Deleted:          false,
		AditionalInfoMap: aditionalInfoMap,
	}

	stepAsBytes, err := json.Marshal(step)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("STEP_"+stepID, stepAsBytes)
}

// UpdateAssetItem updates an existing AssetItem to the world state with given details
func (s *AssetTransferSmartContract) UpdateAssetItem(ctx contractapi.TransactionContextInterface, assetItemID string, ownerID string, stepID string, parentID string, deliveryDate string, orderPrice string, shippingPrice string, status string, quantity string, aditionalInfoMap map[string]string) error {
	assetItemJSON, err := ctx.GetStub().GetState("ASSET_ITEM_" + assetItemID)

	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	if assetItemJSON == nil {
		return fmt.Errorf("The assetItem %s does not exists", assetItemID)
	}

	assetItem := AssetItem{
		AssetItemID:      assetItemID,
		OwnerID:          ownerID,
		StepID:           stepID,
		ParentID:         parentID,
		ProcessDate:      time.Now().Format("2006-01-02 15:04:05"),
		DeliveryDate:     deliveryDate,
		OrderPrice:       orderPrice,
		ShippingPrice:    shippingPrice,
		Status:           status,
		Quantity:         quantity,
		Deleted:          false,
		AditionalInfoMap: aditionalInfoMap,
	}

	assetItemAsBytes, err := json.Marshal(assetItem)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("ASSET_ITEM_"+assetItemID, assetItemAsBytes)
}

// UpdateAsset updates an existing Asset to the world state with given details
func (s *AssetTransferSmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, assetID string, assetName string, assetItems []AssetItem, actors []Actor, steps []Step, aditionalInfoMap map[string]string) error {
	assetJSON, err := ctx.GetStub().GetState("ASSET_" + assetID)

	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	if assetJSON == nil {
		return fmt.Errorf("The asset %s does not exists", assetID)
	}

	asset := Asset{
		AssetID:          assetID,
		AssetName:        assetName,
		AssetItems:       assetItems,
		Actors:           actors,
		Steps:            steps,
		Deleted:          false,
		AditionalInfoMap: aditionalInfoMap,
	}

	assetAsBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("ASSET_"+assetID, assetAsBytes)
}

// AddActor updates an existing Asset by adding a new actor to the world state with given details
func (s *AssetTransferSmartContract) AddActor(ctx contractapi.TransactionContextInterface, actorID string, assetID string) error {
	log.Print("[AddActor] called with actorID: ", actorID, " and assetId: ", assetID)
	/* 	assetJSON, err := ctx.GetStub().GetState("ASSET_" + assetID)
	   	if err != nil {
	   		return fmt.Errorf("Failed to read the data from world state: %s", err)
	   	}

	   	if assetJSON == nil {
	   		return fmt.Errorf("The asset %s does not exists", assetID)
	   	} */

	asset := new(Asset)
	asset, err := s.QueryAsset(ctx, assetID)
	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	actor := new(Actor)
	actor, err = s.QueryActor(ctx, actorID)
	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}
	asset.Actors = append(asset.Actors, *actor)

	assetAsBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	log.Print("[AddActor] Saving Asset data after add actor and returning its response")
	return ctx.GetStub().PutState("ASSET_"+assetID, assetAsBytes)
}

// AddStep updates an existing Asset by adding a new step to the world state with given details
func (s *AssetTransferSmartContract) AddStep(ctx contractapi.TransactionContextInterface, stepID string, assetID string) error {
	log.Print("[AddStep] called with stepID: ", stepID, " and assetId: ", assetID)
	/* assetJSON, err := ctx.GetStub().GetState("ASSET_" + assetID)

	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	if assetJSON == nil {
		return fmt.Errorf("The asset %s does not exists", assetID)
	} */

	asset := new(Asset)
	asset, err := s.QueryAsset(ctx, assetID)
	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	step := new(Step)
	step, err = s.QueryStep(ctx, stepID)
	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}
	asset.Steps = append(asset.Steps, *step)

	assetAsBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	log.Print("[AddStep] Saving Asset data after add step and returning its response")
	return ctx.GetStub().PutState("ASSET_"+assetID, assetAsBytes)
}

// AddAssetItem updates an existing Asset by adding a new assetItem to the world state with given details
func (s *AssetTransferSmartContract) AddAssetItem(ctx contractapi.TransactionContextInterface, assetItemID string, assetID string) error {
	log.Print("[AddAssetItem] called with AssetItemID: ", assetItemID, " and assetId: ", assetID)
	/* assetJSON, err := ctx.GetStub().GetState("ASSET_" + assetID)

	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	if assetJSON == nil {
		return fmt.Errorf("The asset %s does not exists", assetID)
	} */

	asset := new(Asset)
	asset, err := s.QueryAsset(ctx, assetID)
	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	assetItem := new(AssetItem)
	assetItem, err = s.QueryAssetItem(ctx, assetItemID)
	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}
	asset.AssetItems = append(asset.AssetItems, *assetItem)

	assetAsBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	log.Print("[AddAssetItem] Saving Asset data after add asset item and returning its response")
	return ctx.GetStub().PutState("ASSET_"+assetID, assetAsBytes)
}

// DeleteActor sets the deleted flag as true for the given actor
func (s *AssetTransferSmartContract) DeleteActor(ctx contractapi.TransactionContextInterface, actorID string) error {
	actor, err := s.QueryActor(ctx, "ACTOR_"+actorID)
	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	actor.Deleted = true

	actorAsBytes, err := json.Marshal(actor)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("ACTOR_"+actorID, actorAsBytes)
}

// DeleteStep sets the deleted flag as true for the given step
func (s *AssetTransferSmartContract) DeleteStep(ctx contractapi.TransactionContextInterface, stepID string) error {
	step, err := s.QueryActor(ctx, "STEP_"+stepID)
	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	step.Deleted = true

	stepAsBytes, err := json.Marshal(step)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("STEP_"+stepID, stepAsBytes)
}

// DeleteAssetItem sets the deleted flag as true for the given assetItem
func (s *AssetTransferSmartContract) DeleteAssetItem(ctx contractapi.TransactionContextInterface, assetItemID string) error {
	assetItem, err := s.QueryActor(ctx, "ASSET_ITEM_"+assetItemID)
	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	assetItem.Deleted = true

	assetItemAsBytes, err := json.Marshal(assetItem)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("ASSET_ITEM_"+assetItemID, assetItemAsBytes)
}

// DeleteAsset sets the deleted flag as true for the given ssset
func (s *AssetTransferSmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, assetID string) error {
	asset, err := s.QueryActor(ctx, "ASSET_"+assetID)
	if err != nil {
		return fmt.Errorf("Failed to read the data from world state: %s", err)
	}

	asset.Deleted = true

	assetAsBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("ASSET_"+assetID, assetAsBytes)
}

// MoveAssetItem updates the owner field of assetItem with given id in world state
func (s *AssetTransferSmartContract) MoveAssetItem(ctx contractapi.TransactionContextInterface, assetItemID string, newAssetItemID string, stepID string, newOwnerID string, orderPrice string, shippingPrice string, status string, quantity string, aditionalInfo map[string]string) error {
	_, err := s.QueryAssetItem(ctx, assetItemID) // _ to oldAssetItem
	if err != nil {
		return err
	}

	newAssetItem := AssetItem{
		AssetItemID:      newAssetItemID,
		OwnerID:          newOwnerID,
		StepID:           stepID,
		ParentID:         assetItemID,
		ProcessDate:      time.Now().Format("2006-01-02 15:04:05"),
		OrderPrice:       orderPrice,
		ShippingPrice:    shippingPrice,
		Status:           status,
		Quantity:         quantity,
		Deleted:          false,
		AditionalInfoMap: aditionalInfo,
	}

	assetItemAsBytes, err := json.Marshal(newAssetItem)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState("ASSET_ITEM_"+newAssetItemID, assetItemAsBytes)
}

// QueryAssetItem returns the AssetItem stored in the world state with given id
func (s *AssetTransferSmartContract) TrackAssetItem(ctx contractapi.TransactionContextInterface, assetItemID string) ([]*AssetItem, error) {
	assetItem, err := s.QueryAssetItem(ctx, assetItemID)
	log.Print("tracking info from assetItem id: ", assetItem.AssetItemID)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if assetItem == nil {
		return nil, fmt.Errorf("%s does not exist", assetItemID)
	}

	trackedItems := make([]*AssetItem, 0)

	trackedItems = append(trackedItems, assetItem)
	for {
		currentParentId, err := strconv.Atoi(assetItem.ParentID)
		log.Print("currentParentId: ", currentParentId)
		if currentParentId <= 0 {
			log.Print("oldParentId is equals or less than 0. break it")
			break
		}
		parentAssetItem, err := s.QueryAssetItem(ctx, assetItem.ParentID)
		if err != nil {
			return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
		}
		newParentId, err := strconv.Atoi(parentAssetItem.ParentID)
		log.Print("newParentId: ", newParentId)
		trackedItems = append(trackedItems, parentAssetItem)
		assetItem = parentAssetItem
	}
	return trackedItems, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(AssetTransferSmartContract))

	if err != nil {
		fmt.Printf("Error create Arion chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting Arion chaincode: %s", err.Error())
	}
}
