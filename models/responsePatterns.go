package models

type Message struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
}

type FilmRollResponse struct {
	Success  bool      `json:"success"`
	Errors   []Message `json:"errors"`
	Messages []Message `json:"messages"`
	Result   []struct {
		Roll_id     int    `json:"rollId"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Type_id     int    `json:"typeId"`
	} `json:"result"`
}

type RollTypeResponse struct {
	Success  bool      `json:"success"`
	Errors   []Message `json:"errors"`
	Messages []Message `json:"messages"`
	Result   struct {
		Type_id   int    `json:"typeId"`
		StockName string `json:"stockName"`
		Format    string `json:"format"`
		M_id      int    `json:"mId"`
	} `json:"result"`
}
