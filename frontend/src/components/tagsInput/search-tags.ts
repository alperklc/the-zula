import { useEffect, useState } from 'react'
// import { useQuery } from 'urql'
import useDebounce from '../../utils/useDebounce'

export const searchTags = (typeOfParent: string) => (textInputValue: string) => {
  const debouncedInput = useDebounce(textInputValue, 500)
  const [searchKeyword, setSearchKeyword] = useState<string>(textInputValue)

  useEffect(() => {
    setSearchKeyword(textInputValue)
  }, [debouncedInput])

 // const [{ fetching, data }] = useQuery({
 //   pause: !searchKeyword,
 //   variables: { q: searchKeyword, typeOfParent: typeOfParent || '' },
 //   query: `
 //       query($typeOfParent: String, $q: String!) {
 //         tags(typeOfParent: $typeOfParent, q: $q) {
 //           value
 //         }
 //       }
 //     `,
 // })
//
  return { fetching: false, foundTags: [] }
}
