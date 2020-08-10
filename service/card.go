package service

type CardInfoPayload struct {
	Balance float32 `json:"balance"`
	Status  string  `json:"status"`
}

// 获取校园卡信息，包括余额、校园卡在用状态
func GetCardInfo(sid, password string) (*CardInfoPayload, error) {

	data, err := MakeCardInfoRequest(sid, password)
	if err != nil {
		return nil, err
	}

	payload := &CardInfoPayload{
		Balance: data.Balance,
		Status:  data.StatusDesc,
	}

	return payload, nil
}

type DealRecord struct{}

// 获取校园卡消费流水
func GetConsumeList(sid, password, limit, page, start, end string) ([]*DealRow, error) {

	records, err := MakeConsumesRequest(sid, password, limit, page, start, end)
	if err != nil {
		return nil, err
	}

	return records, nil
}
