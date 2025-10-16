import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import Payment from './components/Payment'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="min-h-screen bg-neutral-50">
      {/* Header Section - Green Background */}
      <header className="bg-web-green-600 text-white px-6 py-16">
        <div className="max-w-6xl mx-auto text-center">
          <div className="flex justify-center items-center gap-4 mb-8">
            <div className="w-12 h-12 bg-white bg-opacity-20 rounded-lg flex items-center justify-center">
              <img src={viteLogo} className="w-8 h-8" alt="Vite logo" />
            </div>
            <div className="w-12 h-12 bg-white bg-opacity-20 rounded-lg flex items-center justify-center">
              <img src={reactLogo} className="w-8 h-8" alt="React logo" />
            </div>
          </div>
          <h1 className="text-5xl font-bold mb-4">PayFlow Demo</h1>
          <p className="text-xl text-white/90 max-w-2xl mx-auto">
            Experience seamless payment processing with modern React components and beautiful UI design
          </p>
        </div>
      </header>

      {/* Modern Payment Solutions Section */}
      <section className="py-20 px-6 bg-neutral-50">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-neutral-900 mb-4">Modern Payment Solutions</h2>
            <p className="text-lg text-neutral-600 max-w-3xl mx-auto">
              Built with React and Vite for lightning-fast performance. Secure, responsive, and user-friendly payment forms that adapt to any device.
            </p>
          </div>

          {/* Feature Cards */}
          <div className="grid md:grid-cols-3 gap-8 mb-16">
            {/* Secure Processing */}
            <div className="bg-white p-8 rounded-lg shadow-sm text-center">
              <div className="w-16 h-16 bg-web-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
                <svg className="w-8 h-8 text-web-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
              <h3 className="text-xl font-semibold text-neutral-900 mb-3">Secure Processing</h3>
              <p className="text-neutral-600">
                End-to-end encryption ensures your payment data is always protected and secure.
              </p>
            </div>

            {/* Lightning Fast */}
            <div className="bg-white p-8 rounded-lg shadow-sm text-center">
              <div className="w-16 h-16 bg-web-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
                <svg className="w-8 h-8 text-web-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
              <h3 className="text-xl font-semibold text-neutral-900 mb-3">Lightning Fast</h3>
              <p className="text-neutral-600">
                Powered by Vite and React for instant load times and smooth interactions.
              </p>
            </div>

            {/* Mobile Ready */}
            <div className="bg-white p-8 rounded-lg shadow-sm text-center">
              <div className="w-16 h-16 bg-web-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
                <svg className="w-8 h-8 text-web-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
                </svg>
              </div>
              <h3 className="text-xl font-semibold text-neutral-900 mb-3">Mobile Ready</h3>
              <p className="text-neutral-600">
                Fully responsive design that works perfectly on all devices and screen sizes.
              </p>
            </div>
          </div>

          {/* Interactive Demo */}
          <div className="bg-white rounded-lg shadow-sm p-12 text-center">
            <h3 className="text-2xl font-semibold text-neutral-900 mb-6">Interactive Demo</h3>
            <button 
              className="bg-web-green-500 hover:bg-web-green-600 text-white font-semibold py-4 px-8 rounded-lg transition-colors duration-200 shadow-sm text-lg"
              onClick={() => setCount((count) => count + 1)}
            >
              Click Counter: {count}
            </button>
            <p className="text-neutral-500 text-sm mt-6">
              Edit <code className="bg-neutral-100 text-neutral-700 px-2 py-1 rounded-sm font-mono text-xs">src/App.jsx</code> and save to test HMR
            </p>
          </div>
        </div>
      </section>

      {/* Payment Form Section */}
      <section className="py-20 px-6 bg-neutral-800">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-white mb-4">Try Our Payment Form</h2>
            <p className="text-xl text-neutral-300 max-w-3xl mx-auto">
              Experience our beautifully designed payment form with real-time validation, secure input handling, and smooth animations.
            </p>
          </div>
          <div className="flex justify-center">
            <Payment />
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-neutral-900 text-neutral-300 py-12 px-6">
        <div className="max-w-6xl mx-auto">
          <div className="grid md:grid-cols-4 gap-8 mb-8">
            <div>
              <h4 className="text-white font-semibold mb-4">PayFlow Demo</h4>
              <p className="text-sm text-neutral-400">
                A modern payment form built with React, Vite, and Tailwind CSS. Showcasing best practices for UI/UX in payment processing.
              </p>
            </div>
            <div>
              <h4 className="text-white font-semibold mb-4">Technologies</h4>
              <ul className="space-y-2 text-sm text-neutral-400">
                <li>React 18</li>
                <li>Vite</li>
                <li>Tailwind CSS</li>
                <li>Modern JavaScript</li>
              </ul>
            </div>
            <div>
              <h4 className="text-white font-semibold mb-4">Learn More</h4>
              <ul className="space-y-2 text-sm text-neutral-400">
                <li><a href="https://react.dev" className="hover:text-white transition-colors">React Docs</a></li>
                <li><a href="https://vitejs.dev" className="hover:text-white transition-colors">Vite Guide</a></li>
                <li><a href="https://tailwindcss.com" className="hover:text-white transition-colors">Tailwind CSS</a></li>
                <li>Security Guidelines</li>
              </ul>
            </div>
            <div>
              <h4 className="text-white font-semibold mb-4"></h4>
              <p className="text-xs text-neutral-500">
                Built with ❤️ using modern web technologies
              </p>
            </div>
          </div>
          <div className="border-t border-neutral-800 pt-8 text-center">
            <p className="text-sm text-neutral-500">
              © 2025 PayFlow Demo. This is a demo form built for showcase purposes.
            </p>
          </div>
        </div>
      </footer>
    </div>
  )
}

export default App
