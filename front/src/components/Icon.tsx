import type { CSSProperties } from 'react'

const paths = {
  home: 'm3 10 9-7 9 7M5 9v11h5v-6h4v6h5V9',
  plus: 'M12 5v14M5 12h14',
  search: 'M21 21l-5-5M18 10a8 8 0 1 1-16 0 8 8 0 0 1 16 0Z',
  target: 'm16 3-4 9 9-4M21 12a9 9 0 1 1-9-9M17 12a5 5 0 1 1-5-5',
  chart: 'M4 20V10h3v10M10 20V5h3v15M16 20V2h3v18M2 21h20',
  settings: 'm9 3-1 3-3 1-2 3 2 2-1 3 3 3 3-1 2 2 3-2 1-3 3-1v-4l-3-1-1-3-3-1-2 1ZM15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z',
  clipboard: 'M9 4H5v18h14V4h-4M9 2h6v5H9zM8 12h1m3 0h4m-8 5h1m3 0h4',
  bell: 'M18 8a6 6 0 0 0-12 0c0 7-3 7-3 9h18c0-2-3-2-3-9M10 21h4',
  chevron: 'm8 10 4 4 4-4',
  cube: 'm12 2 9 5v10l-9 5-9-5V7l9-5ZM3 7l9 5 9-5M12 12v10M7 4.8l9 5.1',
  rocket: 'M14 5c3-3 7-3 7-3s0 4-3 7l-7 7-4-4 7-7ZM7 12H3l3-6h7M11 16v5l6-3v-6M7 17l-4 4M5 15l-3 3M9 19l-3 3M17 6h.01',
  arrowLeft: 'm10 5-7 7 7 7M3 12h18',
  arrowRight: 'm14 5 7 7-7 7M3 12h18',
  check: 'm5 12 4 4L19 6',
  info: 'M12 11v6M12 7h.01M22 12a10 10 0 1 1-20 0 10 10 0 0 1 20 0Z',
  warning: 'm12 3 10 18H2L12 3ZM12 9v5M12 17h.01',
  close: 'm6 6 12 12M6 18 18 6',
  tag: 'M3 3h9l9 9-9 9-9-9V3ZM8 8h.01',
  trophy: 'M8 3h8v7a4 4 0 0 1-8 0V3ZM8 5H4v3a4 4 0 0 0 4 4m8-7h4v3a4 4 0 0 1-4 4M12 14v6M8 21h8',
  list: 'M8 6h13M8 12h13M8 18h13M3 6h.01M3 12h.01M3 18h.01',
  calendar: 'M4 5h16v16H4V5ZM4 10h16M8 2v5M16 2v5',
  user: 'M20 21a8 8 0 0 0-16 0M16 7a4 4 0 1 1-8 0 4 4 0 0 1 8 0Z',
} as const

export type IconName = keyof typeof paths

export function Icon({ name, size = 20, className, style }: {
  name: IconName
  size?: number
  className?: string
  style?: CSSProperties
}) {
  return <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.7" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true" className={className} style={style}><path d={paths[name]} /></svg>
}

export function Logo() {
  return <span className="logo">
    <svg width="36" height="27" viewBox="0 0 42 30" fill="none" aria-hidden="true">
      <path d="m24 9-7 10a8 8 0 1 1-12-9l5-5a8 8 0 0 1 12 1" stroke="#0062f6" strokeWidth="4" strokeLinecap="round" />
      <path d="m18 21 7-10a8 8 0 1 1 12 9l-5 5a8 8 0 0 1-12-1" stroke="#0037ad" strokeWidth="4" strokeLinecap="round" />
    </svg>
    <span>Conecta<span>Inova</span></span>
  </span>
}
