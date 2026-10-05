import { Authentication } from "@/pages/authn/authn";
import { Home } from "@/pages/home/home";
import { BrowserRouter, Route, Routes } from "react-router-dom";

function App() {
  return (
    <BrowserRouter>
      <Routes>
          <Route
            path="/editor"
            element={<Editor />}
          />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
