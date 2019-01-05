package XML_IO

import (
	"TicketSystem/config"
	"encoding/xml"
	"io/ioutil"
)

//struct for mails
type Mail struct {
	Mail    string `xml:"Mail"`
	Caption string `xml:"Caption"`
	Message string `xml:"Message"`
	MailId  int    `xml:"MailId"`
}

type Maillist struct {
	MailIdCounter int    `xml:"MailIdCounter"`
	Maillist      []Mail `xml:"mails>mail"`
}

//creating or merging a ticket that was send by mail
func CreateTicketFromMail(mail string, reference string, message string) (Ticket, error) {
	tickets := GetTicketsByClient(mail)
	for _, actTicket := range tickets {
		if actTicket.Reference == reference {
			ChangeStatus(actTicket.Id, InProcess)
			return AddMessage(actTicket, mail, message)
		}

		//TODO: schauen ob Ticket bereits vorhanden ist

	}
	return CreateTicket(mail, reference, message)
}

//delete all mails in the xml file which are already sent
func DeleteMails(mailIds []int) error {
	maillist, err := readMailsFile()
	if err != nil {
		return err
	}

	mailMap := make(map[int]Mail)
	for _, actMail := range maillist.Maillist {
		mailMap[actMail.MailId] = actMail
	}

	for _, actId := range mailIds {
		delete(mailMap, actId)
	}

	var newMaillist Maillist
	for _, actMail := range mailMap {
		newMaillist.Maillist = append(newMaillist.Maillist, actMail)
	}
	newMaillist.MailIdCounter = maillist.MailIdCounter - len(mailIds)
	return WriteToXML(newMaillist, config.MailFilePath())
}

//store the message as a mail in the specific xml file
func SendMail(mail string, caption string, message string) error {
	maillist, err := readMailsFile()
	if err != nil {
		return err
	}
	nextMailId := maillist.MailIdCounter + 1
	newMail := Mail{Mail: mail, Caption: caption, Message: message, MailId: nextMailId}
	maillist.Maillist = append(maillist.Maillist, newMail)
	maillist.MailIdCounter = nextMailId
	return WriteToXML(maillist, config.MailFilePath())
}

//get all mails from the xml file
func readMailsFile() (Maillist, error) {
	file, err := ioutil.ReadFile(config.MailFilePath())
	if err != nil {
		return Maillist{}, err
	}
	var maillist Maillist
	err = xml.Unmarshal(file, &maillist)
	if err != nil {
		return Maillist{}, err
	}
	return maillist, nil
}
