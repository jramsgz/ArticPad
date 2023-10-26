interface ErrorMap {
  [key: string]: AppError;
}

interface AppError {
  title: string;
  message: string;
  setMessage: boolean;
  setFooter: boolean;
}

const errorMap: ErrorMap = {
  /* Internal frontend-defined errors */
  NET_ERR: {
    title: "errors.net_err",
    message: "errors.net_err_msg",
    setMessage: false,
    setFooter: false,
  },
  LOG_IN_AGAIN: {
    title: "errors.unauthorized",
    message: "errors.log_in_again",
    setMessage: true,
    setFooter: false,
  },
  MISSING_TOKEN: {
    title: "errors.missing_token",
    message: "errors.missing_token_msg",
    setMessage: false,
    setFooter: false,
  },

  /* External backend-defined errors */
  bad_request: {
    title: "errors.bad_request",
    message: "errors.bad_request_msg",
    setMessage: true,
    setFooter: false,
  },
  account_not_found: {
    title: "errors.account_not_found",
    message: "errors.account_not_found_msg",
    setMessage: false,
    setFooter: false,
  },
  invalid_credentials: {
    title: "errors.invalid_credentials",
    message: "errors.invalid_credentials_msg",
    setMessage: false,
    setFooter: false,
  },
  username_length_less_than_3: {
    title: "errors.invalid_username",
    message: "errors.username_too_short",
    setMessage: false,
    setFooter: false,
  },
  username_length_more_than_32: {
    title: "errors.invalid_username",
    message: "errors.username_too_long",
    setMessage: false,
    setFooter: false,
  },
  username_contains_invalid_characters: {
    title: "errors.invalid_username",
    message: "errors.username_contains_invalid_chars",
    setMessage: false,
    setFooter: false,
  },
  password_length_less_than_8: {
    title: "errors.invalid_password",
    message: "errors.password_too_short",
    setMessage: false,
    setFooter: false,
  },
  password_length_more_than_64: {
    title: "errors.invalid_password",
    message: "errors.password_too_long",
    setMessage: false,
    setFooter: false,
  },
  password_similarity: {
    title: "errors.invalid_password",
    message: "errors.password_too_similar",
    setMessage: false,
    setFooter: false,
  },
  invalid_email: {
    title: "errors.invalid_email",
    message: "errors.invalid_email_msg",
    setMessage: false,
    setFooter: false,
  },
  email_already_exists: {
    title: "errors.invalid_email",
    message: "errors.email_in_use",
    setMessage: false,
    setFooter: false,
  },
  username_already_exists: {
    title: "errors.invalid_username",
    message: "errors.username_in_use",
    setMessage: false,
    setFooter: false,
  },
  password_strength: {
    title: "errors.invalid_password",
    message: "errors.password_too_weak",
    setMessage: false,
    setFooter: false,
  },
  email_not_verified: {
    title: "errors.email_not_verified",
    message: "errors.email_not_verified_msg",
    setMessage: false,
    setFooter: false,
  },
  email_already_verified: {
    title: "errors.already_verified",
    message: "",
    setMessage: false,
    setFooter: false,
  },
  cannot_send_verification_email: {
    title: "errors.error_sending_email",
    message: "errors.cannot_send_verification_email",
    setMessage: false,
    setFooter: false,
  },
  mail_not_enabled: {
    title: "errors.error_sending_email",
    message: "errors.mail_not_enabled",
    setMessage: false,
    setFooter: false,
  },
  cannot_send_password_reset_email: {
    title: "errors.error_sending_email",
    message: "errors.cannot_send_password_reset_email",
    setMessage: false,
    setFooter: false,
  },
  password_reset_token_expired: {
    title: "errors.token_expired",
    message: "errors.password_reset_token_expired",
    setMessage: false,
    setFooter: false,
  },
  invalid_jwt: {
    title: "errors.invalid_token",
    message: "errors.invalid_jwt",
    setMessage: false,
    setFooter: false,
  },
  invalid_password_reset_token: {
    title: "errors.invalid_token",
    message: "errors.invalid_password_reset_token",
    setMessage: false,
    setFooter: false,
  },
  invalid_verification_token: {
    title: "errors.invalid_token",
    message: "errors.invalid_verification_token",
    setMessage: false,
    setFooter: false,
  },
  unknown_error: {
    title: "errors.unknown",
    message: "errors.unknown",
    setMessage: true,
    setFooter: true,
  },
};

export default errorMap;
