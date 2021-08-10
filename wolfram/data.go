package wolfram

import "encoding/xml"

type QueryResult struct {
	XMLName    xml.Name `xml:"queryresult"`
	HasSuccess bool     `xml:"success,attr"`
	HasError   bool     `xml:"error,attr"`

	Pods        []Pod       `xml:"pod"`
	Tips        Tips        `xml:"tips"`
	DidYouMeans DidYouMeans `xml:"didyoumeans"`
	Assumptions Assumptions `xml:"assumptions"`
	Error       Error       `xml:"error"`
}

type Pod struct {
	XMLName xml.Name `xml:"pod"`
	Title   string   `xml:"title,attr"`
	Subpods []Subpod `xml:"subpod"`
}

type Subpod struct {
	Image Image  `xml:"img"`
	Text  string `xml:"plaintext"`
}

type Image struct {
	XMLName xml.Name `xml:"img"`
	URL     string   `xml:"src,attr"`
}

type Tips struct {
	XMLName xml.Name `xml:"tips"`
	Data    []Tip    `xml:"tip"`
}

type Tip struct {
	XMLName xml.Name `xml:"tip"`
	Text    string   `xml:"text,attr"`
}

type DidYouMeans struct {
	XMLName xml.Name `xml:"didyoumeans"`
	Data    []string `xml:"didyoumean"`
}

type Assumptions struct {
	XMLName xml.Name     `xml:"assumptions"`
	Data    []Assumption `xml:"assumption"`
}

type Assumption struct {
	XMLName xml.Name          `xml:"assumption"`
	Word    string            `xml:"word"`
	Values  []AssumptionValue `xml:"value"`
}

type AssumptionValue struct {
	XMLName xml.Name `xml:"value"`
	Content string   `xml:"desc"`
	Input   string   `xml:"input,attr"`
}

type Error struct {
	XMLName xml.Name `xml:"error"`
	Content string   `xml:"msg"`
}
