import React from 'react'
import Icons from '../icons'
import Animation from '../animations'
import styles from './index.module.css'

interface BaseModalProps {
  className?: string
  onModalClosed?: () => void
}

type ModalComponent<T> = (_: T & BaseModalProps) => React.ReactElement<T & BaseModalProps> | null

type useModalReturn<T, U> = [
  // Modal
  ModalComponent<T>,
  // openModalWith
  (_?: U) => void,
  // closeModal
  () => void,
]

/**
 * @template T - ModalProps
 * @template U - OpenWith
 *
 * @param ModalContent
 *
 * @returns an array with the modal component to render and two methods for opening/closing
 *
 * ModalProps defines the props that we pass to the ModalComponent that's being rendered.
 * OpenWith defines the interface that we can pass to the modal while opening. IOpenWith should be used for the cases where we display the same modal multiple times on a page, for example we can pass the a particular item on a list to the deletion confirmation modal.
 */
function useModal<T, U = unknown>(
  ModalContent: ModalComponent<T & Partial<U> & BaseModalProps>,
): useModalReturn<T, U> {
  const [modalVisible, setModalVisibility] = React.useState<boolean>(false)
  const openWith = React.useRef<U>()

  const openModalWith = (additionalProps?: U) => {
    if (additionalProps) {
      openWith.current = additionalProps
    }
    setModalVisibility(true)
  }

  const closeModal = () => setModalVisibility(false)

  const Modal = (componentProps: T & BaseModalProps) => {
    const handleKeydown = (event: KeyboardEvent) => event.keyCode === 27 && closeModal()

    React.useEffect(() => {
      document.addEventListener('keydown', handleKeydown, false)

      return () => {
        document.removeEventListener('keydown', handleKeydown, false)
      }
    }, [])

    const modalProps = {
      ...componentProps,
      ...(openWith.current ? { ...openWith.current } : null),

      // onModalClosed is passed to the modal content, for being able to close from the modals content itself
      onModalClosed: closeModal,
    }

    return (
      <Animation type='fadeIn' visible={modalVisible}>
        <div className={styles.container}>
          <span className={styles.innerContainer}>
            <div className={styles.overlayBackground} onClick={closeModal} />
            <div className={`${styles.modal} ${componentProps?.className}`}>
              <span onClick={closeModal} className={styles.closeIcon}>
                <Icons.X />
              </span>

              <ModalContent {...modalProps as T & Partial<U> & BaseModalProps} />
            </div>
          </span>
        </div>
      </Animation>
    )
  }

  return [Modal, openModalWith, closeModal]
}

export { styles as modalStyles }

export default useModal
