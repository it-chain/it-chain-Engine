package service

import (
	"errors"
	"fmt"

	"encoding/json"

	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/it-chain-Engine/icode"
	"github.com/it-chain/midgard"
	"github.com/it-chain/tesseract"
	"github.com/it-chain/tesseract/cellcode/cell"
)

type TesseractContainerService struct {
	tesseract      *tesseract.Tesseract
	repository     icode.ReadOnlyMetaRepository
	containerIdMap map[icode.ID]string // key : iCodeId, value : containerId
}

func NewTesseractContainerService(config tesseract.Config, repository icode.ReadOnlyMetaRepository) *TesseractContainerService {
	tesseractObj := &TesseractContainerService{
		tesseract:      tesseract.New(config),
		repository:     repository,
		containerIdMap: make(map[icode.ID]string, 0),
	}
	tesseractObj.InitContainers()
	return tesseractObj
}

func (cs TesseractContainerService) StartContainer(meta icode.Meta) error {
	tesseractIcodeInfo := tesseract.ICodeInfo{
		Name:      meta.RepositoryName,
		Directory: meta.Path,
	}
	containerId, err := cs.tesseract.SetupContainer(tesseractIcodeInfo)
	if err != nil {
		return err
	}
	cs.containerIdMap[meta.ICodeID] = containerId
	return nil
}

func (cs TesseractContainerService) ExecuteTransaction(tx icode.Transaction) (*icode.Result, error) {
	containerId, found := cs.containerIdMap[tx.TxData.ICodeID]

	if !found {
		return nil, errors.New(fmt.Sprintf("no container for iCode : %s", tx.TxData.ICodeID))
	}

	tesseractTxInfo := cell.TxInfo{
		Method: tx.TxData.Method,
		ID:     tx.TxData.ID,
		Params: cell.Params{
			Function: tx.TxData.Params.Function,
			Args:     tx.TxData.Params.Args,
		},
	}

	res, err := cs.tesseract.QueryOrInvoke(containerId, tesseractTxInfo)

	if err != nil {
		return nil, err
	}
	var data map[string]string
	var isSuccess bool

	switch res.Result {
	case "Success":
		isSuccess = true
		err = json.Unmarshal(res.Data, data)
		if err != nil {
			return nil, err
		}
	case "Error":
		isSuccess = false
		data = nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown pb response result %s", res.Result))
	}

	result := &icode.Result{
		Data:    data,
		TxId:    tx.TxId,
		Success: isSuccess,
	}
	return result, nil
}

func (cs TesseractContainerService) StopContainer(id icode.ID) error {
	cs.tesseract.Clients[cs.containerIdMap[id]].Close()
	delete(cs.containerIdMap, id)
	deletedEvent := icode.MetaDeletedEvent{
		EventModel: midgard.EventModel{
			ID:   id,
			Type: "meta.deleted",
		},
	}
	return eventstore.Save(id, deletedEvent)
}

// start containers in repos
func (cs *TesseractContainerService) InitContainers() error {
	metas, err := cs.repository.FindAll()
	if err != nil {
		return err
	}
	for _, meta := range metas {
		cs.StartContainer(*meta)
	}
	return nil
}
