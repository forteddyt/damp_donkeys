package resp

import(
	"github.com/damp_donkeys/internal/pkg/dbutil"
)

type GetCompanyCheckIn struct {
	CompanyName string `json:"company_name"`
	Students []dbutil.Interview `json:"students"`
	JWT string `json:"jwt"`
}

type GetStudent struct {
	DisplayName string `json:"display_name"`
	Major string `json:"major"`
	Class string `json:"class"`
}

type GetLogin struct {
	JWT string `json:"jwt"`
}

type GetCareerFairList struct {
	CareerFairList []string `json:"career_fair_list"`
	JWT string `json:"jwt"`
}

type GetCompanyList struct {
	CompanyList []string `json:"company_list"`
}

type PutCompany struct {
	CompanyName string `json:"company_name`
	UserCode string `json:"user_code"`
	JWT string `json:"jwt"`
}

// Same as AddCompanyResp, just different name
type PutResetCode PutCompany


