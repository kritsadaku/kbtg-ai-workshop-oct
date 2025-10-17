import React, { useState } from 'react';

const Payment = () => {
  const [formData, setFormData] = useState({
    nameOnCard: 'John Doe',
    cardNumber: '1234 5678 9012 3456',
    cvv: '123',
    month: 'MM',
    year: 'YYYY',
    comments: '',
    sameAsShipping: true,
  });

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData((prevData) => ({
      ...prevData,
      [name]: type === 'checkbox' ? checked : value,
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    console.log('Payment submitted:', { ...formData });
  };

  const handleCancel = () => {
    console.log('Payment cancelled');
    // To reset the form:
    // setFormData({
    //   nameOnCard: '',
    //   cardNumber: '',
    //   cvv: '',
    //   month: 'MM',
    //   year: 'YYYY',
    //   comments: '',
    //   sameAsShipping: true,
    // });
  };

  return (
    <div className="bg-black text-white p-8 rounded-lg max-w-md mx-auto">
      <form onSubmit={handleSubmit}>
        {/* Header */}
        <div className="mb-8">
          <h2 className="text-xl font-semibold mb-2">Payment Method</h2>
          <p className="text-gray-400 text-sm">All transactions are secure and encrypted</p>
        </div>

        {/* Name on Card */}
        <div className="mb-6">
          <label htmlFor="nameOnCard" className="block text-sm font-medium mb-2">Name on Card</label>
          <input
            type="text"
            id="nameOnCard"
            name="nameOnCard"
            value={formData.nameOnCard}
            onChange={handleChange}
            className="w-full bg-gray-800 border border-gray-700 rounded-lg px-4 py-3 text-white placeholder-gray-500 focus:border-blue-500 focus:outline-none"
            placeholder="John Doe"
          />
        </div>

        {/* Card Number and CVV */}
        <div className="grid grid-cols-3 gap-4 mb-4">
          <div className="col-span-2">
            <label htmlFor="cardNumber" className="block text-sm font-medium mb-2">Card Number</label>
            <input
              type="text"
              id="cardNumber"
              name="cardNumber"
              value={formData.cardNumber}
              onChange={handleChange}
              className="w-full bg-gray-800 border border-gray-700 rounded-lg px-4 py-3 text-white placeholder-gray-500 focus:border-blue-500 focus:outline-none"
              placeholder="1234 5678 9012 3456"
            />
          </div>
          <div>
            <label htmlFor="cvv" className="block text-sm font-medium mb-2">CVV</label>
            <input
              type="text"
              id="cvv"
              name="cvv"
              value={formData.cvv}
              onChange={handleChange}
              className="w-full bg-gray-800 border border-gray-700 rounded-lg px-4 py-3 text-white placeholder-gray-500 focus:border-blue-500 focus:outline-none"
              placeholder="123"
              maxLength="3"
            />
          </div>
        </div>

        <p className="text-gray-400 text-xs mb-6">Enter your 16-digit number.</p>

        {/* Month and Year */}
        <div className="grid grid-cols-2 gap-4 mb-8">
          <div>
            <label htmlFor="month" className="block text-sm font-medium mb-2">Month</label>
            <select
              id="month"
              name="month"
              value={formData.month}
              onChange={handleChange}
              className="w-full bg-gray-800 border border-gray-700 rounded-lg px-4 py-3 text-white focus:border-blue-500 focus:outline-none appearance-none"
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
            <label htmlFor="year" className="block text-sm font-medium mb-2">Year</label>
            <select
              id="year"
              name="year"
              value={formData.year}
              onChange={handleChange}
              className="w-full bg-gray-800 border border-gray-700 rounded-lg px-4 py-3 text-white focus:border-blue-500 focus:outline-none appearance-none"
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
          <h3 className="text-lg font-medium mb-2">Billing Address</h3>
          <p className="text-gray-400 text-sm mb-4">The billing address associated with your payment method</p>

          <div className="flex items-center">
            <input
              type="checkbox"
              id="sameAsShipping"
              name="sameAsShipping"
              checked={formData.sameAsShipping}
              onChange={handleChange}
              className="w-4 h-4 bg-blue-600 border-blue-600 rounded focus:ring-blue-500 focus:ring-2"
            />
            <label htmlFor="sameAsShipping" className="ml-2 text-sm">Same as shipping address</label>
          </div>
        </div>

        {/* Comments */}
        <div className="mb-8">
          <label htmlFor="comments" className="block text-sm font-medium mb-2">Comments</label>
          <textarea
            id="comments"
            name="comments"
            value={formData.comments}
            onChange={handleChange}
            rows="4"
            className="w-full bg-gray-800 border border-gray-700 rounded-lg px-4 py-3 text-white placeholder-gray-500 focus:border-blue-500 focus:outline-none resize-none"
            placeholder="Add any additional comments"
          />
        </div>

        {/* Buttons */}
        <div className="flex gap-4">
          <button
            type="submit"
            className="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 px-6 rounded-lg transition-colors"
          >
            Submit
          </button>
          <button
            type="button"
            onClick={handleCancel}
            className="flex-1 bg-gray-700 hover:bg-gray-600 text-white font-medium py-3 px-6 rounded-lg transition-colors"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
};

export default Payment;
