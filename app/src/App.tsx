import { useState } from 'react'
import reactLogo from './assets/react.svg'
import websocketLogo from './assets/websocket.svg'
import golangLogo from './assets/golang.svg'

import './App.css'

function App() {
  const [messages, setMessages] = useState<string[]>([])

  const ws = new WebSocket("ws://localhost:8000/ws")

  ws.onopen = (event) => ws.send(JSON.stringify('hello from frontned!'))

  ws.onmessage = (event) => {
    const message = JSON.parse(event)

    // TODO: continue saving messages from server
  }

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
          <img src={websocketLogo} className="logo websocket" alt="WebSocket logo" />
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

      {
        messages.length ? (
          messages.map((message, index) => {
            return (
              <div className="message" key={index}>
                {message}
              </div>
            )
          })
        ) : <div className="message__empty">Chat is empty.</div>
      }

      <form>
        <input id="username" type="text" placeholder="Your name" />
        <textarea id="usermessage" placeholder="Here type your message"></textarea>
      </form>

      <div className="card">
        <button onClick={() => console.log('sent')}>
          Send message
        </button>
      </div>

      <p className="read-the-docs">
        Click on desired logos to learn more.
      </p>
    </div>
  )
}

export default App
