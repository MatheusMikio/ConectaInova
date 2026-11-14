import { Icon } from './Icon'
import type { IconName } from './Icon'

export function ScreenHeading({ title, subtitle, icon }: { title: string; subtitle: string; icon: IconName }) {
  return <div className="screen-heading">
    <div className="heading-line"><h1><Icon name={icon} size={27} />{title}</h1></div>
    <p>{subtitle}</p>
  </div>
}
