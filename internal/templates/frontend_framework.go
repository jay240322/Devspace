package templates

import (
	"os"
	"path/filepath"
)

func GenerateFrontendFramework(frontendPath string, srcPath string, meta ProjectMetadata) {
	var pkg, appCode, appName, mainJs string
	var mainFileName string = "main.js"

	switch meta.Frontend {
	case "Vue.Js":
		appName = "App.vue"
		pkg = `{ "name": "vue-app", "type": "module", "dependencies": { "vue": "^3.4.0" }, "devDependencies": { "vite": "^5.0.0", "@vitejs/plugin-vue": "^5.0.0" } }`
		mainJs = `import { createApp } from 'vue'
import App from './App.vue'

createApp(App).mount('#root')`
		appCode = `<template>
  <div style="padding:40px; text-align:center;">
    <h1>{{.ServiceName}} Vue App</h1>
  </div>
</template>`

	case "Svelte":
		appName = "App.svelte"
		pkg = `{ "name": "svelte-app", "type": "module", "dependencies": { "svelte": "^4.0.0" }, "devDependencies": { "vite": "^5.0.0", "@sveltejs/vite-plugin-svelte": "^3.0.0" } }`
		mainJs = `import App from './App.svelte'

const app = new App({
  target: document.getElementById('root')
})

export default app`
		appCode = `<script></script>
<div style="padding:40px; text-align:center;">
  <h1>✨ {{.ServiceName}} Svelte App</h1>
</div>`

	default: // React (Vite)
		appName = "App.jsx"
		mainFileName = "main.jsx" 
		pkg = `{ "name": "react-app", "type": "module", "dependencies": { "react": "^18.3.1", "react-dom": "^18.3.1" }, "devDependencies": { "vite": "^5.3.1", "@vitejs/plugin-react": "^4.3.1" } }`
		mainJs = `import React from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.jsx'

createRoot(document.getElementById('root')).render(<App />)`
		appCode = `import React from 'react'

export default function App() {
  return (
    <div style={{padding:'40px', textAlign:'center'}}>
      <h1> {{.ServiceName}} React App</h1>
    </div>
  )
}`
	}

	// Dynamic multi-line HTML blueprint compilation
	html := `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{.ServiceName}} Platform</title>
  </head>
  <body style="background-color: #121214; color: #e1e1e6; font-family: sans-serif;">
    <div id="root"></div>
    <script type="module" src="/src/` + mainFileName + `"></script>
  </body>
</html>`

	_ = writeTemplate(filepath.Join(frontendPath, "package.json"), pkg, meta)
	_ = writeTemplate(filepath.Join(frontendPath, "index.html"), html, meta)
	_ = os.WriteFile(filepath.Join(srcPath, mainFileName), []byte(mainJs), 0644)
	_ = writeTemplate(filepath.Join(srcPath, appName), appCode, meta)
}