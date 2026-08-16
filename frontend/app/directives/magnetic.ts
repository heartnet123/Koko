/**
 * v-magnetic — buttons pull slightly toward the cursor.
 * Direct DOM manipulation via rAF; no reactive state, no re-renders.
 * Transform transitions handled by CSS on the bound element.
 */
export default {
  mounted(el: HTMLElement, binding: { value?: number }) {
    const strength = binding.value ?? 0.25
    let raf = 0

    const onMove = (e: MouseEvent) => {
      cancelAnimationFrame(raf)
      raf = requestAnimationFrame(() => {
        const r = el.getBoundingClientRect()
        const dx = e.clientX - (r.left + r.width / 2)
        const dy = e.clientY - (r.top + r.height / 2)
        el.style.transform = `translate(${dx * strength}px, ${dy * strength}px)`
      })
    }

    const onLeave = () => {
      cancelAnimationFrame(raf)
      el.style.transform = 'translate(0, 0)'
    }

    el.addEventListener('mousemove', onMove)
    el.addEventListener('mouseleave', onLeave)

    ;(el as HTMLElement & { _magneticCleanup?: () => void })._magneticCleanup = () => {
      cancelAnimationFrame(raf)
      el.removeEventListener('mousemove', onMove)
      el.removeEventListener('mouseleave', onLeave)
    }
  },
  unmounted(el: HTMLElement) {
    ;(el as HTMLElement & { _magneticCleanup?: () => void })._magneticCleanup?.()
  },
}
