package emailsender

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"sort"
	"strings"
	"time"

	"github.com/govalues/decimal"
	"github.com/wneessen/go-mail"

	"stori/model"
)

const (
	subject = "Stori - Account Summary"

	numberOfDecimals = 2
	byThree          = 3
)

//go:embed templates
var templates embed.FS

type (
	Sender struct {
		Config
	}

	Config struct {
		Host     string
		Port     int
		Username string
		Password string
	}

	templateData struct {
		TotalBalance        string
		MonthsData          []MonthsData
		AverageDebitAmount  string
		AverageCreditAmount string
	}

	MonthsData struct {
		Month string
		Count int
	}
)

func New(config Config) Sender {
	return Sender{Config: config}
}

func (sender Sender) Send(summary model.AccountSummary) error {
	fail := func(err error) error {
		return fmt.Errorf("emailsender: sender: Send: %w", err)
	}

	message, err := sender.buildSummaryMessage(summary)
	if err != nil {
		return fail(err)
	}

	client, err := mail.NewClient(sender.Host, mail.WithPort(sender.Port), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(sender.Username), mail.WithPassword(sender.Password))
	if err != nil {
		return fail(fmt.Errorf("failed to create mail client: %w", err))
	}

	if errDial := client.DialAndSend(message); errDial != nil {
		return fail(fmt.Errorf("failed to send mail: %w", errDial))
	}

	return nil
}

func (sender Sender) buildSummaryMessage(summary model.AccountSummary) (*mail.Msg, error) {
	fail := func(err error) (*mail.Msg, error) {
		return nil, fmt.Errorf("emailsender: sender: buildSummaryMessage: %w", err)
	}

	msg := mail.NewMsg()
	if err := sender.addMetadata(msg, summary.Email); err != nil {
		return fail(err)
	}

	if err := sender.addBody(msg, summary); err != nil {
		return fail(err)
	}

	return msg, nil
}

func (sender Sender) addMetadata(msg *mail.Msg, receiverAddress string) error {
	fail := func(err error) error {
		return fmt.Errorf("emailsender: sender: addMetadata: %w", err)
	}

	if err := msg.From(sender.Username); err != nil {
		return fail(ErrWrongFormAddress)
	}

	if err := msg.To(receiverAddress); err != nil {
		return fail(ErrWrongTargetAddress)
	}

	msg.Subject(subject)

	return nil
}

func (sender Sender) addBody(msg *mail.Msg, summary model.AccountSummary) error {
	fail := func(err error) error {
		return fmt.Errorf("emailsender: sender: addBody: %w", err)
	}

	templ, err := template.ParseFS(templates, "templates/summary.html")
	if err != nil {
		return fail(err)
	}

	if err = msg.SetBodyHTMLTemplate(templ, buildTemplateData(summary)); err != nil {
		return fail(err)
	}

	return nil
}

func buildTemplateData(summary model.AccountSummary) templateData {
	totalBalance := printableAmount(summary.TotalBalance)
	avgDebit := printableAmount(summary.AverageDebitAmount)
	avgCredit := printableAmount(summary.AverageCreditAmount)
	monthsData := buildMonthData(summary)

	return templateData{
		TotalBalance:        totalBalance,
		MonthsData:          monthsData,
		AverageDebitAmount:  avgDebit,
		AverageCreditAmount: avgCredit,
	}
}

func buildMonthData(summary model.AccountSummary) []MonthsData {
	keys := make([]time.Month, 0, len(summary.TransactionsPerMonth))

	for k := range summary.TransactionsPerMonth {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	var monthsData []MonthsData
	for _, k := range keys {
		monthsData = append(monthsData, MonthsData{Month: k.String(), Count: summary.TransactionsPerMonth[k]})
	}

	return monthsData
}

func printableAmount(amount decimal.Decimal) string {
	return fmt.Sprintf("$ %s", addCommas(amount.Round(numberOfDecimals).Pad(numberOfDecimals).String()))
}

func addCommas(decimalStr string) string {
	buf := &bytes.Buffer{}
	comma := []byte{','}

	parts := strings.Split(decimalStr, ".")
	pos := 0

	if len(parts[0])%byThree != 0 {
		pos += len(parts[0]) % byThree
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}

	for ; pos < len(parts[0]); pos += byThree {
		buf.WriteString(parts[0][pos : pos+byThree])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}

	return buf.String()
}
