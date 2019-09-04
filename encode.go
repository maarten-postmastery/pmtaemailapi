package main

import (
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"time"
)

func encode(w io.Writer, m *message) {

	if m.Sender != "" {
		writeHeaderField(w, "x-sender", m.Sender)
	} else {
		writeHeaderField(w, "x-sender", m.From.Address.Address)
	}
	if m.To != nil {
		for _, a := range m.To.recipients() {
			writeHeaderField(w, "x-receiver", a)
		}
	}
	if m.Cc != nil {
		for _, a := range m.Cc.recipients() {
			writeHeaderField(w, "x-receiver", a)
		}
	}
	if m.Bcc != nil {
		for _, a := range m.Bcc.recipients() {
			writeHeaderField(w, "x-receiver", a)
		}
	}

	writeHeaderField(w, "Date", time.Now().Format(time.RFC1123Z))
	writeHeaderField(w, "From", m.From.String())

	// optional header fields
	if m.To != nil {
		writeHeaderField(w, "To", m.To.String())
	}
	if m.Cc != nil {
		writeHeaderField(w, "Cc", m.Cc.String())
	}
	if m.ReplyTo != nil {
		writeHeaderField(w, "Reply-To", m.ReplyTo.String())
	}
	if m.Subject != "" {
		writeHeaderField(w, "Subject", mime.QEncoding.Encode("utf-8", m.Subject))
	}

	switch {
	case m.Text != "" && m.HTML != "":
		mpw := multipart.NewWriter(w)
		writeMultipartHeader(w, "alternative", mpw.Boundary())

		th := make(textproto.MIMEHeader)
		th.Set("Content-Type", "text/plain; charset=utf-8")
		th.Set("Content-Transfer-Encoding", "quoted-printable")
		tp, _ := mpw.CreatePart(th)
		writeQuotedPrintable(tp, m.Text)

		hh := make(textproto.MIMEHeader)
		hh.Set("Content-Type", "text/html; charset=utf-8")
		hh.Set("Content-Transfer-Encoding", "quoted-printable")
		mp, _ := mpw.CreatePart(hh)
		writeQuotedPrintable(mp, m.HTML)

		mpw.Close()
	case m.Text != "":
		writeTextBodyHeader(w, "plain")
		io.WriteString(w, "\r\n")
		writeQuotedPrintable(w, m.Text)
	case m.HTML != "":
		writeTextBodyHeader(w, "html")
		io.WriteString(w, "\r\n")
		writeQuotedPrintable(w, m.HTML)
	default:
		io.WriteString(w, "\r\n") // empty body
	}
}

func writeHeaderField(w io.Writer, name, body string) {
	io.WriteString(w, name)
	io.WriteString(w, ": ")
	io.WriteString(w, body) // TODO: folding
	io.WriteString(w, "\r\n")
}

func writeQuotedPrintable(w io.Writer, text string) {
	qpw := quotedprintable.NewWriter(w)
	io.WriteString(qpw, text)
	qpw.Close()
}

func writeMultipartHeader(w io.Writer, subtype string, boundary string) {
	writeHeaderField(w, "MIME-Version", "1.0")
	writeHeaderField(w, "Content-Type", "multipart/"+subtype+
		";\r\n boundary=\""+boundary+"\"")
}

func writeTextBodyHeader(w io.Writer, subtype string) {
	writeHeaderField(w, "MIME-Version", "1.0")
	writeHeaderField(w, "Content-Type", "text/"+subtype+"; charset=utf-8")
	writeHeaderField(w, "Content-Transfer-Encoding", "quoted-printable")
}
