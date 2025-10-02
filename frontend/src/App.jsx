// frontend/src/App.jsx
import React, { useState, useRef, useEffect } from 'react';
import EditorComponent from './components/Editor';
import './App.css';

function App() {
  const [code, setCode] = useState("// Type your code here");
  const [output, setOutput] = useState(""); // New state for the output
  const ws = useRef(null);

  useEffect(() => {
    ws.current = new WebSocket("ws://localhost:8080/ws");
    ws.current.onopen = () => console.log("WebSocket connection established");
    ws.current.onclose = () => console.log("WebSocket connection closed");

    ws.current.onmessage = (event) => {
      const message = JSON.parse(event.data);
      switch (message.type) {
        case "code_update":
          setCode(message.payload);
          break;
        case "run_output":
          setOutput(message.payload);
          break;
        default:
          break;
      }
    };

    return () => {
      if (ws.current) ws.current.close();
    };
  }, []);

  const handleEditorChange = (value) => {
    setCode(value);
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      const message = {
        type: "code_update",
        payload: value,
      };
      ws.current.send(JSON.stringify(message));
    }
  };

  const runCode = () => {
    setOutput("Executing code...");
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      const message = {
        type: "run_code",
        payload: code,
      };
      ws.current.send(JSON.stringify(message));
    }
  };

  return (
    <div className="app-container">
      <div className="editor-container">
        <button onClick={runCode} className="run-button">Run</button>
        <EditorComponent code={code} onChange={handleEditorChange} />
      </div>
      <div className="output-container">
        <h2>Output</h2>
        <pre className="output-content">{output}</pre>
      </div>
    </div>
  );
}

export default App;