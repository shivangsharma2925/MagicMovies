import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css";
import { AuthProvider } from "./context/AuthProvider.jsx";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import App from "./App.jsx";
import { WebSocketProvider } from "./context/WebSocketProvider.jsx";
import { BrowserRouter } from "react-router-dom";

const queryClient = new QueryClient();

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <WebSocketProvider>
          <BrowserRouter>
            <App />
          </BrowserRouter>
        </WebSocketProvider>
      </AuthProvider>
    </QueryClientProvider>
  </StrictMode>,
);
