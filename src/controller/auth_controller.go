package controller

import (

	"app/src/model"
	"app/src/response"
	"app/src/service"
	"app/src/validation" 
	"github.com/gofiber/fiber/v2"

)

type AuthController struct {
	AuthService  service.AuthService
	UserService  service.UserService
	TokenService service.TokenService
	EmailService service.EmailService
}

func NewAuthController(
	authService service.AuthService, userService service.UserService,
	tokenService service.TokenService, emailService service.EmailService,
) *AuthController {
	return &AuthController{
		AuthService:  authService,
		UserService:  userService,
		TokenService: tokenService,
		EmailService: emailService,
	}
}

// @Tags         Auth
// @Summary      Register as user
// @Accept       json
// @Produce      json
// @Param        request  body  validation.Register  true  "Request body"
// @Router       /auth/register [post]
// @Success      201  {object}  example.RegisterResponse
// @Failure      409  {object}  example.DuplicateEmail  "Email already taken"
func (a *AuthController) Register(c *fiber.Ctx) error {
	req := new(validation.Register)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := a.AuthService.Register(c, req)
	if err != nil {
		return err
	}

	tokens, err := a.TokenService.GenerateAuthTokens(c, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).
		JSON(response.SuccessWithTokens{
			Code:    fiber.StatusCreated,
			Status:  "success",
			Message: "Register successfully",
			User:    *user,
			Tokens:  *tokens,
		})
}

// @Tags         Auth
// @Summary      Login
// @Accept       json
// @Produce      json
// @Param        request  body  validation.Login  true  "Request body"
// @Router       /auth/login [post]
// @Success      200  {object}  example.LoginResponse
// @Failure      401  {object}  example.FailedLogin  "Invalid email or password"
func (a *AuthController) Login(c *fiber.Ctx) error {
	req := new(validation.Login)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := a.AuthService.Login(c, req)
	if err != nil {
		return err
	}

	tokens, err := a.TokenService.GenerateAuthTokens(c, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithTokens{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Login successfully",
			User:    *user,
			Tokens:  *tokens,
		})
}

// @Tags         Auth
// @Summary      Logout
// @Accept       json
// @Produce      json
// @Param        request  body  example.RefreshToken  true  "Request body"
// @Router       /auth/logout [post]
// @Success      200  {object}  example.LogoutResponse
// @Failure      404  {object}  example.NotFound  "Not found"
func (a *AuthController) Logout(c *fiber.Ctx) error {
	req := new(validation.Logout)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := a.AuthService.Logout(c, req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Logout successfully",
		})
}

// @Tags         Auth
// @Summary      Refresh auth tokens
// @Accept       json
// @Produce      json
// @Param        request  body  example.RefreshToken  true  "Request body"
// @Router       /auth/refresh-tokens [post]
// @Success      200  {object}  example.RefreshTokenResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
func (a *AuthController) RefreshTokens(c *fiber.Ctx) error {
	req := new(validation.RefreshToken)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	tokens, err := a.AuthService.RefreshAuth(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.RefreshToken{
			Code:   fiber.StatusOK,
			Status: "success",
			Tokens: *tokens,
		})
}

// @Tags         Auth
// @Summary      Forgot password
// @Description  An email will be sent to reset password.
// @Accept       json
// @Produce      json
// @Param        request  body  validation.ForgotPassword  true  "Request body"
// @Router       /auth/forgot-password [post]
// @Success      200  {object}  example.ForgotPasswordResponse
// @Failure      404  {object}  example.NotFound  "Not found"
func (a *AuthController) ForgotPassword(c *fiber.Ctx) error {
	req := new(validation.ForgotPassword)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	resetPasswordToken, err := a.TokenService.GenerateResetPasswordToken(c, req)
	if err != nil {
		return err
	}

	if errEmail := a.EmailService.SendResetPasswordEmail(req.Email, resetPasswordToken); errEmail != nil {
		return errEmail
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "A password reset link has been sent to your email address.",
		})
}

// @Tags         Auth
// @Summary      Reset password
// @Accept       json
// @Produce      json
// @Param        token   query  string  true  "The reset password token"
// @Param        request  body  validation.UpdatePassOrVerify  true  "Request body"
// @Router       /auth/reset-password [post]
// @Success      200  {object}  example.ResetPasswordResponse
// @Failure      401  {object}  example.FailedResetPassword  "Password reset failed"
func (a *AuthController) ResetPassword(c *fiber.Ctx) error {
	req := new(validation.UpdatePassOrVerify)
	query := &validation.Token{
		Token: c.Query("token"),
	}

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := a.AuthService.ResetPassword(c, query, req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Update password successfully",
		})
}

// @Tags         Auth
// @Summary      Send verification email
// @Description  An email will be sent to verify email.
// @Security BearerAuth
// @Produce      json
// @Router       /auth/send-verification-email [post]
// @Success      200  {object}  example.SendVerificationEmailResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
func (a *AuthController) SendVerificationEmail(c *fiber.Ctx) error {
	user, _ := c.Locals("user").(*model.User)

	verifyEmailToken, err := a.TokenService.GenerateVerifyEmailToken(c, user)
	if err != nil {
		return err
	}

	if errEmail := a.EmailService.SendVerificationEmail(user.Email, *verifyEmailToken); errEmail != nil {
		return errEmail
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Please check your email for a link to verify your account",
		})
}

// @Tags         Auth
// @Summary      Verify email
// @Produce      json
// @Param        token   query  string  true  "The verify email token"
// @Router       /auth/verify-email [post]
// @Success      200  {object}  example.VerifyEmailResponse
// @Failure      401  {object}  example.FailedVerifyEmail  "Verify email failed"
func (a *AuthController) VerifyEmail(c *fiber.Ctx) error {
	query := &validation.Token{
		Token: c.Query("token"),
	}

	if err := a.AuthService.VerifyEmail(c, query); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Verify email successfully",
		})
}
// @Tags         Auth
// @Summary      Get current user
// @Security     BearerAuth
// @Produce      json
// @Router       /auth/me [get]
// @Success      200  {object}  response.SuccessWithCurrentUser
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
func (a *AuthController) GetCurrentUser(c *fiber.Ctx) error {
	user, err := a.AuthService.GetCurrentUser(c)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithCurrentUser{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "User retrieved successfully",
			User:    *user,
		})
}
