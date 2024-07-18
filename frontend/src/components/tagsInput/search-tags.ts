import { useEffect, useState } from 'react'
import useDebounce from '../../utils/useDebounce'
import { useAuth } from '../../contexts/authContext'
import { Api, Tag } from '../../types/Api.ts'

export const searchTags = (typeOfParent: string) => (textInputValue: string) => {
  const { user } = useAuth()
  const api = new Api({ baseApiParams: { headers: { authorization: `Bearer ${user?.access_token}` } } })
  
  const debouncedInput = useDebounce(textInputValue, 500)
  const [searchKeyword, setSearchKeyword] = useState<string>(textInputValue)
  const [fetching, setFetching] = useState<boolean>(false)
  const [tags, setTags] = useState<Tag[]>([])

  useEffect(() => {
    setSearchKeyword(textInputValue)
  }, [debouncedInput])

  const fetchAndUpdate = async (query: string) => {
    if (!query) {
      return
    }

    try {
      setFetching(true);

      const { data } = await api.api.getTags({ type: typeOfParent, q: query })
      setTags(data)

    } catch (e: unknown) {
      console.error(e);
    }
    setFetching(false);
  };

  useEffect(() => {
    fetchAndUpdate(searchKeyword)
  }, [searchKeyword])

  return { fetching, foundTags: tags }
}
