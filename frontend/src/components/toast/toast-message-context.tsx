import React, { useState } from 'react'

export const ToastMessageContext = React.createContext<{
  isVisible: boolean
  message: string
  type: ToastType

  show: (_: string, __: ToastType) => void
  showAsync: (_: string, __: ToastType) => Promise<any>
  hide: () => void
}>({
  isVisible: false,
  message: '',
  type: 'info',
  show: (_: string, __: ToastType) => ({}),
  showAsync: async (_: string, __: ToastType) => new Promise(() => ({})),
  hide: () => ({}),
})

export type ToastType = null | 'error' | 'success' | 'info'

export const ToastMessageProvider = ({
  children,
}: {
  children: React.ReactElement | React.ReactElement[]
}) => {
  const [isVisible, setIsVisible] = useState(false)
  const [message, setMessage] = useState<string>('')
  const [type, setType] = useState<ToastType>('info')

  const hide = () => {
    setIsVisible(false)
  }

  const show = (message: string, type: ToastType) => {
    setIsVisible(true)
    setMessage(message)
    setType(type)

    setTimeout(hide, 2000)
  }

  const showAsync = async (message: string, type: ToastType) =>
    new Promise((resolve) => {
      setIsVisible(true)
      setMessage(message)
      setType(type)

      setTimeout(resolve, 2000)
    })

  return (
    <ToastMessageContext.Provider
      value={{
        isVisible,
        message,
        type,
        show,
        showAsync,
        hide,
      }}
    >
      {children}
    </ToastMessageContext.Provider>
  )
}

export function useToast() {
  const context = React.useContext(ToastMessageContext)

  if (context === undefined) {
    throw new Error('useToast must be used within an ToastMessageProvider')
  }
  return context
}
