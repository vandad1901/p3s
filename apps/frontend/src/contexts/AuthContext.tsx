import { authnService } from "@/api/authn.service";
import { setAccessToken } from "@/api/client";
import { jwtDecode } from "jwt-decode";
import {
  createContext,
  useContext,
  useEffect,
  useRef,
  useState,
  type ReactNode,
} from "react";

type UIUser = {
  userId: string;
  sessionId: string;
  name: string;

  jwt: string;
};

type AuthContextValue = {
  user: UIUser | null;
  setUserFromAuth: (user: Omit<UIUser, "userId"> | null) => void;
};

export const AuthContext = createContext<AuthContextValue | null>(null);

async function refreshJWT(uiUser: Omit<UIUser, "jwt">): Promise<string> {
  let res = await authnService.RefreshJWT({
    userId: parseInt(uiUser.userId),
    sessionId: parseInt(uiUser.sessionId),
  });

  if (!res.ok) {
    console.error("Failed to refresh JWT", res.code, res.message);
    throw new Error("Failed to refresh JWT");
  }

  return res.jwt;
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<UIUser | null>(null);
  const userRef = useRef<UIUser | null>(null);

  const setUserFromAuth = (authUser: Omit<UIUser, "userId"> | null) => {
    if (!authUser) {
      userRef.current = null;
      setUser(null);
      setAccessToken(null);
      localStorage.removeItem("user");
      return;
    }

    const decoded = jwtDecode(authUser.jwt);

    const user: UIUser = {
      userId: (decoded as { sub?: string }).sub ?? "",
      ...authUser,
    };

    userRef.current = user;
    setUser(user);
    setAccessToken(user.jwt);

    const { jwt, ...userWithoutToken } = user;
    localStorage.setItem("user", JSON.stringify(userWithoutToken));
  };

  const setUserFromRefresh = (user: Omit<UIUser, "jwt">, newJWT: string) => {
    const updatedUser = { ...user, jwt: newJWT };

    userRef.current = updatedUser;
    setUser(updatedUser);
    setAccessToken(newJWT);
  };

  useEffect(() => {
    const storedUser = localStorage.getItem("user");

    if (storedUser) {
      let user = JSON.parse(storedUser) as Omit<UIUser, "jwt">;

      refreshJWT(user)
        .then((newJWT) => {
          setUserFromRefresh(user, newJWT);
        })
        .catch((error) => {
          console.error("Failed to restore session", error);
        });
    }

    const tokenRefresher = setInterval(async () => {
      const currentUser = userRef.current;

      if (!currentUser) return;

      try {
        const newJWT = await refreshJWT(currentUser);
        setUserFromRefresh(currentUser, newJWT);
      } catch (error) {
        console.error("Failed to refresh JWT", error);
      }
    }, 10*60*1000 );

    return () => {
      clearInterval(tokenRefresher);
    };
  }, []);

  return (
    <AuthContext.Provider value={{ user, setUserFromAuth }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }

  return context;
}
