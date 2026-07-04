export interface JikanAnime {
  mal_id: number
  title: string
  title_english?: string | null
  images: {
    jpg: {
      image_url: string
      large_image_url: string
    }
    webp?: {
      image_url: string
      large_image_url: string
    }
  }
  synopsis?: string | null
  score: number | null
  year?: number | null
  episodes?: number | null
  genres?: Array<{ mal_id: number; name: string; type?: string }>
  studios?: Array<{ mal_id: number; name: string }>
  aired?: {
    string?: string
    from?: string
    to?: string
  }
  type?: string | null
}

export type Anime = JikanAnime

