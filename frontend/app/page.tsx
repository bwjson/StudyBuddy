'use client'
import axios from 'axios'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'

export default function Home() {
  const [user, setUser] = useState<any>(null)
  const [name, setName] = useState('')
  const [userName, setUserName] = useState('')
  const [password, setPassword] = useState('')
  const [messageCreate, setMessageCreate] = useState('')
  const [messageUpdate, setMessageUpdate] = useState('')
  const [users, setUsers] = useState<any>([])
  const [idForUpdate, setIdForUpdate] = useState(0)
  const [id, setId] = useState(0)
  const navigation = useRouter()

  const handleCreateUser = async () => {
    const res = await axios.post('http://localhost:8080/user/', {
      name,
      username: userName,
      password_hash: password,
    })
    setMessageCreate(res.data.message)
  }

  const handleGetUser = async (id: number) => {
    const res = await axios.get(`http://localhost:8080/user/${id}/`)
    setUser(res.data)
  }

  const handleDeleteUser = async (id: number) => {
    await axios.delete(`http://localhost:8080/user/${id}/`)
  }

  const handleUpdateUser = async (id: number) => {
    const res = await axios.put(`http://localhost:8080/user/${id}/`, {
      name,
      username: userName,
      password_hash: password,
    })

    setMessageUpdate(res.data.message)
  }

  const handleGetAllUsers = async () => {
    const res = await axios.get('http://localhost:8080/user/')
    setUsers(res.data)
  }

  useEffect(() => {
    if (messageCreate) {
      const timer = setTimeout(() => setMessageCreate(''), 5000)
      return () => clearTimeout(timer)
    }
  }, [messageCreate])

  useEffect(() => {
    if (messageUpdate) {
      const timer = setTimeout(() => setMessageUpdate(''), 5000)
      return () => clearTimeout(timer)
    }
  }, [messageUpdate])

  return (
    <>
      <div className='flex justify-center items-center mt-9 mb-9'>
        <button
          onClick={() => navigation.push('/email')}
          className=' pl-5 pr-5 bg-slate-400 rounded-full'
        >
          Email page
        </button>
      </div>

      <div className='flex justify-center items-center mt-9 mb-9'>
        <button
          onClick={() => navigation.push('/users')}
          className=' pl-5 pr-5 bg-slate-400 rounded-full'
        >
          Users Page
        </button>
      </div>
      <div className='flex flex-column justify-content-center align-items-center gap-10'>
        <div>
          <h1>Create user</h1>
          <input
            type='text'
            placeholder='Enter Name'
            className='w-[100%] mb-2 text-black'
            onChange={(e) => setName(e.target.value)}
          />
          <input
            type='text'
            placeholder='Enter User Name'
            className='w-[100%] mb-2 text-black'
            onChange={(e) => setUserName(e.target.value)}
          />
          <input
            type='text'
            placeholder='Enter Password'
            className='w-[100%] mb-2 text-black'
            onChange={(e) => setPassword(e.target.value)}
          />
          <button
            onClick={handleCreateUser}
            className='pl-5 pr-5 bg-slate-400 rounded-full'
          >
            Create
          </button>
          {messageCreate && <p>{'User succesfully created'}</p>}
        </div>

        <hr />

        <div>
          <h1>Update user</h1>

          <input
            type='text'
            placeholder='Enter Name'
            className='w-[100%] mb-2  text-black'
            onChange={(e) => setName(e.target.value)}
          />
          <input
            type='text'
            placeholder='Enter User Name'
            className='w-[100%] mb-2 text-black'
            onChange={(e) => setUserName(e.target.value)}
          />
          <input
            type='text'
            placeholder='Enter Password'
            onChange={(e) => setPassword(e.target.value)}
            className='w-[100%] mb-2 text-black'
          />
          <input
            type='text'
            placeholder='ID'
            onChange={(e) => setIdForUpdate(Number(e.target.value))}
            className='w-[100%] mb-2 text-black'
          />
          <button
            onClick={() => handleUpdateUser(idForUpdate)}
            className='pl-5 pr-5 bg-slate-400 rounded-full'
          >
            Update
          </button>
          {messageUpdate && <p>{messageUpdate}</p>}
        </div>

        <hr />

        <div>
          <h1>Get all users</h1>
          <button
            onClick={handleGetAllUsers}
            className='pl-5 pr-5 bg-slate-400 rounded-full'
          >
            Get
          </button>
          {users !== null && (
            <div className=''>
              {users.data?.map((user: any) => (
                <div>
                  <h1>{user.id}</h1>
                  <h1>{user.name}</h1>
                  <h1>{user.username}</h1>
                  <h1>{user.password_hash}</h1>
                  <hr className='bg-white' />
                </div>
              ))}
            </div>
          )}
        </div>

        <hr />

        <div>
          <h1>Get user by id</h1>
          <input
            type='text'
            placeholder='ID'
            onChange={(e) => setId(Number(e.target.value))}
            className='w-[100%] mb-2 text-black'
          />
          <button
            onClick={() => handleGetUser(id)}
            className='pl-5 pr-5 bg-slate-400 rounded-full'
          >
            Get
          </button>
          {user && (
            <div className='text-white'>
              <h1>{user.data.id}</h1>
              <h1>{user.data.name}</h1>
              <h1>{user.data.username}</h1>
              <h1>{user.data.password_hash}</h1>
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
            className='w-[100%] mb-2 text-black'
          />
          <button
            onClick={() => handleDeleteUser(id)}
            className='pl-5 pr-5 bg-slate-400 rounded-full'
          >
            Delete
          </button>
        </div>
      </div>
    </>
  )
}
