import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import App from './App.jsx'
import faviconUrl from '../assets/favicon.ico'

const faviconLink =
  document.querySelector("link[rel*='icon']") ??
  document.head.appendChild(document.createElement('link'))

faviconLink.rel = 'icon'
faviconLink.type = 'image/x-icon'
faviconLink.href = faviconUrl

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </StrictMode>
)
