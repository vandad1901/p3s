import { Authentication } from "@/pages/authn/authn";
import { Home } from "@/pages/home/home";
import { BrowserRouter, Route, Routes } from "react-router-dom";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<Authentication mode="login" />} />
        <Route path="/signup" element={<Authentication mode="register" />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
