import React from 'react'
import { act, render } from '@testing-library/react'

import TagsDisplay from './index'

describe('TagsDisplay', () => {
  it('displays all tags correctly', async () => {
    // arrange
    const tags = ['this', 'is', 'for', 'testing']

    // act
    const tagsDisplay = render(<TagsDisplay tags={tags} />)
    const renderedTags = (await tagsDisplay.findAllByTestId('tag-value')).map(
      (element: HTMLElement) => element.textContent,
    )

    // assert
    tags.forEach((tag) => {
      expect(renderedTags).toContain(tag)
    })
  })

  it('displays only first N tags, if maxNumberOfTagsToDisplay is given', async () => {
    // arrange
    const tags = ['this', 'is', 'for', 'testing']
    const maxNumberOfTagsToDisplay = 2

    // act
    const tagsDisplay = render(
      <TagsDisplay tags={tags} maxNumberOfTagsToDisplay={maxNumberOfTagsToDisplay} />,
    )
    const renderedTags = (await tagsDisplay.findAllByTestId('tag-value')).map(
      (element: HTMLElement) => element.textContent,
    )

    // assert
    expect(renderedTags).toContain('this')
    expect(renderedTags).toContain('is')
    expect(renderedTags).not.toContain('for')
    expect(renderedTags).not.toContain('testing')
  })

  it('displays an icon for removing tags if a removeTag prop is given', async () => {
    // arrange
    const tags = ['this', 'is', 'for', 'testing']

    // act
    const tagsDisplay = render(<TagsDisplay tags={tags} removeTag={() => ({})} />)
    const removeIcons = await tagsDisplay.findAllByTestId('remove-tag-icon')

    // assert
    expect(removeIcons).toHaveLength(tags.length)
  })

  it('calls removeTag function with the correct tag value if user clicks on remove-tag-icon', async () => {
    // arrange
    const tags = ['this', 'is', 'for', 'testing']
    const mockRemoveTag = vi.fn()
    const indexOfTagToBeRemoved = 2

    // act
    const tagsDisplay = render(<TagsDisplay tags={tags} removeTag={mockRemoveTag} />)
    const removeIcons = await tagsDisplay.findAllByTestId('remove-tag-icon')
    act(() => {
      removeIcons[indexOfTagToBeRemoved].click()
    })

    // assert
    expect(removeIcons).toHaveLength(tags.length)
    expect(mockRemoveTag).toHaveBeenCalledWith(tags[indexOfTagToBeRemoved])
  })
})
