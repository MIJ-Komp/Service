package enum

type EOrderStatus string

const (
	OrderStatusPending   EOrderStatus = "pending"
	OrderStatusPaid      EOrderStatus = "paid"
	OrderStatusPreparing EOrderStatus = "preparing"
	OrderStatusShipped   EOrderStatus = "shipped"
	OrderStatusArrived   EOrderStatus = "arrived"
	OrderStatusCompleted EOrderStatus = "completed"
	OrderStatusCancelled EOrderStatus = "cancelled"
)

func (e EOrderStatus) DisplayString() string {
	switch e {
	case OrderStatusPending:
		return "Pending"
	case OrderStatusPaid:
		return "Dibayar"
	case OrderStatusPreparing:
		return "Di Proses"
	case OrderStatusShipped:
		return "Dalam Pengiriman"
	case OrderStatusArrived:
		return "Tiba Di Tujuan"
	case OrderStatusCompleted:
		return "Selesai"
	default:
		return "Dibatalkan"
	}
}

func (s EOrderStatus) IsValid() bool {
	switch s {
	case OrderStatusPending,
		OrderStatusPaid,
		OrderStatusPreparing,
		OrderStatusShipped,
		OrderStatusArrived,
		OrderStatusCompleted,
		OrderStatusCancelled:
		return true
	}
	return false
}
