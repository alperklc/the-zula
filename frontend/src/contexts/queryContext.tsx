/* eslint-disable @typescript-eslint/no-unused-vars */
import React from 'react'
import { useSearchParams } from 'react-router-dom'

export type Query = {
  tags: string
  q: string
  page: number
  pageSize: number
  sortBy: string
  sortDirection: string
}

type Filter = {
  tags: string[]
}

type Sort = {
  sortBy: string
  sortDirection: string
}

export interface QueryState {
  query: Query
  tagsAsArray: string[]
  updateQuery: (_: string) => void
  paginate: (_: number) => void
  changePageSize: (_: number) => void
  changeSort: (_: string, __: string) => void
  changeFilter: (_: Filter) => void
  applyFilterAndSort: (_: Filter & Sort) => void
}

const initialQuery: Query = {
  tags: '',
  q: '',
  page: 1,
  pageSize: 10,
  sortBy: 'updatedAt',
  sortDirection: 'desc',
}

export const QueryContext = React.createContext<QueryState>({
  query: initialQuery,
  tagsAsArray: [],
  updateQuery: (_: string) => ({}),
  paginate: (_: number) => ({}),
  changePageSize: (_: number) => ({}),
  changeSort: (_: string, __: string) => ({}),
  changeFilter: (_: Filter) => ({}),
  applyFilterAndSort: (_: Filter & Sort) => ({}),
})

export const useQuery = () => React.useContext(QueryContext)

export const QueryContextProvider = ({ query, children }: { query: Query, children: React.ReactNode }) => {
  const [_, setSearchParams] = useSearchParams();

  // This updates querystring on the page, once state changes
  const updateQueryString = (newQuery: Query) => {
    const nextQueryString = `q=${newQuery.q}&tags=${newQuery.tags}&page=${newQuery.page}&pageSize=${newQuery.pageSize}&sortBy=${newQuery.sortBy}&sortDirection=${newQuery.sortDirection}`

    setSearchParams(`${nextQueryString}`, {
      state: nextQueryString,
      replace: true,
    })
  }

  const updateQuery = (q: string) => {
    updateQueryString({ ...query, q, page: 1 })
  }

  const paginate = (page: number) => {
    updateQueryString({ ...query, page })
  }

  const changePageSize = (pageSize: number) => {
    updateQueryString({ ...query, pageSize })
  }

  const changeSort = (sortBy: string, sortDirection: string) => {
    updateQueryString({ ...query, sortBy, sortDirection })
  }

  const changeFilter = ({ tags }: Filter) => {
    updateQueryString({ ...query, tags: tags.join(',') })
  }

  const applyFilterAndSort = ({
    tags,
    sortBy,
    sortDirection,
  }: {
    tags: string[]
    sortBy: string
    sortDirection: string
  }) => {
    updateQueryString({
      ...query,
      tags: tags.join(','),
      sortBy: sortBy || query.sortBy,
      sortDirection: sortDirection || query.sortDirection,
    })
  }

  const tagsAsArray = query?.tags.length > 0 ? query?.tags.split(',') : []

  return (
    <QueryContext.Provider
      value={{
        query,
        tagsAsArray,
        updateQuery,
        paginate,
        changePageSize,
        changeSort,
        changeFilter,
        applyFilterAndSort,
      }}
    >
      {children}
    </QueryContext.Provider>
  )
}
