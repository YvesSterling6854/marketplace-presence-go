package marketplace

type PresenceMember struct {
	ID     string `json:"id"`
	Role   string `json:"role"`
	Online bool   `json:"online"`
}
type Handoff struct{ OrderID, SellerID, BuyerID string }

func ReadyForHandoff(ms []PresenceMember, h Handoff) bool {
	s, b := false, false
	for _, m := range ms {
		if m.Online && m.ID == h.SellerID && m.Role == "seller" {
			s = true
		}
		if m.Online && m.ID == h.BuyerID && m.Role == "buyer" {
			b = true
		}
	}
	return s && b
}
