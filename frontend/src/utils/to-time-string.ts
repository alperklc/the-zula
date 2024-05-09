const toTwoDigitsString = (input: number) => {
  return `${input + 100}`.slice(1)
}

export const toTimeString = (isoDate: string) => {
  const dateOfGivenDay = new Date(isoDate)

  if (Date.now() - dateOfGivenDay.getTime() > 24 * 60 * 60 * 1000) {
    const date = dateOfGivenDay.getDate()
    const month = dateOfGivenDay.getMonth() + 1

    return `${toTwoDigitsString(date)}.${toTwoDigitsString(month)}.${dateOfGivenDay.getFullYear()}`
  }

  const hours = dateOfGivenDay.getHours()
  const minutes = dateOfGivenDay.getMinutes()
  return `${toTwoDigitsString(hours)}:${toTwoDigitsString(minutes)}`
}
