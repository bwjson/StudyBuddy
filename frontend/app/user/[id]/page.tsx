'use client'

import axios from 'axios'
import { useParams } from 'next/navigation'
import React, { useEffect, useState } from 'react'

const UserPage = () => {
  const params = useParams()
  const { id } = params
  const [user, setUser] = useState<any>(null)
  const [error, setError] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  const [tags, setTags] = useState<any>([])

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const res = await axios.get(
          `${process.env.NEXT_PUBLIC_API_URL}/user/${id}`
        )
        setUser(res.data)
      } catch (err) {
        console.error(err)
      } finally {
        setLoading(false)
      }
    }

    const fetchTags = async () => {
      try {
        const res = await axios.get(
          `${process.env.NEXT_PUBLIC_API_URL}user/usertags/${id}`
        )
        setTags(res.data.data)
      } catch (err) {
        console.error(err)
      } finally {
        setLoading(false)
      }
    }

    if (id) {
      fetchUser()
      fetchTags()
    }
  }, [id])

  console.log(tags)

  if (loading)
    return (
      <div className='flex w-100 h-[100vh] justify-center items-center'>
        Loading...
      </div>
    )
  if (error) return <div>{error}</div>

  return (
    <div className='flex items-center justify-center'>
      <div className='flex-col'>
        <h1 className='text-lg mt-5'>User Page</h1>
        {user.data ? (
          <div className='mt-5'>
            <p>ID: {user.data.id}</p>
            <p>Name: {user.data.name}</p>
            <p>UserName: {user.data.username}</p>
            <p>Password Hash: {user.data.password_hash}</p>
          </div>
        ) : (
          <p>User not found.</p>
        )}

        <div className='flex-col mt-10'>
          <p>{user.data.name} tags: </p>
          <div className='flex gap-2 mt-2'>
            {tags ? (
              tags.map((tag: any) => (
                <div
                  key={tag.id}
                  className='border-black border-2 p-2 rounded-2xl cursor-pointer'
                >
                  <p className=''>{tag.title}</p>
                </div>
              ))
            ) : (
              <p>No tags found.</p>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default UserPage
