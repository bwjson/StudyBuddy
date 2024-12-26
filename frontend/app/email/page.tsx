'use client'
import axios from 'axios'
import React, { useEffect, useState } from 'react'

const EmailPage: React.FC = () => {
  const [email, setEmail] = useState('')
  const [message, setMessage] = useState('')
  const [subject, setSubject] = useState('')
  const [files, setFiles] = useState<File[]>([]) // Для хранения выбранных файлов
  const [rest, setRest] = useState('')

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setFiles(Array.from(e.target.files)) // Преобразуем объект FileList в массив
    } else {
      setFiles([])
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault() // Предотвращаем стандартное поведение формы

    const formData = new FormData()
    formData.append('email', email)
    formData.append('message', message)
    formData.append('subject', subject)

    // Добавляем файлы в FormData
    files.forEach((file) => {
      formData.append('attachments', file)
    })

    try {
      const response = await axios.post('http://localhost:8080/user/email', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })
      setRest(response.data.message)
      setEmail("")
      setFiles([])
      setSubject("")
      setMessage("")
    } catch (err: any) {
      setRest(err.response?.data?.message || 'Error sending email')
    }
  }

  useEffect(() => {
    if (rest) {
      const timer = setTimeout(() => setRest(''), 5000)
      return () => clearTimeout(timer)
    }
  }, [rest])

  return (
      <div className="flex w-100 h-[100vh]">
        <div className="flex flex-col w-1/2 h-1/2 m-auto space-y-4">
          <input
              type="email"
              placeholder="Enter email"
              className="text-black"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
          />

          <input
              type="text"
              placeholder="Enter subject"
              className="text-black"
              value={subject}
              onChange={(e) => setSubject(e.target.value)}
          />

          <textarea
              placeholder="Enter message"
              className="text-black"
              value={message}
              onChange={(e) => setMessage(e.target.value)}
          />

          <input
              type="file"
              multiple
              className="text-black"
              onChange={handleFileChange}
          />

          <button
              type="submit"
              className="bg-cyan-500 text-cyan-50"
              onClick={handleSubmit}
          >
            SEND
          </button>

          {rest && <p className="text-black text-lg">{rest}</p>}
        </div>
      </div>
  )
}

export default EmailPage
