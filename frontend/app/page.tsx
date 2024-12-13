'use client'
import axios from 'axios'
import { useEffect, useState } from 'react'

export default function Home() {
  const [user, setUser] = useState<any>(null)
  const [name, setName] = useState('')
  const [userName, setUserName] = useState('')
  const [password, setPassword] = useState('')
  const [messageCreate, setMessageCreate] = useState('')
  const [messageUpdate, setMessageUpdate] = useState('')
  const [users, setUsers] = useState<any>([])

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
    setMessageCreate('User succesfully created')
  }

  const handleGetUser = async (id: number) => {
    const res = await axios.get(`http://localhost:8080/auth/user/${id}`)
    setUser(res.data)
  }

  const handleDeleteUser = async (id: number) => {
    await axios.delete(`http://localhost:8080/auth/user/${id}`)
  }

  const handleUpdateUser = async (id: number) => {
    const res = await axios.put(`http://localhost:8080/auth/user/${id}`, {
      name,
      username: userName,
      passwordhash: password,
    })

    setMessageUpdate(res.data.message)
  }

  const handleGetAllUsers = async () => {
    const res = await axios.get('http://localhost:8080/auth/users')
    setUsers(res.data)
  }

  return (
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
          className='w-[100%] mb-2'
          onChange={(e) => setUserName(e.target.value)}
        />
        <input
          type='text'
          placeholder='Enter Password'
          className='w-[100%] mb-2'
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
          className='w-[100%] mb-2'
          onChange={(e) => setName(e.target.value)}
        />
        <input
          type='text'
          placeholder='Enter User Name'
          className='w-[100%] mb-2'
          onChange={(e) => setUserName(e.target.value)}
        />
        <input
          type='text'
          placeholder='Enter Password'
          onChange={(e) => setPassword(e.target.value)}
          className='w-[100%] mb-2'
        />
        <button
          onClick={() => handleUpdateUser(id)}
          className='pl-5 pr-5 bg-slate-400 rounded-full'
        >
          Update
        </button>
        {messageUpdate && <p>{'User succesfully created'}</p>}
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
            {users.users?.map((user: any) => (
              <div>
                <h1>{user.name}</h1>
                <h1>{user.username}</h1>
                <h1>{user.passwordhash}</h1>
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
          className='w-[100%] mb-2'
        />
        <button
          onClick={() => handleGetUser(id)}
          className='pl-5 pr-5 bg-slate-400 rounded-full'
        >
          Get
        </button>
        {user !== null && (
          <div>
            <h1>{user.name}</h1>
            <h1>{user.username}</h1>
            <h1>{user.passwordhash}</h1>
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
          className='w-[100%] mb-2'
        />
        <button
          onClick={() => handleDeleteUser(id)}
          className='pl-5 pr-5 bg-slate-400 rounded-full'
        >
          Delete
        </button>
      </div>
    </div>
  )
}
