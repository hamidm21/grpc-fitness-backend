package zarinpal

import (
	"encoding/json"

	"github.com/sinabakh/go-zarinpal-checkout"
	"gitlab.com/mefit/mefit-server/utils/log"
)

func NewZarinpal() {
	zarinPay, err := zarinpal.NewZarinpal("f787e5d0-6fe2-11e9-88b7-000c29344814", false)
	if err != nil {
		log.Logger().Error(err)
	}
	zarinPayTest, err := zarinpal.NewZarinpal("f787e5d0-6fe2-11e9-88b7-000c29344814", false)
	if err != nil {
		log.Logger().Error(err)
	}
	log.Logger().Print(zarinPay)
	log.Logger().Print(zarinPayTest)
}

func NewPaymentRequest(Amount int, CallbackURL string) (string, string, int) {
	zarinPay, err := zarinpal.NewZarinpal("f787e5d0-6fe2-11e9-88b7-000c29344814", false)
	if err != nil {
		log.Logger().Error(err)
	}
	paymentURL, authority, statusCode, err := zarinPay.NewPaymentRequest(Amount, CallbackURL, "اشتراک فیتکس", "hordimad21@gmail.com", "09166599516")
	if err != nil {
		if statusCode == -3 {
			log.Logger().Print("Amount is not accepted in banking system")
		}
		log.Logger().Error(err)
	}
	// log.Println(authority)  // Save authority in DB
	// log.Println(paymentURL) // Send user to paymentURL
	return paymentURL, authority, statusCode
}

func PaymentVerification(Authority string, Price int) (bool, string, int, error) {
	zarinPay, err := zarinpal.NewZarinpal("f787e5d0-6fe2-11e9-88b7-000c29344814", false)
	if err != nil {
		log.Logger().Error(err)
	}
	log.Logger().Print(Authority)
	authority := Authority // Read authority from your storage (DB) or callback request
	amount := Price        // The amount of payment in Tomans
	verified, refID, statusCode, err := zarinPay.PaymentVerification(amount, authority)
	if err != nil {
		if statusCode == 101 {
			log.Logger().Print("Payment is already verified")
		}
		log.Logger().Error("error:", err, "statusCode:", statusCode)
	}
	log.Logger().Print("verified:", verified, "refID:", refID)
	return verified, refID, statusCode, nil
}

func UnverifiedTransactions() {
	zarinPay, err := zarinpal.NewZarinpal("f787e5d0-6fe2-11e9-88b7-000c29344814", false)
	if err != nil {
		log.Logger().Error(err)
	}
	authorities, statusCode, err := zarinPay.UnverifiedTransactions()
	if err != nil {
		log.Logger().Error("statusCode:", statusCode, "error:", err)
	}
	marshaledJSON, _ := json.Marshal(authorities)
	log.Logger().Print(string(marshaledJSON))
	// Output:
	// [{"Authority":"XXXX","Amount":100,"Channel":"WebGate","CallbackURL":"http://localhost:3000","Referer":"/","Email":"","CellPhone":"","Date":"2017-12-27 22:12:59"}]
}

func RefreshAuthority() {
	zarinPay, err := zarinpal.NewZarinpal("f787e5d0-6fe2-11e9-88b7-000c29344814", false)
	if err != nil {
		log.Logger().Error(err)
	}
	statusCode, err := zarinPay.RefreshAuthority("XXXX", 2000)
	if err != nil {
		log.Logger().Error("statusCode:", statusCode, "error:", err)
	}
}
