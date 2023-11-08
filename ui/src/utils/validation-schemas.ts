import { z } from "zod";

export const emailSchema = z
  .string({
    required_error: "errors.required",
  })
  .email("errors.invalid_email")
  .max(100, "errors.invalid_email");

export const usernameSchema = z
  .string({
    required_error: "errors.required",
  })
  .min(3, "errors.username_too_short")
  .max(32, "errors.username_too_long")
  .regex(
    new RegExp("^[a-zA-Z0-9_.-]+$"),
    "errors.username_contains_invalid_chars"
  );

export const passwordSchema = z
  .string({
    required_error: "errors.required",
  })
  .min(8, "errors.password_too_short")
  .max(64, "errors.password_too_long")
  .regex(
    new RegExp("^\\P{Ll}*\\p{Ll}[\\s\\S]*$", "u"),
    "errors.password_must_contain_lowercase"
  )
  .regex(
    new RegExp("^\\P{Lu}*\\p{Lu}[\\s\\S]*$", "u"),
    "errors.password_must_contain_uppercase"
  )
  .regex(
    new RegExp("^\\P{N}*\\p{N}[\\s\\S]*$", "u"),
    "errors.password_must_contain_number"
  )
  .regex(
    new RegExp("^[\\p{L}\\p{N}]*[^\\p{L}\\p{N}][\\s\\S]*$", "u"),
    "errors.password_must_contain_special"
  );

export const requiredStringSchema = z
  .string({
    required_error: "errors.required",
  })
  .min(1, "errors.required");
