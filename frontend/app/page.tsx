'use client'
import axios from 'axios'
import { useEffect, useState } from 'react'

export default function Home() {
  const [data, setData] = useState<any>(null)
  useEffect(() => {
    const getData = async () => {
      const res: any = await axios.get('http://localhost:8080/auth/user/24')
      console.log(res)
      setData(res.data)
    }
    getData()
  }, [])
  return (
    <div>
      <h1 style={{ color: 'white' }}>{data?.user?.name}</h1>
      <p></p>
    </div>
  )
}
