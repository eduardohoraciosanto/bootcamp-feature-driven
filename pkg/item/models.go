package item

type Item struct {
	ID       string
	Name     string
	Quantity int
	Price    float32
}

type TransportItem struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity,omitempty"`
	Price    float32 `json:"price"`
}

type ExternalItem struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Price string `json:"price,omitempty"`
}

type ExternalHealth struct {
	Status string `json:"status,omitempty"`
}

type ExternalMeta struct {
	Version string `json:"version"`
}
type ExternalHealthResponse struct {
	Meta ExternalMeta   `json:"meta"`
	Data ExternalHealth `json:"data,omitempty"`
}

type ExternalGetItemResponse struct {
	Meta ExternalMeta `json:"meta"`
	Data ExternalItem `json:"data,omitempty"`
}

type ExternalGetAllItemsResponse struct {
	Meta ExternalMeta   `json:"meta"`
	Data []ExternalItem `json:"data,omitempty"`
}
