import { Link } from "react-router-dom";
import { useAuth } from "../../contexts/AuthContext";
import { Button } from "@base-ui/react/button";

export function Home() {
  const { user, setUserFromAuth } = useAuth();

  return (
    <>
      <header className="bg-gray-300 p-4 flex gap-3">
        {user === null ? (
          <Link to="/login">login</Link>
        ) : (
          <Button onClick={() => setUserFromAuth(null)}>logout</Button>
        )}
        <p>
          Welcome to the home page!{" "}
          {user ? `Hello, ${user.name}` : "You are not logged in."}
        </p>
      </header>
    </>
  );
}
