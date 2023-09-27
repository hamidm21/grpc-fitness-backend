package main

import (
	"fmt"
	"html/template"
	"net/http"
	"gitlab.com/mefit/mefit-server/controller"
	"gitlab.com/mefit/mefit-server/utils"
	"gitlab.com/mefit/mefit-server/utils/config"
	"gitlab.com/mefit/mefit-server/utils/initializer"
	"gitlab.com/mefit/mefit-server/utils/log"
)

var callbackTmpl *template.Template

func main() {
	// Initalize
	defer initializer.Initialize()()
	http.HandleFunc("/payment", callbackFunc)
	fs := http.FileServer(http.Dir("/static"))
	http.Handle("/payment/static/", http.StripPrefix("/payment/static", fs))
	// Mount admin interface to mux
	port := config.Config().GetString(utils.KeyPort)
	log.Logger().Infof("Listening on: %s", port)
	panic(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func callbackFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET params were:", r.URL.Query())
	// if only one expected
	Authority := r.URL.Query().Get("Authority")
	Status := r.URL.Query().Get("Status")
	log.Logger().Info(Authority)
	title, productName, price, refrence, paymentSuccedd := controller.PaymentVerification(Authority, Status)
	logo := "success.png"
	msg := title
	if !paymentSuccedd {
		logo = "failure.png"
	}
	callbackTmpl := template.Must(template.ParseFiles("templates/payment.html"))
	// if multiples possible, or to process empty values like param1 in
	// ?param1=&param2=something
	// param1s := r.URL.Query()["param1"]
	// if len(param1s) > 0 {
	// ... process them ... or you could just iterate over them without a check
	// this way you can also tell if they passed in the parameter as the empty string
	// it will be an element of the array that is the empty string
	// }
	type pay struct {
		Title   string
		Logo    string
		Ref     string
		Price   uint
		Type    string
		Message string
	}
	data := pay{
		Title:   title,
		Logo:    logo,
		Ref:     refrence,
		Price:   price,
		Type:    productName,
		Message: msg,
	}
	callbackTmpl.Execute(w, data)

}