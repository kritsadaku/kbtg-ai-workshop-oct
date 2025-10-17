import { useState } from 'react';
import Payment from './components/Payment';

const LockIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8 mx-auto mb-4 text-[var(--color-web-green-500)]" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
    <path strokeLinecap="round" strokeLinejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
  </svg>
);

const LightningIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8 mx-auto mb-4 text-[var(--color-web-green-500)]" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
    <path strokeLinecap="round" strokeLinejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
  </svg>
);

const MobileIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8 mx-auto mb-4 text-[var(--color-web-green-500)]" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
    <path strokeLinecap="round" strokeLinejoin="round" d="M12 18h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
  </svg>
);

const ViteLogo = () => (
    <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8" viewBox="0 0 256 257">
        <path fill="#646CFF" d="M255.2 63.6L132.5 2.1c-2.6-1.3-5.6-1.3-8.2 0L1.3 63.6C.5 64 0 64.9 0 65.8v124.9c0 .9.5 1.8 1.3 2.2l123 61.5c1.3.6 2.6.9 4.1.9s2.8-.3 4.1-.9l123-61.5c.8-.4 1.3-1.3 1.3-2.2V65.8c0-.9-.5-1.8-1.3-2.2zM128.4 243.2L12.2 185.8V72.4l116.2 57.4v113.4zM38.3 59.9l89.9-44.2l90.4 44.7l-89.9 44.2L38.3 59.9zm204.6 125.9L132.2 243.2V130.3l110.7-54.9v109.5z"/>
    </svg>
)

const ReactLogo = () => (
    <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" viewBox="0 0 24 24">
        <circle cx="12" cy="12" r="2" fill="currentColor"/>
        <g>
            <ellipse cx="12" cy="12" rx="11" ry="4" stroke="currentColor" strokeWidth="1" fill="none"/>
            <ellipse cx="12" cy="12" rx="11" ry="4" transform="rotate(60 12 12)" stroke="currentColor" strokeWidth="1" fill="none"/>
            <ellipse cx="12" cy="12" rx="11" ry="4" transform="rotate(120 12 12)" stroke="currentColor" strokeWidth="1" fill="none"/>
        </g>
    </svg>
)

function App() {
  const [count, setCount] = useState(0);

  return (
    <div className="bg-[var(--color-neutral-50)] text-[var(--color-neutral-700)]">
      {/* Header */}
      <header className="bg-[var(--color-web-green-600)] text-white py-[var(--spacing-16)] px-[var(--spacing-8)] text-center">
        <div className="container mx-auto">
          <div className="flex justify-center items-center gap-4 mb-4">
            <ViteLogo />
            <h1 className="text-[var(--text-5xl)] font-bold">PayFlow Demo</h1>
          </div>
          <p className="text-[var(--text-lg)] text-[var(--color-web-green-100)]">
            Experience seamless payment processing with modern React components and beautiful UI design
          </p>
        </div>
      </header>

      <main className="container mx-auto py-[var(--spacing-16)] px-[var(--spacing-8)]">
        {/* Modern Payment Solutions Section */}
        <section className="text-center mb-[var(--spacing-24)]">
          <h2 className="text-[var(--text-4xl)] font-bold mb-[var(--spacing-4)]">Modern Payment Solutions</h2>
          <p className="text-[var(--color-neutral-500)] mb-[var(--spacing-12)] max-w-3xl mx-auto">
            Built with React and Vite for lightning-fast performance. Secure, responsive, and user-friendly payment forms that adapt to any device.
          </p>
          <div className="grid md:grid-cols-3 gap-[var(--spacing-8)]">
            <div className="bg-white p-[var(--spacing-8)] rounded-[var(--radius-lg)] shadow-[var(--shadow)]">
              <LockIcon />
              <h3 className="text-[var(--text-xl)] font-semibold mb-[var(--spacing-2)]">Secure Processing</h3>
              <p className="text-[var(--color-neutral-500)]">End-to-end encryption ensures your payment data is always protected and secure.</p>
            </div>
            <div className="bg-white p-[var(--spacing-8)] rounded-[var(--radius-lg)] shadow-[var(--shadow)]">
              <LightningIcon />
              <h3 className="text-[var(--text-xl)] font-semibold mb-[var(--spacing-2)]">Lightning Fast</h3>
              <p className="text-[var(--color-neutral-500)]">Powered by Vite and React for instant load times and smooth interactions.</p>
            </div>
            <div className="bg-white p-[var(--spacing-8)] rounded-[var(--radius-lg)] shadow-[var(--shadow)]">
              <MobileIcon />
              <h3 className="text-[var(--text-xl)] font-semibold mb-[var(--spacing-2)]">Mobile Ready</h3>
              <p className="text-[var(--color-neutral-500)]">Fully responsive design that works perfectly on all devices and screen sizes.</p>
            </div>
          </div>
        </section>

        {/* Interactive Demo Section */}
        <section className="text-center mb-[var(--spacing-24)]">
            <div className="bg-white p-[var(--spacing-8)] rounded-[var(--radius-lg)] shadow-[var(--shadow)] max-w-md mx-auto">
                <h3 className="text-[var(--text-xl)] font-semibold mb-[var(--spacing-4)]">Interactive Demo</h3>
                <button
                    onClick={() => setCount((count) => count + 1)}
                    className="bg-[var(--color-web-green-500)] text-white font-bold py-2 px-4 rounded-[var(--radius-md)] hover:bg-[var(--color-web-green-600)] transition-colors"
                >
                    Click Counter: {count}
                </button>
                <p className="text-[var(--text-sm)] text-[var(--color-neutral-400)] mt-2">Edit <code>src/App.jsx</code> and save to test Hot Module Replacement</p>
            </div>
        </section>

        {/* Payment Form Section */}
        <section className="bg-[var(--color-neutral-800)] py-[var(--spacing-24)] px-[var(--spacing-8)] text-center rounded-[var(--radius-lg)]">
          <h2 className="text-[var(--text-4xl)] font-bold text-white mb-[var(--spacing-4)]">Try Our Payment Form</h2>
          <p className="text-[var(--color-neutral-300)] mb-[var(--spacing-12)]">
            Experience our beautifully designed payment form with real-time validation, secure input handling, and smooth animations.
          </p>
          <Payment />
           <div className="text-center mt-8 text-sm text-gray-400 flex justify-center items-center gap-4">
                <span>✨ This is a demo form - no real payments will be processed</span>
                <span>🔒 Secure & Encrypted</span>
                <span>📱 Mobile Optimized</span>
                <span>⚡ Real-time Validation</span>
            </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="bg-[var(--color-neutral-800)] text-[var(--color-neutral-300)] py-[var(--spacing-12)] px-[var(--spacing-8)] mt-[var(--spacing-16)]">
        <div className="container mx-auto grid md:grid-cols-3 gap-[var(--spacing-8)]">
          <div>
            <h4 className="font-bold text-white mb-2">PayFlow Demo</h4>
            <p className="text-sm">A modern payment form built with React, Vite, and Tailwind CSS, showcasing best practices in UI/UX design.</p>
          </div>
          <div>
            <h4 className="font-bold text-white mb-2">Technologies</h4>
            <ul className="text-sm space-y-1">
              <li className="flex items-center gap-2"><ReactLogo /> React 19</li>
              <li className="flex items-center gap-2"><ViteLogo /> Vite</li>
              <li>+ Tailwind CSS</li>
              <li>+ Modern JavaScript</li>
            </ul>
          </div>
          <div>
            <h4 className="font-bold text-white mb-2">Learn More</h4>
            <ul className="text-sm space-y-1">
              <li><a href="#" className="hover:text-white">React Docs</a></li>
              <li><a href="#" className="hover:text-white">Vite Guide</a></li>
            </ul>
          </div>
        </div>
        <div className="text-center text-xs text-[var(--color-neutral-400)] mt-[var(--spacing-8)] border-t border-[var(--color-neutral-700)] pt-[var(--spacing-8)]">
          Built with ❤️ using modern web technologies
        </div>
      </footer>
    </div>
  );
}

export default App;

