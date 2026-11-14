import type { ReactNode } from 'react'

export function Field({ label, children, required = false }: { label: string; children: ReactNode; required?: boolean }) {
  return <label className="field"><span>{label}{required && <span className="required" aria-hidden="true"> *</span>}</span>{children}</label>
}
