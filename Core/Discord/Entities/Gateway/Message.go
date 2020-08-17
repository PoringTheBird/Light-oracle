package Gateway

type Message struct {
	Op 	int 		`json:"op"`
	D 	interface{} `json:"d"`
	S 	*int 		`json:"s"`
	T 	*string		`json:"t"`
}