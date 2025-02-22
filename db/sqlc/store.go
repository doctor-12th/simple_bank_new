package db

import (
	"context"
	"database/sql"
	"fmt"
)
type Store interface{
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult,error)
	Querier
}

// store provides all fucntions to execute db queries and transactions
type SQLStore struct{
	*Queries
	db *sql.DB
}

//NewStore creates a new store
func NewStore(db *sql.DB) Store{
	return &SQLStore{
		db: db,
		Queries: New(db),
	}
}

//execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error{
	tx,err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback();rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
type TransferTxParams struct{
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
	Currency string `json:"currency"`
}

type TransferTxResult struct{
	Transfer Transfers `json:"transfer"`
	FromAccount Accounts `json:"from_account"`
	ToAccount Accounts `json:"to_account"`
	FromEntry Entries `json:"from_entry"`
	ToEntry Entries `json:"to_entry"`
}

var txKey = struct{}{}
//transferTx performs a moeny transfer from on to the other
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult,error){
	var result TransferTxResult
	err := store.execTx(ctx,func(q *Queries) error{
		var err error
		txName := ctx.Value(txKey)

		fmt.Println(txName,"create tansfer")
		result.Transfer, err = q.CreateTransfer(ctx,CreateTransferParams{
			FromAccountID:arg.FromAccountID,
			ToAccountID:arg.ToAccountID,
			Amount:arg.Amount,
		})
		if err !=nil{
			return err
		}
		fmt.Println(txName,"create fromentry")
		result.FromEntry,err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID:arg.FromAccountID,
			Amount:-arg.Amount,
		})
		if err !=nil{
			return err
		}
		fmt.Println(txName,"create toentry")
		result.ToEntry,err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID:arg.ToAccountID,
			Amount:arg.Amount,
		})
		if err !=nil{
			return err
		}
		if arg.FromAccountID < arg.ToAccountID{
				//get account ->update its balance
			result.FromAccount,result.ToAccount,err = addMoney(ctx,q,arg.FromAccountID,arg.ToAccountID,-arg.Amount,arg.Amount)
			if err != nil{
				return err
			}
		}else{
			result.ToAccount,result.FromAccount,err = addMoney(ctx,q,arg.ToAccountID,arg.FromAccountID,arg.Amount,-arg.Amount)
			if err != nil{
				return err
			}
				
		}
		return nil
	})
	return result,err
}

func addMoney(ctx context.Context, q *Queries,accountID1 int64,accountID2 int64,amount1 int64,amount2 int64) (account1 Accounts,account2 Accounts,err error){
	account1,err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
		ID: accountID1,
		Amount: amount1,
	})
	if err != nil{
		return
	}
	account2,err = q.AddAccountBalance(ctx,AddAccountBalanceParams{
		ID: accountID2,
		Amount:amount2,
	})
	if err != nil{
		return
	}
	return
}