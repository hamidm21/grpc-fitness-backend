package controller

// import (
// 	"fmt"

// 	pb "gitlab.com/mefit/mefit-api/proto"
// 	"gitlab.com/mefit/mefit-server/entity"
// 	"gitlab.com/mefit/mefit-server/services/notification"
// 	"gitlab.com/mefit/mefit-server/services/notification/mail"
// 	"gitlab.com/mefit/mefit-server/utils/config"
// 	"gitlab.com/mefit/mefit-server/utils/log"
// )

// //TODO add function for query users data
// func listItemsForListReq(usrID uint, query entity.Entity, in *pb.ListReq) (*pb.ListStatusRes, error) {
// 	log.Logger().Printf("%T", query)
// 	names, err := entity.Crud(query, usrID).Names()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &pb.ListStatusRes{
// 		Names:    names,
// 		Count:    int32(len(names)),
// 		Next:     -1,
// 		Previous: -1,
// 	}, nil
// }

// //Send mail/sms if needed
// func sendMail(u *entity.User) {
// 	confirmLink := fmt.Sprintf("https://%s/auth/confirm?token=%s",
// 		config.Config().Get("host_name"), u.ActivationToken)
// 	notif := notification.NewNotification(notification.MailType)
// 	notif.Subject = fmt.Sprintf("Confirm %s @ Mefit", u.Phone)
// 	notif.To = []string{u.Phone}
// 	notif.Body = mail.ParseTemplate(mail.WelcomeMail, struct {
// 		Name        string
// 		ConfirmLink string
// 	}{*u.Email, confirmLink})
// 	mail.Send(notif)
// }
