package duckduckgo

// DuckDuckGoResult represents a single search result from DuckDuckGo
type DuckDuckGoResult struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

// DuckDuckGoResponse represents the response from a DuckDuckGo search
type DuckDuckGoResponse struct {
	Results []DuckDuckGoResult `json:"results"`
}

// DuckDuckGoJsonResponse represents the JSON response from the DuckDuckGo API
type DuckDuckGoJsonResponse struct {
	Abstract         string      `json:"Abstract"`
	AbstractSource   string      `json:"AbstractSource"`
	AbstractText     string      `json:"AbstractText"`
	AbstractURL      string      `json:"AbstractURL"`
	Answer           string      `json:"Answer"`
	AnswerType       string      `json:"AnswerType"`
	Definition       string      `json:"Definition"`
	DefinitionSource string      `json:"DefinitionSource"`
	DefinitionURL    string      `json:"DefinitionURL"`
	Entity           string      `json:"Entity"`
	Heading          string      `json:"Heading"`
	Image            string      `json:"Image"`
	ImageHeight      interface{} `json:"ImageHeight"` // Can be string or int
	ImageIsLogo      interface{} `json:"ImageIsLogo"` // Can be string or int
	ImageWidth       interface{} `json:"ImageWidth"`  // Can be string or int
	Infobox          interface{} `json:"Infobox"`     // Can be string or object
	Redirect         string      `json:"Redirect"`
	RelatedTopics    []struct {
		FirstURL string `json:"FirstURL,omitempty"`
		Icon     struct {
			Height string `json:"Height"`
			URL    string `json:"URL"`
			Width  string `json:"Width"`
		} `json:"Icon,omitempty"`
		Result string `json:"Result,omitempty"`
		Text   string `json:"Text,omitempty"`
		Name   string `json:"Name,omitempty"`
		Topics []struct {
			FirstURL string `json:"FirstURL"`
			Icon     struct {
				Height string `json:"Height"`
				URL    string `json:"URL"`
				Width  string `json:"Width"`
			} `json:"Icon"`
			Result string `json:"Result"`
			Text   string `json:"Text"`
		} `json:"Topics,omitempty"`
	} `json:"RelatedTopics"`
	Results []interface{} `json:"Results"`
	Type    string        `json:"Type"`
	Meta    struct {
		Attribution  interface{} `json:"attribution"`
		Blockgroup   interface{} `json:"blockgroup"`
		CreatedDate  interface{} `json:"created_date"`
		Description  string      `json:"description"`
		Designer     interface{} `json:"designer"`
		DevDate      interface{} `json:"dev_date"`
		DevMilestone string      `json:"dev_milestone"`
		Developer    []struct {
			Name string `json:"name"`
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"developer"`
		ExampleQuery    string      `json:"example_query"`
		ID              string      `json:"id"`
		IsStackexchange interface{} `json:"is_stackexchange"`
		JsCallbackName  string      `json:"js_callback_name"`
		LiveDate        interface{} `json:"live_date"`
		Maintainer      struct {
			Github string `json:"github"`
		} `json:"maintainer"`
		Name            string      `json:"name"`
		PerlModule      string      `json:"perl_module"`
		Producer        interface{} `json:"producer"`
		ProductionState string      `json:"production_state"`
		Repo            string      `json:"repo"`
		SignalFrom      string      `json:"signal_from"`
		SrcDomain       string      `json:"src_domain"`
		SrcID           int         `json:"src_id"`
		SrcName         string      `json:"src_name"`
		SrcOptions      struct {
			Directory         string `json:"directory"`
			IsFanon           int    `json:"is_fanon"`
			IsMediawiki       int    `json:"is_mediawiki"`
			IsWikipedia       int    `json:"is_wikipedia"`
			Language          string `json:"language"`
			MinAbstractLength string `json:"min_abstract_length"`
			SkipAbstract      int    `json:"skip_abstract"`
			SkipAbstractParen int    `json:"skip_abstract_paren"`
			SkipEnd           string `json:"skip_end"`
			SkipIcon          int    `json:"skip_icon"`
			SkipImageName     int    `json:"skip_image_name"`
			SkipQr            string `json:"skip_qr"`
			SourceSkip        string `json:"source_skip"`
			SrcInfo           string `json:"src_info"`
		} `json:"src_options"`
		SrcURL interface{} `json:"src_url"`
		Status string      `json:"status"`
		Tab    string      `json:"tab"`
		Topic  []string    `json:"topic"`
		Unsafe int         `json:"unsafe"`
	} `json:"meta"`
}
