package policy

import (
	"banking-system-backend/constants"
)

var ApprovalPolicyMap = map[constants.EntityType][]string{
	constants.EntityAccountNominee:     {constants.PolicyNomineeApproval},
	constants.EntityAccountBeneficiary: {constants.PolicyBeneficiaryApproval},
}
