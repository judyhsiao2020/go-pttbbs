package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func Register(
	userID string,
	passwd string,
	ip string,
	email string,

	nickname string,
	realname string,
	career string,
	address string,
	over18 bool,
) (user *Userec, err error) {
	userIDRaw := &ptttype.UserID_t{}
	copy(userIDRaw[:], []byte(userID))

	passwdRaw := []byte(passwd)

	ipRaw := &ptttype.IPv4_t{}
	copy(ipRaw[:], []byte(ip))

	emailRaw := &ptttype.Email_t{}
	copy(emailRaw[:], []byte(email))

	nicknameRaw := &ptttype.Nickname_t{}
	copy(nicknameRaw[:], []byte(nickname))

	realnameRaw := &ptttype.RealName_t{}
	copy(realnameRaw[:], []byte(realname))

	careerRaw := &ptttype.Career_t{}
	copy(careerRaw[:], []byte(career))

	addressRaw := &ptttype.Address_t{}
	copy(addressRaw[:], []byte(address))

	userRaw, err := ptt.NewRegister(
		userIDRaw,
		passwdRaw,
		ipRaw,
		emailRaw,
		false,
		false,

		nicknameRaw,
		realnameRaw,
		careerRaw,
		addressRaw,
		over18,
	)
	if err != nil {
		return nil, err
	}

	user = NewUserecFromRaw(userRaw)

	return user, nil
}
