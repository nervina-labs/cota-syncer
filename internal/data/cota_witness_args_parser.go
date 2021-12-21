package data

import (
	"context"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/biz"
	"github.com/nervina-labs/cota-nft-entries-syncer/internal/data/blockchain"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
)

type CotaWitnessArgsParser struct {
	client *CkbNodeClient
}

func NewCotaWitnessArgsParser(client *CkbNodeClient) CotaWitnessArgsParser {
	return CotaWitnessArgsParser{
		client: client,
	}
}

type cotaCell struct {
	output *ckbTypes.CellOutput
	index  int
}

func (c CotaWitnessArgsParser) Parse(tx *ckbTypes.Transaction, cotaType SystemScript) ([]biz.Entry, error) {
	if !c.hasCotaCell(tx.Outputs, cotaType) {
		return nil, nil
	}
	return c.cotaEntries(tx, cotaType)
}

func (c CotaWitnessArgsParser) isCotaCell(output *ckbTypes.CellOutput, cotaType SystemScript) bool {
	if output.Type == nil {
		return false
	}
	return output.Type.CodeHash == cotaType.CodeHash && output.Type.HashType == cotaType.HashType
}

// inputs 中 cota cells 的个数一定与 outputs 中 cota cells 的个数相等
// 批量注册多个 cota cell 的时候 input 里可能没有 cota cell
func (c CotaWitnessArgsParser) cotaEntries(tx *ckbTypes.Transaction, cotaType SystemScript) ([]biz.Entry, error) {
	inputCotaCellGroups, err := c.inputCotaCellGroups(tx.Inputs, cotaType)
	if err != nil {
		return nil, err
	}
	outputCotaCellGroups, err := c.outputCotaCellGroups(tx.Outputs, cotaType)
	if err != nil {
		return nil, err
	}
	cotaCells := make([]cotaCell, len(inputCotaCellGroups))
	var index int
	for typeHash := range inputCotaCellGroups {
		cotaCell := outputCotaCellGroups[typeHash]
		cotaCells[index] = cotaCell
		index += 1
	}

	var entries []biz.Entry
	for _, cotaCell := range cotaCells {
		witness := tx.Witnesses[cotaCell.index]
		bytes, err := blockchain.WitnessArgsFromSliceUnchecked(witness).InputType().IntoBytes()
		if err != nil {
			return nil, err
		}
		entries = append(entries, biz.Entry{
			Witness:    bytes.RawData(),
			LockScript: cotaCell.output.Lock,
		})
	}
	return entries, nil
}

func (c CotaWitnessArgsParser) inputCotaCellGroups(inputs []*ckbTypes.CellInput, cotaType SystemScript) (map[string]cotaCell, error) {
	group := make(map[string]cotaCell)
	cotaCells, err := c.inputCotaCells(inputs, cotaType)
	if err != nil {
		return group, err
	}
	for _, cell := range cotaCells {
		typeHash, err := cell.output.Type.Hash()
		if err != nil {
			return group, err
		}
		if _, ok := group[typeHash.String()]; !ok {
			group[typeHash.String()] = cell
		}
	}

	return group, nil
}

func (c CotaWitnessArgsParser) hasCotaCell(outputs []*ckbTypes.CellOutput, cotaType SystemScript) (result bool) {
	for _, output := range outputs {
		if result = c.isCotaCell(output, cotaType); result {
			break
		}
	}
	return result
}

func (c CotaWitnessArgsParser) inputCotaCells(inputs []*ckbTypes.CellInput, cotaType SystemScript) ([]cotaCell, error) {
	var cotaCells []cotaCell
	for i := 0; i < len(inputs); i++ {
		prevOutpoint := inputs[i].PreviousOutput
		prevTx, err := c.client.Rpc.GetTransaction(context.TODO(), prevOutpoint.TxHash)
		if err != nil {
			return nil, err
		}
		prevCellOutput := prevTx.Transaction.Outputs[prevOutpoint.Index]
		if c.isCotaCell(prevCellOutput, cotaType) {
			cotaCells = append(cotaCells, cotaCell{
				output: prevCellOutput,
				index:  i,
			})
		}
	}
	return cotaCells, nil
}

func (c CotaWitnessArgsParser) outputCotaCells(outputs []*ckbTypes.CellOutput, cotaType SystemScript) ([]cotaCell, error) {
	var cotaCells []cotaCell
	for i := 0; i < len(outputs); i++ {
		if c.isCotaCell(outputs[i], cotaType) {
			cotaCells = append(cotaCells, cotaCell{
				output: outputs[i],
				index:  i,
			})
		}
	}
	return cotaCells, nil
}

func (c CotaWitnessArgsParser) outputCotaCellGroups(outputs []*ckbTypes.CellOutput, cotaType SystemScript) (map[string]cotaCell, error) {
	group := make(map[string]cotaCell)
	cotaCells, err := c.outputCotaCells(outputs, cotaType)
	if err != nil {
		return group, err
	}
	for _, cell := range cotaCells {
		typeHash, err := cell.output.Type.Hash()
		if err != nil {
			return group, err
		}
		if _, ok := group[typeHash.String()]; !ok {
			group[typeHash.String()] = cell
		}
	}

	return group, nil
}
