package entity

type OrderStatus string

const (
	en_preparacion      OrderStatus = "en_preparacion"
	listo_para_despacho OrderStatus = "listo_para_despacho"
	en_transito         OrderStatus = "en_transito"
	entregado           OrderStatus = "entregado"
	no_entregado        OrderStatus = "no_entregado"
	devuelto            OrderStatus = "devuelto"
)
