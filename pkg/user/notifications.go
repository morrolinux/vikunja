// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-present Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public Licensee as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public Licensee for more details.
//
// You should have received a copy of the GNU Affero General Public Licensee
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package user

import (
	"strconv"

	"code.vikunja.io/api/pkg/config"
	"code.vikunja.io/api/pkg/i18n"
	"code.vikunja.io/api/pkg/notifications"
)

// EmailConfirmNotification represents a EmailConfirmNotification notification
type EmailConfirmNotification struct {
	User         *User
	IsNew        bool
	ConfirmToken string
}

// ToMail returns the mail notification for EmailConfirmNotification
func (n *EmailConfirmNotification) ToMail() *notifications.Mail {

	subject := i18n.TWithParams(n.User.Language, "notifications.email_confirm.subject", n.User.GetName())
	if n.IsNew {
		subject = i18n.TWithParams(n.User.Language, "notifications.email_confirm.subject_new", n.User.GetName())
	}

	nn := notifications.NewMail().
		Subject(subject).
		Greeting(i18n.TWithParams(n.User.Language, "notifications.greeting", n.User.GetName()))

	if n.IsNew {
		nn.Line(i18n.T(n.User.Language, "notifications.email_confirm.welcome"))
	}

	return nn.
		Line(i18n.T(n.User.Language, "notifications.email_confirm.confirm")).
		Action("Confirm your email address", config.ServicePublicURL.GetString()+"?userEmailConfirm="+n.ConfirmToken).
		Line(i18n.T(n.User.Language, "notifications.common.have_nice_day"))
}

// ToDB returns the EmailConfirmNotification notification in a format which can be saved in the db
func (n *EmailConfirmNotification) ToDB() interface{} {
	return nil
}

// Name returns the name of the notification
func (n *EmailConfirmNotification) Name() string {
	return ""
}

// PasswordChangedNotification represents a PasswordChangedNotification notification
type PasswordChangedNotification struct {
	User *User
}

// ToMail returns the mail notification for PasswordChangedNotification
func (n *PasswordChangedNotification) ToMail() *notifications.Mail {
	return notifications.NewMail().
		Subject(i18n.T(n.User.Language, "notifications.password.changed.subject")).
		Greeting(i18n.TWithParams(n.User.Language, "notifications.greeting", n.User.GetName())).
		Line(i18n.T(n.User.Language, "notifications.password.changed.success")).
		Line(i18n.T(n.User.Language, "notifications.password.changed.warning"))
}

// ToDB returns the PasswordChangedNotification notification in a format which can be saved in the db
func (n *PasswordChangedNotification) ToDB() interface{} {
	return nil
}

// Name returns the name of the notification
func (n *PasswordChangedNotification) Name() string {
	return ""
}

// ResetPasswordNotification represents a ResetPasswordNotification notification
type ResetPasswordNotification struct {
	User  *User
	Token *Token
}

// ToMail returns the mail notification for ResetPasswordNotification
func (n *ResetPasswordNotification) ToMail() *notifications.Mail {
	return notifications.NewMail().
		Subject(i18n.T(n.User.Language, "notifications.password.reset.subject")).
		Greeting(i18n.TWithParams(n.User.Language, "notifications.greeting", n.User.GetName())).
		Line(i18n.T(n.User.Language, "notifications.password.reset.instructions")).
		Action("Reset your password", config.ServicePublicURL.GetString()+"?userPasswordReset="+n.Token.Token).
		Line(i18n.T(n.User.Language, "notifications.password.reset.valid_duration")).
		Line(i18n.T(n.User.Language, "notifications.common.have_nice_day"))
}

// ToDB returns the ResetPasswordNotification notification in a format which can be saved in the db
func (n *ResetPasswordNotification) ToDB() interface{} {
	return nil
}

// Name returns the name of the notification
func (n *ResetPasswordNotification) Name() string {
	return ""
}

// InvalidTOTPNotification represents a InvalidTOTPNotification notification
type InvalidTOTPNotification struct {
	User *User
}

// ToMail returns the mail notification for InvalidTOTPNotification
func (n *InvalidTOTPNotification) ToMail() *notifications.Mail {
	return notifications.NewMail().
		Subject(i18n.T(n.User.Language, "notifications.totp.invalid.subject")).
		Greeting(i18n.TWithParams(n.User.Language, "notifications.greeting", n.User.GetName())).
		Line(i18n.T(n.User.Language, "notifications.totp.invalid.message")).
		Line(i18n.T(n.User.Language, "notifications.totp.invalid.warning")).
		Action("Reset your password", config.ServicePublicURL.GetString()+"get-password-reset")
}

// ToDB returns the InvalidTOTPNotification notification in a format which can be saved in the db
func (n *InvalidTOTPNotification) ToDB() interface{} {
	return nil
}

// Name returns the name of the notification
func (n *InvalidTOTPNotification) Name() string {
	return "totp.invalid"
}

// PasswordAccountLockedAfterInvalidTOTOPNotification represents a PasswordAccountLockedAfterInvalidTOTOPNotification notification
type PasswordAccountLockedAfterInvalidTOTOPNotification struct {
	User *User
}

// ToMail returns the mail notification for PasswordAccountLockedAfterInvalidTOTOPNotification
func (n *PasswordAccountLockedAfterInvalidTOTOPNotification) ToMail() *notifications.Mail {
	resetURL := config.ServicePublicURL.GetString() + "get-password-reset"
	return notifications.NewMail().
		Subject(i18n.T(n.User.Language, "notifications.totp.account_locked.subject")).
		Greeting(i18n.TWithParams(n.User.Language, "notifications.greeting", n.User.GetName())).
		Line(i18n.T(n.User.Language, "notifications.totp.account_locked.message")).
		Line(i18n.T(n.User.Language, "notifications.totp.account_locked.disabled")).
		Line(i18n.TWithParams(n.User.Language, "notifications.totp.account_locked.reset_instructions", resetURL, resetURL))
}

// ToDB returns the PasswordAccountLockedAfterInvalidTOTOPNotification notification in a format which can be saved in the db
func (n *PasswordAccountLockedAfterInvalidTOTOPNotification) ToDB() interface{} {
	return nil
}

// Name returns the name of the notification
func (n *PasswordAccountLockedAfterInvalidTOTOPNotification) Name() string {
	return "password.account.locked.after.invalid.totop"
}

// FailedLoginAttemptNotification represents a FailedLoginAttemptNotification notification
type FailedLoginAttemptNotification struct {
	User *User
}

// ToMail returns the mail notification for FailedLoginAttemptNotification
func (n *FailedLoginAttemptNotification) ToMail() *notifications.Mail {
	return notifications.NewMail().
		Subject(i18n.T(n.User.Language, "notifications.login.failed.subject")).
		Greeting(i18n.TWithParams(n.User.Language, "notifications.greeting", n.User.GetName())).
		Line(i18n.T(n.User.Language, "notifications.login.failed.message")).
		Line(i18n.T(n.User.Language, "notifications.login.failed.warning")).
		Line(i18n.T(n.User.Language, "notifications.login.failed.enhance_security")).
		Action("Go to settings", config.ServicePublicURL.GetString()+"user/settings")
}

// ToDB returns the FailedLoginAttemptNotification notification in a format which can be saved in the db
func (n *FailedLoginAttemptNotification) ToDB() interface{} {
	return nil
}

// Name returns the name of the notification
func (n *FailedLoginAttemptNotification) Name() string {
	return "failed.login.attempt"
}

// AccountDeletionConfirmNotification represents a AccountDeletionConfirmNotification notification
type AccountDeletionConfirmNotification struct {
	User         *User
	ConfirmToken string
}

// ToMail returns the mail notification for AccountDeletionConfirmNotification
func (n *AccountDeletionConfirmNotification) ToMail() *notifications.Mail {
	return notifications.NewMail().
		Subject(i18n.T(n.User.Language, "notifications.account.deletion.confirm.subject")).
		Greeting(i18n.TWithParams(n.User.Language, "notifications.greeting", n.User.GetName())).
		Line(i18n.T(n.User.Language, "notifications.account.deletion.confirm.request")).
		Action("Confirm the deletion of my account", config.ServicePublicURL.GetString()+"?accountDeletionConfirm="+n.ConfirmToken).
		Line(i18n.T(n.User.Language, "notifications.account.deletion.confirm.valid_duration")).
		Line(i18n.T(n.User.Language, "notifications.account.deletion.confirm.schedule_info")).
		Line(i18n.T(n.User.Language, "notifications.account.deletion.confirm.consequences")).
		Line(i18n.T(n.User.Language, "notifications.account.deletion.confirm.changed_mind")).
		Line(i18n.T(n.User.Language, "notifications.common.have_nice_day"))
}

// ToDB returns the AccountDeletionConfirmNotification notification in a format which can be saved in the db
func (n *AccountDeletionConfirmNotification) ToDB() interface{} {
	return nil
}

// Name returns the name of the notification
func (n *AccountDeletionConfirmNotification) Name() string {
	return "user.deletion.confirm"
}

// AccountDeletionNotification represents a AccountDeletionNotification notification
type AccountDeletionNotification struct {
	User               *User
	NotificationNumber int
}

// ToMail returns the mail notification for AccountDeletionNotification
func (n *AccountDeletionNotification) ToMail() *notifications.Mail {
	var subject string
	var deletionTimeLine string

	if n.NotificationNumber == 1 {
		subject = i18n.T(n.User.Language, "notifications.account.deletion.scheduled.subject_tomorrow")
		deletionTimeLine = i18n.T(n.User.Language, "notifications.account.deletion.scheduled.deletion_time_tomorrow")
	} else {
		days := strconv.Itoa(n.NotificationNumber)
		subject = i18n.TWithParams(n.User.Language, "notifications.account.deletion.scheduled.subject_days", days)
		deletionTimeLine = i18n.TWithParams(n.User.Language, "notifications.account.deletion.scheduled.deletion_time_days", days)
	}

	return notifications.NewMail().
		Subject(subject).
		Greeting(i18n.TWithParams(n.User.Language, "notifications.greeting", n.User.GetName())).
		Line(i18n.T(n.User.Language, "notifications.account.deletion.scheduled.request_reminder")).
		Line(deletionTimeLine).
		Line(i18n.T(n.User.Language, "notifications.account.deletion.scheduled.changed_mind")).
		Action("Abort the deletion", config.ServicePublicURL.GetString()).
		Line(i18n.T(n.User.Language, "notifications.common.have_nice_day"))
}

// ToDB returns the AccountDeletionNotification notification in a format which can be saved in the db
func (n *AccountDeletionNotification) ToDB() interface{} {
	return nil
}

// Name returns the name of the notification
func (n *AccountDeletionNotification) Name() string {
	return "user.deletion"
}

// AccountDeletedNotification represents a AccountDeletedNotification notification
type AccountDeletedNotification struct {
	User *User
}

// ToMail returns the mail notification for AccountDeletedNotification
func (n *AccountDeletedNotification) ToMail() *notifications.Mail {
	return notifications.NewMail().
		Subject(i18n.T(n.User.Language, "notifications.account.deletion.completed.subject")).
		Greeting(i18n.TWithParams(n.User.Language, "notifications.greeting", n.User.GetName())).
		Line(i18n.T(n.User.Language, "notifications.account.deletion.completed.confirmation")).
		Line(i18n.T(n.User.Language, "notifications.account.deletion.completed.permanent")).
		Line(i18n.T(n.User.Language, "notifications.common.have_nice_day"))
}

// ToDB returns the AccountDeletedNotification notification in a format which can be saved in the db
func (n *AccountDeletedNotification) ToDB() interface{} {
	return nil
}

// Name returns the name of the notification
func (n *AccountDeletedNotification) Name() string {
	return "user.deleted"
}
