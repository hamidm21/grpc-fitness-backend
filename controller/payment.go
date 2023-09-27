package controller

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	// "fmt"
	"github.com/jinzhu/gorm"
	pb "gitlab.com/mefit/mefit-api/proto"
	"gitlab.com/mefit/mefit-server/entity"
	"gitlab.com/mefit/mefit-server/services/zarinpal"
	"gitlab.com/mefit/mefit-server/utils"
	"gitlab.com/mefit/mefit-server/utils/config"
	"gitlab.com/mefit/mefit-server/utils/log"
)

func (s *Controller) PaymentRequest(ctx context.Context, in *pb.PayReq) (*pb.PayRes, error) {
	usrID := uint(ctx.Value(utils.KeyEmail).(uint))
	product := entity.Product{}
	product.ID = uint(in.ProductID)
	if err := entity.SimpleCrud(product).Get(&product, "Plan"); err != nil {
		return nil, utils.ErrNotFound
	}
	pro := &entity.Profile{UserID: ctx.Value(utils.KeyEmail).(uint)}
	if err := entity.SimpleCrud(*pro).Get(pro); err != nil {
		return nil, err
	}
	Purchased := entity.PurchasedProduct{}
	Purchased.ProductID = product.ID
	Purchased.ProfileID = pro.ID
	if err := entity.SimpleCrud(Purchased).Get(&Purchased); err != nil {
		if err == utils.ErrNotFound {
			if in.IsBazaar == true {
				Bpay := entity.BazaarPayment{}
				Bpay.ProductID = product.ID
				Bpay.ProductName = product.Name
				Bpay.UserID = usrID
				Bpay.DevPayload = strconv.FormatInt(time.Now().Unix(), 10) + "/" + strconv.FormatUint(uint64(ctx.Value(utils.KeyEmail).(uint)), 10)
				Bpay.Price = product.Price
				Bpay.PlanID = product.Plan.ID
				Bpay.RSA = "MIHNMA0GCSqGSIb3DQEBAQUAA4G7ADCBtwKBrwC2iCaroSbzqC2LrDH9S51PYlUzRyMHNXts3/QOexCN4x4spF5cSz4cj9NMzFTP0EKR4fRicBIBzrHZfoKn4QTeFrC1wP1GyWPRRzcf0eVVSsTV7ozgounOQNyrqV+T+5nJnkUDZha4tJD7jP8ejC77ePvr7EQmqmnmFNEEs7WtNVSVcL6nDqMig+bI/QUDjDQQ+QrUmw4eVIYxqPr+NOzqg213G9hPxqLz9K9OND8CAwEAAQ=="
				Bpay.SKU = product.SKU
				if err := entity.SimpleCrud(entity.BazaarPayment{ProductID: product.ID, UserID: usrID}).FirstOrCreate(&Bpay); err != nil {
					return nil, err
				}
				log.Logger().Print("SKU is ....  ", product.SKU, " Sku is ... ", Bpay.SKU)
				return &pb.PayRes{
					Type: &pb.PayRes_Bazaar{
						Bazaar: &pb.BazaarPayload{
							RSA:        Bpay.RSA,
							SKU:        product.SKU,
							DevPayload: Bpay.DevPayload,
						},
					},
				}, nil
			}
			price := product.Price - uint(float32(product.Off*product.Price)/100.0)
			paymentURl, authority, status := zarinpal.NewPaymentRequest(int(price), config.Config().GetString(utils.CallbackURL))

			payment := entity.Payment{
				ProductID:   uint(in.ProductID),
				ProductName: product.Name,
				UserID:      uint(ctx.Value(utils.KeyEmail).(uint)),
				Price:       price,
				Plan:        product.Plan,
				PaymentURL:  paymentURl,
				Authority:   authority,
				Status:      uint(status),
			}

			if err := entity.SimpleCrud(&payment).Create(); err != nil {
				return nil, err
			}
			return &pb.PayRes{
				Type: &pb.PayRes_Zarin{
					Zarin: &pb.ZarinPayload{
						TransactionId: int32(payment.ID),
						Uri:           paymentURl,
					},
				},
			}, nil
		}
		return nil, err
	}
	return &pb.PayRes{
		Type: &pb.PayRes_Paid{
			Paid: true,
		},
	}, nil
}

//RESP is a type struct for cafe bazaar verification response
type RESP struct {
	consumptionState int
	purchaseState    int
	kind             string
	developerPayload string
	purchaseTime     time.Time
}

func (s *Controller) BazaarPaymentCheck(ctx context.Context, in *pb.BazaarReq) (*pb.BazaarRes, error) {

	product := entity.Product{}
	product.SKU = in.ProductID
	if err := entity.SimpleCrud(product).Get(&product); err != nil {
		return nil, err
	}

	payment := entity.BazaarPayment{
		Token: in.PurchaseID,
	}
	if err := entity.SimpleCrud(payment).Get(&payment); err != nil && err != utils.ErrNotFound {
		return nil, err
	}
	if payment.Paid == true {
		return &pb.BazaarRes{
			Paid: true,
		}, nil
	}
	bazaarURL := "https://pardakht.cafebazaar.ir/devapi/v2/api/validate/" + in.PackageName + "/inapp/" + in.ProductID + "/purchases/" + in.PurchaseID + "/"
	client := &http.Client{}
	request, err := http.NewRequest("GET", bazaarURL, nil)
	if err != nil {
		return nil, err
	}
	accessToken, err := GetNewAccessToken()
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", accessToken)
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response RESP

	if resp.StatusCode == http.StatusOK {
		log.Logger().Info("respnse of bazaar payment is ok")
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Logger().Error(err)
			return nil, err
		}
		bodyString := string(bodyBytes)
		json.Unmarshal([]byte(bodyString), &response)
		if response.purchaseState == 0 {
			err := DoInTransaction(func(tx *gorm.DB) error {
				usrID := ctx.Value(utils.KeyEmail).(uint)
				p := entity.BazaarPayment{ProductID: product.ID, UserID: usrID}
				if err := entity.SimpleCrud(p).Get(&p); err != nil {
					return err
				}
				paymentUpdate := entity.BazaarPayment{
					PaymentURL: bazaarURL,
					Token:      in.PurchaseID,
					Paid:       true,
					Status:     uint(response.purchaseState),
					ProductID:  product.ID,
					UserID:     usrID,
				}
				q := entity.BazaarPayment{}
				q.ID = p.ID
				if err := entity.SimpleCrud(q).WithTransaction(tx).Updates(&paymentUpdate); err != nil {
					return utils.ErrInternal
				}

				//FIXME: fix this gorm related bullshiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiit
				prof := entity.Profile{}
				prof.UserID = usrID
				if err := entity.SimpleCrud(prof).Get(&prof); err != nil {
					return err
				}
				purchased := entity.PurchasedProduct{
					PaymentType: "bazaar",
					PaymentID:   p.ID,
					ProfileID:   prof.ID,
					ProductID:   product.ID,
				}
				if err := tx.Where(entity.PurchasedProduct{ProfileID: prof.ID, ProductID: product.ID}).FirstOrCreate(&purchased).Error; err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return nil, err
			}

			log.Logger().Info("purchase is paid")
			return &pb.BazaarRes{
				Paid: true,
			}, nil
		}
		log.Logger().Info("e 1")
		return &pb.BazaarRes{
			Paid: false,
		}, nil
	}
	log.Logger().Info("e 2")
	return &pb.BazaarRes{
		Paid: false,
	}, nil
}

//PaymentVerification indicates if a payment is successful or not
func PaymentVerification(authority string, status string) (string, string, uint, string, bool) {
	if status == "OK" {
		payment := entity.Payment{}
		payment.Authority = authority
		if err1 := entity.SimpleCrud(payment).Get(&payment); err1 != nil || payment.Paid {
			log.Logger().Error(err1)
			return "پرداخت نامعتبر است یا قبلا تکمیل شده است", "", 0, "", false
		}
		refrence := strings.TrimLeft(authority, "0")
		verified, refID, statusCode, err := zarinpal.PaymentVerification(payment.Authority, int(payment.Price))
		if err != nil {
			return "پرداخت توسط بانک تایید نشد", "", 0, "", false
		}
		if verified {
			payment.Paid = true
			payment.Refrence = refID
			payment.Status = uint(statusCode)
			query := entity.Payment{Authority: authority}
			if err := entity.SimpleCrud(query).Updates(payment); err != nil {
				log.Logger().Error(err)
				return "مشکلات درونی اتفاق افتاده است لطفا شکیبا باشید", payment.ProductName, payment.Price, refrence, false
			}
			prof := entity.Profile{
				UserID: payment.UserID,
			}
			if err := entity.SimpleCrud(prof).Get(&prof); err != nil {
				log.Logger().Error(err)
				return "مشکل در پیدا کردن اطلاعات کاربری", payment.ProductName, payment.Price, refrence, false
			}

			return "پرداخت با موفقیت انجام شد", payment.ProductName, payment.Price, refrence, true
		}
		return "پرداخت نامعتبر", payment.ProductName, payment.Price, refrence, false

	}
	return "پرداخت نامعتبر", "", 0, "", false
}

func GetNewAccessToken() (string, error) {
	bazaarURL := "https://pardakht.cafebazaar.ir/devapi/v2/auth/token/"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", "FsmTKs30YS1ijDeeMIOv6umyzQKvAJmSK12LLyKO")
	data.Set("client_secret", "Fln0CWGzupMEKwN64tAvULPTDhjl60YBzWmxNIPxZKi0o0DVWOoFgNkDV0pg")
	data.Set("refresh_token", "GyO3hJKtAwEMWyeZmqNORiqSaE9NWm")

	client := &http.Client{}
	r, _ := http.NewRequest("POST", bazaarURL, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	type RESP struct {
		Token string `json:"access_token"`
	}

	var response RESP
	resp, _ := client.Do(r)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Logger().Error(err)
		}
		bodyString := string(bodyBytes)
		json.Unmarshal([]byte(bodyString), &response)
		return response.Token, nil
	}

	return "", utils.ErrInternal
}

//gorm related Bullshit

type InTransaction func(tx *gorm.DB) error

func DoInTransaction(fn InTransaction) error {
	tx := entity.GetDB().Begin()
	if tx.Error != nil {
		return tx.Error
	}
	err := fn(tx)
	if err != nil {
		xerr := tx.Rollback().Error
		if xerr != nil {
			log.Logger().Errorf("While rollback build: %v", xerr)
			return utils.ErrUnavailable
		}
		return err
	}
	if err = tx.Commit().Error; err != nil {
		log.Logger().Errorf("While commit build: %v", err)
		return utils.ErrUnavailable
	}
	return nil
}
