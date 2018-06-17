package rtm

type TransactionsService struct {
	HTTP *HTTP
}

type Transaction struct {
	ID       string
	Undoable string
}

func (s *TransactionsService) Undo(timeline string, transaction Transaction) error {
	query := map[string]string{}
	query["timeline"] = timeline
	query["transaction_id"] = transaction.ID

	err := s.HTTP.Request("rtm.transactions.undo", query, nil)
	if err != nil {
		return err
	}

	return nil
}
