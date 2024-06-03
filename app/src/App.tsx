import { useState, useRef, useEffect } from "react";
import reactLogo from "./assets/react.svg";
import websocketLogo from "./assets/websocket.svg";
import golangLogo from "./assets/golang.svg";

import useWebSocket, { ReadyState } from "react-use-websocket";

import "./App.css";

const WebSocketURL = "ws://localhost:8000/ws";

// following https://stackoverflow.com/questions/62768520/reconnecting-web-socket-using-react-hooks

function App() {
  const [chatMessages, setChatMessages] = useState<MessageEvent<any>[]>([]);
  const [fieldMessage, setFieldMessage] = useState("");
  const [fieldUser, setFieldUser] = useState("");

  const { sendMessage, lastMessage, readyState } = useWebSocket(WebSocketURL, {
    share: false,
    shouldReconnect: () => true,
  });

  const connectionStatus = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Open",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Closed",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];

  // run when the connection state changes
  useEffect(() => {
    console.log("connection state changed", connectionStatus);

    if (readyState === ReadyState.OPEN) {
      sendMessage("Some message");
    }
  });

  // ws message received
  useEffect(() => {
    console.log(`got message: ${lastMessage}`);
    if (lastMessage !== null) {
      setChatMessages((prev) => prev.concat(lastMessage));
    }
  }, [lastMessage]);

  const sendChatMessage = () => {
    const str = JSON.stringify(`${fieldUser}: ${fieldMessage}`);
  };

  const handleUsernameChange = (event) => {
    const val = event.target.value;
    setFieldUser(val);
  };

  const handleFieldMessageChange = (event) => {
    const val = event.target.value;
    setFieldMessage(val);
  };

  return (
    <div className="App">
      <div>
        <a href="https://go.dev" target="_blank">
          <img src={golangLogo} className="logo golang" alt="React logo" />
        </a>
        <a href="https://reactjs.org" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
        <a href="https://www.rfc-editor.org/rfc/rfc6455" target="_blank">
          <img
            src={websocketLogo}
            className="logo websocket"
            alt="WebSocket logo"
          />
        </a>
        <a href="https://vitejs.dev" target="_blank">
          <img src="/vite.svg" className="logo" alt="Vite logo" />
        </a>
      </div>
      <h1>Chat with Go, React, and WebSockets</h1>

      <p>
        This is a simple chat application that allow users share messages
        beteween them.
      </p>

      <h2>connection state: {connectionStatus}</h2>

      {chatMessages.length ? (
        chatMessages.map((message, index) => {
          return (
            <div className="message" key={index}>
              {message.data}
            </div>
          );
        })
      ) : (
        <div className="message__empty">Chat is empty.</div>
      )}

      <form>
        <input
          id="username"
          type="text"
          placeholder="Your name"
          onChange={handleUsernameChange}
        />
        <textarea
          id="usermessage"
          placeholder="Here type your message"
          onChange={handleFieldMessageChange}
        ></textarea>
      </form>

      <div className="card">
        <button
          onClick={sendChatMessage}
          disabled={readyState !== ReadyState.OPEN}
        >
          Send message
        </button>
      </div>

      <p className="read-the-docs">Click on desired logos to learn more.</p>
    </div>
  );
}

export default App;
