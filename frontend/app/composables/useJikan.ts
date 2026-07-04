import { ref, watch, isRef, type Ref } from 'vue'

export function useJikan<T = any>(
  path: string | (() => string) | Ref<string>,
  options: { immediate?: boolean; [key: string]: any } = {}
) {
  const data = ref<T | null>(null)
  const pending = ref(false)
  const error = ref<any>(null)

  const getPath = (): string => {
    if (typeof path === 'function') {
      return path()
    }
    if (isRef(path)) {
      return path.value
    }
    return path
  }

  const execute = async () => {
    pending.value = true
    error.value = null
    const currentPath = getPath()
    const url = `http://localhost:8080/api${currentPath}`

    const fetchWithRetry = async (isRetry = false): Promise<T> => {
      try {
        return await $fetch<T>(url, {
          ...options
        })
      } catch (err: any) {
        const status = err?.status
        if (!isRetry && (status === 429 || (status >= 500 && status < 600))) {
          // Wait 2 seconds and retry once
          await new Promise((resolve) => setTimeout(resolve, 2000))
          return await fetchWithRetry(true)
        }
        throw err
      }
    }

    try {
      const res = await fetchWithRetry()
      data.value = res
    } catch (err: any) {
      error.value = err
      data.value = null
    } finally {
      pending.value = false
    }
  }

  const immediate = options.immediate !== false
  if (immediate) {
    execute()
  }

  if (typeof path === 'function') {
    watch(path, () => {
      execute()
    })
  } else if (isRef(path)) {
    watch(path, () => {
      execute()
    })
  }

  return {
    data,
    pending,
    error,
    refresh: execute
  }
}
