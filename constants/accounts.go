package constants

type EntityType string

const (
	EntityAccountBeneficiary EntityType = "ACCOUNT_BENEFICIARY"
	EntityAccountNominee     EntityType = "ACCOUNT_NOMINEE"
)

type AccountType string

const (
	AccountTypeBusiness       AccountType = "BUSINESS"
	AccountTypeCommodity      AccountType = "COMMODITY"
	AccountTypeCreditCard     AccountType = "CREDIT_CARD"
	AccountTypeCryptocurrency AccountType = "CRYPTOCURRENCY"
	AccountTypeCurrent        AccountType = "CURRENT"
	AccountTypeDemat          AccountType = "DEMAT"
	AccountTypeFixed          AccountType = "FIXED"
	AccountTypeForex          AccountType = "FOREX"
	AccountTypeInsurance      AccountType = "INSURANCE"
	AccountTypeInvestment     AccountType = "INVESTMENT"
	AccountTypeJoint          AccountType = "JOINT"
	AccountTypeLoan           AccountType = "LOAN"
	AccountTypeMicro          AccountType = "MICRO"
	AccountTypeMutualFund     AccountType = "MUTUAL_FUND"
	AccountTypeNPS            AccountType = "NPS"
	AccountTypeNRI            AccountType = "NRI"
	AccountTypeOther          AccountType = "OTHER"
	AccountTypePension        AccountType = "PENSION"
	AccountTypePPF            AccountType = "PPF"
	AccountTypeRecurring      AccountType = "RECURRING"
	AccountTypeSalary         AccountType = "SALARY"
	AccountTypeSavings        AccountType = "SAVINGS"
	AccountTypeSeniorCitizen  AccountType = "SENIOR_CITIZEN"
	AccountTypeSIP            AccountType = "SIP"
	AccountTypeStudent        AccountType = "STUDENT"
	AccountTypeZeroBalance    AccountType = "ZERO_BALANCE"
)

type AccountStatus string

const (
	AccountStatusActive           AccountStatus = "ACTIVE"
	AccountStatusArchived         AccountStatus = "ARCHIVED"
	AccountStatusBlacklisted      AccountStatus = "BLACKLISTED"
	AccountStatusBlocked          AccountStatus = "BLOCKED"
	AccountStatusClosed           AccountStatus = "CLOSED"
	AccountStatusDeactivated      AccountStatus = "DEACTIVATED"
	AccountStatusDeceased         AccountStatus = "DECEASED"
	AccountStatusDeleted          AccountStatus = "DELETED"
	AccountStatusDormant          AccountStatus = "DORMANT"
	AccountStatusFlagged          AccountStatus = "FLAGGED"
	AccountStatusFrozen           AccountStatus = "FROZEN"
	AccountStatusInactive         AccountStatus = "INACTIVE"
	AccountStatusOperational      AccountStatus = "OPERATIONAL"
	AccountStatusPending          AccountStatus = "PENDING"
	AccountStatusRejected         AccountStatus = "REJECTED"
	AccountStatusRestored         AccountStatus = "RESTORED"
	AccountStatusRestricted       AccountStatus = "RESTRICTED"
	AccountStatusSuspended        AccountStatus = "SUSPENDED"
	AccountStatusUnarchived       AccountStatus = "UNARCHIVED"
	AccountStatusUndeactivated    AccountStatus = "UNDEACTIVATED"
	AccountStatusUnderMaintenance AccountStatus = "UNDER_MAINTENANCE"
	AccountStatusUnflagged        AccountStatus = "UNFLAGGED"
	AccountStatusUnrestricted     AccountStatus = "UNRESTRICTED"
	AccountStatusUnsuspended      AccountStatus = "UNSUSPENDED"
)

type KYCStatus string

const (
	KYCStatusAadhaarVerified      KYCStatus = "KYC_AADHAAR_VERIFIED"
	KYCStatusAddressVerified      KYCStatus = "KYC_ADDRESS_VERIFIED"
	KYCStatusApproved             KYCStatus = "KYC_APPROVED"
	KYCStatusAutoApproved         KYCStatus = "KYC_AUTO_APPROVED"
	KYCStatusAutoRejected         KYCStatus = "KYC_AUTO_REJECTED"
	KYCStatusBiometricVerified    KYCStatus = "KYC_BIOMETRIC_VERIFIED"
	KYCStatusCompleted            KYCStatus = "KYC_COMPLETED"
	KYCStatusDocumentRejected     KYCStatus = "KYC_DOCUMENT_REJECTED"
	KYCStatusDocumentVerified     KYCStatus = "KYC_DOCUMENT_VERIFIED"
	KYCStatusEmailVerified        KYCStatus = "KYC_EMAIL_VERIFIED"
	KYCStatusExpired              KYCStatus = "KYC_EXPIRED"
	KYCStatusFaceVerified         KYCStatus = "KYC_FACE_VERIFIED"
	KYCStatusFailed               KYCStatus = "KYC_FAILED"
	KYCStatusInPersonVerified     KYCStatus = "KYC_IN_PERSON_VERIFIED"
	KYCStatusInProgress           KYCStatus = "KYC_IN_PROGRESS"
	KYCStatusManualReview         KYCStatus = "KYC_MANUAL_REVIEW"
	KYCStatusMobileVerified       KYCStatus = "KYC_MOBILE_VERIFIED"
	KYCStatusNotSubmitted         KYCStatus = "KYC_NOT_SUBMITTED"
	KYCStatusPanVerified          KYCStatus = "KYC_PAN_VERIFIED"
	KYCStatusPending              KYCStatus = "KYC_PENDING"
	KYCStatusRejected             KYCStatus = "KYC_REJECTED"
	KYCStatusResubmissionRequired KYCStatus = "KYC_RESUBMISSION_REQUIRED"
	KYCStatusThirdPartyVerified   KYCStatus = "KYC_THIRD_PARTY_VERIFIED"
	KYCStatusUnderReview          KYCStatus = "KYC_UNDER_REVIEW"
	KYCStatusVerified             KYCStatus = "KYC_VERIFIED"
	KYCStatusVideoVerified        KYCStatus = "KYC_VIDEO_VERIFIED"
)

type NomineeStatus string

const (
	NomineeStatusActive   NomineeStatus = "ACTIVE"
	NomineeStatusApproved NomineeStatus = "APPROVED"
	NomineeStatusDeleted  NomineeStatus = "DELETED"
	NomineeStatusPending  NomineeStatus = "PENDING"
	NomineeStatusRejected NomineeStatus = "REJECTED"
)

type BeneficiaryStatus string

const (
	BeneficiaryStatusActive   BeneficiaryStatus = "ACTIVE"
	BeneficiaryStatusApproved BeneficiaryStatus = "APPROVED"
	BeneficiaryStatusDeleted  BeneficiaryStatus = "DELETED"
	BeneficiaryStatusPending  BeneficiaryStatus = "PENDING"
	BeneficiaryStatusRejected BeneficiaryStatus = "REJECTED"
)

type Action string

const (
	ActionApprove Action = "APPROVE"
	ActionCreate  Action = "CREATE"
	ActionDelete  Action = "DELETE"
	ActionReject  Action = "REJECT"
	ActionUpdate  Action = "UPDATE"
)

type Status string

const (
	StatusApproved Status = "APPROVED"
	StatusRejected Status = "REJECTED"
)

type RelationType string

const (
	RelationDaughter RelationType = "DAUGHTER"
	RelationFather   RelationType = "FATHER"
	RelationMother   RelationType = "MOTHER"
	RelationSon      RelationType = "SON"
	RelationSpouse   RelationType = "SPOUSE"
)
