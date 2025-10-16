import { useState } from 'react'

const Payment = () => {
  const [formData, setFormData] = useState({
    nameOnCard: 'John Doe',
    cardNumber: '1234 5678 9012 3456',
    cvv: '123',
    month: 'MM',
    year: 'YYYY',
    comments: '',
    sameAsShipping: true
  })

  const handleSubmit = (e) => {
    e.preventDefault()
    console.log('Payment submitted:', { ...formData })
  }

  const handleCancel = () => {
    console.log('Payment cancelled')
    // ถ้าต้องการ reset ค่า ให้ทำแบบนี้:
    // setFormData({
    //   nameOnCard: '',
    //   cardNumber: '',
    //   cvv: '',
    //   month: 'MM',
    //   year: 'YYYY',
    //   comments: '',
    //   sameAsShipping: true
    // })
  }

  const handleInputChange = (e) => {
    const { name, value, type, checked } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }))
  }

  return (
    <div className="bg-neutral-900 text-white p-8 rounded-lg max-w-md mx-auto shadow-lg border border-neutral-800">
      <form onSubmit={handleSubmit}>
        {/* Header */}
        <div className="mb-8">
          <h2 className="text-xl font-semibold mb-2 text-white">Payment Method</h2>
          <p className="text-neutral-400 text-sm">All transactions are secure and encrypted</p>
        </div>

        {/* Name on Card */}
        <div className="mb-6">
          <label className="block text-sm font-medium mb-2 text-neutral-200">Name on Card</label>
          <input
            type="text"
            name="nameOnCard"
            value={formData.nameOnCard}
            onChange={handleInputChange}
            className="w-full bg-neutral-800 border border-neutral-700 rounded-md px-4 py-3 text-white placeholder-neutral-500 focus:border-web-green-500 focus:outline-none focus:ring-1 focus:ring-web-green-500 transition-colors"
            placeholder="John Doe"
          />
        </div>

        {/* Card Number and CVV */}
        <div className="grid grid-cols-3 gap-4 mb-4">
          <div className="col-span-2">
            <label className="block text-sm font-medium mb-2 text-neutral-200">Card Number</label>
            <input
              type="text"
              name="cardNumber"
              value={formData.cardNumber}
              onChange={handleInputChange}
              className="w-full bg-neutral-800 border border-neutral-700 rounded-md px-4 py-3 text-white placeholder-neutral-500 focus:border-web-green-500 focus:outline-none focus:ring-1 focus:ring-web-green-500 transition-colors"
              placeholder="1234 5678 9012 3456"
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-neutral-200">CVV</label>
            <input
              type="text"
              name="cvv"
              value={formData.cvv}
              onChange={handleInputChange}
              className="w-full bg-neutral-800 border border-neutral-700 rounded-md px-4 py-3 text-white placeholder-neutral-500 focus:border-web-green-500 focus:outline-none focus:ring-1 focus:ring-web-green-500 transition-colors"
              placeholder="123"
              maxLength="3"
            />
          </div>
        </div>

        <p className="text-neutral-400 text-xs mb-6">Enter your 16-digit number.</p>

        {/* Month and Year */}
        <div className="grid grid-cols-2 gap-4 mb-8">
          <div>
            <label className="block text-sm font-medium mb-2 text-neutral-200">Month</label>
            <select
              name="month"
              value={formData.month}
              onChange={handleInputChange}
              className="w-full bg-neutral-800 border border-neutral-700 rounded-md px-4 py-3 text-white focus:border-web-green-500 focus:outline-none focus:ring-1 focus:ring-web-green-500 appearance-none transition-colors"
            >
              <option value="MM">MM</option>
              <option value="01">01</option>
              <option value="02">02</option>
              <option value="03">03</option>
              <option value="04">04</option>
              <option value="05">05</option>
              <option value="06">06</option>
              <option value="07">07</option>
              <option value="08">08</option>
              <option value="09">09</option>
              <option value="10">10</option>
              <option value="11">11</option>
              <option value="12">12</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-neutral-200">Year</label>
            <select
              name="year"
              value={formData.year}
              onChange={handleInputChange}
              className="w-full bg-neutral-800 border border-neutral-700 rounded-md px-4 py-3 text-white focus:border-web-green-500 focus:outline-none focus:ring-1 focus:ring-web-green-500 appearance-none transition-colors"
            >
              <option value="YYYY">YYYY</option>
              <option value="2024">2024</option>
              <option value="2025">2025</option>
              <option value="2026">2026</option>
              <option value="2027">2027</option>
              <option value="2028">2028</option>
              <option value="2029">2029</option>
            </select>
          </div>
        </div>

        {/* Billing Address */}
        <div className="mb-6">
          <h3 className="text-lg font-medium mb-2 text-white">Billing Address</h3>
          <p className="text-neutral-400 text-sm mb-4">The billing address associated with your payment method</p>

          <div className="flex items-center">
            <input
              type="checkbox"
              name="sameAsShipping"
              checked={formData.sameAsShipping}
              onChange={handleInputChange}
              className="w-4 h-4 text-web-green-500 bg-neutral-800 border border-neutral-700 rounded focus:ring-web-green-500 focus:ring-2"
            />
            <label className="ml-2 text-sm text-neutral-200">Same as shipping address</label>
          </div>
        </div>

        {/* Comments */}
        <div className="mb-8">
          <label className="block text-sm font-medium mb-2 text-neutral-200">Comments</label>
          <textarea
            name="comments"
            value={formData.comments}
            onChange={handleInputChange}
            rows="4"
            className="w-full bg-neutral-800 border border-neutral-700 rounded-md px-4 py-3 text-white placeholder-neutral-500 focus:border-web-green-500 focus:outline-none focus:ring-1 focus:ring-web-green-500 resize-none transition-colors"
            placeholder="Add any additional comments"
          />
        </div>

        {/* Buttons */}
        <div className="flex gap-4">
          <button
            type="submit"
            className="flex-1 bg-web-green-500 hover:bg-web-green-600 text-white font-medium py-3 px-6 rounded-md transition-colors duration-200 shadow-sm"
          >
            Submit
          </button>
          <button
            type="button"
            onClick={handleCancel}
            className="flex-1 bg-neutral-700 hover:bg-neutral-600 text-white font-medium py-3 px-6 rounded-md transition-colors duration-200 shadow-sm"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  )
}

export default Payment