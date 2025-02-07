'use client'
import React, { useState } from 'react'
import axios from 'axios'
import { useRouter } from 'next/navigation'

const regUser = async (
  email: string,
  password: string,
  name: string,
  username: string
) => {
  try {
    const res = await axios.post(
      `${process.env.NEXT_PUBLIC_API_URL}/auth/sign-up`,
      {
        email,
        password_hash: password,
        name,
        username,
      }
    )
    return res.status
  } catch (err: any) {
    return err.response?.status || 500
  }
}

const RegistrationPage: React.FC = () => {
  const navigation = useRouter()
  const [email, setEmail] = useState<string>('')
  const [password, setPassword] = useState<string>('')
  const [errors, setErrors] = useState<{
    email?: string
    password?: string
    name?: string
    userName?: string
  }>({})
  const [name, setName] = useState<string>('')
  const [userName, setUserName] = useState<string>('')
  const [loading, setLoading] = useState<boolean>(false)

  const validate = () => {
    const newErrors: {
      email?: string
      password?: string
      name?: string
      userName?: string
    } = {}

    if (!name) {
      newErrors.name = 'Name is required'
    }

    if (!userName) {
      newErrors.userName = 'Username is required'
    }

    if (!email) {
      newErrors.email = 'Email is required'
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
      const status = await regUser(email, password, name, userName)
      if (status === 200) {
        navigation.push('/login')
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
        <h2 className='pb-5 text-2xl font-semibold'>Registration</h2>

        <form
          className='flex flex-col gap-4 w-80 text-black'
          onSubmit={handleSubmit}
        >
          <div className='flex flex-col'>
            <input
              type='text'
              placeholder='Enter Name'
              className={`border rounded-lg p-2 outline-none focus:ring-2 ${
                errors.name
                  ? 'border-red-500 focus:ring-red-400'
                  : 'border-gray-300 focus:ring-blue-400'
              }`}
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
            {errors.name && (
              <span className='text-red-500 text-sm'>{errors.name}</span>
            )}
          </div>
          <div className='flex flex-col'>
            <input
              type='text'
              placeholder='Enter Username'
              className={`border rounded-lg p-2 outline-none focus:ring-2 ${
                errors.userName
                  ? 'border-red-500 focus:ring-red-400'
                  : 'border-gray-300 focus:ring-blue-400'
              }`}
              value={userName}
              onChange={(e) => setUserName(e.target.value)}
            />
            {errors.userName && (
              <span className='text-red-500 text-sm'>{errors.userName}</span>
            )}
          </div>

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
            Register
          </button>
        </form>
      </div>
    </div>
  )
}

export default RegistrationPage
