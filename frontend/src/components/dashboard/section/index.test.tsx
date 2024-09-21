import { describe, it, expect } from 'vitest'
import { BrowserRouter } from 'react-router-dom'
import { render } from '@testing-library/react'
import { AllContentSection, DashboardSection } from './index'

describe('DashboardSection', () => {
  it('WHEN an array of content is given THEN they are rendered correctly', () => {
    // arrange
    const rows = [
      {
          typename: 'NOTE',
          id: 'aabel2cGoibqDbkDIJf78',
          title: 'test note',
        },
        {  
          typename: 'BOOKMARK',
          id: 'AJ0K4sAPIqrG84AvDLOSL',
          name: 'kek.jpg',
          title: 'kek.jpg',
        },
        {
          typename: 'BOOKMARK',
          id: '1FO5kHFBYZiFdxsdeblYu',
          title: 'Documentation - The Go Programming Language',
          name: 'Documentation - The Go Programming Language',
      },
    ]

    // act
    const { queryByTestId, queryAllByTestId } = render(
      <BrowserRouter><DashboardSection title='test' rows={rows} /></BrowserRouter>
    )
    const renderedContent = queryAllByTestId('content')

    // assert
    expect(queryByTestId('title')?.textContent).toBe('test')
    expect(renderedContent.length).toBe(rows.length)
    expect(renderedContent[0].textContent).toBe(rows[0].title)
    expect(renderedContent[1].textContent).toBe(rows[1].name)
    expect(renderedContent[2].textContent).toBe(rows[2].title)
  })
})

describe('AllContentSection', () => {
  it('WHEN all given numbers are 0 THEN the numbers are not displayed', () => {
    // act
    const { queryByTestId } = render(
      <BrowserRouter><AllContentSection data={{ numberOfBookmarks: 0, numberOfNotes: 0 }} /></BrowserRouter>,
    )
    const [numberOfNotes, numberOfBookmarks] = [
      'numberOfNotes',
      'numberOfBookmarks',
    ].map((testId) => queryByTestId(testId))

    // assert
    expect(numberOfNotes?.textContent).toBe('')
    expect(numberOfBookmarks?.textContent).toBe('')
  })

  it('WHEN some of given numbers are NOT 0 THEN they are displayed', () => {
    // act
    const { queryByTestId } = render(
      <BrowserRouter><AllContentSection data={{ numberOfBookmarks: 5, numberOfNotes: 4 }} /></BrowserRouter>,
    )
    const [numberOfNotes, numberOfBookmarks] = [
      'numberOfNotes',
      'numberOfBookmarks',
    ].map((testId) => queryByTestId(testId))

    // assert
    expect(numberOfNotes?.textContent).toBe(' (4)')
    expect(numberOfBookmarks?.textContent).toBe(' (5)')
  })
})
