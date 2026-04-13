import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.jsx";
import "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css";
import { AuthProvider } from "./context/AuthProvider.jsx";
import { BrowserRouter } from "react-router-dom";
import AuthSync from "./context/AuthSync.jsx";

createRoot(document.getElementById("root")).render(
  <StrictMode>
      <AuthProvider>
        <BrowserRouter>
          <AuthSync />
          <App />
        </BrowserRouter>
      </AuthProvider>
  </StrictMode>,
);
