package id_generator_service

import "github.com/devjefster/GoShortUniqueID/idgen"

const (
	timestampFormat = "2006010215"
	idCharset = "1234567890ABCDEFG"
	idLength = 5
)

func GenerateProductID() string {
	productId := idgen.New(idLength, idCharset, timestampFormat)
	return productId.Generate()
}

func GenerateOrderID() string {
	orderId := idgen.New(idLength, idCharset, timestampFormat)
	return orderId.Generate()
}