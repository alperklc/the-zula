import React from 'react'
import { render } from '@testing-library/react'

import { ActivityGraph } from './index'

vi.mock('react-intl', () => ({
  useIntl: () => ({ formatMessage: vi.fn() }),
}))

describe('ActivityGraph', () => {
  it('displays a square for each day in a year', () => {
    // arrange
    const date = new Date('2021-09-04')

    // act
    const { queryAllByTestId, queryByTestId } = render(
      <ActivityGraph data={[]} currentDate={date} />,
    )

    // assert
    expect(queryAllByTestId('emptySquare')).toHaveLength(4)
    expect(queryByTestId('squares')?.childNodes).toHaveLength(370)
  })

  it('displays day squares with correct level', () => {
    // arrange
    const date = new Date('2021-09-03')
    const data = [
      { date: '2021-04-07', count: 4 },
      { date: '2021-04-08', count: 1 },
      { date: '2021-04-11', count: 21 },
      { date: '2021-04-12', count: 7 },
      { date: '2021-04-16', count: 3 },
      { date: '2021-04-17', count: 5 },
    ]

    // act
    const { queryAllByTestId, queryByTestId } = render(
      <ActivityGraph data={data} currentDate={date} />,
    )

    // assert
    expect(queryAllByTestId('emptySquare')).toHaveLength(3)
    expect(queryByTestId('daySquare-2021-04-07')?.getAttribute('data-level')).toBe('1')
    expect(queryByTestId('daySquare-2021-04-11')?.getAttribute('data-level')).toBe('4')
    expect(queryByTestId('daySquare-2021-04-12')?.getAttribute('data-level')).toBe('2')
    expect(queryByTestId('daySquare-2021-04-16')?.getAttribute('data-level')).toBe('1')
  })
})
