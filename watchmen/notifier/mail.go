package notifier

import (
    "log"
    "net/smtp"
    "strconv"
)


type Mail struct {
    User string
    Password string
    Host string
    Port int
    From string
    Recipients []string
    Subject string
}


func NewMail(User, Password, Host string, Port int,
             From string, Recipients []string) *Mail {

    mail := &Mail {
        User: User,
        Password: Password,
        Host: Host,
        Port: Port,
        From: From,
        Recipients: Recipients,
        Subject: "Watchmen alert!",
    }
    return mail
}

func (m *Mail) SendAsync(msgBody *string) {
    m.Send(msgBody)
}

func (m *Mail) Send(msgBody *string) {
    log.Printf("[notifier] Sending mail to: %v", m.Recipients)
    auth := smtp.PlainAuth("", m.User, m.Password, m.Host)

    recsStr := getRecipientsStr(m.Recipients)
    msg := []byte("To: " + *recsStr + "\r\n" +
        "Subject: " + m.Subject + "\r\n" +
        "\r\n" + *msgBody + "\r\n")

    err := smtp.SendMail(m.Host + ":" + strconv.Itoa(m.Port), auth, m.From, m.Recipients, msg)
    if err != nil {
        log.Println(err)
    }
}

func getRecipientsStr(recipients []string) *string {
    recsStr := ""
    for key, rec := range recipients{
        recsStr += rec
        // add comma separator between recipients
        if key + 1 < len(recipients) {
            recsStr += ","
        }
    }
    return &recsStr
}