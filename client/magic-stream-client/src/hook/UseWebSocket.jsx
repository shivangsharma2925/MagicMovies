import { useContext } from "react"
import WebSocketContext from "../context/WebSocketProvider"

const UseWebSocket = () => {
    return useContext(WebSocketContext);
}

export default UseWebSocket;