package rtm

type TransactionsService struct {
	HTTP *HTTP
}

type Transaction struct {
	ID       string
	Undoable string
}

type transactionOperationResultContent struct {
	Stat string
	Err  ErrorResponse
}

type transactionOperationResult struct {
	Rsp transactionOperationResultContent
}

func (s *TransactionsService) Undo(timeline string, transaction Transaction) error {
	result := new(transactionOperationResult)

	query := map[string]string{}
	query["timeline"] = timeline
	query["transaction_id"] = transaction.ID

	err := s.HTTP.Request("rtm.transactions.undo", query, &result)
	err = s.HTTP.VerifyResponse(err, result.Rsp.Stat, result.Rsp.Err)
	if err != nil {
		return err
	}

	return nil
}
