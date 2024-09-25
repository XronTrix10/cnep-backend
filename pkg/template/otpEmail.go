package template

import (
    "bytes"
    "html/template"
)

type OTPEmailData struct {
    Name string
    OTP  string
}

const OTPEmailTemplate = `
<!DOCTYPE html>
<html>
<body>
    <h1>Your OTP Code</h1>
    <br>
    <p>Hello {{.Name}},</p>
    <br>
    <p>Your One-Time Password (OTP) for authentication is:</p>
    <br>
    <h2>{{.OTP}}</h2>
    <br>
    <p>This OTP will expire in 15 minutes. Please do not share this code with anyone.</p>
    <br>
    <p>If you didn't request this OTP, please ignore this email.</p>
    <br>
    <p>This is an automated message, please do not reply.</p>
</body>
</html>
`

func GenerateOTPEmail(data OTPEmailData) (string, error) {
    tmpl, err := template.New("otpEmail").Parse(OTPEmailTemplate)
    if err != nil {
        return "", err
    }

    var body bytes.Buffer
    if err := tmpl.Execute(&body, data); err != nil {
        return "", err
    }

    return body.String(), nil
}