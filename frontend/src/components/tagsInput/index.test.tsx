import React from 'react'
import { render, fireEvent, waitFor, act } from '@testing-library/react'

import { TagsInput } from './index'
import SearchInput from '../layout/list-content/search-input'

describe('TagsInput', () => {
  const onChange = jest.fn()
  const searchTags = jest.fn()

  it('displays given tags above text input', async () => {
    const tags = ['testing', 'tags', 'input']
    const searchTags = (_: string) => ({ fetching: false, foundTags: [] })

    const tagsInput = render(
      <TagsInput
        tags={tags}
        onSearch={searchTags}
        onChange={onChange}
        label='Tags'
        placeholder='placeholder'
      />,
    )

    const displayedTagValues = await tagsInput.findAllByTestId('tag-value')
    expect(tags[0]).toEqual(displayedTagValues[0].textContent)
    expect(tags[1]).toEqual(displayedTagValues[1].textContent)
    expect(tags[2]).toEqual(displayedTagValues[2].textContent)
  })

  it('performs a search with the entered keyword', async () => {
    // arrange
    const searchKeyword = 'aaa'

    // act
    const tagsInput = render(
      <TagsInput
        tags={[]}
        onSearch={searchTags}
        onChange={onChange}
        label='Tags'
        placeholder='placeholder'
      />,
    )

    const searchInput = await tagsInput.findByTestId('search-input')
    fireEvent.focus(searchInput)
    fireEvent.change(searchInput, { target: { value: 'aaa' } })

    // expect
    await waitFor(() => {
      expect(searchTags).toHaveBeenLastCalledWith(searchKeyword)
    })
  })
})
