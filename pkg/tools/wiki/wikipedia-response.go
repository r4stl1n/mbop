package wiki

type WikipediaNormalized struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type WikipediaPages struct {
	Num4269567 struct {
		Pageid  int    `json:"pageid"`
		Ns      int    `json:"ns"`
		Title   string `json:"title"`
		Extract string `json:"extract"`
	} `json:"4269567"`
}

type WikipediaPage struct {
	Pageid  int    `json:"pageid"`
	Ns      int    `json:"ns"`
	Title   string `json:"title"`
	Extract string `json:"extract"`
}
type WikipediaQuery struct {
	Normalized []WikipediaNormalized    `json:"normalized"`
	Pages      map[string]WikipediaPage `json:"pages"`
}

type WikipediaResponse struct {
	BatchComplete string         `json:"batchcomplete"`
	Query         WikipediaQuery `json:"query"`
}
