package marketplace

import "testing"

func TestReadyForHandoff(t *testing.T) {
	h := Handoff{SellerID: "s", BuyerID: "b"}
	cases := []struct {
		n string
		m []PresenceMember
		w bool
	}{{"both", []PresenceMember{{"s", "seller", true}, {"b", "buyer", true}}, true}, {"offline", []PresenceMember{{"s", "seller", true}, {"b", "buyer", false}}, false}, {"role", []PresenceMember{{"s", "buyer", true}, {"b", "buyer", true}}, false}}
	for _, c := range cases {
		t.Run(c.n, func(t *testing.T) {
			if ReadyForHandoff(c.m, h) != c.w {
				t.Fail()
			}
		})
	}
}
