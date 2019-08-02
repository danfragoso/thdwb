package structs

//NodeDOM "DOM Node Struct definition"
type NodeDOM struct {
	Element    string       `json:"element"`
	Content    string       `json:"content"`
	Children   []*NodeDOM   `json:"children"`
	Attributes []*Attribute `json:"attributes"`
	Style      *Stylesheet  `json:"style"`
	Parent     *NodeDOM
}

//Resource "HTTP resource struct definition"
type Resource struct {
	Body string
	Code int
}

//Attribute "Generic key:value attribute definition"
type Attribute struct {
	Name  string
	Value string
}

//Stylesheet "Stylesheet definition for DOM Nodes"
type Stylesheet struct {
	Color    string
	FontSize int
}
