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

export interface JikanEpisode {
  mal_id: number
  url?: string
  title: string
  title_japanese?: string | null
  title_romanji?: string | null
  aired?: string | null
  score?: number | null
  filler?: boolean
  recap?: boolean
  forum_url?: string | null
}

export interface JikanPagination {
  last_visible_page: number
  has_next_page: boolean
}

export interface JikanEpisodesResponse {
  data: JikanEpisode[]
  pagination: JikanPagination
}


