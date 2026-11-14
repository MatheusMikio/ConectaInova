import { useEffect, useId, useRef } from 'react'
import type { ReactNode } from 'react'
import { Icon } from './Icon'

export function FormDialog({ open, title, busy, onClose, children }: {
  open: boolean
  title: string
  busy: boolean
  onClose: () => void
  children: ReactNode
}) {
  const dialog = useRef<HTMLDialogElement>(null)
  const titleId = useId()

  useEffect(() => {
    const element = dialog.current
    if (open) {
      element?.showModal()
      const previousOverflow = document.body.style.overflow
      document.body.style.overflow = 'hidden'
      return () => { document.body.style.overflow = previousOverflow }
    }
    element?.close()
  }, [open])

  return <dialog ref={dialog} className="form-dialog" aria-labelledby={titleId}
    onCancel={event => { if (busy) event.preventDefault() }}
    onClose={onClose}
    onClick={event => { if (event.target === dialog.current && !busy) onClose() }}>
    <div className="form-dialog-heading"><h2 id={titleId}>{title}</h2><button type="button" className="icon-button" aria-label="Fechar formulário" disabled={busy} onClick={onClose}><Icon name="close" /></button></div>
    {children}
  </dialog>
}
