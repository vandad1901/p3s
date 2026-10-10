import { useAuth } from "@/contexts/AuthContext";
import { Navigate } from "react-router-dom";
import { LoginForm } from "./LoginForm";
import { RegisterForm } from "./RegisterForm";

const RegistrationClosed: boolean = import.meta.env.VITE_CLOSED_REGISTRATION === "true";

export function Authentication({ mode }: { mode: "login" | "register" }) {
  const { user } = useAuth();

  if (user !== null) {
    return <Navigate to="/"></Navigate>;
  }

  if (mode === "register" && RegistrationClosed) {
    return (
      <Navigate
        to="/login"
        replace
      ></Navigate>
    );
  }

  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-1">
      <div className="w-full max-w-sm">{mode === "login" ? <LoginForm /> : <RegisterForm />}</div>
    </div>
  );
}
