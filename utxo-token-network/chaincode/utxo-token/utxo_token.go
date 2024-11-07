package main

import (
    "fmt"
    "strconv"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// UTXOTokenContract 구조체 정의
type UTXOTokenContract struct {
    contractapi.Contract
}

// 토큰 발행 메서드
func (u *UTXOTokenContract) IssueToken(ctx contractapi.TransactionContextInterface, recipient string, amount int) error {
    utxoKey := "utxo_" + recipient + "_" + strconv.Itoa(amount)
    err := ctx.GetStub().PutState(utxoKey, []byte(strconv.Itoa(amount)))
    if err != nil {
        return fmt.Errorf("토큰 발행 실패: %s", err.Error())
    }
    return nil
}

// 토큰 전송 메서드
func (u *UTXOTokenContract) TransferToken(ctx contractapi.TransactionContextInterface, sender string, recipient string, amount int) error {
    senderUTXO := "utxo_" + sender
    senderBalanceAsBytes, err := ctx.GetStub().GetState(senderUTXO)
    if err != nil || senderBalanceAsBytes == nil {
        return fmt.Errorf("잔액 조회 실패")
    }

    senderBalance, _ := strconv.Atoi(string(senderBalanceAsBytes))
    if senderBalance < amount {
        return fmt.Errorf("잔액 부족")
    }

    newSenderBalance := senderBalance - amount
    err = ctx.GetStub().PutState(senderUTXO, []byte(strconv.Itoa(newSenderBalance)))
    if err != nil {
        return fmt.Errorf("전송 실패: %s", err.Error())
    }

    recipientUTXO := "utxo_" + recipient
    recipientBalanceAsBytes, _ := ctx.GetStub().GetState(recipientUTXO)
    recipientBalance, _ := strconv.Atoi(string(recipientBalanceAsBytes))

    newRecipientBalance := recipientBalance + amount
    err = ctx.GetStub().PutState(recipientUTXO, []byte(strconv.Itoa(newRecipientBalance)))
    if err != nil {
        return fmt.Errorf("수신자 잔액 업데이트 실패: %s", err.Error())
    }

    return nil
}

// 잔액 조회 메서드
func (u *UTXOTokenContract) BalanceOf(ctx contractapi.TransactionContextInterface, owner string) (int, error) {
    utxoKey := "utxo_" + owner
    balanceAsBytes, err := ctx.GetStub().GetState(utxoKey)
    if err != nil || balanceAsBytes == nil {
        return 0, fmt.Errorf("잔액 조회 실패")
    }
    balance, _ := strconv.Atoi(string(balanceAsBytes))
    return balance, nil
}

// 메인 함수
func main() {
    chaincode, err := contractapi.NewChaincode(new(UTXOTokenContract))
    if err != nil {
        fmt.Printf("체인코드 생성 오류: %s", err.Error())
        return
    }
    if err := chaincode.Start(); err != nil {
        fmt.Printf("체인코드 시작 오류: %s", err.Error())
    }
}