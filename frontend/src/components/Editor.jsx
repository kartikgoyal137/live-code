import React from 'react';
import Editor from '@monaco-editor/react';

function EditorComponent({ code, onChange }) {
  return (
    <Editor
      height="90vh"
      width="100%" 
      language="javascript"
      value={code}
      theme="vs-dark"
      onChange={onChange}
    />
  );
}

export default EditorComponent;