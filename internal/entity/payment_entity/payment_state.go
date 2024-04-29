package payment_entity

type PaymentState int

const (
	None               PaymentState = iota
	WaitingForApproval              // When the payment request is sent to the payment gateway
	Approved                        // When the payment is approved by the payment gateway
	Rejected                        // When the payment is rejected by the payment gateway
)

func (s PaymentState) String() string {
	switch s {
	case None:
		return "None"
	case WaitingForApproval:
		return "WaitingForApproval"
	case Approved:
		return "Approved"
	case Rejected:
		return "Rejected"
	default:
		return "Unknown"
	}
}
