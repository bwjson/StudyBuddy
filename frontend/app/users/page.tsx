'use client'
import axios from 'axios'
import { useRouter } from 'next/navigation'
import React, { useEffect, useState } from 'react'

const UsersPage = () => {
  const [allTags, setAllTags] = useState<any[]>([])
  const [selectedTag, setSelectedTag] = useState<any>(null)
  const [users, setUsers] = useState<any[]>([])
  const navigate = useRouter()
  const [page, setPage] = useState(1)
  const [limit, setLimit] = useState(10)
  const [totalCount, setTotalCount] = useState(0)
  const [totalPages, setTotalPages] = useState(0)
  const [loading, setLoading] = useState<boolean>(false)
  const [error, setError] = useState<string | null>(null)

  const [filters, setFilters] = useState<{
    sort_by: string
    sort_order: string
  }>({
    sort_by: 'name', // Default sort by 'name'
    sort_order: 'asc', // Default sort order 'ascending'
  })

  const API_BASE_URL =
    process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

  useEffect(() => {
    const fetchTags = async () => {
      try {
        setLoading(true)
        const res = await axios.get(`${API_BASE_URL}/user/tags`)
        setAllTags(res.data.data)
      } catch (err) {
        console.error('Ошибка при получении тегов:', err)
        setError('Не удалось загрузить теги.')
      } finally {
        setLoading(false)
      }
    }

    fetchTags()
  }, [API_BASE_URL])

  useEffect(() => {
    const fetchUsers = async () => {
      setLoading(true)
      setError(null)
      try {
        let response
        if (selectedTag) {
          response = await axios.get(
            `${API_BASE_URL}/user/tag/${selectedTag.id}`,
            {
              params: {
                limit,
                page,
                ...(filters.sort_by && { sort_by: filters.sort_by }),
                ...(filters.sort_order && { sort_order: filters.sort_order }),
              },
            }
          )
        } else {
          response = await axios.get(`${API_BASE_URL}/user/`, {
            params: {
              ...(filters.sort_by && { sort_by: filters.sort_by }),
              ...(filters.sort_order && { sort_order: filters.sort_order }),
              limit,
              page,
            },
          })
        }
        setUsers(response.data.data.users)
        setTotalCount(response.data.data.totalCount)
        setTotalPages(Math.ceil(response.data.data.totalCount / limit))
      } catch (err) {
        console.error('Ошибка при получении пользователей:', err)
      } finally {
        setLoading(false)
      }
    }

    fetchUsers()
  }, [selectedTag, page, limit, filters, API_BASE_URL])

  console.log(totalPages)

  const handleTagSelect = (tag: any) => {
    setSelectedTag(tag)
    setPage(1)
  }

  const handleResetFilter = () => {
    setSelectedTag(null)
    setPage(1)
  }

  const handleSortChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const { name, value } = e.target
    setFilters((prev) => ({
      ...prev,
      [name]: value,
    }))
    setPage(1)
  }

  const getPageNumbers = () => {
    const maxPageNumbersToShow = 5
    let start = Math.max(page - Math.floor(maxPageNumbersToShow / 2), 1)
    let end = start + maxPageNumbersToShow - 1

    if (end > totalPages) {
      end = totalPages
      start = Math.max(end - maxPageNumbersToShow + 1, 1)
    }

    const pages = []
    for (let i = start; i <= end; i++) {
      pages.push(i)
    }
    return pages
  }

  return (
    <div className='flex items-center justify-center p-4'>
      <div className='flex flex-col w-full max-w-4xl h-auto'>
        <h1 className='mt-10 text-center mb-10 text-2xl font-bold'>
          Users Page
        </h1>

        <div className='flex flex-wrap gap-4 justify-center'>
          {allTags.slice(0, 8).map((tag: any) => (
            <div
              key={tag.id}
              className={`border-2 p-2 rounded-2xl cursor-pointer ${
                selectedTag?.id === tag.id ? 'bg-white text-black' : ''
              }`}
              onClick={() => handleTagSelect(tag)}
            >
              <p className=''>{tag.title}</p>
            </div>
          ))}
          {selectedTag && (
            <div
              className='border-2 border-black p-2 rounded-2xl cursor-pointer bg-white text-black'
              onClick={handleResetFilter}
            >
              <p>Сбросить фильтр</p>
            </div>
          )}
        </div>

        <div className='mt-6 flex flex-wrap gap-4 justify-center items-center'>
          <div className='flex items-center'>
            <label htmlFor='sort-by' className='mr-2 font-medium'>
              Сортировать по:
            </label>
            <select
              id='sort-by'
              name='sort_by'
              className='border-black border-2 p-2 rounded-2xl cursor-pointer bg-white text-black'
              value={filters.sort_by}
              onChange={handleSortChange}
            >
              <option value='name'>Имя</option>
              <option value='email'>Email</option>
            </select>
          </div>
          <div className='flex items-center'>
            <label htmlFor='sort-order' className='mr-2 font-medium'>
              Порядок:
            </label>
            <select
              id='sort-order'
              name='sort_order'
              className='border-black border-2 p-2 rounded-2xl cursor-pointer bg-white text-black'
              value={filters.sort_order}
              onChange={handleSortChange}
            >
              <option value='asc'>Возрастание</option>
              <option value='desc'>Убывание</option>
            </select>
          </div>
        </div>

        <div className='mt-6 flex justify-end items-center'>
          <label htmlFor='limit-select' className='mr-2 font-medium'>
            Показать:
          </label>
          <select
            id='limit-select'
            className='border-black border-2 p-2 rounded-2xl cursor-pointer bg-white text-black'
            value={limit}
            onChange={(e) => {
              setLimit(parseInt(e.target.value))
              setPage(1) // Сброс страницы на 1 при изменении лимита
            }}
          >
            <option value={10}>10</option>
            <option value={20}>20</option>
            <option value={50}>50</option>
          </select>
        </div>

        <div className='mt-6'>
          {loading ? (
            <p className='text-center'>Загрузка...</p>
          ) : error ? (
            <p className='text-center text-red-500'>{error}</p>
          ) : users.length > 0 ? (
            users.map((user: any) => (
              <div
                key={user.id}
                className='border-black border-2 p-4 rounded-2xl mt-2 cursor-pointer hover:bg-black transition'
                onClick={() => navigate.push(`/user/${user.id}`)}
              >
                <p className='font-semibold'>{user.name}</p>
                <p className='text-gray-600'>{user.email}</p>
              </div>
            ))
          ) : (
            <p className='text-center'>Пользователи не найдены.</p>
          )}
        </div>

        <div className='flex justify-center items-center mt-6 gap-2'>
          <button
            onClick={() => setPage((prev) => Math.max(prev - 1, 1))}
            disabled={page === 1}
            className={`px-4 py-2 border rounded ${
              page === 1
                ? 'opacity-50 cursor-not-allowed'
                : 'hover:bg-black text-black hover:text-white'
            }`}
          >
            Назад
          </button>

          {getPageNumbers().map((pageNumber) => (
            <button
              key={pageNumber}
              onClick={() => setPage(pageNumber)}
              className={`px-4 py-2 border rounded ${
                pageNumber === page
                  ? 'bg-white text-black'
                  : 'hover:bg-black hover:text-white'
              }`}
            >
              {pageNumber}
            </button>
          ))}

          <button
            onClick={() => setPage((prev) => Math.min(prev + 1, totalPages))}
            disabled={page === totalPages || totalPages === 0}
            className={`px-4 py-2 rounded border ${
              page === totalPages || totalPages === 0
                ? 'opacity-50 cursor-not-allowed'
                : 'hover:bg-black text-black hover:text-white'
            }`}
          >
            Вперед
          </button>
        </div>

        <div className='flex justify-center items-center mt-2'>
          <p>
            Страница {page} из {totalPages}
          </p>
        </div>
      </div>
    </div>
  )
}

export default UsersPage
