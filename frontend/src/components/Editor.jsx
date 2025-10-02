import React, {useEffect, useRef} from "react";
import Editor from "@monaco-editor/react"

export default function EditorComponent() {

    const ws = useRef(null)
    useEffect(() => {
        ws.current = new WebSocket("ws://localhost:8080/ws")
        ws.current.onopen = () => {
            console.log("Websocket is connected")
        }

        ws.current.onmessage = (event) => {
            console.log("Received message from server:", event.data);
        }

        ws.current.onclose = () => {
            console.log("WebSocket connection closed");
        };

        return () => {
            if (ws.current) {
                ws.current.close();
            }
        };

    }, [])

    const handleEditorChange = (value, event) => {
        if (ws.current && ws.current.readyState === WebSocket.OPEN) {
            ws.current.send(value);
        }
    };

    return(
        <>
        <Editor
        height="90vh"
        width="90vw"
        defaultLanguage="javascript"
        defaultValue="type here"
        theme="vs-dark"
        onChange={handleEditorChange}
        />
        </>
    )
}