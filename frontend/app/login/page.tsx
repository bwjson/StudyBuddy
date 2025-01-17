'use client'
import React, { useState } from 'react'
import axios from 'axios'
import { useRouter } from 'next/navigation'
import Cookies from 'js-cookie'

const loginUser = async (email: string, password: string) => {
  try {
    const res = await axios.post(
      `${process.env.NEXT_PUBLIC_API_URL}/auth/sign-in`,
      {
        email,
        password: password,
      }
    )

    if (res.data?.data) {
      Cookies.set('accToken', res.data.data.access_token)
      Cookies.set('refToken', res.data.data.refresh_token)
    }
    return res.status
  } catch (err: any) {
    return err.response?.status || 500
  }
}

const LoginPage: React.FC = () => {
  const navigation = useRouter()
  const [email, setEmail] = useState<string>('')
  const [password, setPassword] = useState<string>('')
  const [errors, setErrors] = useState<{ email?: string; password?: string }>(
    {}
  )
  const [loading, setLoading] = useState<boolean>(false)

  const validate = () => {
    const newErrors: { email?: string; password?: string } = {}
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

    if (!email) {
      newErrors.email = 'Email is required'
    } else if (!emailRegex.test(email)) {
      newErrors.email = 'Invalid email format'
    }

    if (!password) {
      newErrors.password = 'Password is required'
    } else if (password.length < 6) {
      newErrors.password = 'Password must be at least 6 characters'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (validate()) {
      setLoading(true)
      const status = await loginUser(email, password)
      if (status === 200) {
        navigation.push('/')
        setLoading(false)
      } else {
        alert('Registration failed. Please try again.')
        setLoading(false)
      }
    }
  }

  return (
    <div className='flex w-100 h-[100vh] justify-center items-center'>
      <div className='flex flex-col border rounded-xl p-10 items-center'>
        <h2 className='pb-5 text-2xl font-semibold'>Login</h2>

        <form
          className='flex flex-col gap-4 w-80 text-black'
          onSubmit={handleSubmit}
        >
          <div className='flex flex-col'>
            <input
              type='email'
              placeholder='Enter email'
              className={`border rounded-lg p-2 outline-none focus:ring-2 ${
                errors.email
                  ? 'border-red-500 focus:ring-red-400'
                  : 'border-gray-300 focus:ring-blue-400'
              }`}
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
            {errors.email && (
              <span className='text-red-500 text-sm'>{errors.email}</span>
            )}
          </div>

          <div className='flex flex-col'>
            <input
              type='password'
              placeholder='Enter password'
              className={`border rounded-lg p-2 outline-none focus:ring-2 ${
                errors.password
                  ? 'border-red-500 focus:ring-red-400'
                  : 'border-gray-300 focus:ring-blue-400'
              }`}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            {errors.password && (
              <span className='text-red-500 text-sm'>{errors.password}</span>
            )}
          </div>

          <button
            type='submit'
            className={`bg-blue-500 text-white p-2 rounded-lg hover:bg-blue-600 transition duration-200 ${
              loading ? 'bg-slate-500 cursor-not-allowed' : ''
            }`}
            disabled={loading}
          >
            Login
          </button>
        </form>
      </div>
    </div>
  )
}

export default LoginPage
