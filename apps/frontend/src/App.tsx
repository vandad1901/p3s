import { Home } from "@/pages/home/home";
import { Suspense, lazy } from "react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Layout } from "./layout";
import { PostPage } from "./pages/editor/post";

const Authentication = lazy(() =>
  import("@/pages/authn/authn").then((m) => ({ default: m.Authentication })),
);

function App() {
  return (
    <BrowserRouter>
      <Suspense fallback={null}>
        <Routes>
          <Route element={<Layout />}>
            <Route
              path="/"
              element={<Home />}
            />
            <Route
              path="/editor/:slug"
              element={<PostPage editMode={true} />}
            />
            <Route
              path="/editor/"
              element={<PostPage editMode={true} />}
            />
            <Route
              path="/post/:slug"
              element={<PostPage editMode={false} />}
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
      </Suspense>
    </BrowserRouter>
  );
}

export default App;
