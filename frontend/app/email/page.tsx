'use client'
import axios from 'axios'
import React, { useEffect, useState } from 'react'

const EmailPage: React.FC = () => {
  const [email, setEmail] = useState('')
  const [message, setMessage] = useState('')
  const [subject, setSubject] = useState('')
  const [rest, setRest] = useState('')

  const handleSubmit = async () => {
    await axios
      .post('http://localhost:8080/user/email', {
        email,
        message,
        subject,
      })
      .then((res) => {
        setRest(res.data.message)
      })
      .catch((err) => {
        setRest(err.response.data.message)
      })
  }

  useEffect(() => {
    if (rest) {
      const timer = setTimeout(() => setRest(''), 5000)
      return () => clearTimeout(timer)
    }
  }, [rest])

  return (
    <div className='flex w-100 h-[100vh]'>
      <div className='flex flex-col w-1/2 h-1/2 m-auto space-y-4'>
        <input
          type='email'
          placeholder='Enter email'
          className='text-black'
          onChange={(e) => setEmail(e.target.value)}
        />

        <input
          type='text'
          placeholder='Enter subject'
          className='text-black'
          onChange={(e) => setSubject(e.target.value)}
        />

        <textarea
          placeholder='Enter message'
          className='text-black'
          onChange={(e) => setMessage(e.target.value)}
        />
        <button
          type='submit'
          className='bg-cyan-500 text-cyan-50'
          onClick={handleSubmit}
        >
          SEND
        </button>
        {rest && <p className='text-black text-lg'>{rest}</p>}
      </div>
    </div>
  )
}

export default EmailPage
