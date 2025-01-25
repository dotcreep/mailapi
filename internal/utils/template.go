package utils

type SendTemplate struct {
	Welcome     Welcome
	Reject      Reject
	Approved    Approved
	Registered  Registered
	Invoice     Invoice
	InvoicePaid Invoice
	Credentials Credentials
}

type Welcome struct {
	To  string `json:"to" form:"to" example:"mail@example.com"`
	OTP int    `json:"otp" form:"otp" example:"1234"`
}

type Reject struct {
	To         string `json:"to" form:"to" example:"mail@example.com"`
	ClientName string `json:"client_name" form:"client_name" example:"John Doe"`
	Reason     string `json:"reason" form:"reason" example:"this is my reason"`
}

type Approved struct {
	To         string `json:"to" form:"to" example:"mail@example.com"`
	ClientName string `json:"client_name" form:"client_name" example:"John Doe"`
}

type Registered struct {
	To         string `json:"to" form:"to" example:"mail@example.com"`
	ClientName string `json:"client_name" form:"client_name" example:"John Doe"`
}

type Invoice struct {
	To         string `json:"to" form:"to" example:"mail@example.com"`
	Attach     string `json:"attach" form:"attach" example:"path/to/file.pdf or file"`
	ClientName string `json:"client_name" form:"client_name" example:"John Doe"`
	URLUpload  string `json:"url_upload" form:"url_upload" example:"https://example.com/upload"`
}

type InvoicePaid struct {
	To            string `json:"to" form:"to" example:"mail@example.com"`
	Attach        string `json:"attach" form:"attach" example:"path/to/file.pdf or file"`
	ClientName    string `json:"client_name" form:"client_name" example:"John Doe"`
	URLUpload     string `json:"url_upload" form:"url_upload" example:"https://example.com/upload"`
	InvoiceNumber int    `json:"invoice_number" form:"invoice_number" example:"1234567890"`
	Total         int    `json:"total" form:"total" example:"100000"`
}

type Credentials struct {
	To                    string `json:"to" form:"to" example:"mail@example.com"`
	ClientName            string `json:"client_name" form:"client_name" example:"John Doe"`
	UsernameAdminMerchant string `json:"user_superadmin" form:"user_superadmin" example:"exampleusername"`
	PasswordAdminMerchant string `json:"password_superadmin" form:"password_superadmin" example:"examplepassword"`
	UsernameMerchant      string `json:"user_merchant" form:"user_merchant" example:"exampleusername"`
	PasswordMerchant      string `json:"password_merchant" form:"password_merchant" example:"examplepassword"`
	AppMobile             string `json:"app_mobile_url" form:"app_mobile_url" example:"https://fs.example.com/getfile/user/abcd123"`
	Website               string `json:"website" form:"website" example:"https://example.com"`
}
