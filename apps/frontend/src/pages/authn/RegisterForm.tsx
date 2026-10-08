"use client";

import { useForm } from "@tanstack/react-form";
import * as z from "zod";

import { authnService } from "@/api/authn.service";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { useAuth } from "@/contexts/AuthContext";
import { useState } from "react";
import { Link } from "react-router";

const formSchema = z
  .object({
    email: z.string().email("Please enter a valid email address."),
    username: z.string().min(5, "Username must be at least 5 characters."),
    password: z.string().min(12, "Password must be at least 12 characters."),
    passwordConfirm: z.string(),
  })
  .refine((data) => data.password === data.passwordConfirm, {
    message: "Passwords do not match.",
    path: ["passwordConfirm"],
    params: { code: "passwordsDoNotMatch" },
  });

export function RegisterForm() {
  const { setUserFromAuth } = useAuth();
  const [serverErrors, setServerErrors] = useState<{
    email?: string;
    username?: string;
  }>({});

  const form = useForm({
    defaultValues: {
      email: "",
      username: "",
      password: "",
      passwordConfirm: "",
    },
    validators: {
      onSubmit: formSchema,
    },
    onSubmit: async ({ value }) => {
      setServerErrors({});

      let res = await authnService.Register({
        user: {
          email: value.email,
          username: value.username,
          id: "0",
        },
        password: value.password,
      });
      if (!res.ok) {
        if (res.message === "identity.DuplicateEmail") {
          setServerErrors({
            email: "This email is already registered.",
          });
        }

        if (res.message === "identity.DuplicateUsername") {
          setServerErrors({
            username: "This username is already registered.",
          });
        }

        return;
      }

      setUserFromAuth({
        sessionId: res.session?.sessionId.toString() ?? "",
        name: value.username,

        jwt: res.session!.jwt,
      });
    },
  });

  return (
    <Card className="w-full sm:max-w-md">
      <CardContent>
        <form
          id="register-form"
          onSubmit={(e) => {
            e.preventDefault();
            form.handleSubmit();
          }}
        >
          <FieldGroup>
            <form.Field
              name="email"
              children={(field) => {
                const isInvalid =
                  !field.state.meta.isValid || !!serverErrors.email;
                return (
                  <Field data-invalid={isInvalid}>
                    <FieldLabel htmlFor={field.name} required>
                      Email
                    </FieldLabel>
                    <Input
                      required
                      id={field.name}
                      name={field.name}
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      aria-invalid={isInvalid}
                      autoComplete="off"
                    />
                    {isInvalid && (
                      <FieldError
                        errors={
                          serverErrors.email
                            ? [{ message: serverErrors.email }]
                            : field.state.meta.errors
                        }
                      />
                    )}
                  </Field>
                );
              }}
            />

            <form.Field
              name="username"
              children={(field) => {
                const isInvalid =
                  !field.state.meta.isValid || !!serverErrors.username;
                return (
                  <Field data-invalid={isInvalid}>
                    <FieldLabel htmlFor={field.name} required>
                      Username
                    </FieldLabel>
                    <Input
                      required
                      id={field.name}
                      name={field.name}
                      value={field.state.value}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      aria-invalid={isInvalid}
                      autoComplete="off"
                    />
                    {isInvalid && (
                      <FieldError
                        errors={
                          serverErrors.username
                            ? [{ message: serverErrors.username }]
                            : field.state.meta.errors
                        }
                      />
                    )}
                  </Field>
                );
              }}
            />

            <form.Field
              name="password"
              children={(field) => {
                const isInvalid = !field.state.meta.isValid;
                return (
                  <Field data-invalid={isInvalid}>
                    <FieldLabel htmlFor={field.name} required>
                      Password
                    </FieldLabel>
                    <Input
                      required
                      id={field.name}
                      name={field.name}
                      value={field.state.value}
                      type="password"
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      aria-invalid={isInvalid}
                    />
                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                );
              }}
            />

            <form.Field
              name="passwordConfirm"
              children={(field) => {
                const isInvalid = !field.state.meta.isValid;
                return (
                  <Field data-invalid={isInvalid}>
                    <FieldLabel htmlFor={field.name} required>
                      Confirm Your Password
                    </FieldLabel>
                    <Input
                      required
                      id={field.name}
                      name={field.name}
                      value={field.state.value}
                      type="password"
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      aria-invalid={isInvalid}
                    />
                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                );
              }}
            />

            <Field>
              <Button type="submit" form="register-form" size={"lg"}>
                Sign Up
              </Button>
            </Field>
            <Separator />
            <FieldDescription className="text-center">
              Already have an account?
              <Link to="/login">Click here to log in</Link>
            </FieldDescription>
          </FieldGroup>
        </form>
      </CardContent>
    </Card>
  );
}
