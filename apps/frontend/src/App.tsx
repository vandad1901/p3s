import { Authentication } from "@/pages/authn/authn";
import { Home } from "@/pages/home/home";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Layout } from "./layout";
import { Editor } from "./pages/editor/editor";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route
            path="/"
            element={<Home />}
          />
          <Route
            path="/editor"
            element={<Editor />}
          />
        </Route>
        <Route
          path="/login"
          element={<Authentication mode="login" />}
        />
        <Route
          path="/signup"
          element={<Authentication mode="register" />}
        />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
