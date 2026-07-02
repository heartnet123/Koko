export function useGenreQuery() {
  const route = useRoute()
  
  const genreIds = computed(() => {
    const genresVal = route.query.genres
    const genreVal = route.query.genre
    const val = genresVal || genreVal
    
    if (!val) return []
    const items = Array.isArray(val) ? val : [val]
    return items
      .filter((item): item is string => typeof item === 'string' && item.trim() !== '')
      .flatMap(item => item.split(','))
      .map(id => id.trim())
      .filter(id => id !== '' && !isNaN(Number(id)))
      .map(Number)
  })
  
  return {
    genreIds
  }
}
