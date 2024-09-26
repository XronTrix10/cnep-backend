package template

import (
    "bytes"
    "html/template"
)

type OTPEmailData struct {
    OTP  string
}

const OTPEmailTemplate = `
<!DOCTYPE html>
<html>
<body>
    <br>
    <p>Hey,</p>
    <p>Please verify your email using the following code</p>
    <br>
    <h2><strong>{{.OTP}}</strong></h2>
    <br>
    <p>This OTP will expire in 15 minutes. Please do not share this code with anyone.</p>
    <p>If you didn't request this OTP, please ignore this email.</p>
    <br>
    <p>This is an automated message, please do not reply.</p>
    <p>Best Regards,</p>
    <p>CNEP Team</p>
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
