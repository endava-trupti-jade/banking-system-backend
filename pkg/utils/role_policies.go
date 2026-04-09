package utils

import (
	"banking-system-backend/constants"
)

var RolePolicies = map[string][]string{
	constants.RoleAdmin: {
		constants.PolicyAccountFullAccess,
		constants.PolicyAccountReadOnly,
		constants.PolicyAdminManage,
		constants.PolicyBeneficiaryApproval,
		constants.PolicyBeneficiaryFullAccess,
		constants.PolicyCustomerManage,
		constants.PolicyNomineeApproval,
		constants.PolicyNomineeFullAccess,
		constants.PolicyTransactionReadOnly,
	},
	constants.RoleManager: {
		constants.PolicyBeneficiaryApproval,
		constants.PolicyNomineeApproval,
	},
	constants.RoleCustomer: {
		constants.PolicyAccountBeneficiaryFullAccess,
		constants.PolicyAccountNomineeFullAccess,
		constants.PolicyAccountReadOnly,
		constants.PolicyAccountTransfer,
		constants.PolicyNomineeFullAccess,
		constants.PolicyTransactionReadOnly,
	},
	constants.RoleAuditor: {
		constants.PolicyTransactionReadOnly,
	},
}
