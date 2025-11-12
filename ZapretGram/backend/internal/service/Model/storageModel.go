package model

type Server struct {
	Address   string `json:"Address"`
	PublicKey string `json:"PublicKey"`
}

type StorageServer struct {
	Servers map[string]Server `json:"servers"`
}
