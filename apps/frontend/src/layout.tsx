import { Button, buttonVariants } from "@/components/ui/button";
import { useAuth } from "@/contexts/AuthContext";
import { cn } from "cn";
import { LogOut, PenLine } from "lucide-react";
import { Link, NavLink, Outlet, useNavigate } from "react-router-dom";

const YEAR = new Date().getFullYear();

export function Layout() {
  const { user, setUserFromAuth } = useAuth();
  const navigate = useNavigate();

  function logout() {
    setUserFromAuth(null);
    navigate("/");
  }

  return (
    <div className="flex min-h-svh flex-col">
      <header className="border-b">
        <nav
          aria-label="Main"
          className="mx-auto flex h-12 max-w-5xl items-center gap-2 px-4"
        >
          <Link
            to="/"
            className="mr-auto text-sm font-semibold tracking-tight"
          >
            p3s
          </Link>

          {user ? (
            <>
              <NavLink
                to="/editor"
                className={({ isActive }) =>
                  cn(buttonVariants({ variant: isActive ? "secondary" : "ghost", size: "lg" }))
                }
              >
                <PenLine /> Write
              </NavLink>

              <div className="flex items-center gap-2 ps-2">
                <span
                  aria-hidden
                  className="bg-primary text-primary-foreground grid size-7 place-items-center rounded-full text-xs font-medium"
                >
                  {Array.from(user.name)[0]?.toLocaleUpperCase()}
                </span>
                <span
                  dir="auto"
                  className="hidden max-w-32 truncate text-xs sm:inline"
                >
                  {user.name}
                </span>
                <Button
                  variant="ghost"
                  size="icon"
                  aria-label="Log out"
                  onClick={logout}
                >
                  <LogOut />
                </Button>
              </div>
            </>
          ) : (
            <>
              <Link
                to="/login"
                className={buttonVariants({ variant: "ghost", size: "lg" })}
              >
                Log in
              </Link>
              <Link
                to="/signup"
                className={buttonVariants({ size: "lg" })}
              >
                Sign up
              </Link>
            </>
          )}
        </nav>
      </header>

      <main className="flex flex-1 flex-col">
        <Outlet />
      </main>
      <footer className="border-t">
        <div className="text-muted-foreground mx-auto flex h-12 max-w-5xl items-center justify-center gap-2 px-4 text-sm">
          <span>© {YEAR} p3s.vandaddelavari.ir - All Rights Reserved</span>
        </div>
      </footer>
    </div>
  );
}
