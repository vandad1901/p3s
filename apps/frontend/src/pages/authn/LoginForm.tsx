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

const formSchema = z.object({
  username: z.string().min(5, "Username must be at least 5 characters."),
  password: z.string().min(12, "Password must be at least 12 characters."),
});

export function LoginForm() {
  const { setUserFromAuth } = useAuth();
  const [serverErrors, setServerErrors] = useState<{
    username?: string;
    password?: string;
  }>({});

  const form = useForm({
    defaultValues: {
      username: "",
      password: "",
    },
    validators: {
      onSubmit: formSchema,
    },
    onSubmit: async ({ value }) => {
      setServerErrors({});

      let res = await authnService.Login({
        username: value.username,
        password: value.password,
      });
      if (!res.ok) {
        if (res.message === "authn.invalidAuth") {
          setServerErrors({
            password: "Incorrect username or password.",
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
          id="login-form"
          onSubmit={(e) => {
            e.preventDefault();
            form.handleSubmit();
          }}
        >
          <FieldGroup>
            <form.Field
              name="username"
              children={(field) => {
                const isInvalid =
                  !field.state.meta.isValid || !!serverErrors.username;
                return (
                  <Field data-invalid={isInvalid}>
                    <FieldLabel htmlFor={field.name}>Username</FieldLabel>
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
                const isInvalid =
                  !field.state.meta.isValid || !!serverErrors.password;
                return (
                  <Field data-invalid={isInvalid}>
                    <FieldLabel htmlFor={field.name}>password</FieldLabel>
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
                      <FieldError
                        errors={
                          serverErrors.password
                            ? [{ message: serverErrors.password }]
                            : field.state.meta.errors
                        }
                      />
                    )}
                  </Field>
                );
              }}
            />
            <Field>
              <Button type="submit" form="login-form" size={"lg"}>
                Login
              </Button>
            </Field>
            <Separator />
            <FieldDescription className="text-center">
              Don&apos;t have an account? <Link to="/signup">Sign up</Link>
            </FieldDescription>
          </FieldGroup>
        </form>
      </CardContent>
    </Card>
  );
}
