package main

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain string
	Host string
	Port int
	Username string
	Password string
	Encryption string
	FromAddress string
	FromName string
}

type Message struct {
	From string
	FromName string
	To string
	Subject string
	Attachments []string
	Data any
	DataMap map[string]any
}


func (m *Mail) sendSMTPMessage(msg Message) error {
	
	if msg.From == "" {
		msg.From = m.FromAddress;
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName;
	}

	data := map[string]any {
		"message": msg.Data,
	}
	msg.DataMap = data;

	formattedMessage, err := m.buildHTMLMessage(msg);


	if err != nil {
		fmt.Printf("HTML ERROR: %v\n",err);
		return err
	}
	
	plainMessage, err := m.buildPlainMessage(msg);
	
	if err != nil {
		fmt.Printf("PLAIN ERROR: %v\n",err);
		return err
	}

	server := mail.NewSMTPClient()

	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10*time.Second
	server.SendTimeout = 10*time.Second
	
	smtpClient, err := server.Connect();
	
	
	if err != nil {
		fmt.Printf("SMTP CONNECT ERROR: %v\n",err);
		return err
	}

	email := mail.NewMSG()

	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, plainMessage);
	email.AddAlternative(mail.TextHTML, formattedMessage);


	if len(msg.Attachments) > 0 {
		for _, atc := range msg.Attachments {
			email.AddAttachment(atc);
		}
	}

	err  = email.Send(smtpClient);
	if err != nil {
		return err
	}
	return nil;
}

func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
		case "tls":
			return mail.EncryptionTLS
		case "ssl":
			return mail.EncryptionSSLTLS
		case "none", "":
			return mail.EncryptionNone
		default:
			return mail.EncryptionNone
	}
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	var templateToRender = "./templates/mail.html.gohtml";

	t, err := template.New("email-html").ParseFiles(templateToRender);

	if err != nil {
		fmt.Printf("PARSE ERROR %v\n",err);
		return "", err
	}

	var tpl bytes.Buffer

	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String();
	formattedMessage, err = m.inlineCSS(formattedMessage);

	if err != nil {
		return "", err
	}

	return formattedMessage, nil;
}

func (m *Mail) buildPlainMessage(msg Message) (string, error) {
	var templateToRender = "./templates/mail.plain.gohtml";

	t, err := template.New("email-plain").ParseFiles(templateToRender);

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String();

	if err != nil {
		return "", err
	}

	return formattedMessage, nil;
}


func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses: false,
		CssToAttributes: false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options);
	
	if err != nil {
		return "",err
	}

	html, err := prem.Transform()

	if err != nil {
		return "", err
	}

	return html, nil
	
}