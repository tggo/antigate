package antigate

const (
	balanceBasePath = "getBalance"
)

type BalanceService struct {
	client *Client
}

type Balance struct {
	Balacne float64  `json:"balance"`
	ErrorId int64  `json:"errorId"`
}

func (s BalanceService) Get() (float64, error) {
	// build url
	path := balanceBasePath

	var balance = Balance{}
	req, err := s.client.NewRequest(path, 0, Task{})
	if err != nil {
		return 0, nil
	}
	resp, err := s.client.Do(req, &balance)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	return balance.Balacne, nil
}