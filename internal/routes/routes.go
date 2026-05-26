package routes

import (
	"banking-system-backend/constants"
	approvalRepoPkg "banking-system-backend/internal/approval/repository"
	approvalServicePkg "banking-system-backend/internal/approval/service"
	"banking-system-backend/internal/config"
	"banking-system-backend/internal/executors"
	"banking-system-backend/internal/handlers"
	"banking-system-backend/internal/kafka/producer"
	"banking-system-backend/internal/middlewares"
	"banking-system-backend/internal/repositories/mongorepo"
	"banking-system-backend/internal/repositories/redisrepo"
	"banking-system-backend/internal/services"
	"banking-system-backend/internal/validators"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	userRepo := mongorepo.NewUserRepository(config.DB)
	customerRepo := mongorepo.NewCustomerRepository(config.DB)
	accountRepo := mongorepo.NewAccountRepository(config.DB)
	counterRepo := mongorepo.NewCounterRepository(config.DB)
	beneficiaryRepo := mongorepo.NewBeneficiaryRepository(config.DB)
	accountBeneficiaryRepo := mongorepo.NewAccountBeneficiaryRepository(config.DB)
	nomineeRepo := mongorepo.NewNomineeRepository(config.DB)
	accountNomineeRepo := mongorepo.NewAccountNomineeRepository(config.DB)
	approvalRepo := approvalRepoPkg.NewApprovalRepository(config.DB)

	cacheRepo := redisrepo.NewCacheRepository()
	authService := services.NewAuthService(userRepo)
	ownershipService := services.NewOwnershipService(accountRepo, cacheRepo)
	accountService := services.NewAccountService(accountRepo, counterRepo, userRepo, customerRepo, ownershipService)
	nomineeService := services.NewNomineeService(nomineeRepo, accountNomineeRepo, customerRepo)
	beneficiaryService := services.NewBeneficiaryService(beneficiaryRepo, customerRepo)

	accountNomineeValidator := validators.NewAccountNomineeValidator(accountNomineeRepo, nomineeRepo)
	accountBeneficiaryValidator := validators.NewAccountBeneficiaryValidator(accountBeneficiaryRepo, beneficiaryRepo)

	nomineeProducer := producer.NewNomineeProducer()
	accountNomineeExecutor := executors.NewAccountNomineeExecutor(
		accountNomineeRepo,
		accountRepo,
		nomineeRepo,
		accountNomineeValidator,
		nomineeProducer,
	)

	beneficiaryProducer := producer.NewBeneficiaryProducer()
	accountBeneficiaryExecutor := executors.NewAccountBeneficiaryExecutor(
		accountRepo,
		accountBeneficiaryRepo,
		beneficiaryProducer,
	)
	executorsMP := map[string]approvalServicePkg.Executor{
		"ACCOUNT_NOMINEE:CREATE":     accountNomineeExecutor,
		"ACCOUNT_NOMINEE:UPDATE":     accountNomineeExecutor,
		"ACCOUNT_NOMINEE:DELETE":     accountNomineeExecutor,
		"ACCOUNT_BENEFICIARY:CREATE": accountBeneficiaryExecutor,
		"ACCOUNT_BENEFICIARY:UPDATE": accountBeneficiaryExecutor,
		"ACCOUNT_BENEFICIARY:DELETE": accountBeneficiaryExecutor,
	}

	approvalService := approvalServicePkg.NewApprovalService(approvalRepo, executorsMP)

	accountNomineeService := services.NewAccountNomineeService(accountNomineeRepo, accountRepo, nomineeRepo, customerRepo, approvalService, accountNomineeValidator)
	accountBeneficiaryService := services.NewAccountBeneficiaryService(accountBeneficiaryRepo, accountRepo, beneficiaryRepo, customerRepo, approvalService, accountBeneficiaryValidator)

	authHandler := handlers.NewAuthHandler(authService)
	accountHandler := handlers.NewAccountHandler(accountService)
	nomineeHandler := handlers.NewNomineeHandler(nomineeService)
	accountNomineeHandler := handlers.NewAccountNomineeHandler(accountNomineeService)
	beneficiaryHandler := handlers.NewBeneficiaryHandler(beneficiaryService)
	accountBeneficiaryHandler := handlers.NewAccountBeneficiaryHandler(accountBeneficiaryService)
	approvalHandler := handlers.NewApprovalHandler(approvalService)

	// AUTH
	auth := r.Group("/api/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// ACCOUNT
	account := r.Group("/api/accounts")
	account.Use(middlewares.AuthMiddleware())
	account.POST("/", middlewares.PolicyMiddleware(constants.PolicyAccountFullAccess), accountHandler.CreateAccount)
	account.GET("/:accountNumber", middlewares.PolicyMiddleware(constants.PolicyAccountFullAccess), accountHandler.GetAccount)
	account.PUT("/:accountNumber", middlewares.PolicyMiddleware(constants.PolicyAccountFullAccess), accountHandler.UpdateAccount)
	account.DELETE("/:accountNumber", middlewares.PolicyMiddleware(constants.PolicyAccountFullAccess), accountHandler.DeleteAccount)

	// Nominee Routes
	nominee := r.Group("/api/nominees")
	nominee.Use(middlewares.AuthMiddleware())
	nominee.POST("/",
		middlewares.PolicyMiddleware(constants.PolicyNomineeFullAccess),
		nomineeHandler.CreateNominee,
	)
	nominee.GET("/:nomineeId", middlewares.PolicyMiddleware(constants.PolicyNomineeFullAccess), nomineeHandler.GetNominee)
	nominee.PUT("/:nomineeId", middlewares.PolicyMiddleware(constants.PolicyNomineeFullAccess), nomineeHandler.UpdateNominee)
	nominee.DELETE("/:nomineeId", middlewares.PolicyMiddleware(constants.PolicyNomineeFullAccess), nomineeHandler.DeleteNominee)
	nominee.GET("/", middlewares.PolicyMiddleware(constants.PolicyNomineeFullAccess), nomineeHandler.ListNominees)

	// Account Nominee (Mapping) Routes
	account.POST("/:accountNumber/nominees", middlewares.PolicyMiddleware(constants.PolicyAccountNomineeFullAccess), accountNomineeHandler.AddAccountNominee)
	account.GET("/:accountNumber/nominees", middlewares.PolicyMiddleware(constants.PolicyAccountNomineeFullAccess), accountNomineeHandler.ListAccountNomineesByAccountNumber)

	account.PUT("/:accountNumber/nominees/:mappingId", middlewares.PolicyMiddleware(constants.PolicyAccountNomineeFullAccess), accountNomineeHandler.UpdateAccountNominee)
	account.DELETE("/:accountNumber/nominees/:mappingId", middlewares.PolicyMiddleware(constants.PolicyAccountNomineeFullAccess), accountNomineeHandler.SoftDeleteAccountNominee)

	// Beneficiary Routes
	beneficiary := r.Group("/api/beneficiaries")
	beneficiary.Use(middlewares.AuthMiddleware())
	beneficiary.POST("/",
		middlewares.PolicyMiddleware(constants.PolicyBeneficiaryFullAccess),
		beneficiaryHandler.CreateBeneficiary,
	)
	beneficiary.GET("/:beneficiaryId", middlewares.PolicyMiddleware(constants.PolicyBeneficiaryFullAccess), beneficiaryHandler.GetBeneficiary)
	beneficiary.PUT("/:beneficiaryId", middlewares.PolicyMiddleware(constants.PolicyBeneficiaryFullAccess), beneficiaryHandler.UpdateBeneficiary)
	beneficiary.DELETE("/:beneficiaryId", middlewares.PolicyMiddleware(constants.PolicyBeneficiaryFullAccess), beneficiaryHandler.DeleteBeneficiary)
	beneficiary.GET("/", middlewares.PolicyMiddleware(constants.PolicyBeneficiaryFullAccess), beneficiaryHandler.ListBeneficiaries)

	// Account Beneficiary (Mapping) Routes
	account.POST("/:accountNumber/beneficiaries", middlewares.PolicyMiddleware(constants.PolicyAccountBeneficiaryFullAccess), accountBeneficiaryHandler.AddAccountBeneficiary)
	account.GET("/:accountNumber/beneficiaries", middlewares.PolicyMiddleware(constants.PolicyAccountBeneficiaryFullAccess), accountBeneficiaryHandler.ListAccountBeneficiariesByAccountNumber)
	account.PUT("/:accountNumber/beneficiaries/:mappingId", middlewares.PolicyMiddleware(constants.PolicyAccountBeneficiaryFullAccess), accountBeneficiaryHandler.UpdateAccountBeneficiary)
	account.DELETE("/:accountNumber/beneficiaries/:mappingId", middlewares.PolicyMiddleware(constants.PolicyAccountBeneficiaryFullAccess), accountBeneficiaryHandler.SoftDeleteAccountBeneficiary)

	// Maker-Checker approval-requests Routes
	approval := r.Group("/api/approval-requests")
	approval.Use(middlewares.AuthMiddleware())

	approval.GET("/",
		middlewares.PolicyMiddleware(constants.PolicyNomineeApproval, constants.PolicyBeneficiaryApproval),
		approvalHandler.ListRequests) // optional
	approval.GET("/:requestID",
		middlewares.PolicyMiddleware(constants.PolicyNomineeApproval, constants.PolicyBeneficiaryApproval),
		approvalHandler.GetRequest) // optional

	approval.PUT("/:requestID/approve",
		middlewares.PolicyMiddleware(constants.PolicyNomineeApproval, constants.PolicyBeneficiaryApproval),
		approvalHandler.Approve,
	)

	approval.PUT("/:requestID/reject",
		middlewares.PolicyMiddleware(constants.PolicyNomineeApproval, constants.PolicyBeneficiaryApproval),
		approvalHandler.Reject,
	)
}
