package constants

type Policy struct {
	Code              string
	Description       string
	Permissions       []string
	Conditions        []string
	OwnershipRequired bool
}

var Policies = map[string]Policy{
	// Admin — Full account control
	PolicyAccountFullAccess: {
		Code:              PolicyAccountFullAccess,
		Description:       "Full control over accounts",
		Permissions:       []string{PermAccRead, PermAccWrite, PermAccDelete},
		Conditions:        []string{CondCustomerOrAdmin, CondAccountActive},
		OwnershipRequired: false, // Admin bypass
	},

	// Account Read — Owner or Admin
	PolicyAccountReadOnly: {
		Code:              PolicyAccountReadOnly,
		Description:       "Account read only access",
		Permissions:       []string{PermAccRead},
		Conditions:        []string{},
		OwnershipRequired: true, // Only enforced for non-admin
	},

	PolicyAccountApproval: {
		Code:              PolicyAccountApproval,
		Description:       "Account approve or reject access",
		Permissions:       []string{PermAccRead},
		Conditions:        []string{},
		OwnershipRequired: true, // Only enforced for non-admin
	},

	// Nominee
	PolicyNomineeFullAccess: {
		Code:              PolicyNomineeFullAccess,
		Description:       "Full control over nominee",
		Permissions:       []string{PermNomineeRead, PermNomineeWrite, PermNomineeDelete},
		Conditions:        []string{},
		OwnershipRequired: true,
	},

	// Account Nominee Approval
	PolicyNomineeApproval: {
		Code:              PolicyNomineeApproval,
		Description:       "Nominee approve or reject access",
		Permissions:       []string{PermNomineeApprove, PermNomineeReject, PermNomineeOnHold},
		Conditions:        []string{},
		OwnershipRequired: false,
	},

	// Account Nominee
	PolicyAccountNomineeFullAccess: {
		Code:              PolicyAccountNomineeFullAccess,
		Description:       "Full control over account nominee",
		Permissions:       []string{PermAccountNomineeRead, PermAccountNomineeWrite, PermAccountNomineeDelete},
		Conditions:        []string{},
		OwnershipRequired: true,
	},

	// Beneficiary
	PolicyBeneficiaryFullAccess: {
		Code:              PolicyBeneficiaryFullAccess,
		Description:       "Full control over beneficiary",
		Permissions:       []string{PermBeneficiaryRead, PermBeneficiaryWrite, PermBeneficiaryDelete},
		Conditions:        []string{},
		OwnershipRequired: true,
	},

	// Account Beneficiary Approval
	PolicyBeneficiaryApproval: {
		Code:              PolicyBeneficiaryApproval,
		Description:       "Beneficiary approve or reject access",
		Permissions:       []string{PermBeneficiaryApprove, PermBeneficiaryReject, PermBeneficiaryOnHold},
		Conditions:        []string{},
		OwnershipRequired: false,
	},

	// Account Beneficiary
	PolicyAccountBeneficiaryFullAccess: {
		Code:              PolicyAccountBeneficiaryFullAccess,
		Description:       "Full control over account beneficiary",
		Permissions:       []string{PermAccountBeneficiaryRead, PermAccountBeneficiaryWrite, PermAccountBeneficiaryDelete},
		Conditions:        []string{},
		OwnershipRequired: true,
	},

	// Transfer — Owner only
	PolicyAccountTransfer: {
		Code:              PolicyAccountTransfer,
		Description:       "Allow transfers only from owned active accounts",
		Permissions:       []string{PermAccRead, PermAccTransfer, PermTransactionWrite},
		Conditions:        []string{CondTransferAllowed, CondAccountActive},
		OwnershipRequired: true,
	},

	// Transactions Read-only
	PolicyTransactionReadOnly: {
		Code:              PolicyTransactionReadOnly,
		Description:       "Read-only access to transactions",
		Permissions:       []string{PermTransactionRead},
		Conditions:        []string{},
		OwnershipRequired: false,
	},

	// Customer Manage — Admin only
	PolicyCustomerManage: {
		Code:              PolicyCustomerManage,
		Description:       "Full control over customers",
		Permissions:       []string{PermCustManage},
		Conditions:        []string{},
		OwnershipRequired: false,
	},

	// Customer Read — Admin only
	PolicyCustomerRead: {
		Code:              PolicyCustomerRead,
		Description:       "Read-only access to customers",
		Permissions:       []string{PermCustRead},
		Conditions:        []string{},
		OwnershipRequired: false,
	},

	// Admin Manage — Admin only
	PolicyAdminManage: {
		Code:              PolicyAdminManage,
		Description:       "Full control over admin users",
		Permissions:       []string{PermAdminManage},
		Conditions:        []string{},
		OwnershipRequired: false,
	},

	// Admin Read — Admin only
	PolicyAdminRead: {
		Code:              PolicyAdminRead,
		Description:       "Read-only access to admin users",
		Permissions:       []string{PermAdminRead},
		Conditions:        []string{},
		OwnershipRequired: false,
	},

	// Admin Write — Admin only
	PolicyAdminWrite: {
		Code:              PolicyAdminWrite,
		Description:       "Write access to admin users",
		Permissions:       []string{PermAdminWrite},
		Conditions:        []string{},
		OwnershipRequired: false,
	},

	// Admin Delete — Admin only
	PolicyAdminDelete: {
		Code:              PolicyAdminDelete,
		Description:       "Delete access to admin users",
		Permissions:       []string{PermAdminDelete},
		Conditions:        []string{},
		OwnershipRequired: false,
	},

	// Admin Activate — Admin only
	PolicyAdminActivate: {
		Code:              PolicyAdminActivate,
		Description:       "Activate admin users",
		Permissions:       []string{PermAdminActivate},
		Conditions:        []string{},
		OwnershipRequired: false,
	},

	// Admin Deactivate — Admin only
	PolicyAdminDeactivate: {
		Code:              PolicyAdminDeactivate,
		Description:       "Deactivate admin users",
		Permissions:       []string{PermAdminDeactivate},
		Conditions:        []string{},
		OwnershipRequired: false,
	},

	// Admin Suspend — Admin only
	PolicyAdminSuspend: {
		Code:              PolicyAdminSuspend,
		Description:       "Suspend admin users",
		Permissions:       []string{PermAdminSuspend},
		Conditions:        []string{},
		OwnershipRequired: false,
	},
}
