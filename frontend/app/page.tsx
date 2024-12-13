'use client'
import axios from 'axios'
import { useEffect, useState } from 'react'

export default function Home() {
  const [data, setData] = useState<any>(null)
  const [name, setName] = useState('')
  const [userName, setUserName] = useState('')
  const [password, setPassword] = useState('')
  const [message, setMessage] = useState('')

  const [id, setId] = useState(0)
  // useEffect(() => {
  //   const getData = async () => {
  //     const res: any = await axios.get('http://localhost:8080/auth/user/24')
  //     setData(res.data)
  //   }
  //   getData()
  // }, [])

  const handleCreateUser = async () => {
    const res = await axios.post('http://localhost:8080/auth/create', {
      name,
      username: userName,
      passwordhash: password,
    })
    setMessage('User succesfully created')
  }

  const handleGetUser = async (id: number) => {
    const res = await axios.get(`http://localhost:8080/auth/user/${id}`)
    setData(res.data)
  }

  const handleDeleteUser = async (id: number) => {
    await axios.delete(`http://localhost:8080/auth/user/${id}`)
  }

  const handleUpdateUser = async (id: number) => {
    await axios.put(`http://localhost:8080/auth/user/${id}`, {
      name,
      username: userName,
      passwordhash: password,
    })
  }

  return (
    <div className='flex flex-column justify-content-center align-items-center gap-3'>
      <div>
        <h1>Create user</h1>
        <input
          type='text'
          placeholder='Enter Name'
          onChange={(e) => setName(e.target.value)}
        />
        <input
          type='text'
          placeholder='Enter User Name'
          onChange={(e) => setUserName(e.target.value)}
        />
        <input
          type='text'
          placeholder='Enter Password'
          onChange={(e) => setPassword(e.target.value)}
        />
        <button
          onClick={handleCreateUser}
          className='pl-5 pr-5 bg-slate-400 rounded-full mt-5'
        >
          Create
        </button>
        {message && <p>{'User succesfully created'}</p>}
      </div>

      <hr />

      <div>
        <h1>Update user</h1>

        <input
          type='text'
          placeholder='Enter Name'
          onChange={(e) => setName(e.target.value)}
        />
        <input
          type='text'
          placeholder='Enter User Name'
          onChange={(e) => setUserName(e.target.value)}
        />
        <input
          type='text'
          placeholder='Enter Password'
          onChange={(e) => setPassword(e.target.value)}
        />
        <button
          onClick={() => handleUpdateUser(id)}
          className='pl-5 pr-5 bg-slate-400 rounded-full mt-5'
        >
          Update
        </button>
        {message && <p>{'User succesfully created'}</p>}
      </div>

      <hr />

      <div>
        <h1>Get user by id</h1>
        <input
          type='text'
          placeholder='ID'
          onChange={(e) => setId(Number(e.target.value))}
        />
        <button
          onClick={() => handleGetUser(id)}
          className='pl-5 pr-5 bg-slate-400 rounded-full mt-5'
        >
          Get
        </button>
        {data !== null && (
          <div>
            <h1>{data.name}</h1>
            <h1>{data.username}</h1>
            <h1>{data.passwordhash}</h1>
          </div>
        )}
      </div>

      <hr />

      <div>
        <h1>Delete user</h1>
        <input
          type='text'
          placeholder='ID'
          onChange={(e) => setId(Number(e.target.value))}
        />
        <button
          onClick={() => handleDeleteUser(id)}
          className='pl-5 pr-5 bg-slate-400 rounded-full mt-5'
        >
          Delete
        </button>
      </div>
    </div>
  )
}
