import { createContext, useEffect, useRef } from "react";

const WebSocketContext = createContext(null);

export const WebSocketProvider = ({ children }) => {

  const socket = useRef(null);

  useEffect(() => {

    socket.current = new WebSocket(
      import.meta.env.VITE_WS_URL ||
      "ws://localhost:8080/ws"
    );

    socket.current.onopen = () => {
      console.log("WebSocket Connected");
    };

    socket.current.onerror = (err) => {
      console.log("WebSocket Error:", err);
    };

    socket.current.onclose = () => {
      console.log("WebSocket Closed");
    };

    return () => {
      socket.current?.close();
    };

  }, []);

  return (
    <WebSocketContext.Provider value={socket}>
      {children}
    </WebSocketContext.Provider>
  );
};

export default WebSocketContext;