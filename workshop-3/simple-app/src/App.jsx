import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import Payment from './components/Payment'

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
      <Payment />
    </>
  )
}

export default App
